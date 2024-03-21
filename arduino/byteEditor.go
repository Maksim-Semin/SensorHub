package arduino

import (
	"fmt"
	database "main/db"
	"strconv"
)

const polynomial uint8 = 0xD5

func ByteReplacing(data []byte, command string, duration string) {
	db := database.DB{}

	switch command {
	//			START BYTE -> COMMAND -> LEN-> DATA WITH TIME -> CRC -> END BYTE
	case "Turn on the water sensor":
		crcData := []byte{data[0], data[1], data[2], data[3], data[4]}
		crc := getCRC(crcData)

		message := []byte{0xFF, 0x2F, 5, data[0], data[1], data[2], data[3], data[4], crc, 0xFF}

		db.UpdateCard(data, command, message)
	case "Turn on the relay":

		intTime, _ := strconv.Atoi(duration)
		byteTime := byte(intTime)

		crcData := []byte{byteTime, data[1], data[2], data[3], data[4]}
		crc := getCRC(crcData)

		message := []byte{0xFF, 0x1F, 5, byteTime, data[1], data[2], data[3], data[4], crc, 0xFF}
		db.UpdateCard(data, fmt.Sprintf("%v on %v sec", command, duration), message)
	}
}

func getCRC(data []byte) uint8 {
	var crc uint8 = 0
	for _, b := range data {
		crc ^= b
		for i := 0; i < 8; i++ {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}
