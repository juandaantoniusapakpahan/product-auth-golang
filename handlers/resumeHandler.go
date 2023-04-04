package handlers

import (
	"io/ioutil"
	"net/http"
	"product-auth/database"
	"product-auth/entity"
	"product-auth/tokenCookie"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		/** Entity AddResume Here **/
		err, statusCode := entity.AddResume(u.Name, u.Email, payload)
		if err != nil {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return c.JSON(statusCode, M{"status": "success"})
	}
}

func EditResumeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		claims := token.Claims.(*tokenCookie.Claims)

		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token was valid")
		}
		payloadUpdate, err := ioutil.ReadAll(c.Request().Body)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		statusCode, err := entity.EditResume(claims.Email, claims.Name, c.Param("id"), payloadUpdate)

		if err != nil {
			return echo.NewHTTPError(statusCode, err.Error())
		}

		return c.JSON(http.StatusOK, M{"status": "success"})
	}
}

func GetResumeByPrivateHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
		}
		claims := token.Claims.(*tokenCookie.Claims)
		resumeResult, statusCode, err := entity.GetResume(claims.Name, claims.Email)
		if err != nil {
			return echo.NewHTTPError(statusCode, err)
		}
		return c.JSON(statusCode, resumeResult)
	}
}

func DeleteResumeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := c.Get("user").(*jwt.Token)
		id := c.Param("id")
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
		}
		claims := token.Claims.(*tokenCookie.Claims)
		err, user := entity.CheckUser(claims.Name, claims.Email)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		db, err := database.Connect()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		resumeId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		userId, err := primitive.ObjectIDFromHex(user.ID.Hex())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		filter := bson.M{}
		filter["_id"] = resumeId
		filter["owner._id"] = userId
		_, err = db.Collection("resumes").DeleteOne(ctx, filter)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return c.JSON(http.StatusOK, M{"status": "success"})
	}
}
