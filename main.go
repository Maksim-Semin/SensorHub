package main

import (
	"main/pkg/arduino"
	"main/pkg/storage"
	"main/pkg/web"
)

func main() {

	db := storage.DB{}
	go db.DBInit()

	go web.StartWeb()

	arduino.Receiver()

}
