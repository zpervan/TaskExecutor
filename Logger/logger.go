package logger

import (
	"log"
	"os"
)

var (
	warnings    *log.Logger
	information *log.Logger
	errors      *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	warnings = log.New(file, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	information = log.New(file, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	errors = log.New(file, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Warn(msg string) {
	warnings.Println(msg)
	log.Println(warnings.Prefix() + " " + msg)
}

func Info(msg string) {
	information.Println(msg)
	log.Println(information.Prefix() + " " + msg)
}

func Error(msg string) {
	errors.Println(msg)
	log.Println(errors.Prefix() + " " + msg)
}
