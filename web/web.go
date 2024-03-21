package web

import (
	"fmt"
	"html/template"
	"log"
	"main/arduino"
	database "main/db"
	"net/http"
	"strconv"
	"strings"
)

func getInfo(w http.ResponseWriter, r *http.Request) {
	db := database.DB{}
	data, err := db.GetCards()
	if err != nil {
		http.Error(w, "Failed to retrieve card information", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/template.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

func updateCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cardNumber := r.FormValue("bytes")
	newCommand := r.FormValue("newCommand")
	relayTime := r.FormValue("relayTime")

	byteData, err := BytConvert(cardNumber)
	if err != nil {
		log.Fatal(err)
	}

	arduino.ByteReplacing(byteData, newCommand, relayTime)
	getInfo(w, r)

}

func deleteCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cardNumber := r.FormValue("bytes")

	byteData, err := BytConvert(cardNumber)
	if err != nil {
		log.Fatal(err)
	}

	db := database.DB{}
	err = db.DeleteCard(byteData)

	if err != nil {
		fmt.Println(err)
	}
	getInfo(w, r)

}

func BytConvert(cardNumber string) ([]byte, error) {
	byteStrs := strings.Split(strings.Trim(cardNumber, "[]"), " ")

	cardByte := make([]byte, len(byteStrs))

	for i, byteStr := range byteStrs {
		byteValue, err := strconv.ParseUint(byteStr, 10, 8)
		if err != nil {
			return cardByte, fmt.Errorf("Error converting a string %s to byte: %v\n", byteStr, err)
		}
		cardByte[i] = byte(byteValue)
	}
	return cardByte, nil
}

func StartWeb() {
	mux := http.NewServeMux()
	mux.HandleFunc("/getCards", getInfo)
	mux.HandleFunc("/change-command", updateCard)
	mux.HandleFunc("/delete-card", deleteCard)
	err := http.ListenAndServe(":8888", mux)
	if err != nil {
		log.Fatal(err)
	}
}
