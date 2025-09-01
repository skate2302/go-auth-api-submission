package handlers

import (
	"context"
	"net/http"
	"time"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Skate2302/go-auth-api/internal/models"
)

// new validator instance
var validate = validator.New()

type SignUpInput struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

// SignUp godoc
// @Summary      Register a new user
// @Description  Create a new user account with name, email, and password.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      SignUpInput  true  "Sign Up User"
// @Success      201   {object}  models.User
// @Failure      400   {object}  object{error=string}
// @Failure      409   {object}  object{error=string}
// @Failure      500   {object}  object{error=string}
// @Router       /signup [post]
func SignUp(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var input SignUpInput 

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(&input); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		collection := client.Database("go-auth-db").Collection("users")
		count, err := collection.CountDocuments(ctx, bson.M{"email": input.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for existing user"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "A user with this email already exists"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		newUser := models.User{
			ID:        primitive.NewObjectID(),
			Name:      input.Name,
			Email:     input.Email,
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, insertErr := collection.InsertOne(ctx, newUser)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, newUser)
	}
}

// Login godoc
// @Summary      Log in a user
// @Description  Authenticate a user with email and password to get a JWT token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      LoginInput  true  "Login User"
// @Success      200   {object}  object{token=string}
// @Failure      400   {object}  object{error=string}
// @Failure      401   {object}  object{error=string}
// @Failure      500   {object}  object{error=string}
// @Router       /login [post]
func Login(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var input LoginInput
		var user models.User
		collection := client.Database("go-auth-db").Collection("users")

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(&input); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error: JWT secret not set"})
			return
		}

		claims := jwt.MapClaims{
			"id":  user.ID,
			"email": user.Email,
			"exp": time.Now().Add(time.Hour * 24).Unix(), // Token is valid for 24 hours
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create authentication token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}