package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	fmt.Println(args)
	if len(args) <= 1 {
		panic("请输入端口号")
	}
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(args[1]))
	})
	http.ListenAndServe("localhost:"+args[1], nil)
}