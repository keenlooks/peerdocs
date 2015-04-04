package main

import (
    //"os"
    "fmt"
    //"net"
    //"encoding/gob"
)

type Token struct {
    DocID string
    TokenData string
}

func listDocs()(string){
    var docList string

	return docList 
}

func receiveCommand()(string, string){
	var command string
	var argument string
	
	//TODO listen for commands from client
	command = "LIST"
	argument = ""

	
	return command, argument
}

func send(payload string){
	//TODO - send payload string to requester
}

func main() {
    fmt.Println("start")
    
    var command string
    var argument string

    //receive command and store in "command" - can be UPDATE, FETCH, LIST, JOIN, LEAVE, or CREATE
    command, argument = receiveCommand()

    switch command {
    	case "UPDATE":
    		//requires DocID as argument
    		//used to push updates that user is providing to document. append changes to change array in token
    		if argument == ""{
    			//TODO send error message back to requester
    			break
    		}

    	case "FETCH":
    		//requires DocID as argument
    		//used by UI to ask for official changes to document (made official when token has passed through node and other users' changes have been added)
    		if argument == ""{
    			//TODO send error message back to requester
    			break
    		}

    		
    	case "LIST":
    		//used by UI to see what Docs this node has access to (has access keys for)
    	
    	case "JOIN":
    		//requires DocID as argument
    		//used by UI with argument of base64 encoded string of DocID:groupList:GroupKey
    			//DocID is a string in form <creating host MAC>.<number identifier>
    			//groupList is array of IPs in string form to check for bootstrapping
    		if argument == ""{
    			//TODO send error message back to requester
    			break
    		}

    	case "LEAVE":
    		//requires DocId as argument
    		//used by UI to delete a key and remove itself from a group sends an update to all nodes to remove it from the list.
    		if argument == ""{
    			//TODO send error message back to requester
    			break
    		}

    	case "CREATE":
    		//used by UI to initiate a Doc, returns a DocID
    		
    }
}


