package models

import (
    	// "encoding/json"
	_"gopkg.in/pg.v5"
)

type Achievement struct{
	ID 		int 	`json:"id,omitempty"`
	Condition   	int 	`json:"condition,omitempty"`
	Icon        	string 	`json:"icon,omitempty"`
	TagID		int	`json:"tag_id,omitempty"`
	Tag        	*Tag    `json:"tag,omitempty"`
	Title       	string 	`json:"title,omitempty"`
	TypeOf      	string 	`json:"type_of,omitempty"`
	XP          	int 	`json:"xp,omitempty"`
}
