// TENDRESSE controllers/achievements_controllers.go

package controllers

import (
	"encoding/json"
	"strings"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/models"

	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


type AchievementsController struct {
}

var achievements_dao dao.Achievement

// wrap with admin rights
func (c AchievementsController) GetAchievements() string {
	achievements := []models.Achievement{}
	if err := achievements_dao.GetFullAchievements(&achievements); err != nil {
		log.Error(err)
		return `{"success":false, "error":"no achievements found"}`
	}
	b, err := json.Marshal(achievements)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal achievements json error"}`
	}
	return strings.Join([]string{`{"success":true, "achievements":`,string(b),"}"}, "")
}

// wrap with admin rights
func (c AchievementsController) GetAchievement(achievement_id int) string {
	achievement := models.Achievement{}
	if err := achievements_dao.GetFullAchievement(&achievement); err != nil {
		log.Error(err)
		return `{"success":false, "error":"achievement not found"}`
	}
	b, err := json.Marshal(achievement)
	if err != nil {
		log.Error(err)
		return `{"success":false, "error":"marshal achievement json error"}`
	}
	return strings.Join([]string{`{"success":true, "achievement":`,string(b),"}"}, "")
}

// wrap with admin rights
func (c AchievementsController) CreateAchievement(achievement_json string) string {
	achievement := models.Achievement{}
	if err := json.Unmarshal([]byte(achievement_json), &achievement); err != nil{
		log.Error(err)
		return `{"success":false, "error":"marshal achievement json error"}`
	}
	if err := achievements_dao.GetAchievementByTitle(achievement.Title, &achievement); err != nil {
		for i,_ := range achievement.Tags {
			if err := config.Global.DB.Where(models.Tag{Name: achievement.Tags[i].Name}).FirstOrCreate(&achievement.Tags[i]).Error; err != nil {
				log.Error(err)
				return `{"success":false, "error":"error while fetching Tags"}`
			}
		}
		config.Global.DB.Create(&achievement)
		return strings.Join([]string{`{"success":true, "achievement":` , string(achievement.ID) , "}"} , "")
	}
	return `{"success":false, "error":"achievement already exists"}`
}

// wrap with admin rights
func (c AchievementsController) UpdateAchievement(achievement_json string) string {
	achievement := models.Achievement{}
	updated_achievement := models.Achievement{}
	if err := json.Unmarshal([]byte(achievement_json), &updated_achievement); err != nil{
		log.Error(err)
		return `{"success":false, "error":"marshal achievement json error"}`
	}
	// First is used without Preloading Tags in order to replace them
	if err := config.Global.DB.First(&achievement, updated_achievement.ID).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"achievement not found"}`
	}
	// we check if the new Title is not already taken
	if strings.Compare(updated_achievement.Title, achievement.Title) != 0 {
		if err := config.Global.DB.Where("Title = ?",updated_achievement.Title).First(&achievement).Error; err != nil {
			log.Error(err)
			return `{"success":false, "error":"achievement title already taken"}`
		}
	}
	for i,_ := range achievement.Tags {
		if err := config.Global.DB.Where(models.Tag{Name: achievement.Tags[i].Name}).FirstOrCreate(&achievement.Tags[i]).Error; err != nil {
			log.Error(err)
			return `{"success":false, "error":"erro while fetching Tags"}`
		}
	}
	config.Global.DB.Save(&achievement)
	config.Global.DB.Model(&achievement).Association("Tags").Replace(&updated_achievement.Tags)
	return strings.Join([]string{`{"success":true, "achievement":` , string(achievement.ID) , "}"} , "")
}

// wrap with admin rights
func (c AchievementsController) DeleteAchievement(achievement_id int) string {
	achievement := models.Achievement{}
	if err := config.Global.DB.Delete(&achievement, achievement_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"achievement already deleted"}`
	}
	return `{"success":true}`
}