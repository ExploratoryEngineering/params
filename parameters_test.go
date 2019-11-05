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
	"testing"
	"time"
)

// Parameter types: endpoint, file, bool, string, int, group, duration
type deviceIOconfig struct {
	Endpoint     string `param:"desc=Server endpoint;default=:0"`
	CertFile     string `param:"desc=Server cert file;file"`
	KeyFile      string `param:"desc=Server key file;file"`
	HostOverride string `param:"desc=Client host override"`
	CAFile       string `param:"desc=Client CA file;file"`
	Standalone   bool   `param:"desc=Standalone server;default=false"`
}

type httpConfig struct {
	Endpoint      string `param:"desc=Server endpoint;default=:8080"`
	TLSCertFile   string `param:"desc=TLS cert file;file"`
	TLSKeyFile    string `param:"desc=TLS key file;file"`
	ACMECert      bool   `param:"desc=Let's Encrypt ACME certs;default=false"`
	ACMEHosts     string `param:"desc=ACME host names"`
	ACMESecretDir string `param:"desc=ACME secret dir"`
}

type radiusConfig struct {
	AuthEndpoint string `param:"desc=RADIUS auth endpoint;default=:1812"`
	APN          string `param:"desc=RADIUS APN name;default=mda.ee"`
	CIDR         string `param:"desc=RADIUS IP address pool;default=10.0.0.0/13"`
	SharedSecret string `param:"desc=RADIUS shared secret;default=radiussharedsecret"`
}

type parameters struct {
	DeviceIO    deviceIOconfig
	HTTP        httpConfig
	RADIUS      radiusConfig
	LogType     string        `param:"desc=Log type;options=plain,syslog,fancy,ansi,full;default=plain"`
	MyURL       string        `param:"desc=URL parameter 1;default=https://example.com/"`
	MyVal       int           `param:"desc=Int parameter;default=1"`
	MyUint      uint          `param:"desc=Uint parameter;default=2;min=0;max=12000"`
	MyFloat     float64       `param:"desc=Float parameter;default=0.5;max=9999.0"`
	MyDuration  time.Duration `param:"desc=Float parameter;default=10ns"`
	MyBool      bool          `param:"desc=Boolean parameter;default=false"`
	MyOtherBool bool          `param:"default=true"`
}

func TestInvalidParameters(t *testing.T) {
	var cfg struct {
		f  uint
		F1 string `param:"invalid;tags;for;field=1"`
	}
	if _, err := newConfigParameters(&cfg); err == nil {
		t.Fatal("Expected error")
	}
}

func TestUnexportedFields(t *testing.T) {
	var cfg struct {
		f0 string
		f1 string `param:"desc=Foo;default=Bar"`
	}
	if _, err := newConfigParameters(&cfg); err == nil {
		t.Fatal("Expected error for unexported methods")
	}

}

func TestInvalidDefault(t *testing.T) {
	var cfg struct {
		F0 int `param:"desc=Foo;default=Bar"`
	}
	if _, err := newConfigParameters(&cfg); err == nil {
		t.Fatal("Expected error")
	}
}

func TestNonStruct(t *testing.T) {
	num := 0
	if _, err := newConfigParameters(&num); err == nil {
		t.Fatal("Expected error")
	}
}
