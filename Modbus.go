package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/c-robinson/iplib/v2"
	"github.com/goburrow/modbus"
)

// Config struct to hold the JSON configuration
type Config struct {
	Protocol           string        `json:"protocol"`
	IPmask             string        `json:"ip_mask"`
	IPDevice           string        `json:"ip_device"`
	SlaveIDRange       IDRange       `json:"slave_id_range"`
	KnownRegisterRange RegisterRange `json:"known_register_range"`
	LengthOfEachRead   int           `json:"length_of_each_read"`
	Port               string        `json:"port"`
	FunctionCode       int           `json:"function_code"`
	BaudRates          []int         `json:"baud_rates"`
	Endianness         string        `json:"endianness"`
}
type RegisterRange struct {
	first int `json:"startvalue"`
	last  int `json:"endvalue"`
}

// IDRange struct to hold start and end IDs
type IDRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func loadjson(filename string) (Config, error) {
	var config Config
	configfile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer configfile.Close()
	jsonParser := json.NewDecoder(configfile)
	err = jsonParser.Decode(&config)
	return config, err
}

// ////////////////////////////Start of MAIN//////////////////////////////////////////////
func main() {
	//Load json file
	json_data, _ := loadjson("config.json")
	check_IP_connection("10.6.70.5", "502")
	value := json_data.LengthOfEachRead
	fmt.Println(value)
	//Actual code
	ip_mask := convert_IP(json_data.IPmask)
	ip_device := convert_IP(json_data.IPDevice)
	Total_adresses := ip_mask ^ (0b11111111111111111111111111111111)
	fmt.Printf("Total number of adresses to scan: %d\n", Total_adresses)
	ip_counter := ip_mask & ip_device
	port_number := json_data.Port

	for i := 1; i <= Total_adresses; i++ {
		// fmt.Printf("%b\n", ip_counter+i)
		ip_checker := iplib.Uint32ToIP4(uint32((ip_counter) + i))
		bool_port_connection := check_IP_connection(ip_checker.String(), port_number)

		if bool_port_connection && (ip_counter != ip_device) { // check if it is the Host
			for j := json_data.SlaveIDRange.Start; j <= json_data.SlaveIDRange.End; j++ {
				data, err := modbusmaker(ip_checker.String(), byte(j), 101, 1)
				fmt.Printf(ip_checker.String()+" %v %d\n", data, err)
			}

		}
		// if !bool_port_connection {
		// 	fmt.Printf("number %d the data is %v\n", i, "lost")
		// }
	}
	fmt.Println("Made it")
	//connect to the deepsea
	fuel_level, _ := modbusmaker("10.6.70.5", 0, 1027, 1)
	fmt.Printf("Read Fuel level of deepsea: %d%%\n", fuel_level[1])

	//Connect to new Bluelog
	// Bluelog `10.6.70.15` Power = 254, Freq = 98, Voltage = 100
	power_bluelog, _ := modbusmaker("10.6.70.15", 1, 254, 2)
	fmt.Printf("Read Power of bluelog: %v\n", bytesToFloat32(power_bluelog))

	// SMA Inverter `10.6.70.28` Power = 199, Freq = 201
	power_SMA, _ := modbusmaker("10.6.70.28", 0, 199, 2)
	fmt.Printf("Read power of SMA %v\n", power_SMA)
}

///////////////////////////////////END OF MAIN//////////////////////////////////////////////////////////////

func bytesToFloat32(bytes []byte) float32 {
	// Assuming Big-endian byte order, adjust accordingly if needed
	bits := binary.BigEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func modbusmaker(IP_mask string, slaveID byte, register_value uint16, number_bytes uint16) ([]byte, int) {
	// i = 1 connection i=2 read and i=3 both error
	i := 0
	address := IP_mask + ":502"
	handler := modbus.NewTCPClientHandler(address)

	handler.Timeout = time.Second / 10 // Set your timeout value
	handler.SlaveId = slaveID
	// Deepsea `10.6.70.5`

	// Create Modbus client using the handler
	client := modbus.NewClient(handler)

	// Connect to the Modbus server or return a error message
	err := handler.Connect()
	if err != nil {
		fmt.Println(err)
		// fmt.Println("Connection error")
		i++
	}
	defer handler.Close()

	results, err := client.ReadHoldingRegisters(register_value, number_bytes)
	if err != nil {
		fmt.Println(err)
		// fmt.Println("Reading error")
		i += 2
	}
	return results, i
	// Display the results
	// fmt.Printf("Read Fuel level of deepsea: %d%%\n", results[1])
}

func check_IP_connection(host string, port string) bool {
	timeout := time.Second / 50
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	boolean := false
	if err != nil {
		// fmt.Print("Closed\n")
		boolean = false
	}
	if conn != nil {
		// fmt.Print("Open\n")
		boolean = true
	}
	return boolean
}

func convert_IP(IP string) int {
	octets := strings.Split(IP, ".")
	octet0, _ := strconv.Atoi(octets[0])
	octet1, _ := strconv.Atoi(octets[1])
	octet2, _ := strconv.Atoi(octets[2])
	octet3, _ := strconv.Atoi(octets[3])
	return (octet0 << 24) | (octet1 << 16) | (octet2 << 8) | octet3
}
