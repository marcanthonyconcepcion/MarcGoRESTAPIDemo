/*
 * Copyright (c) 2021.
 * Marc Concepcion
 * marcanthonyconcepcion@gmail.com
 */

package MarcGoRESTAPIDemo

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type SubscriberModelTestFixture struct {
	dut             Records
	expectedRecords []Subscriber
}

func setupSubscriberModelTestFixture() SubscriberModelTestFixture {
	expectedRecords := []Subscriber{
		{
			Index:          1,
			EmailAddress:   "marcanthonyconcepcion@gmail.com",
			FirstName:      "Marc Anthony",
			LastName:       "Concepcion",
			ActivationFlag: false,
		},
		{
			Index:          2,
			EmailAddress:   "marcanthonyconcepcion@email.com",
			FirstName:      "Marc",
			LastName:       "Concepcion",
			ActivationFlag: false,
		},
		{
			Index:          3,
			EmailAddress:   "kevin.andrews@email.com",
			FirstName:      "Kevin",
			LastName:       "Andrews",
			ActivationFlag: false,
		},
	}
	dut := makeDatabaseRecords()
	noDBConnection := dut.database.Ping()
	if noDBConnection != nil {
		panic(noDBConnection.Error())
	}
	for _, subscriber := range expectedRecords {
		_, createError := dut.create(subscriber)
		if createError != nil {
			panic(createError.Error())
		}
	}
	return SubscriberModelTestFixture{dut, expectedRecords}
}

func (fixture SubscriberModelTestFixture) tearDown() {
	fixture.expectedRecords = nil
	_, truncateFail := fixture.dut.database.Exec("truncate table `subscribers`")
	if truncateFail != nil {
		panic(truncateFail.Error())
	}
	dbCloseFail := fixture.dut.database.Close()
	if dbCloseFail != nil {
		panic(dbCloseFail.Error())
	}
}

func TestCreateModel(t *testing.T) {
	fixture := setupSubscriberModelTestFixture()
	newRecord := Subscriber{
		Index:          4,
		EmailAddress:   "riseofskywalker@starwars.com",
		FirstName:      "Palpatine",
		LastName:       "Rey",
		ActivationFlag: false,
	}
	_, createFail := fixture.dut.create(Subscriber{
		EmailAddress: "riseofskywalker@starwars.com",
		FirstName:    "Palpatine",
		LastName:     "Rey",
	})
	if createFail != nil {
		t.Errorf("ERROR creating database records. %s", createFail.Error())
	}
	fetchedRecords, listFail := fixture.dut.list()
	if listFail != nil {
		t.Errorf("ERROR fetching database records. %s", listFail.Error())
	}
	updatedExpectedRecords := append(fixture.expectedRecords, newRecord)
	for index := range updatedExpectedRecords {
		if updatedExpectedRecords[index] != fetchedRecords[index] {
			t.Errorf("ERROR fetching database records. Expected %v != Actual %v",
				updatedExpectedRecords[index], fetchedRecords[index])
		}
	}
	fixture.tearDown()
}

func TestRetrieveModel(t *testing.T) {
	fixture := setupSubscriberModelTestFixture()
	fetchedRecords, listFail := fixture.dut.list()
	if listFail != nil {
		t.Errorf("ERROR fetching database records. %s", listFail.Error())
	}
	for index := range fixture.expectedRecords {
		if fixture.expectedRecords[index] != fetchedRecords[index] {
			t.Errorf("ERROR fetching database record. Expected %v != Actual %v",
				fixture.expectedRecords[index], fetchedRecords[index])
		}
		retrievedRecord, retrieveFail := fixture.dut.retrieve(uint8(index + 1))
		if retrieveFail != nil {
			t.Errorf("ERROR retrieving database record at index %d. %s", index+1, retrieveFail.Error())
		}
		if fixture.expectedRecords[index] != *retrievedRecord {
			t.Errorf("ERROR fetching database record. Expected %v != Actual %v",
				fixture.expectedRecords[index], retrievedRecord)
		}
	}
	fixture.tearDown()
}

func TestUpdateModel(t *testing.T) {
	fixture := setupSubscriberModelTestFixture()
	form := Subscriber{}
	form.Index = uint8(rand.Intn(len(fixture.expectedRecords)) + 1)
	form.FirstName = "Handsome Marc"
	form.EmailAddress = "marchandsome@yeahmail.com"
	_, updateFail := fixture.dut.update(form)
	if updateFail != nil {
		t.Errorf("ERROR updating database records. %s", updateFail.Error())
	}
	fetchedRecords, listFail := fixture.dut.list()
	if listFail != nil {
		t.Errorf("ERROR fetching database records. %s", listFail.Error())
	}
	updatedExpectedRecords := fixture.expectedRecords
	updatedExpectedRecords[form.Index-1].FirstName = form.FirstName
	updatedExpectedRecords[form.Index-1].EmailAddress = form.EmailAddress
	for index := range updatedExpectedRecords {
		if updatedExpectedRecords[index] != fetchedRecords[index] {
			t.Errorf("ERROR fetching database records. Expected %v != Actual %v",
				updatedExpectedRecords[index], fetchedRecords[index])
		}
	}
	fixture.tearDown()
}

func TestDeleteModel(t *testing.T) {
	fixture := setupSubscriberModelTestFixture()
	index := rand.Intn(len(fixture.expectedRecords)) + 1
	_, deleteFail := fixture.dut.delete(uint8(index))
	if deleteFail != nil {
		t.Errorf("ERROR deleting database records. %s", deleteFail.Error())
	}
	fetchedRecords, listFail := fixture.dut.list()
	if listFail != nil {
		t.Errorf("ERROR fetching database records. %s", listFail.Error())
	}
	updatedExpectedRecords := fixture.expectedRecords[1:]
	if index > 1 {
		updatedExpectedRecords = append(fixture.expectedRecords[:index-1], fixture.expectedRecords[index-1+1:]...)
	}

	for index, updatedExpectedRecord := range updatedExpectedRecords {
		if updatedExpectedRecord != fetchedRecords[index] {
			t.Errorf("ERROR fetching database records. Expected %v != Actual %v",
				updatedExpectedRecord, fetchedRecords[index])
		}
	}
	fixture.tearDown()
}

func TestActivateModel(t *testing.T) {
	fixture := setupSubscriberModelTestFixture()
	index := rand.Intn(len(fixture.expectedRecords)) + 1
	activate := rand.Intn(1) == 1
	_, activateFail := fixture.dut.activate(uint8(index), activate)
	if activateFail != nil {
		t.Errorf("ERROR activating a subscriber. %s", activateFail.Error())
	}
	fetchedRecords, listFail := fixture.dut.list()
	if listFail != nil {
		t.Errorf("ERROR fetching database records. %s", listFail.Error())
	}
	fixture.expectedRecords[index-1].ActivationFlag = activate
	for index, expectedRecord := range fixture.expectedRecords {
		if expectedRecord != fetchedRecords[index] {
			t.Errorf("ERROR fetching database records. Expected %v != Actual %v",
				expectedRecord, fetchedRecords[index])
		}
	}
	fixture.tearDown()
}
