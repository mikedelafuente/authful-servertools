package logger

import (
	"log"

	"github.com/mikedelafuente/authful-servertools/pkg/config"
)

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	log.Fatal(v...)
}

func Println(v ...interface{}) {
	log.Println(v...)
}

func Error(v ...interface{}) {
	log.Println("ERROR: ")
	log.Println(v...)
}

func Debug(v ...interface{}) {
	if config.GetConfig().IsDebug {
		log.Println(v...)
	}
}
