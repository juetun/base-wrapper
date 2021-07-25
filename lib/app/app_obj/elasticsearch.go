// Package app_obj
package app_obj

import (
	"github.com/elastic/go-elasticsearch/v7"
)

var ElasticSearchV7Maps = make(map[string]*elasticsearch.Client)

// GetElasticSearchMaps 获取ElasticSearchMaps操作实例
func GetElasticSearchMaps(nameSpace ...string) (res *elasticsearch.Client, keyName string) {

	switch l := len(nameSpace); l {
	case 0:
		keyName = "default"
	case 1:
		keyName = nameSpace[0]
	default:
	}
	if _, ok := ElasticSearchV7Maps[keyName]; ok {
		res = ElasticSearchV7Maps[keyName]
		return
	}
	return
}
