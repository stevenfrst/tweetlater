package delivery

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tweetlater/appUtils/appHttpParser"
	"tweetlater/appUtils/appHttpResponse"
	"tweetlater/appUtils/appStatus"
	"tweetlater/infra"
	"tweetlater/models"
	usecase "tweetlater/usecases"
)

const tweetAppRoute = "/app"


type AppDelivery struct {
	router    *mux.Router
	parser    *appHttpParser.JsonParser
	responder appHttpResponse.IResponder
	service   usecase.IAppUseCase
	infra infra.Infra
}

func NewAppDelivery(router *mux.Router, parser *appHttpParser.JsonParser, responder appHttpResponse.IResponder, service usecase.IAppUseCase, infra infra.Infra) *AppDelivery {
	return &AppDelivery{
		router, parser, responder, service,infra,
	}
}

func (d *AppDelivery) InitRoute(mdw ...mux.MiddlewareFunc) {
	userRouter := d.router.PathPrefix(tweetAppRoute).Subrouter()
	userRouter.Use(mdw...)

	userRouter.HandleFunc("", d.GetAllRoute).Methods("GET")
	userRouter.HandleFunc("", d.PostBasicTweet).Methods("POST")
	userRouter.HandleFunc("/post", d.PostPremiumTweet).Methods("POST")
	userRouter.HandleFunc("", d.DeleteTweetLater).Methods("DELETE")
	userRouter.HandleFunc("/tweet", d.SendTweet).Methods("POST")
	log.Println("CHECKPOINT 8")

}

func (d *AppDelivery) SendTweet(w http.ResponseWriter, r *http.Request) {
	var newTweet models.Tweet
	if err := d.parser.Parse(r, &newTweet); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	log.Println("CHECKPOINT 9")

	if err := d.service.SendTweet(&newTweet,d.infra); err != nil {
		d.responder.Error(w, appStatus.UnknownError, err.Error())
	}
	log.Println("CHECKPOINT 10")

	d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), newTweet)

}

func (d *AppDelivery) GetAllRoute(w http.ResponseWriter, r *http.Request) {
	tweet,err := d.service.GetAll()
	if err != nil {
		log.Println(err)
	}
	log.Println(tweet)
	d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), tweet)
}

func (d *AppDelivery) PostBasicTweet(w http.ResponseWriter,r *http.Request) {
	var newTweet models.Tweet
	if err := d.parser.Parse(r, &newTweet); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := d.service.AddTweetBasic(&newTweet); err != nil {
		d.responder.Error(w, appStatus.UnknownError, err.Error())
	}
	d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), newTweet)
}

func (d *AppDelivery) PostPremiumTweet(w http.ResponseWriter,r *http.Request) {
	var newTweet models.Tweet
	username, _:= r.URL.Query()["id"]

	if err := d.parser.Parse(r, &newTweet); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := d.service.AddTweetPremium(&newTweet,username[0]); err != nil {
		d.responder.Error(w, appStatus.UnknownError, err.Error())
	}
	d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), newTweet)
}

func (d *AppDelivery) DeleteTweetLater(w http.ResponseWriter, r *http.Request) {
	tweet, isExist := r.URL.Query()["id"]
	if isExist {
		if err := d.service.DeleteTweetLater(tweet[0]); err != nil {
			d.responder.Error(w, appStatus.ErrorLackInfo, err.Error())
		}
		d.responder.Data(w, appStatus.Success, appStatus.StatusText(appStatus.Success), nil)
	} else {
		msg := appStatus.StatusText(appStatus.ErrorLackInfo)
		d.responder.Error(w, appStatus.ErrorLackInfo, fmt.Sprintf(msg, "ID"))
	}

}

func (d *AppDelivery) SendTweets(w http.ResponseWriter, r *http.Request) {

}

