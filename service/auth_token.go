package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"sample-project/structs"
)

func (api *APIv1) tokenAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userUUID, err := api.verifyToken(c.Request())
		if err != nil {
			return c.JSON(http.StatusUnauthorized, api.httpRespUnsuccessful(err.Error()))
		}

		c.Set("user_uuid", userUUID)
		return next(c)
	}
}

func (api *APIv1) verifyToken(r *http.Request) (string, error) {
	bearerToken := r.Header.Get("Authorization")
	headerAuth := strings.Split(bearerToken, " ")
	if len(headerAuth) < 2 {
		return "", fmt.Errorf("Extract token: invalid Authorization header")
	}

	token, err := jwt.Parse(headerAuth[1], func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	accessDetails := &structs.AccessDetails{}
	var userUUID string
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return "", fmt.Errorf("Invalid token")
		}

		userUUID, ok = claims["user_uuid"].(string)
		if !ok {
			return "", fmt.Errorf("Invalid token")
		}

		accessDetails.AccessUuid = accessUuid
		accessDetails.UserUUID = userUUID
	}

	err = api.redisClient.FetchAuth(accessDetails)
	if err != nil {
		return "", fmt.Errorf("Unauthorized")
	}

	return userUUID, nil
}
