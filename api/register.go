/*
---
title: Koksmat Webhooks API
----
Check out the documentation at
https://pkg.go.dev/github.com/swaggest/rest#section-readme
*/
package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v4emb"
)

func addCoreEndpoints(s *web.Service) {

	s.Method(http.MethodPost, "/api/v1/github", nethttp.NewHandler(webhook_GitHub()))
	s.MethodFunc(http.MethodPost, "/api/v1/officegraph/notify", webhook_MicrosoftGraph)

	s.Mount("/debug/core", middleware.Profiler())
}

func Start(port string) {
	service := web.NewService(openapi31.NewReflector())

	service.OpenAPISchema().SetTitle("Koksmat Webhooks API")
	service.OpenAPISchema().SetDescription("This service provides API to expose web hooks")
	service.OpenAPISchema().SetVersion("v1.0.0")

	addCoreEndpoints(service)
	service.Docs("/docs", swgui.New)
	log.Printf("Server starting, view documentation at http://localhost%s/docs", port)
	go http.ListenAndServe(port, service)

}
