package api

import (
	"IPYP/api/groupLedger"
	"IPYP/auth/authorization"
	"IPYP/database"
	"github.com/gin-gonic/gin"
)

//All the routes created by the package nested in
// api/v1/*
func Routes(r *gin.RouterGroup, db *database.DB) {
	groupLedgerRoutes(r, db)
	resourceRoute(r, db)
}

func resourceRoute(r *gin.RouterGroup, db *database.DB) {
	r.GET("/resource",
		authorization.LoadPolicy(db, "role"))

}

func groupLedgerRoutes(r *gin.RouterGroup, db *database.DB) {
	r.PUT("/school", authorization.ValidSession(db),
		authorization.LoadPolicy(db, "school"),
		groupLedger.CreateGroup(db))
}
