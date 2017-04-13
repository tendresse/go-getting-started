package dao

import (
	"github.com/tendresse/go-getting-started/app/models"

	"gopkg.in/pg.v5"
)

type Role struct {
	DB 	*pg.DB
}


func (c Role) CreateRole(role *models.Role) error {
	return c.DB.Insert(role)
}
func (c Role) CreateRoles(roles []*models.Role) error {
	for _,role := range roles {
		if err := c.CreateRole(role); err != nil {
			return err
		}
	}
	return nil
}

func (c Role) GetRole(role *models.Role) error {
	return c.DB.Select(&role)
}
func (c Role) GetRoles(roles []*models.Role) error {
	return c.DB.Model(&roles).Select()
}


func (c Role) GetOrCreateRole(role *models.Role) error {
	_,err := c.DB.Model(&role).
	Column("id").
	Where("title = ?", role.Title).
	OnConflict("DO NOTHING"). // OnConflict is optional
	Returning("id").
	SelectOrInsert()
	return err
}


func (c Role) GetAllRoles(roles *[]models.Role) (error) {
	count, err := c.DB.Model(&models.Role{}).Count()
	if err != nil {
		return err
	}
	return c.DB.Model(&roles).Limit(count).Select()
}


func (c Role) UpdateRole(role *models.Role) error {
	return c.DB.Update(role)
}


func (c Role) GetFullRole(role *models.Role) error {
	return c.DB.Model(&role).Column("role.*", "Users").First()
}
func (c Role) GetFullRoles(roles []*models.Role) error {
	return c.DB.Model(&roles).Column("role.*", "Users").Select()
}


func (c Role) DeleteRole(role *models.Role) error {
	return c.DB.Delete(&role)
}
func (c Role) DeleteRoles(roles []*models.Role) error {
	for _,role := range roles {
		if err := c.DeleteRole(role); err != nil {
			return err
		}
	}
	return nil
}