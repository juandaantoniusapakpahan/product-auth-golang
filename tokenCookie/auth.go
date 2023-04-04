package tokenCookie

import (
	"net/http"
	"product-auth/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	accessTokenName  = "access-token"
	tokenKeySecret   = "a$10$hkVz03JAN0WjDoYWaWJohuxiyIjVoVuEOuaEGKVcwkJcbaObf2zYO"
	refreshTokenName = "refresh-token"
	refreshTokenKey  = "2a$10$HEDbL6/TLfkFGjnyhGgSxOLt.jvSb9BPiTIwm07JbVuhzXf9XWgZ"
)

type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GetTokenKeySecret() string {
	return tokenKeySecret
}

func GetRefreshKeySecret() string {
	return refreshTokenKey
}
func GenerateTokenAndSetCookie(user *entity.User, c echo.Context) error {
	accessToken, expir, err := generateAccessToken(user)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenName, accessToken, expir, c)
	setUserCookie(user, expir, c)

	refreshToken, expi, err := generateRefreshToken(user)
	if err != nil {
		return err
	}
	setTokenCookie(refreshTokenName, refreshToken, expi, c)

	return nil
}

func generateAccessToken(user *entity.User) (string, time.Time, error) {
	expiredTime := time.Now().Add(1 * time.Hour)
	return generateToken(user, expiredTime, []byte(GetTokenKeySecret()))
}

func generateToken(user *entity.User, expr time.Time, secret []byte) (string, time.Time, error) {

	claims := &Claims{
		Name:  user.Name,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expr.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expr, nil
}

func setTokenCookie(tokenName, token string, expir time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = tokenName
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}

func setUserCookie(user *entity.User, expir time.Time, c echo.Context) {

	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Email
	cookie.Expires = expir
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func generateRefreshToken(user *entity.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshKeySecret()))
}

func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Get("user") == nil {
			return next(c)
		}

		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "JWT token missing or invalid")
		}

		claims := token.Claims.(*Claims)

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15*time.Minute {
			rc, err := c.Cookie(refreshTokenName)
			if err == nil && rc != nil {
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(t *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshKeySecret()), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}
				if tkn != nil && tkn.Valid {
					if err != nil {
						panic(err)
					}

					if err != nil {
						panic(err)
					}
					if err != nil {
						c.Response().Writer.WriteHeader(http.StatusInternalServerError)
					}
					_ = GenerateTokenAndSetCookie(&entity.User{
						Name:  claims.Name,
						Email: claims.Email,
					}, c)
				}
			}
		}

		return nil
	}
}
