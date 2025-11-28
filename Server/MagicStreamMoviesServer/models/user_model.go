package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID              bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID          string        `json:"user_id" bson:"user_id"`
	FirstName       string        `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName        string        `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	Email           string        `json:"email" bson:"email" validate:"required,email"`
	Password        string        `json:"password" bson:"password" validate:"required,min=6"`
	Role            string        `json:"role" bson:"role" validate:"oneof=ADMIN USER"`
	CreatedAt       time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" bson:"updated_at"`
	Token           string        `json:"token" bson:"token"`
	RefreshToken    string        `json:"refresh_token" bson:"refresh_token"`
	FavouriteGenres []Genre       `json:"favourite_genres" bson:"favourite_genres" validate:"required,dive"`
}

/*
JWT for authentication and authorization. Compact URL safe token use to securelyt transmit information between parties.
It has 3 parts separated by dots, a header part which specifies the type of token and the signing algorithm for example HS, 256

Payload, contains the claim user data& meta data

Signature verifies the token authenticity using a secret key or public private key

The user logs in with credential, the server verifies it and generates the token jwt, then it is sent back to the client to be
stored in local storage on the client or a cookie. (Local storage not recommended as it can be accessed by client Javascript code
and less susceptible to XSS attacks), for future request, the jwt is included on the authorization header.

Benefits:
Stateles: No session storage needed on the server
Scalable: Easily used in distributed systems

Refresh token is done so the user does not need to log again when token expiracy is reached
*/

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type UserResponse struct {
	UserId          string  `json:"user_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	Role            string  `json:"role"`
	Token           string  `json:"token"`
	RefreshToken    string  `json:"refresh_token"`
	FavouriteGenres []Genre `json:"favourite_genres"`
}

/*
UserResponse is a DTO which stands for Data Transfer Object. It's a design pattern used to transfer data
between software application layers such as between backend and fronted or between services in a distributed system
Key characterisitcs:
-It holds data only: contain fields, properties, getters, setters, it does not contain business logic
-Flat structure: Usually simpler and flatter than domain models,
-Serializable: Designed to be easily serializable to JSON, XML, etc. for transmission over the network.const

Common cases are:
API responses
Returning data to client without exposing the internal domain model (in this case User),
Data aggregation, combining data from multiple sources into a single object
Input Validation: Accepting structured input from external sources for example form submissions+

By using DTO we're only exposing data that needs to be exposed to the client
*/
