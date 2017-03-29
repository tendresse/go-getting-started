package dao

import (
	// "encoding/json"
	"github.com/tendresse/go-getting-started/app/models"
	"gopkg.in/pg.v5"
	"fmt"
)

type Tendresse struct {
	DB 	*pg.DB
}


func (c Tendresse) CreateTendresse(tendresse *models.Tendresse) error {
	return c.DB.Insert(&tendresse)
}
func (c Tendresse) CreateTendresses(tendresses []*models.Tendresse) error {
	for _,tendresse := range tendresses {
		if err := c.CreateTendresse(tendresse); err != nil {
			return err
		}
	}
	return nil
}


func (c Tendresse) GetTendresse(tendresse *models.Tendresse) error {
	return c.DB.Model(&tendresse).First()
}
func (c Tendresse) GetTendresses(tendresses []*models.Tendresse) error {
	return c.DB.Model(&tendresses).Select()
}


func (c Tendresse) GetFullTendresse(tendresse *models.Tendresse) error {
	return c.DB.Model(&tendresse).
		Column("tendresse.*","Sender").
		Column("tendresse.*","Receiver").
		Column("tendresse.*","Gif").
		First()
}
func (c Tendresse) GetFullTendresses(tendresses []*models.Tendresse) error {
	return c.DB.Model(&tendresses).
		Column("tendresse.*","Sender").
		Column("tendresse.*","Receiver").
		Column("tendresse.*","Gif").
		Select()
}


func (c Tendresse) GetPendingTendresses(user *models.User) ([]models.Tendresse,error) {
	var ids[] int
	if _, err := c.DB.Query(&ids, `SELECT t.id from tendresses t where receiver_id = ?`, user.ID); err != nil {
		return nil,err
	}
	tendresses := []models.Tendresse{}
	for _,v := range ids {
		t := models.Tendresse{ID:v}
		if err := c.GetFullTendresse(&t); err != nil {
			return nil,err
		}
		fmt.Println(t)
		tendresses = append(tendresses,t)
	}
	return tendresses,nil
}


func (c Tendresse) DeleteTendresse(tendresse *models.Tendresse) error {
	return c.DB.Delete(&tendresse)
}
func (c Tendresse) DeleteTendresses(tendresses []*models.Tendresse) error {
	for _,tendresse := range tendresses {
		if err := c.DeleteTendresse(tendresse); err != nil {
			return err
		}
	}
	return nil
}