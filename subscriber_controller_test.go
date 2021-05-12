/*
* Copyright (c) 2021.
* Marc Concepcion
* marcanthonyconcepcion@gmail.com
 */

package MarcGoRESTAPIDemo

import (
	"encoding/json"
	"github.com/gorilla/mux"
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
	model := MakeDatabaseRecords()
	dut := MakeSubscriberController(model)
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
	subscriberForm := Subscriber{
		EmailAddress: "riseofskywalker@starwars.com",
		FirstName:    "Rey",
		LastName:     "Palpatine",
	}
	request, fault := http.NewRequest("POST", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	query := request.URL.Query()
	query.Add("email_address", subscriberForm.EmailAddress)
	query.Add("first_name", subscriberForm.FirstName)
	query.Add("last_name", subscriberForm.LastName)
	request.URL.RawQuery = query.Encode()
	createResponse := httptest.NewRecorder()
	createHandler := http.HandlerFunc(fixture.dut.create)
	createHandler.ServeHTTP(createResponse, request)
	if status := createResponse.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	recordCreatedMessage := ConvertToJson(Update{"Record created", subscriberForm})
	if createResponse.Body.String() != recordCreatedMessage {
		t.Errorf("createHandler returned unexpected body: got %v want %v", createResponse.Body.String(), recordCreatedMessage)
	}

	newSubscriber := Subscriber{
		Index:          4,
		EmailAddress:   "riseofskywalker@starwars.com",
		FirstName:      "Rey",
		LastName:       "Palpatine",
		ActivationFlag: false,
	}
	updatedExpectedRecords := append(fixture.expectedRecords, newSubscriber)
	expectedMessage := ConvertToJson(updatedExpectedRecords)
	listResponse := httptest.NewRecorder()
	listHandler := http.HandlerFunc(fixture.dut.list)
	listHandler.ServeHTTP(listResponse, request)
	if status := listResponse.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if listResponse.Body.String() != expectedMessage {
		t.Errorf("createHandler returned unexpected body: got %v want %v", listResponse.Body.String(), expectedMessage)
	}
	fixture.tearDown()
}

func TestRetrieveController(t *testing.T) {
	fixture := setupSubscriberControllerTestFixture()
	request, fault := http.NewRequest("GET", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	index := rand.Intn(len(fixture.expectedRecords)) + 1
	vars := map[string]string{
		"index": strconv.Itoa(index),
	}
	request = mux.SetURLVars(request, vars)
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
	vars := map[string]string{
		"index": strconv.Itoa(int(form.Index)),
	}
	request = mux.SetURLVars(request, vars)
	query := request.URL.Query()
	query.Add("first_name", form.FirstName)
	query.Add("last_name", form.LastName)
	request.URL.RawQuery = query.Encode()

	updateResponse := httptest.NewRecorder()
	handler := http.HandlerFunc(fixture.dut.update)
	handler.ServeHTTP(updateResponse, request)
	if status := updateResponse.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	recordUpdatedMessage := ConvertToJson(Update{"Record updated", form})
	if updateResponse.Body.String() != recordUpdatedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			updateResponse.Body.String(), recordUpdatedMessage)
	}

	listResponse := httptest.NewRecorder()
	listHandler := http.HandlerFunc(fixture.dut.list)
	listHandler.ServeHTTP(listResponse, request)
	if status := listResponse.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fixture.expectedRecords[form.Index-1].FirstName = form.FirstName
	fixture.expectedRecords[form.Index-1].LastName = form.LastName
	expectedMessage := ConvertToJson(fixture.expectedRecords)
	if listResponse.Body.String() != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			listResponse.Body.String(), expectedMessage)
	}
	fixture.tearDown()
}

func TestDeleteController(t *testing.T) {
	fixture := setupSubscriberControllerTestFixture()
	request, fault := http.NewRequest("DELETE", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	index := rand.Intn(len(fixture.expectedRecords)) + 1
	vars := map[string]string{
		"index": strconv.Itoa(index),
	}
	request = mux.SetURLVars(request, vars)
	deleteResponse := httptest.NewRecorder()
	handler := http.HandlerFunc(fixture.dut.delete)
	handler.ServeHTTP(deleteResponse, request)
	if status := deleteResponse.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	recordDeletedMessage := ConvertToJson(Message{"success", "Deleted record of subscriber #" + strconv.Itoa(index)})
	if deleteResponse.Body.String() != recordDeletedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			deleteResponse.Body.String(), recordDeletedMessage)
	}

	listResponse := httptest.NewRecorder()
	listHandler := http.HandlerFunc(fixture.dut.list)
	listHandler.ServeHTTP(listResponse, request)
	if status := listResponse.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	updatedExpectedRecords := fixture.expectedRecords[1:]
	if index > 1 {
		updatedExpectedRecords = append(fixture.expectedRecords[:index-1], fixture.expectedRecords[index-1+1:]...)
	}
	expectedMessage := ConvertToJson(updatedExpectedRecords)
	if listResponse.Body.String() != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v", listResponse.Body.String(), expectedMessage)
	}
	fixture.tearDown()
}

func TestActivateController(t *testing.T) {
	fixture := setupSubscriberControllerTestFixture()
	request, fault := http.NewRequest("PATCH", "/subscribers", nil)
	if fault != nil {
		t.Fatal(fault)
	}
	index := rand.Intn(len(fixture.expectedRecords)-1) + 1
	vars := map[string]string{
		"index": strconv.Itoa(index),
	}
	request = mux.SetURLVars(request, vars)
	query := request.URL.Query()
	query.Add("activation_flag", "true")
	request.URL.RawQuery = query.Encode()
	activateResponse := httptest.NewRecorder()
	handler := http.HandlerFunc(fixture.dut.activate)
	handler.ServeHTTP(activateResponse, request)
	if status := activateResponse.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expectedResponse := ConvertToJson(Message{"success", "Record #" + strconv.Itoa(index) + " activated."})
	if activateResponse.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", activateResponse.Body.String(), expectedResponse)
	}

	listHandler := http.HandlerFunc(fixture.dut.list)
	listResponse := httptest.NewRecorder()
	listHandler.ServeHTTP(listResponse, request)
	if status := listResponse.Code; status != http.StatusOK {
		t.Errorf("createHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	fixture.expectedRecords[index-1].ActivationFlag = true
	expectedMessage := ConvertToJson(fixture.expectedRecords)
	if listResponse.Body.String() != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v", listResponse.Body.String(), expectedMessage)
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
