package params

//
//Copyright 2019 Telenor Digital AS
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// configParameters is the internal (flattened) representation of the
// configuration parameters.
type configParameters struct {
	params []parameter
}

// newConfigParameters creates a configuration parameter set based on
// the supplied pointer to a configuration struct
func newConfigParameters(config interface{}) (*configParameters, error) {
	if config == nil {
		return nil, errors.New("config must be non-nil")
	}
	if reflect.TypeOf(config).Kind() != reflect.Ptr {
		return nil, errors.New("needs pointer to configuration")
	}
	ret := configParameters{
		params: make([]parameter, 0),
	}
	var err error
	if ret.params, err = readParameters("", config, ret.params); err != nil {
		return nil, err
	}

	return &ret, nil
}

func readParameters(prefix string, value interface{}, params []parameter) ([]parameter, error) {
	ct := reflect.TypeOf(value)
	vt := reflect.ValueOf(value)
	if ct.Kind() == reflect.Ptr {
		ct = ct.Elem()
		vt = vt.Elem()
	}
	if ct.Kind() != reflect.Struct {
		return nil, errors.New("needs struct type for configuration")
	}
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Field(i)
		// Skip private fields
		if unicode.IsLower(rune(field.Name[0])) {
			if _, ok := field.Tag.Lookup(tagName); ok {
				return nil, fmt.Errorf("field %s is unexported but has tag", field.Name)
			}
			continue
		}
		if !vt.Field(i).CanInterface() {
			return nil, fmt.Errorf("cannot set field %s", field.Name)
		}
		if vt.Field(i).Kind() == reflect.Struct {
			var err error
			params, err = readParameters(prefix+field.Name+".", vt.Field(i).Interface(), params)
			if err != nil {
				return nil, err
			}
			continue
		}
		param, err := newParameter(prefix, field, vt.Field(i).Interface(), vt.Field(i))
		if err != nil {
			return nil, err
		}
		if param == nil {
			continue
		}
		params = append(params, *param)
	}
	return params, nil
}

func (c *configParameters) getParameter(name string) *parameter {
	for i := range c.params {
		if strings.EqualFold(c.params[i].name, name) {
			return &c.params[i]
		}
	}
	return nil
}

// AssignValues assigns the current parameter config to the struct.
func (c *configParameters) AssignValues(config interface{}) error {
	for _, v := range c.params {
		fields := strings.Split(v.name, ".")
		f := reflect.ValueOf(config).Elem().FieldByName(fields[0])
		for n := 1; n < len(fields); n++ {
			f = f.FieldByName(fields[n])
		}
		if toInternalType(f.Interface()) != v.paramtype {
			return fmt.Errorf("invalid type for %s: %v (%T)", v.name, v.value, v.value)
		}
		if v.value == nil {
			continue
		}
		switch v.paramtype {
		case stringType:
			f.SetString(v.value.(string))
		case intType:
			f.SetInt(int64(v.value.(int)))
		case uintType:
			f.SetUint(uint64(v.value.(uint)))
		case boolType:
			f.SetBool(v.value.(bool))
		case durationType:
			f.SetInt(int64(v.value.(time.Duration)))
		case floatType:
			f.SetFloat(v.value.(float64))
		}
	}
	return nil
}

func (c *configParameters) Validate() error {
	for _, v := range c.params {
		if err := v.validate(); err != nil {
			return err
		}
	}
	return nil
}
