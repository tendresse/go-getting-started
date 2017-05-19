package dao

import (
	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"
)

type User struct {}


func (c User) CreateUser(user *models.User) error {
	return config.Global.DB.Insert(user)
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
	return config.Global.DB.Select(&user)
}
func (c User) GetAllUsers(users *[]models.User) error {
	count, err := config.Global.DB.Model(&models.User{}).Count()
	if err != nil {
		return err
	}
	return config.Global.DB.Model(&users).Limit(count).Select()
}


func (c User) UpdateUser(user *models.User) error {
	return config.Global.DB.Update(&user)
}
func (c User) UpdateDeviceUser(user *models.User) error {
	_, err :=  config.Global.DB.Model(&user).Column("device").Returning("*").Update()
	return err
}
func (c User) UpdateTokenUser(user *models.User) error {
	_, err :=  config.Global.DB.Model(&user).Column("token").Returning("*").Update()
	return err
}


func (c User) DeleteUser(user *models.User) error {
	return config.Global.DB.Delete(&user)
}


func (c User) GetFullUser(user *models.User) error {
	return config.Global.DB.Model(&user).Column("user.*", "TendressesSent").Column("user.*", "TendressesReceived").Column("user.*", "Roles").Column("user.*", "Friends").Column("user.*", "Achievements").Where("id = ?",user.Id).First()
}

func (c User) GetProfile(user *models.User) error {
	return config.Global.DB.Model(&user).Column("user.*", "Achievements").Where("id = ?",user.Id).First()

}


func (c User) GetUserByUsername(user *models.User, username string) error {
	return config.Global.DB.Model(&user).Where("username = ?",username).Select()
}


func (c User) AddRoleToUser(role *models.Role, user *models.User) error {
	return config.Global.DB.Insert(&models.UsersRoles{RoleId: role.Id, UserId: user.Id})
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
	return config.Global.DB.Delete(&models.UsersRoles{RoleId: role.Id, UserId: user.Id})
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
	return config.Global.DB.Insert(&models.UsersFriends{UserId: user.Id, FriendId: friend.Id})
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
	return config.Global.DB.Delete(&models.UsersFriends{UserId: user.Id, FriendId: friend.Id})
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
	if err := config.Global.DB.Select(&r); err != nil {
		return config.Global.DB.Insert(&r)
	}
	return nil
}

func (c User) AddAchievementToUser(achievement *models.Achievement, user *models.User) error {
	return config.Global.DB.Insert(&models.UsersAchievements{AchievementId: achievement.Id, UserId: user.Id, Score:1})
}

func (c User) UpdateAchievementToUser(r *models.UsersAchievements) error {
	return config.Global.DB.Update(&r)
}


func (c User) DeleteAchievementFromUser(achievement *models.Achievement, user *models.User) error {
	return config.Global.DB.Delete(&models.UsersAchievements{AchievementId: achievement.Id, UserId: user.Id})
}
