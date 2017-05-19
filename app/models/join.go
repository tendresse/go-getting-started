package models

import (
	// "encoding/json"
	_"github.com/go-pg/pg"
)

type GifsTags struct {
	TagId int `sql:",pk"`
	GifId int `sql:",pk"`
}


type UsersAchievements struct {
	AchievementId int `sql:",pk"`
	UserId        int `sql:",pk"`
	Score         int
	Unlocked      bool
}

type UsersFriends struct {
	UserId   int `sql:",pk"`
	FriendId int `sql:",pk"`
}

type UsersRoles struct {
	RoleId int `sql:",pk"`
	UserId int `sql:",pk"`
}