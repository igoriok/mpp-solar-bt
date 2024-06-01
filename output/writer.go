package output

import (
	"fmt"
	"io"
	"text/tabwriter"
)

var workMode map[string]string = map[string]string{
	"B": "Battery",
	"C": "Charge",
	"D": "Shutdown",
	"E": "Eco",
	"F": "Fault",
	"H": "Power saving",
	"L": "Line",
	"P": "Power on",
	"S": "Stand by",
	"Y": "Bypass",
}

type fieldInfo struct {
	key    string
	title  string
	format string
	unit   string
}

var fields []fieldInfo = []fieldInfo{
	{key: "ac_input_voltage", title: "AC voltage", format: "%.1f", unit: "V"},
	{key: "ac_input_frequency", title: "AC frequency", format: "%.1f", unit: "Hz"},
	{key: "ac_input_power", title: "AC input power", format: "%v", unit: "W"},
	{key: "ac_output_voltage", title: "AC output voltage", format: "%.1f", unit: "V"},
	{key: "ac_output_frequency", title: "AC output frequency", format: "%.1f", unit: "Hz"},
	{key: "ac_output_apparent_power", title: "AC output apparent power", format: "%v", unit: "VA"},
	{key: "ac_output_active_power", title: "AC output active power", format: "%v", unit: "W"},
	{key: "ac_output_load", title: "AC output load", format: "%v", unit: "%"},
	{key: "nominal_ac_input_voltage", title: "Nominal AC input voltage", format: "%.1f", unit: "V"},
	{key: "nominal_ac_input_current", title: "Nominal AC input current", format: "%.1f", unit: "A"},
	{key: "nominal_ac_output_voltage", title: "Nominal AC output voltage", format: "%.1f", unit: "V"},
	{key: "nominal_ac_output_frequency", title: "Nominal AC output frequency", format: "%.1f", unit: "Hz"},
	{key: "nominal_ac_output_current", title: "Nominal AC output current", format: "%.1f", unit: "A"},
	{key: "nominal_ac_output_apparent_power", title: "Nominal AC output apparent power", format: "%v", unit: "VA"},
	{key: "nominal_ac_output_active_power", title: "Nominal AC output active power", format: "%v", unit: "W"},
	{key: "rated_battery_voltage", title: "Rated battery voltage", format: "%.1f", unit: "V"},
	{key: "battery_voltage", title: "Battery voltage", format: "%.1f", unit: "V"},
	{key: "battery_capacity", title: "Battery capacity", format: "%v", unit: "%"},
	{key: "battery_discharge_current", title: "Battery discharge current", format: "%v", unit: "A"},
	{key: "battery_charging_current", title: "Battery charging current", format: "%v", unit: "A"},
	{key: "bus_voltage", title: "Bus voltage", format: "%v", unit: "V"},
	{key: "inverter_heat_sink_temperature", title: "Heat sink temperature", format: "%v", unit: "C"},
}

type writerOutput struct {
	writer io.Writer
}

func NewWriterOuput(writer io.Writer) Output {
	return &writerOutput{
		writer: writer,
	}
}

func (screen *writerOutput) Write(data map[string]interface{}) error {

	w := tabwriter.NewWriter(screen.writer, 1, 1, 1, ' ', 0)

	for _, fieldInfo := range fields {

		value, ok := data[fieldInfo.key]

		if !ok {
			continue
		}

		fmt.Fprintf(w, "%s\t %s %s\n", fieldInfo.title, fmt.Sprintf(fieldInfo.format, value), fieldInfo.unit)
	}

	w.Flush()

	return nil
}
