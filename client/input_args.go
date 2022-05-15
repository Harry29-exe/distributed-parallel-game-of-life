package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	host               = "localhost"
	port               = "3333"
	threadCount uint32 = 4
)

func readInputArgs() {
	osArgs := os.Args[1:]
	args := map[string]string{}

	for _, arg := range osArgs {
		eqSignIndex := strings.Index(arg, "=")
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

	argValue, ok = args["threadCount"]
	if ok {
		count, err := strconv.Atoi(argValue)
		if err != nil || count < 1 || count > math.MaxInt16 {
			fmt.Println("Variable boardW (i.e. board width) must be integer between 1 and 2^16-1")
			os.Exit(1)
		}
		threadCount = uint32(count)
	}

}
