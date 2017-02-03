package dao

import (
    // "encoding/json"

    //"github.com/tendresse/go-getting-started/app/config"
    "github.com/tendresse/go-getting-started/app/models"
    //
    //"gopkg.in/pg.v5"
    //log "github.com/Sirupsen/logrus"
)

type Tendresse struct {
}

func (c Tendresse) CreateTendresse(tendresse *models.Tendresse) error {

    return nil
}

func (c Tendresse) GetTendresse(id int64, tendresse *models.Tendresse) error {

    return nil
}

func (c Tendresse) GetTendresses(tendresses []*models.Tendresse) error {

    return nil
}

func (c Tendresse) GetTendressesByIds(ids []int64, tendresses []*models.Tendresse) (error) {

    return nil
}
