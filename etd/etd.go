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
		lba        int
		op         string
		popularity int
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

func (lru *LRU) Put(data *Node) (exists bool) {
	// if B in Q
	// update ET

	// elif B in Qf
	// NB++
	// if NB > 20
	//  insert B to Qc

	// elif B in Qc
	//  insert B to Qc

	if _, ok := lru.orderedList.Get(data.lba); ok {
		lru.hit++

		if ok := lru.orderedList.MoveLast(data.lba); !ok {
			fmt.Printf("Failed to move LBA %d to MRU position\n", data.lba)
		}
		return true
		// else put B to Qf
	} else {
		//  RNG: 33% chance for block to enter Qf
		lru.miss++
		lru.readCount++

		var qfThreshold = 33
		if rand.Intn(100) <= qfThreshold {
			node := &Node{
				lba:        data.lba,
				op:         data.op,
				popularity: 1,
			}
			if lru.available > 0 {
				lru.available--
				lru.orderedList.Set(data.lba, node)
			} else {
				lru.pageFault++
				if _, firstValue, ok := lru.orderedList.GetFirst(); ok {
					lruLba := firstValue.(*Node)

					if lruLba.op == "W" {
						lru.writeCount++
					}
					lru.orderedList.PopFirst()
				} else {
					fmt.Println("No elements found to remove")
				}
				lru.orderedList.Set(data.lba, node)
			}
		} else if data.op == "W" {
			lru.writeCount++
		}

		return false
	}
}

func (lru *LRU) Get(trace simulator.Trace) (err error) {
	obj := new(Node)
	obj.lba = trace.Addr
	obj.op = trace.Op
	lru.Put(obj)

	return nil
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
