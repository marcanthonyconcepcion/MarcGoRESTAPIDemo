package MarcGoRESTAPIDemo

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var expectedRecords []Subscriber
var databaseDUT Records

func init() {
	rand.Seed(time.Now().UnixNano())
}

func setUp() {
	expectedRecords = []Subscriber{
		{
			index:          1,
			emailAddress:   "marcanthonyconcepcion@gmail.com",
			firstName:      "Marc Anthony",
			lastName:       "Concepcion",
			activationFlag: false,
		},
		{
			index:          2,
			emailAddress:   "marcanthonyconcepcion@email.com",
			firstName:      "Marc",
			lastName:       "Concepcion",
			activationFlag: false,
		},
		{
			index:          3,
			emailAddress:   "kevin.andrews@email.com",
			firstName:      "Kevin",
			lastName:       "Andrews",
			activationFlag: false,
		},
	}
	databaseDUT = makeDatabaseRecords()
	noDBConnection := databaseDUT.database.Ping()
	if noDBConnection != nil {
		panic(noDBConnection.Error())
	}
	for _, subscriber := range expectedRecords {
		_, createError := databaseDUT.create(subscriber)
		if createError != nil {
			panic(createError.Error())
		}
	}
}

func tearDown() {
	expectedRecords = nil
	_, truncateFail := databaseDUT.database.Exec("truncate table `subscribers`")
	if truncateFail != nil {
		panic(truncateFail.Error())
	}
	dbCloseFail := databaseDUT.database.Close()
	if dbCloseFail != nil {
		panic(dbCloseFail.Error())
	}
}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func TestCreate(t *testing.T) {
	newRecord := Subscriber{
		index:          4,
		emailAddress:   "riseofskywalker@starwars.com",
		firstName:      "Palpatine",
		lastName:       "Rey",
		activationFlag: false,
	}
	_, createFail := databaseDUT.create(newRecord)
	if createFail != nil {
		t.Errorf("ERROR creating database records. %s", createFail.Error())
	}
	fetchedRecords, listFail := databaseDUT.list()
	if listFail != nil {
		t.Errorf("ERROR fetching database records. %s", listFail.Error())
	}
	updatedExpectedRecords := append(expectedRecords, newRecord)
	for index := range updatedExpectedRecords {
		if updatedExpectedRecords[index] != fetchedRecords[index] {
			t.Errorf("ERROR fetching database records. Expected %v != Actual %v",
				updatedExpectedRecords[index], fetchedRecords[index])
		}
	}
}

func TestRetrieve(t *testing.T) {
	fetchedRecords, listFail := databaseDUT.list()
	if listFail != nil {
		t.Errorf("ERROR fetching database records. %s", listFail.Error())
	}
	for index := range expectedRecords {
		if expectedRecords[index] != fetchedRecords[index] {
			t.Errorf("ERROR fetching database record. Expected %v != Actual %v",
				expectedRecords[index], fetchedRecords[index])
		}
		retrievedRecord, retrieveFail := databaseDUT.retrieve(uint8(index + 1))
		if retrieveFail != nil {
			t.Errorf("ERROR retrieving database record at index %d. %s", index+1, retrieveFail.Error())
		}
		if expectedRecords[index] != *retrievedRecord {
			t.Errorf("ERROR fetching database record. Expected %v != Actual %v",
				expectedRecords[index], retrievedRecord)
		}
	}
}

func TestUpdate(t *testing.T) {
	index := rand.Intn(len(expectedRecords) - 1)
	activationFlag := rand.Intn(1)
	emailAddress := "marcanthonyconcepcion@yeahmail.com"

	updateValues := make(map[string]string)
	updateValues["activation_flag"] = strconv.Itoa(activationFlag)
	updateValues["email_address"] = emailAddress
	_, updateFail := databaseDUT.update(uint8(index+1), updateValues)
	if updateFail != nil {
		t.Errorf("ERROR updating database records. %s", updateFail.Error())
	}
	fetchedRecords, listFail := databaseDUT.list()
	if listFail != nil {
		t.Errorf("ERROR fetching database records. %s", listFail.Error())
	}
	updatedExpectedRecords := expectedRecords
	updatedExpectedRecords[index].activationFlag = activationFlag != 0
	updatedExpectedRecords[index].emailAddress = emailAddress
	for index := range updatedExpectedRecords {
		if updatedExpectedRecords[index] != fetchedRecords[index] {
			t.Errorf("ERROR fetching database records. Expected %v != Actual %v",
				updatedExpectedRecords[index], fetchedRecords[index])
		}
	}
}

func TestDelete(t *testing.T) {
	index := rand.Intn(len(expectedRecords) - 1)
	_, deleteFail := databaseDUT.delete(uint8(index + 1))
	if deleteFail != nil {
		t.Errorf("ERROR deleting database records. %s", deleteFail.Error())
	}
	fetchedRecords, listFail := databaseDUT.list()
	if listFail != nil {
		t.Errorf("ERROR fetching database records. %s", listFail.Error())
	}
	updatedExpectedRecords := append(expectedRecords[:index-1], expectedRecords[index+1:]...)
	for index := range updatedExpectedRecords {
		if updatedExpectedRecords[index] != fetchedRecords[index] {
			t.Errorf("ERROR fetching database records. Expected %v != Actual %v",
				updatedExpectedRecords[index], fetchedRecords[index])
		}
	}
}
