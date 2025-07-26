package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/atkosX/go-bittorrent/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/client/main.go <torrent-file>")
		return
	}
	torrentFilePath := os.Args[1]
	data, err := ioutil.ReadFile(torrentFilePath)
	if err != nil {
		panic(err)
	}

	torrent, err := parser.ParseTorrentFile(data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name: %s\n", torrent.Info.Name)
	fmt.Printf("Announce: %s\n", torrent.Announce)
	fmt.Printf("Info Hash: %s\n", torrent.GetInfoHashString())
	fmt.Printf("Total Length: %d bytes\n", torrent.GetTotalLength())
}