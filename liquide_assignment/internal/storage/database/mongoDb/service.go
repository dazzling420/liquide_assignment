package mongodb

import (
	"liquide_assignment/internal/service/login"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *MongoRepository) AddUser(signupRequest login.SignupRequest) error {
	collection := s.db.Database("liquide").Collection("users")
	_, err := collection.InsertOne(s.ctx, signupRequest)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoRepository) GetUser(loginRequest login.LoginRequest) (login.SignupRequest, error) {
	collection := s.db.Database("liquide").Collection("users")
	var result login.SignupRequest
	filter := bson.M{"email": loginRequest.Email}
	err := collection.FindOne(s.ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *MongoRepository) DoesUserExist(pan string) bool {
	collection := s.db.Database("liquide").Collection("users")
	var result login.SignupRequest
	filter := bson.M{"pan": pan}
	err := collection.FindOne(s.ctx, filter).Decode(&result)
	if err != nil {
		return false
	}
	if result == (login.SignupRequest{}) {
		return false
	}
	return true
}
