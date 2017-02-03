// controllers/tendresses_controllers.go

package controllers

import (
	"encoding/json"
	"strings"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/models"

	_ "github.com/googollee/go-socket.io"
	log "github.com/Sirupsen/logrus"
	_ "github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os/user"
)


type TendressesController struct {
}

func (c TendressesController) SendTendresse(user_id string) string {
	friend := models.User{}
	if err := config.Global.DB.First(&friend, user_id).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"friend not found"}`
	}
	gif := GifsController{}.RandomGif()
	tendresse := models.Tendresse{SenderID:config.Global.CurrentUser.ID,ReceivedID:friend.ID,GifID:gif.ID}
	if err := config.Global.DB.Create(&tendresse).Error; err != nil {
		log.Error(err)
		return `{"success":false, "error":"tendresse not created"}`
	}

	//if err := friend.Notify(); err != nil {
	//	log.Error(err)
	//}

	// sender achivements
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

//authorized
so.On("tendresse seen", func(tendresse_id int){
    tendresse = Tendresse.query.get(tendresse_id)
    if tendresse is not None:
        if tendresse.receiver is current_user:
            tendresse.state_viewed = true
            db.session.commit()
            return true
}

func (c TendressesController) Notify(username_friend string) string {
	config.Global.Server.
}