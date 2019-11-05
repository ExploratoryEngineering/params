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
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// TagName is the name of the struct tags used by the package
const tagName = "param"

type internalType int

const (
	stringType internalType = iota
	intType
	uintType
	boolType
	durationType
	floatType
	invalidType
)

// Parameter is a single config parameter
type parameter struct {
	name         string
	paramtype    internalType
	description  string
	defaultValue string
	value        interface{}
	minvalue     string
	maxvalue     string
	file         bool
	options      []string
	required     bool
	isSet        bool
}

// hyphenName converts name into a lowercase string with hyphens. Hyphens are
// inserted when case transitions from LUL (as in "NameName"), UUL (as in "TLAName")
// and LUU (as in "NameTLA")
func (p *parameter) hyphenName() string {
	var ret []rune
	prevChar := 'X'
	prevPrevChar := prevChar
	changes := 0
	for i, ch := range strings.Replace(p.name, ".", "-", -1) {
		// first char is always included
		if i < 2 {
			ret = append(ret, ch)
			prevPrevChar = prevChar
			prevChar = ch
			continue
		}
		caseChange := (unicode.IsUpper(prevPrevChar) && unicode.IsUpper(prevChar) && unicode.IsLower(ch)) ||
			(unicode.IsLower(prevPrevChar) && unicode.IsUpper(prevChar) && unicode.IsLower(ch)) ||
			(unicode.IsLower(prevPrevChar) && unicode.IsUpper(prevChar) && unicode.IsUpper(ch))

		if caseChange {
			ret[i-1+changes] = '-'
			ret = append(ret, prevChar)
			changes++
		}

		ret = append(ret, ch)
		prevPrevChar = prevChar
		prevChar = ch
	}
	return strings.ToLower(string(ret))
}

func (p *parameter) envName() string {
	return strings.Replace(strings.ToUpper(p.hyphenName()), "-", "_", -1)
}
func toInternalType(value interface{}) internalType {
	if _, ok := value.(time.Duration); ok {
		return durationType
	}
	if _, ok := value.(uint); ok {
		return uintType
	}
	if _, ok := value.(int); ok {
		return intType
	}
	if _, ok := value.(float64); ok {
		return floatType
	}
	if _, ok := value.(bool); ok {
		return boolType
	}
	if _, ok := value.(string); ok {
		return stringType
	}
	return invalidType
}

func (p *parameter) SetValueAsString(val string) error {

	switch p.paramtype {
	case stringType:
		p.value = val
	case uintType:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil || v < 0 {
			return fmt.Errorf("invalid value for field %s: %s", p.name, val)
		}
		p.value = uint(v)

	case intType:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value for field %s: %s", p.name, val)
		}
		p.value = int(v)
	case boolType:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("invalid value for field %s: %s", p.name, val)
		}
		p.value = v
	case durationType:
		v, err := time.ParseDuration(val)
		if err != nil {
			return fmt.Errorf("invalid value for field %s: %s", p.name, val)
		}
		p.value = v
	case floatType:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("invalid value for field %s: %s", p.name, val)
		}
		p.value = v

	default:
		return fmt.Errorf("unknown parameter type: %v", p.paramtype)
	}
	p.isSet = true
	return nil
}

// SetValue sets the parameter value
func (p *parameter) SetValue(value interface{}) error {
	if toInternalType(value) != p.paramtype {
		return fmt.Errorf("can't set %s to %v since type is %T", p.name, value, value)
	}
	switch p.paramtype {
	case stringType:
		p.value = value.(string)
	case uintType:
		p.value = value.(uint)
	case intType:
		p.value = value.(int)
	case floatType:
		p.value = value.(float64)
	case boolType:
		p.value = value.(bool)
	case durationType:
		v := value.(time.Duration)
		p.value = v
	default:
		return fmt.Errorf("unknown type: %v", p.paramtype)
	}
	p.isSet = true
	return nil
}

func newParameter(prefix string, field reflect.StructField, value interface{}, fieldval reflect.Value) (*parameter, error) {
	ret := parameter{}
	tagValue, ok := field.Tag.Lookup(tagName)
	if !ok {
		// no tag - ignore it
		return nil, nil
	}
	ret.name = prefix + field.Name
	attribs := strings.Split(tagValue, ";")
	ret.paramtype = toInternalType(value)
	if ret.paramtype == invalidType {
		return nil, fmt.Errorf("field %s has an unknown field type", ret.name)
	}
	ret.value = nil
	for _, v := range attribs {
		tv := strings.Split(v, "=")
		if len(tv) > 2 {
			return nil, fmt.Errorf("invalid format for parameter %s: %s", ret.name, v)
		}
		if tv[0] == "" {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(tv[0])) {
		case "desc":
			ret.description = tv[1]
		case "default":
			ret.defaultValue = tv[1]
		case "min":
			ret.minvalue = tv[1]
			if ret.paramtype != intType &&
				ret.paramtype != uintType &&
				ret.paramtype != floatType {
				return nil, fmt.Errorf("field %s must be a numeric type if min parameter is set", ret.name)
			}
		case "max":
			ret.maxvalue = tv[1]
			if ret.paramtype != intType &&
				ret.paramtype != uintType &&
				ret.paramtype != floatType {
				return nil, fmt.Errorf("field %s must be a numeric type if max parameter is set", ret.name)
			}
		case "file":
			if ret.paramtype != stringType {
				return nil, fmt.Errorf("field %s must be of string type if file flag is set", ret.name)
			}
			ret.file = true

		case "required":
			ret.required = true
		case "options":
			if ret.paramtype != stringType {
				return nil, fmt.Errorf("field %s must be of string type if options flag is set", ret.name)
			}
			ret.options = strings.Split(tv[1], ",")
			if len(ret.options) == 1 && ret.options[0] == "" {
				return nil, fmt.Errorf("field %s does not contain any options", ret.name)
			}
		default:
			return nil, fmt.Errorf("field %s has invalid tags", ret.name)
		}
	}
	if ret.minvalue != "" {
		// ensure value is legal
		if _, err := strconv.ParseFloat(ret.minvalue, 64); err != nil {
			return nil, fmt.Errorf("invalid min value for field %s", ret.name)
		}
	}
	if ret.maxvalue != "" {
		if _, err := strconv.ParseFloat(ret.maxvalue, 64); err != nil {
			return nil, fmt.Errorf("invalid max value for field %s", ret.name)
		}
	}
	if ret.defaultValue != "" {
		if err := ret.SetValueAsString(ret.defaultValue); err != nil {
			return nil, err
		}
	}
	ret.isSet = false
	return &ret, nil
}

func (p *parameter) toFloat() float64 {
	switch p.paramtype {
	case intType:
		return float64(p.value.(int))
	case uintType:
		return float64(p.value.(uint))
	case floatType:
		return p.value.(float64)
	}
	return 0
}

func (p *parameter) validate() error {
	if p.required && !p.isSet {
		return fmt.Errorf("missing required parameter: %s", p.name)
	}
	if p.minvalue != "" {
		v, _ := strconv.ParseFloat(p.minvalue, 64)

		if p.toFloat() < v {
			return fmt.Errorf("value for %s is below the minimum", p.name)
		}
	}
	if p.maxvalue != "" {
		v, _ := strconv.ParseFloat(p.maxvalue, 64)
		if p.toFloat() > v {
			return fmt.Errorf("value for %s is above the minimum", p.name)
		}
	}
	if len(p.options) > 0 {
		found := false
		for _, v := range p.options {
			if p.value == nil {
				continue
			}
			if strings.EqualFold(v, p.value.(string)) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("option %v is not valid for %s", p.value, p.name)
		}
	}

	if p.file && p.value != nil && p.value.(string) != "" {
		if _, err := os.Stat(p.value.(string)); err != nil {
			return err
		}
	}
	return nil
}
