package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Fernando0743/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/Fernando0743/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/Fernando0743/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/tmc/langchaingo/llms/openai"
)

// Validator
var validate = validator.New()

// Function that returns type Gin Handler Function.
// gin.Context is a central and crucial component that encapsulates the entire HTTP request and response cycle for a single request.
// It can read and create HTTP responses
// It acts as a bridge between the Gin framework and your application's handlers and middleware.
// The function begins with capital letter so it is exportable we can acces it from main method once controllers package is imported
func GetMovies(client *mongo.Client) gin.HandlerFunc {
	//c is context of http request, ctx is request of database operation/query
	return func(c *gin.Context) {
		//Execute Query with Timeout of 100 seconds (these 2 lines are for memory management)
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// Get collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		//Create array of type Movie struct defined on models package and inside movie_model file
		var movies []models.Movie

		//Cursor that queries the db
		cursor, err := movieCollection.Find(ctx, bson.M{})

		//Error that occurs when we can't fecth movies from db
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies."})
		}
		//Line for memory management
		defer cursor.Close(ctx)

		//Error that occurs when we can covert movies from db into movies array structure elements
		if err := cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies."})
		}

		c.JSON(http.StatusOK, movies)
	}
}

// Function that returns a single movie from DB given IMDB_ID
func GetMovie(client *mongo.Client) gin.HandlerFunc {
	//c is context of http request, ctx is request of database operation/query
	return func(c *gin.Context) {
		//Create a new context ctx that cancells after 100 seconds.
		//This functions returns this new ctx which is the context that carries the timeout
		//Cancel is a function to manually cancel the context before the timeout if necessary
		//Defer delays the execution of cancel until the surrounding function returns. This ensures that context is properly cleaned up and resources are released even
		// if the function exits early due to an error
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		//Get IMDB ID
		movieID := c.Param("imdb_id")

		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		// Get collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		var movie models.Movie

		//Fetch movie from db and store it in movie var
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)

	}
}

func AddMovie(client *mongo.Client) gin.HandlerFunc {
	//c is context of http request, ctx is request of database operation/query
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		// Get collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		var movie models.Movie

		//Code that checks if an error occurs when attempting to pass the incoming JSON payload from the request body
		//into the movie struct.
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
			return
		}

		//Use go playground validator to validate the input received accodring to the declaratives rules from our Movie Model
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
			return
		}

		//Insert input into movie collection in db
		result, err := movieCollection.InsertOne(ctx, movie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
			return
		}

		c.JSON(http.StatusCreated, result)

	}
}

// Function to update movie review on db
func AdminReviewUpdate(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Check if user has role admin

		role, err := utils.GetRoleFromContext(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found in context"})
			return
		}

		if role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be part of the admin role"})
			return
		}

		//Get movieId from route
		movieId := c.Param("imdb_id")

		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie ID is required"})
			return
		}

		//Information from client its only the review for the movie
		var req struct {
			AdminReview string `json:"admin_review"`
		}
		//"admin_review"  : "Clint Eastwood as always was magnificent. What an amazing cast and movie"
		//Response from the admin is it's review and the ranking name coming from OpenAI sentiment analysis
		var resp struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		//Get response from AI given movie review
		sentiment, rankVal, err := GetReviewRanking(req.AdminReview, client, c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking", "details": err.Error()})
			return
		}

		//Filter to find movie document in movies collection
		filter := bson.M{"imdb_id": movieId}

		//Bson (Database) content to be updated on document
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_value": rankVal,
					"ranking_name":  sentiment,
				},
			},
		}

		var ctx, cancel = context.WithTimeout(c, time.Second*100)
		defer cancel()

		// Get collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		result, err := movieCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		resp.RankingName = sentiment
		resp.AdminReview = req.AdminReview

		c.JSON(http.StatusOK, resp)

	}
}

// Function that prompts AI given a movie review and returns the ranking name and ranking vale and an error if occurs
func GetReviewRanking(admin_review string, client *mongo.Client, c *gin.Context) (string, int, error) {
	//Get rakings collection from db
	rankings, err := GetRankings(client, c)

	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	//Get all ranking names in one string
	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited = sentimentDelimited + ranking.RankingName + ","
		}
	}

	//Format
	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	//Load env file and prompt variable
	err = godotenv.Load(".env")

	if err != nil {
		log.Println("Warning: .env file not found")
	}

	//Get open ai api key from env
	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")

	if OpenAiApiKey == "" {
		return "", 0, errors.New("could not read open ai api key")
	}

	//Sign to open ai with our key
	llm, err := openai.New(openai.WithToken(OpenAiApiKey))

	if err != nil {
		return "", 0, err
	}

	//Get base prompt from env
	base_prompt_template := os.Getenv("BASE_PROMPT_TEMPLATE")

	//include rankings in base prompt
	base_prompt := strings.Replace(base_prompt_template, "{rankings}", sentimentDelimited, 1)

	//pass prompt to llm and concatenate movie review to base_prompt var
	response, err := llm.Call(context.Background(), base_prompt+admin_review)

	if err != nil {
		return "", 0, err
	}
	//Get ranking from prompt
	rank_val := 0

	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rank_val = ranking.RankingValue
			break
		}
	}

	return response, rank_val, nil

}

// Function that queries and returns all rankings from rankings collection from db
func GetRankings(client *mongo.Client, c *gin.Context) ([]models.Ranking, error) {
	var rankings []models.Ranking

	ctx, cancel := context.WithTimeout(c, time.Second*100)
	defer cancel()

	var rankingsCollection *mongo.Collection = database.OpenCollection("rankings", client)

	//Cursor that queries the rankings collection
	cursor, err := rankingsCollection.Find(ctx, bson.M{})

	//Error that occurs when we can't fecth rankings from db
	if err != nil {
		return nil, err
	}
	//Line for memory management
	defer cursor.Close(ctx)

	//Error that occurs when we can covert rankings from db into rankings array structure elements
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, err
}

// Function that queries and returns recommended movies for a user
func GetRecommendedMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Get userID from current context
		userId, err := utils.GetUserIdFromContext(c)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User Id not found in context"})
			return
		}

		//Query db to get user Favourite Genres
		favourite_genres, err := GetUsersFavouriteGenres(userId, client, c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = godotenv.Load("env")
		if err != nil {
			log.Println("Warning: .env file not found")
		}

		var recommendedMovieLimitVal int64 = 5

		recommendedMoviesLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")

		//Convert RecommendedMoviesLimitStr string to int
		if recommendedMoviesLimitStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMoviesLimitStr, 10, 64)
		}

		//Query 5 movies given user favourite genres and ranking values
		findOptions := options.Find()
		//Sort movies by ranking value
		findOptions.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})

		findOptions.SetLimit(recommendedMovieLimitVal)

		//Filter only by favourite genres from user
		filter := bson.M{"genre.genre_name": bson.M{"$in": favourite_genres}}

		var ctx, cancel = context.WithTimeout(c, time.Second*100)
		defer cancel()

		// Get collection
		var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

		cursor, err := movieCollection.Find(ctx, filter, findOptions)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}

		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie

		//Error that occurs when we can covert query from db into recommendedMovies array structure elements
		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recommendedMovies)

	}
}

// Function that returns a user favourite genres given user id
func GetUsersFavouriteGenres(userId string, client *mongo.Client, c *gin.Context) ([]string, error) {
	//Query on DB to find user document given userId
	var ctx, cancel = context.WithTimeout(c, time.Second*100)
	defer cancel()

	filter := bson.M{"user_id": userId}

	//Structure that will be returned when querying for document on db
	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}
	//Set projection in options mongodb
	opts := options.FindOne().SetProjection(projection)

	var result bson.M

	// Get collection
	var userCollection *mongo.Collection = database.OpenCollection("users", client)

	//Query document and store it in result variable with projection structure
	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)

	if err != nil {
		//Document does not exist on db
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
	}

	//Store favourite_genres array
	favGenresArray, ok := result["favourite_genres"].(bson.A)

	if !ok {
		return []string{}, errors.New("unable to retrieve favourite genres for user")
	}

	var genreNames []string

	//Iterate thorugh favourite_genres array and only extract the genre name since on db every genre item has the genre name and genre id
	for _, item := range favGenresArray {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil

}

// Function that returns genres collection
func GetGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, time.Second*100)
		defer cancel()

		var genres []models.Genre

		var genreCollection *mongo.Collection = database.OpenCollection("genres", client)

		cursor, err := genreCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movie genres"})
			return
		}
		//Free resources if query fails
		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, genres)

	}
}
