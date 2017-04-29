package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gocql/gocql"
)

// mapScanner is a workaround for buggy gocql.Iter.MapScan()
type mapScanner struct {
	Columns      []string
	Values       []interface{}
	isCollection []bool
}

func newMapScanner(columns []gocql.ColumnInfo) *mapScanner {
	ms := &mapScanner{}

	for _, column := range columns {
		switch v := column.TypeInfo.(type) {
		case gocql.TupleTypeInfo:
			for j := 0; j < len(v.Elems); j++ {
				ms.Values = append(ms.Values, v.Elems[j].New())
				ms.Columns = append(ms.Columns, fmt.Sprintf("%s[%d]", column.Name, j))
				ms.isCollection = append(ms.isCollection, false)
			}
		case gocql.NativeType:
			ms.Values = append(ms.Values, column.TypeInfo.New())
			ms.Columns = append(ms.Columns, column.Name)
			ms.isCollection = append(ms.isCollection, false)
		default:
			ms.Values = append(ms.Values, column.TypeInfo.New())
			ms.Columns = append(ms.Columns, column.Name)
			ms.isCollection = append(ms.isCollection, true)
		}
	}

	return ms
}

func (t *mapScanner) getHeaders() []string {
	return t.Columns
}

func (t *mapScanner) getStrings() []string {
	set := make([]string, len(t.Columns))
	for i, v := range t.Values {
		if v == nil {
			continue
		}

		if t.isCollection[i] {
			if b, err := json.Marshal(v); err == nil {
				set[i] = string(b)
			} else {
				loggr.Fatal(err)
			}

			continue
		}

		switch vt := v.(type) {
		case *gocql.UUID:
			set[i] = vt.String()
		case *[]uint8:
			set[i] = fmt.Sprintf("0x%X", *vt)
		default:
			switch nv := reflect.ValueOf(vt).Elem().Interface().(type) {
			default:
				set[i] = fmt.Sprint(nv)
			}
		}
	}
	return set
}
