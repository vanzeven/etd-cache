package main

import (
	"bufio"
	"fmt"
	"lirs2/simulator"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		traces    []simulator.Trace = make([]simulator.Trace, 0)
		simulator simulator.Simulator
		timeStart time.Time
		out       *os.File
		fs        os.FileInfo•
	filePath  string
	outPath   string
	algorithm string
	err       error
	cacheList []int
	)

	if len(os.Args) < 4 {
		fmt.Println("program [algorithm (LFU|LRU|LRUWSR)] [file trace] [cache size]...")
		os.Exit(1)
	}

	algorithm = os.Args[1]

	filePath = os.Args[2]
	if fs, err = os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("%v does not exists\n", filePath)
		os.Exit(1)
	}

	cacheList, err = validateCacheSize(os.Args[3:])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	traces, err = readFile(filePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	outPath = fmt.Sprintf("%v_%v_%v.txt", time.Now().Unix(), algorithm, fs.Name())

	out, err = os.Create("output" + "/" + outPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()

	for _, cache := range cacheList {
		switch strings.ToLower(algorithm) {
		case "lfu":
			simulator = lfu.NewLFU(cache)
		case "lru":
			simulator = lru.NewLRU(cache)
		case "lruwsr":
			simulator = lruwsr.NewLRUWSR(cache)
		default:
			log.Fatal("algorithm not supported")
		}

		timeStart = time.Now()

		for _, trace := range traces {
			err = simulator.Get(trace)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		simulator.PrintToFile(out, timeStart)
	}
	fmt.Println("Done")
}

func validateCacheSize(tracesize []string) (sizeList []int, err error) {
	var (
		cacheList []int
		cache     int
	)

	for _, size := range tracesize {
		cache, err = strconv.Atoi(size)
		if err != nil {
			return sizeList, err
		}
		cacheList = append(cacheList, cache)
	}
	return cacheList, nil
}

func readFile(filePath string) (traces []simulator.Trace, err error) {
	var (
		file    *os.File
		scanner *bufio.Scanner
		row     []string
		address int
	)
	file, err = os.Open(filePath)
	if err != nil {
		return traces, err
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		row = strings.Split(scanner.Text(), ",")
		address, err = strconv.Atoi(row[0])
		if err != nil {
			return traces, err
		}
		traces = append(traces,
			simulator.Trace{
				Addr: address,
				Op:   row[1],
			},
		)
	}

	return traces, nil
}