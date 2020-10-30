package common

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

type operateElasticSearch struct {
	es *elasticsearch.Client
}

func NewOperateElasticSearch(es *elasticsearch.Client) (res *operateElasticSearch) {
	return &operateElasticSearch{
		es: es,
	}
}

type EsInfo struct {
	ClientVersion string                 `json:"client_version"`
	Result        map[string]interface{} `json:"result"`
}

func (r *operateElasticSearch)Write(jsonByte []byte)  {
	
}
//获取当前服务器和SDK的信息
func (r *operateElasticSearch) GetClusterInfo() (result EsInfo, err error) {
	var rObj map[string]interface{}
	result = EsInfo{
		ClientVersion: elasticsearch.Version,
	}
	res, err := r.es.Info()
	if err != nil {
		return
	}
	// Check response status
	if res.IsError() {
		err = fmt.Errorf(res.String())
		return
	}

	if err = json.NewDecoder(res.Body).Decode(&rObj); err != nil {
		err = fmt.Errorf("Error parsing the response body: %s ", err.Error())
		return
	}
	result.Result = rObj
	//log.Printf("Client: %s", elasticsearch.Version)
	//log.Printf("Server: %s", rObj["version"].(map[string]interface{})["number"])
	//log.Println(strings.Repeat("~", 37))

	return
}
