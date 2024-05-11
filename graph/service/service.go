package service

import (
	"context"
	"errors"
	"go-simple-graphql/database"
	"go-simple-graphql/graph/model"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogService struct{}

const blog_collection = "blogs"

func (b *BlogService) GetAllBlogs() []*model.Blog {
	var query primitive.D = bson.D{{}}
	var findOptions *options.FindOptions = options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := database.GetCollection(blog_collection).Find(context.TODO(), query, findOptions)
	if err != nil {
		return []*model.Blog{}
	}

	var blogs []*model.Blog = make([]*model.Blog, 0)
	if err := cursor.All(context.TODO(), &blogs); err != nil {
		return []*model.Blog{}
	}

	return blogs
}

func (b *BlogService) GetBlogByID(id string) (*model.Blog, error) {
	blogId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &model.Blog{}, errors.New("id is invalid")
	}

	var query primitive.D = bson.D{{Key: "_id", Value: blogId}}

	var collection *mongo.Collection = database.GetCollection(blog_collection)
	var blogData *mongo.SingleResult = collection.FindOne(context.TODO(), query)

	if blogData.Err() != nil {
		return &model.Blog{}, errors.New("blog not found")
	}

	var blog *model.Blog = &model.Blog{}
	blogData.Decode(blog)

	return blog, nil
}

func (b *BlogService) CreateBlog(input model.NewBlog, user model.User) (*model.Blog, error) {
	var blog model.Blog = model.Blog{
		ID:        uuid.New().String(),
		Title:     input.Title,
		Content:   input.Content,
		Author:    &user,
		CreatedAt: time.Now(),
	}

	var collection = database.GetCollection(blog_collection)
	result, err := collection.InsertOne(context.TODO(), blog)
	if err != nil {
		return &model.Blog{}, errors.New("Failed to create blog")
	}

	var filter primitive.D = bson.D{{Key: "_id", Value: result.InsertedID}}

	var created *mongo.SingleResult = collection.FindOne(context.TODO(), filter)
	var newBlog *model.Blog = &model.Blog{}
	created.Decode(&newBlog)

	return newBlog, nil
}

func (b *BlogService) EditBlog(input model.EditBlog, user model.User) (*model.Blog, error) {
	blogID, err := primitive.ObjectIDFromHex(input.BlogID)

	if err != nil {
		return &model.Blog{}, errors.New("id is invalid")
	}

	var query primitive.D = bson.D{
		{Key: "_id", Value: blogID},
		{Key: "author.id", Value: user.ID},
	}

	var update primitive.D = bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "title", Value: input.Title},
			{Key: "content", Value: input.Content},
			{Key: "updatedAt", Value: time.Now()},
		},
	}}

	var collection *mongo.Collection = database.GetCollection(blog_collection)

	var updateResult *mongo.SingleResult = collection.FindOneAndUpdate(
		context.TODO(),
		query,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if updateResult.Err() != nil {
		return &model.Blog{}, errors.New("blog not found")
	}

	var editedBlog *model.Blog = &model.Blog{}

	updateResult.Decode(editedBlog)

	return editedBlog, nil
}

func (b *BlogService) DeleteBlog(input model.DeleteBlog, user model.User) bool {
	blogID, err := primitive.ObjectIDFromHex(input.BlogID)

	if err != nil {
		return false
	}

	var query primitive.D = bson.D{
		{Key: "_id", Value: blogID},
		{Key: "author.id", Value: user.ID},
	}

	var collection *mongo.Collection = database.GetCollection(blog_collection)

	result, err := collection.DeleteOne(context.TODO(), query)

	return !(err != nil || result.DeletedCount < 1)
}
