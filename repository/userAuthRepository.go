package repository

import (
	"database/sql"
	"log"
	"tweetlater/models"
)

var (
	authQueries = map[string]string{
		"updatePassword":"UPDATE user SET password=? WHERE id=?",
		"findOneUserAuthByUserNameAndPassword":"SELECT id,username FROM user WHERE username=? AND password=?",
	}
)

type IUserAuthRepository interface {
	FindOneByUserNameAndPassword(userName string, password string) *models.User
	UpdatePassword(id string, newPassword string) error
}

type userAuthRepository struct {
	db *sql.DB
	ps map[string]*sql.Stmt
}

func (u *userAuthRepository) FindOneByUserNameAndPassword(userName string, password string) *models.User {
	row := u.ps["findOneUserAuthByUserNameAndPassword"].QueryRow(userName, password)
	res := new(models.User)
	err := row.Scan(&res.Id, &res.Username)
	log.Println(res)
	if err != nil {
		return nil
	}
	return res
}

func (u *userAuthRepository) UpdatePassword(id string, newPassword string) error {
	_, err := u.ps["updateUserPassword"].Exec(newPassword, id)
	if err != nil {
		return err
	}
	return nil
}

func NewUserAuthRepository(db *sql.DB) IUserAuthRepository {
	ps := make(map[string]*sql.Stmt, len(authQueries))

	for n, v := range authQueries {
		p, err := db.Prepare(v)
		if err != nil {
			panic(err)
		}
		ps[n] = p
	}

	return &userAuthRepository{
		db, ps,
	}
}