package main

import (
	"bufio"
	"fmt"
	"lirs2/lirs"
	"lirs2/simulator"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		traces    []simulator.Trace
		simulator simulator.Simulator
		algorithm string
		filePath  string
		outPath   string
		cacheList []int
		timeStart time.Time
		output    *os.File
		fileInfo  os.FileInfo
		err       error
	)

	if len(os.Args) < 4 {
		fmt.Println("Usage: program_name <algorithm> <trace file> <trace size>")
		os.Exit(1)
	}

	algorithm = os.Args[1]

	filePath = os.Args[2]
	if fileInfo, err = os.Stat(filePath); os.IsExist(err) {
		fmt.Printf("%v does not exist\n", filePath)
		os.Exit(1)
	}

	if traces, err = readFile(filePath); err != nil {
		log.Fatalf("Error reading file %v: %v\n", filePath, err)
	}

	if cacheList, err = validateTraceSize(os.Args[3:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	outPath = fmt.Sprintf("%v_%v_%v.txt", time.Now().Unix(), algorithm, fileInfo.Name())
	if output, err = os.Create(outPath); err != nil {
		log.Fatalf("Error creating file %v: %v\n", outPath, err)
	}

	defer func(output *os.File) {
		err := output.Close()
		if err != nil {
			log.Fatalf("Error closing file %v: %v\n", outPath, err)
		}
	}(output)

	for _, cacheSize := range cacheList {
		switch strings.ToLower(algorithm) {
		case "lirs":
			simulator = lirs.NewLIRS(cacheSize, 10)
		//case "lirs2":
		//	simulator = lirs2.NewLIRS2(cacheSize, 10)
		case "etd":
			print("to be implemented")
		default:
			log.Fatal("Algorithm not supported")
		}

		timeStart = time.Now()

		for _, trace := range traces {
			err = simulator.Get(trace)
			if err != nil {
				log.Fatal(err)
			}
		}

		if err = simulator.PrintToFile(output, timeStart); err != nil {
			log.Fatal(err)
		}

	}
	fmt.Println("Simulation Done")
}

func readFile(filePath string) (traces []simulator.Trace, err error) {
	var (
		file    *os.File
		scanner *bufio.Scanner
		row     []string
		address int
	)

	if file, err = os.Open(filePath); err != nil {
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
				Address:   address,
				Operation: row[1],
			},
		)
	}

	return traces, nil
}

func validateTraceSize(traceSizeList []string) (cacheList []int, err error) {
	var cache int
	for _, size := range traceSizeList {
		cache, err = strconv.Atoi(size)
		if err != nil {
			return cacheList, err
		}
		cacheList = append(cacheList, cache)
	}
	return cacheList, nil
}
