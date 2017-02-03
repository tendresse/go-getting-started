package dao

import (
	// "encoding/json"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"

	_ "gopkg.in/pg.v5"
)

type Achievement struct{
}


func (c Achievement) GetAchievement(achievement *models.Achievement) error {
	return config.Global.DB.Select(&achievement)
}

func (c Achievement) GetFullAchievement(achievement *models.Achievement) error {
	return config.Global.DB.Model(&achievement).Column("achievement.*", "Tag").First()
}

func (c Achievement) GetAchievements(achievements *[]models.Achievement) error {
	return config.Global.DB.Model(&achievements).Select()
}

func (c Achievement) GetFullAchievements(achievements *[]models.Achievement) error {
	return config.Global.DB.Model(&achievements).Column("achievement.*", "Tag").Select()
}

func (c Achievement) GetAchievementByTitle(title string, achievement *models.Achievement) error {
	return config.Global.DB.Model(&achievement).Where("title = ?",title).Select()
}

func (c Achievement) CreateAchievement(achievement *models.Achievement) error {
	return config.Global.DB.Insert(achievement)
}


func (c Achievement) GetAchievementByAchievementname(Achievementname string, Achievement *models.Achievement) error {
	return config.Global.DB.Model(&Achievement).Where("Achievementname = ?",Achievementname).Select()
}


func (c Achievement) GetAchievementsByIds(ids []int, Achievements []*models.Achievement) error {
	return nil
}

func (c Achievement) GetAchievementFriends(Achievement *models.Achievement, friends []*models.Achievement) error {
	return nil
}


func (c Achievement) GetAchievementAchievements(id int, Achievement *models.Achievement) error {

	return nil
}
