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
      message[2] = 5; //LEN
      for (int i = 0; i < 5; i++) {
        message[3 + i] = rfid.serNum[i]; // CARD UID
      }
      message[8] = 56; //CRC
      message[9] = 0xAE; // END BYTE
      Serial.write(message, 10);
    }
    delay(1000);
  }
  byte receivedData[bufferSize];
  if (Serial.available() > 0) {
    while (Serial.available() > 0 && bytesRead < bufferSize) {
      receivedData[bytesRead] = Serial.read(); 
      bytesRead++; 
    }

    if (receivedData[bytesRead - 1] == 0xFF) {
      if (receivedData[bytesRead - 9] == 0x1F) {
        digitalWrite(relayPin, HIGH);
      delay(receivedData[bytesRead - 7] * 1000);
      digitalWrite(relayPin, LOW);
      } else if (receivedData[bytesRead - 9] == 0x2F){
        startSensor = true;
      }
      }
    }
    int waterVal = analogRead(A0);
    if (startSensor &&  waterVal > 0) {
      Serial.write(0xFF);
  }
}
