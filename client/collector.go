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

func lambdaStateDiscovery(v CPUTempObj) (float64, string, string) {
	cpu_temp := v.CPUTemp
	cpu_temp_state := "CPU_TEMP_NONDETERMINISTIC"
	host_address := v.HostAddress

	if cpu_temp <= 3 || cpu_temp >= 98 {
		cpu_temp_state = "CPU_TEMP_CRITICAL"
	} else if cpu_temp >= 93 && cpu_temp < 98 {
		cpu_temp_state = "CPU_TEMP_HIGH"
	} else if cpu_temp > 3 && cpu_temp < 93 {
		cpu_temp_state = "CPU_TEMP_OK"
	}
	return cpu_temp, cpu_temp_state, host_address

}

func collectCPUTemperature(nodeIP string) {

	stream, err := eventsource.Subscribe("http://"+nodeIP+"/redfish/v1/Chassis/1/Thermal", "")
	if err != nil {
		return
	}

	for {
		ev := <-stream.Events
		var result CPUTempObj
		json.Unmarshal([]byte(ev.Data()), &result)
		cpu_temp, cpu_temp_state, host_address := lambdaStateDiscovery(result)
		fmt.Printf("%s %.2fC %s\n", host_address, cpu_temp, cpu_temp_state)
	}
}
func main() {

	var nodeList [100]string

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
