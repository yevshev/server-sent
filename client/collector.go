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
		fmt.Printf("\nNode : %s\n", host_address)
		fmt.Printf("CPU Temp: %.2fC\n", cpu_temp)
		fmt.Printf("CPU State: %s\n", cpu_temp_state)
	}
}
func main() {

	nodeList := [5]string{
		"server1:8000", "server2:8000", "server3:8000", "server4:8000",
		"server5:8000", "server6:8000", "server7:8000", "server8:8000",
		"server9:8000", "server10:8000", "server11:8000", "server12:8000",
		"server13:8000", "server14:8000", "server15:8000", "server16:8000",
		"server17:8000", "server18:8000", "server19:8000", "server20:8000",
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
