package buffer

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"datalistener/src/logging"
	"datalistener/src/models"
	"fmt"
	"os"

	"github.com/valyala/fasthttp"
)

type HttpsMtlsBufferConfig struct {
	Protocol                 string
	Address                  string
	Port                     int
	Endpoint                 string
	ContentType              string
	ItemSeparator            string
	ClientKeyFile            string
	ClientCertFile           string
	CACertFile               string
	SkipHostNameVerification bool
}

func (cfg HttpsMtlsBufferConfig) Notify(entries []models.EntryData, convertToJSONL bool) error {
	logging.Logger.Debug("HTTPS MTLS Buffering")

	if len(entries) == 0 {
		return nil
	}

	var combinedBody bytes.Buffer

	for _, entry := range entries {
		combinedBody.Write(entry.Body)
		combinedBody.WriteString(cfg.ItemSeparator)
	}

	url := fmt.Sprintf("%s://%s:%d%s", cfg.Protocol, cfg.Address, cfg.Port, cfg.Endpoint)

	clientCert, err := tls.LoadX509KeyPair(cfg.ClientCertFile, cfg.ClientKeyFile)
	if err != nil {
		return err
	}

	caCert, err := os.ReadFile(cfg.CACertFile)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: cfg.SkipHostNameVerification,
	}

	client := &fasthttp.Client{
		TLSConfig: tlsConfig,
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType(cfg.ContentType)
	req.SetBody(combinedBody.Bytes())

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = client.Do(req, resp)
	if err != nil {
		return err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("HTTPS request failed with status: %d", resp.StatusCode())
	}

	return nil
}

func (cfg HttpsMtlsBufferConfig) String() string {
	return fmt.Sprintf("%s://%s:%d%s - %s", cfg.Protocol, cfg.Address, cfg.Port, cfg.Endpoint, cfg.ContentType)
}
