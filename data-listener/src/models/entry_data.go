package models

import (
	"encoding/json"
	"net/url"
)

type EntryData struct {
	Headers map[string][]string
	Body    []byte
	Url     url.URL
}

func (e *EntryData) String() string {
	data := struct {
		Headers map[string][]string `json:"headers"`
		Body    string              `json:"body"`
		Url     string              `json:"url"`
	}{
		Headers: e.Headers,
		Body:    string(e.Body),
		Url:     e.Url.String(),
	}

	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}
