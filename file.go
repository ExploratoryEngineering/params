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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

// Flatten nested JSON structs into a single level, ie to internal representation
// of config.
func flattenMap(prefix string, in, out map[string]interface{}) {
	for k, v := range in {
		submap, ok := v.(map[string]interface{})
		if ok {
			flattenMap(prefix+strings.ToLower(k)+".", submap, out)
			continue
		}
		out[prefix+strings.ToLower(k)] = v
	}
}

// NewFile populates a config struct with values from a config file
func NewFile(config interface{}, reader io.Reader) error {
	params, err := newConfigParameters(config)
	if err != nil {
		return err
	}
	if reader == nil {
		return errors.New("reader must be non-nil")
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal(bytes, &jsonMap)
	if err != nil {
		return err
	}
	// Flatten config into keys, all lowercase
	configMap := make(map[string]interface{})
	flattenMap("", jsonMap, configMap)
	for k, v := range configMap {
		// Fix the types if int or uint is used
		para := params.getParameter(k)
		if para == nil {
			continue
		}
		if para.paramtype == intType {
			tmp := v.(float64)
			if err := para.SetValue(int(tmp)); err != nil {
				return err
			}
			continue
		}
		if para.paramtype == uintType {
			tmp := v.(float64)
			if err := para.SetValue(uint(tmp)); err != nil {
				return err
			}
			continue
		}
		if para.paramtype == durationType {
			tmp, ok := v.(string)
			if !ok {
				return fmt.Errorf("can't parse duration field %s", para.name)
			}
			d, err := time.ParseDuration(tmp)
			if err != nil {
				return fmt.Errorf("field %s isn't properly formatted", para.name)
			}
			para.SetValue(d)
			continue
		}
		if err := para.SetValue(v); err != nil {
			return err
		}
	}
	if err := params.AssignValues(config); err != nil {
		return err
	}
	return params.Validate()
}
