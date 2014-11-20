package main

import (
	// "encoding/json"
	"fmt"
	linuxproc "github.com/chilts/gocigar/Godeps/_workspace/src/github.com/c9s/goprocinfo/linux"
	ini "github.com/chilts/gocigar/Godeps/_workspace/src/github.com/vaughan0/go-ini"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

func doDiskSpace() {
	disk, err := linuxproc.ReadDisk("/")
	if err != nil {
		fmt.Printf("Can't update document %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%d %d %d\n", disk.All, disk.Used, disk.Free)
}

func getLoadAvg() (*linuxproc.LoadAvg, error) {
	loadAvg, err := linuxproc.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		return nil, err
	}
	return loadAvg, nil
}

func main() {
	file, err := ini.LoadFile("cigar.ini")
	if err != nil {
		fmt.Printf("Error loading config file : %v\n", err)
		os.Exit(1)
	}

	fmt.Println("----------------------------------------------------------------------------")
	for key, value := range file[""] {
		fmt.Printf("%s => %s\n", key, value)
	}
	for name, section := range file {
		fmt.Printf("Section name: %s\n", name)
		for key, value := range section {
			fmt.Printf("  %s => %s\n", key, value)
		}
	}
	fmt.Println("----------------------------------------------------------------------------")

	doMemInfo()

	var memInfo MemInfo
	memInfo.populate()
	// fmt.Printf("%q\n", memInfo)
	fmt.Printf("Memory usage is %d   [memTotal: %d, memFree: %d]\n", memInfo.memTotal-memInfo.memFree, memInfo.memTotal, memInfo.memFree)

	doDiskSpace()
	loadAvg, err := getLoadAvg()
	if err != nil {
		fmt.Printf("Can't get LoadAvg : %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%q\n", loadAvg)
	fmt.Printf("LoadAvg: OneMin=%f, FiveMin=%f, FifteenMin=%f\n", loadAvg.Last1Min, loadAvg.Last5Min, loadAvg.Last15Min)
	fmt.Println("\n")

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for t := range ticker.C {
			iso8601 := time.Now().UTC().Format(time.RFC3339)
			fmt.Println("UTC: ", iso8601)
			fmt.Println("UTC: ", time.Now().UTC().Format(time.RFC3339Nano))
			fmt.Println("Tick at", t.UTC().Format(time.RFC3339))
			loadAvg, err := getLoadAvg()
			if err != nil {
				fmt.Printf("Can't get LoadAvg : %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("LoadAvg: OneMin=%f, FiveMin=%f, FifteenMin=%f\n", loadAvg.Last1Min, loadAvg.Last5Min, loadAvg.Last15Min)
			fmt.Println()
		}
	}()

	// sleep for 22 seconds to allow 4 ticks
	time.Sleep(22 * time.Second)
	ticker.Stop()

	fmt.Print("Got to the end\n")
}
