package httpapi

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"MrFandore/Prac4/internal/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path
		method := r.Method

		metrics.HttpRequestsTotal.WithLabelValues(method, path).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)

		if lrw.StatusCode() >= 400 {
			metrics.HttpErrorsTotal.WithLabelValues(
				method,
				path,
				strconv.Itoa(lrw.StatusCode()),
			).Inc()
		}

		if strings.HasPrefix(path, "/students/") && path != "/students/" {
			parts := strings.Split(strings.TrimPrefix(path, "/students/"), "/")
			if len(parts) > 0 && parts[0] != "" {
				studentID := parts[0]
				metrics.StudentRequestDuration.WithLabelValues(method, studentID).Observe(duration)
			}
		}
	})
}
