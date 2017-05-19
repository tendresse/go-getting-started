package dao

import (
	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"
)

type Tendresse struct {}


func (c Tendresse) CreateTendresse(tendresse *models.Tendresse) error {
	return config.Global.DB.Insert(&tendresse)
}
func (c Tendresse) CreateTendresses(tendresses []*models.Tendresse) error {
	for _,tendresse := range tendresses {
		if err := c.CreateTendresse(tendresse); err != nil {
			return err
		}
	}
	return nil
}

func (c Tendresse) UpdateTendresse(tendresse *models.Tendresse) error {
	return config.Global.DB.Update(&tendresse)
}


func (c Tendresse) GetTendresse(tendresse *models.Tendresse) error {
	return config.Global.DB.Model(&tendresse).First()
}
func (c Tendresse) GetTendresses(tendresses []*models.Tendresse) error {
	return config.Global.DB.Model(&tendresses).Select()
}


func (c Tendresse) GetFullTendresse(tendresse *models.Tendresse) error {
	return config.Global.DB.Model(&tendresse).
		Column("tendresse.*","Sender").
		Column("tendresse.*","Receiver").
		Column("tendresse.*","Gif").
		Where("tendresse.id = ?",tendresse.Id).
		First()
}
func (c Tendresse) GetFullTendresses(tendresses []*models.Tendresse) error {
	return config.Global.DB.Model(&tendresses).
		Column("tendresse.*","Sender").
		Column("tendresse.*","Receiver").
		Column("tendresse.*","Gif").
		Select()
}


func (c Tendresse) GetPendingTendresses(user *models.User) ([]models.Tendresse,error) {
	var ids[] int
	if _, err := config.Global.DB.Query(&ids, `SELECT t.id from tendresses t where t.receiver_id = ?`, user.Id); err != nil {
		return nil,err
	}
	tendresses := []models.Tendresse{}
	if len(ids) > 0 {
		for _, v := range ids {
			t := models.Tendresse{Id: v}
			if err := c.GetFullTendresse(&t); err != nil {
				return nil, err
			}
			tendresses = append(tendresses, t)
		}
	}
	return tendresses,nil
}

func (c Tendresse) CountSenderTendresses(sender_id int) (int,error) {
	return config.Global.DB.Model(&models.Tendresse{}).Where("sender_id = ?",sender_id).Count()
}


func (c Tendresse) DeleteTendresse(tendresse *models.Tendresse) error {
	return config.Global.DB.Delete(&tendresse)
}
func (c Tendresse) DeleteTendresses(tendresses []*models.Tendresse) error {
	for _,tendresse := range tendresses {
		if err := c.DeleteTendresse(tendresse); err != nil {
			return err
		}
	}
	return nil
}


