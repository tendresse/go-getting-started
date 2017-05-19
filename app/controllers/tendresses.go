// controllers/tendresses.go

package controllers

import (
	"encoding/json"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/dao"
	"github.com/tendresse/go-getting-started/app/models"

	log "github.com/Sirupsen/logrus"
	"github.com/graarh/golang-socketio"
)


type Tendresse struct {
}

func (c Tendresse) SendTendresse(friend_id int, current_user_id *int, so *gosocketio.Channel) string {
	users_dao        := dao.User{}
	tags_dao         := dao.Tag{}
	gifs_dao         := dao.Gif{}
	tendresses_dao   := dao.Tendresse{}
	achievements_dao := dao.Achievement{}

	// gestion de la tendresse
	friend := models.User{Id: friend_id}
	gif := models.Gif{}
	if err := gifs_dao.GetRandomGif(&gif); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while getting gif"}`
	}
	tendresse := models.Tendresse{SenderId: *current_user_id, ReceiverId: friend.Id, GifId: gif.Id}
	if err := tendresses_dao.CreateTendresse(&tendresse); err != nil {
		log.Error(err)
		return `{"success":false, "error":"tendresse not created"}`
	}
	log.Println("user_id:",*current_user_id,"sent a tendresse with gif_id:",gif.Id,"to user_id:",friend_id)

	// receiver achievements
	achievements := []models.Achievement{}
	for _,tag := range gif.Tags {
		// get Tag's achievements
		if err := tags_dao.GetFullTag(&tag); err != nil {
			log.Error(err)
			continue
		}
		for _,achievement := range tag.Achievements {
			ua := models.UsersAchievements{UserId: friend_id, AchievementId: achievement.Id}
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
					achievements = append(achievements,models.Achievement{Id: achievement.Id})
				}
			}
			if err := users_dao.UpdateAchievementToUser(&ua); err != nil {
				log.Error(err)
			}
		}
	}

	for _, achievement := range achievements {
		ach := models.Achievement{Id: achievement.Id}
		if err := achievements_dao.GetAchievement(&ach); err != nil {
			log.Error(err)
			continue
		}
		b, err := json.Marshal(ach)
		if err != nil {
			log.Error(err)
			continue
		}
		if err := users_dao.GetProfile(&friend); err != nil {
			log.Println(err)
			continue
		}
		log.Println(friend)
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
				log.Println(achievement)
				ua := models.UsersAchievements{UserId: friend_id, AchievementId: achievement.Id}
				if err := users_dao.GetOrCreateUserWithAchievement(&ua); err != nil {
					log.Error(err)
					continue
				}
				ua.Score++
				if ua.Score >= achievement.Condition {
					ua.Unlocked = true
					achievements = append(achievements,models.Achievement{Id: achievement.Id})
				}
				if err := config.Global.DB.Update(&ua); err != nil {
					log.Error(err)
				}
			}
		}
	}
	current_user := models.User{Id: *current_user_id}
	if err := users_dao.GetUser(&current_user); err != nil {
		log.Error(err)
	} else {
		for _, achievement := range achievements {
			ach := models.Achievement{Id: achievement.Id}
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
	tendressesDAO := dao.Tendresse{}
	tendresse := models.Tendresse{Id: tendresse_id}
	if err := tendressesDAO.GetTendresse(&tendresse); err != nil {
		log.Error(err)
		return `{"success":false, "error":"tendresse not found"}`
	}
	if tendresse.ReceiverId != *current_user_id{
		return `{"success":false, "error":"you are not the receiver of this tendresse"}`
	}
	tendresse.Viewed = true
	if err := tendressesDAO.UpdateTendresse(&tendresse); err != nil {
		log.Error(err)
		return `{"success":false, "error":"error while updating tendresse"}`
	}
	log.Println("user_id:",*current_user_id,"saw tendresse_id:",tendresse_id)
	return `{"success":true}`
}