package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	fmt.Println(args)
	fmt.Println("args")
	if len(args) <= 1 {
		panic("请输入端口号")
	}
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(args)
		w.Write([]byte(args[1]))
	})
	http.ListenAndServe(":"+args[1], nil)
}
