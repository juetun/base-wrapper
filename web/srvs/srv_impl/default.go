/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:13 下午
 */

// @Copyright (c) 2020.
// @Author ${USER}
// @Date ${DATE}
package srv_impl

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/juetun/base-wrapper/lib/utils/identifying_code_pkg"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/juetun/base-wrapper/lib/app/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/web/daos/dao_impl"
	"github.com/juetun/base-wrapper/web/srvs"
	"github.com/juetun/base-wrapper/web/wrapper"
)

type ServiceDefaultImpl struct {
	base.ServiceBase
}

func NewServiceDefaultImpl(context ...*base.Context) (res srvs.ServiceDefault) {
	p := &ServiceDefaultImpl{}
	p.SetContext(context...)
	res = p
	return
}

//验证码生成
func (r *ServiceDefaultImpl) Auth(arg *wrapper.ArgumentDefault) (result interface{}, err error) {
	type Result struct {
		IdKey         string `json:"id_key"`
		Base64Captcha string `json:"base_64_captcha"`
	}
	var res Result
	// 生成验证码逻辑
	res.IdKey, res.Base64Captcha, err = identifying_code_pkg.NewIdentifyingCode(
		identifying_code_pkg.Context(&identifying_code_pkg.CustomizeRdsStore{
			Context: r.Context,
		})).CreateAndGetImgBase64Message()
	result = res
	return
}

//验证码校验逻辑
func (r *ServiceDefaultImpl) AuthRes(arg *wrapper.ArgumentDefault) (result interface{}, err error) {
	// 校验逻辑
	result = identifying_code_pkg.NewIdentifyingCode(identifying_code_pkg.Context(&identifying_code_pkg.CustomizeRdsStore{
		// 参数
	})).Context.Verify(arg.IdKey, "anwser", true)

	return
}
func (r *ServiceDefaultImpl) Index(arg *wrapper.ArgumentDefault) (res *wrapper.ResultDefault, err error) {
	res = &wrapper.ResultDefault{}
	res.Users, err = dao_impl.NewDaoUserImpl(r.Context).
		GetUser(arg)
	return
}
func (r *ServiceDefaultImpl) TestEs(arg *wrapper.ArgumentDefault) (result interface{}, err error) {
	log.SetFlags(0)

	var (
		rObj map[string]interface{}
		wg   sync.WaitGroup
	)

	// Initialize a client with the default settings.
	//
	// An `ELASTICSEARCH_URL` environment variable will be used when exported.
	//
	es := app_obj.GetElasticSearchMaps()
	opEs := common.NewOperateElasticSearch(es)
	esInfo, err := opEs.GetClusterInfo()
	log.Printf("%#v", esInfo)

	var data = struct {
		Title string `json:"title"`
	}{
		Title: "Test One",
	}
	rByte, err := json.Marshal(data)

	// 添加数据
	opEs.Write(rByte)

	// 2. Index documents concurrently
	//
	for i, title := range []string{"Test One", "Test Two"} {
		wg.Add(1)

		go func(i int, title string) {
			defer wg.Done()

			// Build the request body.
			var b strings.Builder
			b.WriteString(`{"title" : "`)
			b.WriteString(title)
			b.WriteString(`"}`)

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "test",
				DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(b.String()),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, _ := req.Do(context.Background(), es)
			// if err != nil {
			//	log.Fatalf("Error getting response: %s", err)
			// }
			defer res.Body.Close()
			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(i, title)
	}
	wg.Wait()

	log.Println(strings.Repeat("-", 37))

	// 3. Search for the indexed documents
	//
	// Build the request body.
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "test",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&rObj); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(rObj["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(rObj["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range rObj["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))
	return
}
func (r *ServiceDefaultImpl) Tmain(arg *wrapper.ArgumentDefault) (result interface{}, err error) {

	return
}
