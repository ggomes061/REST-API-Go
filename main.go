package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:*****@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", createAlbum)

	router.Run("localhost:8080")
}

func createAlbum(c *gin.Context) {
	var awesomeAlbum album
	if err := c.BindJSON(&awesomeAlbum); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	stmt, err := db.Prepare("INSERT INTO albums (id, title, artist, price) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	if _, err := stmt.Exec(awesomeAlbum.ID, awesomeAlbum.Title, awesomeAlbum.Artist, awesomeAlbum.Price); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusCreated, awesomeAlbum)
}

func getAlbums(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, title, artist, price FROM albums")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var albums []album
	for rows.Next() {
		var a album
		err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
		if err != nil {
			log.Fatal(err)
		}

		albums = append(albums, a)

	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.IndentedJSON(http.StatusOK, albums)

}
