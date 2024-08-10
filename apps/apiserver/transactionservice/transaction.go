package transactionservice

import (
	"github.com/gin-gonic/gin"
	"my.assignment/apiserver/transactionservice/sum"
	"my.assignment/apiserver/transactionservice/transaction"
	"my.assignment/apiserver/transactionservice/transactionTypes"
	postgres "my.assignment/utils/postgresDB"
)

func Init(r *gin.Engine, pgLocal *postgres.Postgres) {
	userRoute := r.Group("/transactionservice")
	transaction.Init(userRoute, pgLocal)
	sum.Init(userRoute, pgLocal)
	transactionTypes.Init(userRoute, pgLocal)
}
