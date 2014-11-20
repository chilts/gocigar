package main

import (
	// "encoding/json"
	"fmt"
	linuxproc "github.com/chilts/cigar/Godeps/_workspace/src/github.com/c9s/goprocinfo/linux"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type MemInfo struct {
	memTotal   uint64
	memFree    uint64
	buffers    uint64
	cached     uint64
	swapCached uint64
}

func (memInfo *MemInfo) populate() {
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
		// fmt.Printf("Fields are : %q\n", fields)
		var val uint64
		if fields[0] == "MemTotal:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			memInfo.memTotal = val * 1024
		}
		if fields[0] == "MemFree:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			memInfo.memFree = val * 1024
		}
		if fields[0] == "Buffers:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			memInfo.buffers = val * 1024
		}
		if fields[0] == "Cached:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			memInfo.cached = val * 1024
		}
		if fields[0] == "SwapCached:" {
			val, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return
			}
			memInfo.swapCached = val * 1024
		}
	}
	return
}

func doMemInfo() {
	meminfo, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Fatal(err)
	}

	var memTotal = meminfo["MemTotal"] * 1024
	var memFree = meminfo["MemFree"] * 1024
	var memUsed = memTotal - memFree
	fmt.Printf("Memory usage is %d   [memTotal: %d, memFree: %d]\n", memUsed, memTotal, memFree)

	// let's do some JSON
	// str, err := json.Marshal(meminfo)
	// fmt.Printf("json=%s\n", str)
}

func main() {
	doMemInfo()

	var memInfo MemInfo
	memInfo.populate()
	// fmt.Printf("%q\n", memInfo)
	fmt.Printf("Memory usage is %d   [memTotal: %d, memFree: %d]\n", memInfo.memTotal-memInfo.memFree, memInfo.memTotal, memInfo.memFree)
}
