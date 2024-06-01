package main

import (
	"log"
	"os"
	"watch-power-bt/input"
	"watch-power-bt/output"
)

func main() {

	log.SetOutput(os.Stderr)

	bluetooth := input.NewBluetoothInput("3C:E4:B0:8C:6C:B1")
	mqtt := output.NewMqttOutput("tcp://192.168.0.30:1883", "device/inverter/status")
	writer := output.NewWriterOuput(os.Stdout)

	data, err := bluetooth.Read()

	if err != nil {
		log.Fatal(err)
	}

	err = writer.Write(data)

	if err != nil {
		log.Fatal(err)
	}

	err = mqtt.Write(data)

	if err != nil {
		log.Fatal(err)
	}
}
