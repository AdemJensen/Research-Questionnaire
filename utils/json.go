package utils

import "encoding/json"

func Json(obj interface{}) string {
	v, _ := json.Marshal(obj)
	return string(v)
}
