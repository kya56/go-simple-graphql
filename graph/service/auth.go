package service

import (
	"context"
	"errors"
	"go-simple-graphql/database"
	"go-simple-graphql/graph/model"
	"go-simple-graphql/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

const user_collection = "users"

func (userService *UserService) Register(input model.NewUser) string {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	var password string = string(encryptedPass)

	var newUser model.User = model.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	var collection *mongo.Collection = database.GetCollection(user_collection)

	res, err := collection.InsertOne(context.TODO(), newUser)

	if err != nil {
		return ""
	}

	var userId string = res.InsertedID.(primitive.ObjectID).Hex()

	token, err := utils.GenerateAccessToken(userId)

	if err != nil {
		return ""
	}

	return token
}

func (userService *UserService) Login(input model.LoginInput) string {
	var collection *mongo.Collection = database.GetCollection(user_collection)

	var user *model.User = &model.User{}
	filter := bson.M{"email": input.Email}

	var res *mongo.SingleResult = collection.FindOne(context.TODO(), filter)

	if err := res.Decode(user); err != nil {
		return ""
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		return ""
	}

	token, err := utils.GenerateAccessToken(user.ID)

	if err != nil {
		return ""
	}

	return token
}

func (userService *UserService) GetUser(id string) (*model.User, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &model.User{}, errors.New("id is invalid")
	}

	var query primitive.D = bson.D{{Key: "_id", Value: userID}}
	var collection *mongo.Collection = database.GetCollection(user_collection)

	var userData *mongo.SingleResult = collection.FindOne(context.TODO(), query)

	if userData.Err() != nil {
		return &model.User{}, errors.New("user not found")
	}

	var user *model.User = &model.User{}
	userData.Decode(user)

	return user, nil
}
