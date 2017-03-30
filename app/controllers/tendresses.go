// controllers/tendresses_controllers.go

package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"

	_ "github.com/googollee/go-socket.io"
	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tendresse/go-getting-started/app/dao"
)


type TendressesController struct {
}

func (c TendressesController) SendTendresse(user_id string) string {
	users_dao        := dao.User{DB:config.Global.DB}
	gifs_dao         := dao.Gif{DB:config.Global.DB}
	tendresses_dao   := dao.Tendresse{DB:config.Global.DB}
	achievements_dao := dao.Achievement{DB:config.Global.DB}

	/*
 	 *  gestion tendresse
 	*/
	friend_id,err := strconv.Atoi(user_id)
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
	tendresse := models.Tendresse{SenderID:config.Global.CurrentUser.ID,ReceiverID:friend.ID,GifID:gif.ID}
	if err := tendresses_dao.CreateTendresse(&tendresse); err != nil {
		log.Error(err)
		return `{"success":false, "error":"tendresse not created"}`
	}

	//if err := friend.Notify(); err != nil {
	//	log.Error(err)
	//}

	// receiver achievements
	for _,tag := range gif.Tags {
		for _,achievement := range tag.Achievements {
			ua := models.UsersAchievements{UserID:friend_id,AchievementID:achievement.ID}
			if err := users_dao.GetOrCreateUserWithAchievement(&ua); err != nil {
				log.Error(err)
			}
			if ua.Score < achievement.Condition {
				ua.Score++
			}
			if ua.Unlocked != true {
				if ua.Score >= achievement.Condition {
					ua.Unlocked = true
				}
			}
			if err := config.Global.DB.Update(&ua); err != nil {
				log.Error(err)
			}
		}
	}

	// sender achievements
	sender_achievements := []models.Achievement{}

	if err := achievements_dao.GetSenderAchievementsWithCondition(&sender_achievements,)

	user_sender_achievements := []models.Achievement{}
	nb_tendresses_sent := config.Global.DB.Model(&config.Global.CurrentUser).Association("Friends").Count()
	if err := config.Global.DB.Model(&config.Global.CurrentUser).Select("id").Where("TypeOf = ?","send").Association("Achievements").Find(&user_sender_achievements).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"error getting sender achievements"}`
	}
	sender_achievements := []models.Achievement{}
	nb_tendresses_sent = config.Global.DB.Model(&config.Global.CurrentUser).Association("Friends").Count()
	if err := config.Global.DB.Select("id").Where("Condition <= ",nb_tendresses_sent).Where("TypeOf = ","send").Find(&sender_achievements).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"could not get current sender achievements"}`
	}
	new_achievements := []models.Achievement{}
	SenderAchievements:
	for _,send_achievement := range sender_achievements {
		for _,user_achievement := range user_sender_achievements {
			if user_achievement.ID == send_achievement.ID {
				continue SenderAchievements
			}
		}
		append(config.Global.CurrentUser.Achievements,send_achievement)
		append(new_achievements,send_achievement)
	}
	if err := config.Global.DB.Save(&config.Global.CurrentUser).Error; err != nil {
		log.Error(err)
	}
	b, err := json.Marshal(new_achievements)
	if err != nil {
		log.Error(err)
	}else{
		config.Global.Socket.BroadcastTo(config.Global.CurrentUser.Username,"new achievements", string(b))
	}


        // receiver achievements
	user_receiver_achievements := []models.Achievement{}
	if err := config.Global.DB.Model(&config.Global.CurrentUser).Select("id").Where("TypeOf = ?","receive").Association("Achievements").Find(&user_sender_achievements).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"error getting sender achievements"}`
	}
	receiver_achievements := []models.Achievement{}
	for _,tag := range gif.Tags {

	}
	if err := config.Global.DB.Select("id").Where("TypeOf = ","send").Find(&sender_achievements).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"could not get current sender achievements"}`
	}
	new_achievements := []models.Achievement{}
	ReceiverAchievements:
	for _,send_achievement := range sender_achievements {
		for _,user_achievement := range user_sender_achievements {
			if user_achievement.ID == send_achievement.ID {
				continue ReceiverAchievements
			}
		}
		append(config.Global.CurrentUser.Achievements,send_achievement)
		append(new_achievements,send_achievement)
	}
	if err := config.Global.DB.Save(&config.Global.CurrentUser).Error; err != nil {
		log.Error(err)
	}
	b, err := json.Marshal(new_achievements)
	if err != nil {
		log.Error(err)
	}else{
		config.Global.Socket.BroadcastTo(config.Global.CurrentUser.Username,"new achievements", string(b))
	}
	return `{"success":true}`
}


func (c TendressesController) Notify(username_friend string) string {
	return ""
}