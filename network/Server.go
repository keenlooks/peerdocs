package main

import (
    "os"
    "fmt"
    "net"
    "encoding/gob"
)

type Token struct {
    DocID string
    TokenData string
}

type NetworkPacket struct {
    Payload Token
}

var nodeconns map[string]net.Conn
var nodeaddr map[string]string
var nodeports map[string]string

var myname string

func handleConnection(conn net.Conn) {
    dec := gob.NewDecoder(conn)
    p := &NetworkPacket{}

    dec.Decode(p)

    fmt.Printf("Received : DocID: %s : TokenData: %s \n", 
        p.Payload.DocID, p.Payload.TokenData);

    // Probably handle the token here by calling appropriate middleware function

    // Then forward it to the next node in the ring
    forwardToken(&(p.Payload))
}

func forwardToken(payload *Token) {
    var next string

    // hardcoding the ring logic for now
    if myname == "server1" {
        next = "server2"
    }

    if myname == "server2" {
        next = "server3"
    }

    if myname == "server3" {
        next = "server1"
    }

    conn, ok := nodeconns[next]
    var err error

    if ok == false {
        // Create new connection
        addr := nodeaddr[next]
        fmt.Printf("Creating connection to address: %s\n", addr);
        conn, err = net.Dial("tcp", addr)
        if err != nil {
            // handle error
            fmt.Printf("Error in creating connection to address: %s\n", addr);
        }
        // nodeconns[next] = conn
    }

    encoder := gob.NewEncoder(conn)
    
    p := &NetworkPacket {
         Payload: Token {
                   DocID: "Document1",
                   TokenData: "HelloWorld!",
         },
    }
    encoder.Encode(p)
    conn.Close()
}

func main() {
    fmt.Println("start");

    nodeconns = make(map[string]net.Conn)
    nodeaddr = make(map[string]string)
    nodeports = make(map[string]string)

    nodeaddr["server1"] = "localhost:8080" 
    nodeaddr["server2"] = "localhost:8081" 
    nodeaddr["server3"] = "localhost:8082" 

    nodeports["server1"] = ":8080";
    nodeports["server2"] = ":8081";
    nodeports["server3"] = ":8082";

    myname = os.Args[1]
    myport := nodeports[myname]
    ln, err := net.Listen("tcp", myport)

    if err != nil {
        // handle error here later
    }

    for {
        conn, err := ln.Accept() // Try to accept a connection

        if err != nil {
            // handle error
            continue
        }

        go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
    }
}
