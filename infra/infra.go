package infra

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"tweetlater/models"
)

var (
	cfg *viper.Viper
)

type Infra interface {
	SqlDb() *sql.DB
	Config() *viper.Viper
	ApiServer() string
	TweetApi() *models.Credentials
}

type infra struct{}

func NewInfra() Infra {
	return &infra{}
}

func (i *infra) SqlDb() *sql.DB {
	dbUser := i.Config().GetString("DBUSER")
	dbPassword := i.Config().GetString("DBPASSWORD")
	dbHost := i.Config().GetString("DBHOST")
	dbPort := i.Config().GetString("DBPORT")
	schema := i.Config().GetString("DBSCHEMA")
	dbEngine := i.Config().GetString("DBENGINE")

	db, err := sql.Open(dbEngine, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, schema))
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}

func (i *infra) ApiServer() string {
	host := i.Config().GetString("HTTPHOST")
	port := i.Config().GetString("HTTPPORT")
	return fmt.Sprintf("%s:%s", host, port)
}

func (i *infra) Config() *viper.Viper {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	cfg = viper.GetViper()
	return cfg
}

func (i *infra) TweetApi() *models.Credentials{
	acc_token := i.Config().GetString("ACCESS_TOKEN")
	acc_token_scr := i.Config().GetString("ACCESS_TOKEN_SECRET")
	cons_key := i.Config().GetString("CONSUMER_KEY")
	cons_scr := i.Config().GetString("CONSUMER_SECRET")
	creds := models.Credentials{
		AccessToken: acc_token,
		AccessTokenSecret:acc_token_scr,
		ConsumerKey: cons_key,
		ConsumerSecret: cons_scr,
	}

	return &creds
}


