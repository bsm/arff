package arff

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Writer instances can write ARFF data
type Writer struct {
	attrs int
	buf   *writeBuffer
	dst   io.Writer
	own   io.Closer
}

// Create creates a new relation file in fname and returns a writer
func Create(fname string, r *Relation) (*Writer, error) {
	file, err := os.Create(fname)
	if err != nil {
		return nil, err
	}

	w, err := NewWriter(file, r)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	w.own = file
	return w, nil
}

// NewWriter creates a new writer from a generic io.WriteCloser
func NewWriter(dst io.Writer, r *Relation) (*Writer, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	w := &Writer{
		attrs: len(r.Attributes),
		buf:   new(writeBuffer),
		dst:   dst,
	}

	if err := w.buf.WriteRelation(r.Name); err != nil {
		return nil, err
	}
	if err := w.buf.WriteAttributes(r.Attributes); err != nil {
		return nil, err
	}
	if _, err := w.buf.WriteString("@DATA\n"); err != nil {
		return nil, err
	}
	if err := w.buf.FlushTo(dst); err != nil {
		return nil, err
	}
	return w, nil
}

// Append appends a DataRow
func (w *Writer) Append(row *DataRow) error {
	if len(row.Values) != w.attrs {
		return errAttrMismatch
	}

	for i, v := range row.Values {
		if i != 0 {
			if err := w.buf.WriteByte(','); err != nil {
				return err
			}
		}
		if err := w.buf.WriteRowValue(v); err != nil {
			return err
		}
	}

	if row.Weight != 0 {
		if err := w.buf.WriteByte(','); err != nil {
			return err
		}
		if err := w.buf.WriteByte('{'); err != nil {
			return err
		}
		if err := w.buf.WriteFloat(row.Weight); err != nil {
			return err
		}
		if err := w.buf.WriteByte('}'); err != nil {
			return err
		}
	}
	if err := w.buf.WriteByte('\n'); err != nil {
		return err
	}
	return w.buf.FlushTo(w.dst)
}

// Close closes the underlying writer
func (w *Writer) Close() error {
	if w.own != nil {
		return w.own.Close()
	}
	return nil
}

// --------------------------------------------------------------------

type writeBuffer struct {
	bytes.Buffer
}

func (w *writeBuffer) WriteFloat(f float64) error {
	str := strings.TrimRight(strconv.FormatFloat(f, 'f', -1, 64), "0")
	_, err := w.WriteString(str)
	return err
}

func (w *writeBuffer) WriteInt(n int64) error {
	_, err := w.WriteString(strconv.FormatInt(n, 10))
	return err
}

func (w *writeBuffer) WriteUint(u uint64) error {
	_, err := w.WriteString(strconv.FormatUint(u, 10))
	return err
}

func (w *writeBuffer) WriteTime(t time.Time) error {
	_, err := w.WriteString(t.UTC().Format(iso8691DateFormat))
	return err
}

func (w *writeBuffer) WriteRelation(name string) error {
	if _, err := w.WriteString("@RELATION"); err != nil {
		return err
	} else if err := w.WriteByte(' '); err != nil {
		return err
	} else if err := w.WriteQuoted(name); err != nil {
		return err
	} else if err := w.WriteByte('\n'); err != nil {
		return err
	} else if err := w.WriteByte('\n'); err != nil {
		return err
	}
	return nil
}

func (w *writeBuffer) WriteAttributes(attrs []Attribute) error {
	for _, attr := range attrs {
		if err := w.WriteAttribute(&attr); err != nil {
			return err
		}
	}
	return w.WriteByte('\n')
}

func (w *writeBuffer) WriteAttribute(attr *Attribute) (err error) {
	if _, err = w.WriteString("@ATTRIBUTE"); err != nil {
		return
	} else if err = w.WriteByte(' '); err != nil {
		return
	} else if err = w.WriteQuoted(attr.Name); err != nil {
		return
	} else if err = w.WriteByte(' '); err != nil {
		return
	}

	switch attr.DataType {
	case DataTypeString:
		if _, err = w.WriteString("STRING"); err != nil {
			return
		}
	case DataTypeDate:
		if _, err = w.WriteString("DATE"); err != nil {
			return
		}
	case DataTypeNumeric:
		if _, err = w.WriteString("NUMERIC"); err != nil {
			return
		}
	case DataTypeNominal:
		if err = w.WriteByte('{'); err != nil {
			return
		}
		for i, v := range attr.NominalValues {
			if i != 0 {
				if err = w.WriteByte(','); err != nil {
					return
				}
			}
			if err = w.WriteQuoted(v); err != nil {
				return
			}
		}
		if err = w.WriteByte('}'); err != nil {
			return
		}
	default:
		return errInvalidAttrType
	}

	return w.WriteByte('\n')
}

func (w *writeBuffer) WriteRowValue(v interface{}) (err error) {
	if v == nil {
		err = w.WriteByte('?')
		return
	}

	switch vv := v.(type) {
	case float64:
		err = w.WriteFloat(vv)
	case int:
		err = w.WriteInt(int64(vv))
	case int8:
		err = w.WriteInt(int64(vv))
	case int16:
		err = w.WriteInt(int64(vv))
	case int32:
		err = w.WriteInt(int64(vv))
	case int64:
		err = w.WriteInt(int64(vv))
	case uint:
		err = w.WriteUint(uint64(vv))
	case uint8:
		err = w.WriteUint(uint64(vv))
	case uint16:
		err = w.WriteUint(uint64(vv))
	case uint32:
		err = w.WriteUint(uint64(vv))
	case uint64:
		err = w.WriteUint(uint64(vv))
	case time.Time:
		err = w.WriteTime(vv)
	case string:
		err = w.WriteQuoted(vv)
	default:
		err = fmt.Errorf("invalid value %v (%T)", v, v)
	}
	return
}

func (w *writeBuffer) WriteQuoted(s string) error {
	_, err := w.WriteString(quote(s))
	return err
}

func (w *writeBuffer) FlushTo(to io.Writer) error {
	_, err := w.WriteTo(to)
	w.Reset()
	return err
}
