package cmd

import (
	"fmt"
	"os"
	"watch-power-bt/input"
	"watch-power-bt/output"

	"github.com/spf13/cobra"
)

var (
	mqttHost, mqttTopic string

	rootCmd = &cobra.Command{
		Args: cobra.ExactArgs(1),
		RunE: run,
	}
)

func init() {
	rootCmd.Flags().StringVar(&mqttHost, "mqtt-host", "localhost", "MQTT server")
	rootCmd.Flags().StringVar(&mqttTopic, "mqtt-topic", "mpp-solar-bt", "MQTT topic")
}

func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {

	bluetooth := input.NewBluetoothInput(args[0])
	mqtt := output.NewMqttOutput(fmt.Sprintf("tcp://%s:1883", mqttHost), mqttTopic)
	writer := output.NewWriterOuput(os.Stdout)

	data, err := bluetooth.Read()

	if err != nil {
		return err
	}

	err = writer.Write(data)

	if err != nil {
		return err
	}

	err = mqtt.Write(data)

	if err != nil {
		return err
	}

	return nil
}
