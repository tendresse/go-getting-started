package dao

import (
	// "encoding/json"

	"github.com/tendresse/go-getting-started/app/models"

	"gopkg.in/pg.v5"
)


type Achievement struct{
	DB 	*pg.DB
}


func (c Achievement) CreateAchievement(achievement *models.Achievement) error {
	return c.DB.Insert(achievement)
}

func (c Achievement) CreateAchievements(achievements []*models.Achievement) error {
	for _,achievement := range achievements {
		if err := c.CreateAchievement(achievement); err != nil {
			return err
		}
	}
	return nil
}

func (c Achievement) UpdateAchievement(achievement *models.Achievement) error {
	return c.DB.Update(achievement)
}


func (c Achievement) GetAchievement(achievement *models.Achievement) error {
	return c.DB.Select(&achievement)
}
func (c Achievement) GetAchievements(achievements []*models.Achievement) error {
	return c.DB.Model(&achievements).Select()
}


func (c Achievement) GetAllAchievements(achievements *[]models.Achievement) error {
	count, err := c.DB.Model(&models.Achievement{}).Count()
	if err != nil {
		return err
	}
	return c.DB.Model(&achievements).Limit(count).Select()
}

func (c Achievement) GetFullAchievement(achievement *models.Achievement) error {
	return c.DB.Model(&achievement).Column("achievement.*", "Tag").First()
}
func (c Achievement) GetFullAchievements(achievements []*models.Achievement) error {
	return c.DB.Model(&achievements).Column("achievement.*", "Tag").Select()
}


func (c Achievement) GetAchievementByTitle(title string, achievement *models.Achievement) error {
	return c.DB.Model(&achievement).Where("title = ?",title).Select()
}


func (c Achievement) GetOrCreateAchievement(achievement *models.Achievement) error {
	return c.DB.Select(&achievement)
}


func (c Achievement) DeleteAchievement(achievement *models.Achievement) error {
	return c.DB.Delete(&achievement)
}
func (c Achievement) DeleteAchievements(achievements []*models.Achievement) error {
	for _,achievement := range achievements {
		if err := c.DeleteAchievement(achievement); err != nil {
			return err
		}
	}
	return nil
}

func (c Achievement) GetSenderAchievementsWithCondition(achievements *[]models.Achievement, condition int) error {
	return c.DB.Model(achievements).Where("type_of = ?","sender").Where("condition > ?",condition).Select()
}