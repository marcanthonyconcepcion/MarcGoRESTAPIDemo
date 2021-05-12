/*
 * Copyright (c) 2021.
 * Marc Concepcion
 * marcanthonyconcepcion@gmail.com
 */

package MarcGoRESTAPIDemo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

type SubscriberController struct {
	model Records
}

type Message struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

type Update struct {
	Message string     `json:"message"`
	Updates Subscriber `json:"updates"`
}

func (controller SubscriberController) list(response http.ResponseWriter, request *http.Request) {
	subscribers, recordsError := controller.model.list()
	if recordsError != nil {
		controller.sendErrorMessage(http.StatusInternalServerError, response, recordsError.Error())
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
	if 0 == len(request.URL.Query()) {
		controller.sendErrorMessage(http.StatusMethodNotAllowed, response,
			"HTTP command POST without providing parameters is not allowed. Please provide an acceptable HTTP command.")
		return
	}
	emailAddress := request.URL.Query().Get("email_address")
	firstName := request.URL.Query().Get("first_name")
	lastName := request.URL.Query().Get("last_name")
	subscriber := Subscriber{}
	subscriber.LastName = lastName
	subscriber.FirstName = firstName
	subscriber.EmailAddress = emailAddress
	_, recordsError := controller.model.create(subscriber)
	if recordsError != nil {
		controller.sendErrorMessage(http.StatusInternalServerError, response, recordsError.Error())
		return
	}

	jsonSubscriber, jsonError := json.Marshal(Update{"Record created", subscriber})
	if jsonError != nil {
		log.Panic(jsonError)
		return
	}
	_, ioError := io.WriteString(response, string(jsonSubscriber))
	if ioError != nil {
		log.Panic(ioError)
	}
}

func (controller SubscriberController) retrieve(response http.ResponseWriter, request *http.Request) {
	index, indexError := strconv.Atoi(mux.Vars(request)["index"])
	if indexError != nil {
		controller.sendErrorMessage(http.StatusBadRequest, response, indexError.Error())
		return
	}
	subscriber, recordsError := controller.model.retrieve(uint8(index))
	if recordsError != nil {
		if errors.Is(recordsError, sql.ErrNoRows) {
			controller.sendErrorMessage(http.StatusNotFound, response, "Subscriber does not exist.")
		} else {
			controller.sendErrorMessage(http.StatusInternalServerError, response, recordsError.Error())
		}
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
		controller.sendErrorMessage(http.StatusBadRequest, response, indexError.Error())
		return
	}
	subscriber.Index = uint8(index)
	if 0 == len(request.URL.Query()) {
		controller.sendErrorMessage(http.StatusMethodNotAllowed, response,
			"HTTP command PUT without providing parameters is not allowed. Please provide an acceptable HTTP command.")
		return
	}
	emailAddress := request.URL.Query().Get("email_address")
	subscriber.EmailAddress = emailAddress
	firstName := request.URL.Query().Get("first_name")
	subscriber.FirstName = firstName
	lastName := request.URL.Query().Get("last_name")
	subscriber.LastName = lastName
	_, recordsError := controller.model.update(subscriber)
	if recordsError != nil {
		controller.sendErrorMessage(http.StatusInternalServerError, response, recordsError.Error())
		return
	}

	jsonSubscriber, jsonError := json.Marshal(Update{"Record updated", subscriber})
	if jsonError != nil {
		log.Panic(jsonError)
		return
	}
	_, ioError := io.WriteString(response, string(jsonSubscriber))
	if ioError != nil {
		log.Panic(ioError)
	}
}

func (controller SubscriberController) delete(response http.ResponseWriter, request *http.Request) {
	index, indexError := strconv.Atoi(mux.Vars(request)["index"])
	if indexError != nil {
		controller.sendErrorMessage(http.StatusBadRequest, response, indexError.Error())
		return
	}
	_, recordsError := controller.model.delete(uint8(index))
	if recordsError != nil {
		controller.sendErrorMessage(http.StatusInternalServerError, response, recordsError.Error())
		return
	}

	jsonSubscriber, jsonError := json.Marshal(Message{"success", "Deleted record of subscriber #" + strconv.Itoa(index)})
	if jsonError != nil {
		log.Panic(jsonError)
		return
	}
	_, ioError := io.WriteString(response, string(jsonSubscriber))
	if ioError != nil {
		log.Panic(ioError)
	}
}

func (controller SubscriberController) activate(response http.ResponseWriter, request *http.Request) {
	index, indexError := strconv.Atoi(mux.Vars(request)["index"])
	if indexError != nil {
		controller.sendErrorMessage(http.StatusBadRequest, response, indexError.Error())
		return
	}
	activateString := request.URL.Query().Get("activation_flag")
	activate := false
	if activateString == "true" {
		activate = true
	} else {
		controller.sendErrorMessage(http.StatusBadRequest, response,
			"Only activating a subscriber is allowed. Please set the activation_flag to 'true'.")
		return
	}
	_, recordsError := controller.model.activate(uint8(index), activate)
	if recordsError != nil {
		controller.sendErrorMessage(http.StatusInternalServerError, response, recordsError.Error())
		return
	}

	jsonSubscriber, jsonError := json.Marshal(Message{"success", "Record #" + strconv.Itoa(index) + " activated."})
	if jsonError != nil {
		log.Panic(jsonError)
		return
	}
	_, ioError := io.WriteString(response, string(jsonSubscriber))
	if ioError != nil {
		log.Panic(ioError)
	}
}

func (controller SubscriberController) sendErrorMessage(httpStatusCode int, response http.ResponseWriter, errorMessage string) {
	jsonErrorMessage, jsonError := json.Marshal(Message{"error", errorMessage})
	if jsonError != nil {
		log.Panic(jsonError)
	}
	http.Error(response, string(jsonErrorMessage), httpStatusCode)
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
