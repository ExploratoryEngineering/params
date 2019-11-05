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
	"os"
	"testing"
	"time"
)

func TestNameConversion(t *testing.T) {
	tf := func(name, expected string) {
		p := parameter{name: name}
		if p.envName() != expected {
			t.Fatal(expected, " != ", p.envName())
		}
	}
	tf("Some", "SOME")
	tf("Some.Name", "SOME_NAME")
	tf("a.B.c", "A_B_C")
}
func TestEnvironmentVariables(t *testing.T) {
	os.Setenv("HTTP_ENDPOINT", ":1234")
	os.Setenv("DEVICE_IO_ENDPOINT", "localhost:4711")
	os.Setenv("RADIUS_AUTH_ENDPOINT", "localhost:9999")
	os.Setenv("MY_BOOL", "true")
	os.Setenv("MY_DURATION", "451ms")

	config := parameters{}

	if err := NewEnv(&config); err != nil {
		t.Fatalf("Error creating config: %v", err)
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
		t.Fatalf("Bool isn't set")
	}
	if config.MyDuration != 451*time.Millisecond {
		t.Fatalf("Config isn't set for duration: %+v", config)
	}
}

func TestEnvironmentInvalidValue(t *testing.T) {
	var c struct {
		F0 int `param:"desc=Some desc"`
	}
	os.Setenv("F0", "foo")
	if err := NewEnv(&c); err == nil {
		t.Fatal("Expected error")
	}
}

func TestEnvironmentInvalidType(t *testing.T) {
	if err := NewEnv(nil); err == nil {
		t.Fatal("Expected error")
	}
}

func TestEnvironment(t *testing.T) {
	var cfg struct {
		F1 string `param:"default=F1"`
		F2 string `param:"default=F2"`
		F3 string `param:""`
	}

	os.Setenv("F2", "env2")
	os.Setenv("F3", "env3")
	if err := NewEnv(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.F1 != "F1" || cfg.F2 != "env2" || cfg.F3 != "env3" {
		t.Fatalf("Config not set properly: %+v", cfg)
	}
}
