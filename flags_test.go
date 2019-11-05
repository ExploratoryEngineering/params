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
	"testing"
)

func TestCommandLineParameters(t *testing.T) {
	if err := NewFlag(nil, nil); err == nil {
		t.Fatal("Expected error with nil arguments")
	}
	config := parameters{}
	if err := NewFlag(config, []string{}); err == nil {
		t.Fatal("Expected error wit non-pointer parameters")
	}

	if err := NewFlag(&config, []string{
		"--http-endpoint", ":1234",
		"-device-io-endpoint=localhost:4711",
		"--radius-auth-endpoint", "localhost:9999",
		"-my-bool",
	}); err != nil {
		t.Fatal(err)
	}
	if config.HTTP.Endpoint != ":1234" {
		t.Fatalf("Config isn't set for http endpoint: %+v", config.HTTP)

	}
	if config.DeviceIO.Endpoint != "localhost:4711" {
		t.Fatalf("Config isn't set for deviceIO: %+v", config.DeviceIO)

	}
	if config.RADIUS.AuthEndpoint != "localhost:9999" {
		t.Fatalf("Config isn't set for radius: %+v", config.RADIUS)
	}

	if !config.MyBool {
		t.Fatalf("Bool flag isn't set: %+v", config)
	}
}

func TestInvalidFlags(t *testing.T) {
	config := parameters{}
	if err := newFlagWithErrorHandling(&config, []string{
		"--incorrect-parameter-value=12",
	}, flag.ContinueOnError, false); err == nil {
		t.Fatal("Expected error")
	}
}
