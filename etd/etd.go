package etd

import (
	"fmt"
	"github.com/secnot/orderedmap"
	"lirs2/simulator"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type (
	Node struct {
		lba int
		op  string
	}

	LRU struct {
		maxLen        int
		available     int
		qcAvailable   int
		qetdAvailable int
		hit           int
		miss          int
		pageFault     int
		writeCount    int
		readCount     int
		writeCost     float32
		readCost      float32
		eraseCost     float32

		qf    *orderedmap.OrderedMap
		qc    *orderedmap.OrderedMap
		qetd  [][3]int
		clock int
	}
)

func Etd(value int) *LRU {
	lru := &LRU{
		maxLen:        value,
		available:     value,
		qcAvailable:   int(math.Ceil(float64(value) / 10)),
		qetdAvailable: value,
		clock:         0,
		hit:           0,
		miss:          0,
		pageFault:     0,
		writeCount:    0,
		readCount:     0,
		writeCost:     0.25,
		readCost:      0.025,
		eraseCost:     2,
		qf:            orderedmap.NewOrderedMap(),
		qc:            orderedmap.NewOrderedMap(),
		qetd:          make([][3]int, value),
	}
	return lru
}

func checkValueExists(qetd [][3]int, value int) bool {
	for i := 0; i < len(qetd); i++ {
		if qetd[i][0] == value {
			return true
		}
	}
	return false
}

func (lru *LRU) Get(trace simulator.Trace) (err error) {
	// if B in Q
	// update ET
	existInQ := checkValueExists(lru.qetd, trace.Addr)
	if existInQ {
		i := 0
		for i < lru.maxLen {

			i++
		}

		// elif B in Qf
		// NB++
		// if NB > 20
		//  insert B to Q or Qc
	} else if op, ok := lru.qf.Get(trace.Addr); ok {
		lru.hit++

		op2 := op.(string)
		op3 := op2[:1]
		pop3, _ := strconv.Atoi(op2[1:])
		pop3++

		if pop3 > 20 {
			print("\nblock number ", trace.Addr, " reached threshold, moving to Qc")
			if lru.qetdAvailable > 0 {
				lru.qetdAvailable--
				i := 0
				for lru.qetd[i][0] != 0 {
					i++
				}
				lru.qetd[i][0] = trace.Addr
				if op3 == "W" {
					lru.qetd[i][1] = 2
				} else {
					lru.qetd[i][1] = 4
				}
				// temp: set 3 for debugging, use 100 in prod
				lru.qetd[i][2] = 3
				print("\ninserting block ", trace.Addr, " to Qc with ETP of ", lru.qetd[i][2])

			} else if lru.qcAvailable > 0 {
				lru.qcAvailable--
				lru.qc.Set(trace.Addr, op3)
				print("\ninserting block ", trace.Addr, " to Qc")
			} else {
				lru.pageFault++
				if _, op, ok := lru.qc.GetFirst(); ok {

					if op == "W" {
						lru.writeCount++
					}
					lru.qc.PopFirst()
				} else {
					fmt.Println("No elements found to remove")
				}
				lru.qc.Set(trace.Addr, op3)
				print("\npopping Qc then inserting: ", trace.Addr)
			}
			lru.qf.Delete(trace.Addr)
			lru.available++
		} else {
			print("\nblock number ", trace.Addr, " found in Qf, pop: ", pop3)
			op3 = op3 + strconv.Itoa(pop3)
			lru.qf.Set(trace.Addr, op3)
			if ok := lru.qf.MoveLast(trace.Addr); !ok {
				fmt.Printf("Failed to move LBA %d to MRU position\n", trace.Addr)
			}
		}
	} else if _, ok := lru.qc.Get(trace.Addr); ok {
		// elif B in Qc
		//  insert B to Qc
		lru.hit++
		print("\nblock number ", trace.Addr, " found in Qc")
		if ok := lru.qc.MoveLast(trace.Addr); !ok {
			fmt.Printf("Failed to move LBA %d to MRU position\n", trace.Addr)
		}

		return nil
	} else {

		// else put B to Qf
		//  RNG: 33% chance for block to enter Qf
		lru.miss++
		lru.readCount++

		// TODO: change to nonrandom 1/3
		var qfThreshold = 100
		if rand.Intn(100) <= qfThreshold {
			if lru.available > 0 {
				lru.available--
				pop := trace.Op + "1"
				lru.qf.Set(trace.Addr, pop)
				print("\ninserting block ", trace.Addr, " to Qf, pop: ", pop)
			} else {
				lru.pageFault++
				if _, op, ok := lru.qf.GetFirst(); ok {

					op2 := op.(string)
					op3 := op2[:1]

					if op3 == "W" {
						lru.writeCount++
					}
					lru.qf.PopFirst()
				} else {
					fmt.Println("No elements found to remove")
				}
				pop := trace.Op + "1"
				lru.qf.Set(trace.Addr, pop)
				print("\npopping Qf then inserting: ", trace.Addr)
			}
		} else if trace.Op == "W" {
			lru.writeCount++
		}

		return nil
	}
	lru.clock++
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
