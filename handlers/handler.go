package handlers

import (
	"buffer-handler/logging"
	"buffer-handler/metrics"
	"buffer-handler/models"
	"net/url"

	"github.com/valyala/fasthttp"
)

func RequestHandler(ctx *fasthttp.RequestCtx, bufferChannel chan models.EntryData, streamerChannel chan models.EntryData) {
	if !ctx.IsPost() {
		ctx.Error("Only POST requests are allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	body := ctx.PostBody()

	headers := make(map[string][]string)
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = append(headers[string(key)], string(value))
	})

	parsedUrl, err := url.Parse(ctx.URI().String())
	if err != nil {
		ctx.Error("Failed to parse request URL", fasthttp.StatusInternalServerError)
		return
	}

	requestData := models.EntryData{
		Headers: headers,
		Body:    body,
		Url:     *parsedUrl,
	}

	if streamerChannel != nil {
		streamerChannel <- requestData
	}
	if bufferChannel != nil {
		bufferChannel <- requestData
	}

	// logging.Logger.Sugar().Debugf("Received POST request with body: %s\n", string(body))
	// logging.Logger.Sugar().Debugf("Headers: %+v\n", headers)
	// logging.Logger.Sugar().Debugf("URL: %s\n", r.URL.String())

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Request data processed successfully")
	metrics.RecordRequestMetrics(string(ctx.Method()), "RequestHandler", len(requestData.Body))
	logging.Logger.Debug("Request data processed successfully")
}
