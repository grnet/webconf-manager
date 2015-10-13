package main

import (
	"flag"
	"gopkg.in/gcfg.v1"
)

var flConfig = flag.String("conf", "", "specify configuration file")
var flServerIP = flag.String("ip", "", "ip address the server will bind to")
var flServerPort = flag.Int("port", 0, "specify the port to listen on")

type Config struct {
	Server struct {
		Bindip string
		Port   int
	}
}

const defaultConfig = `
    [server]
    bindip = ""
    port = 8081
`

func LoadConfiguration() Config {

	flag.Parse()

	cfg := Config{}
	if *flConfig != "" {
		err := gcfg.ReadFileInto(&cfg, *flConfig)
		if err != nil {
			panic (err)
		}
	} else {
		_ = gcfg.ReadStringInto(&cfg, defaultConfig)
	}
	if *flServerIP != "" {
		cfg.Server.Bindip = *flServerIP
	}
	if *flServerPort != 0 {
		cfg.Server.Port = *flServerPort
	}

	return cfg
}
