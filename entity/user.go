package entity

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"product-auth/database"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var ctx = context.Background()

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name     string             `json:"name" bson:"name" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required, email"`
	Password string             `json:"password" bson:"password" validate:"required"`
	Role     string             `json:"role" bson:"role" default:"user"`
}

func AddUser(c echo.Context) (error, int) {
	db, err := database.Connect()
	body, err := ioutil.ReadAll(c.Request().Body)
	var userData User

	err = json.Unmarshal([]byte(body), &userData)
	if err != nil {
		return err, http.StatusBadRequest
	}

	passworHash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 10)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	if userData.Role == "" {
		userData.Role = "user"
	}

	userData.Password = string(passworHash)

	err = PossibleDoubleEmail(userData.Email)
	if err != nil {
		return err, http.StatusBadRequest
	}
	_, err = db.Collection("users").InsertOne(ctx, userData)

	if err != nil {
		return err, http.StatusBadRequest
	}
	return nil, 0
}

func PossibleDoubleEmail(email string) error {
	db, err := database.Connect()
	data := new(User)

	if err != nil {
		return err
	}
	err = db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&data)

	if err == nil {
		return errors.New("Email has been regiester")
	}
	return nil
}

func SignIn(user User) (*User, error, int) {
	db, err := database.Connect()
	var data User
	if err != nil {
		return &User{}, err, http.StatusInternalServerError
	}

	err = db.Collection("users").FindOne(ctx, bson.M{"email": user.Email}).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &data, err, http.StatusNotFound
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(user.Password)); err != nil {
		return &User{}, err, http.StatusUnauthorized
	}

	return &data, nil, 0
}

func CheckUser(name string, email string) (error, *User) {
	db, err := database.Connect()
	if err != nil {
		return err, &User{}
	}

	u := new(User)

	filter := bson.M{"name": name, "email": email}
	err = db.Collection("users").FindOne(ctx, filter).Decode(&u)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err, &User{}
		}
		panic(err)
	}

	return nil, u
}
