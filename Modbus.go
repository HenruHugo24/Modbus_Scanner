package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/goburrow/modbus"
)

// Config struct to hold the JSON configuration
type Config struct {
	Protocol           string        `json:"protocol"`
	IPmask             string        `json:"ip_mask"`
	IPRange            string        `json:"ip_device"`
	SlaveIDRange       IDRange       `json:"slave_id_range"`
	KnownRegisterRange RegisterRange `json:"known_register_range"`
	LengthOfEachRead   int           `json:"length_of_each_read"`
	Port               int           `json:"port"`
	FunctionCode       int           `json:"function_code"`
	BaudRates          []int         `json:"baud_rates"`
	Endianness         string        `json:"endianness"`
}

// IDRange struct to hold start and end IDs
type IDRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type RegisterRange struct {
	Start int `json:"end"`
	End   int `json:"end"`
}

func main() {

	// client := modbus.TCPClient("10.6.70.5")
	// //client := modbus.NewTCPClient("10.6.70.5")
	// // Read input register 1027
	// results, _ := client.ReadInputRegisters(1027, 1)
	// fmt.Printf("Read  register: %016b\n", results)
	fuel_level := modbusmaker("10.6.70.5", 0, 1027, 1)
	fmt.Printf("Read Fuel level of deepsea: %d%%\n", fuel_level[1])
	//Connect to new network (Bluelog)
	power_bluelog := modbusmaker("10.6.70.15", 1, 254, 2)
	fmt.Printf("Read Power of bluelog: %v\n", power_bluelog)

	handler1 := modbus.NewTCPClientHandler("10.6.70.15:502")
	handler1.SlaveId = 1
	handler1.Timeout = 5 * time.Second
	defer handler1.Close()
	Bluelog := modbus.NewClient(handler1)

	err1 := handler1.Connect()
	if err1 != nil {
		fmt.Println("Error connecting to Modbus server:", err1)
		return
	}
	defer handler1.Close()

	//Read values
	// Bluelog `10.6.70.15` Power = 254, Freq = 98, Voltage = 100
	Bluelog_power, _ := Bluelog.ReadHoldingRegisters(254, 2)
	Bluelog_freq, _ := Bluelog.ReadHoldingRegisters(98, 2)
	Bluelog_voltage, _ := Bluelog.ReadHoldingRegisters(100, 2)
	fmt.Println(Bluelog_power)
	//Convert values
	b := []byte{0x3f, 0x9d, 0x70, 0xa4} // Represents 1.23456789 in IEEE 754 representation
	fmt.Println(b)
	f := bytesToFloat32(b)
	fmt.Println(f)
	power := bytesToFloat32(Bluelog_power)
	freq := bytesToFloat32(Bluelog_freq)
	Voltage := bytesToFloat32(Bluelog_voltage)

	//Display
	fmt.Println(power)
	fmt.Printf("Read  Power of bluelog: %f\n", power)
	fmt.Printf("Read  Frequency of bluelog: %f\n", freq)
	fmt.Printf("Read  Voltage of bluelog: %f\n", Voltage)

	//New network (SMA)
	handler2 := modbus.NewTCPClientHandler("10.6.70.28:502")
	handler2.Timeout = 5 * time.Second // Set your timeout value
	defer handler1.Close()

	SMA := modbus.NewClient(handler2)

	err2 := handler2.Connect()
	if err2 != nil {
		fmt.Println("Error connecting to Modbus server:", err2)
		return
	}
	defer handler2.Close()

	//Read values
	// SMA Inverter `10.6.70.28` Power = 199, Freq = 201
	SMA_power, err3 := SMA.ReadInputRegisters(233, 1)
	if err3 != nil {
		fmt.Println("Error reading input registers:", err3)
		return
	}
	fmt.Printf("Read SMA power %v\n", SMA_power)

	SMA_freq, _ := SMA.ReadHoldingRegisters(233, 1)
	fmt.Printf("Read SMA power %v", SMA_freq)
}

func bytesToFloat32(bytes []byte) float32 {
	// Assuming Big-endian byte order, adjust accordingly if needed
	bits := binary.BigEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func mask_ip(IP_mask string, IP_device string) {

}

func modbusmaker(IP_mask string, slaveID int, register_value uint16, number_bytes int) []byte {
	address := IP_mask + ":502"
	handler := modbus.NewTCPClientHandler(address)

	handler.Timeout = 5 * time.Second // Set your timeout value
	handler.SlaveId = byte(slaveID)
	// Deepsea `10.6.70.5`

	// Create Modbus client using the handler
	client := modbus.NewClient(handler)

	// Connect to the Modbus server or return a error message
	err := handler.Connect()
	if err != nil {
		fmt.Println("Error connecting to Modbus server:", err)
	}
	defer handler.Close()

	results, err := client.ReadHoldingRegisters(register_value, uint16(number_bytes))
	if err != nil {
		fmt.Println("Error reading input registers:", err)
	}
	return results
	// Display the results
	// fmt.Printf("Read Fuel level of deepsea: %d%%\n", results[1])
}
