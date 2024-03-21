package arduino

import (
	"fmt"
	"log"
	database "main/db"
	"main/mySerial"

	"github.com/lib/pq"
)

type Message struct {
	FirstByte byte
	Command   byte
	Length    uint8
	UID       []byte
	CRC       uint8
	EndByte   byte
}

func Receiver() {

	buf := make([]byte, 256)
	var message Message
	for {
		n, err := mySerial.SP.ReadData(buf)
		if err != nil {
			log.Fatalf("Error port reading: %v", err)
		}
		if n > 0 {
			if n == 1 {
				message.FirstByte = buf[0]
				if message.FirstByte == 0xFF {
					fmt.Println("WATER DETECTED")
				}
			} else if n == 9 {
				message.Command = buf[0]
				message.Length = buf[1]
				message.UID = buf[2 : 2+message.Length]
				message.CRC = buf[7]
				message.EndByte = buf[8]

				err := distribution(message)
				if err != nil {
					fmt.Println(err)
				}

			}
		}
	}
}

func distribution(message Message) error {
	db := database.DB{}

	if message.FirstByte != 0xAA {
		return fmt.Errorf("unexpected FirstByte value: %02X", message.FirstByte)
	}

	card, err := db.GetCard(message.UID)
	if err != nil {
		return err
	}

	if card.Command != "UNDEFINED" && len(card.Message) > 0 {
		err := sendMessage(card.Message)
		if err != nil {
			return err
		}

	} else {
		err := db.CreateCard(message.UID)
		if err != nil {
			pqErr, ok := err.(*pq.Error)
			if ok && pqErr.Code == "23505" {
				fmt.Println("The card has already been saved, visit the http://127.0.0.1:8888/getCards to assign an action")
				return nil
			}
			return err
		}
	}

	return nil
}

func sendMessage(data []byte) error {
	_, err := mySerial.SP.WriteData(data)
	if err != nil {
		return err
	}
	return nil

}
