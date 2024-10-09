package main

import(
	"fmt"
	"bytes"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const(
	P_HOST = "127.0.0.1"
	P_PORT = "9001"
	L_HOST = "127.0.0.1"
	L_PORT = "9000"
	L_SSL_PORT = "10000"
)

type Request struct {
	Method string `json:"method"`
	URL *url.URL `json:"url"`
	Host string `json:"host"`
	Header http.Header `json:"header"`
	Body string `json:"body"`
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	nr := Request{}
	nr.Method = r.Method
	nr.URL = r.URL
	nr.Header = r.Header
	nr.Host = r.Host
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	nr.Body = fmt.Sprintf("%s", body)
	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(nr)
	if err != nil {
		fmt.Println(err)
	}
	res, err := http.Post(fmt.Sprintf("http://%s:%s", P_HOST, P_PORT), "application/json", &buf)
	if err != nil {
		fmt.Println(err)
	}
	page, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", res.Header)
	for k, v := range res.Header {
		w.Header().Add(k, strings.Join(v, " "))
	}
	w.Write([]byte(page))
}

func proxyHandlerTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok\n"))
}

func main() {
	http.HandleFunc("/", proxyHandlerTest)
	//http.ListenAndServe(fmt.Sprintf("%s:%v", L_HOST, L_PORT), nil)
	http.ListenAndServeTLS(fmt.Sprintf("%s:%v", L_HOST, L_SSL_PORT), "certs/proxy.crt", "certs/proxy.key", nil)
}
