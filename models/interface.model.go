// Class for interface/basic manipulation
// This class processing basic data or calling functions
// Actions from controller check more specific by data and put to correct function
package models

import (
	"os"
	"net/http"
	"io"
	"strings"
	"fmt"
)

// Initialisation
const avatarPath = "./toHuus/conf/avatar/"
const avatarPathHTML = "./avatar/"
const avatarNameAdditive = "avatar_"
const avatarButton = "avatarFile"

// Function to add an new data by calling specific functions
//  Params: Request, dataType(String) -> Item data from form
func AddData(r *http.Request, dataType string){
	switch dataType{
	case "devices":
		// Add a device
		AddDevice(r.FormValue("dName"), r.FormValue("dRoom"), r.FormValue("dType"))
		break
	case "types":
		// Add a type
		AddType(r.FormValue("tName"), r.FormValue("tKind"), r.FormValue("tMin"), r.FormValue("tMax"))
		break
	case "events":
		// Add a event
		deviceItems := getDeviceFormForEvent(r)
		AddEvent(r.FormValue("eName"), r.FormValue("eTime"), r.FormValue("eOffset"), deviceItems, GetUserData(r).Id)
		break
	default:
		UserMessage = "Error: Can not add data"
	}
}

// Function to update data by calling specific functions
//  Params: Request, dataType(String) -> Item data from form
func UpdateData(r *http.Request, dataType string){
	switch dataType{
	case "devices":
		// Update a device
		UpdateDevice(r.FormValue("dName"), r.FormValue("dRoom"), r.FormValue("dType"))
		break
	case "types":
		// Update a type
		UpdateType(r.FormValue("tName"), r.FormValue("tKind"), r.FormValue("tMin"), r.FormValue("tMax"))
		break
	case "events":
		// Update an event
		deviceItems := getDeviceFormForEvent(r)
		UpdateEvent(r.FormValue("eName"), r.FormValue("eTime"), r.FormValue("eOffset"), deviceItems, GetUserData(r).Id)
		break
	default:
		UserMessage = "Error: Can not update data"
	}
}

// Helper function to get all devices with values from a form
//  Params: Request -> Process From data
func getDeviceFormForEvent(r *http.Request) Items{
	deviceItems := Items{}
	// Get all devices out of form for the event
	for k, v := range r.Form {
		if k == "eDevice" {
			deviceItems.Name = v
		}
		if k == "to" {
			deviceItems.Value = v
		}
	}
	return deviceItems
}

// Function to delete data by calling specific functions
//  Params: Request, Name(Strign), dataType(String) -> Item data from form and by name
func DelData(r *http.Request, name string, dataType string){
	switch dataType{
	case "devices":
		// Delete an device
		DelDevice(name)
		break
	case "types":
		// Delete an type
		DelType(name)
		break
	case "events":
		// Delete an event
		DelEvent(name, GetUserData(r).Id)
		break
	default:
		UserMessage = "Error: Can not delete data"
	}
}

// Function to upload an avatar image and write to config files
//  Params: ResponseWriter, Request -> Image by ParseMultipartForm
func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// Get
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile(avatarButton)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	user := GetUserData(r)
	f, err := os.Create(avatarPath + handler.Filename)
	// Close
	if err != nil {
		fmt.Println(err)
		return
	}
	io.Copy(f, file)
	f.Close()
	// Rename
	ending := strings.SplitAfter(handler.Filename, ".")
	newPath := avatarPath + avatarNameAdditive + user.Username + "." + ending[len(ending)-1]
	newPathHTML := avatarPathHTML + avatarNameAdditive + user.Username + "." + ending[len(ending)-1]
	os.Rename(avatarPath + handler.Filename, newPath)
	// Update User
	SetAvatar(r, newPathHTML)
}

// Function to delete an avatar image
//  Params: Path(String) -> Path to file
func DeleteAvatar(path string) {
	f, err := os.Create(path)
	f.Close()
	// Close
	if err != nil {
		fmt.Println(err)
		return
	}else{
		err = os.Remove(path)
		if err != nil{
			UserMessage = "Error: " + err.Error()
			return
		}
	}
}




