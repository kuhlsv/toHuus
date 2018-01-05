package models

import (
	"os"
	"io"
	"fmt"
	"net/http"
	"encoding/xml"
	"strings"
	"io/ioutil"
	"toHuus/db"
)

// Initialisation
const tmpPath = "./toHuus/conf/tmp/"
const importButton = "dataFile"

// Function to import an XML to DB
//  Params: ResponseWriter, Request -> Image by ParseMultipartForm
func Import(w http.ResponseWriter, r *http.Request) {
	// Get
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile(importButton)
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := os.Create(tmpPath + handler.Filename)
	// Close
	if err != nil {
		fmt.Println(err)
		return
	}
	// Write and read (rename)
	io.Copy(f, file)
	f.Close()
	file.Close()
	newPath := tmpPath + "toHuus.xml"
	os.Rename(tmpPath + handler.Filename, newPath)
	bytes, _ := ioutil.ReadFile(newPath)
	// Parse
	data := xmlData{}
	xml.Unmarshal(bytes, &data)
	// DB
	database := db.OpenConnection()
	// Insert
	for _, e := range data.Devices {
		database.C(dbCollDevices).Insert(e)
	}
	database.C(dbCollTypes).DropCollection() // Reset types
	for _, e := range data.Types {
		database.C(dbCollTypes).Insert(e)
	}
	for _, e := range data.Relations {
		database.C(dbCollRelEvents).Insert(e)
	}
	for _, e := range data.Events {
		database.C(dbCollEvents).Insert(e)
	}
	database.C(dbCollSim).DropCollection() // Reset simulator
	for _, e := range data.Simulator {
		database.C(dbCollSim).Insert(e)
	}
	for _, e := range data.Users {
		database.C(DbCollUserdata).Insert(e)
	}
	// Maybe want to remove file
	//os.Remove(newPath)
}

// Function to export an XML from DB
//  Params: ResponseWriter, Request -> Return xml for download
func Export(w http.ResponseWriter, r *http.Request) {
	// Create
	data := xmlData{}
	data.Devices = GetAllDevices()
	data.Types = GetAllTypes()
	data.Events = GetAllEvents("")
	data.Relations = GetAllRelation()
	data.Simulator = GetSimData()
	data.Users = GetAllUserData()
	// Define
	w.Header().Set("Content-Disposition", "attachment; filename=toHuus.xml")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	// Parse and send
	xml, _ := xml.Marshal(data)
	io.Copy(w, strings.NewReader(string(xml)))
}


