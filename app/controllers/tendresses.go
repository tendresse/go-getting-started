// controllers/tendresses.go

package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/models"

	log "github.com/Sirupsen/logrus"
	"github.com/graarh/golang-socketio"
)


type Tendresse struct {
}

func (c Tendresse) SendTendresse(id string, current_user_id *int, so *gosocketio.Channel) string {
	users_dao        := dao.User{DB:config.Global.DB}
	gifs_dao         := dao.Gif{DB:config.Global.DB}
	tendresses_dao   := dao.Tendresse{DB:config.Global.DB}
	achievements_dao := dao.Achievement{DB:config.Global.DB}

	// gestion de la tendresse
	friend_id,err := strconv.Atoi(id)
	if err != nil {
		return `{"success":false, "error":"error while parsing friend id"}`
	}
	friend := models.User{ID:friend_id}
	if err := users_dao.GetUser(&friend).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"friend not found"}`
	}
	gif := models.Gif{}
	if err := gifs_dao.GetRandomGif(&gif); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while getting gif"}`
	}
	tendresse := models.Tendresse{SenderID:*current_user_id,ReceiverID:friend.ID,GifID:gif.ID}
	if err := tendresses_dao.CreateTendresse(&tendresse); err != nil {
		log.Error(err)
		return `{"success":false, "error":"tendresse not created"}`
	}
	log.Println("user_id ",*current_user_id," sent a tendresse with gif_id = ",gif.ID," to user_id = ",friend_id)

	// receiver achievements
	achievements := []models.Achievement{}
	for _,tag := range gif.Tags {
		for _,achievement := range tag.Achievements {
			ua := models.UsersAchievements{UserID:friend_id,AchievementID:achievement.ID}
			if err := users_dao.GetOrCreateUserWithAchievement(&ua); err != nil {
				log.Error(err)
				continue
			}
			if ua.Score < achievement.Condition {
				ua.Score++
			}
			if ua.Unlocked != true {
				if ua.Score >= achievement.Condition {
					ua.Unlocked = true
					achievements = append(achievements,models.Achievement{ID:achievement.ID})
				}
			}
			if err := config.Global.DB.Update(&ua); err != nil {
				log.Error(err)
			}
		}
	}

	for _, achievement := range achievements {
		ach := models.Achievement{ID: achievement.ID}
		if err := achievements_dao.GetAchievement(&ach); err != nil {
			log.Error(err)
			continue
		}
		b, err := json.Marshal(ach)
		if err != nil {
			log.Error(err)
			continue
		}
		so.BroadcastTo(friend.Username, "new achievement", string(b))
	}

	// sender achievements
	sender_achievements := []models.Achievement{}
	achievements = []models.Achievement{}
	count,err := tendresses_dao.CountSenderTendresses(*current_user_id)
	if err != nil {
		log.Error(err)
	}else {
		if err := achievements_dao.GetSenderAchievementsWithCondition(&sender_achievements,count); err != nil {
			log.Error(err)
		} else {
			for _,achievement := range sender_achievements {
				ua := models.UsersAchievements{UserID:friend_id,AchievementID:achievement.ID}
				if err := users_dao.GetOrCreateUserWithAchievement(&ua); err != nil {
					log.Error(err)
					continue
				}
				ua.Score++
				if ua.Score >= achievement.Condition {
					ua.Unlocked = true
					achievements = append(achievements,models.Achievement{ID:achievement.ID})
				}
				if err := config.Global.DB.Update(&ua); err != nil {
					log.Error(err)
				}
			}
		}
	}
	current_user := models.User{ID:*current_user_id}
	if err := users_dao.GetUser(&current_user); err != nil {
		log.Error(err)
	} else {
		for _, achievement := range achievements {
			ach := models.Achievement{ID: achievement.ID}
			if err := achievements_dao.GetAchievement(&ach); err != nil {
				log.Error(err)
				continue
			}
			b, err := json.Marshal(ach)
			if err != nil {
				log.Error(err)
				continue
			}
			so.BroadcastTo(current_user.Username, "new achievement", string(b))
		}
	}
	return `{"success":true}`
}

func (c Tendresse) SetTendresseAsSeen(tendresse_id int, current_user_id *int) string {
	tendressesDAO := dao.Tendresse{DB:config.Global.DB}
	tendresse := models.Tendresse{ID:tendresse_id}
	if err := tendressesDAO.GetTendresse(&tendresse); err != nil {
		log.Error(err)
		return `{"success":false, "error":"tendresse not found"}`
	}
	if tendresse.ReceiverID != *current_user_id{
		return `{"success":false, "error":"you are not the receiver of this tendresse"}`
	}
	tendresse.Viewed = true
	if err := tendressesDAO.UpdateTendresse(&tendresse); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while updating tendresse"}`
	}
	log.Println("user_id = ",*current_user_id," saw tendresse_id = ",tendresse_id)
	return `{"success":true}`
}