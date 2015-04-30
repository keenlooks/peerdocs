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
var numPastLocalChanges = map[string]int{}
var listenPort = ":8080"
var tokenListenPort = ":12345"
var backspacestring = "\\b"
var cursorPos = map[string]int{}
var selfName = ""
var selfAddr = ""
var selfPort = ""
var docmodified = map[string]bool{}

type Hostarray struct{
    Hosts []Host            `json:"hosts"`
}

type HostInvite struct {
    TypeRequest string      `json:"type"`
    Name string             `json:"name"`
    Address string          `json:"address"`
    DocID int               `json:"docid"`
    DocKey string           `json:"dockey"`
}

type Host struct {
    Name string             `json:"name"`
    Address string          `json:"address"`
    DocID string            `json:"docid"`
    DocKey string           `json:"dockey"`
}

type Token struct {
    DocID       string
    Updates string     //used to update groupList with new member
    Changes []Change
    Key         string
    TokenData   string
    NodeDetails Node
}

type DocdeltsNoID struct{
  Doccgs []Change           `json:"doccgs"`
}

type Docdelts struct{
  Id int                    `json:"docid"`
  Cursorpos int             `json:"cursor"`
  Doccgs []Change           `json:"doccgs"`
}

type ChangeNoID struct {
    Position int            `json:"location"`
    Charstoappend string    `json:"mod"`
}

type Change struct {
//    Id int                  `json:"id"`
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
    Cursorpos int           `json:"cursor"`
    Title string            `json:"title"`
    Ctext string            `json:"ctext"`
}

type Doccreate struct {
    Title string            `json:"title"`
    Ctext string            `json:"ctext"`
    Cursorpos int           `json:"cursor"`
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
    p := &response
    responseB, _ := json.Marshal(p)
    responsestring = string(responseB)

    //Check for invitations


    responsestring="{\"docmetas\":"+responsestring+"}"
    io.WriteString(w,responsestring)
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
            fmt.Println("cannot open "+f.Name()+" in list")
        }
        buf := make([]byte, 8192)
        count, _ := fopened.Read(buf)
        //fmt.Println(string(buf))
        if count == 0 {return []Docmeta{}}
        dm.Title = strings.Split(strings.Split(string(buf), "<Title>")[1], "</Title>")[0]
        if(docmodified[f.Name()]){
            dm.Lastmod = "true" //f.ModTime().String()
        }else{
            dm.Lastmod = "false"
        }
		docList = append(docList, *dm)
	}

    //check for invites

    for k,_ := range invites {
        dm := &Docmeta{}
        dm.Id, _ = strconv.Atoi(k)
        dm.Title = "<not joined>"
        dm.Lastmod = "pending"
        docList = append(docList, *dm)
    }
    return docList
}

func inviteNodeHttp(w http.ResponseWriter, req *http.Request){

    if(req.Method == "OPTIONS"){
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with,Content-Type")
        w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,GET,POST")
        w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
        w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")
        return
    }

    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with,Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")

    req.ParseForm()
    buf := make([]byte, 14)
    req.Body.Read(buf)

    hostinvite := &HostInvite{}
    decoder := json.NewDecoder(req.Body)
    decoder.Decode(hostinvite)
    

    //Call backend
    if(hostinvite.TypeRequest == "invite"){
        fmt.Println("inviting "+hostinvite.Name)
        sendInvitation(hostinvite.Address, hostinvite.Name, strconv.Itoa(hostinvite.DocID))
    }
    if(hostinvite.TypeRequest == "join"){
        fmt.Println("joining "+ hostinvite.Name)
        joinGroup("", hostinvite.Address, strconv.Itoa(hostinvite.DocID), hostinvite.DocKey, false)
    }

    responseB, _ := json.Marshal(hostinvite)
    responsestring := string(responseB)
    responsestring="{\"invitation\": {\"id\":"+strconv.Itoa(rand.Int()%int(math.Pow(2,float64(32))))+","+responsestring[1:]+"}";
    io.WriteString(w, responsestring)
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
    buf := make([]byte, 7)
    req.Body.Read(buf)
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

    f.WriteString("<DocID>"+strconv.Itoa(macaddr+counter)+"</DocID>\n<Title>"+dc.Title+"</Title>\n<GroupKey>"+"TODO"/*generate secure key and make it base64*/+"</GroupKey>\n<GroupList></GroupList>\n<Text>"+dc.Ctext+"</Text>")

    host := Host{}
    host.Name = selfName
    host.Address = selfAddr
    host.DocID = strconv.Itoa(macaddr+counter)
    host.DocKey = strconv.Itoa(rand.Int()%int(math.Pow(2,float64(32))))


    numPastLocalChanges[host.DocID] = 0
    updateDocNodeWithId(host)
    createDocument(host.DocID, host.DocKey)

    dm := &Docfetch{}
    dm.Cursorpos = dc.Cursorpos
    dm.Id = macaddr + counter
    //fstat, _ := f.Stat()
    dm.Title = dc.Title
    
    dm.Ctext = dc.Ctext
    f.Close()

    return *dm
}

func createDocWithId(dc Docfetch)(){
    //create a file in the correct XML format with sections: <DocID>, <GroupKey>, <GroupList>, <Text>
    
    //create the file
    f, err := os.Create(docFolderPath+strconv.Itoa(dc.Id))
    if err != nil{
        fmt.Println("Could not create file "+strconv.Itoa(dc.Id))
    }
    numPastLocalChanges[strconv.Itoa(dc.Id)] = 0

    f.WriteString("<DocID>"+strconv.Itoa(dc.Id)+"</DocID>\n<Title>"+dc.Title+"</Title>\n<GroupKey></GroupKey>\n<GroupList>"+/*put yourself in group list*/"</GroupList>\n<Text>"+dc.Ctext+"</Text>")
    f.Close()
}


func updateDocNodeWithId(host Host){
    fopened, err := os.Open(docFolderPath+host.DocID)
    if(err != nil){
        fmt.Println("cannot open "+host.DocID+" for reading")
    }
    buf := make([]byte, 4096)
    count, _ := fopened.Read(buf)
    if count == 0 {fmt.Println("cannot read file "+host.DocID)}
    fopened.Close()
    inputstring := string(buf)
    inputstringtext := strings.Split(strings.Split(inputstring, "<GroupList>")[1], "</GroupList>")[0]

    grouplist := Hostarray{}
    err = json.Unmarshal([]byte(inputstringtext), &grouplist)
    grouplist.Hosts = append(grouplist.Hosts, host)
    responseB, _ := json.Marshal(&grouplist)
    inputstringtext = string(responseB)

    fopened, err = os.Open(docFolderPath+host.DocID)
    if(err != nil){
        fmt.Println("cannot open "+host.DocID+" for writing")
    }
    buf = make([]byte, 4096)
    count, _ = fopened.Read(buf)
    fopened.Close()
    inputstring = string(buf)
    inputstringbeforetext := strings.Split(inputstring, "<GroupKey>")[0]
    inputstringaftertext := strings.Split(inputstring, "</GroupList>")[1]

    outputstring := inputstringbeforetext+ "<GroupKey>"+host.DocKey+"</GroupKey><GroupList>"+inputstringtext+"</GroupList>"+inputstringaftertext
    ioutil.WriteFile(docFolderPath+host.DocID,[]byte(outputstring), 0666)


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

func updateFile(DocID string)(bool){

    fopened, err := os.Open(docFolderPath+DocID)
    if(err != nil){
        fmt.Println("cannot open "+DocID+" for reading")
        return false
    }
    buf := make([]byte, 8192)
    count, _ := fopened.Read(buf)
    if count == 0 {return false}
    fopened.Close()
    inputstring := string(buf)
    inputstringtext := strings.Split(strings.Split(inputstring, "<Text>")[1], "</Text>")[0]
    changes := officialChanges[DocID]
    for _, change := range changes {
        if len(inputstringtext) != 0{
            if !(change.Position < 0 || change.Position > len(inputstringtext)-1){
                inputstringtext = inputstringtext[:change.Position]+change.Charstoappend+inputstringtext[change.Position:]
            }else{
                if(change.Position == len(inputstringtext)){
                    inputstringtext = inputstringtext + change.Charstoappend
                }else{
                    fmt.Println("received location out of bounds "+strconv.Itoa(change.Position)+" mod: "+change.Charstoappend)
                }
            }
        }else{
            inputstringtext = change.Charstoappend
        }
    }

    //backspace out any backspaces in there
    for strings.Index(inputstringtext,backspacestring) != -1 {
        if strings.Index(inputstringtext,backspacestring)-1 >= 0{
            inputstringtext = inputstringtext[:strings.Index(inputstringtext,backspacestring)-1] + inputstringtext[strings.Index(inputstringtext,backspacestring)+len(backspacestring):]
        }else{
            inputstringtext = inputstringtext[:strings.Index(inputstringtext,backspacestring)] + inputstringtext[strings.Index(inputstringtext,backspacestring)+len(backspacestring):]
        }
       // cursorPos[DocID] -= len(backspacestring)
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
    officialChanges[token.DocID]=officialChanges[token.DocID][:0]
    //update own changes and files with changes in token, update token
    for _, change := range token.Changes {
        for _,localchange := range localChanges[token.DocID]{
            if change.Position <= localchange.Position {
                localchange.Position += len(change.Charstoappend)//-(strings.Count(change.Charstoappend,backspacestring)*(len(backspacestring)+1)))
            } 
        }
        if change.Position <= cursorPos[token.DocID] {
            cursorPos[token.DocID] += (len(change.Charstoappend)-strings.Count(change.Charstoappend,backspacestring)*(len(backspacestring)+1))
        }
    }
    token.Changes = append(token.Changes[numPastLocalChanges[token.DocID]:], localChanges[token.DocID]...)
    numPastLocalChanges[token.DocID] = len(localChanges[token.DocID])
    
    //clear local changes
    localChanges[token.DocID] = localChanges[token.DocID][:0]

    //make all changes official
    officialChanges[token.DocID] = append(officialChanges[token.DocID], token.Changes...)

    //TODO - change GroupList based on updates string in token

    //update local files with changes
    if len(officialChanges[token.DocID]) == 0 {
        //fmt.Println("no changes to "+token.DocID)
        return token
    }

    if updateFile(token.DocID){
        //once file is updated clear official list
        docmodified[token.DocID] = true
        return token
    }

    fmt.Println("could not update "+token.DocID)
    return token
}


func updateChangesHttpGet(w http.ResponseWriter, req *http.Request){
    req.ParseForm()
    
    dd := &Docdelts{}
    //fmt.Println(string(buf))
    buf := make([]byte, 11)
    req.Body.Read(buf)
    decoder := json.NewDecoder(req.Body)
    decoder.Decode(dd)
    //fmt.Println(dd)
    DocID := strconv.Itoa(dd.Id) //strings.Split(req.URL.Path, "docdelts/")[1]
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with,Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")
    if req.Method == "PUT"{
        //fmt.Println("updating "+DocID)
        if(docmodified[DocID]){ //if the handleToken was just called cursorPos is now the delta and not the actual position
            //cursorPos[DocID] += dd.Cursorpos
        }else{
            cursorPos[DocID] = dd.Cursorpos
        }
        if !updateChanges(DocID, dd.Doccgs){ fmt.Println("update changes for "+DocID+" didnt work")}



        //THIS PART NEEDS TO BE REMOVED AFTER INTEGRATION WITH TOKEN PASSING
       /* officialChanges[DocID] = append(officialChanges[DocID], localChanges[DocID]...)
        localChanges[DocID] = nil
        if updateFile(DocID){

            //once file is updated clear official list
            officialChanges[DocID]=nil
        }else{
            fmt.Println("update failed")
        }*/    
        //THIS PART NEEDS TO BE REMOVED AFTER INTEGRATION WITH TOKEN PASSING

        p := &dd
        responseB, _ := json.Marshal(p)
        responsestring := string(responseB)
        
        if docmodified[DocID] {
            responsestring = "{\"docdelt\": {\"id\":11223344,\"docid\":"+DocID+",\"cursor\":"+strconv.Itoa(cursorPos[DocID])+",\"doccgs\":[{\"id\":0,\"location\":0, \"mod\":\"a\"}]}}"
            //docmodified[DocID] = false
        }else{
            responsestring = "{\"docdelt\": {\"id\":11223344,\"docid\":"+DocID+",\"cursor\":"+strconv.Itoa(cursorPos[DocID])+",\"doccgs\":[]}}"
        }
        io.WriteString(w,responsestring)
    }else{
    io.WriteString(w,"{\"docdelt\": {\"docid\":"+strings.Split(req.URL.Path, "docdelts/")[1]+",\"cursor\":"+strconv.Itoa(cursorPos[DocID])+",\"doccgs\":[]}}")
}
}

func updateChangesHttp(w http.ResponseWriter, req *http.Request){

    req.ParseForm()
    dd := &Docdelts{}
    //fmt.Println(string(buf))
    buf := make([]byte, 11)
    req.Body.Read(buf)
    decoder := json.NewDecoder(req.Body)
    decoder.Decode(dd)
    //fmt.Println(dd)
    //dc.Title = req.URL.Query().Get("title")
    //dc.Ctext = req.FormValue("ctext")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Origin,x-requested-with,Content-Type")
    w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")

    if req.Method == "POST"{
        io.WriteString(w,"{\"docdelt\": {\"id\":11223344,\"docid\":"+strconv.Itoa(dd.Id)+",\"cursor\":"+strconv.Itoa(dd.Cursorpos)+",\"doccgs\":[]}}")
    }else if req.Method == "OPTIONS"{
        //just headers
    }

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
    //fmt.Println(req.URL.Path)
    //fmt.Println(strings.Split(req.URL.Path, "docs/")[1])
    response := fetchDoc(strings.Split(req.URL.Path, "docs/")[1])
    /*DocID := strconv.Itoa(response.Id)

    recentlocalchanges := []Change{}
    copy(recentlocalchanges, localChanges[DocID])

    for _, change := range officialChanges[DocID] {
        for _,localchange := range recentlocalchanges{
            if change.Position <= localchange.Position {
                localchange.Position += len(change.Charstoappend)//-(strings.Count(change.Charstoappend,backspacestring)*(len(backspacestring)+1)))
            } 
        }
        if change.Position <= cursorPos[DocID] {
            cursorPos[DocID] += (len(change.Charstoappend)-strings.Count(change.Charstoappend,backspacestring)*(len(backspacestring)+1))
        }
    }
    inputstringtext := response.Ctext
    for _, change := range recentlocalchanges {
        if len(inputstringtext) != 0{
            if !(change.Position < 0 || change.Position > len(inputstringtext)-1){
                inputstringtext = inputstringtext[:change.Position]+change.Charstoappend+inputstringtext[change.Position:]
            }else{
                if(change.Position == len(inputstringtext)){
                    inputstringtext = inputstringtext + change.Charstoappend
                }else{
                    fmt.Println("received location out of bounds "+strconv.Itoa(change.Position)+" mod: "+change.Charstoappend)
                }
            }
        }else{
            inputstringtext = change.Charstoappend
        }
    }*/
    
/*
    for strings.Index(inputstringtext,backspacestring) != -1 {
        if strings.Index(inputstringtext,backspacestring)-1 >= 0{
            inputstringtext = inputstringtext[:strings.Index(inputstringtext,backspacestring)-1] + inputstringtext[strings.Index(inputstringtext,backspacestring)+len(backspacestring):]
        }else{
            inputstringtext = inputstringtext[:strings.Index(inputstringtext,backspacestring)] + inputstringtext[strings.Index(inputstringtext,backspacestring)+len(backspacestring):]
        }
       // cursorPos[DocID] -= len(backspacestring)
    }

    response.Ctext = inputstringtext*/

    //if response has not changed
    if(docmodified[strconv.Itoa(response.Id)]){
        
    }else{
        response.Title = "None"
        docmodified[strconv.Itoa(response.Id)] = false
    }

    //add headers
    responsestring := ""
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers","Content-Type,Origin,x-requested-with")
    w.Header().Set("Access-Control-Allow-Methods", "PUT,PATCH,GET,POST")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length")

    p := &response
    responseB, _ := json.Marshal(p)
    responsestring = string(responseB)
    responsestring="{\"doc\":"+responsestring+"}"
    io.WriteString(w,responsestring)
}

func fetchDoc(DocID string)(Docfetch){

    fopened, _ := os.Open("docs/"+DocID)
    buf := make([]byte, 8192)
    count, _ := fopened.Read(buf)
    if count == 0 {
        fmt.Println("cannot open "+DocID+" for fetching")
        return Docfetch{}
    }

    df := Docfetch{}
    df.Id,_ = strconv.Atoi(DocID)
    df.Cursorpos = cursorPos[DocID]
    df.Title = strings.Split(strings.Split(string(buf), "<Title>")[1], "</Title>")[0]
    df.Ctext = strings.Split(strings.Split(string(buf), "<Text>")[1], "</Text>")[0]
    return df
}



func joinGroupsFromDoc(){
    doclist := listDocs()
    for _, doc := range doclist {
        if(doc.Lastmod == "pending"){continue}
        fopened, err := os.Open(docFolderPath+strconv.Itoa(doc.Id))
        if(err != nil){
            fmt.Println("cannot open doc")
        }
        buf := make([]byte, 256)
        fopened.Read(buf)
        //fmt.Println(string(buf))
        grouplist := Hostarray{}
        groupliststring := strings.Split(strings.Split(string(buf), "<GroupList>")[1], "</GroupList>")[0]
        if groupliststring == "" {continue}
        err = json.Unmarshal([]byte(groupliststring), &grouplist)
        for _, host := range grouplist.Hosts {
            if(host.Name != selfName){
                fmt.Println("joining "+host.Name+" at "+host.Address + " with ID "+host.DocID)
                if joinGroup(host.Name, host.Address, host.DocID, host.DocKey, true){
                    break
                }
            }
        }
    }
}

func main() {
   

    fmt.Println("Server running...")
    selfName = os.Args[1]
    selfAddr = os.Args[2]
    selfPort = os.Args[3]
    go initializeNetworkServer(selfName, selfAddr, selfPort)
    joinGroupsFromDoc()
    //start listening for clients
    http.HandleFunc("/api/docmeta", listDocsHttp)
    http.HandleFunc("/api/docs", createDocHttp)
    http.HandleFunc("/api/docs/", fetchDocHttp)
    http.HandleFunc("/api/docdelts", updateChangesHttp)
    http.HandleFunc("/api/docdelts/", updateChangesHttpGet)
    http.HandleFunc("/api/invitations", inviteNodeHttp)
    err := http.ListenAndServe(os.Args[4], nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
