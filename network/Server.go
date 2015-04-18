package main

import (
    "os"
    "fmt"
    "net"
    "encoding/gob"
    "sync"
)

type Token struct {
    DocID       string
    Key         string
    TokenData   string
    NodeDetails Node
}

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

var tokenring  map[string](map[string]*RingInfo)
var nodes      map[string]*Node

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
        fmt.Printf("Enter Choice: (1) for create (2) for join (3) to print rings\n");
        fmt.Scanf("%d", &input);

        switch input {
        case 1: fmt.Printf("Enter DOC ID: ");
                docID = ""
                fmt.Scanln(&docID)
                fmt.Printf("Entered Details: DOC ID=%s\n", docID);
                createDoc(docID)
        case 2: fmt.Printf("Enter DOC ID: ");
                docID = ""
                fmt.Scanln(&docID)

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
                joinGroup(joinNodeName, joinNodeAddr, docID, key)
        case 3: for key, value := range tokenring { 
                    fmt.Printf("=====Printing ring for doc %s=====\n", key);
                    printTokenRing(value)
                }
        }
    }

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

        /* TODO: Call to the middleware here to update the token */

        if(p.Ptype == "FORWARD-TOKEN") {
            forwardToken(conn, &(p.Payload))
            mutex.Unlock()
        }
    }

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

    for key, value := range ring {
        np.RingEntry.PrevNode = value.PrevNode
        np.RingEntry.NextNode = value.NextNode
        np.Src = key
        enc.Encode(np)

        n = nodes[key]
        np.NodeEntry.NodeName = n.NodeName
        np.NodeEntry.NodeAddr = n.NodeAddr
        np.NodeEntry.conn = nil
        enc.Encode(np)
    }
   
    np.Ptype = "FETCH_RING_DONE"
    enc.Encode(np)

    return 
}

func fetchRing(conn net.Conn, dec *gob.Decoder, 
               docID string, key string) {
    fmt.Printf("Trying to fetch ring details for doc ID %s\n", docID)
    var re *RingInfo

    ring, ok := tokenring[docID]

    resp := &NetworkPacket{}
    var node *Node

    fmt.Printf("Waiting for ring details..\n")
    for {
        err := dec.Decode(resp)
        if err != nil {
            break;
        }

        if resp.Ptype == "FETCH_RING_DONE" {
            break
        }

        re = new(RingInfo)
        re.NextNode = resp.RingEntry.NextNode
        re.PrevNode = resp.RingEntry.PrevNode
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
}

func forwardToken(conn net.Conn, payload *Token) {
    var next string
    var curElmt *RingInfo
    var ok bool
    ring,ok := tokenring[payload.DocID]
    if ok == false {
       fmt.Printf("Cannot forward to the provided DOC ID\n"); 
       return
    }
   
    curElmt, ok = ring[myname]
    if ok == false {
         fmt.Printf("Could not retrieve the  next node name\n");
         return
    }

    nextNode,ok := nodes[curElmt.NextNode]
    if ok == false {
        fmt.Printf("Could not retrieve the  next node\n");
        return;
    }
    
    np := new(NetworkPacket);
    np.Src = myname
    np.Dst = next
    np.DstAddr = nextNode.NodeAddr;
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

        // Broadcast the ring to all other nodes here later
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
               docID string, key string) {
    np := new(NetworkPacket)
    np.Ptype = "JOIN"
    np.Src = myname
    np.SrcAddr = myaddr
    np.DstAddr = joinNodeAddr
    np.AckNeeded = false
    np.Payload.DocID = docID
    np.Payload.Key = key
    np.Payload.NodeDetails.NodeName = myname
    np.Payload.NodeDetails.NodeAddr = myaddr

    var conn net.Conn
    var err error

    joinNode, ok := nodes[joinNodename]
    if ok == false {
        conn, err = net.Dial("tcp", joinNodeAddr)
        if err != nil {
            return
        }
    } else {
        conn = joinNode.conn
    }

    enc := gob.NewEncoder(conn)
    enc.Encode(np)
    dec := gob.NewDecoder(conn)

    createDoc(docID)
    delete(tokenring[docID], myname)
    fetchRing(conn, dec, docID, key)
    printTokenRing(tokenring[docID])
    conn.Close()

    return
}

func printTokenRing(ring map[string]*RingInfo) {
    for key, value := range ring {
        fmt.Printf("Key:%s, next-node:%s, prev-node:%s\n", 
          key, value.NextNode, value.PrevNode)
    }

    return;
}

func createDoc(docID string) {
    var ring map[string]*RingInfo;
    ring = make(map[string]*RingInfo);
    new_ring_node := new(RingInfo);
    ring[myname] = new_ring_node

    tokenring[docID] = ring

    return;
}


func initialize(myname string, myaddr string) {
    tokenring = make(map[string](map[string]*RingInfo))
    nodes = make(map[string]*Node)

    mutex = &sync.Mutex{};

    this_node := new(Node)
    this_node.NodeName = myname
    this_node.NodeAddr = myaddr
    nodes[myname] = this_node;
}

func main() {
    fmt.Println("start");

    myname = os.Args[1]
    myaddr = os.Args[2]
    myport := os.Args[3]
    fmt.Printf("myname=%s, myaddr=%s, myport=%s\n", 
        myname, myaddr, myport);
    initialize(myname, myaddr)

    // Start taking inputs
    go readConsoleInput()

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
