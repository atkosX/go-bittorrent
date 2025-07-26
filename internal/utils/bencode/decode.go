package bencode

import(
    "bytes"
    "errors"
    "io"
    "strconv"
    "bufio"
)

type Bvalue interface {}

func Decode(data []byte)(Bvalue,error){
    reader:=bufio.NewReader(bytes.NewReader(data))
    return decodeValue(reader)
}

func decodeValue(reader *bufio.Reader) (Bvalue, error) {
    b, err := reader.Peek(1)
    if err != nil {
        return nil, err
    }
    switch b[0] {
    case 'i':
        return decodeInteger(reader)
    case 'l':
        return decodeList(reader)
    case 'd':
        return decodeDict(reader)
    default:
        if b[0] >= '0' && b[0] <= '9' {
            return decodeString(reader)
        }
        return nil, errors.New("invalid bencode data")
    }
}

func decodeString(reader *bufio.Reader) (Bvalue, error){

    strLen:=[]byte{}
    for{
        b,err:=reader.ReadByte()
        if err != nil {
            return nil, err
        }
        if b==':'{
            break
        }
        if b >= '0' && b <= '9' {
            strLen = append(strLen, b)
        } else if b == ':' {
            break
        } else {
            return nil, errors.New("invalid character in string length")
        }
    }

    length, err := strconv.Atoi(string(strLen))
    if err != nil {
        return nil, err
    }

    buf:=make([]byte,length)
    _,err=io.ReadFull(reader,buf)
    if err != nil {
        return nil, err
    }
    return string(buf), nil
}


func decodeInteger(reader *bufio.Reader) (Bvalue, error) {
    if _, err := reader.ReadByte(); err != nil {
        return nil, err
    }
    intStr := []byte{}
    for {
        b, err := reader.ReadByte()
        if err != nil {
            return nil, err
        }
        if b == 'e' {
            break
        }
        intStr = append(intStr, b)
    }
    value, err := strconv.Atoi(string(intStr))
    if err != nil {
        return nil, err
    }
    return value, nil
}

func decodeList(reader *bufio.Reader) (Bvalue, error) {
    if _,err:= reader.ReadByte(); err != nil {
        return nil, err
    }
    list:=[]Bvalue{}
    for {
        b,err:=reader.Peek(1)
        if err != nil {
            return nil, err
        }
        if b[0] == 'e' {
            _, err := reader.ReadByte()
            if err != nil {
                return nil, err
            }
            break
        }
        value, err := decodeValue(reader)
        if err != nil {
            return nil, err
        }
        list = append(list, value)
    }
    return list, nil
}

func decodeDict(reader *bufio.Reader) (Bvalue, error) {
	if _, err := reader.ReadByte(); err != nil {
		return nil, err
	}
	dict := make(map[string]Bvalue)
	for {
		b, err := reader.Peek(1)
		if err != nil {
			return nil, err
		}
		if b[0] == 'e' {
			_, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			break
		}

		key, err := decodeString(reader)
		if err != nil {
			return nil, err
		}

		keyStr, ok := key.(string)
		if !ok {
			return nil, errors.New("key is not a string")
		}

		value, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		dict[keyStr] = value
	}
	return dict, nil
}

