package main

/* import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
) */

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

/* func main() {
    r := gin.Default()

    r.GET("/hello", func(c *gin.Context) {
        c.String(200, "Hello, World!")
    })

	r.Use(static.Serve("/", static.LocalFile("./", true)))
    r.Run()
} */

func main() {
	port := GetPort()
	http.HandleFunc("/", hello)
	log.Print("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func GetPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	return ":" + port
}