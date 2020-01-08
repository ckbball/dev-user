package v1

import (
  "context"
  //"errors"
  v1 "github.com/ckbball/dev-user/pkg/api/v1"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  //"go.mongodb.org/mongo-driver/mongo/readpref"
)

type repository interface {
  Create(*v1.User) (string, error)
  Update(*v1.User, string) (int64, int64, error)
  Delete(string) (int64, error)
  GetById(string) (*User, error)
  FilterUsers(*v1.FindRequest) ([]*User, error)
}

type UserRepository struct {
  ds *mongo.Collection //
}

func NewUserRepository(client *mongo.Collection) *UserRepository {
  return &UserRepository{
    ds: client,
  }
}

func (s *UserRepository) makeFilter(fieldMap map[string]string) *bson.D {
  output := bson.D{}

  for key, value := range fieldMap {
    if key == "id" {
      filter := bson.E{
        Key: "_id",
        Value: &bson.A{
          "$eq",
          value,
        },
      }
      output = append(output, filter)
    }

    if key == "language" {
      filter := bson.E{
        Key: "Languages",
        Value: &bson.A{
          "$in",
          value,
        },
      }
      output = append(output, filter)
    }

    if key == "experience" {
      filter := bson.E{
        Key: "Experience",
        Value: &bson.A{
          "$eq",
          value,
        },
      }
      output = append(output, filter)
    }
  }
  return &output
}

func (s *UserRepository) GetById(id string) (*User, error) {
  primitiveId, _ := primitive.ObjectIDFromHex(id)

  var user User
  err := s.ds.FindOne(context.TODO(), User{Id: primitiveId}).Decode(&user)
  if err != nil {
    return nil, err
  }

  return &user, nil
}

func (repository *UserRepository) Create(user *v1.User) (string, error) {
  // add a duplicate email and a duplicate username check

  insertUser := bson.D{
    {"email", user.Email},
    {"password", user.Password},
    {"username", user.Username},
    {"last_active", user.LastActive},
    {"experience", user.Experience},
    {"languages", user.Languages},
  }

  result, err := repository.ds.InsertOne(context.TODO(), insertUser)

  if err != nil {
    return "", err
  }

  id := result.InsertedID
  w, _ := id.(primitive.ObjectID)

  out := w.Hex()

  return out, err

}

func (repository *UserRepository) Update(user *v1.User, id string) (int64, int64, error) {
  // add a duplicate email and a duplicate username check

  primitiveId, _ := primitive.ObjectIDFromHex(id)

  result, err := repository.ds.UpdateOne(context.Background(),
    bson.D{
      {"_id", primitiveId},
    },
    bson.D{
      {"$set", bson.D{
        {"email", user.Email},
        {"password", user.Password},
        {"username", user.Username},
        {"last_active", user.LastActive},
        {"experience", user.Experience},
        {"languages", user.Languages},
        // in the future add other fields
      }},
    },
  )

  if err != nil {
    return -1, -1, err
  }

  return result.MatchedCount, result.ModifiedCount, nil
}

func (repository *UserRepository) Delete(id string) (int64, error) {
  primitiveId, _ := primitive.ObjectIDFromHex(id)
  filter := bson.D{{"_id", primitiveId}}

  result, err := repository.ds.DeleteOne(context.Background(), filter)
  if err != nil {
    return -1, err
  }
  return result.DeletedCount, nil
}

func (s *UserRepository) FilterUsers(req *v1.FindRequest) ([]*User, error) {

  fieldMap := map[string]string{"language": req.Language, "experience": req.Experience}

  // make bson filter object here by calling makeFilter()
  filter := s.makeFilter(fieldMap)

  findOptions := options.Find()
  findOptions.SetLimit(int64(req.Limit))
  findOptions.SetSort(bson.D{{"_id", -1}})
  findOptions.SetSkip(int64(req.Page))

  var users []*User
  cur, err := s.ds.Find(context.TODO(), filter, findOptions)
  if err != nil {
    return nil, err
  }
  defer cur.Close(context.TODO())

  for cur.Next(context.TODO()) {
    var elem *User
    err := cur.Decode(&elem)
    if err != nil {
      return nil, err
    }

    users = append(users, elem)
  }

  if err := cur.Err(); err != nil {
    return users, err
  }

  return users, nil
}
