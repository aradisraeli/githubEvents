package routeHandlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"githubEvents/api/controllers"
	"githubEvents/shared"
	_ "githubEvents/shared/models"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

// GetAllEventsHandler godoc
// @Summary List all collected events.
// @Description List all collected events.
// @Tags events
// @Accept */*
// @Produce json
// @Param        size    query     int64  true  "size of page."
// @Param        page    query     int64  true  "number of page."
// @Success 200 {object} models.Page[models.Event]
// @Router /api/v1/events [get]
func GetAllEventsHandler(c echo.Context) error {
	client, ok := c.Get("mongoClient").(*mongo.Client)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	queryParams := c.QueryParams()
	params := make(map[string]interface{})
	for key, val := range queryParams {
		if len(val) > 0 {
			params[key] = val[0]
		}
	}
	size, err := strconv.ParseInt(c.QueryParam("size"), shared.DecimalBase, shared.BitSize)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	} else if size < 1 {
		message := fmt.Sprintln("Size cannot be lower then 1. Got ", size)
		log.Errorf(message)
		return c.JSON(http.StatusBadRequest, message)
	}
	delete(params, "size")
	page, err := strconv.ParseInt(c.QueryParam("page"), shared.DecimalBase, shared.BitSize)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	} else if page < 1 {
		message := fmt.Sprintln("Page cannot be lower then 1. Got ", page)
		log.Errorf(message)
		return c.JSON(http.StatusBadRequest, message)
	}
	delete(params, "page")

	paginatedResult, err := controllers.GetAllEvents(client, page, size, params)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, paginatedResult)
}

// CountEventsHandler godoc
// @Summary Count all collected events.
// @Description Count all collected events.
// @Tags events
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]int
// @Router /api/v1/events/count [get]
func CountEventsHandler(c echo.Context) error {
	client, ok := c.Get("mongoClient").(*mongo.Client)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	queryParams := c.QueryParams()
	params := make(map[string]interface{})
	for key, val := range queryParams {
		if len(val) > 0 {
			params[key] = val[0]
		}
	}

	totalEvents, err := controllers.CountEvents(client, params)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"count": totalEvents,
	})
}

// RecentActorsHandler godoc
// @Summary Gets the actors of recent events.
// @Description Gets the actors of recent events.
// @Tags events
// @Accept */*
// @Produce json
// @Param        amount    query     int64  true  "amount of recent actors."
// @Success 200 {object} []models.Actor
// @Router /api/v1/events/actors [get]
func RecentActorsHandler(c echo.Context) error {
	client, ok := c.Get("mongoClient").(*mongo.Client)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	amount, err := strconv.ParseInt(c.QueryParam("amount"), shared.DecimalBase, shared.BitSize)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	} else if amount < 1 {
		message := fmt.Sprintln("Amount cannot be lower then 1. Got ", amount)
		log.Errorf(message)
		return c.JSON(http.StatusBadRequest, message)
	}

	result, err := controllers.RecentActors(client, amount)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, result)
}

// RecentReposHandler godoc
// @Summary Gets the repos of recent events.
// @Description Gets the repos of recent events.
// @Tags events
// @Accept */*
// @Produce json
// @Param        amount    query     int64  true  "amount of recent repos."
// @Success 200 {object} []models.EventRepo
// @Router /api/v1/events/repos [get]
func RecentReposHandler(c echo.Context) error {
	client, ok := c.Get("mongoClient").(*mongo.Client)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	amount, err := strconv.ParseInt(c.QueryParam("amount"), shared.DecimalBase, shared.BitSize)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, err)
	} else if amount < 1 {
		message := fmt.Sprintln("Amount cannot be lower then 1. Got ", amount)
		log.Errorf(message)
		return c.JSON(http.StatusBadRequest, message)
	}

	result, err := controllers.RecentRepos(client, amount)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, result)
}
