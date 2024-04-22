package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"github.com/valyala/fasthttp"
)

func main() {
	clientCrt := "../datalistener/https/mtls/client-certificate.pem"
	clientKey := "../datalistener/https/mtls/client-private-key.pem"
	caCrt := "../datalistener/https/mtls/ca-certificate.pem"
	addr := "https://localhost:10443"

	clientCert, err := tls.LoadX509KeyPair(clientCrt, clientKey)
	if err != nil {
		log.Fatalf("Error loading client certificate and key: %s\n", err)
	}

	// Load the CA certificate for server verification
	caCert, err := os.ReadFile(caCrt)
	if err != nil {
		log.Fatalf("Error loading CA certificate: %s\n", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a TLS configuration with mTLS settings
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, // Skip hostname verification
	}

	// Create an HTTP client with the TLS configuration
	client := &fasthttp.Client{
		TLSConfig: tlsConfig,
	}

	// Send a GET request to the server
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	req.SetRequestURI(addr) // Replace with the correct server URL
	req.Header.SetContentType("application/json")

	// Set the request body to JSON content
	jsonStr := `{"key1": "value1", "key2": "value2"}` // Replace with your JSON content
	req.SetBodyString(jsonStr)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(addr) // Replace with the correct server URL
	err = client.Do(req, resp)
	if err != nil {
		log.Fatalf("Error sending request: %s\n", err)
		return
	}

	// Read and print the response
	body := resp.Body()
	fmt.Printf("Response from server: %s\n", body)
}
