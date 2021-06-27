package service

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sample-project/structs"
)

func (api *APIv1) createDrone(c *gin.Context) {
	drone := structs.Drone{}
	err := c.Bind(&drone)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.httpRespUnsuccessful(err.Error()))
		return
	}

	userUUID, isExists := c.Get("user_uuid")
	if !isExists {
		c.JSON(http.StatusUnauthorized, api.httpRespUnsuccessful("Can't identify user"))
		return
	}

	drone.UserUUID = userUUID.(string)

	err = api.store.CreateDrone(&drone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.httpRespUnsuccessful(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, api.httpRespSuccessful("Created"))
}

func (api *APIv1) getAllDrones(c *gin.Context) {
	drones, err := api.store.GetAllDrones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.httpRespUnsuccessful(err.Error()))
		return
	}

	c.JSON(http.StatusOK, api.httpRespSuccessful(drones))
}
