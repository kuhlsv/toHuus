// Class for data manipulation
// This class processing data of items (Add/Edit/Delete/Get)
// Package db is needed
package models

import (
	"gopkg.in/mgo.v2/bson"
	"toHuus/db"
	"strconv"
	"strings"
)

// Initialisation
const dbCollDevices = "Devices"
const dbCollEvents = "Events"
const dbCollTypes = "Types"
const dbCollRelEvents = "RelationEventsDevices"
const dbCollSim = "Simulator"

// Function to add an new device to db
//  Params: Name(String), Room(String), dataType(String) -> Device data
func AddDevice(name string, room string, dataType string){
	database := db.OpenConnection()
	coll := database.C(dbCollDevices)
	device := Device{}
	coll.Find(bson.M{"Name" : name}).One(&device)
	// Check not present
	if len(device.Name) < 1 && len(name) > 0 && len(room) > 0 {
		// Write to db
		device = Device{
			bson.NewObjectId(),
			name,
			room,
			dataType,
			0,
		}
		coll.Insert(device)
	} else {
		UserMessage = "Error: Device already exist"
	}
	db.CloseConnection()
}

// Function to add an new event to db for an specific user
//  Params: Name(String), Time(String), dataType(Items), userid(ObjectId) -> Event data
func AddEvent(name string, time string, offset string, deviceItems Items, userid bson.ObjectId){
	database := db.OpenConnection()
	coll := database.C(dbCollEvents)
	coll2 := database.C(dbCollRelEvents)
	event := Event{}
	eventid := bson.NewObjectId()
	coll.Find(bson.M{"Name" : name}).One(&event)
	// Check not present
	if len(event.Name) < 1 && len(name) > 0 {
		// Write Devices/Relation
		for n := range deviceItems.Name{
			state, _ := strconv.Atoi(deviceItems.Value[n])
			name := strings.Split(deviceItems.Name[n], " | ")[0]
			rel := RelationEventDevice{
				bson.NewObjectId(),
				eventid,
				GetDeviceByName(name).Id,
				state,
			}
			coll2.Insert(rel)
		}
		// Write to db
		event = Event{
			eventid,
			userid,
			name,
			time,
			offset,
		}
		coll.Insert(event)
	} else {
		UserMessage = "Error: Event already exist"
	}
	db.CloseConnection()
}

// Function to add an new tye to db
//  Params: Name(String), Kind(String), Min(String), Max(String) -> Type data
func AddType(name string, kind string, min string, max string){
	database := db.OpenConnection()
	coll := database.C(dbCollTypes)
	// Convert to int
	iMin,err := strconv.Atoi(min)
	iMax,err2 := strconv.Atoi(max)
	if err2 != nil && err != nil {
		// Set switch values
		iMin = 0
		iMax = 1
	}
	// Check not present
	ntype := Type{}
	coll.Find(bson.M{"Name" : name}).One(&ntype)
	if len(ntype.Name) < 1 && len(name) > 0 && iMin < iMax {
		// Write to db
		ntype = Type{
			bson.NewObjectId(),
			name,
			kind,
			iMin,
			iMax,
		}
		coll.Insert(ntype)
	} else {
		UserMessage = "Error: Type already exist"
	}
	db.CloseConnection()
}

// Function to update an device to db
//  Params: Name(String), Room(String), dataType(String) -> Device data
func UpdateDevice(name string, room string, dataType string){
	database := db.OpenConnection()
	coll := database.C(dbCollDevices)
	device := Device{}
	coll.Find(bson.M{"Name" : name}).One(&device)
	// Write to db
	device2 := Device{
		device.Id,
		device.Name,
		room,
		dataType,
		device.State,
	}
	coll.Update(device, device2)
	db.CloseConnection()
}

// Function to update an device to db
//  Params: Name(String), Room(String), dataType(String) -> Device data
func UpdateState(name string, state string){
	database := db.OpenConnection()
	coll := database.C(dbCollDevices)
	device := Device{}
	coll.Find(bson.M{"Name" : name}).One(&device)
	newState, err := strconv.Atoi(state)
	if err == nil {
		// Write to db
		coll.Update(bson.M{"_id" : device.Id}, bson.M{"$set" : bson.M{"State" : newState}})
	}
	db.CloseConnection()
}

// Function to update an event to db for an specific user
//  Params: Name(String), Time(String), dataType(String), userid(ObjectId) -> Event data
func UpdateEvent(name string, time string, offset string, deviceItems Items, userid bson.ObjectId){
	database := db.OpenConnection()
	coll := database.C(dbCollEvents)
	coll2 := database.C(dbCollRelEvents)
	event := Event{}
	coll.Find(bson.M{"Name" : name}).One(&event)
	// Devices/Relation
	for n := range deviceItems.Name{
		state, _ := strconv.Atoi(deviceItems.Value[n])
		name := strings.Split(deviceItems.Name[n], " | ")[0]
		rel := RelationEventDevice{}
		coll2.Find(bson.M{"EventId" : event.Id, "DeviceId" : GetDeviceByName(name).Id}).One(&rel)
		// Update relation from edit
		if rel.Id != "" {
			rel2 := RelationEventDevice{
				rel.Id,
				event.Id,
				GetDeviceByName(name).Id,
				state,
			}
			coll2.Update(rel, rel2)
		}else{
			// Add new relation from edit
			rel2 := RelationEventDevice{
				bson.NewObjectId(),
				event.Id,
				GetDeviceByName(name).Id,
				state,
			}
			coll2.Insert(rel2)
		}
	}
	// Write to db
	event2 := Event{
		event.Id,
		userid,
		name,
		time,
		offset,
	}
	coll.Update(event, event2)
	db.CloseConnection()
}

// Function to update an tye to db
//  Params: Name(String), Kind(String), Min(String), Max(String) -> Type data
func UpdateType(name string, kind string, min string, max string){
	database := db.OpenConnection()
	coll := database.C(dbCollTypes)
	// Check not present
	ntype := Type{}
	coll.Find(bson.M{"Name" : name}).One(&ntype)
	iMin,err := strconv.Atoi(min)
	iMax,err2 := strconv.Atoi(max)
	if err2 != nil && err != nil {
		// Set switch values
		iMin = ntype.Min
		iMax = ntype.Max
	}
	// Write to db
	ntype2 := Type{
		ntype.Id,
		name,
		kind,
		iMin,
		iMax,
	}
	coll.Update(ntype, ntype2)
	db.CloseConnection()
}

// Function to delete an device from db
//  Params: Name(String) -> Device name to find item
func DelDevice(name string){
	database := db.OpenConnection()
	coll := database.C(dbCollDevices)
	coll.Remove(bson.M{"Name" : name})
	db.CloseConnection()
}

// Function to delete an event from db
//  Params: Name(String), userId(ObjectId) -> Event name to find item of an specific user
func DelEvent(name string, userid bson.ObjectId){
	database := db.OpenConnection()
	// Device/Relation
	coll := database.C(dbCollRelEvents)
	coll.RemoveAll(bson.M{"EventId" : GetEventByName(name, userid).Id})
	// Events
	coll2 := database.C(dbCollEvents)
	coll2.Remove(bson.M{"Name" : name})
	db.CloseConnection()
}

// Function to delete an type from db
//  Params: Name(String) -> Type name to find item
func DelType(name string){
	database := db.OpenConnection()
	coll := database.C(dbCollTypes)
	coll.Remove(bson.M{"Name" : name})
	db.CloseConnection()
}

// Function to delete multiple events from db
//  Params: userId(ObjectId) -> Delete all events by user
func DelEventsById(userid bson.ObjectId){
	database := db.OpenConnection()
	coll := database.C(dbCollRelEvents)
	coll.RemoveAll(bson.M{"UserId" : userid})
	db.CloseConnection()
}

// Function to get every items from all data
//  Params: userId(ObjectId) -> Delete all events by user
//  Return: AllDTE(type from model) -> Struct with all items
func GetAllDTE(userid bson.ObjectId) AllDTE{
	result := AllDTE{
		GetAllDevices(),
		GetAllEvents(userid),
		GetAllTypes(),
		GetRelationByUser(userid),
	}
	return result
}

// Function to get all devices
//  Return: Devices(type from model) -> Struct with arrays of Device(type)
func GetAllDevices() []Device{
	result := []Device{}
	database := db.OpenConnection()
	coll := database.C(dbCollDevices)
	coll.Find(nil).All(&result)
	db.CloseConnection()
	return result
}

// Function to get all events by user
//  Params: userId(ObjectId) -> Events limited by active user
//  Return: Events(type from model) -> Struct with arrays of Event(type)
func GetAllEvents(userid bson.ObjectId) []Event{
	result := []Event{}
	database := db.OpenConnection()
	coll := database.C(dbCollEvents)
	if userid == ""{
		coll.Find(nil).All(&result)
	}else{
		coll.Find(bson.M{"UserId" : userid}).All(&result)
	}
	db.CloseConnection()
	return result
}

// Function to get all types
//  Return: Types(type from model) -> Struct with arrays of Type(type)
func GetAllTypes() []Type{
	result := []Type{}
	database := db.OpenConnection()
	coll := database.C(dbCollTypes)
	coll.Find(nil).All(&result)
	db.CloseConnection()
	return result
}

// Function to get all devices by type
//  Params: Type(String) -> Kind of type to find
//  Return: Devices(type from model) -> Struct with arrays of Device(type)
func GetAllDevicesByType(ntype string) []Device{
	result := []Device{}
	database := db.OpenConnection()
	coll := database.C(dbCollDevices)
	coll.Find(bson.M{"Type" : ntype}).All(&result)
	db.CloseConnection()
	return result
}

// Function to get all relations
//  Return: RelationEventDevice(type from model) -> Struct all relation data
func GetAllRelation() []RelationEventDevice{
	relation := []RelationEventDevice{}
	database := db.OpenConnection()
	coll := database.C(dbCollRelEvents)
	coll.Find(nil).All(&relation)
	db.CloseConnection()
	return relation
}

// Function to get all relations to an event
//  Params: EventId(String) -> Event to query
//  Return: Devices(type from model) -> Struct with arrays of Device(type)
func GetRelationToEvent(id string) []Device{
	relation := []RelationEventDevice{}
	resultD := []Device{}
	database := db.OpenConnection()
	coll := database.C(dbCollRelEvents)
	coll.Find(bson.M{"EventId" : id}).All(&relation)
	coll = database.C(dbCollDevices)
	// Check all devices to the id took from relation
	for i := 0; i < len(relation); i++ {
		coll.Find(bson.M{"_id" : relation[i].DeviceId}).One(&resultD[i])
	}
	db.CloseConnection()
	return resultD
}

// Function to get new state for event
//  Params: EventId(String), DeviceId(String) -> Event to query
//  Return: int -> State
func GetNewState(evt string, dev string) int{
	relation := RelationEventDevice{}
	database := db.OpenConnection()
	coll := database.C(dbCollRelEvents)
	coll.Find(bson.M{"EventId" : evt , "DeviceId" : dev}).One(&relation)
	return relation.NewState
}

// Function to get all relations to an user
//  Params: UserId(ObjectId) -> User to query
//  Return: Item(type from model) -> Struct with relation and device data
func GetRelationByUser(id bson.ObjectId) []Item{
	relation := []Item{}
	events := GetAllEvents(id)
	database := db.OpenConnection()
	coll := database.C(dbCollRelEvents)
	// Loop through all events from the user
	for i := 0; i < len(events); i++ {
		buffer := []RelationEventDevice{}
		coll.Find(bson.M{"EventId" : events[i].Id}).All(&buffer)
		// Loop through all relations to get the correct relation in an array
		for j := 0; j < len(buffer); j++ {
			dev := GetDeviceById(buffer[j].DeviceId)
			relation = append(relation,
				Item{ events[i].Id, dev.Name, dev.Room, buffer[j].NewState })
		}
	}
	db.CloseConnection()
	return relation
}

// Function to get a device by name
//  Params: Name(String) -> Device to query
//  Return: Device(type from model) -> Struct witch device data
func GetDeviceByName(name string) Device{
	result := Device{}
	database := db.OpenConnection()
	coll := database.C(dbCollDevices)
	coll.Find(bson.M{"Name" : name}).One(&result)
	db.CloseConnection()
	return result
}

// Function to get a event by name
//  Params: Name(String), UserId(ObjectId) -> Event to query for a user
//  Return: Event(type from model) -> Struct witch event data
func GetEventByName(name string, userid bson.ObjectId) Event{
	result := Event{}
	database := db.OpenConnection()
	coll := database.C(dbCollEvents)
	coll.Find(bson.M{"Name" : name, "UserId" : userid}).One(&result)
	db.CloseConnection()
	return result
}

// Function to get a device by id
//  Params: Id(String) -> Device to query
//  Return: Device(type from model) -> Struct witch device data
func GetDeviceById(id bson.ObjectId) Device{
	result := Device{}
	database := db.OpenConnection()
	coll := database.C(dbCollDevices)
	coll.Find(bson.M{"_id" : id}).One(&result)
	db.CloseConnection()
	return result
}

// Function to get the kind of an type
//  Params: ntype(String) -> Type name
//  Return: String -> Kind
func GetKindByType(ntype string) string{
	result := Type{}
	database := db.OpenConnection()
	coll := database.C(dbCollTypes)
	coll.Find(bson.M{"Name" : ntype}).One(&result)
	db.CloseConnection()
	return result.Kind
}

// Function to get a event by id
//  Params: Id(String) -> Event to query
//  Return: Event(type from model) -> Struct witch event data
func GetEventById(id string, userid bson.ObjectId) Event{
	result := Event{}
	database := db.OpenConnection()
	coll := database.C(dbCollEvents)
	coll.Find(bson.M{"_id" : id, "UserId" : userid}).One(&result)
	db.CloseConnection()
	return result
}

// Function to get the current sim states as array
//  Return: SimState(type from model) -> Struct witch simulator data
func GetSimData() []SimState{
	result := []SimState{}
	database := db.OpenConnection()
	coll := database.C(dbCollSim)
	coll.Find(nil).All(&result)
	db.CloseConnection()
	return result
}

// Function to get the current sim states
//  Return: SimState(type from model) -> Struct witch simulator data
func SetSimData(states SimState) {
	result := []SimState{}
	database := db.OpenConnection()
	coll := database.C(dbCollSim)
	coll.Find(nil).All(&result)
	if len(result)>0{
		if result[0].Id == "" {
			coll.Insert(states)
		}else{
			states.Id = result[0].Id
			coll.Update(result[0], states)
		}
	}else{
		coll.Insert(states)
	}
	db.CloseConnection()
}

