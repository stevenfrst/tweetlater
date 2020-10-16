package models

type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status int `json:"status"`
	PostTweet int `json:"postTweet"`
	SendMsg int `json:"sendMsg"`
}