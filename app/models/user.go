package models

import (
	_ "gopkg.in/pg.v5"
)

type User struct {
	ID 			int 		`json:"id,omitempty"`
	Achievements 		[]Achievement 	`pg:",many2many:users_achievements,joinFK:Achievement" json:"achievements,omitempty"`
	Device 			string		`json:"device,omitempty"`
	Friends 		[]User 		`pg:",many2many:users_friends,joinFK:Friend" json:"friends,omitempty"`
	NSFW			bool		`json:"nsfw,omitempty"`
	Username 		string		`json:"username,omitempty"`
	Passhash 		string		`json:"password,omitempty"`
	Premium			bool		`json:"premium,omitempty"`
	Roles			[]Role		`pg:",many2many:users_roles,joinFK:Role" json:"roles,omitempty"`
	TendressesReceived	[]*Tendresse	`pg:",fk:Receiver" json:"tendresses_received,omitempty" `
	TendressesSent		[]*Tendresse	`pg:",fk:Sender" json:"tendresses_sent,omitempty"`
}


