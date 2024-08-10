package transactionTypes

import (
	"net/http"

	postgres "my.assignment/utils/postgresDB"

	"github.com/gin-gonic/gin"
)

var pg *postgres.Postgres

func Init(r *gin.RouterGroup, pgLocal *postgres.Postgres) {
	pg = pgLocal
	userRoute := r.Group("/types")
	userRoute.GET("/:type", transactionFromType)
}

func transactionFromType(c *gin.Context) {
	typeIds := []string{}
	txList, err := pg.SelectTransactionFromType(c.Param("type"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	for _, tx := range txList {
		typeIds = append(typeIds, tx.ID)
	}
	c.IndentedJSON(http.StatusOK, typeIds)
}
