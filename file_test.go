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
	"strings"
	"testing"
)

func TestParameterFile(t *testing.T) {
	config := parameters{}

	if err := NewFile(nil, nil); err == nil {
		t.Fatal("Expected invalid config here")
	}

	if err := NewFile(&config, nil); err == nil {
		t.Fatal("Expected invalid config here")
	}

	if err := NewFile(&config, strings.NewReader("invalidjson}")); err == nil {
		t.Fatal("Expected invalid file format")
	}

	if err := NewFile(&config, strings.NewReader("{}")); err != nil {
		t.Fatal("Did not expect and error when reading empty JSON: ", err)
	}

	file := `{
		"http": {
			"endpoint": ":1234"
		},
		"deviceIO": {
			"endpoint": "localhost:4711"
		},
		"radius": {
			"authEndpoint": "localhost:9999"
		},
		"myUrl": "https://example.com/",
		"myVal": 12,
		"myUint": 11,
		"myBool": true,
		"myUint": 4711,
		"myDuration": "12ms"
	}`

	if err := NewFile(&config, strings.NewReader(file)); err != nil {
		t.Fatal("Did not expect an error when reading JSON: ", err)
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
	if config.MyURL != "https://example.com/" {
		t.Fatalf("Config isn't set for URL: %+v", config)
	}
}

func TestDebugging(t *testing.T) {
	config := parameters{}
	file := `{
		"http": {
			"endpoint": ":1234"
		},
		"deviceIO": {
			"endpoint": "localhost:4711"
		},
		"radius": {
			"authEndpoint": "localhost:9999"
		},
		"myUrl": "https://example.com/",
		"myUInt": 12,
		"myInt": 13,
		"myBool": true,
		"myUint": 4711,
		"myDuration": "12ms"
	}`

	if err := NewFile(&config, strings.NewReader(file)); err != nil {
		t.Fatal("Did not expect an error when reading JSON: ", err)
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
	if config.MyURL != "https://example.com/" {
		t.Fatalf("Config isn't set for URL: %+v", config)
	}
}

type invalidReader struct {
}

func (i *invalidReader) Read([]byte) (int, error) {
	return 0, errors.New("i'm not working")
}
func TestInvalidReader(t *testing.T) {
	config := parameters{}
	if err := NewFile(&config, &invalidReader{}); err == nil {
		t.Fatal("Expected error")
	}
}

func TestInvalidConfig(t *testing.T) {
	config := parameters{}
	file := `{
		"http": {
			"endpoint": 1
		}
	}`

	if err := NewFile(&config, strings.NewReader(file)); err == nil {
		t.Fatal("Expected error")
	}
}
