package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"tweetlater/appUtils/appCookieStore"
	"tweetlater/appUtils/appHttpParser"
	"tweetlater/appUtils/appHttpResponse"
	"tweetlater/delivery"
	"tweetlater/infra"
	"tweetlater/manager"
	appMiddleware "tweetlater/middleware"
)

type appRouter struct {
	app                       *tweetLater
	cookieStore               *sessions.CookieStore
	logRequestMiddleware      *appMiddleware.LogRequestMiddleware
	tokenValidationMiddleware *appMiddleware.TokenValidationMiddleware
	parser                    *appHttpParser.JsonParser
	responder                 appHttpResponse.IResponder
	infra                     infra.Infra
}

type appRoutes struct {
	del delivery.IDelivery
	mdw []mux.MiddlewareFunc
}

func (ar *appRouter) InitMainRouter() {
	ar.app.router.Use(ar.logRequestMiddleware.Log)
	var serviceManager = manager.NewServiceManger(ar.infra)
	appRoutes := []appRoutes{
		{
			del: delivery.NewAuthDelivery(ar.app.router, ar.cookieStore, ar.parser, ar.responder, serviceManager.UserAuthUseCase()),
			mdw: nil,
		},
		{
			del: delivery.NewUserDelivery(ar.app.router, ar.parser, ar.responder, serviceManager.UserUseCase()),
			mdw: []mux.MiddlewareFunc{
				ar.tokenValidationMiddleware.Validate,
			},
		},
		{
			del: delivery.NewAppDelivery(ar.app.router, ar.parser, ar.responder, serviceManager.AppUseCase(),ar.infra),
			mdw: []mux.MiddlewareFunc{
				ar.tokenValidationMiddleware.Validate,
			},
		},
	}
	for _, r := range appRoutes {
		r.del.InitRoute(r.mdw...)
	}
}

func NewAppRouter(app *tweetLater) *appRouter {
	var cookieStore = appCookieStore.NewAppCookieStore().Store
	return &appRouter{
		app,
		cookieStore,
		appMiddleware.NewLogRequestMiddleware(),
		appMiddleware.NewTokenValidationMiddleware(cookieStore),
		appHttpParser.NewJsonParser(),
		appHttpResponse.NewJSONResponder(),
		app.infra,
	}
}