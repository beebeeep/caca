package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"gopkg.in/resty.v0"
	"gopkg.in/yaml.v2"
)

var config CacaConfig

func fail(txt string, arg ...interface{}) {
	fmt.Printf(txt+"\n", arg...)
	os.Exit(-1)
}

func cacaerr(txt string, arg ...interface{}) error {
	return errors.New(fmt.Sprintf(txt+"\n", arg...))
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
	case "show":
		showDistro(args[1:])
	case "search":
		searchPackages(args[1:])
	case "copy":
		copyPackage(args[1:])
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
	cacusInstanceName := flag.String("instance", "", "Cacus instance")
	debug := flag.Bool("d", false, "Cacus instance")
	flag.Parse()

	config, err = loadConfig(*cfgFile)
	if err != nil {
		fail("Cannot load config file: %v", err)
	}

	var cacus *CacusInstance
	for name, instance := range config.Instances {
		if instance.Default && *cacusInstanceName == "" {
			cacus = &instance
			break
		}
		if *cacusInstanceName != "" && name == *cacusInstanceName {
			cacus = &instance
			break
		}
	}

	if cacus != nil {
		// found configuration, set up base URL and auth header
		resty.SetHostURL(fmt.Sprintf("%s/api/v1", cacus.BaseURL))
		resty.SetHeader("Authorization", "Bearer "+cacus.Token)
		if len(cacus.CaCert) > 0 {
			resty.SetRootCertificate(cacus.CaCert)
		}
		resty.SetDebug(*debug)
	} else {
		fail("Cannot find cacus instance")
	}

	if len(flag.Args()) < 1 {
		fail("Specify command: upload, show, search, copy")
	}
	processCommand(flag.Args())
}
