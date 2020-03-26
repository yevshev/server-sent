package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/donovanhide/eventsource"
)

type CPUTempObj struct {
	TimeStamp   time.Time
	HostAddress string
	CPUTemp     float64
}

func lambdaStateDiscovery(v CPUTempObj) (string, float64, string, string) {
	cpu_temp := v.CPUTemp
	cpu_temp_state := "CPU_TEMP_NONDETERMINISTIC"
	host_address := v.HostAddress
	timestamp := v.TimeStamp.Format(time.StampNano)

	if cpu_temp <= 3 || cpu_temp >= 98 {
		cpu_temp_state = "CPU_TEMP_CRITICAL"
	} else if cpu_temp >= 93 && cpu_temp < 98 {
		cpu_temp_state = "CPU_TEMP_HIGH"
	} else if cpu_temp > 3 && cpu_temp < 93 {
		cpu_temp_state = "CPU_TEMP_OK"
	}
	return timestamp, cpu_temp, cpu_temp_state, host_address

}

func collectCPUTemperature(hostname string) {

	stream, err := eventsource.Subscribe("http://"+hostname+"/redfish/v1/Chassis/1/Thermal", "")
	if err != nil {
		return
	}

	for {
		ev := <-stream.Events
		var result CPUTempObj
		json.Unmarshal([]byte(ev.Data()), &result)
		timestamp, cpu_temp, cpu_temp_state, host_address := lambdaStateDiscovery(result)
		fmt.Printf("%s %s %.2fC %s\n", timestamp, host_address, cpu_temp, cpu_temp_state)
	}
}

func main() {

	// Poll 50 servers
	var nodeList [50]string

	// Fill array with server hostnames
	for i := range nodeList {
		var nodeNum = strconv.Itoa(i)
		nodeList[i] = "server" + nodeNum + ":8000"
	}

	var wg sync.WaitGroup

	for _, node := range nodeList {
		wg.Add(1)
		go func(nodeAddress string) {
			defer wg.Done()
			collectCPUTemperature(nodeAddress)
		}(node)
	}

	wg.Wait()
}
