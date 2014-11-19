package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// type MemInfo struct {
//     memTotal uint64
//     memFree uint64
//     buffers uint64
//     cached uint64
//     swapCached uint64
// }

func getMemorySample() (memTotal, memFree, buffers, cached, swapCached uint64, err error) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			return
		}
		fmt.Printf("Fields are : %q\n", fields)
		var val uint64
		if fields[0] == "MemTotal:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			memTotal = val * 1024
		}
		if fields[0] == "MemFree:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			memFree = val * 1024
		}
		if fields[0] == "Buffers:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			buffers = val * 1024
		}
		if fields[0] == "Cached:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			cached = val * 1024
		}
		if fields[0] == "SwapCached:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			swapCached = val * 1024
		}
	}
	return
}

func main() {
	memTotal, memFree, _, _, _, err := getMemorySample()
	if err != nil {
		//     t.Fatal("/proc/meminfo read fail")
	}
	fmt.Printf("Memory usage is %d   [memTotal: %d, memFree: %d]\n", memTotal-memFree, memTotal, memFree)
}
