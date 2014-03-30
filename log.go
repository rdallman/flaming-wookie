package main

import (
    "io"
    "io/ioutil"
    "log"
    "os"
)

var (
    TRACE   *log.Logger //standard stuff
    INFO    *log.Logger //special information
    WARNING *log.Logger //there is something you need to know about
    ERROR   *log.Logger //something has failed
)

func init() {
	//create log file
    file, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file", err)
    }

    TRACE = log.New(io.MultiWriter(file, ioutil.Discard),
        "TRACE: ",
        log.LstdFlags)

    INFO = log.New(io.MultiWriter(file, os.Stdout),
        "INFO: ",
        log.LstdFlags)

    WARNING = log.New(io.MultiWriter(file, os.Stdout),
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    ERROR = log.New(io.MultiWriter(file, os.Stdout),
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}