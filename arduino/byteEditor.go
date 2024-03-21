package arduino

import (
	"fmt"
	database "main/db"
	"strconv"
)

func ByteReplacing(data []byte, command string, duration string) {
	db := database.DB{}

	switch command {
	//			START BYTE -> COMMAND -> LEN-> DATA WITH TIME -> CRC -> END BYTE
	case "Turn on the water sensor":

		message := []byte{0xFF, 0x2F, 5, data[0], data[1], data[2], data[3], data[4], 56, 0xFF}

		db.UpdateCard(data, command, message)
	case "Turn on the relay":

		intTime, _ := strconv.Atoi(duration)
		byteTime := byte(intTime)

		message := []byte{0xFF, 0x1F, 5, byteTime, data[1], data[2], data[3], data[4], 56, 0xFF}
		db.UpdateCard(data, fmt.Sprintf("%v on %v sec", command, duration), message)
	}
}
