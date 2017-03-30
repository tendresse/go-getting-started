package dao

import (
	// "encoding/json"
	"github.com/tendresse/go-getting-started/app/models"
	"gopkg.in/pg.v5"
)

type User struct {
	DB 	*pg.DB
}


func (c User) CreateUser(user *models.User) error {
	return c.DB.Insert(user)
}
func (c User) CreateUsers(users []*models.User) error {
	for _,user := range users {
		if err := c.CreateUser(user); err != nil {
			return err
		}
	}
	return nil
}


func (c User) GetUser(user *models.User) error {
	return c.DB.Select(&user)
}
func (c User) GetAllUsers(users *[]models.User) error {
	count, err := c.DB.Model(&models.User{}).Count()
	if err != nil {
		return err
	}
	return c.DB.Model(&users).Limit(count).Select()
}


func (c User) UpdateUser(user *models.User) error {
	return c.DB.Update(&user)
}


func (c User) DeleteUser(user *models.User) error {
	return c.DB.Delete(&user)
}


func (c User) GetFullUser(user *models.User) error {
	return c.DB.Model(&user).Column("user.*", "TendressesSent").Column("user.*", "TendressesReceived").Column("user.*", "Roles").Column("user.*", "Friends").Column("user.*", "Achievements").First()
}


func (c User) GetUserByUsername(user *models.User, username string) error {
	return c.DB.Model(&user).Where("username = ?",username).Select()
}


func (c User) AddRoleToUser(role *models.Role, user *models.User) error {
	return c.DB.Insert(&models.UsersRoles{RoleID:role.ID, UserID:user.ID})
}
func (c User) AddRolesToUser(roles []*models.Role, user *models.User) error {
	for _,role := range roles {
		if err := c.AddRoleToUser(role, user); err != nil {
			return err
		}
	}
	return nil
}


func (c User) DeleteRoleFromUser(role *models.Role, user *models.User) error {
	return c.DB.Delete(&models.UsersRoles{RoleID:role.ID, UserID:user.ID})
}
func (c User) DeleteRolesFromUser(roles []*models.Role, user *models.User) error {
	for _,role := range roles {
		if err := c.DeleteRoleFromUser(role, user); err != nil {
			return err
		}
	}
	return nil
}


func (c User) AddFriendToUser(friend *models.User, user *models.User) error {
	return c.DB.Insert(&models.UsersFriends{UserID:user.ID, FriendID:friend.ID})
}
func (c User) AddFriendsToUser(friends []*models.User, user *models.User) error {
	for _,friend := range friends {
		if err := c.AddFriendToUser(friend, user); err != nil {
			return err
		}
	}
	return nil
}


func (c User) DeleteFriendFromUser(friend *models.User, user *models.User) error {
	return c.DB.Delete(&models.UsersFriends{UserID:user.ID, FriendID:friend.ID})
}
func (c User) DeleteFriendsFromUser(friends []*models.User, user *models.User) error {
	for _,friend := range friends {
		if err := c.DeleteFriendFromUser(friend, user); err != nil {
			return err
		}
	}
	return nil
}

func (c User) GetOrCreateUserWithAchievement(r *models.UsersAchievements) error {
	if err := c.DB.Model(&r).First(); err != nil {
		return c.DB.Insert(&r)
	}
	return nil
}

func (c User) AddAchievementToUser(achievement *models.Achievement, user *models.User) error {
	return c.DB.Insert(&models.UsersAchievements{AchievementID:achievement.ID, UserID:user.ID})
}


func (c User) DeleteAchievementFromUser(achievement *models.Achievement, user *models.User) error {
	return c.DB.Delete(&models.UsersAchievements{AchievementID:achievement.ID, UserID:user.ID})
}
