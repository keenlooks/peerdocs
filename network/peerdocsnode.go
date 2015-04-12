package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "net"
    "strconv"
    //"xml"
    "strings"
    "encoding/gob"
    "encoding/json"
    //"sort"
    "io"
    "net/http"
    "log"
)

var docFolderPath = "./docs"
var localChanges = map[string][]Change{}
var officialChanges = map[string][]Change{}
var listenPort = ":8080"
var tokenListenPort = ":12345"

type Token struct {
    DocID string
    Updates string     //used to update groupList with new member
    Changes []Change
}

type Change struct {
    Id int                  `json:"id"`
    Position int            `json:"location"`
    Charstoappend string    `json:"mod"`
}

type Docmeta struct {
    Id int                  `json:"id"`
    Title string            `json:"title"`
    Lastmod string          `json:"lastmod"`
}

type FrontEndRequest struct {
    Command 		string	`json:"command"`
    Argument 		string	`json:"argument"`
    Changearray []Change	`json:"changearray"`
}

type FrontEndResponse struct {
	Responsestring 	string	`json:"responsestring"`
	Changearray []Change	`json:"changearray"`
}

func listDocsHttp(w http.ResponseWriter, req *http.Request) {
    response := listDocs()
    //convert response to correct structure
    responsestring := ""
    w.Header().Set("Access-Control-Allow-Origin", "lvh.me")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with")
    w.Header().Set("Access-Control-Allow-Methods", "PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
    //encoder := json.NewEncoder(w)
    p := &response
    //encoder.Encode(p)
    responseB, _ := json.Marshal(p)
    responsestring = string(responseB)
    //responsestring = "Access-Control-Allow-Credentials:true\nAccess-Control-Allow-Headers:Origin,x-requested-with\nAccess-Control-Allow-Methods:PUT,PATCH,GET,POST\nAccess-Control-Allow-Origin:*\nAccess-Control-Expose-Headers:Content-Length" + responsestring
    io.WriteString(w,responsestring)
    //encoder.Encode(p)
    //io.WriteString(w, response)
}

func listDocs()([]Docmeta){    
    //check whether the given file or directory exists or not
    _, err := os.Stat(docFolderPath)
    if os.IsNotExist(err) { return []Docmeta{} }
    files, _ := ioutil.ReadDir(docFolderPath)
    if len(files) == 0 { return []Docmeta{} }
    counter := 1
    docList := []Docmeta{}
	for _, f := range files {
        dm := &Docmeta{}
        dm.Id = counter
        counter += 1
        dm.Title = f.Name()
        dm.Lastmod = f.ModTime().String()
		docList = append(docList, *dm)
	}	
    return docList
}

func createDocHttp(w http.ResponseWriter, req *http.Request){
    response := createDoc()
    //convert response to correct structure
    fmt.Println(req)
    encoder := json.NewEncoder(w)
    p := &response
    encoder.Encode(p)
}

func createDoc()(Docmeta){
    //create a file in the correct XML format with sections: <DocID>, <GroupKey>, <GroupList>, <Text>
    
    //get MAC address
    interfaces, err := net.Interfaces()
    macaddr := interfaces[1].HardwareAddr.String()

    //check for doc directory, if does not exist, create it
    _, err = os.Stat(docFolderPath)
    if os.IsNotExist(err) { if (os.Mkdir(docFolderPath, 0xFFF) != nil) { return Docmeta{} }}
    files, _ := ioutil.ReadDir(docFolderPath)
    
    counter := 0

    for _, filename := range files {if macaddr+strconv.Itoa(counter) == filename.Name(){counter+=1}}

    //create the file
    f, err := os.Create("docs/"+macaddr+strconv.Itoa(counter))
    if err != nil{
        fmt.Println("Could not create file "+macaddr+strconv.Itoa(counter))
        return Docmeta{}
    }

    f.WriteString("<DocID>"+macaddr+strconv.Itoa(counter)+"</DocID>\n<GroupKey>"+"TODO"/*generate secure key and make it base64*/+"</GroupKey>\n<GroupList>"+"TODO"/*put yourself in group list*/+"</GroupList>\n<Text></Text>")

    dm := &Docmeta{}
    dm.Id = 1
    fstat, _ := f.Stat()
    dm.Title = fstat.Name()
    
    dm.Lastmod = fstat.ModTime().String()
    f.Close()

    return *dm
}

func joinGroup(argument string)(bool){
	//will connect to token ring of group described by base64 encoded "argument"
    
    //create doc with contents 
    f, err := os.Create("docs/"+strings.Split(strings.Split(argument, "<DocID>")[1], "</DocID>")[0])
    if err != nil{
        return false
    }

    f.WriteString(argument)

    f.Close()

	return true
}

func leaveGroup(argument string)(bool){

    //TODO tell all other nodes in group contained in "argument" to delete your IP address by adding something to the token

    //then delete file
    return os.Remove("doc/"+argument) == nil
}

//takes connection as argument and decodes using json the command and arguments from FrontEnd
func receiveCommand(conn net.Conn)(string, string, []Change){
	dec := json.NewDecoder(conn)
    p := &FrontEndRequest{}
    dec.Decode(p)
	return p.Command, p.Argument, p.Changearray
}

func sendResponse(conn net.Conn, response FrontEndResponse){
	//send payload to requester
	encoder := json.NewEncoder(conn)
    p := &response
    encoder.Encode(p)
    conn.Close()
}

/*func handleConn(conn net.Conn){

// hello world, the web server
    
	//receive command and argument - command can be UPDATE, FETCH, LIST, JOIN, LEAVE, or CREATE
    command, argument, changearray := receiveCommand(conn)

    //process command and argument
    response := process(command, argument, changearray)

    //send response back to client
    sendResponse(conn, response)
}*/

func listenToken(){
    ln, err := net.Listen("tcp", tokenListenPort)

    if err != nil {
        fmt.Println("error listening for connection")
    }

    for {
        conn, err := ln.Accept() // Try to accept a connection
        if err != nil {
            fmt.Println("error accepting connection")
        }

        go handleToken(conn)
    }
}

func handleToken(conn net.Conn){
    dec := gob.NewDecoder(conn)
    token := &Token{}
    dec.Decode(token)
    conn.Close()
    //update own changes and files with changes in token, update token
    token.Changes = append(token.Changes, localChanges[token.DocID]...)
    
    //clear local changes
    localChanges[token.DocID] = nil

    //make all changes official
    officialChanges[token.DocID] = append(officialChanges[token.DocID], token.Changes...)

    //TODO - change GroupList based on updates string in token

    //TODO - update local files with changes


    return // will be removed out once rest is implemented
    
    //TODO - find next available node by looking through group list and trying to connect to each one in order (or some other way)
        //TODO - try first host after you
    conn2, err := net.Dial("tcp", "<host>"+tokenListenPort)
    for err != nil {
        //TODO - if doesnt work, increment to next possible host and try in for loop
        conn2, err = net.Dial("tcp", "<host>"+tokenListenPort)
    }

    encoder := gob.NewEncoder(conn2)
    encoder.Encode(token)
    conn2.Close()
}

func updateChanges(DocID string, updates []Change)(bool){
    //updates must be sorted with oldest entry as 0th elementcccc
    localChanges[DocID] = append(localChanges[DocID] , updates...)
    return true
}

func fetchOfficialChanges(DocID string)(official []Change){
    return officialChanges[DocID]
}



/*func process(command string, argument string, changearray []Change)(FrontEndResponse){   

    response := ""
    var changes = []Change{}
    
    switch command {
    	case "UPDATE":
    		//requires DocID as argument
    		//used to push updates that user is providing to document. append changes to change array in token
    		if argument == "" || len(changearray) == 0 {
    			response = "command requires DocID as argument and changearray to have changes"
    			break
    		}

            //append changes to global localChanges slice
            updateChanges(argument, changearray)

    	case "FETCH":
    		//requires DocID as argument
    		//used by UI to ask for official changes to document (made official when token has passed through node and other users' changes have been added)
    		if argument == "" {
    			//TODO send error message back to requester
    			response = "command requires argument"
    			break
    		}

            //return official changes contained in officialChanges array
            changes = fetchOfficialChanges(argument)

    		
    	case "LIST":
    		//used by UI to see what Docs this node has access to (has access keys for)
   			response = listDocs()

    	case "JOIN":
    		//requires DocID as argument
    		//used by UI with argument of base64 encoded string of DocID:groupList:GroupKey
    			//DocID is a string in form <creating host MAC>.<number identifier>
    			//groupList is array of IPs in string form to check for bootstrapping
    		if argument == ""{
    			response = "command requires argument"
    			break
    		}

    		if joinGroup(argument){
    			response = "success"
    		}else{
    			response = "fail"
    		}

    	case "LEAVE":
    		//requires DocID as argument
    		//used by UI to delete a key (the doc file) and remove itself from a group sends an update to all nodes to remove it from the list.
    		if argument == ""{
    			response = "command requires argument"
    			break
    		}
            if leaveGroup(argument){
                response = "success"
            }else{
                response = "fail"
            }

    	case "CREATE":
    		//used by UI to initiate a Doc, returns a DocID
            response = createDoc()

    }


    return FrontEndResponse{
    	Responsestring: response,
    	Changearray: changes,
    }
}*/

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
    

    /*ln, err := net.Listen("tcp", listenPort)

    if err != nil {
        fmt.Println("error listening for connection")
    }*/
    
    //start listening for tokens
    go listenToken()
    fmt.Println("Server running...")

    //start listening for clients
    http.HandleFunc("/api/docmeta", listDocsHttp)
    http.HandleFunc("/api/createdoc", createDocHttp)
    err := http.ListenAndServe(listenPort, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

    /*
    for {
        conn, err := ln.Accept() // Try to accept a connection
        if err != nil {
            fmt.Println("error accepting connection")
        }

        go handleConn(conn)
    }
    */
}


