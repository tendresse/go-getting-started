// controllers/achievements_controllers.go

package controllers

import (
	"encoding/json"
	"strings"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/models"

	log "github.com/Sirupsen/logrus"
)


type AchievementsController struct {
}

// wrap with admin rights
func (c AchievementsController) GetAchievements() string {
	achievements_dao := dao.Achievement{DB:config.Global.DB}
	achievements := []models.Achievement{}
	if err := achievements_dao.GetAllAchievements(&achievements); err != nil {
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
	achievements_dao := dao.Achievement{DB:config.Global.DB}
	achievement := models.Achievement{ID:achievement_id}

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
	achievements_dao := dao.Achievement{DB:config.Global.DB}
	tags_dao         := dao.Tag{DB:config.Global.DB}

	achievement      := models.Achievement{}
	if err := json.Unmarshal([]byte(achievement_json), &achievement); err != nil{
		log.Error(err)
		return `{"success":false, "error":"marshal achievement json error"}`
	}
	if err := achievements_dao.GetAchievementByTitle(achievement.Title, &achievement); err != nil {
		log.Error(err)
		return `{"success":false, "error":"achievement already exists"}`
	}
	tag    := models.Tag{Title:achievement.Tag.Title}
	if err := tags_dao.GetOrCreateTag(&tag); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while fetching Tags"}`
	}
	achievement.TagID = tag.ID
	achievements_dao.CreateAchievement(&achievement)
	return strings.Join([]string{`{"success":true, "achievement":` , string(achievement.ID) , "}"} , "")
}

// wrap with admin rights
func (c AchievementsController) UpdateAchievement(achievement_json string) string {
	achievements_dao    := dao.Achievement{DB:config.Global.DB}
	tags_dao            := dao.Tag{DB:config.Global.DB}

	updated_achievement := models.Achievement{}
	if err := json.Unmarshal([]byte(achievement_json), &updated_achievement); err != nil{
		log.Error(err)
		return `{"success":false, "error":"marshal achievement json error"}`
	}

	achievement := models.Achievement{ID:updated_achievement.ID}
	if err := achievements_dao.GetAchievement(&achievement); err != nil {
		log.Error(err)
		return `{"success":false, "error":"achievement not found"}`
	}
	// if Title has changed, we check if the new Title is not already taken
	if strings.Compare(updated_achievement.Title, achievement.Title) != 0 {
		if err := achievements_dao.GetAchievementByTitle(updated_achievement.Title, nil); err != nil {
			log.Error(err)
			return `{"success":false, "error":"achievement title already taken"}`
		}
	}
	// if Tag has changed, we fetch or create the Tag
	if strings.Compare(updated_achievement.Tag.Title, achievement.Tag.Title) != 0 {
		tag := models.Tag{Title:achievement.Tag.Title}
		if err := tags_dao.GetOrCreateTag(&tag); err != nil {
			log.Error(err)
			return `{"success":false, "error":"error while fetching Tags"}`
		}
		updated_achievement.TagID = tag.ID
	}
	if err := achievements_dao.UpdateAchievement(&updated_achievement); err != nil{
		log.Error(err)
		return `{"success":false, "error":"error while updating Achievement"}`
	}
	return strings.Join([]string{`{"success":true, "achievement":` , string(achievement.ID) , "}"} , "")
}

// wrap with admin rights
func (c AchievementsController) DeleteAchievement(achievement_id int) string {
	achievements_dao := dao.Achievement{DB:config.Global.DB}
	if err := achievements_dao.DeleteAchievement(&models.Achievement{ID:achievement_id}).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"achievement already deleted"}`
	}
	return `{"success":true}`
}