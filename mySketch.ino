#include <SPI.h>
#include <RFID.h>

#define SS_PIN 10
#define RST_PIN 9

RFID rfid(SS_PIN, RST_PIN);

const int relayPin = 2;

byte message[10];
const int bufferSize = 64; 
byte receivedData[bufferSize]; 
int bytesRead = 0;

bool startSensor = false;

void setup() {
  Serial.begin(9600);
  pinMode(relayPin, OUTPUT);
  pinMode(A0, INPUT);
  SPI.begin();
  rfid.init();
}



void loop() {
  if (rfid.isCard()) {
    if (rfid.readCardSerial()) {
      message[0] = 0xAA; // FIRST BYTE
      message[1] = 0xEE; // DEFAULT VALUE FOR EMPTY COMMAND

      int msgLen = 0;

      for (int i = 0; i < 5; i++) {
        message[3 + i] = rfid.serNum[i]; // CARD UID
        msgLen++;
      }

      message[2] = msgLen;
      message[8] = calculateCRC(&message[3], msgLen); //CRC
      message[9] = 0xAE; // END BYTE
      Serial.write(message, 10);
    }
    delay(1000);
  }
  byte receivedData[bufferSize];
  if (Serial.available() > 0) {
  // read data from buffer
    while (Serial.available() > 0 && bytesRead < bufferSize) {
      receivedData[bytesRead] = Serial.read();
      bytesRead++;
    }

    int msgLen = receivedData[bytesRead - 8];
    int crc = receivedData[bytesRead - 2];

    if (crc != calculateCRC(&receivedData[bytesRead - 7], msgLen)) {
      return;
    }
    //if last byte is correct
    if (receivedData[bytesRead - 1] == 0xFF) {
    // if command on reley
      if (receivedData[bytesRead - 9] == 0x1F) {
        digitalWrite(relayPin, HIGH);
      delay(receivedData[bytesRead - 7] * 1000);
      digitalWrite(relayPin, LOW);

      // if command activate water sensor
      } else if (receivedData[bytesRead - 9] == 0x2F){
        startSensor = true;
      }
      }
    }
    // read value from water sensor
    int waterVal = analogRead(A0);
    if (startSensor &&  waterVal > 0) {
      Serial.write(0xFF);
  }
}

// calculate CRC from message in card
uint8_t calculateCRC(uint8_t* data, size_t len) {
  uint8_t crc = 0;
  uint8_t polynomial = 0xD5;
  for (size_t i = 0; i < len; i++) {
    crc ^= data[i];
    for (uint8_t j = 0; j < 8; j++) {
      if (crc & 0x80) {
        crc = (crc << 1) ^ polynomial;
      } else {
        crc <<= 1;
      }
    }
  }
  return crc;
}