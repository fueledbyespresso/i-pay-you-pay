package authorization

import (
	"IPYP/database"
	"github.com/gin-gonic/gin"
)

type GroupLedger struct {
	GroupID int `json:"id"`
}

type Role struct {
	ID        int                    `json:"id"`
	Name      string                 `json:"name"`
	Desc      string                 `json:"desc"`
	Resources map[string]GroupLedger `json:"resources"`
}

func ValidSession(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := db.SessionStore.Get(c.Request, "session")
		if err != nil {
			c.AbortWithStatusJSON(500, "The server was unable to retrieve this session")
			return
		}

		if session.ID == "" {
			c.AbortWithStatusJSON(401, "This user has no current session. Use of this endpoint is thus unauthorized")
			return
		}
	}
}
