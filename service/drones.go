package service

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"sample-project/structs"
)

func (api *APIv1) createDrone(c echo.Context) error {
	drone := structs.Drone{}
	err := c.Bind(&drone)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.httpRespUnsuccessful(err.Error()))
	}

	err = drone.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.httpRespUnsuccessful(err.Error()))
	}

	userUUID := c.Get("user_uuid")
	if userUUID == nil {
		return c.JSON(http.StatusUnauthorized, api.httpRespUnsuccessful("Can't identify user"))
	}

	drone.UserUUID = userUUID.(string)

	err = api.store.CreateDrone(&drone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.httpRespUnsuccessful(err.Error()))
	}

	return c.JSON(http.StatusCreated, api.httpRespSuccessful("Created"))
}

func (api *APIv1) getAllDrones(c echo.Context) error {
	drones, err := api.store.GetAllDrones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.httpRespUnsuccessful(err.Error()))
	}

	return c.JSON(http.StatusOK, api.httpRespSuccessful(drones))
}
