package device

/*
    "devices": [
      {"type": "feather_m0_wifi", "board":"adafruit_feather_m0_usb", "template":"device-tpl/arduino.tmpl", "source":"arduino.ino" }
    ],
 */
/*    "libs" : [
{"name": "PubSubClient", "id":89},
{"name": "WiFi101", "id":299}
]*/

func init() {
	feather_m0 := &DeviceDetails{Type:"feather_m0_wifi", Board:"adafruit_feather_m0_usb"}
	files := map[string]string{
		"arduino.ino": tplArduinoMain,
		"mqtt.ino": tplArduinoMqtt,
		"wifi.ino": tplArduinoWifi,
	}
	libs := map[string]int{
		"PubSubClient": 89,
		"WiFi101": 299,
	}
	feather_m0.Files = files
	feather_m0.Libs = libs

	Register(feather_m0)
}


var tplArduinoMain = `#include <SPI.h>
#include <WiFi101.h>
#include <PubSubClient.h>

String endpointId = "0"; //optional endpoint indentifier

uint8_t InPin = A3;  //set input pin
uint8_t OutPin = A4; //set ouput pin

bool digitalIn = true;
bool digitalOut = true;

unsigned long lastDebounceTime = 0;
unsigned long debounceDelay = 50;

//////////////////////

WiFiClient wifiClient;
PubSubClient client(wifiClient);

char in_msg_buff[100];

bool lastCondition = false;

String fireTrigger() {
    int value;

    if (digitalIn) {
        value = digitalRead(InPin);
    } else {
        value = analogRead(InPin);
    }

    // create custom condition
    bool condition = value == 1;

    String ret = "";

    if (condition && !lastCondition) {

        if ((millis() - lastDebounceTime) > debounceDelay) {

            Serial.print("fire endpoint message");

            //value changed so publish trigger event
            ret = String(value);
        }

        lastDebounceTime = millis();
    }

    lastCondition = condition;

    return ret;
}

void setup() {
    Serial.begin(115200);

    while (!Serial) {
        delay(10);
    }

    setup_wifi();

    pinMode(InPin, INPUT);
    pinMode(OutPin, OUTPUT);

    setup_mqtt();
}

void loop() {

    if (!client.connected()) {
        mqtt_reconnect();
    }

    // MQTT client loop processing
    client.loop();

    String value = fireTrigger();

    if (value.length() > 0) {

        publishMQTT(value);
    }
}

void callback(char *topic, byte *payload, unsigned int length) {
    Serial.print("Message arrived [");
    Serial.print(topic);
    Serial.print("] ");
    for (int i=0; i < length; i++) {
        Serial.print((char) payload[i]);
    }
    Serial.println();

    if (digitalOut) {
        if ((char) payload[0] == '1') {
            digitalWrite(OutPin, HIGH);
        } else {
            digitalWrite(OutPin, LOW);
        }
    } else {
        int i=0;

        for(i=0; i<length; i++) {
            in_msg_buff[i] = payload[i];
        }
        in_msg_buff[i] = '\0';

        int value = atoi(in_msg_buff);

        analogWrite(OutPin, value);
    }
}
`

var tplArduinoWifi = `#include <SPI.h>
#include <WiFi101.h>
#include <PubSubClient.h>

char ssid[] = "{{setting . "wifi:ssid"}}";
const char *password = "{{setting . "wifi:password"}}";

void setup_wifi() {

    //Configure pins for Adafruit ATWINC1500 Feather
    WiFi.setPins(8,7,4,2);

    // check for the presence of the shield:
    if (WiFi.status() == WL_NO_SHIELD) {
        Serial.println("WiFi shield not present");
        // don't continue:
        while (true);
    }

    delay(10);

    Serial.println();
    Serial.print("Connecting to ");
    Serial.println(ssid);

    WiFi.begin(ssid, password);

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    //WiFi.lowPowerMode();

    randomSeed(micros());

    Serial.println("");
    Serial.println("WiFi connected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());
}
`

var tplArduinoMqtt = `#include <SPI.h>
#include <PubSubClient.h>

const char *mqtt_server = "{{setting . "mqtt:server"}}";
const char *mqtt_user = "{{setting . "mqtt:user"}}";
const char *mqtt_pass = "{{setting . "mqtt:pass"}}";
const char *mqtt_pubTopic = "flogo/{{setting . "device:name"}}/out";
const char *mqtt_subTopic = "flogo/{{setting . "device:name"}}/in";

const char *mqtt_readyMsg = "{\"status\": \"READY\"}";

char out_msg_buff[100];

//////////////////////

void setup_mqtt() {
    client.setServer(mqtt_server, 1883);
    client.setCallback(callback);
}

void mqtt_reconnect() {
    // Loop until we're reconnected
    while (!client.connected()) {
        Serial.print("Attempting MQTT connection...");
        // Create a random client ID
        String clientId = "device-{{setting . "device:name"}}-";
        clientId += String(random(0xffff), HEX);
        // Attempt to connect
        if (client.connect(clientId.c_str(), mqtt_user, mqtt_pass)) {
            Serial.println("connected");
            client.publish(mqtt_pubTopic, mqtt_readyMsg);
            client.subscribe(mqtt_subTopic);
        } else {
            Serial.print("failed, rc=");
            Serial.print(client.state());
            Serial.println(" try again in 5 seconds");
            // Wait 5 seconds before retrying
            delay(5000);
        }
    }
}

void publishMQTT(String value) {
	String payload = "{\"ep\":\"" + endpointId +  "\", \"value\": \"" + value + "\"}";
	payload.toCharArray(out_msg_buff, payload.length() + 1);
	client.publish(mqtt_pubTopic, out_msg_buff);
}
`