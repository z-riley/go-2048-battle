// Package log wraps the standard library log package. This allows the logging
// functionality to be easily disabled.
package log

import (
	"log"

	"github.com/z-riley/go-2048-battle/config"
)

func Print(v ...any) {
	if config.Debug {
		log.Print(v...)
	}
}

func Printf(format string, v ...any) {
	if config.Debug {
		log.Printf(format, v...)
	}
}

func Println(v ...any) {
	if config.Debug {
		log.Println(v...)
	}
}

func Fatal(v ...any) {
	if config.Debug {
		log.Fatal(v...)
	}
}

func Fatalf(format string, v ...any) {
	if config.Debug {
		log.Fatalf(format, v...)
	}
}

func Fatalln(v ...any) {
	if config.Debug {
		log.Fatalln(v...)
	}
}

func Panic(v ...any) {
	if config.Debug {
		log.Panic(v...)
	}
}

func Panicf(format string, v ...any) {
	if config.Debug {
		log.Panicf(format, v...)
	}
}

func Panicln(v ...any) {
	if config.Debug {
		log.Panicln(v...)
	}
}
