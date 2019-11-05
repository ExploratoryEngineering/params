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
	"strings"
	"testing"
)

func TestHyphenName(t *testing.T) {
	testString := func(s, e string) {
		p := parameter{name: s}
		out := p.hyphenName()
		if out != e {
			t.Fatalf("Can't hypen: %s = %s", p.name, out)
		}
	}
	testString("Na", "na")
	testString("NaNa", "na-na")
	testString("NaNaNa", "na-na-na")
	testString("BATMAN", "batman")

	testString("NanNan", "nan-nan")

	testString("NameName", "name-name")
	testString("A", "a")
	testString("AA", "aa")
	testString("Ab", "ab")
	testString("NameTLA", "name-tla")
	testString("TLAName", "tla-name")
	testString("TLANameTLA", "tla-name-tla")
	testString("NameTLAName", "name-tla-name")
}
func TestRequiredParameter(t *testing.T) {
	var requiredOpts struct {
		ParamA string `param:"desc=Description;required"`
	}

	if err := NewFile(&requiredOpts, strings.NewReader("{}")); err == nil {
		t.Fatalf("Expected error with required set but got none")
	}

	if err := NewFile(&requiredOpts, strings.NewReader(`{"paramA": "yes"}`)); err != nil {
		t.Fatalf("Did not expect error when parameter is set: %v (%+v)", err, requiredOpts)
	}

	if err := NewEnv(&requiredOpts); err == nil {
		t.Fatalf("Expected error with required set but got none")
	}
}

func TestMinMaxParameter(t *testing.T) {
	var pc1 struct {
		Val int `param:"desc=description;min=100;max=200"`
	}

	if err := NewFlag(&pc1, []string{"--val", "100"}); err != nil || pc1.Val != 100 {
		t.Fatalf("Got error or value not set for v1: %v %+v", err, pc1)
	}

	var pc2 struct {
		Val uint `param:"desc=description;min=100;max=200"`
	}
	if err := NewFlag(&pc2, []string{"--val=99"}); err == nil {
		t.Fatalf("Expected error when value is outside range: %+v", pc1)
	}

	var pc3 struct {
		Val float64 `param:"desc=description;min=100;max=200"`
	}
	if err := NewFlag(&pc3, []string{"--val=201"}); err == nil {
		t.Fatal("Expected error when value is above max")
	}

	var pc4 struct {
		ValOne uint `param:"desc=description;min=1"`
	}
	if err := NewFlag(&pc4, []string{"--val-one", "9999"}); err != nil || pc4.ValOne != 9999 {
		t.Fatalf("Got error or value not set: %v %+v", err, pc4)
	}
}

func TestMinMaxDefaults(t *testing.T) {
	var pc1 struct {
		Val        uint `param:"desc=foo;min=1;max=200;default=2"`
		privateVal uint
	}
	if err := NewFlag(&pc1, []string{}); err != nil {
		t.Fatalf("Default config should work but got error: %v", err)
	}
}
func TestInvalidMinMax(t *testing.T) {
	var pc1 struct {
		Val string `param:"desc=foo;min=10"`
	}
	if err := NewEnv(&pc1); err == nil {
		t.Fatal("Expected error when min is set for string parameter")
	}
	var pc2 struct {
		Val string `param:"desc=foo;max=12"`
	}
	if err := NewEnv(&pc2); err == nil {
		t.Fatal("Expected error when max is set for string parameter")
	}
	var pc3 struct {
		Val int `param:"desc=foo;max=one"`
	}
	if err := NewEnv(&pc3); err == nil {
		t.Fatal("Expected error when max is invalid")
	}
	var pc4 struct {
		Val int `param:"desc=foo;min=zero"`
	}
	if err := NewEnv(&pc4); err == nil {
		t.Fatal("Expected error when min is invalid")
	}
}

func TestOptions(t *testing.T) {
	var pc1 struct {
		Val string `param:"desc=foo;options=a,B,C,d"`
	}
	if err := NewFlag(&pc1, []string{"--val", "A"}); err != nil {
		t.Fatal("Can't assign option a even if it is valid: ", err)
	}
	if err := NewFlag(&pc1, []string{"--val", "q"}); err == nil {
		t.Fatal("Should not be allowed to use option q")
	}
	if err := NewFlag(&pc1, []string{"--val", "b"}); err != nil {
		t.Fatal("Can't assign option b even if it is valid: ", err)
	}
}

func TestInvalidOptions(t *testing.T) {
	var pc1 struct {
		Val int `param:"desc=foo;options=op1,op2,op3,op4"`
	}
	if err := NewFile(&pc1, strings.NewReader("{}")); err == nil {
		t.Fatal("Expected error when int is used for options")
	}
	var pc2 struct {
		Val string `param:"desc=foo;options="`
	}
	if err := NewFile(&pc2, strings.NewReader("{}")); err == nil {
		t.Fatal("Expected error with zero options")
	}
}

func TestFileType(t *testing.T) {
	var pc1 struct {
		File1 string `param:"desc=foo;file"`
	}
	os.Remove("foo.txt")
	if err := NewFile(&pc1, strings.NewReader(`{"file1": "foo.txt"}`)); err == nil {
		t.Fatal("Expected error when file does not exist")
	}
	f, _ := os.Create("foo.txt")
	f.Close()
	defer os.Remove("foo.txt")
	if err := NewFile(&pc1, strings.NewReader(`{"file1": "foo.txt"}`)); err != nil {
		t.Fatal("Expected no error when file exists: ", err)
	}
}
