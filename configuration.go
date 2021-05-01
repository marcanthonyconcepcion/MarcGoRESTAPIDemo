package MarcGoRESTAPIDemo

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Configuration struct {
	Database struct {
		Host     string
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

func readConfiguration(fileName string) (*Configuration, error) {
	buffer, fault := ioutil.ReadFile(fileName)
	if fault != nil {
		return nil, fault
	}
	configuration := &Configuration{}
	fault = yaml.Unmarshal(buffer, configuration)
	if fault != nil {
		return nil, fmt.Errorf("error reading file %q: %v", fileName, fault)
	}
	return configuration, nil
}
