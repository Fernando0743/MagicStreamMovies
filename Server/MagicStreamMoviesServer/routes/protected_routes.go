package routes

import (
	controller "github.com/Fernando0743/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/Fernando0743/MagicStreamMovies/Server/MagicStreamMoviesServer/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Setup protected routes
func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	//Protect relevant routes (Auth Middleware is a  Gin handler function used to validate incoming access tokens
	// and grant/prohibt access to protected endpoints)
	router.Use(middleware.AuthMiddleware())

	//PROTECTED ROUTES

	//Route that returns a single movie from DB given IMDB id
	router.GET("/movie/:imdb_id", controller.GetMovie(client))

	//Route that creates and insert one movie to movies collection in DB
	router.POST("/addmovie", controller.AddMovie(client))

	//Route that updates movie review
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate(client))

	//Route that fecthes recommended movies for user
	router.GET("/recommendedmovies", controller.GetRecommendedMovies(client))
}
