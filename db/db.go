// Package for a database connection
// Call the Database function at start once
// Call OpenConnection to get the database pointer
package db

import (
	"gopkg.in/mgo.v2"
	"fmt"
)

// Initialisation
var url = ""
var name = ""
var db *mgo.Database = nil

// Basic function for this package to initialise the connection
// Call once at program start
//  Params: String, String -> DB Name and DB URL
func Database(dbName string, dbUrl string){
	setDB(dbName)
	setUrl(dbUrl)
}

// Function to open an connection from the db
//  Return: *Database -> Pointer of DB connection
func OpenConnection() *mgo.Database{
	// Open
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println(err)
	}
	db = session.DB(name)
	//defer
	return db
}

// Function to close an connection from the db
func CloseConnection(){
	db.Session.Close()
}

// Function to get the connection
//  Return: *Database -> Pointer of DB connection
func GetConnection() *mgo.Database{
	return db
}

// Helper-Function to set name
//  Params: String -> DB Name
func setDB(dbName string){
	name = dbName
}

// Helper-Function to set url
//  Params: String -> DB  url
func setUrl(dbUrl string){
	url = dbUrl
}