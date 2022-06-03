package ReadConfig

import (
	"os"
	"testing"

	"github.com/pschlump/dbgo"
)

// Test with embeded struct and new default for environment variables.
type Test2EmbedTypeDflt struct {
	AaaEmbeded string `default:"aaa-embeded-value"`
	AbbEmbeded string `default:"$ENV$Abb"`
	AbcEmbeded string `default:"$ENV$Abc=ADefaultEnvNotSet"`
	AbdEmbeded string `default:"$ENV$Abd="`
	AbeEmbeded string `default:"$ENV$Abe"`
}

type Test2TypeDflt struct {
	Test2EmbedTypeDflt
	ExampeWithDefault string `default:"dflt-1"`
	SomePassword      string `default:"dflt-2"`
	CheckDefault      string `default:"dflt-3"`
}

func TestEnvWithDefault(t *testing.T) {

	dbgo.DbPfb(db8, "--------------- New test -------------------\n")

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
	"AbcEmbeded": "ADefaultEnvNotSet",
	"AbdEmbeded": "",
	"AbeEmbeded": "",
	"ExampeWithDefault": "dflt-1",
	"SomePassword": "dflt-2",
	"CheckDefault": "dflt-3"
}`,
		},
	}

	db1 = false // turn on output for debuging in ReadFile
	db2 = false // turn on output for debuging in SetFromEnv
	db3 = false //

	var test2 Test2TypeDflt

	for ii, test := range tests {
		os.Setenv(test.SetEnvName, test.SetEnvVal)
		ReadFile(test.FileName, &test2)
		dbgo.DbPfb(db8, "Result: %s\n", dbgo.SVarI(test2))
		got := dbgo.SVarI(test2)
		if got != test.Expected {
			t.Errorf("Test %d, expected %s got %s\n", ii, test.Expected, got)
		}
	}

}

const db8 = false
