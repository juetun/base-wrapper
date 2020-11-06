/**
* @Author:changjiang
* @Description:
* @File:default
* @Version: 1.0.0
* @Date 2020/8/18 6:13 下午
 */
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/juetun/base-wrapper/lib/app_obj"
	"github.com/juetun/base-wrapper/lib/base"
	"github.com/juetun/base-wrapper/lib/common"
	"github.com/juetun/base-wrapper/web/daos"
	"github.com/juetun/base-wrapper/web/pojos"
	"log"
	"strconv"
	"strings"
	"sync"
)

type ServiceDefault struct {
	base.ServiceBase
}

func NewServiceDefault(context ...*base.Context) (p *ServiceDefault) {
	p = &ServiceDefault{}
	p.SetContext(context...)
	return
}
func (r *ServiceDefault) Index(arg *pojos.ArgumentDefault) (res *pojos.ResultDefault, err error) {
	res = &pojos.ResultDefault{}
	dao := daos.NewDaoUser(r.Context)
	res.Users, err = dao.GetUser(arg)
	return
}
func (r *ServiceDefault) TestEs(arg *pojos.ArgumentDefault) (result interface{}, err error) {
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

	//添加数据
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
			//if err != nil {
			//	log.Fatalf("Error getting response: %s", err)
			//}
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
func (r *ServiceDefault) Tmain(arg *pojos.ArgumentDefault) (result interface{}, err error) {

	return
}