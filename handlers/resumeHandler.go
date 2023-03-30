package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"product-auth/database"
	"product-auth/entity"
	"product-auth/tokenCookie"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2/bson"
)

func CreateResumHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, ok)
		}
		u := token.Claims.(*tokenCookie.Claims)
		payload, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		rsm := new(entity.Resume)
		err = json.Unmarshal(payload, &rsm)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		rsm.UserId = bson.ObjectId(u.ID)
		db, err := database.Connect()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		_, err = db.Collection("resumes").InsertOne(ctx, rsm)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusCreated, M{"status": "success", "message": "Success added resume"})
	}
}
