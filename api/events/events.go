package events

import (
	"IPYP/database"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"strconv"
	"time"
)

type event struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Dates []date `json:"dates"`
}

type date struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func CreateEvent(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event event
		err := c.BindJSON(&event)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid request.")
			return
		}

		row := db.Db.QueryRow(`INSERT INTO event (title, description) VALUES ($1, $2) returning id`, event.Title, event.Desc)
		err = row.Scan(&event.ID)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid event title and description.")
			return
		}

		c.JSON(200, event)
	}
}

func GetEvent(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event event
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid id. Must be an integer.")
			return
		}

		row := db.Db.QueryRow(`SELECT title, description FROM event WHERE id=$1`, id)
		err = row.Scan(&event.Title, &event.Desc)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid event ID.")
			return
		}

		c.JSON(200, event)
	}
}

func GetEvents(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var events []event
		rows, err := db.Db.Query(`SELECT title, description FROM event`)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid event ID.")
			return
		}

		for rows.Next() {
			var event event
			err = rows.Scan(&event.Title, &event.Desc)
			if err != nil {
				c.AbortWithStatusJSON(500, "Invalid event ID.")
				return
			}
			events = append(events, event)
		}
		c.JSON(200, events)
	}
}

func UpdateEvent(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid id. Must be an integer")
			return
		}

		var event event
		err = c.BindJSON(&event)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid request.")
			return
		}

		row := db.Db.QueryRow(`UPDATE event SET title = $1, description=$2 WHERE id=$3`, event.Title, event.Desc, id)
		if row.Err() != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}
		c.JSON(200, event.ID)
	}
}

func DeleteEvent(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid id. Must be an integer.")
			return
		}

		row := db.Db.QueryRow(`DELETE FROM event WHERE id=$1`, id)
		if row.Err() != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}

		c.JSON(200, ``)
	}
}
