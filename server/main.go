package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/valyala/fasthttp"
)

func main() {
	serverCrt := "../datalistener/https/mtls/server-certificate.pem"
	serverKey := "../datalistener/https/mtls/server-private-key.pem"
	caCrt := "../datalistener/https/mtls/ca-certificate.pem"
	addr := ":10443"

	serverCert, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		fmt.Printf("Error loading server certificate and key: %s\n", err)
		return
	}

	caCert, err := os.ReadFile(caCrt)
	if err != nil {
		fmt.Printf("Error loading CA certificate: %s\n", err)
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	server := &fasthttp.Server{
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		clientIP := r.RemoteAddr
		fmt.Printf("Received POST request from client IP: %s\n", clientIP)
		fmt.Printf("Received Body: %s\n", string(body))

		fmt.Fprintf(w, "Hello, client! You have successfully connected to the server over mTLS.")
	})

	fmt.Printf("Server is listening on %s\n", addr)
	err = server.ListenAndServeTLS(addr, "", "")
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
