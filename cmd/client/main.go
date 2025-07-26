package main

import (
    "fmt"
    "log"
    "github.com/atkosX/go-bittorrent/internal/utils/bencode"
)

func main() {
    data := []byte("d3:foo3:bar5:helloi52ee")

    decoded, err := bencode.Decode(data)
    if err != nil {
        log.Fatalf("Failed to decode: %v", err)
    }
    fmt.Printf("Decoded value: %#v\n", decoded)
}
