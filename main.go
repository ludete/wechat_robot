package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ludete/wechat_robot/app"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	bz := make([]byte, 1024*10)
	r.Body.Read(bz)
	fmt.Printf("url : %s, header : %v\n ", r.URL, r.Header)
	fmt.Printf("body : %s\n", bz)
	fmt.Fprintf(w, "nihao")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
	app := app.NewRobotApp("")
	_ = app

	route := mux.NewRouter()
	_ = route
}
