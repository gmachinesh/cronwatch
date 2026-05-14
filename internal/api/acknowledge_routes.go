package api

import "net/http"

func registerAcknowledgeRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	mux.Handle("/api/v1/jobs/acknowledge",
		authMiddleware(
			loggingMiddleware(
				methodGuard(http.MethodPost,
					http.HandlerFunc(handleAcknowledgeJob),
				),
			),
		),
	)

	mux.Handle("/api/v1/jobs/unacknowledge",
		authMiddleware(
			loggingMiddleware(
				methodGuard(http.MethodPost,
					http.HandlerFunc(handleUnacknowledgeJob),
				),
			),
		),
	)

	mux.Handle("/api/v1/jobs/acknowledged",
		authMiddleware(
			loggingMiddleware(
				methodGuard(http.MethodGet,
					http.HandlerFunc(handleAcknowledgedJobs),
				),
			),
		),
	)
}
