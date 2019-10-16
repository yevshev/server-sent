package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/alexandrevicenzi/go-sse"
)

type CPUTempObj struct {
	TimeStamp   time.Time
	HostAddress string
	CPUTemp     float64
}

func main() {
	// Create the server.
	s := sse.NewServer(nil)
	defer s.Shutdown()

	// Register with /events endpoint.
	http.Handle("/redfish/v1/", s)

	// Dispatch messages to channel-1.
	go func() {
		for {
			cpu_thermal := GetCPUTemp()
			s.SendMessage("/redfish/v1/Chassis/1/Thermal", sse.SimpleMessage(string(cpu_thermal)))
			time.Sleep(1 * time.Second)
		}
	}()

	// // Dispatch messages to channel-2
	// go func() {
	// 	i := 0
	// 	for {
	// 		i++
	// 		s.SendMessage("/events/channel-2", sse.SimpleMessage(strconv.Itoa(i)))
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()
	hostIP := GetNodeIPAddress()
	http.ListenAndServe(hostIP+":8000", nil)
}

func GetNodeIPAddress() string {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("could not obtain host IP address: %v", err)
	}
	ip := ""
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	return ip
}

func randTemperature(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return math.Floor((min+rand.Float64()*(max-min))*100) / 100
}

// CPU temperature
func GetCPUTemp() []byte {
	hostIP := GetNodeIPAddress()
	log.Println("\nCPU temperature\n")

	// Its a mockup CPU temperature
	cpuTempObj := new(CPUTempObj)
	cpuTempObj.TimeStamp = time.Now()
	cpuTempObj.HostAddress = hostIP
	cpuTempObj.CPUTemp = randTemperature(3.0, 98.0)

	jsonObj, err := json.Marshal(cpuTempObj)
	if err != nil {
		log.Println(fmt.Sprintf("Could not marshal the response data: %v", err))
	}
	return jsonObj

}
