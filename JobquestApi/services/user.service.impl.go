package services

import (
	"JobquestApi/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx context.Context
}

func NewUserServiceImpl(ctx context.Context, usercollection *mongo.Collection) UserService {
	return &UserServiceImpl {
		userCollection: usercollection,
		ctx: ctx,
	}
}


func (u *UserServiceImpl) CreateUser(user *models.User) error {
	_, err := u.userCollection.InsertOne(u.ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServiceImpl) GetUser(id *string) (*models.User, error) {
	var user *models.User
	query := bson.D{{Key: "_id", Value: *id}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}


func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	opts := options.Find().SetProjection(bson.M{"_id": 0}).SetNoCursorTimeout(true)
	cursor, err := u.userCollection.Find(u.ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(u.ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserServiceImpl) UpdateUser(email string, user *models.User) (*models.User, error) {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"first_name": user.FirstName, "last_name": user.LastName}}
	opts := options.Update().SetUpsert(false)
	res, err := u.userCollection.UpdateOne(u.ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	if res.MatchedCount < 1 {
		return nil, errors.New("no user found")
	}
	return user, nil
}


func (u *UserServiceImpl) DeleteUser(email *string) error {
	filter := bson.M{"email": *email}
	res, err := u.userCollection.DeleteOne(u.ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
