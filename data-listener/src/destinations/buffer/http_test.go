package buffer

import (
	"bytes"
	"datalistener/src/models"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupMockServer(t *testing.T) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}
		defer r.Body.Close()

		assert.Equal(t, "POST", r.Method)

		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		if len(body) != 0 {
			assert.Contains(t, string(body), "pippo", "body should contain pippo")
		}

		if bytes.Contains(body, []byte("error")) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
	return httptest.NewServer(handler)
}

func TestNotify(t *testing.T) {
	mockServer := setupMockServer(t)
	defer mockServer.Close()

	u, err := url.Parse(mockServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse mock server URL: %v", err)
	}

	host, portStr, err := net.SplitHostPort(u.Host)
	if err != nil {
		t.Fatalf("Failed to split host and port: %v", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("Failed to convert port to integer: %v", err)
	}

	tests := []struct {
		name    string
		cfg     HttpBufferConfig
		entries []models.EntryData
		wantErr bool
	}{
		{
			name: "Empty Entries",
			cfg: HttpBufferConfig{
				Protocol:      "http",
				Address:       host,
				Port:          port,
				Endpoint:      "/test",
				ContentType:   "application/json",
				ItemSeparator: "\n",
			},
			entries: []models.EntryData{},
			wantErr: false,
		},
		{
			name: "Valid Entries",
			cfg: HttpBufferConfig{
				Protocol:      "http",
				Address:       host,
				Port:          port,
				Endpoint:      "/test",
				ContentType:   "application/json",
				ItemSeparator: "\n",
			},
			entries: []models.EntryData{
				{Body: []byte(`{"key1": "value1", "key2": "pippo1"}`)},
				{Body: []byte(`{"key3": "value3", "key4": "pippo2"}`)},
			},
			wantErr: false,
		},
		{
			name: "Server Error",
			cfg: HttpBufferConfig{
				Protocol:      "http",
				Address:       host,
				Port:          port,
				Endpoint:      "/error",
				ContentType:   "application/json",
				ItemSeparator: "\n",
			},
			entries: []models.EntryData{
				{Body: []byte(`{"error!": "value1", "key2": "pippo1"}`)},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Notify(tc.entries, false)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
