/*
 * Copyright (c) 2021.
 * Marc Concepcion
 * marcanthonyconcepcion@gmail.com
 */

package MarcGoRESTAPIDemo

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

type SubscriberController struct {
	model Records
}

type ErrorMessage struct {
	Error string
}

func (controller SubscriberController) list(response http.ResponseWriter, request *http.Request) {
	subscribers, recordsError := controller.model.list()
	if recordsError != nil {
		controller.sendErrorMessage(response, recordsError.Error())
		return
	}
	jsonSubscribers, jsonError := json.Marshal(subscribers)
	if jsonError != nil {
		log.Panic(jsonError)
		return
	}
	_, ioError := io.WriteString(response, string(jsonSubscribers))
	if ioError != nil {
		log.Panic(ioError)
	}
}

func (controller SubscriberController) create(response http.ResponseWriter, request *http.Request) {
	emailAddress := request.URL.Query().Get("email_address")
	firstName := request.URL.Query().Get("first_name")
	lastName := request.URL.Query().Get("last_name")
	subscriber := Subscriber{}
	subscriber.LastName = lastName
	subscriber.FirstName = firstName
	subscriber.EmailAddress = emailAddress
	_, recordsError := controller.model.create(subscriber)
	if recordsError != nil {
		controller.sendErrorMessage(response, recordsError.Error())
		return
	}
}

func (controller SubscriberController) retrieve(response http.ResponseWriter, request *http.Request) {
	index, indexError := strconv.Atoi(mux.Vars(request)["index"])
	if indexError != nil {
		controller.sendErrorMessage(response, indexError.Error())
	}
	subscriber, recordsError := controller.model.retrieve(uint8(index))
	if recordsError != nil {
		controller.sendErrorMessage(response, recordsError.Error())
		return
	}
	jsonSubscriber, jsonError := json.Marshal(subscriber)
	if jsonError != nil {
		log.Panic(jsonError)
		return
	}
	_, ioError := io.WriteString(response, string(jsonSubscriber))
	if ioError != nil {
		log.Panic(ioError)
	}
}

func (controller SubscriberController) update(response http.ResponseWriter, request *http.Request) {
	subscriber := Subscriber{}
	index, indexError := strconv.Atoi(mux.Vars(request)["index"])
	if indexError != nil {
		controller.sendErrorMessage(response, indexError.Error())
	}
	subscriber.Index = uint8(index)
	emailAddress := request.URL.Query().Get("email_address")
	subscriber.EmailAddress = emailAddress
	firstName := request.URL.Query().Get("first_name")
	subscriber.FirstName = firstName
	lastName := request.URL.Query().Get("last_name")
	subscriber.LastName = lastName
	_, recordsError := controller.model.update(subscriber)
	if recordsError != nil {
		controller.sendErrorMessage(response, recordsError.Error())
		return
	}
}

func (controller SubscriberController) delete(response http.ResponseWriter, request *http.Request) {
	index, indexError := strconv.Atoi(mux.Vars(request)["index"])
	if indexError != nil {
		controller.sendErrorMessage(response, indexError.Error())
	}
	_, recordsError := controller.model.delete(uint8(index))
	if recordsError != nil {
		controller.sendErrorMessage(response, recordsError.Error())
		return
	}
}

func (controller SubscriberController) activate(response http.ResponseWriter, request *http.Request) {
	index, indexError := strconv.Atoi(mux.Vars(request)["index"])
	if indexError != nil {
		controller.sendErrorMessage(response, indexError.Error())
	}
	activateString := request.URL.Query().Get("activation_flag")
	activate := false
	if activateString == "true" {
		activate = true
	} else if activateString == "false" {
		activate = false
	} else {
		controller.sendErrorMessage(response, "Invalid activation flag. "+
			"Please provide either 'true' to activate the subscriber or 'false' to deactivate the subscriber.")
		return
	}
	_, recordsError := controller.model.activate(uint8(index), activate)
	if recordsError != nil {
		controller.sendErrorMessage(response, recordsError.Error())
		return
	}
}

func (controller SubscriberController) sendErrorMessage(response http.ResponseWriter, error string) {
	jsonErrorMessage, jsonError := json.Marshal(ErrorMessage{error})
	if jsonError != nil {
		log.Panic(jsonError)
	}
	_, ioError := io.WriteString(response, string(jsonErrorMessage))
	if ioError != nil {
		log.Panic(ioError)
	}
}

func MakeSubscriberController(model Records) SubscriberController {
	return SubscriberController{model}
}

func (controller SubscriberController) ViewHandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/subscribers", controller.list).Methods("GET")
	router.HandleFunc("/subscribers", controller.create).Methods("POST")
	router.HandleFunc("/subscribers/{index}", controller.update).Methods("PUT")
	router.HandleFunc("/subscribers/{index}", controller.activate).Methods("PATCH")
	router.HandleFunc("/subscribers/{index}", controller.delete).Methods("DELETE")
	router.HandleFunc("/subscribers/{index}", controller.retrieve).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
