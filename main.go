package main

import (
	"flag"
	"os"

	"github.com/labstack/gommon/color"
)

var Colorer = color.New()
var opts *Options

func init() {
	// Read Flags
	port := flag.Int("port", DefaultParams.ProxyPort, "Reverse proxy port. (Required)")
	prefix := flag.String("prefix", DefaultParams.ProxyPrefix, "route prefix (ex: http://localhost:3000/prefix/")
	flag.Parse()

	// Read RemoteURL
	tail := flag.Args()
	if len(tail) != 1 {
		PrintUsage()
	}

	var err error
	opts, err = NewOptions(tail[0], *port, *prefix)
	if err != nil {
		PrintUsage()
	}
}

func main() {
	Run(opts)
}

func PrintUsage() {
	Colorer.Println("Usage: revproxy [OPTIONS]... <REMOTE_URL>")
	Colorer.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}
