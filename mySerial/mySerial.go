package mySerial

import (
	"github.com/tarm/serial"
	"log"
	"time"
)

type SerialPort struct {
	Port   *serial.Port
	Config *serial.Config
}

func NewSerialPort() (*SerialPort, error) {
	config := &serial.Config{
		Name:        "COM3",
		Baud:        9600,
		ReadTimeout: time.Second * 5,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	return &SerialPort{Port: port, Config: config}, nil
}

func (sp *SerialPort) ClosePort() error {
	return sp.Port.Close()
}

func (sp *SerialPort) ReadData(data []byte) (int, error) {
	return sp.Port.Read(data)
}

func (sp *SerialPort) WriteData(data []byte) (int, error) {
	n, err := sp.Port.Write(data)
	if err != nil {
		return 0, err
	}
	return n, nil
}

var SP = initialization()

func initialization() *SerialPort {
	SP, err := NewSerialPort()
	if err != nil {
		log.Fatal(err)
	}
	return SP
}
