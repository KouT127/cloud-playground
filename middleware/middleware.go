package middleware

import (
	"context"
	"fmt"
	"github.com/KouT127/pub-sub-practice/config"
	"github.com/KouT127/pub-sub-practice/model"
	"log"
	"net/http"
	"strings"
)

var CloudTraceContext = struct{}{}

func CloudTraceMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var trace string
		if config.ProjectID != "" {
			traceHeader := r.Header.Get("X-Cloud-Trace-Context")
			traceParts := strings.Split(traceHeader, "/")
			if len(traceParts) > 0 && len(traceParts[0]) > 0 {
				trace = fmt.Sprintf("projects/%s/traces/%s", config.ProjectID, traceParts[0])
			}
		}
		log.Println(model.Entry{
			Message:   "middleware trace",
			Component: "trace",
			Trace:     trace,
		})
		r = r.WithContext(context.WithValue(ctx, CloudTraceContext, trace))
		next.ServeHTTP(w, r)
	}
}
