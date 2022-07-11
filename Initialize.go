package gcplog

import (
	"context"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/trace"
	"net/http"
)

func Initialize(name, projID string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if projID == "" {
				panic("projectID is required")
			}
			if name == "" {
				panic("name is required")
			}
			projectID = projID
			ctx, done := WithContext(r, name)
			defer done()
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type LoggerParam struct {
	ProjectID string
	SpanID    string
	TraceID   string
}

var projectID string

func WithContext(r *http.Request, label string) (context.Context, func()) {
	ctx := r.Context()
	span := new(trace.Span)
	httpFormat := propagation.HTTPFormat{}
	if sc, ok := httpFormat.SpanContextFromRequest(r); ok {
		ctx, span = trace.StartSpanWithRemoteParent(ctx, label, sc)
	} else {
		ctx, span = trace.StartSpan(ctx, label)
	}

	return ctx, span.End
}
