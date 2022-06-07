package transactions

import (
	"IPYP/database"
	"github.com/docker/distribution/uuid"
	"github.com/gin-gonic/gin"
	"time"
)

type transaction struct {
	ID         int       `json:"id"`
	Total      int       `json:"total"`
	Desc       string    `json:"desc"`
	Recorder   uuid.UUID `json:"recorder"`
	TransTime  time.Time `json:"time_of_transaction"`
	RecordTime time.Time `json:"time_of_record"`
}

func AddContact(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var trans transaction
		err := c.BindJSON(&trans)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid request.")
			return
		}

		row := db.Db.QueryRow(`INSERT INTO transaction (total, description, recorder, time_of_record, time_of_transaction) 
																VALUES ($1, $2, $3, $4, $5) returning id`, trans.Total, trans.Desc, trans.Recorder, trans.RecordTime, trans.TransTime)
		err = row.Scan(&trans.ID)
		if err != nil {
			c.AbortWithStatusJSON(400, "User does not exist.")
			return
		}

		c.JSON(200, trans)
	}
}

func GetContacts(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var trans transaction
		err := c.BindJSON(&trans)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid request.")
			return
		}

		row := db.Db.QueryRow(``)
		err = row.Scan(&trans.ID)
		if err != nil {
			c.AbortWithStatusJSON(400, "User does not exist.")
			return
		}

		c.JSON(200, trans)
	}
}
