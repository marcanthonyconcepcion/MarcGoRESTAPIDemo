/*
* Copyright (c) 2021.
* Marc Concepcion
* marcanthonyconcepcion@gmail.com
 */

package MarcGoRESTAPIDemo

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type SubscriberControllerTestFixture struct {
	dut             SubscriberController
	model           Records
	expectedRecords []Subscriber
}

func setupSubscriberControllerTestFixture() SubscriberControllerTestFixture {
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
	model := makeDatabaseRecords()
	dut := makeSubscriberController(model)
	for _, subscriber := range expectedRecords {
		_, createError := model.create(subscriber)
		if createError != nil {
			panic(createError.Error())
		}
	}
	return SubscriberControllerTestFixture{dut, model, expectedRecords}
}

func (fixture SubscriberControllerTestFixture) tearDown() {
	fixture.expectedRecords = nil
	_, truncateFail := fixture.model.database.Exec("truncate table `subscribers`")
	if truncateFail != nil {
		panic(truncateFail.Error())
	}
	dbCloseFail := fixture.model.database.Close()
	if dbCloseFail != nil {
		panic(dbCloseFail.Error())
	}
}

func TestListController(t *testing.T) {
	fixture := setupSubscriberControllerTestFixture()
	request, fault := http.NewRequest("GET", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(fixture.dut.list)
	handler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expectedMessage := ConvertToJson(fixture.expectedRecords)
	if response.Body.String() != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Body.String(), expectedMessage)
	}
	fixture.tearDown()
}

func TestCreateController(t *testing.T) {
	fixture := setupSubscriberControllerTestFixture()
	newSubscriber := Subscriber{
		Index:          4,
		EmailAddress:   "riseofskywalker@starwars.com",
		FirstName:      "Rey",
		LastName:       "Palpatine",
		ActivationFlag: false,
	}
	request, fault := http.NewRequest("POST", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	query := request.URL.Query()
	query.Add("email_address", newSubscriber.EmailAddress)
	query.Add("first_name", newSubscriber.FirstName)
	query.Add("last_name", newSubscriber.LastName)
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	createHandler := http.HandlerFunc(fixture.dut.create)
	createHandler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	updatedExpectedRecords := append(fixture.expectedRecords, newSubscriber)
	expectedMessage := ConvertToJson(updatedExpectedRecords)
	listHandler := http.HandlerFunc(fixture.dut.list)
	listHandler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if response.Body.String() != expectedMessage {
		t.Errorf("createHandler returned unexpected body: got %v want %v", response.Body.String(), expectedMessage)
	}
	fixture.tearDown()
}

func TestRetrieveController(t *testing.T) {
	fixture := setupSubscriberControllerTestFixture()
	request, fault := http.NewRequest("GET", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	query := request.URL.Query()
	index := rand.Intn(len(fixture.expectedRecords)-1) + 1
	query.Add("index", strconv.Itoa(index))
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(fixture.dut.retrieve)
	handler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expectedMessage := ConvertToJson(fixture.expectedRecords[index-1])
	if response.Body.String() != expectedMessage {
		t.Errorf("createHandler returned unexpected body: got %v want %v",
			response.Body.String(), expectedMessage)
	}
	fixture.tearDown()
}

func TestUpdateController(t *testing.T) {
	fixture := setupSubscriberControllerTestFixture()
	request, fault := http.NewRequest("POST", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	form := Subscriber{}
	form.Index = uint8(rand.Intn(len(fixture.expectedRecords)) + 1)
	form.FirstName = "Handsome Marc"
	form.LastName = "Immaculate Conception"

	query := request.URL.Query()
	query.Add("index", strconv.Itoa(int(form.Index)))
	query.Add("first_name", form.FirstName)
	query.Add("last_name", form.LastName)
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(fixture.dut.update)
	handler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	listHandler := http.HandlerFunc(fixture.dut.list)
	listHandler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fixture.expectedRecords[form.Index-1].FirstName = form.FirstName
	fixture.expectedRecords[form.Index-1].LastName = form.LastName
	expectedMessage := ConvertToJson(fixture.expectedRecords)
	if response.Body.String() != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Body.String(), expectedMessage)
	}
	fixture.tearDown()
}

func TestDeleteController(t *testing.T) {
	fixture := setupSubscriberControllerTestFixture()
	index := rand.Intn(len(fixture.expectedRecords)) + 1
	request, fault := http.NewRequest("DELETE", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	query := request.URL.Query()
	query.Add("index", strconv.Itoa(index))
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(fixture.dut.delete)
	handler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	listHandler := http.HandlerFunc(fixture.dut.list)
	listHandler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	updatedExpectedRecords := fixture.expectedRecords[1:]
	if index > 1 {
		updatedExpectedRecords = append(fixture.expectedRecords[:index-1], fixture.expectedRecords[index-1+1:]...)
	}
	expectedMessage := ConvertToJson(updatedExpectedRecords)
	if response.Body.String() != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Body.String(), expectedMessage)
	}
	fixture.tearDown()
}

func TestActivateController(t *testing.T) {
	fixture := setupSubscriberControllerTestFixture()
	index := rand.Intn(len(fixture.expectedRecords)) + 1
	request, fault := http.NewRequest("PATCH", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	query := request.URL.Query()
	query.Add("index", strconv.Itoa(index))
	query.Add("activation_flag", "true")
	request.URL.RawQuery = query.Encode()
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(fixture.dut.activate)
	handler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	listHandler := http.HandlerFunc(fixture.dut.list)
	listHandler.ServeHTTP(response, request)
	if status := response.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fixture.expectedRecords[index-1].ActivationFlag = true
	expectedMessage := ConvertToJson(fixture.expectedRecords)
	if response.Body.String() != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Body.String(), expectedMessage)
	}
	fixture.tearDown()
}

func ConvertToJson(object interface{}) string {
	jsonObject, jsonError := json.Marshal(object)
	if jsonError != nil {
		log.Panic(jsonError)
		return jsonError.Error()
	}
	return string(jsonObject)
}
