package handlers

import (
	"datalistener/src/metrics"
	"datalistener/src/models"
	"encoding/json"
	"net/url"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func RequestHandler(ctx *fasthttp.RequestCtx, bufferChannel chan models.EntryData, streamerChannel chan models.EntryData, logger *zap.Logger, validateJSON *bool) {

	if !ctx.IsPost() {
		ctx.Error("Only POST requests are allowed", fasthttp.StatusMethodNotAllowed)
		return
	}
	body := ctx.PostBody()

	// Check if the body is a valid JSON
	if *validateJSON {
		var jsonCheck map[string]interface{}
		if err := json.Unmarshal(body, &jsonCheck); err != nil {
			ctx.Error("Invalid JSON", fasthttp.StatusBadRequest)
			return
		}
	}

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

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Success")
	metrics.RecordRequestMetrics(string(ctx.Method()), "RequestHandler", len(requestData.Body))
	logger.Debug("Request data handled successfully")
}
