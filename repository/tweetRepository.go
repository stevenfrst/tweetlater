package repository

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	//guuid "github.com/google/uuid"
	"tweetlater/models"
)

type IAppRepository interface {
	FindAll() ([]*models.Tweet, error)
	CreateBasic(tweet *models.Tweet) error
	CreatePremium(tweet *models.Tweet) error
	Delete(id string) error
	IsPremium(username string) *models.User
	FindOneById(tweet *models.Tweet) (*models.Tweet, error)
}

var (
	appQueries = map[string]string{
		"addTweetLaterBasic":"INSERT into draft_tweet(text) values(?)",
		"addTweetLaterPremium":"INSERT into draft_tweet(text,timeSend) values(?,?)",
		"deleteTweetLaterByID":"DELETE FROM draft_tweet WHERE id=?",
		"getAllDraftTweetLater":"SELECT id,text,isSend,timeSend FROM draft_tweet",
		"GetINFO":"SELECT id,username,status FROM user WHERE username=?",
		"GetTweet":"SELECT id,text,isSend,timeSend FROM draft_tweet WHERE id=?",
		"Sended":"UPDATE draft_tweet SET isSend=1 WHERE id=?",
	}
)

type appRepository struct {
	db *sql.DB
	ps map[string]*sql.Stmt
}


func (u *appRepository) FindOneById(tweet *models.Tweet) (*models.Tweet, error) {
	row := u.ps["GetTweet"].QueryRow(tweet.Id)
	res := new(models.Tweet)
	err := row.Scan(&res.Id, &res.Text, &res.IsSend,&res.TimeSend)
	log.Println("checkpoint 1",res)
	if err != nil {
		log.Println("DB GET 1",err)
		return nil, err
	}
	_, err = u.ps["Sended"].Exec(tweet.Id)
	log.Println("checkpoint 2")
	if err != nil {
		log.Println("DB SEND",err)
		return nil, err
	}

	return res, nil
}

func (u *appRepository) FindAll() ([]*models.Tweet, error) {
	rows,_ := u.db.Query(appQueries["getAllDraftTweetLater"])
	//res := new(models.Tweet)
	var res []*models.Tweet
	for rows.Next(){
		var each = models.Tweet{}
		err := rows.Scan(&each.Id, &each.Text, &each.IsSend, &each.TimeSend)
		if err != nil {
			return nil, err
		}
		res = append(res,&each)
	}
	log.Println(res)
	return res, nil
}

func (u *appRepository) CreateBasic(tweet *models.Tweet) error {
	tx, _ := u.db.Begin()
	_, err := tx.Stmt(u.ps["addTweetLaterBasic"]).Exec(tweet.Text)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}

func (u *appRepository) CreatePremium(tweet *models.Tweet) error {
	tx, _ := u.db.Begin()
	_, err := tx.Stmt(u.ps["addTweetLaterPremium"]).Exec(tweet.Text,tweet.TimeSend)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}

func (u *appRepository) Delete(id string) error {
	_, err := u.ps["deleteTweetLaterByID"].Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *appRepository) IsPremium(username string) *models.User{
	row := u.ps["GetINFO"].QueryRow(username)
	res := new(models.User)
	err := row.Scan(&res.Id, &res.Username, &res.Status)
	if err != nil {
		return nil
	}
	return res
}

func NewAppRepository(db *sql.DB) IAppRepository {
	ps := make(map[string]*sql.Stmt, len(appQueries))

	for n, v := range appQueries {
		p, err := db.Prepare(v)
		if err != nil {
			panic(err)
		}
		ps[n] = p
	}
	return &appRepository{db, ps}
}
