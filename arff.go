package arff

import (
	"fmt"
	"strconv"
	"time"
)

type DataType uint8

const (
	DataTypeNumeric DataType = iota
	DataTypeString
	DataTypeDate
	DataTypeNominal
)

// Relation contains meta-data and attribute definition
type Relation struct {
	// The relation name
	Name string

	// The attributes
	Attributes []Attribute
}

// AddAttribute stores an attribute, avoiding duplicates.
// Include nominalVals for nominal data-types
func (r *Relation) AddAttribute(name string, dataType DataType, nominalVals []string) error {
	for _, attr := range r.Attributes {
		if attr.Name == name {
			return errAttrRedefined
		}
	}

	r.Attributes = append(r.Attributes, Attribute{
		Name:          name,
		DataType:      dataType,
		NominalValues: nominalVals,
	})
	return nil
}

func (r *Relation) validate() error {
	if r.Name == "" {
		return errMissingRelName
	}
	for _, attr := range r.Attributes {
		if err := attr.validate(); err != nil {
			return err
		}
	}
	return nil
}

// Attribute is an attribute of the dataset
type Attribute struct {
	// The attribute name
	Name string

	// DataType represent the attribute data-type
	DataType DataType

	// NominalValues are only populated for nominal types
	NominalValues []string
}

func (a *Attribute) validate() error {
	if a.Name == "" {
		return errMissingAttrName
	}
	return nil
}

func (a *Attribute) parse(s string) (interface{}, error) {
	if s == "?" {
		return nil, nil
	}

	switch a.DataType {
	case DataTypeNumeric:
		num, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("value '%s' is not numeric", s)
		}
		return num, nil
	case DataTypeDate:
		dt, err := time.ParseInLocation(iso8691DateFormat, s, utc)
		if err != nil {
			return nil, fmt.Errorf("value '%s' is not an ISO8601 date", s)
		}
		return dt, nil
	}
	return unquote(s), nil
}

// DataRow represents a parsed data row
type DataRow struct {
	Values []interface{}
	Weight float64
}

// --------------------------------------------------------------------

const iso8691DateFormat = "2006-01-02T15:04:05"

var utc *time.Location

func init() {
	var err error

	utc, err = time.LoadLocation("UTC")
	if err != nil {
		panic("unable to load UTC time-zone information: " + err.Error())
	}
}

type constError string

// Error implements error interface
func (e constError) Error() string { return string(e) }

const (
	errBadSyntax       constError = "bad syntax"
	errMissingAttrName constError = "missing attribute name"
	errMissingAttrType constError = "missing data-type"
	errInvalidAttrType constError = "invalid data-type"
	errAttrRedefined   constError = "redefined attribute"
	errAttrMismatch    constError = "attribute mismatch"
	errMissingRelName  constError = "missing relation name"
	errInvalidWeight   constError = "invalid weight definition"
)
