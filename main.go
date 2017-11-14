package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/BurntSushi/toml"
)

type Node struct {
	Path    string `yaml:"path"`
	Backend struct {
		ServiceName string `yaml:"serviceName"`
		ServicePort int    `yaml:"servicePort"`
	}
}

func loadmeta(configFile string) (meta ConfigMap, err error) {
	if configFile != "" {
		_, err = toml.DecodeFile(configFile, &meta)
		if err != nil {
			return
		}
	}
	return
}

func GetConfigMap(cm ConfigMap, server string) (s service, err error) {
	for _, result := range cm.Service {
		if result.ServiceName == server {
			s = result
			return s, err
		}
	}
	err = errors.New("not found")
	return s, err
}

func FlagSet(opts *Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("cibo", flag.ExitOnError)
	flagSet.Bool("version", false, "print version string")
	flagSet.String("input", "", "path to nginx config file")
	flagSet.String("config", "", "path to config file")
	flagSet.String("output", "result.yaml", "path to config file")
	flagSet.String("outputtmp", "resulttmp.yaml", "path to tmp config file")
	return flagSet
}

func main() {
	opts := NewOptions()
	flagSet := FlagSet(opts)
	flagSet.Parse(os.Args[1:])

	if flagSet.Lookup("version").Value.(flag.Getter).Get().(bool) {
		fmt.Println(VersionString())
		os.Exit(0)
	}

	configFile := flagSet.Lookup("config").Value.String()
	meta, err := loadmeta(configFile)
	if err != nil {
		log.Fatalf("ERROR: failed to load config file %s - %s", configFile, err.Error())
		os.Exit(-1)
	}
	fmt.Println(meta)
	opts.Conf = flagSet.Lookup("input").Value.String()
	opts.OutputFile = flagSet.Lookup("output").Value.String()
	opts.TmpFile = flagSet.Lookup("outputtmp").Value.String()
	takon := NewTakon()
	takon.SetContent(opts.ReadNginxFile(opts.InputFile))
	err = takon.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	var output []interface{}
	for _, location := range takon.Locations {
		if location.TryFiles != "" {
			var node Node
			server, err := GetConfigMap(meta, location.TryFiles)
			if err != nil {
				continue
			}
			node.Path = location.Path
			node.Backend.ServiceName = server.ServiceName
			node.Backend.ServicePort = server.ServicePort
			output = append(output, node)
		}
	}
	d, err := yaml.Marshal(output)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	ioutil.WriteFile(opts.OutputFile, d, 0644)
}
