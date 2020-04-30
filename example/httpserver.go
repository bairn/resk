package main

import "net/http"

func main() {
	http.HandleFunc("/hello",
		func(w http.ResponseWriter, r *http.Request) {
		   s := r.URL.RawQuery
		   w.Write([]byte("hello,world." + s))
	    })

	http.ListenAndServe(":8002", nil)
}