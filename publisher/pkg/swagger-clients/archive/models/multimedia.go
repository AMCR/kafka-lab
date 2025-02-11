// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Multimedia multimedia
//
// swagger:model Multimedia
type Multimedia struct {

	// caption
	Caption string `json:"caption,omitempty"`

	// credit
	Credit string `json:"credit,omitempty"`

	// crop name
	CropName string `json:"crop_name,omitempty"`

	// height
	Height int64 `json:"height,omitempty"`

	// legacy
	Legacy *MultimediaLegacy `json:"legacy,omitempty"`

	// rank
	Rank int64 `json:"rank,omitempty"`

	// subtype
	Subtype string `json:"subtype,omitempty"`

	// type
	Type string `json:"type,omitempty"`

	// url
	URL string `json:"url,omitempty"`

	// width
	Width int64 `json:"width,omitempty"`
}

// Validate validates this multimedia
func (m *Multimedia) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLegacy(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Multimedia) validateLegacy(formats strfmt.Registry) error {
	if swag.IsZero(m.Legacy) { // not required
		return nil
	}

	if m.Legacy != nil {
		if err := m.Legacy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("legacy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("legacy")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this multimedia based on the context it is used
func (m *Multimedia) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLegacy(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Multimedia) contextValidateLegacy(ctx context.Context, formats strfmt.Registry) error {

	if m.Legacy != nil {

		if swag.IsZero(m.Legacy) { // not required
			return nil
		}

		if err := m.Legacy.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("legacy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("legacy")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Multimedia) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Multimedia) UnmarshalBinary(b []byte) error {
	var res Multimedia
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// MultimediaLegacy multimedia legacy
//
// swagger:model MultimediaLegacy
type MultimediaLegacy struct {

	// xlarge
	Xlarge string `json:"xlarge,omitempty"`

	// xlargeheight
	Xlargeheight int64 `json:"xlargeheight,omitempty"`

	// xlargewidth
	Xlargewidth int64 `json:"xlargewidth,omitempty"`
}

// Validate validates this multimedia legacy
func (m *MultimediaLegacy) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this multimedia legacy based on context it is used
func (m *MultimediaLegacy) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MultimediaLegacy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MultimediaLegacy) UnmarshalBinary(b []byte) error {
	var res MultimediaLegacy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
