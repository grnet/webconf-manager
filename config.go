package main

import (
	"flag"
	"gopkg.in/gcfg.v1"
)

var flConfig = flag.String("conf", "", "specify configuration file")
var flServerIP = flag.String("ip", "", "ip address the server will bind to")
var flServerPort = flag.Int("port", 0, "specify the port to listen on")
var flSqlPath = flag.String("dbpath", "", "specify full path to local sqlite db")
var flMonHostIP = flag.String("monhostip", "", "ip address of monitoring host to query")
var flMonHostUsername = flag.String("monhostuser", "", "username to query monitoring host interface")
var flMonHostPassword = flag.String("monhostpass", "", "password to query monitoring host interface")

type Config struct {
	Server struct {
		Bindip string
		Port   int
	}
	Sql struct {
		Path string
	}
	MonHost struct {
		IP       string
		Username string
		Password string
	}
}

const defaultConfig = `
[server]
bindip = ""
port = 8081

[sql]
path = "/var/lib/webconf/inventory.db"

[monhost]
ip = 83.212.170.52
username = webadmin
password = password
`

func LoadConfiguration() Config {

	flag.Parse()

	cfg := Config{}
	if *flConfig != "" {
		err := gcfg.ReadFileInto(&cfg, *flConfig)
		if err != nil {
			panic(err)
		}
	} else {
		err := gcfg.ReadStringInto(&cfg, defaultConfig)
		if err != nil {
			panic(err)
		}
	}
	if *flServerIP != "" {
		cfg.Server.Bindip = *flServerIP
	}
	if *flServerPort != 0 {
		cfg.Server.Port = *flServerPort
	}
	if *flMonHostIP != "" {
		cfg.MonHost.IP = *flMonHostIP
	}
	if *flSqlPath != "" {
		cfg.Sql.Path = *flSqlPath
	}
	if *flMonHostUsername != "" {
		cfg.MonHost.Username = *flMonHostUsername
	}
	if *flMonHostPassword != "" {
		cfg.MonHost.Password = *flMonHostPassword
	}

	return cfg
}
