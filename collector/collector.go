package collector

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"bnbl.io/csp_exporter/csp"
)

type Collector interface {
	Collect(app string, report *csp.Report)
	Handler() http.Handler
}

type collector struct {
	reportCounter *prometheus.CounterVec
}

const (
	labelApp                = "app"
	labelBlockedURI         = "blocked_uri"
	labelDisposition        = "disposition"
	labelDocumentURI        = "document_uri"
	labelEffectiveDirective = "effective_directive"
	labelOriginalPolicy     = "original_policy"
	labelReferrer           = "referrer"
	labelScriptSample       = "script_sample"
	labelStatusCode         = "status_code"
	labelViolatedDirective  = "violated_directive"
	labelSourceFile         = "source_file"
	labelLineNumber         = "line_number"
	labelColumnNumber       = "column_number"
)

func NewCollector() (Collector, error) {
	c := &collector{
		reportCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "csp",
			Name:      "violation_reports_total",
			Help:      "Count of CSP violation reports.",
		}, []string{
			labelApp,
			labelBlockedURI,
			labelDisposition,
			labelDocumentURI,
			labelEffectiveDirective,
			labelOriginalPolicy,
			labelReferrer,
			labelScriptSample,
			labelStatusCode,
			labelViolatedDirective,
			labelSourceFile,
			labelLineNumber,
			labelColumnNumber,
		}),
	}
	if err := prometheus.Register(c.reportCounter); err != nil {
		log.Fatalf("could not register report counter with prometheus: %v", err)
	}
	return c, nil
}

func (c *collector) Collect(app string, report *csp.Report) {
	labels := prometheus.Labels{
		labelApp:                app,
		labelBlockedURI:         report.BlockedURI,
		labelDisposition:        report.Disposition,
		labelDocumentURI:        report.DocumentURI,
		labelEffectiveDirective: report.EffectiveDirective,
		labelOriginalPolicy:     report.OriginalPolicy,
		labelReferrer:           report.Referrer,
		labelScriptSample:       report.ScriptSample,
		labelViolatedDirective:  report.ViolatedDirective,
		labelSourceFile:         report.SourceFile,
		labelStatusCode:         "",
		labelLineNumber:         "",
		labelColumnNumber:       "",
	}

	if report.StatusCode > 0 {
		labels[labelStatusCode] = strconv.Itoa(report.StatusCode)
	}
	if report.LineNumber > 0 {
		labels[labelLineNumber] = strconv.Itoa(report.LineNumber)
	}
	if report.ColumnNumber > 0 {
		labels[labelColumnNumber] = strconv.Itoa(report.ColumnNumber)
	}

	c.reportCounter.With(labels).Inc()
}

func (c *collector) Handler() http.Handler {
	return promhttp.Handler()
}
