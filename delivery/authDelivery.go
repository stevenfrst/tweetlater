package delivery

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"tweetlater/appUtils/appHttpParser"
	"tweetlater/appUtils/appHttpResponse"
	"tweetlater/appUtils/appStatus"
	"tweetlater/models"
	usecase "tweetlater/usecases"
)

const (
	loginRoute  = "/login"
	logoutRoute = "/logout"
)

type AuthDelivery struct {
	router      *mux.Router
	cookieStore *sessions.CookieStore
	parser      *appHttpParser.JsonParser
	responder   appHttpResponse.IResponder
	service     usecase.IUserAuthUseCase // ! Do IT LATER
}

func NewAuthDelivery(router *mux.Router, cookie *sessions.CookieStore, parser *appHttpParser.JsonParser, responder appHttpResponse.IResponder, service usecase.IUserAuthUseCase) *AuthDelivery {
	return &AuthDelivery{router, cookie, parser, responder, service}
}

func (d *AuthDelivery) InitRoute(mdw ...mux.MiddlewareFunc) {
	d.router.HandleFunc(loginRoute, d.authRoute).Methods("POST")
	d.router.HandleFunc(logoutRoute, d.authLogoutRoute).Methods("GET")
}

func (d *AuthDelivery) authRoute(w http.ResponseWriter, r *http.Request) {
	session, _ := d.cookieStore.Get(r, "app-cookie")
	var userAuth models.User
	if err := d.parser.Parse(r, &userAuth); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	userInfo := d.service.UserNamePasswordValidation(userAuth.Username, userAuth.Password)
	//userInfo := d.service.UserNamePasswordValidation("steven", "ono")

	if userInfo == nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	session.Values["authenticated"] = true
	err := session.Save(r, w)
	d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), userInfo)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}


func (d *AuthDelivery) authLogoutRoute(w http.ResponseWriter, r *http.Request) {
	session, _ := d.cookieStore.Get(r, "app-cookie")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), "Logout")
}
