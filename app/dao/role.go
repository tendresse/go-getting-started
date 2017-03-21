package dao

import (
	"github.com/tendresse/go-getting-started/app/models"
	"gopkg.in/pg.v5"
	//log "github.com/Sirupsen/logrus"
)

type Role struct {
	DB 	*pg.DB
}

func (c Role) CreateRole(role *models.Role) error {
	return c.DB.Insert(role)
}

func (c Role) UpdateRole(role *models.Role) error {
	return c.DB.Update(role)
}

func (c Role) GetRole(role *models.Role) error {
	return c.DB.Select(&role)
}

func (c Role) GetFullRole(role *models.Role) error {
	return c.DB.Model(&role).Column("role.*", "Users").First()
}

func (c Role) DeleteRole(role *models.Role) error {
	return c.DB.Delete(&role)
}

func (c Role) GetRoles(roles *[]models.Role) error {
	return c.DB.Model(&roles).Select()
}

func (c Role) GetFullRoles(roles *[]models.Role) error {
	return c.DB.Model(&roles).Column("role.*", "Users").Select()
}

func (c Role) DeleteRoles(roles *[]models.Role) error {
	return c.DB.Delete(&roles)
}