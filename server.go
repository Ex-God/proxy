package main

import (
	"fmt"
	"bytes"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

//TODO
//fix header issue
//add auth
//add https

type Request struct {
	Method string `json:"method"`
	URL *url.URL `json:"url"`
	Host string `json:"host"`
	Header http.Header `json:"header"`
	Body string `json:"body"`
}

func doRequest(r Request) *http.Response {
	if r.URL.Scheme == "" {
		r.URL.Scheme = "http"
	}
	address := fmt.Sprintf("%s://%s%s", r.URL.Scheme, r.Host, r.URL.Path)
	if r.URL.RawQuery != "" {
		address = fmt.Sprintf("%s://%s%s?%s", r.URL.Scheme, r.Host, r.URL.Path, r.URL.RawQuery)
	}
	if r.Host == "" {
		r.Host = "127.0.0.1"
	}
	if r.Host == "127.0.0.1:9000" {
		return nil
	}
	req, err := http.NewRequest(r.Method, address, bytes.NewBuffer([]byte(r.Body)))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for k, v := range r.Header {
		req.Header.Add(k, strings.Join(v, " "))
	}
	fmt.Printf("%+v\n", req.Header)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	req := Request{}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(buf, &req)
	if err != nil {
		fmt.Println(err)
	}
	res := doRequest(req)
	page, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range res.Header {
		w.Header().Add(k, strings.Join(v, " "))
	}
	//fmt.Println(string(page))
	w.Write([]byte(page))
}

func main() {
	http.HandleFunc("/", serverHandler)
	http.ListenAndServe("127.0.0.1:9001", nil)
}

