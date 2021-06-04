// Package app_obj
package app_obj

import (
	"github.com/elastic/go-elasticsearch/v7"
)

var ElasticSearchV7Maps = make(map[string]*elasticsearch.Client)

// 获取ElasticSearchMaps操作实例
func GetElasticSearchMaps(nameSpace ...string) (res *elasticsearch.Client) {

	var s string
	switch len := len(nameSpace); len {
	case 0:
		s = "default"
	case 1:
		s = nameSpace[0]
	default:
	}
	if _, ok := ElasticSearchV7Maps[s]; ok {
		res = ElasticSearchV7Maps[s]
		return
	}
	return
}
