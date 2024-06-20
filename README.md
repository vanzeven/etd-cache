# ETD-cache
This is an implementation of the ETD cache replacement algorithm implemented in Golang, based on the paper by Ningwei Dai, et al

paper: https://dl.acm.org/doi/abs/10.1145/2742854.2742881

## Prerequirements
1. Install Go and a good IDE (GoLand or VSCode are good choices)
2. Install Nim
3. Find some traces

## How To Run ?
1. Choose dataset
2. Compile converter ```split_financial.nim``` and ```split_websearch.num``` programs
```
nim c [nama program].nim
```
3. Make directory ```data``` in root directory

4. Running compile on dataset and choose ```data``` as directory target
```
./[nama program] [dataset] ../data/[output]
```
5. Go get module
```
go get github.com/petar/GoLLRB
go get github.com/secnot/orderedmap
```
6. Run ```main.go``` (algo is case insensitive)
```
go run main.go [algorithm(LRU|LFU|ETD)] [file] [trace size]
```

example command
```
go run .\main.go lru .\converter\fin2-src18m 5000
```
