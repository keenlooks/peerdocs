package main

import (
    "fmt"
    "net"
    "encoding/gob"
    "sync"
)

type NetworkPacket struct {
    Payload   Token
    RingEntry RingInfo
    NodeEntry Node
    Ptype     string
    Src       string
    Dst       string
    SrcAddr   string
    DstAddr   string
    AckNeeded bool
    DocCon    Docfetch
}

type RingInfo struct {
    PrevNode string
    NextNode string
}

type Node struct {
    NodeName string
    NodeAddr string
    conn     net.Conn
}

type Docs struct {
    DocID         string
    Key           string
    cond          *sync.Cond
    packetarrived bool
    Payload       Token
}

type Invitations struct {
    Doc             Docs
    inviteeNodeName string
    inviteeNodeAddr string
}

var tokenring  map[string](map[string]*RingInfo)
var nodes      map[string]*Node
var docs       map[string]*Docs
var invites    map[string]*Invitations

var myname string
var myaddr string
var mutex *sync.Mutex;

func readConsoleInput() {
    var input int
    var docID string
    var joinNodeAddr string
    var joinNodeName string
    var key string

    for {
        input = 0
        fmt.Printf("Enter Choice: (1) Create (2) Join (3) Invite (4) Print rings (5) Display Invites\n");
        fmt.Scanf("%d", &input);

        switch input {
        case 1: fmt.Printf("Enter DOC ID: ");
                docID = ""
                fmt.Scanln(&docID)

                fmt.Printf("Enter key for document: ")
                key = ""
                fmt.Scanln(&key)

                fmt.Printf("Entered Details: DOC ID=%s, key=%s\n", docID, key);
                createDocument(docID, key)
        case 2: fmt.Printf("Enter DOC ID: ");
                docID = ""
                fmt.Scanln(&docID)
/*
                fmt.Printf("Enter node addr to contact: ");
                joinNodeAddr = ""
                fmt.Scanln(&joinNodeAddr)

                fmt.Printf("Enter node name to contact: ");
                joinNodeName = ""
                fmt.Scanln(&joinNodeName)

                fmt.Printf("Enter Doc key: ");
                key = ""
                fmt.Scanln(&key)
                fmt.Printf("Entered Details: DOC ID=%s, Node Addr=%s, Node Name=%s, Key=%s\n", 
                    docID, joinNodeAddr, joinNodeName, key);
*/
                fmt.Printf("Entered Details: DOC ID=%s\n", docID);
                joinGroup(joinNodeName, joinNodeAddr, docID, key, false)
        case 3: fmt.Printf("Enter DOC ID: ");
                docID = ""
                fmt.Scanln(&docID)

                fmt.Printf("Enter invitee addr: ");
                joinNodeAddr = ""
                fmt.Scanln(&joinNodeAddr) 

                fmt.Printf("Enter invitee node name:");
                joinNodeName = ""
                fmt.Scanln(&joinNodeName)

                fmt.Printf("Entered Details: DOC ID=%s, Invitee Node Addr=%s, Invitee Node Name=%s\n", 
                    docID, joinNodeAddr, joinNodeName);
                sendInvitation(joinNodeAddr,joinNodeName, docID)
        case 4: for key, value := range tokenring { 
                    fmt.Printf("=====Printing ring for doc %s=====\n", key);
                    printTokenRing(value)
                }
        case 5: fmt.Printf("===== Invitations======\n");
                for _, value := range invites { 
                    fmt.Printf("DocID: %s, Key: %s, Invitor Name: %s, Invitor Addr: %s\n",
                        value.Doc.DocID, value.Doc.Key, value.inviteeNodeName, value.inviteeNodeAddr);
                }
                fmt.Printf("=======================\n");
        }
    }

    return
}

func sendInvitation(inviteAddr string, inviteNodename string, docID string) {
    _, ok := docs[docID]
    if ok == false {
        fmt.Printf("Document doesnot exist\n")
        return
    }
  
    _, ok = nodes[inviteNodename]
    if ok != false {
        fmt.Printf("Node already part of ring for doc\n")
        return
    }

    conn, err := net.Dial("tcp", inviteAddr)
    if err != nil {
        fmt.Printf("Could not establish connection to given node\n")
        return
    }

    np := new(NetworkPacket)
    np.Ptype = "INVITE"
    np.Src = myname
    np.SrcAddr = myaddr
    np.Dst = inviteNodename
    np.DstAddr = inviteAddr
    np.AckNeeded = false
    np.Payload.DocID = docID
    np.Payload.Key = docs[docID].Key

    enc := gob.NewEncoder(conn)
    enc.Encode(np)

    conn.Close()
    return
}

func handleConnection(conn net.Conn) {
    var err error

    var p *NetworkPacket

    for { 
        dec := gob.NewDecoder(conn)
        enc := gob.NewEncoder(conn)

        p = &NetworkPacket{}
        err = dec.Decode(p);
        if err != nil {
            fmt.Printf("Error in handle connection, exiting. Error = %s\n", err);
            conn.Close()
            return
        }

        fmt.Printf("Received : %+v\n", p);
        mutex.Lock()

        if(p.Ptype == "JOIN") {
            updateTokenRing(conn, enc, dec, p, myname, true);
            sendRing(conn, enc, p)
            mutex.Unlock()
            conn.Close()
        }

        if(p.Ptype == "UPDATE-RING") {
            updateTokenRing(conn, enc, dec, p, p.Src, false)
            mutex.Unlock()
        }

        if(p.Ptype == "FETCH-RING") {
            sendRing(conn, enc, p)
            mutex.Unlock()
        }
        
        if(p.Ptype == "INVITE") {
            receiveInvitation(conn, enc, dec, p) 
            mutex.Unlock()
        }

        /* TODO: Call to the middleware here to update the token */

        if(p.Ptype == "FORWARD-TOKEN") {
            doc, ok := docs[p.Payload.DocID]
            if(ok == false) {
                fmt.Printf("Could not retrieve doc info for %s\n", p.Payload.DocID)
                mutex.Unlock()
                return
            }
            doc.cond.L.Lock()
            doc.packetarrived = true
            doc.Payload = p.Payload
            doc.cond.L.Unlock()
            doc.cond.Signal()
            mutex.Unlock()
        }
    }

    return
}

func receiveInvitation(conn net.Conn, enc *gob.Encoder, 
                       dec *gob.Decoder, np *NetworkPacket) {
    invite := new(Invitations)
    invite.Doc.DocID = np.Payload.DocID
    invite.Doc.Key = np.Payload.Key
    invite.inviteeNodeName = np.Src
    invite.inviteeNodeAddr = np.SrcAddr

    invites[np.Src] = invite
    return
}

func sendRing(conn net.Conn, enc *gob.Encoder, p *NetworkPacket) {
    fmt.Printf("Trying to send ring details for %s\n", p.Payload.DocID)
    ring, ok := tokenring[p.Payload.DocID]

    if ok == false {
        fmt.Printf("Failed to fetch ring details for doc %s\n", p.Payload.DocID)
        return
    }

    np := new(NetworkPacket)
    var n *Node
    var textsent bool
    textsent = false

    for key, value := range ring {
        np.RingEntry.PrevNode = value.PrevNode
        np.RingEntry.NextNode = value.NextNode
        np.Src = key
        if textsent == false {
            np.DocCon = fetchDoc(p.Payload.DocID)
            textsent = true
        }
        enc.Encode(np)

        n = nodes[key]
        np.NodeEntry.NodeName = n.NodeName
        np.NodeEntry.NodeAddr = n.NodeAddr


        np.NodeEntry.conn = nil
        enc.Encode(np)
    }
   
    np.Ptype = "FETCH-RING-DONE"
    enc.Encode(np)
   
    numelem := len(ring) 
    if(numelem == 2) {
        fmt.Printf("Initiating ring forwards for doc %s\n",  p.Payload.DocID)
        doc := docs[p.Payload.DocID]
        doc.cond.L.Lock()
        doc.packetarrived = true
        doc.Payload.DocID = p.Payload.DocID
        doc.cond.L.Unlock()
        doc.cond.Signal() 
    }
    
    return 
}

func fetchRing(conn net.Conn, dec *gob.Decoder, 
               docID string, key string) {
    fmt.Printf("Trying to fetch ring details for doc ID %s\n", docID)
    var re *RingInfo
    var doccreated bool
    var ring map[string]*RingInfo
    var ok bool
    doccreated = false

    ring, ok = tokenring[docID]
    if ok == false {
        ring = make(map[string]*RingInfo)
    }

    resp := &NetworkPacket{}
    var node *Node

    fmt.Printf("Waiting for ring details..\n")
    
    for {
        err := dec.Decode(resp)
        if err != nil {
            fmt.Printf("Error in decoding = %s\n", err)
            break;
        }

        if resp.Ptype == "FETCH-RING-DONE" {
            break
        }

        re = new(RingInfo)
        re.NextNode = resp.RingEntry.NextNode
        re.PrevNode = resp.RingEntry.PrevNode
        if(doccreated == false) {
            createDocWithId(resp.DocCon)
            doccreated = true
        }
        ring[resp.Src] = re

        node, ok = nodes[resp.Src]
        dec.Decode(resp)
        if(ok == false) {
            node = new(Node)
            node.NodeName = resp.NodeEntry.NodeName
            node.NodeAddr = resp.NodeEntry.NodeAddr
            nodes[resp.Src] = node
        }
    }
    tokenring[docID] = ring
    printTokenRing(ring)

    var nodeDet *Node
    var host Host
    for nodename, _ := range ring {
        nodeDet = nodes[nodename]
        host.Name = nodeDet.NodeName
        host.Address = nodeDet.NodeAddr
        host.DocID = docID
        host.DocKey = key
        updateDocNodeWithId(host)
    }
}

func forwardToken(docID string) {
    var next string
    var curElmt *RingInfo
    var ok bool
    var doc *Docs
    doc,ok = docs[docID]
    fmt.Printf("Forward Token thread created for doc %s\n", docID);
    if ok == false {
        fmt.Printf("doc %s does not exist\n", docID)
        return
    }

    for {
        doc.cond.L.Lock()
        for doc.packetarrived == false {
            doc.cond.Wait()
        }
        
        doc.packetarrived = false
        fmt.Printf("[forwardToken] Calling handleToken for docID %s\n", docID)
        newToken := handleToken(doc.Payload)

        ring,ok := tokenring[docID]

        if ok == false {
            fmt.Printf("Cannot forward to the provided DOC ID\n"); 
            doc.cond.L.Unlock()
            continue
        }
   
        curElmt, ok = ring[myname]
        if ok == false {
            fmt.Printf("Could not retrieve the next node name\n");
            doc.cond.L.Unlock()
            continue
        }
        
        nextNode,ok := nodes[curElmt.NextNode]
        if ok == false {
            fmt.Printf("Could not retrieve the  next node\n");
            doc.cond.L.Unlock()
            return;
        }
        
        np := new(NetworkPacket);
        np.Src = myname
        np.SrcAddr = myaddr
        np.Dst = next
        np.DstAddr = nextNode.NodeAddr;
        np.Payload = newToken
        np.Ptype = "FORWARD-TOKEN"

        conn, err := net.Dial("tcp", np.DstAddr)
        if err == nil {
            enc := gob.NewEncoder(conn)
            enc.Encode(np) 
        } else {
            fmt.Printf("Error in dialing to %s. Error = %s\n", np.DstAddr, err)
        }
        doc.cond.L.Unlock()
    }

    return
}

func 
updateTokenRing(conn net.Conn, enc *gob.Encoder, dec *gob.Decoder, 
                p *NetworkPacket, pos string, broadcast bool) int {
    var srcRingEntry, nxtRingEntry, curRingEntry *RingInfo
    var ok bool
    var ring map[string]*RingInfo
    nxtRingEntry = nil;
    var nodeToAdd *Node

    docID := p.Payload.DocID
    ring, ok = tokenring[docID]
    if ok == false {
        fmt.Printf("Doc ID %s not found\n", docID);
        if p.AckNeeded == true {
            p.Ptype = "ACK-FAIL"
            enc.Encode(p)
        }

        _, ok = nodes[p.Src]
        if ok == false {
            // Close this connection since no nodes are using this conenction
        }

        return -1
    }

    srcRingEntry, ok = ring[p.Payload.NodeDetails.NodeName];
    printTokenRing(ring);
    if ok == false {
        curRingEntry = ring[pos]
        numelem := len(ring)
        fmt.Printf("Lnght of ring  is %d\n", numelem)
        fmt.Printf("Inserting at pos: %s\n", pos)
        srcRingEntry = new(RingInfo);
        if(numelem == 1) {
            srcRingEntry.PrevNode = pos;
            srcRingEntry.NextNode = pos;

            curRingEntry.PrevNode = p.Payload.NodeDetails.NodeName;
            curRingEntry.NextNode = p.Payload.NodeDetails.NodeName;

            ring[p.Payload.NodeDetails.NodeName] = srcRingEntry;

        } else {

            curRingEntry = ring[pos]
            nxtRingEntry = ring[curRingEntry.NextNode]

            srcRingEntry.PrevNode = pos;
            srcRingEntry.NextNode = curRingEntry.NextNode

            nxtRingEntry.PrevNode = p.Payload.NodeDetails.NodeName
            curRingEntry.NextNode = p.Payload.NodeDetails.NodeName
 
            ring[p.Payload.NodeDetails.NodeName] = srcRingEntry;
        }

        nodeToAdd, ok = nodes[p.Payload.NodeDetails.NodeName]
        if ok == false {
            fmt.Printf("Creating new node..\n");
            nodeToAdd = new(Node)
            nodeToAdd.NodeName = p.Payload.NodeDetails.NodeName
            nodeToAdd.NodeAddr = p.Payload.NodeDetails.NodeAddr
            nodeToAdd.conn = conn
            nodes[p.Payload.NodeDetails.NodeName] = nodeToAdd
        }

        if(broadcast == true) {
            fmt.Printf("Broadcasting...\n");
            broadcastRingUpdate(ring, nodeToAdd, docID);
        }
    } else {
        fmt.Printf("Node already part of ring\n");
    }

    printTokenRing(ring);
    return 0
}

func 
broadcastRingUpdate(ring map[string]*RingInfo,
                    newnode *Node, docID string) {
    var nodeToAdd Node;

    for key, value := range ring {
        if key == myname {
            continue;
        }

        if key == newnode.NodeName {
            continue;
        }

        value = value; //Placate the compiler
        node := nodes[key]
        np := new(NetworkPacket)

        np.Ptype = "UPDATE-RING"
        np.AckNeeded = false
        np.Src = myname
        np.Dst = key
        np.DstAddr = node.NodeAddr
        np.Payload.DocID = docID

        nodeToAdd.NodeName = newnode.NodeName
        nodeToAdd.NodeAddr = newnode.NodeAddr        
        np.Payload.NodeDetails = nodeToAdd

        fmt.Printf("Sending ring update to %s with addr %s\n", 
            node.NodeName, node.NodeAddr)
        var conn net.Conn
        var err error
        conn,err = net.Dial("tcp", node.NodeAddr)
        if err != nil {
            fmt.Printf("broadcastRingUpdate: Error in net.Dial = %s\n", err)
            continue
        }

        enc := gob.NewEncoder(conn)
        enc.Encode(np)
    }
   
    fmt.Printf("Done with broadcasting..\n") 
    return
}

func joinGroup(joinNodename string, joinNodeAddr string, 
               docID string, key string, bootstrap bool) {
    if(bootstrap == false) {
       invitation, ok := invites[joinNodename]
       if(ok == false) {
           fmt.Printf("No invitation recorded for this document\n")
            return
       }

       joinNodename = invitation.inviteeNodeName
       joinNodeAddr = invitation.inviteeNodeAddr 
       key = invitation.Doc.Key
    }

    np := new(NetworkPacket)
    np.Ptype = "JOIN"
    np.Src = myname
    np.SrcAddr = myaddr
    np.DstAddr = joinNodeAddr
    np.Dst = joinNodename
    np.AckNeeded = false
    np.Payload.DocID = docID
    np.Payload.Key = key
    np.Payload.NodeDetails.NodeName = myname
    np.Payload.NodeDetails.NodeAddr = myaddr

    var conn net.Conn
    var err error

    conn, err = net.Dial("tcp", np.DstAddr)
    if err != nil {
        fmt.Printf("Failed to create connection to np.DstAddr\n")
        return
    }

    enc := gob.NewEncoder(conn)
    enc.Encode(np)
    dec := gob.NewDecoder(conn)

    fetchRing(conn, dec, docID, key)
    printTokenRing(tokenring[docID])
    conn.Close()

    doc := new(Docs)
    doc.DocID = docID
    doc.Key = key
    doc.cond = &sync.Cond{L: &sync.Mutex{}}
    doc.packetarrived = false
    docs[docID] = doc

    go forwardToken(docID)

    return
}

func printTokenRing(ring map[string]*RingInfo) {
    for nodename, value := range ring {
        fmt.Printf("Nodename:%s, next-node:%s, prev-node:%s\n", 
          nodename, value.NextNode, value.PrevNode)
    }

    return;
}

func createDocument(docID string, key string) {
    var ring map[string]*RingInfo;
    ring = make(map[string]*RingInfo);
    new_ring_node := new(RingInfo);
    ring[myname] = new_ring_node
    tokenring[docID] = ring


    for nodename, value := range ring {
        fmt.Printf("NodeName:%s, next-node:%s, prev-node:%s\n", 
          nodename, value.NextNode, value.PrevNode)
    }

    doc := new(Docs)
    doc.DocID = docID
    doc.Key = key
    doc.cond = &sync.Cond{L: &sync.Mutex{}}
    doc.packetarrived = false
    docs[docID] = doc

    go forwardToken(docID)

    return;
}


func 
initializeNetworkServer(name string, addr string, myport string) {
    myname = name
    myaddr = addr
    //myport := port

    tokenring = make(map[string](map[string]*RingInfo))
    nodes     = make(map[string]*Node)
    docs      = make(map[string]*Docs)
    invites   = make(map[string]*Invitations)

    mutex = &sync.Mutex{};

    this_node := new(Node)
    this_node.NodeName = myname
    this_node.NodeAddr = myaddr
    nodes[myname] = this_node;

    fmt.Printf("myname=%s, myaddr=%s, myport=%s\n", 
        myname, myaddr, myport);

    // Start taking inputs
    go readConsoleInput()

    fmt.Printf("%s is starting the network server listen port at %s\n", myname, myaddr);
    ln, err := net.Listen("tcp", myport)

    if err != nil {
        fmt.Printf("Could not create listening socket\n");
        return
    }

    for {
        conn, err := ln.Accept() // Try to accept a connection

        if err != nil {
            // handle error
            continue
        }

        go handleConnection(conn)
    }
}

//func main() {
//    fmt.Println("start");
//
//}
