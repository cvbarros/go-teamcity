package teamcity

import (
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
)

// Href href
// swagger:model href
type Href struct {

	// href
	Href string `json:"href,omitempty" xml:"href"`
}

// Validate validates this href
func (m *Href) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
