package tests

import (
	// "encoding/json"
	"github.com/go-pg/pg"
	log "github.com/Sirupsen/logrus"
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

func init() {
	//initDB()
	//loadDB()
	//testTag()
	//testAchievement()
	//testBlog()
	//testGif()
	//testRole()
	//testUser()
	//testTendresse()
}

func initDB() {
	Global.DB = pg.Connect(&pg.Options{
		User: "postgres",
		Database: "postgres",
		Password: "postgres",
		Addr: "127.0.0.1:32768",
	})
	if err := executeQueries(delete_schema); err != nil {
		fmt.Println("bug ?")
		log.Fatalln(err)
	}
	fmt.Println("deleted schema")

	if err := executeQueries(create_schema); err != nil {
		fmt.Println("bug ?")
		log.Fatalln(err)
	}
	fmt.Println("created schema")
}

func loadDB(){
	// TAGS
	tags_dao := dao.Tag{DB:Global.DB}

	tag_milf := models.Tag{
		Title: "milf",
	}
	tag_boobs := models.Tag{
		Title: "boobs",
	}
	tag_gif := models.Tag{
		Title: "gif",
		Banned: true,
	}
	tag_cowboy := models.Tag{
		Title: "cowboy",
	}
	if err := tags_dao.CreateTags([]*models.Tag{&tag_milf,&tag_cowboy,&tag_boobs,&tag_gif}); err != nil {
		fmt.Println("CreateTags")
		fmt.Println(err)
	}

	// ACHIEVEMENTS
	achievements_dao := dao.Achievement{DB:Global.DB}

	achievement_mommy := models.Achievement{
		Condition: 10,
		Icon:      "ach1.png",
		TagId:     1,
		Title:     "mommy",
		TypeOf:    "send",
		XP:        100,
	}
	achievement_rider := models.Achievement{
		Condition: 10,
		Icon:      "ach2.jpg",
		TagId:     4,
		Title:     "rider cowboy",
		TypeOf:    "receive",
		XP:        5,
	}
	if err := achievements_dao.CreateAchievements([]*models.Achievement{&achievement_mommy, &achievement_rider}); err != nil {
		fmt.Println("CreateAchievements")
		fmt.Println(err)
	}

	// BLOGS
	blog_dao := dao.Blog{DB : Global.DB}

	blog_milf := models.Blog{
		Title: "milf_tumblr",
		Url: "https://milf.tumblr.com",
	}
	blog_cowboy:= models.Blog{
		Title: "super_cowboy",
		Url: "https://super_cowboy.tumblr.com",
	}
	if err := blog_dao.CreateBlogs([]*models.Blog{&blog_milf,&blog_cowboy}); err != nil {
		fmt.Println("CreateBlog")
		fmt.Println(err)
	}

	// GIFS
	gifs_dao := dao.Gif{DB:Global.DB}

	gif_milf_pov := models.Gif{
		BlogId:    1,
		Url:       "https://sc1.tumblr.com/gif1",
		LameScore: 0,
	}
	gif_gays := models.Gif{
		BlogId:    2,
		Url:       "https://sc1.tumblr.com/gif2",
		LameScore: 0,
	}
	if err := gifs_dao.CreateGifs([]*models.Gif{&gif_milf_pov, &gif_gays}); err != nil {
		fmt.Println("CreateGifs")
		fmt.Println(err)
	}

	// RELATION GIF TAGS
	if err := gifs_dao.AddTagToGif(&tag_milf, &gif_milf_pov); err != nil {
		fmt.Println("AddTagsToGif")
		fmt.Println(err)
	}
	if err := gifs_dao.AddTagToGif(&tag_boobs, &gif_milf_pov); err != nil {
		fmt.Println("AddTagsToGif")
		fmt.Println(err)
	}
	if err := gifs_dao.AddTagToGif(&tag_cowboy, &gif_gays); err != nil {
		fmt.Println("AddTagToGif")
		fmt.Println(err)
	}

	// ROLES
	role_dao := dao.Role{DB:Global.DB}
	role_admin := models.Role{
		Title:"admin",
	}
	if err := role_dao.CreateRole(&role_admin); err != nil {
		fmt.Println("CreateRole")
		fmt.Println(err)
	}
	role_modo := models.Role{
		Title:"modo",
	}
	if err := role_dao.CreateRole(&role_modo); err != nil {
		fmt.Println("CreateRole")
		fmt.Println(err)
	}


	// USERS
	users_dao := dao.User{DB:Global.DB}
	user_alice := models.User{
		Device: "token_android_123",
		NSFW: true,
		Username: "alice",
		Passhash: "md5_32_123",
		Premium: false,
	}
	if err := users_dao.CreateUser(&user_alice); err != nil {
		fmt.Println("CreateUser")
		fmt.Println(err)
	}
	user_bob := models.User{
		Device: "token_ios_123",
		NSFW: true,
		Username: "bob",
		Passhash: "md5_32_123",
		Premium: false,
	}
	if err := users_dao.CreateUser(&user_bob); err != nil {
		fmt.Println("CreateUser")
		fmt.Println(err)
	}

	// FRIENDS
	if err := users_dao.AddFriendToUser(&user_bob,&user_alice); err != nil {
		fmt.Println("AddFriendToUser")
		fmt.Println(err)
	}
	if err := users_dao.AddFriendToUser(&user_alice,&user_bob); err != nil {
		fmt.Println("AddFriendToUser")
		fmt.Println(err)
	}

	// RELATIONS USERS ACHIEVEMENTS
	if err := users_dao.AddAchievementToUser(&achievement_mommy,&user_alice); err != nil {
		fmt.Println("AddAchievementsToUser")
		fmt.Println(err)
	}
	if err := users_dao.AddAchievementToUser(&achievement_rider,&user_alice); err != nil {
		fmt.Println("AddAchievementsToUser")
		fmt.Println(err)
	}


	// RELATION USERS ROLES
	if err := users_dao.AddRoleToUser(&role_admin,&user_alice); err != nil {
		fmt.Println("AddRoleToUser")
		fmt.Println(err)
	}
	if err := users_dao.AddRoleToUser(&role_modo,&user_bob); err != nil {
		fmt.Println("AddRoleToUser")
		fmt.Println(err)
	}

	// TENDRESSES
	tendresses_dao := dao.Tendresse{DB:Global.DB}
	alice_bob_1 := models.Tendresse{
		SenderId:   user_alice.Id,
		ReceiverId: user_bob.Id,
		GifId:      gif_gays.Id,
		Viewed:     false,
	}
	alice_bob_2 := models.Tendresse{
		SenderId:   user_alice.Id,
		ReceiverId: user_bob.Id,
		GifId:      gif_milf_pov.Id,
		Viewed:     true,
	}
	bob_alice_1 := models.Tendresse{
		SenderId:   user_bob.Id,
		ReceiverId: user_alice.Id,
		GifId:      gif_gays.Id,
		Viewed:     false,
	}
	if err := tendresses_dao.CreateTendresse(&alice_bob_1); err != nil {
		fmt.Println("CreateTendresse")
		fmt.Println(err)
	}
	if err := tendresses_dao.CreateTendresse(&alice_bob_2); err != nil {
		fmt.Println("CreateTendresse")
		fmt.Println(err)
	}
	if err := tendresses_dao.CreateTendresse(&bob_alice_1); err != nil {
		fmt.Println("CreateTendresse")
		fmt.Println(err)
	}
}

func testTag(){
	tags_dao := dao.Tag{DB:Global.DB}
	gifs_dao := dao.Gif{DB:Global.DB}

	tag_create := models.Tag{
		Title: "bbc",
		Banned: false,
	}
	if err := tags_dao.CreateTag(&tag_create); err != nil {
		fmt.Println("CreateTag")
		fmt.Println(err)
	}

	tag_create.Title = "i like trains"
	gif := models.Gif{Id: 1}
	gifs_dao.GetGif(&gif)
	tag_create.Gifs = []models.Gif{gif}
	if err := tags_dao.UpdateTag(&tag_create); err != nil {
		fmt.Println("UpdateTag")
		fmt.Println(err)
	}

	tag := models.Tag{Id: 1}
	if err := tags_dao.GetTag(&tag); err != nil {
		fmt.Println("GetTag")
		fmt.Println(err)
	}

	tag_full := models.Tag{Id: 1}
	if err := tags_dao.GetFullTag(&tag_full); err != nil {
		fmt.Println("GetFullTag")
		fmt.Println(err)
	}

	tag_by_title := models.Tag{}
	if err := tags_dao.GetTagByTitle("milf", &tag_by_title); err != nil {
		fmt.Println("GetTagByTitle")
		fmt.Println(err)
	}

	tags := []*models.Tag{{Id: 1}, {Id: 2}}
	if err := tags_dao.GetTags(tags); err != nil {
		fmt.Println("GetTags")
		fmt.Println(err)
	}

	tags_full := []*models.Tag{{Id: 1}, {Id: 2}}
	if err := tags_dao.GetFullTags(tags_full); err != nil {
		fmt.Println("GetFullTags")
		fmt.Println(err)
	}
}

func testAchievement() {
	achievements_dao := dao.Achievement{DB:Global.DB}

	achievement_create := models.Achievement{
		Condition: 10,
		Icon:      "ach5.jpg",
		TagId:     4,
		Title:     "i like dolphin",
		TypeOf:    "receive",
		XP:        10,
	}
	if err := achievements_dao.CreateAchievement(&achievement_create); err != nil {
		fmt.Println("CreateAchievement")
		fmt.Println(err)
	}

	achievement_create.Title = "i like trains"
	if err := achievements_dao.UpdateAchievement(&achievement_create); err != nil {
		fmt.Println("UpdateAchievement")
		fmt.Println(err)
	}

	achievement := models.Achievement{Id: 1}
	if err := achievements_dao.GetAchievement(&achievement); err != nil {
		fmt.Println("GetAchievement")
		fmt.Println(err)
	}

	achievement_full := models.Achievement{Id: 1}
	if err := achievements_dao.GetFullAchievement(&achievement_full); err != nil {
		fmt.Println("GetFullAchievement")
		fmt.Println(err)
	}

	achievement_by_title := models.Achievement{}
	if err := achievements_dao.GetAchievementByTitle("mommy", &achievement_by_title); err != nil {
		fmt.Println("GetAchievementByTitle")
		fmt.Println(err)
	}

	achievements := []*models.Achievement{{Id: 1}, {Id: 2}}
	if err := achievements_dao.GetAchievements(achievements); err != nil {
		fmt.Println("GetAchievements")
		fmt.Println(err)
	}

	achievements_full := []*models.Achievement{{Id: 1}, {Id: 2}}
	if err := achievements_dao.GetFullAchievements(achievements_full); err != nil {
		fmt.Println("GetFullAchievements")
		fmt.Println(err)
	}
}

func testBlog(){
	blog_dao := dao.Blog{DB : Global.DB}

	blog_create := models.Blog{Title:"test_blog", Url:"http://test.tumblr.com"}
	if err := blog_dao.CreateBlog(&blog_create); err != nil {
		fmt.Println("CreateBlog")
		fmt.Println(err)
	}

	blog_create.Title = "updated_name"
	if err := blog_dao.UpdateBlog(&blog_create); err != nil {
		fmt.Println("UpdateBlog")
		fmt.Println(err)
	}

	blog_get := models.Blog{Id: 1}
	if err := blog_dao.GetBlog(&blog_get); err != nil {
		fmt.Println("GetBlog")
		fmt.Println(err)
	}


	blog_full := models.Blog{Id: 2}
	if err := blog_dao.GetFullBlog(&blog_full); err != nil {
		fmt.Println("GetFullBlog")
		fmt.Println(err)
	}


	blog_title := models.Blog{}
	if err := blog_dao.GetBlogByTitle("updated_name", &blog_title); err != nil {
		fmt.Println("GetBlogByTitle")
		fmt.Println(err)
	}


	blogs := []*models.Blog{{Id: 1}, {Id: 2}}
	if err := blog_dao.GetBlogs(blogs); err != nil {
		fmt.Println("GetBlogs")
		fmt.Println(err)
	}

	blogs_full := []*models.Blog{{Id: 1}, {Id: 2}}
	if err := blog_dao.GetFullBlogs(blogs_full); err != nil {
		fmt.Println("GetFullBlogs")
		fmt.Println(err)
	}
}

func testGif(){
	gifs_dao := dao.Gif{DB:Global.DB}

	gif_create := models.Gif{
		BlogId:    2,
		Url:       "https://sc1.tumblr.com/gif14",
		LameScore: 0,
	}
	if err := gifs_dao.CreateGif(&gif_create); err != nil {
		fmt.Println("CreateGif")
		fmt.Println(err)
	}

	gif_create.LameScore = 10
	if err := gifs_dao.UpdateGif(&gif_create); err != nil {
		fmt.Println("UpdateGif")
		fmt.Println(err)
	}

	gif := models.Gif{Id: 1}
	if err := gifs_dao.GetGif(&gif); err != nil {
		fmt.Println("GetGif")
		fmt.Println(err)
	}

	gif_full := models.Gif{Id: 1}
	if err := gifs_dao.GetFullGif(&gif_full); err != nil {
		fmt.Println("GetFullGif")
		fmt.Println(err)
	}

	gifs := []*models.Gif{{Id: 1}, {Id: 2}}
	if err := gifs_dao.GetGifs(gifs); err != nil {
		fmt.Println("GetGifs")
		fmt.Println(err)
	}
	all_gifs := []models.Gif{}
	if err := gifs_dao.GetAllGifs(&all_gifs); err != nil {
		fmt.Println("GetGifs")
		fmt.Println(err)
	}
	fmt.Println(all_gifs)

	gifs_full := []*models.Gif{{Id: 1}, {Id: 2}}
	if err := gifs_dao.GetFullGifs(gifs_full); err != nil {
		fmt.Println("GetFullGifs")
		fmt.Println(err)
	}

	random := models.Gif{}
	if err := gifs_dao.GetRandomGif(&random); err != nil {
		fmt.Println("GetRandomGif")
		fmt.Println(err)
	}

	gif_url := models.Gif{Url:"https://sc1.tumblr.com/gif14"}
	if err := gifs_dao.GetGifByUrl(&gif_url); err != nil {
		fmt.Println("GetGifByUrl")
		fmt.Println(err)
	}
}

func testRole(){
	roles_dao := dao.Role{DB : Global.DB}

	role_create := models.Role{Title:"admin"}
	if err := roles_dao.CreateRole(&role_create); err != nil {
		fmt.Println("CreateRole")
		fmt.Println(err)
	}

	role_create.Title = "admin_updated"
	if err := roles_dao.UpdateRole(&role_create); err != nil {
		fmt.Println("UpdateRole")
		fmt.Println(err)
	}

	role := models.Role{Id: 1}
	if err := roles_dao.GetRole(&role); err != nil {
		fmt.Println("GetRole")
		fmt.Println(err)
	}

	role_full := models.Role{Id: 1}
	if err := roles_dao.GetFullRole(&role_full); err != nil {
		fmt.Println("GetFullRole")
		fmt.Println(err)
	}


	roles := []*models.Role{{Id: 1}, {Id: 2}}
	if err := roles_dao.GetRoles(roles); err != nil {
		fmt.Println("GetRoles")
		fmt.Println(err)
	}

	roles_full := []*models.Role{{Id: 1}, {Id: 2}}
	if err := roles_dao.GetFullRoles(roles_full); err != nil {
		fmt.Println("GetFullRoles")
		fmt.Println(err)
	}
}

func testUser(){
	users_dao := dao.User{DB:Global.DB}

	user_create := models.User{
		Device: "token_android_123",
		NSFW: true,
		Username: "create",
		Passhash: "md5_32_123",
		Premium: false,
	}
	if err := users_dao.CreateUser(&user_create); err != nil {
		fmt.Println("CreateUser")
		fmt.Println(err)
	}

	user_create.Username = "updated"
	if err := users_dao.UpdateUser(&user_create); err != nil {
		fmt.Println("UpdateUser")
		fmt.Println(err)
	}

	user := models.User{Id: 1}
	if err := users_dao.GetUser(&user); err != nil {
		fmt.Println("GetUser")
		fmt.Println(err)
	}

	user_full := models.User{Id: 1}
	if err := users_dao.GetFullUser(&user_full); err != nil {
		fmt.Println("GetFullUser")
		fmt.Println(err)
	}


	username := "updated"
	user_alice := models.User{}
	if err := users_dao.GetUserByUsername(&user_alice, username); err != nil {
		fmt.Println("GetUserByUsername")
		fmt.Println(err)
	}


	role := models.Role{Id: 3}
	if err := users_dao.AddRoleToUser(&role, &user_alice); err != nil {
		fmt.Println("AddRoleToUser")
		fmt.Println(err)
	}

	roles := []*models.Role{{Id: 1}, {Id: 2}}
	if err := users_dao.AddRolesToUser(roles, &user_alice); err != nil {
		fmt.Println("AddRolesToUser")
		fmt.Println(err)
	}

	role = models.Role{Id: 3}
	if err := users_dao.DeleteRoleFromUser(&role, &user_alice); err != nil {
		fmt.Println("DeleteRoleFromUser")
		fmt.Println(err)
	}

	if err := users_dao.DeleteRolesFromUser(roles, &user_alice); err != nil {
		fmt.Println("DeleteRolesFromUser")
		fmt.Println(err)
	}

	roles = []*models.Role{{Id: 1}, {Id: 2}}
	if err := users_dao.AddRolesToUser(roles, &user_alice); err != nil {
		fmt.Println("AddRolesToUser")
		fmt.Println(err)
	}

	friend := models.User{Id: 2}
	if err := users_dao.AddFriendToUser(&friend, &user_create); err != nil {
		fmt.Println("AddFriendToUser")
		fmt.Println(err)
	}
	if err := users_dao.DeleteFriendFromUser(&friend, &user_create); err != nil {
		fmt.Println("DeleteFriendFromUser")
		fmt.Println(err)
	}

	friends := []*models.User{{Id: 1}, {Id: 2}}
	if err := users_dao.AddFriendsToUser(friends, &user_create); err != nil {
		fmt.Println("AddFriendsToUser")
		fmt.Println(err)
	}
	if err := users_dao.DeleteFriendsFromUser(friends, &user_create); err != nil {
		fmt.Println("DeleteFriendsFromUser")
		fmt.Println(err)
	}

	ach := models.Achievement{Id: 1}
	if err := users_dao.AddAchievementToUser(&ach, &user_create); err != nil {
		fmt.Println("AddAchievementToUser")
		fmt.Println(err)
	}
	if err := users_dao.DeleteAchievementFromUser(&ach, &user_create); err != nil {
		fmt.Println("DeleteAchievementFromUser")
		fmt.Println(err)
	}

	if err := users_dao.DeleteUser(&user_full); err != nil {
		fmt.Println("DeleteUser")
		fmt.Println(err)
	}
}

func testTendresse(){
	tendresses_dao := dao.Tendresse{DB : Global.DB}

	tendresse := models.Tendresse{Id: 1}
	if err := tendresses_dao.GetTendresse(&tendresse); err != nil {
		fmt.Println("GetTendresse")
		fmt.Println(err)
	}
	tendresses_full := []*models.Tendresse{{Id: 1}, {Id: 2}}
	if err := tendresses_dao.GetFullTendresses(tendresses_full); err != nil {
		fmt.Println("GetFullTendresses")
		fmt.Println(err)
	}

	// users_dao := dao.User{DB:Global.DB}
	alice := models.User{Id: 1}
	tendresses_alice, err := tendresses_dao.GetPendingTendresses(&alice)
	if err != nil {
		fmt.Println("GetPendingTendresses")
		fmt.Println(err)
	}
	fmt.Println(tendresses_alice)

	if err := tendresses_dao.DeleteTendresse(&tendresse); err != nil {
		fmt.Println("DeleteTendresse")
		fmt.Println(err)
	}
	//var count int
	//if err := tendresses_dao.CountSenderTendresses(&count, 1); err != nil {
	//	fmt.Println("DeleteTendresse")
	//	fmt.Println(err)
	//}
	//fmt.Println(count)
}