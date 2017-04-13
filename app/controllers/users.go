// controllers/users_controllers.go

package controllers

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/models"

	"golang.org/x/crypto/bcrypt"
	"github.com/graarh/golang-socketio"
	"github.com/dgrijalva/jwt-go"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/pg.v5"
)

type User struct {
	DB 	*pg.DB
}


func (c User) IsAuthentified(current_user_id int) bool {
	return current_user_id != 0
}

func (c User) HasAtLeastOneOfTheseRoles(roles []string, current_user_id *int) bool {
	userDAO := dao.User{DB:config.Global.DB}
	user := models.User{ID:*current_user_id}
	if err := userDAO.GetFullUser(&user); err != nil {
		return false
	}
	for _,role := range roles {
		for _,user_role := range user.Roles {
			if strings.Compare(role,user_role.Title) == 0 {
				return true
			}
		}
	}
	return false
}

func (c User) Signup(username string, password string, current_user_id *int, so *gosocketio.Channel) string {
	userDAO := dao.User{DB:config.Global.DB}
	*current_user_id = 0
	user := models.User{}
	username = strings.ToLower(username)
	if err := userDAO.GetUserByUsername(&user,username); err != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
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
		*current_user_id  = user.ID
		token,err := GenerateToken(user.ID)
		if err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while creating token"}`
		}
		log.Println("user : ",username," with id = ",*current_user_id," created an account.")
		so.Join(user.Username)
		strings.Join([]string{`{"success":true, "token":`, token, "}"}, "")
	}
	return `{"success":false, "error":"username already taken"}`
}

func (c User) Login(username string, password string, current_user_id *int, so *gosocketio.Channel) string {
	userDAO := dao.User{DB:config.Global.DB}
	*current_user_id = 0
	user := models.User{}
	username = strings.ToLower(username)
	if err := userDAO.GetUserByUsername(&user,username); err != nil {
		log.Error(err)
		return `{"success":false, "error":"username or password incorrect"}`
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Passhash), []byte(password) ); err != nil {
		log.Error(err)
		return `{"success":false, "error":"username or password incorrect"}`
	}
	*current_user_id = user.ID
	token,err := GenerateToken(user.ID)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while creating token"}`
	}
	log.Println("user : ",username," with id = ",*current_user_id," connected.")
	so.Join(user.Username)
	return strings.Join([]string{`{"success":true, "token":`, token, "}"}, "")
}

func (c User) LoginWithToken(user_token string, current_user_id *int, so *gosocketio.Channel) string {
	userDAO := dao.User{DB:config.Global.DB}
	*current_user_id = 0
	token, err := jwt.Parse(user_token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(token.Header["alg"].(string))
		}
		return []byte(config.Global.SecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := models.User{ID: claims["id"].(int)}
		if err := userDAO.GetUser(&user); err != nil {
			log.Error(err)
			return `{"success":false, "error":"invalid token"}`
		}
		*current_user_id = user.ID
		log.Println("user : ",user.Username," with id = ",*current_user_id," connected with token.")
		so.Join(user.Username)
		return `{"success":true}`
	}
	log.Error(err)
	return `{"success":false, "error":"invalid token"}`
}

func (c User) GetUser(user_id int) string {
	userDAO := dao.User{DB:config.Global.DB}
	user := models.User{ID:user_id}
	if err := userDAO.GetFullUser(&user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	b, err := json.Marshal(user)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal user json error"}`
	}
	return strings.Join([]string{`{"success":true, "user":`,string(b),"}"} ,"")
}

func (c User) GetProfile(user_id int) string {
	userDAO := dao.User{DB:config.Global.DB}
	user := models.User{ID:user_id}
	if err := userDAO.GetProfile(&user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	userJSON := models.User{
		ID: user.ID,
		Username: user.Username,
		Achievements: user.Achievements,
	}
	b, err := json.Marshal(userJSON)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while marshaling user's json"}`
	}
	return strings.Join([]string{`{"success":true, "user":`,string(b),"}"} ,"")
}

func (c User) GetPendingTendresses(current_user_id *int) string {
	current_user := models.User{ID:*current_user_id}
	tendressesDAO := dao.Tendresse{DB:config.Global.DB}
	tendresses,err := tendressesDAO.GetPendingTendresses(&current_user)
	if err != nil {
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

func (c User) UpdateDevice(device_token string, current_user_id *int) string {
	usersDAO := dao.User{DB:config.Global.DB}
	current_user := models.User{ID:*current_user_id}
	current_user.Device = device_token
	if err := usersDAO.UpdateUser(&current_user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while updating user"}`
	}
	return `{"success":true}`
}

func (c User) AddFriend(username string, current_user_id *int) string {
	userDAO := dao.User{DB:config.Global.DB}
	user := models.User{}
	current_user := models.User{ID:*current_user_id}
	if err := userDAO.GetFullUser(&current_user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while getting current user"}`
	}
	username = strings.ToLower(username)
	if err := userDAO.GetUserByUsername(&user,username); err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	if err := userDAO.AddFriendToUser(&user, &models.User{ID:*current_user_id}); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error adding friend"}`
	}
	return `{"success":true}`
}

func (c User) DeleteFriend(friend_id int, current_user_id *int) string {
	userDAO := dao.User{DB:config.Global.DB}
	if err := userDAO.DeleteFriendFromUser(&models.User{ID:friend_id}, &models.User{ID:*current_user_id}); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error adding friend"}`
	}
	return `{"success":true}`
}

func (c User) GrantUserRole(user_id int, title string) string {
	userDAO := dao.User{DB:config.Global.DB}
	rolesDAO := dao.Role{DB:config.Global.DB}
	user := models.User{ID:user_id}
	if err := userDAO.GetFullUser(&user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	for _,user_role := range user.Roles {
		if strings.Compare(title,user_role.Title) == 0 {
			return `{"success":false, "error":"user already has this role"}`
		}
	}
	role := models.Role{Title:title}
	if err := rolesDAO.GetOrCreateRole(&role); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while fetching/creating role"}`
	}
	if err := userDAO.AddRoleToUser(&role, &user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while adding role to user"}`
	}
	return `{"success":true}`
}


func (c User) DeleteUserRole(user_id int, title string) string {
	userDAO := dao.User{DB:config.Global.DB}
	rolesDAO := dao.Role{DB:config.Global.DB}
	user := models.User{ID:user_id}
	if err := userDAO.GetFullUser(&user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	found := false
	for _,user_role := range user.Roles {
		if strings.Compare(title,user_role.Title) == 0 {
			found = true
			break
		}
	}
	if found == false {
		role := models.Role{Title:title}
		if err := rolesDAO.GetOrCreateRole(&role); err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while fetching/creating role"}`
		}
		if err := userDAO.DeleteRoleFromUser(&role, &user); err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while adding role to user"}`
		}
	}
	return `{"success":true}`
}

func GenerateToken(user_id int) (string,error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user_id,
	})
	tokenString, err := token.SignedString([]byte(config.Global.SecretKey))
	return tokenString, err
}