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

func (c User) GetUser(user *models.User) error {
	return c.DB.Select(&user)
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

func (c User) GetUsers(users *[]models.User) error {
	return c.DB.Model(&users).Select()
}

func (c User) GetFullUsers(users *[]models.User) error {
	return c.DB.Model(&users).Column("user.*", "Friends").Select()
}

func (c User) DeleteUsers(users *[]models.User) error {
	return c.DB.Delete(&users)
}

func (c User) GetUserByUsername(user *models.User, username string) error {
	return c.DB.Model(&user).Where("username = ?",username).Select()
}

func (c User) AddRoleToUser(role *models.Role, user *models.User) error {
	return c.DB.Insert(&models.UsersRoles{RoleID:role.ID, UserID:user.ID})
}

func (c User) DeleteRoleFromUser(role *models.Role, user *models.User) error {
	return c.DB.Delete(&models.UsersRoles{RoleID:role.ID, UserID:user.ID})
}

func (c User) AddRolesToUser(roles *[]models.Role, user *models.User) error {
	for _,role := range *roles {
		if err := c.AddRoleToUser(&role, user); err != nil {
			return err
		}
	}
	return nil
}

func (c User) DeleteRolesFromUser(roles *[]models.Role, user *models.User) error {
	for _,role := range *roles {
		if err := c.DeleteRoleFromUser(&role, user); err != nil {
			return err
		}
	}
	return nil
}

func (c User) AddFriendToUser(friend *models.User, user *models.User) error {
	return c.DB.Insert(&models.UsersFriends{UserID:user.ID, FriendID:friend.ID})
}

func (c User) DeleteFriendFromUser(friend *models.User, user *models.User) error {
	return c.DB.Delete(&models.UsersFriends{UserID:user.ID, FriendID:friend.ID})
}

func (c User) AddFriendsToUser(friends *[]models.User, user *models.User) error {
	for _,friend := range *friends {
		if err := c.AddFriendToUser(&friend, user); err != nil {
			return err
		}
	}
	return nil
}

func (c User) DeleteFriendsFromUser(friends *[]models.User, user *models.User) error {
	for _,friend := range *friends {
		if err := c.DeleteFriendFromUser(&friend, user); err != nil {
			return err
		}
	}
	return nil
}

func (c User) AddAchievementToUser(achievement *models.Achievement, user *models.User) error {
	return c.DB.Insert(&models.UsersAchievements{AchievementID:achievement.ID, UserID:user.ID})
}

func (c User) AddAchievementsToUser(achievements *[]models.User, user *models.User) error {
	for _,achievement := range *achievements {
		if err := c.DeleteFriendFromUser(&achievement, user); err != nil {
			return err
		}
	}
	return nil
}