package models

type Tweet struct {
	Id int `json:"id"`
	Text string `json:"text"`
	IsSend int `json:"isSend"`
	TimeSend int `json:"timeSend"`
}