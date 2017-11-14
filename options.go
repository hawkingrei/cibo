package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type ConfigMap struct {
	Service []service
}

type service struct {
	ServiceName string `yaml:"ServiceName"`
	ServicePort int    `yaml:"ServicePort"`
}

type Options struct {
	InputFile  string
	Conf       string
	OutputFile string
	TmpFile    string
}

type Backend struct {
	server string
	port   int
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) ReadNginxFile(path string) (results []string) {
	fmt.Println(path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return results
}

func (o *Options) ReadConf() (result []Backend, err error) {
	_, err = toml.DecodeFile(o.Conf, &result)
	if err != nil {
		return
	}
	return
}
