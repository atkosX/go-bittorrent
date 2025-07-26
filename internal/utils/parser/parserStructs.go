package parser

import "encoding/hex"

type TorrentFile struct {
    Announce string
    Info     Info
    InfoHash [20]byte
}

func (t *TorrentFile) GetInfoHashString() string {
	return hex.EncodeToString(t.InfoHash[:])
}

func (t *TorrentFile) GetTotalLength() int {
	if t.Info.Length > 0 {
		return t.Info.Length
	}
	var totalLength int
	for _, file := range t.Info.Files {
		totalLength += file.Length
	}
	return totalLength
}

type Info struct {
    Name        string
    PieceLength int
    Pieces      string 
    Length      int   
    Files       []File 
}

type File struct {
    Length int
    Path   []string
}