package models

import (
	_ "github.com/go-pg/pg"
)

type User struct {
	Id           		int 		`json:"id,omitempty"`
	Achievements 		[]Achievement 	`pg:",many2many:users_achievements,joinFK:Achievement" json:"achievements,omitempty"`
	Device       		string		`json:"-"`
	Friends      		[]User 		`pg:",many2many:users_friends,joinFK:Friend" json:"friends,omitempty"`
	NSFW         		bool		`json:"nsfw,omitempty"`
	Username     		string		`json:"username,omitempty"`
	Passhash     		string		`json:"-"`
	Premium      		bool		`json:"premium,omitempty"`
	Roles		        []Role		`pg:",many2many:users_roles,joinFK:Role" json:"roles,omitempty"`
	TendressesReceived	[]*Tendresse	`pg:",fk:Receiver" json:"tendresses_received,omitempty" `
	TendressesSent		[]*Tendresse	`pg:",fk:Sender" json:"tendresses_sent,omitempty"`
	Token 			string 		`json:"-"`
}