package sum

import (
	"net/http"

	postgres "my.assignment/utils/postgresDB"

	"github.com/gin-gonic/gin"
)

var pg *postgres.Postgres

func Init(r *gin.RouterGroup, pgLocal *postgres.Postgres) {
	pg = pgLocal
	userRoute := r.Group("/sum")
	userRoute.GET("/:id", SumofTransaction)

}

func SumofTransaction(c *gin.Context) {
	tx, err := pg.SelectSingleTransaction(c.Param("id"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	childSum, err := getSum([]string{tx.ID})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, childSum+*tx.Amount)
}

func getSum(ids []string) (float64, error) {
	allChildTransactionList, err := pg.SelectChildTransactions(ids)
	if err != nil {
		return 0, err
	}
	childIds := []string{}
	sum := float64(0)
	for _, tx := range allChildTransactionList {
		childIds = append(childIds, tx.ID)
		sum += *tx.Amount
	}
	if len(childIds) == 0 {
		return sum, nil
	}
	childSum, err := getSum(childIds)

	if err != nil {
		return 0, err
	}
	return childSum + sum, nil

}
