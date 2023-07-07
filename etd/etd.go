package etd

import (
	"fmt"
	"github.com/secnot/orderedmap"
	"lirs2/simulator"
	"math/rand"
	"os"
	"time"
)

type (
	Node struct {
		lba int
		op  string
	}

	LRU struct {
		maxLen     int
		available  int
		hit        int
		miss       int
		pageFault  int
		writeCount int
		readCount  int
		writeCost  float32
		readCost   float32
		eraseCost  float32

		orderedList *orderedmap.OrderedMap
	}
)

func Etd(value int) *LRU {
	lru := &LRU{
		maxLen:      value,
		available:   value,
		hit:         0,
		miss:        0,
		pageFault:   0,
		writeCount:  0,
		readCount:   0,
		writeCost:   0.25,
		readCost:    0.025,
		eraseCost:   2,
		orderedList: orderedmap.NewOrderedMap(),
	}
	return lru
}

//func (lru *LRU) Put(data *Node) (exists bool) {

//}

func (lru *LRU) Get(trace simulator.Trace) (err error) {
	// if B in Q
	// update ET

	// elif B in Qf
	// NB++
	// if NB > 20
	//  insert B to Qc

	if _, ok := lru.orderedList.Get(trace.Addr); ok {
		lru.hit++
		print(trace.Addr)
		print(trace.Op)

		if ok := lru.orderedList.MoveLast(trace.Addr); !ok {
			fmt.Printf("Failed to move LBA %d to MRU position\n", trace.Addr)
		}
		return nil

		//} else if true {
		//	// elif B in Qc
		//	//  insert B to Qc
		//	return false
	} else {
		// else put B to Qf
		//  RNG: 33% chance for block to enter Qf

		print(trace.Addr)
		//data.op = "T"
		print(trace.Op + " ")
		print("miss\n")

		lru.miss++
		lru.readCount++

		var qfThreshold = 100
		if rand.Intn(100) <= qfThreshold {
			if lru.available > 0 {
				lru.available--
				print("insert to lru\n")
				print(trace.Addr)
				print(trace.Op)
				print("\n-----\n")
				lru.orderedList.Set(trace.Addr, trace.Op)
			} else {
				lru.pageFault++
				if _, op, ok := lru.orderedList.GetFirst(); ok {

					op2 := op.(string)
					print("LRU is full do operation: ", op2, "\n")
					if op == "W" {
						lru.writeCount++
					}
					lru.orderedList.PopFirst()
				} else {
					fmt.Println("No elements found to remove")
				}
				lru.orderedList.Set(trace.Addr, trace.Op)
			}
		} else if trace.Op == "W" {
			lru.writeCount++
		}

		return nil
	}
}

func (lru LRU) PrintToFile(file *os.File, timeStart time.Time) (err error) {
	file.WriteString(fmt.Sprintf("cache size: %d\n", lru.maxLen))
	file.WriteString(fmt.Sprintf("cache hit: %d\n", lru.hit))
	file.WriteString(fmt.Sprintf("cache miss: %d\n", lru.miss))
	file.WriteString(fmt.Sprintf("write count: %d\n", lru.writeCount))
	file.WriteString(fmt.Sprintf("read count: %d\n", lru.readCount))
	file.WriteString(fmt.Sprintf("hit ratio: %8.4f\n", (float32(lru.hit)/float32(lru.hit+lru.miss))*100))
	file.WriteString(fmt.Sprintf("runtime: %8.4f\n", float32(lru.readCount)*lru.readCost+float32(lru.writeCount)*(lru.writeCost+lru.eraseCost)))
	file.WriteString(fmt.Sprintf("time execution: %8.4f\n\n", time.Since(timeStart).Seconds()))

	return nil
}
