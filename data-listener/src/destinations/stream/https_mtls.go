package stream

import (
	"crypto/tls"
	"crypto/x509"
	"datalistener/src/models"
	"fmt"
	"os"

	"github.com/valyala/fasthttp"
)

type HttpsMtlsStreamConfig struct {
	Protocol                 string
	Address                  string
	Port                     int
	EndPoint                 string
	ContentType              string
	ClientCertFile           string
	ClientKeyFile            string
	CACertFile               string
	SkipHostNameVerification bool
}

func (cfg HttpsMtlsStreamConfig) Notify(entry models.EntryData) error {

	url := fmt.Sprintf("%s://%s:%d%s", cfg.Protocol, cfg.Address, cfg.Port, cfg.EndPoint)

	// Load the client certificate and private key
	clientCert, err := tls.LoadX509KeyPair(cfg.ClientCertFile, cfg.ClientKeyFile)
	if err != nil {
		return err
	}

	// Load the CA certificate for server verification
	caCert, err := os.ReadFile(cfg.CACertFile)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a TLS configuration with mTLS settings
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: cfg.SkipHostNameVerification, // Skip hostname verification
	}

	// Create an HTTP client with the TLS configuration
	client := &fasthttp.Client{
		TLSConfig: tlsConfig,
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType(cfg.ContentType)
	req.SetBody(entry.Body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = client.Do(req, resp)
	if err != nil {
		return err
	}
	body := resp.Body()
	fmt.Printf("Response from server: %s\n", body)

	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("HTTPS request failed with status: %d", resp.StatusCode())
	}

	return nil
}

func (cfg HttpsMtlsStreamConfig) String() string {
	return fmt.Sprintf("%s://%s:%d%s - %s", cfg.Protocol, cfg.Address, cfg.Port, cfg.EndPoint, cfg.ContentType)
}
