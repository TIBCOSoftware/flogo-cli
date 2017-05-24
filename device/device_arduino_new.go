package device

var triggerTpls = make(map[string]string)
var actionTpls = make(map[string]string)
var activityTpls = make(map[string]string)

func init() {
	feather_m0 := &DeviceDetails{Type:"feather_m0_wifi2", Board:"adafruit_feather_m0_usb"}

	feather_m0.MainFile = tplArduinoMainNew

	files := map[string]string{
		"mqtt.ino": tplArduinoMqttNew,
		"wifi.ino": tplArduinoWifiNew,
	}
	libs := map[string]int{
		"PubSubClient": 89,
		"WiFi101": 299,
	}
	feather_m0.MqttFiles = files
	feather_m0.Libs = libs

	Register(feather_m0)

	triggerTpls["github.com/TIBCOSoftware/flogo-contrib/trigger/device-pin"] = tplTriggerPin
	triggerTpls["github.com/TIBCOSoftware/flogo-contrib/trigger/device-mqtt"] = tplTriggerMqtt
	triggerTpls["github.com/TIBCOSoftware/flogo-contrib/trigger/device-button"] = tplTriggerButton

	actionTpls["github.com/TIBCOSoftware/flogo-contrib/action/device-activity"] = tplActionActivity

	activityTpls["github.com/TIBCOSoftware/flogo-contrib/activity/device-pin"] = tplActivityPin
	activityTpls["github.com/TIBCOSoftware/flogo-contrib/activity/device-mqtt"] = tplActivityMqtt
}


var tplArduinoMainNew = `#include <SPI.h>
{{if .MqttEnabled}}#include <WiFi101.h>
#include <PubSubClient.h>

WiFiClient wifiClient;
PubSubClient client(wifiClient);

{{end}}


void setup() {
    Serial.begin(115200);

    while (!Serial) {
        delay(10);
    }

    {{if .MqttEnabled}}
    setup_wifi();
    setup_mqtt();

	//init mqtt triggers
	{{range .MqttTriggers}}t_{{.}}_init();
	{{end}}
    {{end}}

	//init simple triggers
	{{range .Triggers}}t_{{.}}_init();
	{{end}}

	//init actions
	{{range .Actions}}a_{{.}}_init();
	{{end}}
}

void loop() {

    {{if .MqttEnabled}}
    if (!client.connected()) {
        mqtt_reconnect();
    }

    // MQTT client loop processing
    client.loop();
    {{end}}

	//triggers
	{{range .Triggers}}t_{{.}}();
	{{end}}
}

{{if .MqttEnabled}}
void callback(char *topic, byte *payload, unsigned int length) {

    Serial.print("Message arrived [");
    Serial.print(topic);
    Serial.print("] ");
    for (int i=0; i < length; i++) {
        Serial.print((char) payload[i]);
    }
    Serial.println();

	//mqtt triggers
	{{range .MqttTriggers}}t_{{.}}(topic, payload, length); {{end}}
}
{{end}}
`

var tplArduinoWifiNew = `#include <SPI.h>
#include <WiFi101.h>

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

var tplArduinoMqttNew = `#include <SPI.h>
#include <WiFi101.h>
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

void publishMQTT(String value, String payload) {
	payload.toCharArray(out_msg_buff, payload.length() + 1);
	client.publish(mqtt_pubTopic, out_msg_buff);
}
`

var tplTriggerPin = `
uint8_t t_{{.Id}}_pin = {{setting . "pin"}};    // set input pin
bool t_{{.Id}}_lc = false;    // last condition value

void t_{{.Id}}_init() {
	pinMode(t_{{.Id}}_pin, INPUT);
}

void t_{{.Id}}() {

    int value = {{if settingb . "digital"}}digitalRead(t_{{.Id}}_pin){{else}}analogRead(t_{{.Id}}_pin){{end}};

    // create custom condition
    bool condition = value {{setting . "condition"}};

    if (condition && !t_{{.Id}}_lc) {

        a_{{.ActionId}}(value);
    }

    t_{{.Id}}_lc = condition;
}
`


//triggers will say if they are part of loop or mqtt callback
var tplTriggerMqtt = `

void t_{{.Id}}_init() {
}

void t_{{.Id}}(char *topic, byte *payload, unsigned int length) {

    if (topic == "{{setting . "topic"}}") {

        char buf[8];
        int i=0;

        for(i=0; i<length; i++) {
            buf[i] = payload[i];
        }
        buf[i] = '\0';

        int value = atoi(buf);

        a_{{.ActionId}}(value);
    }
}
`

var tplTriggerButton = `
unsigned long t_{{.Id}}_ldt = 0; // lastDebounceTime

int t_{{.Id}}_bs;          // the current reading from the input pin
int t_{{.Id}}_lbs = LOW; // the previous reading from the input pin

uint8_t t_{{.Id}}_pin = {{setting . "pin"}};  //set input pin

void t_{{.Id}}_init() {
	pinMode(t_{{.Id}}_pin, INPUT);
}

void t_{{.Id}}() {

  int reading = digitalRead(t_{{.Id}}_pin);

  if (reading != t_{{.Id}}_lbs) {
    // reset the debouncing timer
    t_{{.Id}}_ldt = millis();
  }

  if ((millis() - t_{{.Id}}_ldt) > 50) {

    if (reading != t_{{.Id}}_bs) {
      t_{{.Id}}_bs = reading;

      if (t_{{.Id}}_bs == HIGH) {
        a_{{.ActionId}}(HIGH);
      }
    }
  }

  t_{{.Id}}_lbs = reading;
}
`

var tplActionActivity = `

void a_{{.Id}}_init() {
 	{{ template "activity_init" .Data }}
}

void a_{{.Id}}(int value) {
 	{{ template "activity_code" .Data }}
}
`

var tplActivityPin = `
{{ define "activity_init"}}
  pinMode({{setting .Activity "pin"}}, OUTPUT);
{{ end }}

{{ define "activity_code"}}
  int val = {{if .UseTriggerVal}}value{{else}}{{setting .Activity "value"}}{{end}};

  {{if settingb .Activity "digital"}}digitalWrite({{setting .Activity "pin"}}, val){{else}}analogWrite({{setting .Activity "pin"}}, val){{end}};
{{ end }}
`

var tplActivityMqtt = `
{{ define "activity_init"}}{{end}}

{{ define "activity_code"}}
	String payload = {{if .UseTriggerVal}}String(value){{else}}"{{setting .Activity "payload"}}"{{end}};
	publishMQTT("{{setting .Activity "topic"}}", payload);
{{end}}
`

