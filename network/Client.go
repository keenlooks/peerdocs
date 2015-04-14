package main

import (
    "fmt"
    "log"
    "net"
    "encoding/gob"
)

type Token struct {
    DocID       string
    Key         string
    TokenData   string
    NodeDetails Node
}

type Node struct {
    NodeName string
    NodeAddr string
    conn     net.Conn
}

type NetworkPacket struct {
    Payload  Token
    Ptype    string
    Src      string
    Dst      string
    SrcAddr  string
    DstAddr  string
}

func main() {
    fmt.Println("start client");
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatal("Connection error", err)
    }
    encoder := gob.NewEncoder(conn)

    p := new(NetworkPacket)
    p.Ptype = "JOIN";
    p.Src = "client1";
    p.Dst = "server1"
    p.DstAddr = "localhost:8080"

    p.Payload.DocID = "Document1";
    p.Payload.TokenData = "HelloWorld!";
    p.Payload.Key = "TrustMe";

    p.Payload.NodeDetails.NodeName = "server2"
    p.Payload.NodeDetails.NodeAddr = "localhost:8081"

    fmt.Printf("Src: %s, Type: %s, DocID: %s : TokenData: %s \n", 
        p.Src, p.Ptype, p.Payload.DocID, p.Payload.TokenData);
    fmt.Printf("Received : %+v", p);
    encoder.Encode(p)
    conn.Close()
    fmt.Println("done");
}
