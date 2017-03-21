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

func (c Achievement) CreateAchievements(achievements *[]models.Achievement) error {
	return c.DB.Insert(achievements)
}

func (c Achievement) UpdateAchievement(achievement *models.Achievement) error {
	return c.DB.Update(achievement)
}

func (c Achievement) GetAchievement(achievement *models.Achievement) error {
	return c.DB.Select(&achievement)
}

func (c Achievement) DeleteAchievement(achievement *models.Achievement) error {
	return c.DB.Delete(&achievement)
}

func (c Achievement) GetFullAchievement(achievement *models.Achievement) error {
	return c.DB.Model(&achievement).Column("achievement.*", "Tag").First()
}

func (c Achievement) GetAchievementByTitle(title string, achievement *models.Achievement) error {
	return c.DB.Model(&achievement).Where("title = ?",title).Select()
}

func (c Achievement) GetAchievements(achievements *[]models.Achievement) error {
	return c.DB.Model(&achievements).Select()
}

func (c Achievement) GetFullAchievements(achievements *[]models.Achievement) error {
	return c.DB.Model(&achievements).Column("achievement.*", "Tag").Select()
}

func (c Achievement) DeleteAchievements(achievements *[]models.Achievement) error {
	return c.DB.Delete(&achievements)
}

