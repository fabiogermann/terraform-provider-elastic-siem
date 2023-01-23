package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func contains[K comparable](s []K, item K) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

func moveToFirstPositionOfSlice[K comparable](slice []K, item K) []K {
	if len(slice) == 0 || (slice)[0] == item {
		return slice
	}
	if (slice)[len(slice)-1] == item {
		slice = append([]K{item}, (slice)[:len(slice)-1]...)
		return slice
	}
	for p, x := range slice {
		if x == item {
			slice = append([]K{item}, append((slice)[:p], (slice)[p+1:]...)...)
			break
		}
	}
	return slice
}

func ifThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func Convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = Convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = Convert(v)
		}
	}
	return i
}

func ObjectFronJSON(jsonString string, result interface{}) error {
	return json.Unmarshal([]byte(jsonString), &result)
}

func Sha256String(name string) string {
	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:])
}
