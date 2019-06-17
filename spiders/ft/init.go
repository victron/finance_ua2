package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

const daysDefault int = 7
var (
	Debug *log.Logger
	Info *log.Logger
	Warning *log.Logger
	Error *log.Logger
)

var url = flag.String("url", "", "url for getting data (mandatory arg)")
var config = flag.String("config", "", "json template for body request(mandatory arg)")
var days = flag.Int("days", daysDefault, "requested days stat, overriding value in template")

var db = flag.String("db", "", "DB to use (mandatory arg)")
var collectionArg = flag.String("collection", "", "collection to save data (mandatory arg)")
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

	if *db == "" || *collectionArg == "" || *url == "" || *config == ""{
		flag.Usage()
		Error.Fatal("some arguments is mandatory")
	}

	
}