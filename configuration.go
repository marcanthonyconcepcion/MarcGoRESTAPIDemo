package MarcGoRESTAPIDemo

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Configuration struct {
	Database struct {
		Host     string
		Port     uint16
		DBName   string
		User     string
		Password string
	}
	MVC struct {
		Resource string
	}
	Log struct {
		Filename string
	}
}

func readConfiguration(fileName string) *Configuration {
	buffer, fault := ioutil.ReadFile(fileName)
	if fault != nil {
		panic(fault.Error())
	}
	configuration := &Configuration{}
	fault = yaml.Unmarshal(buffer, configuration)
	if fault != nil {
		panic(fault.Error())
	}
	return configuration
}
