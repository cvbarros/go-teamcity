package teamcity

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"
)

// responseDecoder decodes a http response body into the interface.
// To accomplish this, it first marshals the body into json and then unmarshal it into the interface.

type responseDecoder struct{}

func (responseDecoder) Decode(resp *http.Response, v interface{}) error {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if vs, ok := v.(*string); ok {
		*vs = string(bodyBytes)
		return nil
	}
	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(bodyBytes))
	}

	return nil
}
