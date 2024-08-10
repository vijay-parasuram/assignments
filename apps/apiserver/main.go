package main

import (
	"github.com/gin-gonic/gin"
	"my.assignment/apiserver/transactionservice"
	postgres "my.assignment/utils/postgresDB"
)

func main() {
	router := gin.Default()
	pg, err := postgres.NewPostgres()
	if err != nil {
		panic(err)
	}
	defer pg.Db.Close()
	transactionservice.Init(router, pg)
	router.Run(":8080")
}
