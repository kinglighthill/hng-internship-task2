package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"database/sql"
	_ "github.com/lib/pq"
)

type emailSent func (gin.ResponseWriter, string) 

type formData struct {
	firstName string
	lastName string
	email string
	comment string
}

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

	r := c.Request
	w := c.Writer

	if r.Method == "POST" {
		firstName := c.PostForm("first-name")
		lastName := c.PostForm("last-name")
		email := c.PostForm("email")
		comment := c.PostForm("comment")

		data := formData {
			firstName: firstName,
			lastName: lastName,
			email: email,
			comment: comment,
		}
		insertFormData(data, func(w gin.ResponseWriter, script string) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			io.WriteString(w, script)
		}, w)
	}  
}

func showErrorAlert() string {
	script := `<html>
					<body>
						<script type="text/javascript">
							alert("Error submiting form");
						</script>
					</body>
			    </html>`
	return  script
}

func showSuccesAlert() string {
	script := `<html>
					<body>
						<script type="text/javascript">
							alert("Form submitted successfully. \nThank you for contacting me, your contact has been received");
						</script>
					</body>
			    </html>`
	return  script
}
  
func insertFormData(data formData, onEmailSent emailSent, w gin.ResponseWriter) {
	sqlStatement := `INSERT INTO resume (first_name, last_name, email, comment)
	VALUES ($1, $2, $3, $4)`
	_, err := Db.Exec(sqlStatement, data.firstName, data.lastName, data.email, data.comment)

	if err != nil {
		onEmailSent(w, showErrorAlert())
	}
	onEmailSent(w, showSuccesAlert())
}