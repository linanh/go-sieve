package sieve

import (
	"errors"
	"fmt"
	"strings"
)

//ParamsType Params type
type ParamsType map[string]string

func (params ParamsType) String() string {
	strs := make([]string, 4)
	for k, v := range params {
		switch k {
		case "days":
			strs[0] = fmt.Sprintf(":days %s ", v)
		case "addresses":
			strs[1] = fmt.Sprintf(":addresses %s ", v)
		case "subject":
			strs[2] = fmt.Sprintf(`:subject %s `, v)
		case "text":
			strs[3] = fmt.Sprintf(`"%s"`, v)
		}
	}
	return strings.Join(strs, "")
}

//Scan parse Param
func (params *ParamsType) Scan(data []byte) (int, error) {

	pos := 0
	l := len(data)

	if l == 0 {
		return 0, errors.New("Unexpected zero length number")
	}

	keyNum := 0
	valNum := 0
	isKeyStart := false
	keyBytes := []byte{}
	isValueStart := false
	isValueNumber := false
	valueBytes := []byte{}

	for pos < l {

		isEnd := data[pos] == ' ' || data[pos] == '\r' || data[pos] == '\n' || data[pos] == '\t'
		if isEnd && !isKeyStart && !isValueStart {
			pos++
			continue
		}

		if pos < l-1 && !isKeyStart && !isValueStart && data[pos] == ':' && data[pos+1] != ' ' {
			keyNum++
			isKeyStart = true
			keyBytes = []byte{}
			pos++
			continue
		}

		if isKeyStart {
			if isEnd {
				(*params)[string(keyBytes)] = ""
				isKeyStart = false
				pos++
				continue
			}

			keyBytes = append(keyBytes, data[pos])
			pos++
			continue
		}

		if !isKeyStart && !isValueStart {
			valNum++
			if keyNum < valNum {
				keyNum++
				keyBytes = []byte("text")
			}
			isValueStart = true
			valueBytes = []byte{}
		}

		if isValueStart {
			vlen := len(valueBytes)
			if vlen == 0 && isEnd {
				pos++
				continue
			}

			//valueBytes type is number
			if vlen == 0 && data[pos] >= '0' && data[pos] <= '9' {
				isValueNumber = true
			}

			//Number valueBytes end
			if isValueNumber && isEnd {
				isValueNumber = false
				(*params)[string(keyBytes)] = string(valueBytes)
				isValueStart = false
				pos++
				continue
			}

			//multiple string valueBytes end
			if vlen > 0 && valueBytes[0] == '[' && data[pos] == ']' && data[pos-1] != '\\' {
				valueBytes = append(valueBytes, data[pos])
				(*params)[string(keyBytes)] = string(valueBytes)
				isValueStart = false
				pos++
				continue
			}

			//single string valueBytes end
			if vlen > 0 && valueBytes[0] == '"' && data[pos] == '"' && data[pos-1] != '\\' {
				valueBytes = append(valueBytes, data[pos])
				(*params)[string(keyBytes)] = string(valueBytes)
				isValueStart = false
				pos++
				continue
			}

			//single string with dot valueBytes end
			if vlen > 0 && data[pos] == '.' && pos < l-1 && (data[pos+1] == '\r' || data[pos+1] == '\n') && (data[pos-1] == '\r' || data[pos-1] == '\n') {
				(*params)[string(keyBytes)] = strings.TrimPrefix(strings.TrimPrefix(string(valueBytes), "text:\r\n"), "text:\n")
				isValueStart = false
				pos += 2
				break
			}

			valueBytes = append(valueBytes, data[pos])
			pos++
			continue

		}

	}

	return pos, nil
}
