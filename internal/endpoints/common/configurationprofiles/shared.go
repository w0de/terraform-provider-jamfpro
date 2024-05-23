// common/configurationprofiles/shared.go contains the shared functions to process configuration profiles.
package configurationprofiles

import (
	"bytes"
	"log"
	"sort"

	"howett.net/plist"
)

// sortKeys recursively sorts the keys of a nested map in alphabetical order.
func sortKeys(data map[string]interface{}) map[string]interface{} {
	sortedData := make(map[string]interface{})
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	log.Printf("Unsorted keys: %v\n", keys)
	sort.Strings(keys)
	log.Printf("Sorted keys: %v\n", keys)
	for _, k := range keys {
		log.Printf("Processing key: %s\n", k)
		switch v := data[k].(type) {
		case map[string]interface{}:
			log.Printf("Key %s is a nested map, sorting nested keys...\n", k)
			sortedData[k] = sortKeys(v)
		case []interface{}:
			log.Printf("Key %s is an array, processing items...\n", k)
			sortedArray := make([]interface{}, len(v))
			for i, item := range v {
				log.Printf("Processing item %d of array %s\n", i, k)
				if nestedMap, ok := item.(map[string]interface{}); ok {
					sortedArray[i] = sortKeys(nestedMap)
				} else {
					sortedArray[i] = item
				}
			}
			sortedData[k] = sortedArray
		default:
			sortedData[k] = v
		}
	}
	return sortedData
}

// EncodePlist encodes a cleaned map back to plist XML format
func EncodePlist(cleanedData map[string]interface{}) (string, error) {
	log.Printf("Encoding plist data: %v\n", cleanedData)
	var buffer bytes.Buffer
	encoder := plist.NewEncoder(&buffer)
	encoder.Indent("\t") // Optional: for pretty-printing the XML
	if err := encoder.Encode(cleanedData); err != nil {
		log.Printf("Error encoding plist data: %v\n", err)
		return "", err
	}
	return buffer.String(), nil
}