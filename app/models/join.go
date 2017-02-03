package models

import (
	// "encoding/json"
	_"gopkg.in/pg.v5"
)
type GifsTags struct {
	TagID		int
	GifID		int
}


type UsersAchievements struct {
	AchievementID	int
	UserID 		int
	Score		int
}

type UsersFriends struct {
	UserID 		int
	FriendID	int
}

type UsersRoles struct {
	RoleID		int
	UserID		int
}