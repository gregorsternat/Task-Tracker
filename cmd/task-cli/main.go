package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"task-tracker/pkg/config"
)

func main() {
	taskFilePtr := flag.String("file", "", "Overwrite the default json file where task are stored")
	logLevelPtr := flag.String("log", "", "Overwrite the default log level")
	flag.Parse()

	conf, err := config.LoadConfig(*taskFilePtr, *logLevelPtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.SetPrefix("[" + conf.LogLevel.String() + "] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Printf("Using task file: %s\n", conf.TaskFilePath)
	log.Printf("Log level set to: %s\n", conf.LogLevel)
}
