package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetDataParams struct {
	From *time.Time `form:"from" json:"from" binding:"required" time_format:"2006-01-02T15:04:05.000Z07:00"`
	To   *time.Time `form:"to" json:"to" binding:"required" time_format:"2006-01-02T15:04:05.000Z07:00"`
}

type SellParams struct {
	ProductID uint   `json:"id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	//r.Use(gin.Logger())
	r.Use(CORS())

	r.GET("/", func(ctx *gin.Context) {
		var getDataParams GetDataParams

		if err := ctx.BindQuery(&getDataParams); err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, "Missing parameters")
			return
		}

		products, err := GetProducts(db, *getDataParams.From, *getDataParams.To)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, products)
	})
	r.POST("/", func(ctx *gin.Context) {
		sellParams := SellParams{}
		if err := ctx.BindJSON(&sellParams); err != nil {
			ctx.JSON(http.StatusBadRequest, "Missing parameters")
			return
		}
		if sellParams.Password != AppSetting.SellPassword {
			ctx.JSON(http.StatusBadRequest, "Invalid credentials")
			return
		}
		fmt.Println("GOT SELL REQUEST: " + fmt.Sprint(sellParams.ProductID))
		err := SellProduct(db, uint(sellParams.ProductID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, nil)
	})
	return r
}
