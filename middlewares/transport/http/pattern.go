package http

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func marshalMap(pattern string, prefix string, param map[string]interface{}) string {

	for k, v := range param {

		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Map {
			if prefix == "" {
				pattern = marshalMap(pattern, k, v.(map[string]interface{}))
			} else {
				pattern = marshalMap(pattern, prefix+"."+k, v.(map[string]interface{}))
			}
			continue
		} else if t.Kind() == reflect.Array {
			continue
		}

		key := ":" + k
		val := fmt.Sprintf("%v", v)

		/// JSON UnMarshal would convert int/int32/int64 -> float64
		if reflect.TypeOf(v).Kind() == reflect.Float64 {
			val = fmt.Sprintf("%v", int64(v.(float64)))
		}

		if "" != prefix {
			key = fmt.Sprintf(":%v.%v", prefix, k)
		}
		pattern = strings.Replace(pattern, key, val, -1)
	}
	return pattern
}

//// marshal requet to URI by pattern
func MarshalPattern(pattern string, request interface{}) string {

	if nil == request {
		return pattern
	}
	mp := make(map[string]interface{})
	js, _ := json.Marshal(request)
	json.Unmarshal(js, &mp)
	return marshalMap(pattern, "", mp)
}
