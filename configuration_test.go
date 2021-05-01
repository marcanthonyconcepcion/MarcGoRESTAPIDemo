package MarcGoRESTAPIDemo

import (
	"testing"
)

func TestReadYamlFile(t *testing.T) {
	configuration, fault := readConfiguration("resources/MarcGoRESTAPIDemo.yaml")
	if fault != nil {
		t.Errorf("Error reading configuration file %s.", fault)
	}
	if "localhost" != configuration.Database.Host {
		t.Errorf("Value %s is NOT the expected database host from the config file.", configuration.Database.Host)
	}
	if "subscribers_database" != configuration.Database.DBName {
		t.Errorf("Value %s is NOT the expected database name from the config file.", configuration.Database.DBName)
	}
	if "user" != configuration.Database.User {
		t.Errorf("Value %s is NOT the expected database user from the config file.", configuration.Database.User)
	}
	if "password" != configuration.Database.Password {
		t.Errorf("Value %s is NOT the expected database password from the config file.", configuration.Database.Password)
	}
	if "subscribers" != configuration.MVC.Resource {
		t.Errorf("Value %s is NOT the expected mvc resource from the config file.", configuration.MVC.Resource)
	}
	if "error.log" != configuration.Log.Filename {
		t.Errorf("Value %s is NOT the expected log file name from the config file.", configuration.Log.Filename)
	}
}