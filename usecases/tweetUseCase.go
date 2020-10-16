package usecase

import (
	"errors"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"log"
	"time"
	"tweetlater/infra"
	"tweetlater/models"
	"tweetlater/repository"
)

type IAppUseCase interface {
	AddTweetBasic(newTweet *models.Tweet) error
	AddTweetPremium(newTweet *models.Tweet,username string) error
	GetAll() ([]*models.Tweet,error)
	DeleteTweetLater(Id string) error
	SendTweet(tweet *models.Tweet,infra infra.Infra) error
}

type AppUseCase struct {
	appRepo repository.IAppRepository
}

func NewAppUseCase(appRepo repository.IAppRepository) IAppUseCase {
	return &AppUseCase{
		appRepo,
	}
}




func GetClient(creds *models.Credentials) (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}

	log.Printf("User's ACCOUNT:\n%+v\n", user)
	log.Println("checkpoint 3")
	return client, nil
}

func (uc *AppUseCase) SendTweet(tweet *models.Tweet,infra infra.Infra) error {
	log.Println("checkpoint 4")
	client,err := GetClient(infra.TweetApi())
	if err != nil {
		log.Println(err)
	}
	log.Println("checkpoint 5")
	msg,err	 := uc.appRepo.FindOneById(tweet)
	timeLast := msg.TimeSend
	log.Println(timeLast)
	if timeLast > 0 {
		time.Sleep(time.Duration(timeLast)*time.Minute)
		log.Println("Sleeping Wait")
	}

	if err != nil {
		log.Println("USECASE",err)
		log.Println("checkpoint 00")
		return err
	}
	log.Println(msg.Text)
	_, _, err = client.Statuses.Update(msg.Text, nil)
	log.Println("checkpoint 6")
	if err != nil {
		log.Println("USECASE",err)
		return err
	}
	log.Println("checkpoint 6")
	return nil

}

func (uc *AppUseCase) AddTweetBasic(newTweet *models.Tweet) error {
	return uc.appRepo.CreateBasic(newTweet)
}

func (uc *AppUseCase) AddTweetPremium(newTweet *models.Tweet,username string) error {
	userInfo := uc.appRepo.IsPremium(username)
	if userInfo.Status == 0 {
		return errors.New("BJIR") // ! IMPROVE LATER
	}
	return uc.appRepo.CreatePremium(newTweet)
}

func (uc *AppUseCase) GetAll() ([]*models.Tweet,error) {
	return uc.appRepo.FindAll()
}

func (uc *AppUseCase) DeleteTweetLater(Id string) error {
	return uc.appRepo.Delete(Id)
}





