package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("\n CPU Temperature: %.2fC and CPU Temperature State: %s\n", cpu_temp, cpu_temp_state)
		fmt.Printf("Host Address: %s\n", host_address)
	}
}
func main() {

	//nodeList := []string{"10.0.34.71:8000"}
	nodeList := []string{"server1:8000"}
	var wg sync.WaitGroup
	wg.Add(len(nodeList))

	for _, node := range nodeList {
		go func() {
			defer wg.Done()
			collectCPUTemperature(node)
		}()

	}

	wg.Wait()
}
