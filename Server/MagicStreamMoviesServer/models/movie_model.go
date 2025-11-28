/*
	The package declaration defines the package name to which the file belongs

A package in go is a way to group related files and functions together. Files that belong
to the same package can share functions, types and variables

In Go, a package is a fundamental unit for organizing and structuring code.

	It serves as a collection of related Go source files that reside within the same directory.
*/
package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=100"`
}

type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"required"`
}

type Movie struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"_id"`                             //Unique identifier for movie document on DB
	ImdbID      string        `bson:"imdb_id" json:"imdb_id" validate:"required"`           //Unique identifier from IMDB for a specific movie, serie, actor.. Validate "required" is used to validate that this parameter os not empty
	Title       string        `bson:"title" json:"title" validate:"required,min=2,max=500"` //Required is to check that the parameter is not empty and min 2 it to check that the title is at least 2 characters
	PosterPath  string        `bson:"poster_path" json:"poster_path" valdiate:"required,url"`
	YoutubeID   string        `bson:"youtube_id" json:"youtube_id" validate:"required"` //Youtube URL of movie trailer
	Genre       []Genre       `bson:"genre" json:"genre" validate:"required,dive"`      //Movie can have one or more genre
	AdminReview string        `bson:"admin_review" json:"admin_review"`                 //Review of the movie
	Ranking     Ranking       `bson:"ranking" json:"ranking" validate:"required"`
}

/*
Define how we want our fields to look in BSON and JSON format. Bson refers to how the field looks within the relevant
document in the db and JSON is how the field would look when the data is returned from our web API application to calling
client code.

So we can define in our structs how the fields map to our MongoDB as well as to the JSON data will be sent to calling client code
*/
