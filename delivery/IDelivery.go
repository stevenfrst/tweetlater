package delivery

import "github.com/gorilla/mux"

type IDelivery interface {
	InitRoute(mdw ...mux.MiddlewareFunc)
}
