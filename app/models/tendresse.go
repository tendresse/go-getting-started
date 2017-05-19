package models

import (
	_ "github.com/go-pg/pg"
)

type Tendresse struct {
	Id         int 		`json:"id,omitempty"`
	SenderId   int		`json:"sender_id,omitempty"`
	Sender     *User	`json:"sender,omitempty"`
	ReceiverId int		`json:"receiver_id,omitempty"`
	Receiver   *User	`json:"receiver,omitempty"`
	GifId      int		`json:"gif_id,omitempty"`
	Gif        *Gif		`json:"gif,omitempty"`
	Viewed     bool		`json:"viewed,omitempty"`
}
