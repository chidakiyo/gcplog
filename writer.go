package gcplog

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opencensus.io/trace"
	"log"
	"os"
)

type Writer func(context.Context, string, string, ...interface{})

var _ Writer = outputText
var _ Writer = outputJSON

func outputStructure(ctx context.Context, severity string, analysisLog interface{}) {
	sc := trace.FromContext(ctx).SpanContext()
	payload := &struct {
		Severity      string      `json:"severity"`
		StructuredLog interface{} `json:"structure"`
		Trace         string      `json:"logging.googleapis.com/trace"`
		SpanID        string      `json:"logging.googleapis.com/spanId"`
	}{
		Severity:      severity,
		StructuredLog: analysisLog,
		Trace:         fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID.String()),
		SpanID:        sc.SpanID.String(),
	}

	json.NewEncoder(os.Stdout).Encode(payload)
}

func outputJSON(ctx context.Context, severity, format string, a ...interface{}) {
	sc := trace.FromContext(ctx).SpanContext()
	payload := &struct {
		Severity string `json:"severity"`
		Message  string `json:"message"`
		Trace    string `json:"logging.googleapis.com/trace"`
		SpanID   string `json:"logging.googleapis.com/spanId"`
	}{
		Severity: severity,
		Message:  fmt.Sprintf(format, a...),
		Trace:    fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID.String()),
		SpanID:   sc.SpanID.String(),
	}

	json.NewEncoder(os.Stdout).Encode(payload)
}

func outputText(ctx context.Context, severity, format string, a ...interface{}) {
	log.Printf(severity+": "+format, a...)
}
