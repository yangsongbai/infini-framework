/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/cihub/seelog"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

// BytesToUint64 convert bytes to type uint64
func BytesToUint64(b []byte) (v uint64) {
	length := uint(len(b))
	for i := uint(0); i < length-1; i++ {
		v += uint64(b[i])
		v <<= 8
	}
	v += uint64(b[length-1])
	return
}

// BytesToUint32 convert bytes to uint32
func BytesToUint32(b []byte) (v uint32) {
	length := uint(len(b))
	for i := uint(0); i < length-1; i++ {
		v += uint32(b[i])
		v <<= 8
	}
	v += uint32(b[length-1])
	return
}

// Uint64toBytes convert uint64 to bytes
func Uint64toBytes(b []byte, v uint64) {
	for i := uint(0); i < 8; i++ {
		b[7-i] = byte(v >> (i * 8))
	}
}

// Uint32toBytes convert uint32 to bytes, max uint: 4294967295
func Uint32toBytes(b []byte, v uint32) {
	for i := uint(0); i < 4; i++ {
		b[3-i] = byte(v >> (i * 8))
	}
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// DeepCopy return a deep copied object
func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}

/** https://github.com/cloudfoundry/bytefmt/blob/master/bytes.go start
https://github.com/cloudfoundry/bytefmt/blob/master/LICENSE
Apache License  Version 2.0, January 2004
 **/

// ByteSize unit definition
const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

var bytesPattern *regexp.Regexp = regexp.MustCompile(`(?i)^(-?\d+)([KMGT]B?|B)$`)

var errInvalidByteQuantity = errors.New("Byte quantity must be a positive integer with a unit of measurement like M, MB, G, or GB")

// ByteSize returns a human-readable byte string of the form 10M, 12.5K, and so forth.  The following units are available:
//	T: Terabyte
//	G: Gigabyte
//	M: Megabyte
//	K: Kilobyte
//	B: Byte
// The unit that results in the smallest number greater than or equal to 1 is always chosen.
func ByteSize(bytes uint64) string {
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}

// ToMegabytes parses a string formatted by ByteSize as megabytes.
func ToMegabytes(s string) (uint64, error) {
	bytes, err := ToBytes(s)
	if err != nil {
		return 0, err
	}

	return bytes / MEGABYTE, nil
}

// ToBytes parses a string formatted by ByteSize as bytes.
func ToBytes(s string) (uint64, error) {
	parts := bytesPattern.FindStringSubmatch(strings.TrimSpace(s))
	if len(parts) < 3 {
		return 0, errInvalidByteQuantity
	}

	value, err := strconv.ParseUint(parts[1], 10, 0)
	if err != nil || value < 1 {
		return 0, errInvalidByteQuantity
	}

	var bytes uint64
	unit := strings.ToUpper(parts[2])
	switch unit[:1] {
	case "T":
		bytes = value * TERABYTE
	case "G":
		bytes = value * GIGABYTE
	case "M":
		bytes = value * MEGABYTE
	case "K":
		bytes = value * KILOBYTE
	case "B":
		bytes = value * BYTE
	}

	return bytes, nil
}

/** https://github.com/cloudfoundry/bytefmt/blob/master/bytes.go end **/

func BytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// ToLowercase convert string bytes to lowercase
func ToLowercase(str []byte) []byte {
	for i, s := range str {
		if s > 64 && s < 91 {
			str[i] = s + 32
		}
	}
	return str
}

// ToUppercase convert string bytes to uppercase
func ToUppercase(str []byte) []byte {
	for i, s := range str {
		if s > 96 && s < 123 {
			str[i] = s - 32
		}
	}
	return str
}

//TODO optimize performance
//ReplaceByte simply replace old bytes to new bytes, the two bytes should have same length
func ReplaceByte(str []byte, old, new []byte) []byte {
	return []byte(strings.Replace(string(str), string(old), string(new), -1))
}

//ToJSONBytes convert interface to json with byte array
func ToJSONBytes(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

//FromJSONBytes simply do json unmarshal
func FromJSONBytes(b []byte, v interface{}) {
	err := json.Unmarshal(b, v)
	if err != nil {
		panic(err)
	}
}

func EncodeToBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetBytes(key interface{}) []byte {
	return []byte(fmt.Sprintf("%v", key.(interface{})))
}

func GetSplitFunc(split []byte) func(data []byte, atEOF bool) (advance int, token []byte, err error) {

	sLen := len(split)

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		// Return nothing if at end of file and no data passed
		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		// Find next separator and return token
		if i := bytes.Index(data, split); i >= 0 {
			return i + sLen, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return dataLen, data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}

func ExtractFieldFromJson(data *[]byte, fieldStartWith []byte, fieldEndWith []byte, leftMustContain []byte) (bool, []byte) {
	return ExtractFieldFromJsonOrder(data, fieldStartWith, fieldEndWith, leftMustContain, false)
}

func ExtractFieldFromJsonOrder(data *[]byte, fieldStartWith []byte, fieldEndWith []byte, leftMustContain []byte, reverse bool) (bool, []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(*data))
	scanner.Split(GetSplitFunc(fieldEndWith))
	var str []byte
	for scanner.Scan() {
		text := scanner.Bytes()
		if bytes.Contains(text, leftMustContain) {
			str = text
			break
		}
	}

	if len(str) > 0 {
		var offset int
		if reverse {
			offset = bytes.LastIndex(str, fieldStartWith)
		} else {
			offset = bytes.Index(str, fieldStartWith)
		}

		if offset > 0 && offset < len(str) {
			newStr := str[offset+len(fieldStartWith) : len(str)]
			return true, newStr
		}
	} else {
		log.Trace("input data doesn't contain the split bytes")
	}
	return false, nil
}

func ProcessJsonData(data *[]byte, fieldStartWith []byte, fieldEndWith []byte, leftMustContain []byte, reverse bool, handler func(start,end int)) bool {
	scanner := bufio.NewScanner(bytes.NewReader(*data))
	scanner.Split(GetSplitFunc(fieldEndWith))
	var str []byte
	for scanner.Scan() {
		text := scanner.Bytes()
		if bytes.Contains(text, leftMustContain) {
			str = text
			break
		}
	}

	if len(str) > 0 {
		var offset int
		if reverse {
			offset = bytes.LastIndex(str, fieldStartWith)
		} else {
			offset = bytes.Index(str, fieldStartWith)
		}

		if offset > 0 && offset < len(str) {
			handler(offset+len(fieldStartWith),len(str))
			return true
		}
	} else {
		log.Trace("input data doesn't contain the split bytes")
	}
	return false
}

func IsBytesEndingWith(data *[]byte, ending []byte) bool {
	return IsBytesEndingWithOrder(data, ending, false)
}

func IsBytesEndingWithOrder(data *[]byte, ending []byte, reverse bool) bool {
	var offset int
	if reverse {
		offset = bytes.LastIndex(*data, ending)
	} else {
		offset = bytes.Index(*data, ending)
	}
	return len(*data)-offset <= len(ending)
}


func BytesSearchValue(data,startTerm,endTerm,searchTrim []byte) bool  {
	index:=bytes.Index(data,startTerm)
	leftData:=data[index+len(startTerm):]
	endIndex:=bytes.Index(leftData,endTerm)
	lastTerm:=leftData[0:endIndex]

	if bytes.Contains(lastTerm,searchTrim){
		return true
	}
	return false
}