package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	rand.Seed(time.Now().Unix())
	SetupSetting()
	SetupDatabase()
}

func main() {
	gin.SetMode(ServerSetting.RunMode)

	routersInit := InitRouter()
	endPoint := fmt.Sprintf(":%d", ServerSetting.HttpPort)

	go updatePricesTask()
	go sellProductTask()

	server := &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}

	log.Printf("Start http server on Port %s", endPoint)

	err := server.ListenAndServe()
	log.Println(err)
}

func updatePricesTask() {
	frequency := 30 * time.Second
	ticker := time.NewTicker(frequency)
	for range ticker.C {
		updatePrices(frequency)
	}
}

func updatePrices(frequency time.Duration) {
	lastTime := time.Now().Add(-frequency)
	beforeLastTime := time.Now().Add((-frequency) * 2)
	lastProducts, _ := GetProducts(db, lastTime, time.Now())
	beforeLastProducts, _ := GetProducts(db, beforeLastTime, lastTime)

	// CALCULATE GLOBAL DATA
	productCount := len(lastProducts)
	var sells uint = 0
	for _, s := range lastProducts {
		lastSells := len(s.SellData)
		sells = sells + uint(lastSells)
	}
	fmt.Println("All Sells: " + fmt.Sprint(sells))
	for _, s := range lastProducts {
		blProduct := getProductFromList(beforeLastProducts, s.ID)
		if blProduct == nil {
			continue
		}
		productSells := uint(len(s.SellData))
		fmt.Println("Sells of " + s.Name + ": " + fmt.Sprint(productSells))
		lastProductSells := uint(len(blProduct.SellData))
		fmt.Println("Last Sells of " + s.Name + ": " + fmt.Sprint(lastProductSells))
		sellPercentage := float32(0)
		if sells != 0 {
			sellPercentage = float32(productSells) / float32(sells)
		}
		sellPercentageP := sellPercentage * float32(productCount)
		sellChange := float32(0)
		if lastProductSells != 0 {
			sellChange = float32(productSells) / float32(lastProductSells)
		} else {
			sellChange = 1
		}
		sellChangeFactor := calculateInfluence(sellChange, 4)
		otherProductsFactor := float32(0)
		if sellPercentageP > 1 {
			otherProductsFactor = calculateInfluence(sellPercentageP, 4)
		} else {
			otherProductsFactor = 1
		}

		newPrice := uint(float32(s.StartPrice) * sellChangeFactor * otherProductsFactor)
		SetNewPrice(db, s.ID, newPrice)
		fmt.Println("New Price for Product: " + s.Name + ": " + fmt.Sprint(newPrice))
	}
}

func calculateInfluence(factor float32, multiplier float32) float32 {
	return factor/multiplier + 1 - 1/multiplier
}

func getProductFromList(products []Product, id uint) *Product {
	for _, s := range products {
		if s.ID == id {
			return &s
		}
	}
	return nil
}

func getFirstSell(products []Product) *time.Time {
	var firstSell *time.Time
	for _, s := range products {
		if len(s.SellData) == 0 {
			continue
		}
		productFirstSell := s.SellData[0].CreatedAt
		if firstSell == nil || productFirstSell.Before(*firstSell) {
			firstSell = &productFirstSell
		}
	}
	return firstSell
}

func sellProductTask() {

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		now := time.Now().UnixMilli() / 1000 / 120
		fmt.Println(now)
		product1 := int(math.Sin(float64(now)*float64(0.25))*5 + 5)
		product2 := int(math.Sin((float64(now)+4*math.Pi)*float64(0.25))*5 + 5)
		fmt.Println("1: " + fmt.Sprint(product1))
		fmt.Println("2: " + fmt.Sprint(product2))
		probabilities := make([]uint, product1+product2)
		for i := 0; i < product1; i++ {
			probabilities[i] = 1
		}
		for i := product1; i < product1+product2; i++ {
			probabilities[i] = 2
		}

		sellProduct := probabilities[rand.Intn(len(probabilities))]
		fmt.Println("Selling Product: " + fmt.Sprint(sellProduct))
		SellProduct(db, sellProduct)
	}
}
