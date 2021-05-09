/*
 * Copyright (c) 2021.
 * Marc Concepcion
 * marcanthonyconcepcion@gmail.com
 */

package MarcGoRESTAPIDemo

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

var settings = readConfiguration("resources/MarcGoRESTAPIDemo.yaml")

type Records struct {
	database *sql.DB
}

type Subscriber struct {
	Index          uint8  `json:"index,omitempty"`
	EmailAddress   string `json:"email_address,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	ActivationFlag bool   `json:"activation_flag,omitempty"`
}

func makeDatabaseRecords() Records {
	database, dbInstanceFail := sql.Open("mysql", settings.Database.User+":"+settings.Database.Password+
		"@tcp("+settings.Database.Host+":"+strconv.Itoa(int(settings.Database.Port))+")/"+settings.Database.DBName)
	if dbInstanceFail != nil {
		panic(dbInstanceFail.Error())
	}
	return Records{database}
}

func (records Records) create(subscriber Subscriber) (sql.Result, error) {
	result, fault := records.database.Exec(
		"insert into `subscribers` (`email_address`, `last_name`, `first_name`) values (?, ?, ?)",
		subscriber.EmailAddress, subscriber.LastName, subscriber.FirstName)
	return result, fault
}

func (records Records) retrieve(index uint8) (*Subscriber, error) {
	var subscriber Subscriber
	record := records.database.QueryRow("select * from `subscribers` where `index`=?", index)
	if recordModelError := record.Scan(&subscriber.Index, &subscriber.EmailAddress, &subscriber.LastName, &subscriber.FirstName,
		&subscriber.ActivationFlag); recordModelError != nil {
		return &subscriber, recordModelError
	}
	return &subscriber, nil
}

func (records Records) update(subscriber Subscriber) (sql.Result, error) {
	var parametersToUpdate []string
	if "" != subscriber.EmailAddress {
		parametersToUpdate = append(parametersToUpdate, "`email_address` = "+"\""+subscriber.EmailAddress+"\"")
	}
	if "" != subscriber.LastName {
		parametersToUpdate = append(parametersToUpdate, "`last_name` = "+"\""+subscriber.LastName+"\"")
	}
	if "" != subscriber.FirstName {
		parametersToUpdate = append(parametersToUpdate, "`first_name` = "+"\""+subscriber.FirstName+"\"")
	}
	result, updateFail := records.database.Exec("update `subscribers` set "+
		strings.Join(parametersToUpdate, ",")+" where `index`=?", subscriber.Index)
	return result, updateFail
}

func (records Records) activate(index uint8, activate bool) (sql.Result, error) {
	activationFlag := 0
	if activate == true {
		activationFlag = 1
	}
	result, updateFail := records.database.Exec(
		"update `subscribers` set activation_flag=? where `index`=?", activationFlag, index)
	return result, updateFail
}

func (records Records) delete(index uint8) (sql.Result, error) {
	result, deleteError := records.database.Exec("delete from `subscribers` where `index`=?", index)
	return result, deleteError
}

func (records Records) list() ([]Subscriber, error) {
	rows, dbQueryError := records.database.Query("select * from `subscribers`")
	if dbQueryError != nil {
		return nil, dbQueryError
	}
	subscribers := make([]Subscriber, 0)
	for rows.Next() {
		var subscriber Subscriber
		if recordModelError := rows.Scan(&subscriber.Index, &subscriber.EmailAddress, &subscriber.LastName, &subscriber.FirstName,
			&subscriber.ActivationFlag); recordModelError != nil {
			return subscribers, recordModelError
		}
		subscribers = append(subscribers, subscriber)
	}
	rowsCloseError := rows.Close()
	if rowsCloseError != nil {
		return subscribers, rowsCloseError
	}
	return subscribers, nil
}
