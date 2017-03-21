// controllers/users_controllers.go

package controllers

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"
	"github.com/tendresse/go-getting-started/app/dao"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/googollee/go-socket.io"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/pg.v5"
	"golang.org/x/crypto/bcrypt"
)

type UsersController struct {
}


func (c UsersController) IsAuthentified() bool {
	// if config.Global.CurrentUser == models.User{}
	if config.Global.CurrentUser.ID == 0{
		return false
	}
	return true
}

func (c UsersController) HasRole(roles []string) bool {
	for _,role := range roles {
		for _,user_role := range config.Global.CurrentUser.Roles {
			if strings.Compare(role,user_role) == 0 {
				return true
			}
		}
	}
	return false
}

func (c UsersController) Signup(username string, password string) string {
	userDAO := dao.UserDAO{}
	config.Global.CurrentUser = models.User{}
	user := models.User{}
	if err := userDAO.GetUserByUsername(username, &user); err != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), 0)
		if err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while hashing password"}`
		}
		hashed_password := string(hashed[:])
		user = models.User{Username:username, Passhash:hashed_password}
		if err = userDAO.CreateUser(&user); err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while creating User"}`
		}
		config.Global.CurrentUser = user
		token,err := user.GenerateToken()
		if err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while creating token"}`
		}
		strings.Join([]string{`{"success":true, "token":`, token, "}"}, "")
	}
	return `{"success":false, "error":"username already taken"}`
}

func (c UsersController) Login(username string, password string) string {
	userDAO := dao.UserDAO{}
	config.Global.CurrentUser = models.User{}
	user := models.User{}
	if err := userDAO.GetUserByUsername(username, &user); err != nil {
		return `{"success":false, "error":"username or password incorrect"}`
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Passhash), []byte(password) ); err != nil {
		log.Error(err)
		return `{"success":false, "error":"username or password incorrect"}`
	}
	config.Global.CurrentUser = user
	token,err := GenerateToken(&user)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while creating token"}`
	}
	strings.Join([]string{`{"success":true, "token":`, token, "}"}, "")
}

func (c UsersController) LoginWithToken(user_token string) string {
	config.Global.CurrentUser = models.User{}
	userDAO := dao.UserDAO{}
	token, err := jwt.Parse(user_token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(string(token.Header["alg"]))
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.Global.SecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := models.User{ID: claims["id"]}
		if err := userDAO.GetUser(&user); err != nil {
			log.Error(err)
			return `{"success":false, "error":"invalid token"}`
		}
		config.Global.CurrentUser = user
		return `{"success":true}`
	} else {
		log.Error(err)
		return `{"success":false, "error":"invalid token"}`
	}
}

func (c UsersController) GetUser(user_id int) string {
	userDAO := dao.UserDAO{}
	user := models.User{ID:user_id}
	if err := userDAO.GetUser(&user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	// achievements
	// friends
	b, err := json.Marshal(user)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal user json error"}`
	}
	return strings.Join([]string{`{"success":true, "blog":`,string(b),"}"} ,"")
}

func (c UsersController) GetPendingTendresses(user_token string) string {
	tendresses := []models.Tendresse{}
	if err := config.Global.DB.Model(&config.Global.CurrentUser).Where("Viewed = ?",false).Association("Tags").Find(&tendresses).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while fetching pending tendresses"}`
	}
	b, err := json.Marshal(tendresses)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal tendresses json error"}`
	}
	return strings.Join([]string{`{"success":true, "blog":`,string(b),"}"} ,"")
}

func (c UsersController) UpdateDevice(device_token string) string {
	if err := config.Global.DB.Model(&config.Global.CurrentUser).Update("device",device_token).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	return `{"success":true}`
}

func (c UsersController) AddFriend(username string) string {
	// TODO : LOG WHO ADDED WHO
	user := models.User{}
	if err := config.Global.DB.Where("Username = ?", username).First(&user).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	friends := []models.User{}
	if err := config.Global.DB.Model(&config.Global.CurrentUser).Where("Username = ?",username).Association("Friends").Find(&friends).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"error getting friends"}`
	}
	if len(friends) != 0 {
		return `{"success":false, "error":"friend already added"}`
	}
	if err := config.Global.DB.Model(config.Global.CurrentUser).Association("Friends").Append(user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error adding user"}`
	}
	b, err := json.Marshal(user)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal tendresses json error"}`
	}
	return strings.Join([]string{`{"success":true, "blog":`,string(b),"}"} ,"")
}

func (c UsersController) DeleteFriend(user_id string) string {
	// TODO : LOG WHO ADDED WHO
	user := models.User{}
	if err := config.Global.DB.First(&user, user_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	friends := []models.User{}
	if err := config.Global.DB.Model(&config.Global.CurrentUser).Where("id = ?",user.ID).Association("Friends").Find(&friends).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"error getting friends"}`
	}
	if len(friends) == 0 {
		return `{"success":false, "error":"user already unfriend"}`
	}
	if err := config.Global.DB.Model(config.Global.CurrentUser).Association("Friends").Delete(user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error adding user"}`
	}
	return `{"success":true}`
}

func (c UsersController) GrantUserRole(user_id int, role string) string {
	// TODO : LOG WHO ADDED WHO
	user := models.User{}
	if err := config.Global.DB.First(&user, user_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	for _,user_role := range user.Roles {
		if strings.Compare(role,user_role) == 0 {
			return `{"success":false, "error":"user already has this role"}`
		}
	}
	append(user.Roles,role)
	if err := config.Global.DB.Save(&user).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"error saving new user's roles"}`
	}
	return `{"success":true}`
}


func (c UsersController) DeleteUserRole(user_id int, role string) string {
	// TODO : LOG WHO ADDED WHO
	user := models.User{}
	if err := config.Global.DB.First(&user, user_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	role_index := -1
	for i,user_role := range user.Roles {
		if strings.Compare(role,user_role) == 0 {
			role_index = i
			break
		}
	}
	if role_index == -1 {
		return `{"success":false, "error":"user did not have the role"}`
	}
	append(user.Roles[:role_index], user.Roles[role_index+1:]...)
	if err := config.Global.DB.Save(&user).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"error saving new user's roles"}`
	}
	return `{"success":true}`
}

func GenerateToken(user *models.User) (string,error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
	})
	tokenString, err := token.SignedString([]byte(config.Global.SecretKey))
	return tokenString, err
}