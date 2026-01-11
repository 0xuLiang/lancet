package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func Marshal(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	var sliceValue reflect.Value
	var sliceType reflect.Type

	switch {
	case !rv.IsValid():
		return nil, errors.New("v is nil")
	case rv.Kind() == reflect.Ptr && rv.IsNil():
		return nil, errors.New("v is nil")
	case rv.Kind() == reflect.Slice:
		sliceValue = rv
		sliceType = rv.Type().Elem()
	case rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Slice:
		sliceValue = rv.Elem()
		sliceType = rv.Elem().Type().Elem()
	case rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct:
		sliceValue = reflect.New(reflect.SliceOf(rv.Elem().Type())).Elem()
		sliceValue = reflect.Append(sliceValue, rv.Elem())
		sliceType = rv.Elem().Type()
	case rv.Kind() == reflect.Struct:
		sliceValue = reflect.New(reflect.SliceOf(rv.Type())).Elem()
		sliceValue = reflect.Append(sliceValue, rv)
		sliceType = rv.Type()
	default:
		return nil, errors.New("v must be a struct, a struct pointer or a slice of struct")
	}
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
	}
	if sliceType.Kind() != reflect.Struct {
		return nil, errors.New("element must be a struct")
	}

	var headers []string
	for i := 0; i < sliceType.NumField(); i++ {
		field := sliceType.Field(i)
		tag := field.Tag.Get("csv")
		if tag == "" {
			tag = field.Name
		}
		headers = append(headers, tag)
	}

	var records [][]string
	records = append(records, headers)
	for i := 0; i < sliceValue.Len(); i++ {
		var record []string
		rvElem := sliceValue.Index(i)
		if rvElem.Kind() == reflect.Ptr {
			if rvElem.IsNil() {
				return nil, errors.New("slice element is nil")
			}
			rvElem = rvElem.Elem()
		}
		for j := 0; j < rvElem.NumField(); j++ {
			field := rvElem.Field(j)
			var value string
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				value = strconv.FormatInt(field.Int(), 10)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				value = strconv.FormatUint(field.Uint(), 10)
			case reflect.Float32, reflect.Float64:
				value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
			case reflect.Bool:
				value = strconv.FormatBool(field.Bool())
			case reflect.String:
				value = field.String()
			case reflect.Struct:
				if field.Type() == reflect.TypeOf(time.Time{}) {
					t := field.Interface().(time.Time)
					b, err := t.MarshalText()
					if err != nil {
						return nil, err
					}
					value = string(b)
				} else {
					return nil, fmt.Errorf("unsupported struct type: %s", field.Type())
				}
			default:
				return nil, fmt.Errorf("unsupported field type: %s", field.Type())
			}
			record = append(record, value)
		}
		records = append(records, record)
	}

	b := &bytes.Buffer{}
	writer := csv.NewWriter(b)
	if err := writer.WriteAll(records); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	var sliceValue reflect.Value
	var sliceType reflect.Type
	var singleStruct bool

	if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Slice {
		sliceValue = rv.Elem()
		sliceType = rv.Elem().Type().Elem()
	} else if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct {
		sliceValue = reflect.New(reflect.SliceOf(rv.Type().Elem())).Elem()
		sliceType = rv.Elem().Type()
		singleStruct = true
	} else {
		return errors.New("v must be a pointer to a struct or a slice of struct")
	}

	var isPtr bool
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
		isPtr = true
	}
	if sliceType.Kind() != reflect.Struct {
		return errors.New("element must be a struct")
	}

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return errors.New("no records found")
	}

	headers := records[0]
	fieldIndex := make(map[string]int)
	for i := 0; i < sliceType.NumField(); i++ {
		field := sliceType.Field(i)
		tag := field.Tag.Get("csv")
		if tag == "" {
			tag = field.Name
		}
		for _, header := range headers {
			if tag == header {
				fieldIndex[header] = i
			}
		}
	}

	for _, record := range records[1:] {
		newValue := reflect.New(sliceType)
		limit := len(headers)
		if len(record) < limit {
			limit = len(record)
		}
		for i := 0; i < limit; i++ {
			value := record[i]
			header := headers[i]
			if fi, ok := fieldIndex[header]; ok {
				field := newValue.Elem().Field(fi)
				switch field.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					var intValue int64
					if value != "" {
						if intValue, err = strconv.ParseInt(value, 10, 64); err != nil {
							return err
						}
					}
					field.SetInt(intValue)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
					var uintValue uint64
					if value != "" {
						if uintValue, err = strconv.ParseUint(value, 10, 64); err != nil {
							return err
						}
					}
					field.SetUint(uintValue)
				case reflect.Float32, reflect.Float64:
					var floatValue float64
					if value != "" {
						if floatValue, err = strconv.ParseFloat(value, 64); err != nil {
							return err
						}
					}
					field.SetFloat(floatValue)
				case reflect.Bool:
					var boolValue bool
					if value != "" {
						if boolValue, err = strconv.ParseBool(value); err != nil {
							return err
						}
					}
					field.SetBool(boolValue)
				case reflect.String:
					field.SetString(value)
				case reflect.Struct:
					if field.Type() == reflect.TypeOf(time.Time{}) {
						var t time.Time
						if err := t.UnmarshalText([]byte(value)); err != nil {
							return err
						}
						field.Set(reflect.ValueOf(t))
					} else {
						return fmt.Errorf("unsupported struct type: %s", field.Type())
					}
				default:
					return fmt.Errorf("unsupported field type: %s", field.Type())
				}
			}
		}
		if isPtr {
			sliceValue.Set(reflect.Append(sliceValue, newValue))
		} else {
			sliceValue.Set(reflect.Append(sliceValue, newValue.Elem()))
		}
	}

	if singleStruct {
		if sliceValue.Len() == 0 {
			return errors.New("no data rows found")
		}
		rv.Elem().Set(sliceValue.Index(0))
	}

	return nil
}
