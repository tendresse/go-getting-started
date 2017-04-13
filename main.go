package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/tendresse/go-getting-started/app/config"
	"github.com/tendresse/go-getting-started/app/controllers"

	"github.com/Sirupsen/logrus"
	"gopkg.in/pg.v5"
	"golang.org/x/net/websocket"
	"github.com/graarh/golang-socketio"
	"github.com/garyburd/redigo/redis"
	"io"
	"time"
)

var (
	waitTimeout = time.Minute * 10
	log         = logrus.WithField("cmd", "go-websocket-chat-demo")
	rr          redisReceiver
	rw          redisWriter
)

func main() {
	initEnv()
	initDB()
	defer config.Global.DB.Close()

	var current_user_id int
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		current_user_id = 0
		c.Emit("ready","websocket is ready")
		log.Println("New client connected, client id is ", c.Id())
	})
	server.On(gosocketio.OnError, func(c *gosocketio.Channel) {
		log.Println("Error occurs")
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		current_user_id = 0
		log.Println("Disconnected")
	})
	server.On("signup", func(c *gosocketio.Channel, username string, password string) string {
		// return token
		return controllers.User{}.Signup(username, password, &current_user_id, c)
	})
	server.On("login", func(c *gosocketio.Channel, username string, password string) string {
		// return token
		return controllers.User{}.Login(username, password, &current_user_id, c)
	})
	server.On("login with token", func(c *gosocketio.Channel, token string) string {
		return controllers.User{}.LoginWithToken(token, &current_user_id, c)
	})
	server.On("ready", func(c *gosocketio.Channel, msg string) string {
		if current_user_id == 0 {
			return `{"success":false, "error":"you are not connected"}`
		}
		c.Emit("me", controllers.User{}.GetUser(current_user_id))
		c.Emit("tendresses", controllers.User{}.GetPendingTendresses(&current_user_id) )
		return `{"success":true}`
	})
	server.On("get profile", func(c *gosocketio.Channel, id string) string {
		if current_user_id == 0 {
			return `{"success":false, "error":"you are not connected"}`
		}
		user_id,err := strconv.Atoi(id)
		if err != nil {
			log.Error(err)
			return `{"success":false, "error":"user id cannot be converted"}`
		}
		log.Println("user_id ",user_id," profile was requested by user_id ",current_user_id)
		return controllers.User{}.GetProfile(user_id)
	})
	server.On("update device", func(c *gosocketio.Channel, device_token string) string {
		if current_user_id == 0 {
			return `{"success":false, "error":"you are not connected"}`
		}
		log.Println("user_id ",current_user_id," updated his device with device_token = ",device_token)
		return controllers.User{}.UpdateDevice(device_token, &current_user_id)
	})
	server.On("send tendresse", func(c *gosocketio.Channel, friend_id string) string {
		if current_user_id == 0 {
			return `{"success":false, "error":"you are not connected"}`
		}
		return controllers.Tendresse{}.SendTendresse(friend_id, &current_user_id, c)
	})
	server.On("tendresse seen", func(c *gosocketio.Channel, id_tendresse string) string {
		if current_user_id == 0 {
			return `{"success":false, "error":"you are not connected"}`
		}
		tendresse_id,err := strconv.Atoi(id_tendresse)
		if err != nil {
			log.Error(err)
			return `{"success":false, "error":"user id cannot be converted"}`
		}
		return controllers.Tendresse{}.SetTendresseAsSeen(tendresse_id, &current_user_id)
	})
	server.On("add friend", func(c *gosocketio.Channel, username string) string {
		if current_user_id == 0 {
			return `{"success":false, "error":"you are not connected"}`
		}
		return controllers.User{}.AddFriend(username, &current_user_id)
	})
	server.On("delete friend", func(c *gosocketio.Channel, id_friend string) string {
		if current_user_id == 0 {
			return `{"success":false, "error":"you are not connected"}`
		}
		friend_id,err := strconv.Atoi(id_friend)
		if err != nil {
			log.Error(err)
			return `{"success":false, "error":"user id cannot be converted"}`
		}
		return controllers.User{}.DeleteFriend(friend_id, &current_user_id)
	})
	server.On("random", func(c *gosocketio.Channel, ) string {
		return controllers.Gif{}.RandomGif()
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.WithField("PORT", port).Fatal("$PORT must be set")
	}

	redisURL := os.Getenv("REDIS_URL")
	redisPool, err := redis.NewRedisPoolFromURL(redisURL)
	if err != nil {
		log.WithField("url", redisURL).Fatal("Unable to create Redis pool")
	}

	rr = newRedisReceiver(redisPool)
	rw = newRedisWriter(redisPool)

	go func() {
		for {
			waited, err := redis.WaitForAvailability(redisURL, waitTimeout, rr.wait)
			if !waited || err != nil {
				log.WithFields(logrus.Fields{"waitTimeout": waitTimeout, "err": err}).Fatal("Redis not available by timeout!")
			}
			rr.broadcast(availableMessage)
			err = rr.run()
			if err == nil {
				break
			}
			log.Error(err)
		}
	}()

	go func() {
		for {
			waited, err := redis.WaitForAvailability(redisURL, waitTimeout, nil)
			if !waited || err != nil {
				log.WithFields(logrus.Fields{"waitTimeout": waitTimeout, "err": err}).Fatal("Redis not available by timeout!")
			}
			err = rw.run()
			if err == nil {
				break
			}
			log.Error(err)
		}
	}()


	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir("./public")))
	serveMux.HandleFunc("/ws", handleWebsocket)
	static := http.StripPrefix("/static", http.FileServer(http.Dir("static/")))
	serveMux.Handle("/static/", static)
	serveMux.Handle("/socket", websocket.Handler(UserSocket))
	serveMux.Handle("/admin", websocket.Handler(AdminSocket))
	log.Println("Starting server on port 3000 ...")
	err := http.ListenAndServe(":"+port, serveMux)
	if err != nil {
		log.Panic("ListenAndServe: " + err.Error())
	}
}

func initEnv() {
	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" {
		log.Error("env variable secret_key is not set")
	}
	config.Global.SecretKey = secret_key

	database_uri := os.Getenv("DATABASE_URI")
	if database_uri == "" {
		log.Error("env variable database_uri is not set")
		// testing
		database_uri = "postgres://postgres:postgres@localhost:32768/postgres?sslmode=disable"
	}
	config.Global.DatabaseURI = database_uri

	tumblr_api_key := os.Getenv("TUMBLR_API_KEY")
	if tumblr_api_key == "" {
		log.Error("env variable tumblr_api_key is not set")
	}
	config.Global.TumblrAPIKey = tumblr_api_key
}

func initDB() {
	var err error
	db_url,err := pg.ParseURL(config.Global.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	config.Global.DB = pg.Connect(db_url)
	if err != nil {
		log.Fatal(err)
	}
	db := controllers.DB{}
	db.LoadDB()
}

func UserSocket(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func AdminSocket(ws *websocket.Conn) {
	io.Copy(ws, ws)
}
