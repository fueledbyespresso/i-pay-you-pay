package api

import (
	"IPYP/api/groupLedger"
	"IPYP/api/transactions"
	"IPYP/auth/authorization"
	"IPYP/database"
	"github.com/gin-gonic/gin"
)

//All the routes created by the package nested in
// api/v1/*
func Routes(r *gin.RouterGroup, db *database.DB) {
	groupLedgerRoutes(r, db)
	transactionRoutes(r, db)
}

func transactionRoutes(r *gin.RouterGroup, db *database.DB) {
	r.PUT("/transaction/",
		authorization.ValidSession(db),
		transactions.CreateTransaction(db))
	r.GET("/transaction",
		authorization.ValidSession(db),
		transactions.GetTransaction(db))
	r.GET("/transactions",
		authorization.ValidSession(db),
		transactions.GetTransactions(db))
	r.POST("transactions/:id",
		authorization.ValidSession(db),
		transactions.UpdateTransaction(db))
	r.DELETE("transactions/:id",
		authorization.ValidSession(db),
		transactions.DeleteTransaction(db))
}

func groupLedgerRoutes(r *gin.RouterGroup, db *database.DB) {
	r.PUT("/group/",
		authorization.ValidSession(db),
		groupLedger.CreateGroup(db))
	r.GET("/group/",
		authorization.ValidSession(db),
		groupLedger.GetGroup(db))
}
