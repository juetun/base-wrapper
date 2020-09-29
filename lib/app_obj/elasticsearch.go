package app_obj

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

var ElasticSearchV7Maps = make(map[string]*elasticsearch.Client)

// 获取ElasticSearchMaps操作实例
func GetElasticSearchMaps(nameSpace ...string) *elasticsearch.Client {

	var s string
	switch len := len(nameSpace); len {
	case 0:
		s = "default"
	case 1:
		s = nameSpace[0]
	default:
		panic("nameSpace receive at most one parameter")
	}
	if _, ok := ElasticSearchV7Maps[s]; ok {
		return ElasticSearchV7Maps[s]
	}
	panic(fmt.Sprintf("the ElasticSearchMaps connect(%s) is not exist", s))
}
