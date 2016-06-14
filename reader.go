package arff

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Reader struct {
	Relation
	scn *scanner
	src io.Reader
	own io.Closer
	row *DataRow
	err error
}

// Open reads a file at location
func Open(fname string) (*Reader, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	rd, err := NewReader(file)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	rd.own = file
	return rd, nil
}

// NewReader craetes an ARFF reader from any io.Reader
func NewReader(src io.Reader) (*Reader, error) {
	r := &Reader{
		src: src,
		scn: &scanner{Reader: bufio.NewReader(src)},
	}

	if err := r.parseHeader(); err != nil {
		return nil, r.wrapError(err)
	}
	return r, nil
}

// ReadAll reads all data rows at once, the equivalent of:
//     var rows []DataRow
//     for r.Next() {
//       rows = append(rows, r.Row())
//     }
//     err := r.Err()
func (r *Reader) ReadAll() ([]DataRow, error) {
	var rows []DataRow
	for r.Next() {
		rows = append(rows, *r.Row())
	}
	if err := r.Err(); err != nil {
		return nil, err
	}
	return rows, nil
}

// Close closes the reader
func (r *Reader) Close() error {
	if r.own != nil {
		return r.own.Close()
	}
	return nil
}

// Next returns true if can advance the row cursor
func (r *Reader) Next() bool {
	strs, err := r.scn.DataRow()
	if err != nil {
		r.markFailed(err)
		return false
	}

	nvals, nattrs := len(strs), len(r.Attributes)
	if nvals < nattrs {
		r.markFailed(errAttrMismatch)
		return false
	}

	var row DataRow
	for i, attr := range r.Attributes {
		v, err := attr.parse(strs[i])
		if err != nil {
			r.markFailed(err)
			return false
		}
		row.Values = append(row.Values, v)
	}

	// check if there is a weight
	if nvals > nattrs {
		weight := strs[nattrs]
		plast := len(weight) - 1
		if len(weight) < 2 || weight[0] != '{' || weight[plast] != '}' {
			r.markFailed(errInvalidWeight)
			return false
		}

		num, err := strconv.ParseFloat(weight[1:plast], 64)
		if err != nil || num < 0 {
			r.markFailed(errInvalidWeight)
			return false
		}
		row.Weight = num
	}

	r.row = &row
	return true
}

// Row returns the current DataRow
func (r *Reader) Row() *DataRow { return r.row }

// Err returns an error if any
func (r *Reader) Err() error {
	return r.err
}

func (r *Reader) parseHeader() error {
	for {
		fields, err := r.scn.HeaderFields()
		if err != nil {
			return err
		}

		if len(fields) == 0 {
			continue
		}
		switch strings.ToUpper(fields[0]) {
		case "@RELATION":
			if len(fields) < 2 {
				return errMissingRelName
			}
			r.Relation.Name = unquote(fields[1])
		case "@ATTRIBUTE":
			if len(fields) < 2 {
				return errMissingAttrName
			} else if len(fields) < 3 {
				return errMissingAttrType
			}

			attr := Attribute{
				Name: unquote(fields[1]),
			}
			switch strings.ToUpper(fields[2]) {
			case "NUMERIC", "REAL", "INTEGER":
				attr.DataType = DataTypeNumeric
			case "STRING":
				attr.DataType = DataTypeString
			case "DATE":
				attr.DataType = DataTypeDate
			case "NOMINAL":
				attr.DataType = DataTypeNominal
				attr.NominalValues = unquoteAll(fields[3:])
			default:
				return errInvalidAttrType
			}
			r.Relation.Attributes = append(r.Relation.Attributes, attr)
		case "@DATA":
			return nil
		default:
			return errBadSyntax
		}
	}
}

func (r *Reader) markFailed(err error) {
	if err != io.EOF {
		r.err = r.wrapError(err)
	}
	r.row = nil
}

func (r *Reader) wrapError(err error) error {
	if err != nil {
		return fmt.Errorf("LINE %d: %s", r.scn.Lineno, err.Error())
	}
	return nil
}

// --------------------------------------------------------------------

type scanner struct {
	*bufio.Reader
	Lineno int
}

func (s *scanner) DataRow() ([]string, error) {
	for {
		s.Lineno++

		line, err := s.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		if vv := scanCSV(line[:len(line)-1]); len(vv) != 0 {
			return vv, nil
		}
	}
}

func (s *scanner) HeaderFields() ([]string, error) {
	s.Lineno++

	line, err := s.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	min := 0
	prv := rune(0)

	var fields []string
	var inQuote, inBracket bool

MainLoop:
	for i := 0; i < len(line); {
		r, size := utf8.DecodeRune(line[i:])

		switch r {
		case '\'':
			if !inQuote {
				inQuote = true
			} else if prv != '\\' {
				inQuote = false
				if !inBracket {
					fields = append(fields, string(line[min:i+size]))
					min = i + size
				}
			}
		case '{':
			if !inQuote {
				inBracket = true
			}
		case '}':
			if !inQuote {
				inBracket = false
				if min < i+size {
					fields = append(fields, "NOMINAL")
					fields = append(fields, scanCSV(line[min+1:i])...)
					min = i + size
				}
			}
		case ' ', '\t':
			if !inQuote && !inBracket {
				if min < i {
					fields = append(fields, string(line[min:i]))
				}
				min = i + size
			}
		case '%':
			if !inQuote && !inBracket {
				line = line[:i]
				break MainLoop
			}
		}

		prv = r
		i += size
	}

	if min < len(line)-1 {
		fields = append(fields, string(line[min:len(line)-1]))
	}

	return fields, nil
}

func scanCSV(line []byte) []string {
	min := 0
	prv := rune(0)

	var fields []string
	var inQuote bool

MainLoop:
	for i := 0; i < len(line); {
		r, size := utf8.DecodeRune(line[i:])

		switch r {
		case '\'':
			if !inQuote {
				inQuote = true
			} else if prv != '\\' {
				inQuote = false

				fields = append(fields, string(line[min:i+size]))
				min = i + size
			}
		case ',':
			if !inQuote {
				if min < i {
					fields = append(fields, string(line[min:i]))
				}
				min = i + size
			}
		case ' ', '\t':
			if !inQuote {
				min = i + size
			}
		case '%':
			if !inQuote {
				line = line[:i]
				break MainLoop
			}
		}

		prv = r
		i += size
	}
	if min < len(line) {
		fields = append(fields, string(line[min:len(line)]))
	}

	return fields
}
