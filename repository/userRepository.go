package repository

import (
	"database/sql"
	"log"
	"tweetlater/models"
	guuid "github.com/google/uuid"
)

type IUserRepository interface {
	FindOneById(id string) (*models.User, error)
	Create(newUser *models.User) error
	Update(id string, newUser *models.User) error
	Delete(id string) error
	Upgrade(user *models.User) error

}

var (
	userQueries = map[string]string{
		"insertUser":     "INSERT into user(id,username,password) values(?,?,?)",
		"updateUserByID":  "UPDATE user SET password=? WHERE id=?",
		"findOneUserByID":    "SELECT id,username,password,status FROM user WHERE id=?",
		"deleteUserByID":     "DELETE FROM user WHERE id=?",
		"upgradeACC":"UPDATE user SET status=1 WHERE id=?",
	}
)

type userRepository struct {
	db *sql.DB
	ps map[string]*sql.Stmt
}

func (u *userRepository) FindOneById(id string) (*models.User, error) {
	row := u.ps["findOneUserByID"].QueryRow(id)
	res := new(models.User)
	err := row.Scan(&res.Id, &res.Username, &res.Password,&res.Status)
	if err != nil {
		return nil, err
	}
	log.Println(res)
	return res, nil
}

func (u *userRepository) Create(newUser *models.User) error {
	id := guuid.New()
	newUser.Id = id.String()
	tx, _ := u.db.Begin()
	_, err := tx.Stmt(u.ps["insertUser"]).Exec(newUser.Id, newUser.Username, newUser.Password)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}

func (u *userRepository) Update(Id string, newUser *models.User) error {
	_, err := u.ps["updateUserByID"].Exec(newUser.Id, newUser.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) Upgrade(user *models.User) error {
	_, err := u.ps["upgradeACC"].Exec(user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepository) Delete(id string) error {
	_, err := u.ps["deleteUserByID"].Exec(id)
	if err != nil {
		return err
	}
	return nil
}



func NewUserRepository(db *sql.DB) IUserRepository {
	ps := make(map[string]*sql.Stmt, len(userQueries))

	for n, v := range userQueries {
		p, err := db.Prepare(v)
		if err != nil {
			panic(err)
		}
		ps[n] = p
	}
	return &userRepository{db, ps}
}
