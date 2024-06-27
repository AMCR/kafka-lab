// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Person person
//
// swagger:model Person
type Person struct {

	// firstname
	Firstname string `json:"firstname,omitempty"`

	// lastname
	Lastname string `json:"lastname,omitempty"`

	// middlename
	Middlename string `json:"middlename,omitempty"`

	// organization
	Organization string `json:"organization,omitempty"`

	// qualifier
	Qualifier string `json:"qualifier,omitempty"`

	// rank
	Rank int64 `json:"rank,omitempty"`

	// role
	Role string `json:"role,omitempty"`

	// title
	Title string `json:"title,omitempty"`
}

// Validate validates this person
func (m *Person) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this person based on context it is used
func (m *Person) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Person) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Person) UnmarshalBinary(b []byte) error {
	var res Person
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
