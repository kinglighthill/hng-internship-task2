package main

/* import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
) */

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"net/smtp"

	// "crypto/tls"
	/* "github.com/gin-gonic/contrib/static" */
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/spf13/viper"
	// gomail "gopkg.in/mail.v2"
)

// const myEmail string = "testyholyhill@gmail.com"

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "80"
		log.Fatal("$PORT set to 80")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("index.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/", contactMe)

	router.Run(":" + port)
}

/* func formResponse(c *gin.Context) {
	w := c.Writer

	fmt.Fprintf(w, 
		`<html>
            <head>
            </head>
            <body>
            <h1>Go Timer (ticks every second!)</h1>
            <div id="output"></div>
            <script type="text/javascript">
            console.log("`+hello+`");
            </script>
            </body>
        </html>`)
} */

func contactMe(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
	
	script := `<html>
					<body>
						<script type="text/javascript">
							alert("Thank you for contacting me, your email has been received");
						</script>
					</body>
			    </html>`

	r := c.Request
	w := c.Writer

	if r.Method == "POST" {
		firstName := c.PostForm("first-name")
		lastName := c.PostForm("last-name")
		email := c.PostForm("email")
		comment := c.PostForm("comment")

		sendMail(email, "Hello " + firstName + " " +  lastName + ". Thank you for reaching out to me.")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, script)

		fmt.Fprintf(w, "Hello, %q <br>first name:%s <br>last name: %s <br>email: %s <br>comment: %s ", html.EscapeString(r.URL.Path), firstName, lastName, email, comment)
	}  
}

func sendMail(email string, body string) {
	myEmail := getEnv("EMAIL")
	password := getEnv("PASSWORD")

  	// Receiver email address.
	to := []string {email}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	
	// Authentication.
	auth := smtp.PlainAuth("", myEmail, password, smtpHost)
	
	// Message.
	message := []byte(body)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, myEmail, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func getEnv(key string) string {
	viper.SetConfigFile(".env")
  
	err := viper.ReadInConfig()
  
	if err != nil {
	  log.Fatalf("Error while reading config file %s", err)
	}
  
	value, ok := viper.Get(key).(string)
  
	if !ok {
	  log.Fatalf("Invalid type assertion")
	}
  
	return value
}
  

/* func main() {
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
} */