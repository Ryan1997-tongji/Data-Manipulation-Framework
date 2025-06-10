// package utils
// @author: chenzhewei.97
// @create date: 2025/6/9
package utils

import "encoding/json"

func SimpleJson(obj interface{}) string {
	bs, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bs)
}