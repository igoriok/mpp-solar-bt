package input

import (
	"encoding/hex"
	"fmt"
	"log"
	"slices"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

const (
	CHAR_0X2A01 = "00002a01-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A02 = "00002a02-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A03 = "00002a03-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A04 = "00002a04-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A05 = "00002a05-0000-1000-8000-00805f9b34fb" // (no perm)
	CHAR_0X2A06 = "00002a06-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A07 = "00002a07-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A08 = "00002a08-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A09 = "00002a09-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A0B = "00002a0b-0000-1000-8000-00805f9b34fb"

	CHAR_0X2A0C = "00002a0c-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A0D = "00002a0d-0000-1000-8000-00805f9b34fb"
	CHAR_0X2A0E = "00002a0e-0000-1000-8000-00805f9b34fb" // (optional)

	ACTION_POINTER_INDEX_MASK = 65280
)

type bluetoothInput struct {
	addr string
}

func NewBluetoothInput(addr string) Input {
	return &bluetoothInput{
		addr: addr,
	}
}

func (bt *bluetoothInput) Read() (map[string]interface{}, error) {

	chars, err := getCharacteristics(bt.addr, []string{CHAR_0X2A01, CHAR_0X2A02, CHAR_0X2A03, CHAR_0X2A04, CHAR_0X2A05, CHAR_0X2A06, CHAR_0X2A07, CHAR_0X2A08, CHAR_0X2A09, CHAR_0X2A0B})

	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})

	for key, value := range chars {

		log.Printf("%s: %s", key, hex.Dump(value))

		switch key {

		case CHAR_0X2A03:

			data["ac_input_voltage"] = float32(h2l_short(byte2short(value, 0))) / 10
			data["ac_input_frequency"] = float32(h2l_short(byte2short(value, 2))) / 10
			data["ac_output_voltage"] = float32(h2l_short(byte2short(value, 4))) / 10
			data["ac_output_frequency"] = float32(h2l_short(byte2short(value, 6))) / 10
			data["ac_output_apparent_power"] = h2l_short(byte2short(value, 8))
			data["ac_output_active_power"] = h2l_short(byte2short(value, 10))
			data["ac_output_load"] = h2l_short(byte2short(value, 12))
			data["bus_voltage"] = h2l_short(byte2short(value, 14))
			data["battery_voltage"] = float32(h2l_short(byte2short(value, 16))) / 100
			data["battery_charging_current"] = h2l_short(byte2short(value, 18))

		case CHAR_0X2A04:

			data["battery_capacity"] = h2l_short(byte2short(value, 0))
			data["inverter_heat_sink_temperature"] = h2l_short(byte2short(value, 2))
			data["battery_discharge_current"] = h2l_short(byte2short(value, 4))

			if value[13] == 1 {
				data["ac_input_power"] = h2l_short(byte2short(value, 14))
			}

			status := getBitArray(value[6])

			data["status"] = fmt.Sprintf("%v%v%v", status[5], status[6], status[7])
			data["mode"] = string(value[12:13])

		case CHAR_0X2A05:

			data["nominal_ac_input_voltage"] = float32(h2l_short(byte2short(value, 0))) / 10
			data["nominal_ac_input_current"] = float32(h2l_short(byte2short(value, 8))) / 10
			data["rated_battery_voltage"] = float32(h2l_short(byte2short(value, 14))) / 10
			data["nominal_ac_output_voltage"] = float32(h2l_short(byte2short(value, 4))) / 10
			data["nominal_ac_output_frequency"] = float32(h2l_short(byte2short(value, 6))) / 10
			data["nominal_ac_output_current"] = float32(h2l_short(byte2short(value, 8))) / 10
			data["nominal_ac_output_apparent_power"] = h2l_short(byte2short(value, 10))
			data["nominal_ac_output_active_power"] = h2l_short(byte2short(value, 12))
		}
	}

	return data, nil
}

func getCharacteristics(addr string, uuids []string) (map[string][]byte, error) {

	dev, err := connectDevice(addr)

	if err != nil {
		return nil, err
	}

	defer dev.Disconnect()

	log.Print("Reading...")

	chars, err := dev.GetCharacteristics()

	if err != nil {
		return nil, err
	}

	values := make(map[string][]byte)

	for _, char := range chars {

		uuid := char.Properties.UUID

		if slices.Contains(uuids, uuid) {

			value, err := char.ReadValue(make(map[string]interface{}))

			if err != nil {
				continue
			}

			values[uuid] = value
		}
	}

	return values, nil
}

func connectDevice(addr string) (*device.Device1, error) {

	dev, err := findDervice(addr)

	if err != nil {
		return nil, err
	}

	if dev.Properties.Connected {
		return dev, nil
	}

	if !dev.Properties.Paired {

		err = pairDevice(dev)

		if err != nil {
			return dev, err
		}
	}

	log.Print("Connecting...")

	err = dev.Connect()

	if err != nil {
		return dev, nil
	}

	log.Print("Connected!")

	return dev, nil
}

func findDervice(addr string) (*device.Device1, error) {

	a, err := adapter.GetDefaultAdapter()

	if err != nil {
		return nil, err
	}

	//a.SetPairable(false)
	//a.SetDiscoverable(false)
	//a.SetPowered(true)

	dev, _ := a.GetDeviceByAddress(addr)

	if dev != nil {
		return dev, nil
	}

	log.Print("Discovering...")

	err = a.StartDiscovery()

	if err != nil {
		return nil, err
	}

	defer a.StopDiscovery()

	ch, cancel, _ := a.OnDeviceDiscovered()
	defer cancel()

	for e := range ch {

		if e.Type == adapter.DeviceRemoved {
			continue
		}

		dev, _ := device.NewDevice1(e.Path)

		if dev.Properties.Address == addr {

			log.Print("Found!")

			return dev, nil
		}
	}

	return nil, nil
}

func pairDevice(dev *device.Device1) error {

	conn, err := dbus.SystemBus()

	if err != nil {
		return err
	}

	ag := agent.NewSimpleAgent()

	err = agent.ExposeAgent(conn, ag, agent.CapNoInputNoOutput, true)

	if err != nil {
		return err
	}

	defer agent.RemoveAgent(ag)

	log.Print("Pairing...")

	err = dev.Pair()

	if err != nil {
		return err
	}

	err = dev.SetTrusted(true)

	if err != nil {
		return err
	}

	log.Print("Paired!")

	return nil
}

func byte2short(bArr []byte, index int) uint16 {
	var value uint32 = 0
	for i := 0; i < index+2; i++ {
		value += (uint32(bArr[i]) & 255) << (((2 - (i - index)) - 1) * 8)
	}
	return uint16(value)
}

func getBitArray(b byte) []byte {
	var bArr = make([]byte, 8)
	for i := 7; i >= 0; i-- {
		bArr[i] = b & 1
		b = b >> 1
	}
	return bArr
}

func h2l_short(s uint16) uint16 {
	return (((s >> 8) & 255) + (((s & 255) << 8) & ACTION_POINTER_INDEX_MASK))
}
