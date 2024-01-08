package controllers

import (
	"context"
	"fmt"
	"go/token"
	"github.com/mtoha/akhil/database"
	//"golang_jwt_akhil_sharma/helpers"
	helper "github.com/mtoha/akhil/helpers"
	"github.com/mtoha/akhil/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email of password is incorrect")

		check = false
	}

	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := validate.Struct(user)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": validationError.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.email})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error occured when Checking Count Email": err.Error()})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "This Email already useed"})
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.phone})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error occured when Checking Count phone": err.Error()})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "This Phone already useed"})
		}

		user.cratedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.updatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.user_id = user.ID.Hex()

		token, refreshToken, _ := helper.GenerateAllTokens(*user.email, *user.firstName, *user.lastName, *user.userType, *&user.user_id)

		user.token = &token
		user.refreshToken = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}
		defer cancel()
		c.JSON(http.Status0k, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{Herroc: err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.password, *foundUser.password)
		defer cancel()

		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}

		helper.GenerateAllTokens(*foundUser.email)
	}
}

func GetUsers() gin.HandlerFunc{ 
	return func(c *gin.Context){ 
		helper.CheckUserType(c, "ADMIN"); err != nil { 
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout (context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage <1{ 
			recordPerPage = 10
		
		} 

		page, errl := strconv.Atoi(c.Query("page"))
		if errl !=nil || page<1{
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}},
			{"total_count", bson.D{{"$sum", 1}}}, 
			{"data", bson.D{{"push", "$$ROOT"}}} 
		}}}
		
		projectStage := bson.D{
			{"$project", bson.D{
			{"_id", 0}, 
			{"total_count", 1}, 
			{"user items", bson.D{{"$slice", []Jinterface{}{"$data", startIndex, recordPerPage}}}} 
			
			}}}

			result,err := userCollection.Aggregate(ctx, mongo.Pipeline{ 
				matchStage, groupStage, projectStage 
			})
			defer cancel()
			
			if err!= nil{
				c.JSON{http.StatusInternalServerError, gin.H{"error":"error occured while listing user items"}
			}
			
			
			var allUsers []bson.M
			if err = result.All(ctx, &allusers); err !=nil{
				log.Fatal(err)
			}

			c.JSON(http.StatusOK, allusers[0])
		}
	}
}
		

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUld(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.findOne(ctx, bson.M{"user_id": userId}).decode(&user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
