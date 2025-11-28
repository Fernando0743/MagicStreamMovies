package routes

import (
	controller "github.com/Fernando0743/MagicStreamMovies/Server/MagicStreamMoviesServer/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Setup protected routes
func SetupUnProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	//NO MIDDLEWARE BECAUSE UNPROTECTED ROUTES
	//UNPROTECTED ROUTES

	//Route that returns all movies from DB
	router.GET("/movies", controller.GetMovies(client))

	//Route that creates and insert one user to users collection in DB
	router.POST("/register", controller.RegisterUser(client))

	//Route that logins/authenticate user and save its data (including generated tokens) on db
	router.POST("/login", controller.LoginUser(client))

	//Route that returns all genres form genres collection in mongodb
	router.GET("/genres", controller.GetGenres(client))

	//Route that logouts a user
	router.POST("/logout", controller.LogoutHandler(client))

	//Route thar refreshes tokens
	router.POST("/refresh", controller.RefreshTokenHandler(client))

}
