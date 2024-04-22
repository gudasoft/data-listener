package buffer

import (
	"bytes"
	"datalistener/src/logging"
	"datalistener/src/models"
	"fmt"

	"github.com/valyala/fasthttp"
)

type HttpsBufferConfig struct {
	Protocol      string
	Address       string
	Port          int
	Endpoint      string
	ContentType   string
	ItemSeparator string
}

func (cfg HttpsBufferConfig) Notify(entries []models.EntryData, convertToJSONL bool) error {
	logging.Logger.Debug("HTTPS Buffering")

	if len(entries) == 0 {
		return nil
	}

	var combinedBody bytes.Buffer

	for _, entry := range entries {
		combinedBody.Write(entry.Body)
		combinedBody.WriteString(cfg.ItemSeparator)
	}

	url := fmt.Sprintf("%s://%s:%d%s", cfg.Protocol, cfg.Address, cfg.Port, cfg.Endpoint)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType(cfg.ContentType)
	req.SetBody(combinedBody.Bytes())

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

func (cfg HttpsBufferConfig) String() string {
	return fmt.Sprintf("%s://%s:%d%s - %s", cfg.Protocol, cfg.Address, cfg.Port, cfg.Endpoint, cfg.ContentType)
}
