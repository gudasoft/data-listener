package handlers

import (
	"datalistener/src/models"
	"testing"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func TestRequestHandlerValidJSONWithMockedLogger(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/somepath")
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetBody([]byte(`{"key": "value"}`))

	var bufferChannel chan models.EntryData
	var streamerChannel chan models.EntryData
	logger := zap.NewNop()
	validateJSON := true

	RequestHandler(ctx, bufferChannel, streamerChannel, logger, &validateJSON)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Errorf("Expected status code 200, got %d", ctx.Response.StatusCode())
	}

	expectedResponseBody := "Success"
	if string(ctx.Response.Body()) != expectedResponseBody {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponseBody, string(ctx.Response.Body()))
	}
}

func TestRequestHandlerInvalidJSON(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/somepath")
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("application/json")
	ctx.Request.SetBody([]byte(`{"key": "value",}`))

	var bufferChannel chan models.EntryData
	var streamerChannel chan models.EntryData
	logger := zap.NewNop()
	validateJSON := true

	RequestHandler(ctx, bufferChannel, streamerChannel, logger, &validateJSON)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", ctx.Response.StatusCode())
	}

	expectedResponseBody := "Invalid JSON"
	if string(ctx.Response.Body()) != expectedResponseBody {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponseBody, string(ctx.Response.Body()))
	}
}

func TestRequestHandlerWithBinaryData(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/somepath")
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("application/octet-stream")
	binaryData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	ctx.Request.SetBody(binaryData)

	var bufferChannel chan models.EntryData
	var streamerChannel chan models.EntryData
	logger := zap.NewNop()
	validateJSON := true

	RequestHandler(ctx, bufferChannel, streamerChannel, logger, &validateJSON)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", ctx.Response.StatusCode())
	}

	expectedResponseBody := "Invalid JSON"
	if string(ctx.Response.Body()) != expectedResponseBody {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponseBody, string(ctx.Response.Body()))
	}
}

func TestRequestHandlerWithStringBody(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/somepath")
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.Header.SetContentType("text/plain")
	ctx.Request.SetBody([]byte("This is a string"))

	var bufferChannel chan models.EntryData
	var streamerChannel chan models.EntryData
	logger := zap.NewNop()
	validateJSON := true

	RequestHandler(ctx, bufferChannel, streamerChannel, logger, &validateJSON)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", ctx.Response.StatusCode())
	}

	expectedResponseBody := "Invalid JSON"
	if string(ctx.Response.Body()) != expectedResponseBody {
		t.Errorf("Expected response body '%s', got '%s'", expectedResponseBody, string(ctx.Response.Body()))
	}
}

// TODO Test ConnectionClosure
// Test KeepAlive
// Test Timeout
