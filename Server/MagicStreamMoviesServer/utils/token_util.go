//File containing code to generate access tokens

package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Fernando0743/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	jwt "github.com/golang-jwt/jwt/v5"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

// Read Secret Key
var SECRET_KEY string = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY string = os.Getenv("SECRET_REFRESH_KEY")

// Function that generates all tokens. (Access and Refresh Tokens)
func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	//CODE FOR CREATING ACCESS TOKEN
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	//Create and sign token. HS256 is the signing method to verify token authenticity
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	//CODE FOR CREATING REFRESH TOKEN

	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}
	//Create and sign token. ES256 is the signing method to verify token authenticity
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))

	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

/*
jwt.RegisteredClaims is a struct provided by "github.com/golang-jwt/jwt/v5". It contains a set of standard claims defined by the
JWT specification RFC5719. These fields include:

-Issuer: Who issued the token
-Subject: Subject of the token (usually a userid)
-Audience: Intended recipients of the token
-ExpiresAt: When the token should expire
-NotBefore: When the token becomes valid
-IssuedAt: When the token was issued
-ID: Unique identifier for the token (Can be used to prevent replay attacks)
*/

// Update All tokens on db
func UpdateAllTokens(userId, token, refreshToken string, client *mongo.Client) (err error) {
	//Query with time out for inserting user into db
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	//Format time of update
	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	//Asign relevant values to database fields
	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    updateAt,
		},
	}

	// Get collection
	var userCollection *mongo.Collection = database.OpenCollection("users", client)

	//Update document with updated data on db filtering the document with userId
	result, err := userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)

	if err != nil {
		return err
	}

	// Optional: Check if document was actually found and updated
	if result.MatchedCount == 0 {
		return fmt.Errorf("no user found with user_id: %s", userId)
	}

	return nil
}

// Extract the token from the header of the HTTP request
func GetAccessToken(c *gin.Context) (string, error) {
	//Get authorization header
	//authHeader := c.Request.Header.Get("Authorization")

	//if authHeader == "" {
	//return "", errors.New("authorization Header is required")
	//}
	//Read token String
	//tokenString := authHeader[len("Bearer "):]

	//if tokenString == "" {
	//	return "", errors.New("bearer token is required")
	//}
	//Read token from cookies
	tokenString, err := c.Cookie("access_token")

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Function that will validate tokens
func ValidateToken(tokenString string) (*SignedDetails, error) {
	//Signed Details is where token claims will be decoded into because one of its elements is registeredclaims
	claims := &SignedDetails{}

	//Parse Token String to Token object, decode its claims into claims variable
	//Uses a callback function to provide the secret key used to verify the signature, we return it as a bite slice
	//which is required
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	//Check that signed algorithm is HMAC like HS256. This is a critical security step.
	//Attackers can try to spoof token using a different algorithm if we don't check this
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	//Check token expiracy
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil

}

// Function that returns userID from http request
func GetUserIdFromContext(c *gin.Context) (string, error) {
	//Get user id from current context
	userId, exists := c.Get("userId")

	//Error when trying to fetch userId from context
	if !exists {
		return "", errors.New("userId does not exist in this context")
	}

	//Convert user id to string
	id, ok := userId.(string)

	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return id, nil
}

// Function that returns Role from http request
func GetRoleFromContext(c *gin.Context) (string, error) {
	//Get user id from current context
	role, exists := c.Get("role")

	//Error when trying to fetch userId from context
	if !exists {
		return "", errors.New("role does not exist in this context")
	}

	//Convert user id to string
	memberRole, ok := role.(string)

	if !ok {
		return "", errors.New("role to retrieve userId")
	}

	return memberRole, nil
}

func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return []byte(SECRET_REFRESH_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
