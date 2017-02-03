package app

import (
	"net/http"
	"os"

	"github.com/tendresse/go-getting-started/app/controllers"
	"github.com/tendresse/go-getting-started/app/config"

	"gopkg.in/pg.v5"
	log "github.com/Sirupsen/logrus"
	"github.com/rs/cors"
	"github.com/googollee/go-socket.io"
)

func init() {
	var err error

	initEnv()
	initDB()

	config.Global.Server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	config.Global.Server.On("connection", func(so socketio.Socket) {
		config.Global.Socket = &so
		so.Emit("chat", "hello message")
			log.Println("on connection 42")
		so.On("chat", func(msg string) {
			log.Println("recieved message", msg)
		})
		so.On("random", func(msg string) string {
			log.Println("commentaire random", msg)
			// gif := models.Gif{3,"oui.com","ouioui.com",3,"oui"}
			// b, _ := json.Marshal(gif)
			return `{"Id":3,"BlogID":2,"Url":"http://i.giphy.com/xTgeIYpiaiWvWHyWwU.gif","LameScore":4}`
		})
		// Socket.io acknowledgement example
		// The return type may vary depending on whether you will return
		// For this example it is "string" type
		so.On("chat message with ack", func(msg string) string {
			var gifs_controller controllers.GifsController
			return gifs_controller.RandomJSONGif()
		})
		so.On("disconnection", func() {
			log.Println("disconnected from chat")
		})
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})
	mux.Handle("/socket.io/", config.Global.Server)
	mux.Handle("/assets", http.FileServer(http.Dir("./assets")))

	/**
	* DEV MODE
	*  - enable CORS
	*/
	handler := cors.Default().Handler(mux)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	handler = c.Handler(handler)

	log.Println("Serving at localhost:5000...")
}

func initEnv() {
	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" {
		log.Fatal("env variable secret_key is not set")
	}
	config.Global.SecretKey = secret_key

	database_uri := os.Getenv("DATABASE_URI")
	if database_uri == "" {
		log.Fatal("env variable database_uri is not set")
	}
	config.Global.DatabaseURI = database_uri

	tumblr_api_key := os.Getenv("TUMBLR_API_KEY")
	if tumblr_api_key == "" {
		log.Fatal("env variable tumblr_api_key is not set")
	}
	config.Global.TumblrAPIKey = tumblr_api_key
}

func initDB() {
	var err error
	config.Global.DB = pg.Connect(&pg.Options{
		User: "postgres",
		Database: "postgres",
		Password: "postgres",
		Addr: "127.0.0.1:32768",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer config.Global.DB.Close()
}
