package delivery

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"tweetlater/appUtils/appHttpParser"
	"tweetlater/appUtils/appHttpResponse"
	"tweetlater/appUtils/appStatus"
	"tweetlater/models"
	usecase "tweetlater/usecases"
)

const (
	userMainRoute = "/user"
)

type UserDelivery struct {
	router    *mux.Router
	parser    *appHttpParser.JsonParser
	responder appHttpResponse.IResponder
	service   usecase.IUserUseCase
}

func NewUserDelivery(router *mux.Router, parser *appHttpParser.JsonParser, responder appHttpResponse.IResponder, service usecase.IUserUseCase) *UserDelivery {
	return &UserDelivery{
		router, parser, responder, service,
	}
}

func (d *UserDelivery) InitRoute(mdw ...mux.MiddlewareFunc) {
	userRouter := d.router.PathPrefix(userMainRoute).Subrouter()
	userRouter.Use(mdw...)

	userRouter.HandleFunc("", d.userPostRoute).Methods("POST")
	userRouter.HandleFunc("", d.userPutRoute).Methods("PUT")
	userRouter.HandleFunc("", d.userDeleteRoute).Methods("DELETE")
	userRouter.HandleFunc("", d.userRoute).Methods("GET")
}

func (d *UserDelivery) upgradeRoute(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := d.parser.Parse(r, &user); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := d.service.UpgradeAccount(&user); err != nil {
		d.responder.Error(w, appStatus.UnknownError, err.Error())
	}
	d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), user)
}


func (d *UserDelivery) userRoute(w http.ResponseWriter, r *http.Request) {
	userId, isExist := r.URL.Query()["id"]
	if isExist {
		users := d.service.GetUserInfo(userId[0])
		if users != nil {
			d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), users)
		} else {
			d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), "no record found")
		}

	} else {
		d.responder.Error(w, appStatus.ErrorLackInfo, "Please provide some IDs")
	}

}

func (d *UserDelivery) userPostRoute(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := d.parser.Parse(r, &newUser); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := d.service.Register(&newUser); err != nil {
		d.responder.Error(w, appStatus.UnknownError, err.Error())
	}
	d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), newUser)
}

func (d *UserDelivery) userPutRoute(w http.ResponseWriter, r *http.Request) {
	userId, isExist := r.URL.Query()["id"]
	if isExist {
		var usrReq models.User
		if err := d.parser.Parse(r, &usrReq); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		userUpdate := d.service.UpdateInfo(userId[0], &usrReq)
		d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), userUpdate)
	} else {
		msg := appStatus.StatusText(appStatus.ErrorLackInfo)
		d.responder.Error(w, appStatus.ErrorLackInfo, fmt.Sprintf(msg, "ID"))
	}
}

func (d *UserDelivery) userDeleteRoute(w http.ResponseWriter, r *http.Request) {
	userId, isExist := r.URL.Query()["id"]
	if isExist {
		if err := d.service.Unregister(userId[0]); err != nil {
			d.responder.Error(w, appStatus.ErrorLackInfo, err.Error())
		}
		d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), nil)
	} else {
		msg := appStatus.StatusText(appStatus.ErrorLackInfo)
		d.responder.Error(w, appStatus.ErrorLackInfo, fmt.Sprintf(msg, "ID"))
	}

}