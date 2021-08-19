package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
	"database/sql"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "80"
		log.Fatal("$PORT set to 80")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	Db = db

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

func contactMe(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
	
	script := `<html>
					<body>
						<script type="text/javascript">
							alert("Thank you for contacting me, your email has been received \nCheck your mail");
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

		myEmail := getEnv("EMAIL")
		// password := getEnv("PASSWORD")

		message := "Hello " + firstName + " " +  lastName + ". Thank you for reaching out to me."

		insertFormData(firstName, lastName, email, comment)

		sendMailSG(email, message)

		sendMail(email, message)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, script)

		fmt.Fprintf(w, "Hello, %q <br>first name:%s <br>last name: %s <br>email: %s <br>comment: %s <br>my email: %s", 
			html.EscapeString(r.URL.Path), firstName, lastName, email, comment, myEmail)
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

func sendMailSG(email string, body string) {
	myEmail := getEnv("EMAIL")
	from := mail.NewEmail("Kingsley Ugwudinso", myEmail)
	subject := "Hurray!!! I got it."
	to := mail.NewEmail("Kingsley Ugwudinso", email)
	plainTextContent := body
	htmlContent := "<strong>" + body + "</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(getEnv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
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
  
func insertFormData(firstName string, lastName string, email string, comment string) {
	sqlStatement := `INSERT INTO resume (first_name, last_name, email, comment)
	VALUES ($1, $2, $3, $4)`
	_, _ = Db.Exec(sqlStatement, firstName, lastName, email, comment)
}