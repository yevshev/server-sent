package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
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

	// Get hostname
	hostIP, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(hostIP+":8000", nil)
}

func randTemperature(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return math.Floor((min+rand.Float64()*(max-min))*100) / 100
}

// CPU temperature
func GetCPUTemp() []byte {
	//hostIP := GetNodeIPAddress()
	hostIP, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	log.Println("CPU temperature")

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
