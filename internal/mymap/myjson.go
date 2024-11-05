package mymap

import (
	"encoding/json"
)

// OrderedMap holds the map data and preserves key order
//type OrderedMap struct {
//	Data  map[string]interface{}
//	Order []string
//}

// HandlerTemplateData holds request and response as OrderedMaps
type HandlerTemplateData struct {
	Type     string      `json:"type"`
	Name     string      `json:"name"`
	Api      string      `json:"api"`
	Method   string      `json:"method"`
	Header   string      `json:"header"`
	Request  *OrderedMap `json:"request"`
	Response *OrderedMap `json:"response"`
}

// UnmarshalJSON for HandlerTemplateData to handle OrderedMaps
func (h *HandlerTemplateData) UnmarshalJSON(data []byte) error {
	type Alias HandlerTemplateData
	aux := &struct {
		Request  json.RawMessage `json:"request"`
		Response json.RawMessage `json:"response"`
		*Alias
	}{
		Alias: (*Alias)(h),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var omRequest = NewOrderedMap()
	err := json.Unmarshal(aux.Request, omRequest)
	if err != nil {
		return err
	}
	var omResponse = NewOrderedMap()
	err = json.Unmarshal(aux.Request, omResponse)
	if err != nil {
		return err
	}

	h.Request = omRequest
	h.Response = omResponse

	return nil
}

//// Helper to unmarshal JSON into an OrderedMap while preserving order
//func unmarshalOrderedMap(data json.RawMessage) *OrderedMap {
//	orderedMap := &OrderedMap{
//		Data:  make(map[string]interface{}),
//		Order: []string{},
//	}
//
//	dec := json.NewDecoder(bytes.NewReader(data))
//	if _, err := dec.Token(); err != nil { // Read the opening `{`
//		return nil
//	}
//
//	for dec.More() {
//		// Read key
//		key, err := dec.Token()
//		if err != nil {
//			return nil
//		}
//		keyStr := key.(string)
//
//		// Read value
//		var value interface{}
//		if err := dec.Decode(&value); err != nil {
//			return nil
//		}
//
//		// Handle nested map
//		if nestedData, ok := value.(map[string]interface{}); ok {
//			nestedJSON, _ := json.Marshal(nestedData)
//			value = unmarshalOrderedMap(nestedJSON)
//		}
//
//		// Store key-value pair in order
//		orderedMap.Data[keyStr] = value
//		orderedMap.Order = append(orderedMap.Order, keyStr) // Maintain original order
//	}
//
//	return orderedMap
//}

// Recursive helper to build an ordered map[string]interface{} from OrderedMap
//func buildOrderedMap(data *OrderedMap) map[string]interface{} {
//	ordered := make(map[string]interface{})
//	// Iterate in the order of data.Order
//	for _, key := range data.Order {
//		value := data.Data[key]
//		if nestedMap, ok := value.(*OrderedMap); ok {
//			ordered[key] = buildOrderedMap(nestedMap) // Recursive call for nested OrderedMap
//		} else {
//			ordered[key] = value
//		}
//	}
//	return ordered
//}

// GenerateMyJson converts HandlerTemplateData's ordered request and response maps into ordered map[string]interface{}
//func GenerateMyJson(handler HandlerTemplateData) (map[string]interface{}, map[string]interface{}, error) {
//	request := buildOrderedMap(handler.Request)
//	response := buildOrderedMap(handler.Response)
//	return request, response, nil
//}
