package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// InvalidValueError is returned when a value cannot be transformed and should be omitted.
var InvalidValueError = errors.New("invalid value")

var (
	truthyBooleanValues = []string{"1", "t", "T", "true", "TRUE", "True"}
	falsyBooleanValues  = []string{"0", "f", "F", "false", "FALSE", "False"}
)

func main() {
	f, err := os.Open("./input.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	out := TransformData(data)
	fmt.Println(string(out))
}

func TransformData(data []byte) []byte {
	var in map[string]interface{}
	if err := json.Unmarshal(data, &in); err != nil {
		panic(err)
	}

	transformedJSON, err := transformMap(in)
	if err != nil {
		log.Fatalf("Error transforming JSON: %v", err)
	}

	data, err = json.Marshal(transformedJSON)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	return data
}

func transformMap(originalMap map[string]interface{}) (map[string]interface{}, error) {
	transformedMap := make(map[string]interface{})
	for key, value := range originalMap {
		if key == "" {
			continue
		}

		valueMap, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		v, err := transformValue(valueMap)
		if err != nil {
			continue
		}

		transformedMap[key] = v
	}

	return transformedMap, nil
}

// transformValue handles the transformation of individual elements based on their type.
func transformValue(valueMap map[string]interface{}) (interface{}, error) {
	for dataType, val := range valueMap {
		switch dataType {
		case "S":
			valStr, ok := val.(string)
			if !ok {
				continue
			}
			return transformString(valStr)
		case "N":
			valStr, ok := val.(string)
			if !ok {
				continue
			}
			return transformNumber(valStr)
		case "BOOL":
			valStr, ok := val.(string)
			if !ok {
				continue
			}
			return transformBool(valStr)
		case "NULL":
			valStr, ok := val.(string)
			if !ok {
				continue
			}
			return transformNull(valStr)
		case "L":
			return transformList(val)
		case "M":
			{
				v, ok := val.(map[string]interface{})
				if !ok {
					continue
				}

				return transformMap(v)
			}
		}
	}
	return nil, nil
}

// Transforms strings, sanitizing and converting RFC3339 timestamps to Unix.
func transformString(s string) (interface{}, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, InvalidValueError
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.Unix(), nil
	}
	return s, nil
}

// Transforms numbers, stripping leading zeros and converting to appropriate numeric types.
func transformNumber(s string) (interface{}, error) {
	s = strings.TrimSpace(s)
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return n, nil
	}

	return nil, InvalidValueError
}

// Transforms booleans from various string representations to bool.
func transformBool(s string) (interface{}, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if slices.Contains(truthyBooleanValues, s) {
		return true, nil
	} else if slices.Contains(falsyBooleanValues, s) {
		return false, nil
	}

	return false, InvalidValueError
}

// Transforms NULL representations to actual nil values in Go.
func transformNull(s string) (interface{}, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if slices.Contains(truthyBooleanValues, s) {
		return nil, nil
	}

	return nil, InvalidValueError
}

// transformList now calls transformValue for each element,
// allowing for nested transformations.
func transformList(val interface{}) (interface{}, error) {
	items, ok := val.([]interface{})
	if !ok {
		return items, InvalidValueError
	}

	transformedList := make([]interface{}, 0)
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		transformedItem, err := transformValue(itemMap)
		if err != nil {
			continue
		}
		if transformedItem != nil {
			transformedList = append(transformedList, transformedItem)
		}
	}

	return transformedList, nil
}
