package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	m := map[string]interface{}{
		"time":  "time",
		"type":  "payloadType",
		"reqID": "reqID",
		//"all":   string(buf),
	}
	//if len(payloadArr) > 1 {
	//	m["body"] = string(payloadArr[1])
	//}
	j, err := json.Marshal(m)
	if err!=nil{
		fmt.Println(err)
	}
	resp, _ := http.Post("http://localhost:3000/write", "application/json",
		bytes.NewReader(j))
	fmt.Println(resp)
}
