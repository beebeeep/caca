package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"gopkg.in/yaml.v2"
)

var config CacaConfig

func fail(txt string, arg ...interface{}) {
	fmt.Printf(txt+"\n", arg...)
	os.Exit(-1)
}

func loadConfig(filename string) (CacaConfig, error) {
	var cfg CacaConfig
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func processCommand(args []string) {
	/*
		var defaultCacus *CacusInstance
		for _, v := range config.Instances {
			if v.Default {
				defaultCacus = &v
				break
			}
		}*/

	switch args[0] {
	case "upload":
		uploadPackage(args[1:])
		break
	default:
		fail("Unknown command '%s'", args[0])
	}
}

func main() {
	me, err := user.Current()
	if err != nil {
		fail("Cannot get current user: %v", err)
	}
	defaultCfgFile := path.Join(me.HomeDir, ".cacarc")

	cfgFile := flag.String("config", defaultCfgFile, "Config file")
	flag.Parse()

	config, err = loadConfig(*cfgFile)
	if err != nil {
		fail("Cannot load config file: %v", err)
	}

	if len(flag.Args()) < 1 {
		fail("Specify command")
	}
	processCommand(flag.Args())
}
