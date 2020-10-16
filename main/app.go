package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tweetlater/infra"
)


type tweetLater struct {
	infra  infra.Infra
	router *mux.Router
}

func (app *tweetLater) run() {
	h := app.infra.ApiServer()
	log.Println("Listening on", h)
	NewAppRouter(app).InitMainRouter()
	err := http.ListenAndServe(h, app.router)
	if err != nil {
		log.Fatalln(err)
	}
}

func NewTweetLaterApp() *tweetLater {
	r := mux.NewRouter()
	appInfra := infra.NewInfra()
	return &tweetLater{
		infra:  appInfra,
		router: r,
	}
}

func main() {
	NewTweetLaterApp().run()
}
