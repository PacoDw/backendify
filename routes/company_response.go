package routes

import (
	"encoding/json"
	"time"

	"github.com/openlyinc/pointy"
)

// V1LegacyResponse represents the response for legacy providers.
type V1LegacyResponse struct {
	CN        string `json:"cn,omitempty"`
	CreatedOn string `json:"created_on,omitempty"`
	ClosedOn  string `json:"closed_on,omitempty"`
}

// V2LegacyResponse represents the response for legacy providers.
type V2LegacyResponse struct {
	CompanyName string `json:"company_name,omitempty"`
	TIN         string `json:"tin,omitempty"`
	DissolvedOn string `json:"dissolved_on,omitempty"`
}

// CompanyResponse represents the current reply message.
type CompanyResponse struct {
	ID          string     `json:"id,omitempty"`           // the company id requested by a customer
	Name        string     `json:"name,omitempty"`         // the company name, as returned by a backend
	Actived     *bool      `json:"actived,omitempty"`      // indicating if the company is still active according to the active_until date
	ActiveUntil *time.Time `json:"active_until,omitempty"` // RFC 3339 UTC date-time expressed as a string, optional.

	*V1LegacyResponse
	*V2LegacyResponse
}

// UnmarshalJSON helps to customize the data to set only the current reply message.
func (s *CompanyResponse) UnmarshalJSON(data []byte) error {
	// we created another type because it just get only the fields no the
	// methods to avoid get the UnmarshalJSON
	type Alias CompanyResponse

	aux := &struct {
		// passing all attributes without UnmarshalJSON method
		*Alias
	}{
		// pointing attribbutes with the same address
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.V1LegacyResponse != nil {
		s.Name = aux.CN

		if t, err := time.Parse(time.RFC3339, aux.ClosedOn); err == nil {
			s.ActiveUntil = &t

			// if the ActiveUntil time is more recently then the company is currently actived
			s.Actived = pointy.Bool(s.ActiveUntil.After(time.Now()))
		}

		s.V1LegacyResponse = nil
	}

	if aux.V2LegacyResponse != nil {
		s.Name = aux.CompanyName

		if t, err := time.Parse(time.RFC3339, aux.DissolvedOn); err == nil {
			s.ActiveUntil = &t

			// if the ActiveUntil time is more recently then the company is currently actived
			s.Actived = pointy.Bool(s.ActiveUntil.After(time.Now()))
		}

		s.V2LegacyResponse = nil
	}

	return nil
}

// ToJSON transforms the current struct to json.
func (s *CompanyResponse) ToJSON() []byte {
	if res, err := json.Marshal(s); err == nil {
		return res
	}

	return []byte{}
}
