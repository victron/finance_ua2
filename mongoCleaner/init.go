package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	Debug *log.Logger
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
)

var db = flag.String("db", "", "mandatory arg, DB to use")
var collectionArg = flag.String("collection", "", "mandatory arg, collection to use")
var hours = flag.Int("hours", 8, "clean docs older then num hurs")
var debugOn = flag.Bool("v", false, "verbosity level: DEBUG enable ")

func init()  {
	flag.Parse()

	if *debugOn {
		Debug = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)
	} else {
		Debug = log.New(ioutil.Discard, "DEBUG: ", log.Lshortfile)
	}
	Info = log.New(os.Stdout, "INFO: ", log.Lshortfile)
	Warning = log.New(os.Stdout, "WARNING: ", log.Lshortfile)
	Error = log.New(os.Stdout, "ERROR: ", log.Lshortfile)

	if *db == "" || *collectionArg == ""{
		Error.Fatal("--db and --collection arguments is mandatory")
	}

	
}