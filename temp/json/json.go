package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// import (
// 	"encoding/json"
// 	"fmt"
// )

// type ST struct {
// 	Foo string `json:"foo"`
// 	NST struct {
// 		Foo4 string `json:"foo4"`
// 	} `json:"foo3"`
// }

// func main() {
// 	// var body map[string]string
// 	dat := `
// 	{
// 		"foo": "bar",
// 		"foo2": "bar2",
// 		"foo3": {
// 			"foo4": "foo5"
// 		}
// 	}
// 	`
// 	var st ST
// 	err := json.Unmarshal([]byte(dat), &st)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(st)
// }

func JsonGet(obj map[string]interface{}, path string) (interface{}, bool) {
	spt := strings.Split(path, ".")
	temp := obj
	n := len(spt)
	for _, v := range spt[:n-1] {
		switch temp[v].(type) {
		case map[string]interface{}:
			temp = temp[v].(map[string]interface{})
		default:
			return nil, false
		}
	}
	v, ok := temp[spt[n-1]]
	return v, ok
}

func main() {
	var result map[string]interface{}
	dat := `
	{
		"foo": "bar",
		"foo2": "bar2",
		"foo3": {
			"foo4": "foo5",
			"foo5": {
				"foo6": "foo7"
			}
		}
	}
	`
	json.Unmarshal([]byte(dat), &result)
	fmt.Println(JsonGet(result, "foo3.foo5"))
	fmt.Println(JsonGet(result, "foo3.foo5.foo7"))
	fmt.Println(JsonGet(result, "foo3.foo5.foo8"))
	fmt.Println(JsonGet(result, "foo3.foo5.foo6"))
	fmt.Println(JsonGet(result, "foo3.foo5.foo6.foo8"))
	fmt.Println(JsonGet(result, "foo3.foo5.foo6.foo8.foo9"))
}
