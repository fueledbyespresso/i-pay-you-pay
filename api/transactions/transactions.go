package transactions

import (
	"IPYP/database"
	"github.com/docker/distribution/uuid"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"strconv"
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

func CreateTransaction(db *database.DB) gin.HandlerFunc {
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
			c.AbortWithStatusJSON(400, "Invalid transaction data.")
			return
		}

		c.JSON(200, trans)
	}
}

func GetTransaction(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event transaction
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid id. Must be an integer.")
			return
		}

		row := db.Db.QueryRow(`SELECT total, description FROM transaction WHERE id=$1`, id)
		err = row.Scan(&event.Total, &event.Desc)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid transaction ID.")
			return
		}

		c.JSON(200, event)
	}
}

func GetTransactions(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var events []transaction
		session, err := db.SessionStore.Get(c.Request, "session")
		if err != nil {
			c.AbortWithStatusJSON(500, "The server was unable to retrieve this session")
			return
		}
		googleID := session.Values["GoogleId"]

		rows, err := db.Db.Query(`SELECT total, description FROM transaction
   										 JOIN user_transaction_bridge utb on transaction.id = utb.transaction_id
                         				 JOIN account a on transaction.payer = a.user_id
        									WHERE (payer = a.user_id OR utb.user_id = a.user_id)
           									 AND google_id = $1`, googleID)
		if err != nil {
			c.AbortWithStatusJSON(400, "User not properly registered.")
			return
		}

		for rows.Next() {
			var event transaction
			err = rows.Scan(&event.Total, &event.Desc)
			if err != nil {
				c.AbortWithStatusJSON(500, "Invalid transaction ID.")
				return
			}
			events = append(events, event)
		}
		c.JSON(200, events)
	}
}

func UpdateTransaction(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid id. Must be an integer")
			return
		}

		var event transaction
		err = c.BindJSON(&event)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid request.")
			return
		}

		row := db.Db.QueryRow(`UPDATE transaction SET total = $1, description=$2 WHERE id=$3`, event.Total, event.Desc, id)
		if row.Err() != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}
		c.JSON(200, event.ID)
	}
}

func DeleteTransaction(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.AbortWithStatusJSON(400, "Invalid id. Must be an integer.")
			return
		}

		row := db.Db.QueryRow(`DELETE FROM transaction WHERE id=$1`, id)
		if row.Err() != nil {
			database.CheckDBErr(err.(*pq.Error), c)
			return
		}

		c.JSON(200, `Delete successful`)
	}
}
