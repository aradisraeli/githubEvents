package middlewares

import (
	"github.com/labstack/echo/v4"
	"githubEvents/shared"
	"githubEvents/shared/dal"
	"log"
	"net/http"
)

func MongoDBMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Establish MongoDB connection
		client, err := dal.ConnectToMongo(shared.MongoUri)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}

		// Attach the MongoDB client to the context
		c.Set(shared.MongoClientKey, client)

		// Call the next handler
		if err := next(c); err != nil {
			log.Println(err)
			return err
		}

		// Disconnect from MongoDB after the request has been handled
		defer func() {
			if err := dal.DisconnectFromMongo(client); err != nil {
				log.Println(err)
			}
		}()

		return nil
	}
}
