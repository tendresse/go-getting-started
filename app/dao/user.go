package dao

import (
	// "encoding/json"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"

	"gopkg.in/pg.v5"
)

type User struct {
}

func (c User) CreateUser(user *models.User) error {
	return config.Global.DB.Insert(user)
}

func (c User) GetUser(user *models.User) error {
	return config.Global.DB.Select(&user)
}

func (c User) GetFullUser(user *models.User) error {
	return config.Global.DB.Model(&user).Column("user.*", "Friends").First()
}

func (c User) GetUserByUsername(username string, user *models.User) error {
	return config.Global.DB.Model(&user).Where("username = ?",username).Select()
}

func (c User) GetUsersByIds(ids []int, users []*models.User) error {
	_, err := config.Global.DB.Query(&users, `SELECT * FROM users WHERE id IN (?)`, pg.In(ids))
	return err
}

func (c User) GetUserFriends(user *models.User, friends []*models.User) error {
	return nil
}

func (c User) GetUserAchievements(id int, user *models.User) error {
	return nil
}
/*
 var story Story
   err = db.Model(&story).
      Column("story.*", "Author").
      Where("story.id = ?", story1.Id).
      Select()
 */