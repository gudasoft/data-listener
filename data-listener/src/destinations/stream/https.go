package stream

import (
	"datalistener/src/logging"
	"datalistener/src/models"
	"fmt"

	"github.com/valyala/fasthttp"
)

type HttpsStreamConfig struct {
	Protocol    string
	Address     string
	Port        int
	Endpoint    string
	ContentType string
}

func (cfg HttpsStreamConfig) Notify(entry models.EntryData) error {
	logging.Logger.Debug("Https Streaming")

	url := fmt.Sprintf("%s://%s:%d%s", cfg.Protocol, cfg.Address, cfg.Port, cfg.Endpoint)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType(cfg.ContentType)
	req.SetBody(entry.Body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fasthttp.Do(req, resp)
	if err != nil {
		logging.Logger.Error(err.Error())
		return err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		logging.Logger.Sugar().Errorf("HTTPS request failed with status: %d", resp.StatusCode())
		return fmt.Errorf("HTTPS request failed with status: %d", resp.StatusCode())
	}

	return nil
}

func (cfg HttpsStreamConfig) String() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Address, cfg.Port)
}
