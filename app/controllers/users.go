// controllers/users_controllers.go

package controllers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/models"

	"golang.org/x/crypto/bcrypt"
	"github.com/graarh/golang-socketio"
	"github.com/dgrijalva/jwt-go"
	log "github.com/Sirupsen/logrus"
)

type User struct {
}

type MyCustomClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}


func (c User) IsAuthentified(current_user_id int) bool {
	return current_user_id != 0
}

func (c User) HasAtLeastOneOfTheseRoles(roles []string, current_user_id *int) bool {
	userDAO := dao.User{}
	user := models.User{Id : *current_user_id}
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
	userDAO := dao.User{}
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
		*current_user_id  = user.Id
		token,err := GenerateToken(user.Id)
		if err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while creating token"}`
		}
		user.Token = token
		if err = userDAO.UpdateTokenUser(&user); err != nil {
			log.Error("error while updating User's token :",token,"with error :",err)
		}
		log.Println("user",username,"with id:",*current_user_id,"created an account.")
		so.Join(user.Username)
		return strings.Join([]string{`{"success":true, "token":`, token, "}"}, "")
	}
	return `{"success":false, "error":"username already taken"}`
}

func (c User) Login(username string, password string, current_user_id *int, so *gosocketio.Channel) string {
	userDAO := dao.User{}
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
	*current_user_id = user.Id
	token,err := GenerateToken(user.Id)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while creating token"}`
	}
	log.Println("user :",username,"with id =",*current_user_id,"connected.")
	so.Join(user.Username)
	user.Token = token
	if err = userDAO.UpdateTokenUser(&user); err != nil {
		log.Error("error while updating User's token :",token,"with error :",err)
	}
	log.Println("user",username,"with id:",*current_user_id,"just logged in.")
	return strings.Join([]string{`{"success":true, "token":`, token, "}"}, "")
}

func (c User) LoginWithToken(user_token string, current_user_id *int, so *gosocketio.Channel) string {
	userDAO := dao.User{}
	*current_user_id = 0

	the_token, err := jwt.ParseWithClaims(user_token, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Global.SecretKey), nil
	})
	if the_token.Valid {
		if claims, ok := the_token.Claims.(*MyCustomClaims); ok {
			user := models.User{Id: claims.UserID}
			if err := userDAO.GetUser(&user); err != nil {
				log.Error(err)
				return `{"success":false, "error":"token's user not found"}`
			}
			// check if token is in user's DB
			if strings.Compare(user_token,user.Token) != 0 {
				log.Error("token for user",user.Username,"with Id =",claims.UserID,"is blacklisted :",user_token)
				return `{"success":false, "error":"invalid token"}`
			}
			*current_user_id = user.Id
			log.Println("user:",user.Username,"with id:",*current_user_id,"juste logged in with token.")
			so.Join(user.Username)
			return `{"success":true}`
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			log.Error("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			log.Error("Timing is everything")
		} else {
			log.Error("Couldn't handle this token :", err)
		}
	} else {
		log.Error("Couldn't handle this token :", err)
	}
	log.Error("claims is not valid or token is not valid")
	log.Error(err)
	return `{"success":false, "error":"invalid token"}`
}

func (c User) GetUser(user_id int) string {
	userDAO := dao.User{}
	user := models.User{Id: user_id}
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
	userDAO := dao.User{}
	user := models.User{Id: user_id}
	if err := userDAO.GetProfile(&user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	userJSON := models.User{
		Id:           user.Id,
		Username:     user.Username,
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
	current_user := models.User{Id: *current_user_id}
	tendressesDAO := dao.Tendresse{}
	tendresses, err := tendressesDAO.GetPendingTendresses(&current_user)
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
	usersDAO := dao.User{}
	current_user := models.User{Id: *current_user_id}
	current_user.Device = device_token
	if err := usersDAO.UpdateDeviceUser(&current_user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while updating user"}`
	}
	return `{"success":true}`
}

func (c User) AddFriend(username string, current_user_id *int) string {
	userDAO := dao.User{}
	user := models.User{}
	current_user := models.User{Id: *current_user_id}
	if err := userDAO.GetFullUser(&current_user); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while getting current user"}`
	}
	username = strings.ToLower(username)
	if err := userDAO.GetUserByUsername(&user,username); err != nil {
		log.Error(err)
		return `{"success":false, "error":"unknown user"}`
	}
	if err := userDAO.AddFriendToUser(&user, &models.User{Id: *current_user_id}); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error adding friend"}`
	}
	log.Println("user:",current_user.Username,"added user:",username,"as a friend.")
	return `{"success":true}`
}

func (c User) DeleteFriend(friend_id int, current_user_id *int) string {
	userDAO := dao.User{}
	if err := userDAO.DeleteFriendFromUser(&models.User{Id: friend_id}, &models.User{Id: *current_user_id}); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error adding friend"}`
	}
	return `{"success":true}`
}

func (c User) GrantUserRole(user_id int, title string) string {
	userDAO := dao.User{}
	rolesDAO := dao.Role{}
	user := models.User{Id: user_id}
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
	userDAO := dao.User{}
	rolesDAO := dao.Role{}
	user := models.User{Id: user_id}
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
	claims := MyCustomClaims{
		user_id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 7 * 2).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Global.SecretKey))
	return tokenString, err
}