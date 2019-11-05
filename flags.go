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
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type flagDef struct {
	flagName  string
	name      string
	value     interface{}
	paramtype internalType
}

// flagValue returns the value as a scalar value
func (d *flagDef) flagValue() interface{} {
	switch d.paramtype {
	case stringType:
		return *d.value.(*string)

	case boolType:
		return *d.value.(*bool)
	case intType:
		return *d.value.(*int)
	case uintType:
		return *d.value.(*uint)
	case floatType:
		return *d.value.(*float64)
	case durationType:
		return *d.value.(*time.Duration)
	default:
		panic(fmt.Sprintf("can't %v", d.paramtype))
	}
}

func defaultAsInt(defaultVal string) int {
	if defaultVal == "" {
		return 0
	}
	dv, _ := strconv.ParseInt(defaultVal, 10, 64)
	return int(dv)
}

func defaultAsBool(defaultVal string) bool {
	d, _ := strconv.ParseBool(defaultVal)
	return d
}
func defaultAsFloat64(defaultVal string) float64 {
	d, _ := strconv.ParseFloat(defaultVal, 32)
	return d
}
func defaultAsDuration(defaultVal string) time.Duration {
	d, _ := time.ParseDuration(defaultVal)
	return d
}
func makeFlag(fs *flag.FlagSet, p parameter) (*flagDef, error) {
	ret := flagDef{flagName: p.hyphenName(), name: p.name, paramtype: p.paramtype}
	switch p.paramtype {
	case stringType:
		var s string
		ret.value = &s
		fs.StringVar(&s, ret.flagName, p.defaultValue, p.description)
	case boolType:
		var b bool
		ret.value = &b
		fs.BoolVar(&b, ret.flagName, defaultAsBool(p.defaultValue), p.description)
	case intType:
		var v int
		ret.value = &v
		fs.IntVar(&v, ret.flagName, defaultAsInt(p.defaultValue), p.description)
	case uintType:
		var v uint
		ret.value = &v
		fs.UintVar(&v, ret.flagName, uint(defaultAsInt(p.defaultValue)), p.description)
	case floatType:
		var v float64
		ret.value = &v
		fs.Float64Var(&v, ret.flagName, defaultAsFloat64(p.defaultValue), p.description)
	case durationType:
		var v time.Duration
		ret.value = &v
		fs.DurationVar(&v, ret.flagName, defaultAsDuration(p.defaultValue), p.description)
	default:
		return nil, fmt.Errorf("can't make flag for %s:%v", p.name, p.paramtype)
	}
	return &ret, nil
}

// NewFlag parses the command line parameters. This uses the flag package
// internally. Flag names are derived from the names in the configuration structures
// If you have structs within structs the names are prefixed with the name of the
// structure containing the fields:
//
//  type serverConfig struct {
//      HTTP httpConfig  // These parameters will be prefixed with http-
//      HostName string  // This parameter will be named host-name
//  }
//
//  type httpConfig struct {
//      Endpoint string  // This parameter will be named http-endpoint
//      TLS      bool    // This parameter will be named http-tls
//  }
//
func NewFlag(config interface{}, args []string) error {
	return newFlagWithErrorHandling(config, args, flag.ExitOnError, false)
}

// NewEnvFlag returns a struct populated with settings from environment
// variables and command line arguments. The command line arguments overrides
// the environment variables.
func NewEnvFlag(config interface{}, args []string) error {
	return newFlagWithErrorHandling(config, args, flag.ExitOnError, true)
}

// newFlagWithErrorHandling is just for testing; flag.ContinueOnError
// keeps executing but returns an error
func newFlagWithErrorHandling(config interface{}, args []string, opt flag.ErrorHandling, envOverride bool) error {
	params, err := newConfigParameters(config)
	if err != nil {
		return err
	}

	fs := flag.NewFlagSet("parameters", opt)

	flagVars := make([]*flagDef, len(params.params))
	for i, p := range params.params {
		flagVars[i], err = makeFlag(fs, p)
		if err != nil {
			return err
		}
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Check for environment overrides
	if envOverride {
		for i, p := range params.params {
			val, ok := os.LookupEnv(p.envName())
			if ok {
				if err := params.params[i].SetValueAsString(val); err != nil {
					return err
				}
			}
		}
	}

	// Set flags if they're specified
	var flagsToSet []*flagDef
	fs.Visit(func(f *flag.Flag) {
		for i := range flagVars {
			if flagVars[i].flagName == f.Name {
				flagsToSet = append(flagsToSet, flagVars[i])
				return
			}
		}
	})

	for i := range flagsToSet {
		p := params.getParameter(flagsToSet[i].name)
		if p == nil {
			continue
		}
		if err := p.SetValue(flagsToSet[i].flagValue()); err != nil {
			return err
		}
	}
	params.AssignValues(config)
	return params.Validate()
}
