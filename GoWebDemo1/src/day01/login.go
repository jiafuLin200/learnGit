package main

import (
	"net/http"
	"fmt"
	"log"
	"html/template"
	"strings"
)

func helloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析参数,默认是不会解析的
	fmt.Println("body=", r.Body)
	fmt.Println("Form=", r.Form) // 这些信息是输出到服务器端的打印信息
	fmt.Println("path=", r.URL.Path)
	fmt.Println("scheme=", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, GetCurrPath()) // 输出当前的工作路径
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./login.gtpl")
		fmt.Println("--------------------------------------")
		//fmt.Println(err.Error())
		fmt.Println("--------------------------------------")
		t.Execute(w, nil)
	} else {
		// 请求的是登陆数据,那么执行登陆的逻辑判断
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
func main() {
	http.HandleFunc("/", helloName) // 设置访问的路由
	http.HandleFunc("/login", login) // 设置访问的路由
	http.HandleFunc("/upload", Upload)
	err := http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
