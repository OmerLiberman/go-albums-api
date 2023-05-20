package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type album struct {
	Title  string  `bson:"title,omitempty"`
	Artist string  `bson:"artist,omitempty"`
	Price  float64 `bson:"price,omitempty"`
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	albums := []interface{}{
		album{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		album{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		album{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	albums := []interface{}{
		album{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		album{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		album{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// func getAlbumByID(c *gin.Context) {
// 	id := c.Param("id")
// 	albums := []interface{}{
// 		album{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 		album{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 		album{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// 	}

// 	for _, a := range albums {
// 		if a.ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
// }

func connectToDB() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	return client
}

func main() {
	var client = connectToDB()
	albumsCollection := client.Database("AlbumsApp").Collection("albums")

	albums := []interface{}{
		album{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		album{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		album{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}

	_, err := albumsCollection.InsertMany(context.TODO(), albums)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	// router.GET("/albums/:id", getAlbumByID)

	router.Run("localhost:3001")
}
