package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
    DB_USER     = "postgres"
    DB_PASSWORD = "TiCuAhPl08"
    TABLE_NAME     = "users"
)

type User struct {
	ID int `json:"id"`
	UserName string `json:"name"`
	Email string `json:"email"`
}

func getAllUsers(c *gin.Context, db *sql.DB) {
	var res User;
	var users []User;
	rows, err := db.Query("SELECT * FROM users")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		if err := rows.Scan(&res.ID, &res.UserName, &res.Email); err != nil {
			log.Fatalln(err)
		}
		users = append(users, res)
	}
	c.IndentedJSON(http.StatusOK, users)
}

func getUserWithID(c *gin.Context, db *sql.DB) {
	var res User;
	rows, err := db.Query("SELECT * FROM users WHERE id = $1", c.Param("id"))
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
	rows.Next()
	if err := rows.Scan(&res.ID, &res.UserName, &res.Email); err != nil {
		log.Fatalln(err)
	}
	c.IndentedJSON(http.StatusOK, []User{res})
}

func createUser(c *gin.Context, db *sql.DB) {
	var newUser User
	if err := c.Bind(&newUser); err != nil {
		return
	}
	if _, err := db.Exec("INSERT INTO users VALUES ($1, $2, $3)", newUser.ID, newUser.UserName, newUser.Email); err != nil {
		log.Fatalln(err)
	}
	getAllUsers(c, db)
}

func deleteAllUsers(c *gin.Context, db *sql.DB) {
	if _, err := db.Exec("DELETE FROM users"); err != nil {
		log.Fatalln(err)
	}

	c.IndentedJSON(http.StatusOK, []User{})
}

func deleteUserWithID(c *gin.Context, db *sql.DB) {
	if _, err := db.Exec("DELETE FROM users WHERE id = $1", c.Param("id")); err != nil {
		log.Fatalln(err)
	}

	getAllUsers(c, db)
}

func main() {
	// connect to database
	connStr := fmt.Sprintf("postgresql://%v:%v@localhost:5432/%v?sslmode=disable", DB_USER, DB_PASSWORD, TABLE_NAME)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/users", func (c *gin.Context) {
		getAllUsers(c, db)
	})
	router.GET("/users/:id", func (c *gin.Context) {
		getUserWithID(c, db)
	})
	router.POST("/users", func (c *gin.Context) {
		createUser(c, db)
	})
	router.DELETE("/users", func (c *gin.Context) {
		deleteAllUsers(c, db)
	})
	router.DELETE("/users/:id", func (c *gin.Context) {
		deleteUserWithID(c, db)
	})

	router.Run("localhost:8000")
}