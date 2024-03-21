package main

import (
	"log"
	"main/arduino"
	database "main/db"
	"main/web"
)

func main() {
	db := database.DB{}
	err := db.DBInit()
	if err != nil {
		log.Fatal(err)
	}
	go web.StartWeb()

	arduino.Receiver()

}
