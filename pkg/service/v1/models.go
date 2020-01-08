package v1

import (
  // v1 "github.com/ckbball/smurfin-checkout/pkg/api/v1"
  //"go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  //"go.mongodb.org/mongo-driver/mongo"
)

type AccountPurchased struct {
  PurchaseDate         int64  `json:"purchase_date"`
  AccountLoginName     string `json:"account_login_name"`
  AccountLoginPassword string `json:"account_login_password"`
  AccountEmail         string `json:"account_email"`
  AccountEmailPassword string `json:"account_email_password"`
  AccountId            string `json:"account_id"`
  VendorId             string `json:"vendor_id"`
  BuyerId              string `json:"buyer_id"`
  BuyerEmail           string `json:"buyer_email"`
}

type User struct {
  Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Email      string             `json:"email,omitempty" bson:"email,omitempty"`
  Password   string             `json:"password,omitempty" bson:"password,omitempty"`
  Username   string             `json:"username,omitempty" bson:"username,omitempty"`
  LastActive int                `json:"lastActive,omitempty" bson:"lastActive,omitempty"`
  Experience string             `json:"experience,omitempty" bson:"experience,omitempty"`
  Languages  []string           `json:"languages,omitempty" bson:"languages,omitempty"`
}
