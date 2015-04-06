package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "net"
    //"encoding/gob"
    "encoding/json"
)

var docFolderPath = "./docs"
var localChanges = []string{}
var listenPort = ":12345"

type Token struct {
    DocID string
    totalChanges []string
}

type FrontEndRequest struct {
    Command 		string	`json:"command"`
    Argument 		string	`json:"argument"`
    Changearray []string	`json:"changearray"`
}

type FrontEndResponse struct {
	Responsestring 	string	`json:"responsestring"`
	Changearray []string	`json:"changearray"`
}

func listDocs()(string){
    var docList string
    
    //check whether the given file or directory exists or not
    _, err := os.Stat(docFolderPath)
    if os.IsNotExist(err) { return "No Docs" }
    files, _ := ioutil.ReadDir("./docs/")
    if len(files) == 0 { return "No Docs" }
	for _, f := range files {
		docList += f.Name() + "|"	
	}
	return docList[:len(docList)-1]
}

func createDoc()(string){
    //create a file in the correct XML format with sections: <DocID>, <GroupKey>, <GroupList>, <Text>
}

func joinGroup(DocID string)(bool){
	//will connect to token ring of group described in file at path "string"
	//TODO - Handled by Network layer
	return true
}

//takes connection as argument and decodes using json the command and arguments from FrontEnd
func receiveCommand(conn net.Conn)(string, string, []string){
	
	dec := json.NewDecoder(conn)
    p := &FrontEndRequest{}

    dec.Decode(p)
	//TODO listen for commands from client
	
	return p.Command, p.Argument, p.Changearray
}

func sendResponse(conn net.Conn, response FrontEndResponse){
	//send payload to requester
	encoder := json.NewEncoder(conn)
    p := &response
    encoder.Encode(p)
    conn.Close()
}

func handleConn(conn net.Conn){
	//receive command and argument - command can be UPDATE, FETCH, LIST, JOIN, LEAVE, or CREATE
    command, argument, changearray := receiveCommand(conn)

    //process command and argument
    response := process(command, argument, changearray)

    //send response back to client
    sendResponse(conn, response)
}

func process(command string, argument string, changearray []string)(FrontEndResponse){   

    response := ""
    var changes []string
    
    switch command {
    	case "UPDATE":
    		//requires DocID as argument
    		//used to push updates that user is providing to document. append changes to change array in token
    		if argument == "" || len(changearray) == 0 {
    			//TODO send error message back to requester
    			response = "command requires DocID as argument and changearray to have changes"
    			break
    		}

    	case "FETCH":
    		//requires DocID as argument
    		//used by UI to ask for official changes to document (made official when token has passed through node and other users' changes have been added)
    		if argument == "" {
    			//TODO send error message back to requester
    			response = "command requires argument"
    			break
    		}

    		
    	case "LIST":
    		//used by UI to see what Docs this node has access to (has access keys for)
   			response = listDocs()

    	case "JOIN":
    		//requires DocID as argument
    		//used by UI with argument of base64 encoded string of DocID:groupList:GroupKey
    			//DocID is a string in form <creating host MAC>.<number identifier>
    			//groupList is array of IPs in string form to check for bootstrapping
    		if argument == ""{
    			//TODO send error message back to requester
    			response = "command requires argument"
    			break
    		}

    		if joinGroup(argument){
    			response = "success"
    		}else{
    			response = "fail"
    		}

    	case "LEAVE":
    		//requires DocId as argument
    		//used by UI to delete a key (the doc file) and remove itself from a group sends an update to all nodes to remove it from the list.
    		if argument == ""{
    			//TODO send error message back to requester
    			response = "command requires argument"
    			break
    		}

    	case "CREATE":
    		//used by UI to initiate a Doc, returns a DocID
    

    }


    return FrontEndResponse{
    	Responsestring: response,
    	Changearray: changes,
    }
}

func main() {
    //testing json request
    /*conn, err := net.Dial("tcp", "localhost:8080")
    encoder := json.NewEncoder(conn)
    p := &FrontEndRequest {
         Command: "LIST", 
         Argument: "ertqewqr",
         Changearray: []string{"hello"}}
    fmt.Println(p)
    bla, _ := json.Marshal(p)
    fmt.Println(bla)
    encoder.Encode(p)
    conn.Close()
    fmt.Println("done");
    return*/

    
    //requests sent to the server must be in form: 
    //  {"command":"<command>","argument":"<argument>","changearray":<position:string array of changes>}
    fmt.Println("Server running...")

    ln, err := net.Listen("tcp", listenPort)

    if err != nil {
        fmt.Println("error listening for connection")
    }

    for {
        conn, err := ln.Accept() // Try to accept a connection

        if err != nil {
            fmt.Println("error accepting connection")
        }

        go handleConn(conn)
    }
}


