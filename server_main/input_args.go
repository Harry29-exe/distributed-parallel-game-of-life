package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	host                            = "localhost"
	port                            = "3333"
	protocol                        = "tcp"
	boardW            uint32        = 32
	boardH            uint32        = 32
	outputFilePath                  = "/home/kamil/GolandProjects/distributed-parallel-game-of-life/gol.out"
	programIterations               = 10
	delay                           = false
	delayTime         time.Duration = 0 // delay in milliseconds between sending board parts
)

func readInputArgs() {
	osArgs := os.Args[1:]
	args := map[string]string{}

	for _, arg := range osArgs {
		eqSignIndex := strings.Index(arg, "=")
		if eqSignIndex < 0 {
			continue
		}
		argName, argValue := arg[0:eqSignIndex], arg[eqSignIndex+1:]
		args[argName] = argValue
	}

	argValue, ok := args["host"]
	if ok {
		host = argValue
	}

	argValue, ok = args["port"]
	if ok {
		port = argValue
	}

	argValue, ok = args["protocol"]
	if ok {
		protocol = argValue
	}

	argValue, ok = args["boardW"]
	if ok {
		w, err := strconv.Atoi(argValue)
		if err != nil || w < 1 || w > math.MaxUint32 {
			fmt.Println("Variable boardW (i.e. board width) must be integer between 1 and 2^32-1")
			os.Exit(1)
		}
		boardW = uint32(w)
	}

	argValue, ok = args["boardH"]
	if ok {
		h, err := strconv.Atoi(argValue)
		if err != nil || h < 1 || h > math.MaxUint32 {
			fmt.Println("Variable boardH (i.e. board height) must be integer between 1 and 2^32-1")
			os.Exit(1)
		}
		boardH = uint32(h)
	}

	argValue, ok = args["outputFilePath"]
	if ok {
		outputFilePath = argValue
	}

	argValue, ok = args["delay"]
	if ok {
		delayVal, err := strconv.ParseBool(argValue)
		if err != nil {
			fmt.Printf("Delay argument must be valid boolean value, but is: %s\n", argValue)
			os.Exit(1)
		}
		delay = delayVal
	}

	argValue, ok = args["delayTime"]
	if ok {
		if delay == false {
			fmt.Printf("Delay time can only be set when delay arg is set to true")
			os.Exit(1)
		}

		delayTimeVal, err := strconv.Atoi(argValue)
		if err != nil {
			fmt.Printf("Delat time arg must be integer value but is: %s", argValue)
			os.Exit(1)
		}

		delayTime = time.Duration(delayTimeVal) * time.Millisecond
	}

	argValue, ok = args["programIterations"]
	if ok {
		iterations, err := strconv.Atoi(argValue)
		if err != nil || iterations < 1 || iterations > math.MaxUint32 {
			fmt.Println("Variable program iterations (i.e. board height) must be integer between 1 and 2^32-1")
			os.Exit(1)
		}
		programIterations = iterations
	}

}

func printArgs() {
	fmt.Println("Starting server with following arguments:")
	fmt.Println("\thost=" + host)
	fmt.Println("\tport=" + port)
	fmt.Println("\tprotocol=" + protocol)
	fmt.Println("\tboardW=" + strconv.Itoa(int(boardW)))
	fmt.Println("\tboardH=" + strconv.Itoa(int(boardH)))
	fmt.Println("\toutputFilePath=" + outputFilePath)
	fmt.Println("\tprogramIterations=" + strconv.Itoa(programIterations))
	fmt.Println("\tdelay=" + strconv.FormatBool(delay))
	fmt.Println("\tdelayTime=" + strconv.Itoa(int(delayTime)))

	fmt.Print("\n\n")
}
