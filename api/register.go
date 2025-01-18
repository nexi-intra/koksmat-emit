// Package api provides the implementation of the Koksmat Webhooks API.
// This package sets up the necessary endpoints and starts the HTTP server.
//
// The API includes the following endpoints:
// - POST /api/v1/github: Handles GitHub webhooks.
// - POST /api/v1/officegraph/notify: Handles Microsoft Graph notifications.
//
// The service also includes a profiler available at /debug/core and
// documentation available at /docs.
//
// The service is built using the swaggest/rest and go-chi/chi packages.
//
// For more information, check out the documentation at:
// https://pkg.go.dev/github.com/swaggest/rest#section-readme
/*
---
title: Koksmat Webhooks API
----
Check out the documentation at
https://pkg.go.dev/github.com/swaggest/rest#section-readme
*/
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/nexi-intra/koksmat-emit/internal/emitter"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v4emb"
)

func addCoreEndpoints(s *web.Service, app *emitter.App) {

	s.Method(http.MethodPost, "/api/v1/github", nethttp.NewHandler(webhook_GitHub(app)))
	s.MethodFunc(http.MethodPost, "/api/v1/officegraph/notify", webhook_MicrosoftGraph(app))

	s.Mount("/debug/core", middleware.Profiler())
}

func Start(port string, app *emitter.App) {

	service := web.NewService(openapi3.NewReflector())

	service.OpenAPISchema().SetTitle("Koksmat Webhooks API")
	service.OpenAPISchema().SetDescription("This service provides API to expose web hooks")
	service.OpenAPISchema().SetVersion("V1.0.0")

	addCoreEndpoints(service, app)
	service.Docs("/docs", swgui.New)
	app.Obs.Info("Server starting, view documentation at http://localhost" + port + "/docs")
	go http.ListenAndServe(port, service)

}
