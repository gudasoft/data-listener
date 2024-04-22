package stream

import (
	"datalistener/src/logging"
	"datalistener/src/models"
	"fmt"

	"github.com/valyala/fasthttp"
)

type HttpStreamConfig struct {
	Protocol    string
	Address     string
	Port        int
	Endpoint    string
	ContentType string
}

func (cfg HttpStreamConfig) Notify(data models.EntryData) error {
	logging.Logger.Debug("HTTP Streaming with fasthttp")

	url := fmt.Sprintf("%s://%s:%d%s", cfg.Protocol, cfg.Address, cfg.Port, cfg.Endpoint)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType(cfg.ContentType)
	req.SetBody(data.Body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fasthttp.Do(req, resp)
	if err != nil {
		logging.Logger.Sugar().Errorf("Error making HTTP request:", err)
		return err
	}

	statusCode := resp.StatusCode()
	if statusCode >= 200 && statusCode < 300 {
		logging.Logger.Debug("Data sent successfully over HTTP with fasthttp")
	} else {
		logging.Logger.Sugar().Errorf("HTTP request failed with status code:", statusCode)
		return fmt.Errorf("HTTP request failed with status code: %d", statusCode)
	}

	return nil
}

func (cfg HttpStreamConfig) String() string {
	return fmt.Sprintf("%s://%s:%d%s - %s", cfg.Protocol, cfg.Address, cfg.Port, cfg.Endpoint, cfg.ContentType)
}
