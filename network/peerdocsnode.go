package main


import (
    "os"
    "fmt"
    "io/ioutil"
    "net"
    "strconv"
    //"xml"
    "strings"
  //  "encoding/gob"
    "encoding/json"
    //"sort"
    "math"
    "math/rand"
    "io"
    "net/http"
    "log"
)

var docFolderPath = "./docs/"
var localChanges = map[string][]Change{}
var officialChanges = map[string][]Change{}
var listenPort = ":8080"
var tokenListenPort = ":12345"

type Token struct {
    DocID string
    Updates string     //used to update groupList with new member
    Changes []Change
}

type Docdelts struct{
  Id int                    `json:"id"`
  Doccgs []Change           `json:"doccgs"`
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

type Docfetch struct {
    Id int                  `json:"id"`
    Title string            `json:"title"`
    Ctext string            `json:"ctext"`
}

type Doccreate struct {
    Title string            `json:"title"`
    Ctext string            `json:"ctext"`
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
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with,Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")
    //encoder := json.NewEncoder(w)
    p := &response
    //encoder.Encode(p)
    responseB, _ := json.Marshal(p)
    responsestring = string(responseB)
    //responsestring = "Access-Control-Allow-Credentials:true\nAccess-Control-Allow-Headers:Origin,x-requested-with\nAccess-Control-Allow-Methods:PUT,PATCH,GET,POST\nAccess-Control-Allow-Origin:*\nAccess-Control-Expose-Headers:Content-Length" + responsestring
    responsestring="{\"docmetas\":"+responsestring+"}"
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
        nameint, _ := strconv.Atoi(f.Name())
        dm.Id = nameint
        counter += 1
        fopened, err := os.Open(docFolderPath+f.Name())
        if(f.Name()[0]=='.'){continue}
        if(err != nil){
            fmt.Println("cannot open doc")
        }
        buf := make([]byte, 128)
        count, _ := fopened.Read(buf)
        //fmt.Println(string(buf))
        if count == 0 {return []Docmeta{}}
        dm.Title = strings.Split(strings.Split(string(buf), "<Title>")[1], "</Title>")[0]
        dm.Lastmod = f.ModTime().String()
		docList = append(docList, *dm)
	}	
    return docList
}

func createDocHttp(w http.ResponseWriter, req *http.Request){
    if(req.Method == "OPTIONS"){
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with,Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")
        return
    }
    req.ParseForm()
    dc := &Doccreate{}
    //fmt.Println(string(buf))
    decoder := json.NewDecoder(req.Body)
    decoder.Decode(dc)
    //dc.Title = req.URL.Query().Get("title")
    //dc.Ctext = req.FormValue("ctext")
    fmt.Println(dc.Title)

    response := createDoc(*dc)
   responseB, _ := json.Marshal(response)
    responsestring := string(responseB)
    responsestring="{\"doc\":"+responsestring+"}";
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with,Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")
    //convert response to correct structure
    //encoder := json.NewEncoder(w)
   // p := &response
    //encoder.Encode(p)
    io.WriteString(w, responsestring)
}

func createDoc(dc Doccreate)(Docfetch){
    //create a file in the correct XML format with sections: <DocID>, <GroupKey>, <GroupList>, <Text>
    
    //get MAC address
    interfaces, err := net.Interfaces()
    macaddrstring := interfaces[1].HardwareAddr.String()
    macaddrstring = strings.Replace(macaddrstring,":","",-1) 
    macaddrint64,_ := strconv.ParseInt(macaddrstring,16,0)
    rand.Seed(macaddrint64)
    macaddr := rand.Int()%int(math.Pow(2,float64(32)))

    //check for doc directory, if does not exist, create it
    _, err = os.Stat(docFolderPath)
    if os.IsNotExist(err) { if (os.Mkdir(docFolderPath, 0xFFF) != nil) { return Docfetch{} }}
    files, _ := ioutil.ReadDir(docFolderPath)
    
    counter := 0

    for _, filename := range files {if strconv.Itoa(macaddr+counter) == filename.Name(){counter+=1}}

    //create the file
    f, err := os.Create(docFolderPath+strconv.Itoa(macaddr+counter))
    if err != nil{
        fmt.Println("Could not create file "+strconv.Itoa(macaddr+counter))
        return Docfetch{}
    }

    f.WriteString("<DocID>"+strconv.Itoa(macaddr+counter)+"</DocID>\n<Title>"+dc.Title+"</Title>\n<GroupKey>"+"TODO"/*generate secure key and make it base64*/+"</GroupKey>\n<GroupList>"+"TODO"/*put yourself in group list*/+"</GroupList>\n<Text>"+dc.Ctext+"</Text>")

    dm := &Docfetch{}

    dm.Id = macaddr + counter
    //fstat, _ := f.Stat()
    dm.Title = dc.Title
    
    dm.Ctext = dc.Ctext
    f.Close()

    return *dm
}

func joinGroup(argument string)(bool){
	//will connect to token ring of group described by base64 encoded "argument"
    
    //create doc with contents 
    f, err := os.Create(docFolderPath+strings.Split(strings.Split(argument, "<DocID>")[1], "</DocID>")[0])
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
    return os.Remove(docFolderPath+argument) == nil
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

/*func listenToken(){
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
}*/

func updateFile(DocID string)(bool){

    fopened, err := os.Open(docFolderPath+DocID)
    if(err != nil){
        fmt.Println("cannot open "+DocID+" for reading")
        return false
    }
    buf := make([]byte, 4096)
    count, _ := fopened.Read(buf)
    if count == 0 {return false}
    fopened.Close()
    inputstring := string(buf)
    inputstringtext := strings.Split(strings.Split(inputstring, "<Text>")[1], "</Text>")[0]
    changes := officialChanges[DocID]

    for _, change := range changes { 
        inputstringtext = inputstringtext[:change.Position]+change.Charstoappend+inputstringtext[change.Position:]
    }

    fopened, err = os.Open(docFolderPath+DocID)
    if(err != nil){
        fmt.Println("cannot open "+DocID+" for writing")
        return false
    }
    buf = make([]byte, 4096)
    count, _ = fopened.Read(buf)
    fopened.Close()
    inputstring = string(buf)
    inputstringbeforetext := strings.Split(inputstring, "<Text>")[0]

    outputstring := inputstringbeforetext+ "<Text>"+inputstringtext+"</Text>"
    ioutil.WriteFile(docFolderPath+DocID,[]byte(outputstring), 0666)
    return true
}

func handleToken(token Token)(Token){
    //dec := gob.NewDecoder(conn)
    //token := &Token{}
    //dec.Decode(token)
    //conn.Close()

    //update own changes and files with changes in token, update token
    for _, change := range token.Changes {
        for _,localchange := range localChanges[token.DocID]{
            if change.Position <= localchange.Position {
                localchange.Position += len(change.Charstoappend)
            }
        }
    }  
    token.Changes = append(token.Changes, localChanges[token.DocID]...)

    //clear local changes
    localChanges[token.DocID] = nil

    //make all changes official
    officialChanges[token.DocID] = append(officialChanges[token.DocID], token.Changes...)

    //TODO - change GroupList based on updates string in token

    //update local files with changes
    if updateFile(token.DocID){

    //once file is updated clear official list
    officialChanges[token.DocID]=nil
    return token
    }
    fmt.Println("could not update "+token.DocID)
    return token

    /*return // will be removed out once rest is implemented
    
    //TODO - find next available node by looking through group list and trying to connect to each one in order (or some other way)
        //TODO - try first host after you
    conn2, err := net.Dial("tcp", "<host>"+tokenListenPort)
    for err != nil {
        //TODO - if doesnt work, increment to next possible host and try in for loop
        conn2, err = net.Dial("tcp", "<host>"+tokenListenPort)
    }

    encoder := gob.NewEncoder(conn2)
    encoder.Encode(token)
    conn2.Close()*/
}


func updateChangesHttpGet(w http.ResponseWriter, req *http.Request){
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with,Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")

    io.WriteString(w,"{'docdelt': {'id':"+strings.Split(req.URL.Path, "docdelts/")[1]+",'doccgs':[]}}")
}

func updateChangesHttp(w http.ResponseWriter, req *http.Request){

    req.ParseForm()
    dd := &Docdelts{}
    //fmt.Println(string(buf))
    decoder := json.NewDecoder(req.Body)
    decoder.Decode(dd)
    fmt.Println(dd)
    //dc.Title = req.URL.Query().Get("title")
    //dc.Ctext = req.FormValue("ctext")

    if req.Method == "PUT"{
        fmt.Println("updating "+strconv.Itoa(dd.Id))

        if !updateChanges(strconv.Itoa(dd.Id), dd.Doccgs){ fmt.Println("update changes for "+strconv.Itoa(dd.Id)+" didnt work")}



        //THIS PART NEEDS TO BE REMOVED AFTER INTEGRATION WITH TOKEN PASSING
        DocID := strconv.Itoa(dd.Id)
        officialChanges[DocID] = append(officialChanges[DocID], localChanges[DocID]...)
        localChanges[DocID] = nil
        if updateFile(DocID){

            //once file is updated clear official list
            officialChanges[DocID]=nil
        }    
        //THIS PART NEEDS TO BE REMOVED AFTER INTEGRATION WITH TOKEN PASSING


        p := &dd
        //encoder.Encode(p)
        responseB, _ := json.Marshal(p)
        responsestring := string(responseB)
        //responsestring = "Access-Control-Allow-Credentials:true\nAccess-Control-Allow-Headers:Origin,x-requested-with\nAccess-Control-Allow-Methods:PUT,PATCH,GET,POST\nAccess-Control-Allow-Origin:*\nAccess-Control-Expose-Headers:Content-Length" + responsestring
        responsestring = "{\"docdelts\":"+responsestring+"}"
        io.WriteString(w,responsestring)
    }else if req.Method == "POST"{
        io.WriteString(w,"{'docdelt': {'id':"+strconv.Itoa(rand.Int()%int(math.Pow(2,float64(32))))+",'doccgs':[]}}")
    }else if req.Method == "OPTIONS"{
        //just headers
    }

    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with,Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")
    return
}

func updateChanges(DocID string, updates []Change)(bool){
    //updates must be sorted with oldest entry as 0th elementcccc
    localChanges[DocID] = append(localChanges[DocID] , updates...)
    return true
}

func fetchOfficialChanges(DocID string)(official []Change){
    return officialChanges[DocID]
}

func fetchDocHttp(w http.ResponseWriter, req *http.Request){
    fmt.Println(req.URL.Path)
    //fmt.Println(strings.Split(req.URL.Path, "docs/")[1])
    response := fetchDoc(strings.Split(req.URL.Path, "docs/")[1])
    //convert response to correct structure
    responsestring := ""
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Content-Type,Origin,x-requested-with")
    w.Header().Set("Access-Control-Allow-Methods", "PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
    //encoder := json.NewEncoder(w)
    p := &response
    //encoder.Encode(p)
    responseB, _ := json.Marshal(p)
    responsestring = string(responseB)
    //responsestring = "Access-Control-Allow-Credentials:true\nAccess-Control-Allow-Headers:Origin,x-requested-with\nAccess-Control-Allow-Methods:PUT,PATCH,GET,POST\nAccess-Control-Allow-Origin:*\nAccess-Control-Expose-Headers:Content-Length" + responsestring
    responsestring="{\"docs\":"+responsestring+"}"
    io.WriteString(w,responsestring)
}

func fetchDoc(DocID string)(Docfetch){

    fopened, _ := os.Open("docs/"+DocID)
    buf := make([]byte, 2048)
    count, _ := fopened.Read(buf)
    if count == 0 {return Docfetch{}}

    df := Docfetch{}
    df.Id,_ = strconv.Atoi(DocID)
    df.Title = strings.Split(strings.Split(string(buf), "<Title>")[1], "</Title>")[0]
    df.Ctext = strings.Split(strings.Split(string(buf), "<Text>")[1], "</Text>")[0]
    return df
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
    //go listenToken()
    fmt.Println("Server running...")

    //start listening for clients
    http.HandleFunc("/api/docmeta", listDocsHttp)
    http.HandleFunc("/api/docs", createDocHttp)
    http.HandleFunc("/api/docs/", fetchDocHttp)
    http.HandleFunc("/api/docdelts", updateChangesHttp)
    http.HandleFunc("/api/docdelts/", updateChangesHttpGet)
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
/*
http://localhost:8080/api/docdelts

{
 'docdelts':
{
  id:43
  doccgs:[
  {
    id: 66
    location: 43
    mod:'insertted text'
  },
   {
    id: 67
    location: 43
    mod:'insertted text'
  },
  {
    id: 68
    location: 43
    mod:'insertted text'
  },

  ]
}

}
*/

// GET http://localhost:8080/api/docs/ID



