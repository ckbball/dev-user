package v1

import (
  //"go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  //"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
  Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Email      string             `json:"email,omitempty" bson:"email,omitempty"`
  Password   string             `json:"password,omitempty" bson:"password,omitempty"`
  Username   string             `json:"username,omitempty" bson:"username,omitempty"`
  LastActive int                `json:"lastActive,omitempty" bson:"last_active,omitempty"`
  Experience string             `json:"experience,omitempty" bson:"experience,omitempty"`
  Languages  []string           `json:"languages,omitempty" bson:"languages,omitempty"`
}
