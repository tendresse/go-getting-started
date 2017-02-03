package models

import (
	_ "gopkg.in/pg.v5"
)

type Tendresse struct {
	ID 		int 	`json:"id,omitempty"`
	SenderID 	int	`json:"sender_id,omitempty"`
	Sender 		*User	`json:"sender,omitempty"`
	ReceiverID 	int	`json:"receiver_id,omitempty"`
	Receiver 	*User	`json:"receiver,omitempty"`
	GifID 		int	`json:"gif_id,omitempty"`
	Gif 		*Gif	`json:"gif,omitempty"`
	Viewed		bool	`json:"viewed,omitempty"`
}
