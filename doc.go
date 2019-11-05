// Package params is a package to define configuration parameters for servers.
//
// Parameters are defined as tags on structs. Configuration structs might have
// structs within structs for parameters.
//
// A limited number of data types are supported: strings (string), integers
// (int, uint), booleans (bool), duration (time.Duration) and floats (float64)
//
// The tags are set with the keyword "param". The fields must be publicly accessible:
//
//   type config struct {
//      NameOfApp string `param:"desc=This is the parameter description"`
//   }
//
// Keywords are separated by semicolons. There is no escaping so defaults
// can't contain equal or semicolons. The following keywords are supported:
//
//  desc      - a description
//  default   - the default value for the parameter
//  min       - minimum value for parameter. Flag must be int, uint, float or Duration
//  max       - maximum value for parameter. Flag must be int, uint, float or Duration
//  file      - if present the flag points to a file and that file must exist. Flag must be a string.
//  required  - if present the flag must be specfified in a valid config
//  options   - a list of options. Type must be string. Options are case insensitive.
//
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
