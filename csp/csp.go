package csp

import (
	"encoding/json"
	"io"
)

type Report struct {
	BlockedURI         string `json:"blocked-uri"`
	Disposition        string `json:"disposition"`
	DocumentURI        string `json:"document-uri"`
	EffectiveDirective string `json:"effective-directive"`
	OriginalPolicy     string `json:"original-policy"`
	Referrer           string `json:"referrer"`
	ScriptSample       string `json:"script-sample"`
	StatusCode         int    `json:"status-code"`
	ViolatedDirective  string `json:"violated-directive"`
	SourceFile         string `json:"source-file"`
	LineNumber         int    `json:"line-number"`
	ColumnNumber       int    `json:"column-number"`
}

func ReadReport(r io.Reader) (*Report, error) {
	type body struct {
		Report Report `json:"csp-report"`
	}
	d := json.NewDecoder(r)
	var b body
	if err := d.Decode(&b); err != nil {
		return nil, err
	}
	return &b.Report, nil
}
