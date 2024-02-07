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

	"github.com/c-robinson/iplib"
	"github.com/goburrow/modbus"
)

// Config struct to hold the JSON configuration
type Config struct {
	Protocol           string        `json:"protocol"`
	IPmask             string        `json:"ip_mask"`
	IPDevice           string        `json:"gateway"`
	SlaveIDRange       IDRange       `json:"slave_id_range"`
	KnownRegisterRange RegisterRange `json:"known_register_range"`
	LengthOfRead       int           `json:"length_of_each_read"`
	Port               string        `json:"port"`
	FunctionCode       int           `json:"function_code"`
	BaudRates          []int         `json:"baud_rates"`
	Endianness         string        `json:"endianness"`
	SwapBytes          bool          `json:"swap_bytes"`
	SwapWords          bool          `json:"swap_words"`
	Addfunction        bool          `json:"Addfunction"`
}
type RegisterRange struct {
	First int `json:"startvalue"`
	Last  int `json:"endvalue"`
}

// IDRange struct to hold start and end IDs
type IDRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func loadjson(filename string) (Config, error) { //youtube
	var config Config
	configfile, err := os.Open(filename)
	if err != nil {
		return config, err
	}

	jsonParser := json.NewDecoder(configfile)
	err = jsonParser.Decode(&config)
	defer configfile.Close()
	return config, err
}

// ////////////////////////////Start of MAIN//////////////////////////////////////////////
func main() {

	// Load json file
	json_data, _ := loadjson("config.json")

	// couple of importantvariables
	add_function_code := 0 //if register reading should start at 30000-40000
	if json_data.Addfunction {
		add_function_code = json_data.FunctionCode * 10000
	}
	port_number := json_data.Port
	register_length := json_data.LengthOfRead
	register_start_byte := json_data.KnownRegisterRange.First
	register_end_byte := json_data.KnownRegisterRange.Last
	amount_of_registers := register_end_byte - register_start_byte
	fmt.Printf("Length of read %d\nFrom register %d\nTo register %d\n", register_length, register_start_byte+add_function_code, register_end_byte+add_function_code)

	//Starting IP and amount of addresses to read
	ip_mask := convert_IP(json_data.IPmask)
	gateway := convert_IP(json_data.IPDevice) // if gateway can be the device it would be nice
	Total_adresses := ip_mask ^ (0b11111111111111111111111111111111)
	ip_counter := ip_mask & gateway

	for i := 1; i <= Total_adresses; i++ {
		// fmt.Printf("%b\n", ip_counter+i)
		ip_checker := iplib.Uint32ToIP4(uint32((ip_counter) + i))
		bool_port_connection := check_IP_connection(ip_checker.String(), port_number)

		if bool_port_connection && (ip_counter != gateway) { // check if it is the Host
			for j := json_data.SlaveIDRange.Start; j <= json_data.SlaveIDRange.End; j++ {
				for k := 0; k < amount_of_registers; k += register_length {
					data, _, fail := modbusmaker(ip_checker.String(), byte(j), uint16(register_start_byte+k+add_function_code), uint16(register_length))
					if fail == 0 {
						// names, _ := net.LookupAddr(ip_checker.String())
						// fmt.Println("Hostname:", names)
						fmt.Printf("IP addres "+ip_checker.String()+" SlaveID %d Register %d Data %v\n", j, register_start_byte+k, data)
					}

				}

			}

		}
	}
	fmt.Println("End of search hope you found what you are looking for")
	//Deepsea "10.6.70.5" Fuel level 1027
	// fuel_level, _, _ := modbusmaker("10.6.70.5", 0, 1027, 1)
	// fmt.Printf("Read Fuel level of deepsea: %d%%\n", fuel_level[1])

	// // Bluelog `10.6.70.15` Power = 254, Freq = 98, Voltage = 100
	// power_bluelog, _, _ := modbusmaker("10.6.70.15", 1, 254, 2)
	// fmt.Printf("Read Power of bluelog: %v\n", bytesToFloat32(power_bluelog, false, true))

	// SMA Inverter `10.6.70.28` Power = 199, Freq = 201 30775
	power_SMA, _, _ := modbusmaker("10.6.70.28", 126, 40199, 1)
	power := (uint16(power_SMA[0]) << 8) + uint16(power_SMA[1])
	fmt.Printf("Read power of SMA %f\n", float32(power/100))

	// for i := 40500; i < 41000; i += 2 {
	// 	power_SMA, _, _ := modbusmaker("10.6.70.28", 126, uint16(i), 2)
	// 	fmt.Printf("Read power of register %d in SMA %v\n", i, power_SMA)
	// }
}

///////////////////////////////////END OF MAIN//////////////////////////////////////////////////////////////

func bytesToFloat32(bytes []byte, swap_bytes bool, swap_words bool) float32 { //not sure if this works
	if swap_bytes {
		temp_byte := bytes[1]
		bytes[1] = bytes[0]
		bytes[0] = temp_byte
	}
	if swap_words {
		temp_byte1 := bytes[0]
		temp_byte2 := bytes[1]
		bytes[0] = bytes[2]
		bytes[1] = bytes[3]
		bytes[2] = temp_byte1
		bytes[3] = temp_byte2
	}
	bits := binary.BigEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func modbusmaker(IP_mask string, slaveID byte, register_value uint16, number_bytes uint16) ([]byte, error, int) { //his is the man doing the talking
	// i = 1 connection error i=2 read error and i=3 both error
	i := 0
	address := IP_mask + ":502"
	handler := modbus.NewTCPClientHandler(address)

	handler.Timeout = time.Second / 10
	handler.SlaveId = slaveID

	client := modbus.NewClient(handler)

	err := handler.Connect()
	if err != nil {
		// fmt.Println(err)
		i++
	}
	defer handler.Close()

	results, err := client.ReadHoldingRegisters(register_value, number_bytes)
	if err != nil {
		// fmt.Println(err)
		i = i + 2
	}
	return results, err, i
}

func check_IP_connection(host string, port string) bool {
	timeout := time.Second / 50
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout) //not sure thank you google
	boolean := false
	if err != nil {
		boolean = false
	}
	if conn != nil {
		boolean = true
	}
	// name, _ := net.LookupCNAME(host)
	// fmt.Println(name)
	// net.LookupCNAME()
	return boolean
}

func convert_IP(IP string) int { //convert from string to integer
	octets := strings.Split(IP, ".")
	octet0, _ := strconv.Atoi(octets[0])
	octet1, _ := strconv.Atoi(octets[1])
	octet2, _ := strconv.Atoi(octets[2])
	octet3, _ := strconv.Atoi(octets[3])
	return (octet0 << 24) | (octet1 << 16) | (octet2 << 8) | octet3
}
