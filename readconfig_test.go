package ReadConfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/pschlump/dbgo"
)

// GlobalConfigData is the gloal configuration data.
// It holds all the data from the cfg.json file.
type GlobalConfigData struct {
	ExampeWithDefault string `default:"dflt-1"`
	SomePassword      string `default:"dflt-2"`
	CheckDefault      string `default:"dflt-3"`
}

var gCfg GlobalConfigData // global configuration data.

func TestMineBlock(t *testing.T) {

	tests := []struct {
		SetEnvName       string
		SetEnvVal        string
		FileName         string
		ExpectedPassword string
	}{
		{
			SetEnvName:       "MyPassword",
			SetEnvVal:        "xyzzy-3",
			FileName:         "./testdata/a.json",
			ExpectedPassword: "xyzzy-3",
		},
		{
			SetEnvName:       "Test2",
			SetEnvVal:        "xyzzy-2",
			FileName:         "./testdata/b.json",
			ExpectedPassword: "xyzzy-2",
		},
	}

	db1 = false // turn on output for debuging in ReadFile
	db2 = false // turn on output for debuging in SetFromEnv

	home, err := homedir.Dir()
	if err != nil {
		if os.PathSeparator == '/' {
			home = os.Getenv("HOME")
		} else {
			home = "C:\\"
		}
	}

	buf := `{
	"SomePassword": "$ENV$Test2"
}
`
	os.Mkdir(home+"/local", 0755)
	ioutil.WriteFile(home+"/local/b.json", []byte(buf), 0644)

	for ii, test := range tests {
		os.Setenv(test.SetEnvName, test.SetEnvVal)
		ReadFile(test.FileName, &gCfg)
		if gCfg.SomePassword != test.ExpectedPassword {
			t.Errorf("Test %d, expected %s got %s\n", ii, test.ExpectedPassword, gCfg.SomePassword)
		}
		if gCfg.CheckDefault != "dflt-3" {
			t.Errorf("Test %d, expected %s got %s\n", ii, "dflt-3", gCfg.SomePassword)
		}
	}

}

// Test with embeded struct.
type Test2EmbedType struct {
	AaaEmbeded string `default:"aaa-embeded-value"`
	AbbEmbeded string `default:"$ENV$Abb"`
}

type Test2Type struct {
	Test2EmbedType
	ExampeWithDefault string `default:"dflt-1"`
	SomePassword      string `default:"dflt-2"`
	CheckDefault      string `default:"dflt-3"`
}

func Test2(t *testing.T) {

	tests := []struct {
		SetEnvName string
		SetEnvVal  string
		FileName   string
		Expected   string
	}{
		{
			SetEnvName: "Abb",
			SetEnvVal:  "abb-value-from-env",
			FileName:   "./testdata/test2.json",
			Expected: `{
	"AaaEmbeded": "aaa-embeded-value",
	"AbbEmbeded": "abb-value-from-env",
	"ExampeWithDefault": "dflt-1",
	"SomePassword": "dflt-2",
	"CheckDefault": "dflt-3"
}`,
		},
	}

	db1 = false // turn on output for debuging in ReadFile
	db2 = false // turn on output for debuging in SetFromEnv
	db3 = false //

	var test2 Test2Type

	for ii, test := range tests {
		os.Setenv(test.SetEnvName, test.SetEnvVal)
		ReadFile(test.FileName, &test2)
		// fmt.Printf("Result: %s\n", dbgo.SVarI(test2))
		got := dbgo.SVarI(test2)
		if got != test.Expected {
			t.Errorf("Test %d, expected %s got %s\n", ii, test.Expected, got)
		}
	}

}

func TestSetup(t *testing.T) {
	if runtime.GOOS == "windows" {
		// fmt.Println("You are running on Windows")
		if home != "C:/" {
			t.Errorf("Home not set correctly for Windows")
		}
	} else {
		if home != os.Getenv("HOME") {
			t.Errorf("Home not set")
		}
	}
}

// Test with embeded struct using a ~ in environment variable.
type Test3EmbedType struct {
	AbbEmbeded string `default:"$ENV$homeTest"`
}

type Test3Type struct {
	Test3EmbedType
}

func Test3(t *testing.T) {

	tests := []struct {
		me         bool
		SetEnvName string
		SetEnvVal  string
		FileName   string
		Expected   string
	}{
		{
			SetEnvName: "homeTest",
			SetEnvVal:  "~/.keystore",
			FileName:   "./testdata/test2.json",
			Expected:   "%s/.keystore",
		},
		{
			me:         true,
			SetEnvName: "homeTest",
			SetEnvVal:  "~pschlump/.keystore",
			FileName:   "./testdata/test2.json",
			Expected:   "%s/.keystore",
		},
	}

	db1 = false // turn on output for debuging in ReadFile
	db2 = false // turn on output for debuging in SetFromEnv
	db3 = false //
	db4 = false // home dir replacement stuff

	var test3 Test3Type

	for ii, test := range tests {
		skip := false
		if test.me && os.Getenv("HOME") == "/Users/pschlump" || test.me && os.Getenv("HOME") == "/Users/philip" {
			skip = true
		}
		if !skip {
			os.Setenv(test.SetEnvName, test.SetEnvVal)
			exp := fmt.Sprintf(test.Expected, home)
			ReadFile(test.FileName, &test3)
			got := test3.AbbEmbeded
			if got != exp {
				t.Errorf("Test %d, expected %s got %s\n", ii, exp, got)
			}
		}
	}

}
