package models

import (
	// "encoding/json"
	_"gopkg.in/pg.v5"
)

type GifsTags struct {
	TagID		int `sql:",pk"`
	GifID		int `sql:",pk"`
}


type UsersAchievements struct {
	AchievementID	int `sql:",pk"`
	UserID 		int `sql:",pk"`
	Score		int
	Unlocked	bool
}

type UsersFriends struct {
	UserID 		int `sql:",pk"`
	FriendID	int `sql:",pk"`
}

type UsersRoles struct {
	RoleID		int `sql:",pk"`
	UserID		int `sql:",pk"`
}