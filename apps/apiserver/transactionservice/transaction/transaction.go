package transaction

import (
	"net/http"

	postgres "my.assignment/utils/postgresDB"

	"github.com/gin-gonic/gin"
)

var pg *postgres.Postgres

func Init(r *gin.RouterGroup, pgLocal *postgres.Postgres) {
	pg = pgLocal
	userRoute := r.Group("/transaction")
	userRoute.POST("", CreateTransaction)
	userRoute.PUT("/:id", UpdateTransaction)
	userRoute.GET("", GetAllTransaction)
	userRoute.GET("/:id", GetSingleTransaction)

}

func UpdateTransaction(c *gin.Context) {
	t := new(postgres.Transaction)
	c.BindJSON(t)
	t.ID = c.Param("id")
	if t.ID == "" {
		c.String(http.StatusBadRequest, "id is empty")
	}
	err := pg.UpdateSingleTransaction(t)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, t)
}

func CreateTransaction(c *gin.Context) {
	t := new(postgres.Transaction)
	c.BindJSON(t)
	err := pg.InsertSingleTransaction(t)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, t)
}

func GetSingleTransaction(c *gin.Context) {
	txList, err := pg.SelectSingleTransaction(c.Param("id"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, txList)
}

func GetAllTransaction(c *gin.Context) {
	txList, err := pg.SelectAllTransaction()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, txList)
}
