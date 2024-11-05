package handlergen

import (
	"encoding/json"
	"github.com/mrbryside/go-generate/internal/mymap"
)

type HandlerTemplateData struct {
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	Request  *mymap.OrderedMap `json:"request"`
	Response *mymap.OrderedMap `json:"response"`
	Api      string            `json:"api"`
	Method   string            `json:"method"`
}

// UnmarshalJSON for HandlerTemplateData to handle OrderedMaps in Request and Response to maintain user spec json order
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

	var omRequest = mymap.NewOrderedMap()
	err := json.Unmarshal(aux.Request, omRequest)
	if err != nil {
		return err
	}
	var omResponse = mymap.NewOrderedMap()
	err = json.Unmarshal(aux.Response, omResponse)
	if err != nil {
		return err
	}

	h.Request = omRequest
	h.Response = omResponse

	return nil
}
