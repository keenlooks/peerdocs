package main

import (
    "fmt"
    "log"
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

func main() {
    fmt.Println("start client");
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatal("Connection error", err)
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
    fmt.Println("done");
}
