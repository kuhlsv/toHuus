"use strict";
///////////////////////////////////////////////////////////////////////////////////////
////// All
///////////////////////////////////////////////////////////////////////////////////////

// Init
window.addEventListener("load", function(){ init() });

// Helper
let ajaxIntervalOverview = 3000;
let ajaxIntervalState = 1000;
let address = "http://localhost:4242";

function $id(id){
    return document.getElementById(id);
}
function $tag(tag){
    return document.getElementsByTagName(tag);
}
function $name(name){
    return document.getElementsByName(name);
}

// Start function
function init() {
    // Serialization
    serialization();
    // Set timeout for overview if ui is open
    if ($id("home")) {
        // Get device state and write at start/interval
        getDevices("once");
        getDevices();
    }
    // Get states for simulator if sim is open
    if ($id("simUi")) {
        // Get sim states and write at start/interval
        getSimStates("once");
        getSimStates();
    }
    // EventHandler
    loadEventListeners();
    // Show Message
    showMessage();
}

// Function to initialise event listeners
function loadEventListeners() {
    // Add event listener just on interface site
    // if items exist
    if ($id("addMenu")) {
        dialogHandlers();
        actionButtonHandlers();
    }
    // Add event listener to simulater ui
    // if items exist
    if ($id("simUi")) {
        actionButtonHandlersSim();
    }
    // Toggle menu
    $id("headerToggle").addEventListener("click", ()=> showNav())
}

// Function to decode names and ids
// Params: text -> text to be converted
// Return: converted text
function decodeHTMLEntities(text) {
    let entities = [
        ['amp', '&'], ['apos', '\''], ['#x27', '\''],
        ['#x2F', '/'], ['#39', '\''], ['#47', '/'],
        ['lt', '<'], ['gt', '>'], ['nbsp', ' '], ['quot', '"']
    ];
    for (let i = 0, max = entities.length; i < max; ++i)
        text = text.replace(new RegExp('&'+entities[i][0]+';', 'g'), entities[i][1]);
    return text;
}

///////////////////////////////////////////////////////////////////////////////////////
///// Interface
///////////////////////////////////////////////////////////////////////////////////////

// Function to initialise the listeners for the dialog windows
function dialogHandlers() {
    // helper function
    function action(bId, iName, dId, show, read) {
        $name(iName)[0].readOnly = false; // Input
        $id(bId).value = "Add"; // Button
        // Dialog
        $id(dId).style.display = show ? "block" : "none";
        updateDialogTitle(dId, true);
    }
    // to show the dialogs
    $id("addDeviceBtn").addEventListener("click", ()=>{
        action("addD","dName","addDevice", true);
    });
    $id("addTypeBtn").addEventListener("click", ()=>{
        action("addT","tName","addType", true);
    });
    $id("addEventBtn").addEventListener("click", ()=>{
        action("addE","eName","addEvent", true);
    });
    // to hide the dialogs
    $id("cancelD").addEventListener("click", ()=>{
        action("addD","dName","addDevice", false);
    });
    $id("cancelT").addEventListener("click", ()=>{
        action("addT","tName","addType", false);
    });
    $id("cancelE").addEventListener("click", ()=>{
        action("addE","eName","addEvent", false);
    });
    // Listener to add/cancel Item to Event Dialog
    $id("eDevices-btn").addEventListener("click", ()=> {
        addItemEvent();
    });
    $id("cancelE").addEventListener("click", ()=> {
        cancelItemEvent();
    });
    // Listener to update type selection
    $id("tKind").addEventListener("change", ()=>{
        changeMinMax();
    });

    // Check all items set from user to enable add button
    // Just if all inputs are set, something can be added
    let buttons = ["addE", "addEvent", "addD", "addDevice", "addT", "addType"];
    for(let i = 1; i < buttons.length; i=i+2) {
        let elements = $id(buttons[i]).getElementsByTagName("input");
        // Check changes input fields
        for(let e = 0; e < elements.length; e++) {
            elements[e].addEventListener("change", () => {
                let value = true;
                for (let e2 = 0; e2 < elements.length; e2++) {
                    if (elements[e2].value == "") {
                        value = false;
                    }
                }
                $id(buttons[i-1]).disabled = !value;
            });
        }
        // Check at hover dialog
        $id(buttons[i]).addEventListener("mouseover", ()=>{
            for(let e = 0; e < elements.length; e++) {
                let value = true;
                for (let e2 = 0; e2 < elements.length; e2++) {
                    if (elements[e2].value == "") {
                        value = false;
                    }
                }
                $id(buttons[i-1]).disabled = !value;
            }
        });
    }
}

// Function optimize the template
function serialization() {
    let body = $tag('body')[0];
    // Disable transitions until the page has loaded
    body.className += 'is-loading';
    window.addEventListener('load', function () {
        body.className -= 'is-loading';
    });
}

// If there is a error message, it will be shown for short time
function showMessage(){
    let text = $id("message").childNodes[1].textContent;
    if(text.substr(0,5) !== "Error"){
        $id("message").style.backgroundColor = "darkgreen";
    }
    if(text !== "") {
        $id("message").style.display = "block";
        setTimeout( ()=>{
            $id("message").style.display = "none";
            text = ""
        }, 3000)
    }
}

// Funktion to add an Item to Event Dialog for multiple devices
function addItemEvent() {
    let pattern = $id("eDevices1").cloneNode(true);
    $id("eDevices-table").getElementsByTagName('tbody')[0].appendChild(pattern);
}

// Funktion to remove Items from Event Dialog
function cancelItemEvent() {
    let pattern = $id("eDevices1").cloneNode(true);
    let elements = $name("eDevices-tr");
    let table = $id("eDevices-table").getElementsByTagName('tbody')[0];
    for(let i = 0; i < elements.length; i++) {
        table.removeChild(elements[i]);
    }
    table.appendChild(pattern);
}

// Funktion to change Min Max option to Type selection in dialog
function changeMinMax() {
    let check = false;
    let parent = $id("tDevices-table");
    // Disable for switch (Just on/off)
    if($id("tKind").value != "Switch"){
        if(parent.lastChild.textContent != "Max: "){
            // Create Min/Max inputs
            check = true;
            let trMin = document.createElement("tr");
            let trMax = document.createElement("tr");
            trMin.appendChild(document.createTextNode("Min: "));
            trMax.appendChild(document.createTextNode("Max: "));
            let input = document.createElement("input");
            input.type = "number";
            input.name = "tMin";
            input.value = "0";
            trMin.appendChild(input.cloneNode());
            input.name = "tMax";
            input.value= "1";
            trMax.appendChild(input);
            parent.appendChild(trMin);
            parent.appendChild(trMax);
        }
    }else{
        if(parent.lastChild.textContent == "Max: "){
            // Remove Min/Max inputs
            parent.removeChild(parent.lastChild);
            parent.removeChild(parent.lastChild);
        }
    }
    return check;
}

// Funktion to load listeners to edit/remove images
function actionButtonHandlers() {
    let getArrayFromName = function(tagname) {
        // get the NodeList and transform it into an array
        return Array.prototype.slice.call(document.getElementsByName(tagname));
    };
    // Devices
    let ed = getArrayFromName("DeviceItemEdit");
    let del = getArrayFromName("DeviceItemDel");
    // Events and Types
    let e = ed.concat(getArrayFromName("EventItemEdit"), getArrayFromName("TypeItemEdit"));
    let d = del.concat(getArrayFromName("EventItemDel"), getArrayFromName("TypeItemDel"));
    // Edit
    for(let i = 0; i < e.length; i++){
        e[i].addEventListener("click", ()=>{
           editItem(e[i].id, e[i].name.substr(0,1));
        });
    }
    // Delete
    for(let i = 0; i < d.length; i++){
        d[i].addEventListener("click", ()=>{
            removeItem(d[i].id, d[i].name.substr(0,1));
        });
    }
}

// Function to remove an item
// Send del request to url via ajax
// Params: id -> name to remove, item -> item kind
function removeItem(id, item) {
    let xhr = new XMLHttpRequest();
    let data = address + "/ui/del";
    // Remove from Device list
    switch(item) {
        case "D":
            // Device
            data += "?Item="+item+"&Name="+id.substr(2,id.length);
            xhr.addEventListener("load", function () {
                let parent = $id("devices").getElementsByTagName("table")[0].lastChild;
                parent.removeChild($id(id));
            });
            break;
        case "E":
            // Event
            data += "?Item="+item+"&Name="+id.substr(2,id.length);
            xhr.addEventListener("load", function () {
                let parent = $id("events").getElementsByTagName("table")[0].lastChild;
                parent.removeChild($id(id));
            });
            break;
        case "T":
            // Type
            data += "?Item="+item+"&Name="+id.substr(2,id.length);
            xhr.addEventListener("load", function () {
                let parent = $id("types").getElementsByTagName("table")[0].lastChild;
                parent.removeChild($id(id));
            });
            break;
        default:
    }
    // Remove from overview
    let overview = $id("home").getElementsByTagName("article")[0];
    for(let e in overview.childNodes){
        if(overview.childNodes[e].id == "O_"+id.split("D_")[1]){
            overview.removeChild(overview.childNodes[e]);
        }
    }
    // Remove from server
    xhr.open("GET", data);
    xhr.send();
}

// Function to edit an item (call dialog)
// Params: id -> name to edit, item -> item kind
function editItem(id, item) {
    let element;
    let trs = document.getElementsByTagName("tr");
    for(let e in trs){
        if(trs[e].id == id){
            element = trs[e];
        }
    }
    // Get the current data to dialog window and change button
    // Name is unique and can't be updated. Here a new device has to be create
    switch(item) {
        case "D":
            // Device
            $name("dName")[0].value = element.childNodes[1].textContent;
            $name("dName")[0].readOnly = true;
            $name("dRoom")[0].value = element.childNodes[3].textContent;
            $name("dType")[0].value = element.childNodes[5].textContent;
            // Show
            updateDialogTitle("addDevice", false);
            $id("addD").value = "Update";
            $id("addDevice").style.display = "block";
            break;
        case "E":
            // Event
            $name("eName")[0].value = element.childNodes[1].textContent;
            $name("eName")[0].readOnly = true;
            $name("eTime")[0].value = element.childNodes[3].textContent;
            $name("eOffset")[0].value = element.childNodes[5].textContent;
            let devices = element.childNodes[7].textContent.split(", ");
            for(let i = 0; i < devices.length-1; i++){
                if(i != 0){
                    // new input for devices
                    addItemEvent();
                }
                // Get all devices in inputs for event dialog
                let buffer = devices[i].split("(");
                $name("eDevice")[i].value = buffer[0].split("|")[1] + " | " + buffer[0].split("|")[0];
                //$name("eDevice")[i].readOnly = true;
                $name("to")[i].value = buffer[1].substr(0,buffer[1].length-1);
            }
            // Show
            updateDialogTitle("addEvent", false);
            $id("addE").value = "Update";
            $id("addEvent").style.display = "block";
            break;
        case "T":
            // Type
            $name("tName")[0].value = element.childNodes[1].textContent;
            $name("tName")[0].readOnly = true;
            $name("tKind")[0].value = element.childNodes[3].textContent;
            if(changeMinMax()){
                $name("tMin")[0].value = element.childNodes[5].textContent;
                $name("tMax")[0].value = element.childNodes[7].textContent;
            }
            updateDialogTitle("addType", false);
            $id("addT").value = "Update";
            $id("addType").style.display = "block";
            break;
        default:
    }
}

// Helper to quick modify title of dialog
// Params: id -> Specific dialog, val -> true=Update=>Add, false=Add=>Update
function updateDialogTitle(id , val) {
    let title = $id(id).childNodes[1].childNodes[1].childNodes[1];
    if(val){
        if(title.textContent.split(" ")[0] == "Update") {
            title.textContent = title.textContent.replace("Update", "Add");
        }
    }else{
        if(title.textContent.split(" ")[0] == "Add") {
            title.textContent = title.textContent.replace("Add", "Update");
        }
    }
}

// Function show nav on mobile devices
function showNav() {
    // Toggle
    if($id("header").style.transform == "translateX(-275px)"){
        $id("header").style.transform = "translateX(0px)";
        $id("headerToggle").style.transform = "translateX(+275px)";
    }else{
        $id("header").style.transform = "translateX(-275px)";
        $id("headerToggle").style.transform = "translateX(0px)";
    }
}

// Function to set interval for ajax request of devices
function getDevices(val){
    let getString = "AllDevices";
    let data = "?Get="+getString;
    let url = address + "/ui/get" + data;
    let xhr = new XMLHttpRequest();
    xhr.addEventListener("load", ()=>{
        let response = xhr.responseText;
        if(response != "") {
            buildNewOverview(JSON.parse(response));
        }
    });
    if(val == "once"){
        requestDevices(xhr, url)
    }else{
        setInterval(()=>{requestDevices(xhr, url)}, ajaxIntervalOverview);
    }
}

// Function to get device states updated
function requestDevices(xhr, url){
    xhr.open("GET", url);
    xhr.send();
}

// Function to build a new overwiev with device data
function buildNewOverview(data){
    // Remove old content
    let parent = $id("home").getElementsByTagName("article")[0];
    while(parent.hasChildNodes()){
        parent.removeChild(parent.firstChild);
    }
    // Add new Content
    for(let i = 0; i < data.length; i++){
        let elemtent = createElement(data[i]);
        parent.appendChild(elemtent)
    }
    // Helper function to create elements to set (Type=Kind)
    function createElement(data){
        let item = null;
        let outer = document.createElement("div");
        outer.style.backgroundColor = "#333";
        outer.id = "O_" + data.Name;
        outer.className = "oItem";
        switch (data.Type){
            // Create Switch and click event listener
            case "Switch":
                // Color to show state
                // Red = off, Green = on
                let switchState = function(item, state){
                    if(state == null){
                        item.style.backgroundColor = item.style.backgroundColor == "darkred" ? "darkgreen" : "darkred";
                        item.style.color = item.style.color == "white" ? "black" : "white";
                    }else{
                        if(state == "0"){
                            item.style.backgroundColor = "darkred";
                            item.style.color = "white";
                        }else{
                            item.style.backgroundColor = "darkgreen";
                            item.style.color = "black";
                        }
                    }
                }
                item = document.createElement("div");
                item.addEventListener("click", ()=>{
                    // Actions
                    updateItem(data, item.style.backgroundColor == "darkred" ? "1" : "0");
                    switchState(item);
                });
                item.className += data.Type;
                switchState(item, data.State);
                item.appendChild(document.createTextNode(data.Name));
                item.appendChild(document.createElement("br"));
                item.appendChild(document.createTextNode(data.Room));
                break;
            // Create number and click event listener
            case "Number":
                outer.style.padding = "0.3em";
                item = document.createElement("input");
                item.type = "number";
                item.addEventListener("change", ()=>{
                    updateItem(data, item.value);
                });
                item.className += data.Type;
                item.style.backgroundColor = "grey";
                item.style.color = "black";
                item.value = data.State;
                outer.appendChild(document.createTextNode(data.Name + " - "));
                outer.appendChild(document.createTextNode(data.Room));
                outer.appendChild(document.createElement("br"));
                break;
            // Create Range and click event listener
            case "Range":
                outer.style.padding = "0.3em";
                item = document.createElement("input");
                item.type = "range";
                item.addEventListener("change", ()=>{
                    updateItem(data, item.value);
                });
                item.className += data.Type;
                item.style.backgroundColor = "grey";
                item.style.color = "black";
                item.value = data.State;
                outer.appendChild(document.createTextNode(data.Name + " - "));
                outer.appendChild(document.createTextNode(data.Room));
                outer.appendChild(document.createElement("br"));
                break;
            default:
                item = document.createElement("div");
        }
        outer.appendChild(item);
        return outer;
    }
}

// Function to update new values to db
function updateItem(data, state){
    let get = "?State=" + state + "&Name=" +data.Name;
    let url = address + "/ui/set" + get;
    let xhr = new XMLHttpRequest();
    xhr.open("GET", url);
    xhr.send();
}

//////////////////////////////////////////////////////////////////////////////////////
/// Simulator
//////////////////////////////////////////////////////////////////////////////////////

// Funktion to send action to sim data in ui
function actionButtonHandlersSim() {
    let xhr = new XMLHttpRequest();
    // State
    $id("simStateOn").addEventListener("click", ()=>{
        let url = address + "/sim/set" + "?Set=" + "State" + "&Value=true";
        xhr.open("GET", url);
        xhr.send();
    });
    $id("simStateOff").addEventListener("click", ()=>{
        let url = address + "/sim/set" + "?Set=" + "State" + "&Value=false";
        xhr.open("GET", url);
        xhr.send();
    });
    // Time
    $id("simTime").addEventListener("change", ()=>{
        let url = address + "/sim/set" + "?Set=" + "Time" + "&Value=" + $id("simTime").value;
        xhr.open("GET", url);
        xhr.send();
    });
    // Multiplier
    $id("simMultiplier").addEventListener("change", ()=>{
        let url = address + "/sim/set" + "?Set=" + "Multiplier" + "&Value=" + $id("simMultiplier").value;
        xhr.open("GET", url);
        xhr.send();
    });
    // Check stopped to enable import
    $id("data").addEventListener("mouseover", ()=>{
        if($id("currentState")){
            if($id("currentState").textContent == "On"){
                $id("dataFileSub").disabled = true;
            }else{
                $id("dataFileSub").disabled = false;
            }
        }
    });
}

// Function to set interval for ajax request of devices
function getSimStates(val){
    let getString = "States";
    let data = "?Get="+getString;
    let url = address + "/sim/get" + data;
    let xhr = new XMLHttpRequest();
    xhr.addEventListener("load", ()=>{
        let response = xhr.responseText;
        if(response != "") {
            buildNewStatelist(JSON.parse(response));
        }
    });
    if(val == "once"){
        requestStates(xhr, url)
    }else{
        setInterval(()=>{requestDevices(xhr, url)}, ajaxIntervalState);
    }
}

// Function to get device states updated
function requestStates(xhr, url){
    xhr.open("GET", url);
    xhr.send();
}

// Function to build a new overwiev with device data
function buildNewStatelist(data){
    // Remove old content
    let parent = $id("simUi").getElementsByTagName("header")[0];
    while(parent.hasChildNodes()){
        parent.removeChild(parent.firstChild);
    }
    // Add new Content
    let table = createElement(data);
    parent.appendChild(table);

    // Helper function to create
    // elements to set (Type=Kind)
    function createElement(data){
        let outer = document.createElement("table");
        outer.id = "stateTable";
        outer.style.padding = "0.5em";
        let th = document.createElement("tr");
        for(let i in data){
            if(i != "Id"){
                let td = document.createElement("td");
                td.appendChild(document.createTextNode(i));
                th.appendChild(td);
            }
        }
        outer.appendChild(th);
        let tr = document.createElement("tr");
        for(let i in data){
            if(i != "Id"){
                let td = document.createElement("td");
                let val = data[i];
                if(i == "State"){
                    td.id = "currentState";
                    val = val ? "On" : "Off";
                }
                if(i == "Time"){
                    let h = Math.floor(val/60/60);
                    let m = Math.round((val-h*60*60)/60);
                    val = h + ":" + m;
                }
                // Update action inputs
                if(i == "Multiplier"){
                    $id("simMultiplier").value = data[i];
                }
                td.appendChild(document.createTextNode(val));
                tr.appendChild(td);
            }
        }
        outer.appendChild(tr);
        return outer;
    }
}
