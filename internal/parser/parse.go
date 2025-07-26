package parser

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/atkosX/go-bittorrent/internal/utils/bencode"
)

func ParseTorrentFile(data []byte) (*TorrentFile, error) {
	if data == nil || len(data) == 0 {
		return nil, errors.New("empty torrent data")
	}

	decoded, err := bencode.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding torrent data: %w", err)
	}

	torrentDict, ok := decoded.(map[string]bencode.Bvalue)
	if !ok {
		return nil, errors.New("decoded data is not a dictionary")
	}

	tf := &TorrentFile{}

	if announce, ok := torrentDict["announce"].(string); ok {
		tf.Announce = announce
	} else {
		return nil, errors.New("missing announce url")
	}

	infoValue, ok := torrentDict["info"]
	if !ok {
		return nil, errors.New("missing info dictionary")
	}

	infoDict, ok := infoValue.(map[string]bencode.Bvalue)
	if !ok {
		return nil, errors.New("info is not a dictionary")
	}

	info := Info{}
	if name, ok := infoDict["name"].(string); ok {
		info.Name = name
	}
	if pieceLength, ok := infoDict["piece length"].(int); ok {
		info.PieceLength = pieceLength
	}
	if pieces, ok := infoDict["pieces"].(string); ok {
		info.Pieces = pieces
	}

	if length, ok := infoDict["length"].(int); ok {
		info.Length = length
	} else if files, ok := infoDict["files"].([]bencode.Bvalue); ok {
		for _, fileEntry := range files {
			if fileDict, ok := fileEntry.(map[string]bencode.Bvalue); ok {
				file := File{}
				if length, ok := fileDict["length"].(int); ok {
					file.Length = length
				}
				if path, ok := fileDict["path"].([]bencode.Bvalue); ok {
					for _, p := range path {
						if pathStr, ok := p.(string); ok {
							file.Path = append(file.Path, pathStr)
						}
					}
				}
				info.Files = append(info.Files, file)
			}
		}
	} else {
		return nil, errors.New("missing length or files in info dictionary")
	}

	tf.Info = info

	var infoBuf bytes.Buffer
	err = bencode.Encode(&infoBuf, infoValue)
	if err != nil {
		return nil, fmt.Errorf("error encoding info dict: %w", err)
	}
	tf.InfoHash = sha1.Sum(infoBuf.Bytes())

	return tf, nil
}