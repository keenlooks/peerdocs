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
    Payload  Token
    Ptype    string
    Src      string
    Dst      string
    SrcAddr  string
    DstAddr  string
}

type RingInfo struct {
    prevNode string
    nextNode string
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

func handleConnection(conn net.Conn) {
    mutex.Lock()
    dec := gob.NewDecoder(conn)
    p := &NetworkPacket{}
    
    dec.Decode(p);
    fmt.Printf("Received : %+v\n", p);

    if(p.Ptype == "JOIN") {
        mutex.Unlock()
        updateTokenRing(p, myname, true);
        return;
    }
    if(p.Ptype == "CREATE") {
        mutex.Unlock()
        createDoc(p.Payload.DocID) 
        return;
    }
    if(p.Ptype == "UPDATE-RING") {
        mutex.Unlock()
        updateTokenRing(p, p.Src, false)
        return;
    }

    /* TODO: Call to the middleware here to update the token */

    // Now forward the updated token to the next node in the ring
    forwardToken(&(p.Payload))
    conn.Close();
    mutex.Unlock()
    return
}

func forwardToken(payload *Token) {
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

    nextNode,ok := nodes[curElmt.nextNode]
    if ok == false {
        fmt.Printf("Could not retrieve the  next node\n");
        return;
    }
    
    np := new(NetworkPacket);
    np.Src = myname
    np.Dst = next
    np.DstAddr = nextNode.NodeAddr;
     
    forwardPacket(np); 
}

func updateTokenRing(p *NetworkPacket, pos string, broadcast bool) {
    var newnode, nextnode, curnode, prevnode *RingInfo
    var ok bool
    var ring map[string]*RingInfo
    nextnode = nil;
    prevnode = nil;
    nodeToAdd := new(Node)
 
    docID := p.Payload.DocID
    ring, ok = tokenring[docID]
    if ok == false {
        fmt.Printf("Doc ID %s not found\n", docID);
        return
    }

    newnode, ok = ring[p.Src];
    printTokenRing(ring);
    if ok == false {
        curnode = ring[pos]
        if(curnode != nil) {
            nextnode = ring[curnode.nextNode]
            prevnode = ring[curnode.prevNode]
        }

        newnode = new(RingInfo);
        newnode.prevNode = pos;
 
        if nextnode == nil {
            newnode.nextNode = pos
        } else {
            newnode.nextNode = curnode.nextNode
        }

        if(prevnode == nil) {
            curnode.prevNode = p.Src
        }

        curnode.nextNode = p.Src;
        
        if(nextnode != nil) {
            nextnode.prevNode = p.Src
        }

        ring[p.Src] = newnode;
        nodeToAdd.NodeName = p.Src
        nodeToAdd.NodeAddr = p.SrcAddr
        nodes[p.Src] = nodeToAdd

        // Broadcast the ring to all other nodes here later
        if(broadcast == true) {
            fmt.Printf("Broadcasting...\n");
            broadcastRingUpdate(ring, nodeToAdd, docID);
        }
    } else {
        fmt.Printf("Node already part of ring\n");
    }

    printTokenRing(ring);
}

func broadcastRingUpdate(ring map[string]*RingInfo, newnode *Node, docID string) {
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
        np.Src = myname
        np.Dst = key
        np.DstAddr = node.NodeAddr
        np.Payload.DocID = docID

        nodeToAdd.NodeName = newnode.NodeName
        nodeToAdd.NodeAddr = newnode.NodeAddr        
        np.Payload.NodeDetails = nodeToAdd

        forwardPacket(np)
    }
    
    return
}

func joinGroup(joinNodeAddr string, docID string, key string) {
    np := new(NetworkPacket)
    np.Ptype = "JOIN"
    np.Src = myname
    np.DstAddr = joinNodeAddr
    np.Payload.DocID = docID
    np.Payload.Key = key

    forwardPacket(np)
}

func forwardPacket(np *NetworkPacket) {
    var conn net.Conn
    var err error

    conn,err = net.Dial("tcp", np.DstAddr)
    if err != nil {
        fmt.Printf("Could not create connection\n");
        return
    }

    encoder := gob.NewEncoder(conn)
    encoder.Encode(np)
    conn.Close() 
}


func printTokenRing(ring map[string]*RingInfo) {
    for key, value := range ring {
        fmt.Printf("Key:%s, next-node:%s, prev-node:%s\n", 
          key, value.nextNode, value.prevNode)
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
