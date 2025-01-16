package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
)

func addCoreEndpoints(s *web.Service, jwtAuth func(http.Handler) http.Handler) {

	s.Method(http.MethodPost, "/api/v1/github", nethttp.NewHandler(github.GitHubWebhook()))
	//s.Use(rateLimitByAppId(50))
	s.MethodFunc(http.MethodPost, "/api/v1/subscription/notify", validateSubscription)
	s.Route("/v1/webhooks", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			//r.Use(adminAuth, nethttp.HTTPBasicSecurityMiddleware(s.OpenAPICollector, "User", "User access"))
			r.Use(jwtAuth, nethttp.HTTPBearerSecurityMiddleware(s.OpenAPICollector, "Bearer", "", ""))

			r.Method(http.MethodGet, "/", nethttp.NewHandler(getWebHooks()))

		})
	})
	s.Mount("/debug/core", middleware.Profiler())
}
