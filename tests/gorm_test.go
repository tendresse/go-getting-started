package tests

import (
	// "encoding/json"
	"gopkg.in/pg.v5"
	"log"
	"fmt"
	"github.com/tendresse/go-getting-started/app/models"
	"github.com/tendresse/go-getting-started/app/dao"
)

func executeQueries(queries []string) error {
	for _, q := range queries {
		_, err := Global.DB.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

type GlobalStruct struct{
	DB	*pg.DB
}

var Global GlobalStruct
var achievements_dao dao.Achievement

func init() {
	initDB()
	loadDB()
	achievement := models.Achievement{}
	if err := achievements_dao.GetAchievementByTitle("mommy", &achievement); err != nil {
		fmt.Println(err)
	}else{
		fmt.Println(achievement)
	}
	//
}

func initDB() {
	Global.DB = pg.Connect(&pg.Options{
		User: "postgres",
		Database: "postgres",
		Password: "postgres",
		Addr: "127.0.0.1:32768",
	})
	if err := executeQueries(delete_schema); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("deleted schema")

	if err := executeQueries(create_schema); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("created schema")

}

func loadDB(){
	tag_milf := models.Tag{
		Name: "milf",
		Banned: false,
	}
	if err := Global.DB.Insert(&tag_milf); err != nil {
		fmt.Println(err)
	}
	tag_boobs := models.Tag{
		Name: "boobs",
		Banned: false,
	}
	if err := Global.DB.Insert(&tag_boobs); err != nil {
		fmt.Println(err)
	}
	tag_gif := models.Tag{
		Name: "gif",
		Banned: true,
	}
	if err := Global.DB.Insert(&tag_gif); err != nil {
		fmt.Println(err)
	}
	tag_cowboy := models.Tag{
		Name: "cowboy",
		Banned: false,
	}
	if err := Global.DB.Insert(&tag_cowboy); err != nil {
		fmt.Println(err)
	}
	achievement_mommy := models.Achievement{
		Condition: 10,
		Icon: "ach1.png",
		TagID: 1,
		Title: "mommy",
		TypeOf: "send",
		XP: 100,
	}
	if err := Global.DB.Insert(&achievement_mommy); err != nil {
		fmt.Println(err)
	}
	achievement_rider := models.Achievement{
		Condition: 10,
		Icon: "ach2.jpg",
		TagID: 4,
		Title: "rider cowboy",
		TypeOf: "receive",
		XP: 5,
	}
	if err := Global.DB.Insert(&achievement_rider); err != nil {
		fmt.Println(err)
	}
	blog_milf := models.Blog{
		Name: "milf_tumblr",
		Url: "https://milf.tumblr.com",
	}
	if err := Global.DB.Insert(&blog_milf); err != nil {
		fmt.Println(err)
	}
	blog_cowboy:= models.Blog{
		Name: "super_cowboy",
		Url: "https://super_cowboy.tumblr.com",
	}
	if err := Global.DB.Insert(&blog_cowboy); err != nil {
		fmt.Println(err)
	}
	gif_milf_pov := models.Gif{
		BlogID:1,
		Url:"https://sc1.tumblr.com/gif1",
		LameScore:0,
	}
	if err := Global.DB.Insert(&gif_milf_pov); err != nil {
		fmt.Println(err)
	}
	gif_gays := models.Gif{
		BlogID:2,
		Url:"https://sc1.tumblr.com/gif2",
		LameScore:0,
	}
	if err := Global.DB.Insert(&gif_gays); err != nil {
		fmt.Println(err)
	}
	role_admin := models.Role{
		Name:"admin",
	}
	if err := Global.DB.Insert(&role_admin); err != nil {
		fmt.Println(err)
	}
	role_modo := models.Role{
		Name:"modo",
	}
	if err := Global.DB.Insert(&role_modo); err != nil {
		fmt.Println(err)
	}
	user_alice := models.User{
		Device: "token_android_123",
		NSFW: true,
		Username: "alice",
		Passhash: "md5_32_123",
		Premium: false,
	}
	if err := Global.DB.Insert(&user_alice); err != nil {
		fmt.Println(err)
	}
	user_bob := models.User{
		Device: "token_ios_123",
		NSFW: true,
		Username: "bob",
		Passhash: "md5_32_123",
		Premium: false,
	}
	if err := Global.DB.Insert(&user_bob); err != nil {
		fmt.Println(err)
	}
	alice_bob_1 := models.Tendresse{
		SenderID: user_alice.ID,
		ReceiverID: user_bob.ID,
		GifID: gif_gays.ID,
		Viewed: false,
	}
	if err := Global.DB.Insert(&alice_bob_1); err != nil {
		fmt.Println(err)
	}
	alice_bob_2 := models.Tendresse{
		SenderID: user_alice.ID,
		ReceiverID: user_bob.ID,
		GifID: gif_milf_pov.ID,
		Viewed: true,
	}
	if err := Global.DB.Insert(&alice_bob_2); err != nil {
		fmt.Println(err)
	}
	bob_alice_1 := models.Tendresse{
		SenderID: user_bob.ID,
		ReceiverID: user_alice.ID,
		GifID: gif_gays.ID,
		Viewed: false,
	}
	if err := Global.DB.Insert(&bob_alice_1); err != nil {
		fmt.Println(err)
	}
	relation_gif_tags_1 := models.GifsTags{
			GifID: gif_milf_pov.ID,
			TagID: tag_milf.ID,
		}
	if err := Global.DB.Insert(&relation_gif_tags_1); err != nil {
		fmt.Println(err)
	}
	relation_gif_tags_2 := models.GifsTags{
		GifID: gif_milf_pov.ID,
		TagID: tag_boobs.ID,
	}
	if err := Global.DB.Insert(&relation_gif_tags_2); err != nil {
		fmt.Println(err)
	}
	relation_gif_tags_3 := models.GifsTags{
		GifID: gif_gays.ID,
		TagID: tag_cowboy.ID,
	}
	if err := Global.DB.Insert(&relation_gif_tags_3); err != nil {
		fmt.Println(err)
	}
	relation_users_achievements_1 := models.UsersAchievements{
		UserID:user_alice.ID,
		AchievementID:achievement_mommy.ID,
	}
	if err := Global.DB.Insert(&relation_users_achievements_1); err != nil {
		fmt.Println(err)
	}
	relation_users_achievements_2 := models.UsersAchievements{
		UserID:user_alice.ID,
		AchievementID:achievement_rider.ID,
	}
	if err := Global.DB.Insert(&relation_users_achievements_2); err != nil {
		fmt.Println(err)
	}
	relation_users_friends_1 := models.UsersFriends{
		UserID:user_alice.ID,
		FriendID:user_bob.ID,
	}
	if err := Global.DB.Insert(&relation_users_friends_1); err != nil {
		fmt.Println(err)
	}
	relation_users_friends_2 := models.UsersFriends{
		UserID:user_bob.ID,
		FriendID:user_alice.ID,
	}
	if err := Global.DB.Insert(&relation_users_friends_2); err != nil {
		fmt.Println(err)
	}
	relation_users_roles_1 := models.UsersRoles{
		UserID:user_alice.ID,
		RoleID:role_admin.ID,
	}
	if err := Global.DB.Insert(&relation_users_roles_1); err != nil {
		fmt.Println(err)
	}
	relation_users_roles_2 := models.UsersRoles{
		UserID:user_bob.ID,
		RoleID:role_modo.ID,
	}
	if err := Global.DB.Insert(&relation_users_roles_2); err != nil {
		fmt.Println(err)
	}
}