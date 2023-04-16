package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"product-auth/entity"
	"product-auth/tokenCookie"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
)

type M map[string]interface{}

var ctx = context.Background()

func SignUpHandler() echo.HandlerFunc {
	return func(c echo.Context) error {

		err, statusCode := entity.AddUser(c)
		if err != nil {
			return echo.NewHTTPError(statusCode, err.Error())
		}

		return c.JSON(http.StatusCreated, M{"status": "success", "message": "Success added user"})
	}
}

func SignInHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var userPayload entity.User
		payload, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := json.Unmarshal([]byte(payload), &userPayload); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		data, err, statusCode := entity.SignIn(userPayload)
		if err != nil {
			return echo.NewHTTPError(statusCode, err)
		}

		if err := tokenCookie.GenerateTokenAndSetCookie(data, c); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(M{"success": "true"})
	}

}
