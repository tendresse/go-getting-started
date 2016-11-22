package main

import (
    "encoding/json"
	"log"
	"net/http"
	"os"
	"time"
    "github.com/auth0/go-jwt-middleware"
    "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
    "github.com/googollee/go-socket.io"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}


func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	server, err := socketio.NewServer(nil)
    if err != nil {
        log.Fatal(err)
    }

	server.On("connection", func(so socketio.Socket) {
        log.Println("on connection")
        doEvery(20*time.Millisecond, so.Emit("time", time.Now().Format(time.RFC850)))
        so.On("disconnection", func() {
            log.Println("on disconnect")
        })
    })

    server.On("error", func(so socketio.Socket, err error) {
        log.Println("error:", err)
    })

    http.Handle("/socket.io/", server)
    http.Handle("/", http.FileServer(http.Dir("./asset")))
    log.Println("Serving at localhost:5000...")
    log.Fatal(http.ListenAndServe(":5000", nil))

    db, err := gorm.Open("postgres", "host=myhost user=gorm dbname=gorm sslmode=disable password=mypassword")
  
	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()
}
