package main

import "github.com/atr0phy/go-microservices/internal/api/app"

func main() {
	router := app.SetupRouter()
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
