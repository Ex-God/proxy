package main

import (
    "crypto/tls"
    "io"
    "net"
    "net/http"
    "time"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	// Establish a connection to the target server
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	// Hijack the connection to allow for bidirectional communication
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Start copying data between the client and destination connection
	go io.Copy(destConn, clientConn)
	go io.Copy(clientConn, destConn)
}

func main() {
	server := &http.Server{
		Addr: "127.0.0.1:10000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				http.Error(w, "Only CONNECT method is supported", http.StatusMethodNotAllowed)
			}
		}),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	server.ListenAndServeTLS("certs/proxy.crt", "certs/proxy.key")
}
