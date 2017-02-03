package main

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
	"github.com/tendresse/go-getting-started/controllers"
	"github.com/tendresse/go-getting-started/models"
	"log"
	"net/http"
	"os"
	// "github.com/auth0/go-jwt-middleware"
	// "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var DB gorm.DB

func main() {
	var err error

	/**
	 * CONFIGURATION + ENV
	 */
	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" {
		log.Fatal("env variable secret_key is not set")
	}
	database_uri := os.Getenv("DATABASE_URI")
	if database_uri == "" {
		log.Fatal("env variable database_uri is not set")
	}

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {

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
			return gifs_controller.Random()
		})

		so.On("disconnection", func() {
			log.Println("disconnected from chat")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	/**
	 * SERVER MUX
	 * TODO :
	 *  -
	 */

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})
	mux.Handle("/socket.io/", server)
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

	/**
	 * DATABASE
	 * TODO :
	 *  - connection POSTGRE
	 */
	// Gorm, err := gorm.Open("postgres", "host=localhost:32768 user=postgres dbname=postgres sslmode=disable password=postgres"
	DB, err = gorm.Open("sqlite3", ".gorm.db")
	DB.LogMode(true)
	if err != nil {
		log.Println("Failed to connect to DB")
	} else {
		log.Println("DB Connected")
	}
	defer DB.Close()

	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", handler))

}
