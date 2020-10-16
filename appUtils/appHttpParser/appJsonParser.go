package appHttpParser

import (
	"encoding/json"
	"net/http"
	"strings"
)

type JsonParser struct{}

func NewJsonParser() *JsonParser {
	return &JsonParser{}
}

func (p *JsonParser) Parse(r *http.Request, model interface{}) error {
	if r.Body == nil {
		return nil
	}
	defer r.Body.Close()
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(model); err != nil {
			return err
		}
	}
	return nil
}
