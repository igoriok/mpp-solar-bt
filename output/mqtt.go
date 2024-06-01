package output

import (
	"encoding/json"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttOutput struct {
	client mqtt.Client
	topic  string
	qos    byte
}

func init() {
	mqtt.ERROR = log.New(os.Stderr, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stderr, "[CRITICAL] ", 0)
	mqtt.WARN = log.New(os.Stderr, "[WARN] ", 0)
	mqtt.DEBUG = log.New(os.Stderr, "[DEBUG] ", 0)
}

func NewMqttOutput(server string, topic string) Output {
	return &mqttOutput{
		client: newClient(server),
		topic:  topic,
	}
}

func (output *mqttOutput) Write(data map[string]interface{}) error {

	body, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = output.connect()

	if err != nil {
		return err
	}

	defer output.disconnect()

	return output.publish(body)
}

func (output *mqttOutput) connect() error {

	log.Print("Connecting...")

	t := output.client.Connect()

	if t.Wait() && t.Error() != nil {
		return t.Error()
	}

	log.Print("Connected!")

	return nil
}

func (output *mqttOutput) disconnect() {
	output.client.Disconnect(100)
}

func (output *mqttOutput) publish(data interface{}) error {

	log.Print("Publishing...")

	t := output.client.Publish(output.topic, output.qos, false, data)

	if t.Wait() && t.Error() != nil {
		return t.Error()
	}

	log.Print("Published!")

	return nil
}

func newClient(server string) mqtt.Client {

	config := mqtt.NewClientOptions().AddBroker(server)
	client := mqtt.NewClient(config)

	return client
}
