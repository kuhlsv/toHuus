// Class for struct handling
// This class just provide public struct for models and work with private helper structs
// With this help, other classes do not need imports for use of structs
package models

import (
	"gopkg.in/mgo.v2/bson"
)

// An UserData represents a user with additional data
type UserData struct{
	Id 	 bson.ObjectId 	`bson:"_id"`
	Username string 	`bson:"Username"`
	Title string 		`bson:"Title"`
	Password string 	`bson:"Password"`
	SessionId string 	`bson:"Session"`
	Avatar string 		`bson:"Avatar"`
}

// An Device represents a device with additional data
type Device struct{
	Id 	 	bson.ObjectId 	`bson:"_id"`
	Name 	string			`bson:"Name"`
	Room 	string			`bson:"Room"`
	Type 	string			`bson:"Type"`
	State 	int				`bson:"State"`
}

// An Event represents a event with additional data
type Event struct{
	Id 	 	bson.ObjectId 	`bson:"_id"`
	UserId 	bson.ObjectId	`bson:"UserId"`
	Name 	string			`bson:"Name"`
	Time 	string			`bson:"Time"`
	Offset  string			`bson:"Offset"`
}

// An Type represents a type with additional data
type Type struct{
	Id 	 	bson.ObjectId 	`bson:"_id"`
	Name 	string			`bson:"Name"`
	Kind 	string			`bson:"Kind"`
	Min		int				`bson:"Min"`
	Max		int				`bson:"Max"`
}

// An RelationEventDevice represents a relation between Event and Devices and the new state
// With this type joins will be build
type RelationEventDevice struct{
	Id 	 		bson.ObjectId 	`bson:"_id"`
	EventId 	bson.ObjectId 	`bson:"EventId"`
	DeviceId 	bson.ObjectId 	`bson:"DeviceId"`
	NewState	int 			`bson:"NewState"`
}

// An AllDTE is a type arrays of all items
type AllDTE struct{
	Devices []Device
	Events	[]Event
	Types   []Type
	Rel     []Item
}

// An Item is a type with data from devices and event actions
type Item struct{
	Id 		bson.ObjectId
	Name 	string
	Room	string
	Value 	int
}

// An Item is a type with arrays of device names and event actions
type Items struct{
	Name 	[]string
	Value 	[]string
}

type SimState struct {
	Id          bson.ObjectId `bson:"_id" json:"Id" xml:"id"`
	CurrentTime int64	      `json:"Time" xml:"time" bson:"time"`
	Sunrise     string        `json:"Sunrise" xml:"sunrise"`
	Sunset      string        `json:"Sunset" xml:"sunset"`
	State	    bool          `json:"State" xml:"state"`
	Multiplier  int           `json:"Multiplier" xml:"multiplier"`
}

// Contains information to export the whole database
type xmlData struct {
	Devices    	[]Device	         `xml:"devices"`
	Types 		[]Type		         `xml:"types"`
	Events     	[]Event			     `xml:"events"`
	Relations   []RelationEventDevice`xml:"relations"`
	Simulator  	[]SimState		     `xml:"simulator"`
	Users      	[]UserData           `xml:"users"`
}