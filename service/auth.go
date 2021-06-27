package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"

	"sample-project/structs"
)

func (api *APIv1) getAccessDetailsFromReq(r *http.Request) (*structs.AccessDetails, error) {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) < 2 {
		return nil, fmt.Errorf("extract token: invalid Authorization format")
	}

	token, err := jwt.Parse(strArr[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parse: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	userUUID := fmt.Sprintf("%.f", claims["user_uuid"])
	return &structs.AccessDetails{AccessUuid: accessUuid, UserUUID: userUUID}, nil
}

func (api *APIv1) login(c *gin.Context) {
	user := structs.User{}
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.httpRespUnsuccessful(err.Error()))
		return
	}

	userUUID, err := api.store.GetUserUUID(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, api.httpRespUnsuccessful("incorrect login or password"))
		return
	}

	token, err := api.generateToken(userUUID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, api.httpRespUnsuccessful(err.Error()))
		return
	}

	err = api.redisClient.RegisterAuth(userUUID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.httpRespUnsuccessful(err.Error()))
		return
	}

	c.JSON(http.StatusOK, api.httpRespSuccessful(token.AccessToken))
}

func (api *APIv1) logout(c *gin.Context) {
	au, err := api.getAccessDetailsFromReq(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	deleted, err := api.redisClient.DeleteAuth(au.AccessUuid)
	if err != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	c.JSON(http.StatusOK, &Response{Success: true, Message: "Successfully logged out"})
}

func (api *APIv1) generateToken(userUUID string) (*structs.AuthToken, error) {
	token := &structs.AuthToken{}

	token.AccessExpiresAt = time.Now().Add(30 * time.Minute).Unix()
	token.AccessUuid = uuid.NewV4().String()
	token.RefreshExpiresAt = time.Now().Add(24 * 7 * time.Hour).Unix()
	token.RefreshUuid = uuid.NewV4().String()

	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["user_uuid"] = userUUID
	accessTokenClaims["access_uuid"] = token.AccessUuid
	accessTokenClaims["exp"] = token.AccessExpiresAt

	var err error
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	token.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["refresh_uuid"] = token.RefreshUuid
	refreshTokenClaims["user_uuid"] = userUUID
	refreshTokenClaims["exp"] = token.RefreshExpiresAt

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	token.AccessToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	return token, err
}
