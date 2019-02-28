package ReadConfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/fatih/structtag"
	"github.com/pschlump/jsonSyntaxErrorLib"
)

// ReadFile will read a configuration file into the global configuration structure.
func ReadFile(filename string, lCfg interface{}) (err error) {

	// Get the type and value of the argument we were passed.
	ptyp := reflect.TypeOf(lCfg)
	pval := reflect.ValueOf(lCfg)

	// Requries that lCfg is a pointer.
	if ptyp.Kind() != reflect.Ptr {
		fmt.Fprintf(os.Stderr, "Must pass a address of a struct to RedFile\n")
		os.Exit(1)
	}

	var typ reflect.Type
	var val reflect.Value
	typ = ptyp.Elem()
	val = pval.Elem()

	// Create Defaults

	// Make sure we now have a struct
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("SetFromEnv was not passed a struct.\n")
	}

	// Can we set values?
	if val.CanSet() {
		if db1 {
			fmt.Printf("Debug: We can set values.\n")
		}
	} else {
		return fmt.Errorf("SetFromEnv passed a struct that will not allow setting of values\n")
	}

	// The number of fields in the struct is determined by the type of struct it is. Loop through them.
	for i := 0; i < typ.NumField(); i++ {

		// Get the type of the field from the type of the struct. For a struct, you always get a StructField.
		sfld := typ.Field(i)

		// Get the type of the StructField, which is the type actually stored in that field of the struct.
		tfld := sfld.Type

		// Get the Kind of that type, which will be the underlying base type
		// used to define the type in question.
		kind := tfld.Kind()

		// Get the value of the field from the value of the struct.
		vfld := val.Field(i)
		tag := string(sfld.Tag)

		// ... and start using structtag by parsing the tag
		tags, err := structtag.Parse(tag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse structure tag ->%s<- %s\n", tag, err)
			os.Exit(1)
		}

		// Dump out what we've found
		if db1 {
			fmt.Printf("Debug: struct field %d: name %s type %s kind %s value %v tag ->%s<-\n", i, sfld.Name, tfld, kind, vfld, tag)

			// iterate over all tags
			for tn, t := range tags.Tags() {
				fmt.Printf("\t[%d] tag: %+v\n", tn, t)
			}

			// get a single tag
			defaultTag, err := tags.Get("default")
			if err != nil {
				fmt.Printf("`default` Not Set\n")
			} else {
				fmt.Println(defaultTag)         // Output: default:"foo,omitempty,string"
				fmt.Println(defaultTag.Key)     // Output: default
				fmt.Println(defaultTag.Name)    // Output: foo
				fmt.Println(defaultTag.Options) // Output: [omitempty string]
			}
		}

		defaultTag, err := tags.Get("default")
		// Is that field some kind of string, and is the value one we can set?
		// 1. Other tyeps (all ints, floats) - not just strings		xyzzy001-type
		if kind == reflect.String && vfld.CanSet() {
			if err != nil || defaultTag.Name == "" {
				// Ignore error - indicates no "default" tag set.
			} else {
				defaultValue := defaultTag.Name
				if db1 {
					fmt.Printf("Debug: Looking to set field %s to a default value of ->%s<-\n", sfld.Name, defaultValue)
				}
				vfld.SetString(defaultValue)
			}
		} else if kind != reflect.String && err == nil {
			// report errors - defauilt is only implemented with strings.
			fmt.Fprintf(os.Stderr, "default tag on struct is only implemented for String fields that are settable in struct.  Fatal error on %s tag %s\n", sfld.Name, tag)
			os.Exit(1)
		}
	}

	// look for filename in ~/local (C:\local on Winderz)
	var home string
	if os.PathSeparator == '/' {
		home = os.Getenv("HOME")
	} else {
		home = "C:\\"
	}
	homeLocal := path.Join(home, "local")
	base := path.Base(filename)
	if ExistsIsDir(homeLocal) && Exists(path.Join(homeLocal, base)) {
		filename = path.Join(homeLocal, base)
	}
	if db1 {
		fmt.Printf("Debug: File name after checing ~/local [%s]\n", filename)
	}

	var buf []byte
	buf, err = ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read the JSON file [%s]: error %s\n", filename, err)
		os.Exit(1)
	}

	// err = json.Unmarshal(buf, &gCfg)
	err = json.Unmarshal(buf, lCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid initialization - Unable to parse JSON file, %s\n", err)
		PrintErrorJson(string(buf), err) // show line for error
		os.Exit(1)
	}

	// err = SetFromEnv(&gCfg)
	err = SetFromEnv(lCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error pulling from environment: %s\n", err)
		os.Exit(1)
	}

	return err
}

func PrintErrorJson(js string, err error) (rv string) {
	rv = jsonSyntaxErrorLib.GenerateSyntaxError(js, err)
	fmt.Fprintf(os.Stderr, "%s\n", rv)
	return
}

func SetFromEnv(s interface{}) (err error) {

	// Get the type and value of the argument we were passed.
	ptyp := reflect.TypeOf(s)
	pval := reflect.ValueOf(s)
	// We can't do much with the Value (it's opaque), but we need it in order
	// to fetch individual fields from the struct later.

	var typ reflect.Type
	var val reflect.Value

	// If we were passed a pointer, dereference to get the type and value
	// pointed at.
	if ptyp.Kind() == reflect.Ptr {
		if db2 {
			fmt.Printf("Debug: Argument is a pointer, dereferencing.\n")
		}
		typ = ptyp.Elem()
		val = pval.Elem()
	} else {
		if db2 {
			fmt.Printf("Debug: Argument is %s.%s, a %s.\n", ptyp.PkgPath(), ptyp.Name(), ptyp.Kind())
		}
		typ = ptyp
		val = pval
	}

	// Make sure we now have a struct
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("SetFromEnv was not passed a struct.\n")
	}

	// Can we set values?
	if val.CanSet() {
		if db2 {
			fmt.Printf("Debug: We can set values.\n")
		}
	} else {
		return fmt.Errorf("SetFromEnv passed a struct that will not allow setting of values\n")
	}

	// The number of fields in the struct is determined by the type of struct
	// it is. Loop through them.
	for i := 0; i < typ.NumField(); i++ {

		// Get the type of the field from the type of the struct. For a struct, you always get a StructField.
		sfld := typ.Field(i)

		// Get the type of the StructField, which is the type actually stored in that field of the struct.
		tfld := sfld.Type

		// Get the Kind of that type, which will be the underlying base type
		// used to define the type in question.
		kind := tfld.Kind()

		// Get the value of the field from the value of the struct.
		vfld := val.Field(i)

		// Dump out what we've found
		if db2 {
			fmt.Printf("Debug: struct field %d: name %s type %s kind %s value %v\n", i, sfld.Name, tfld, kind, vfld)
		}

		// Is that field some kind of string, and is the value one we can set?
		// 1. Other tyeps (all ints, floats) - not just strings		xyzzy001-type
		if kind == reflect.String && vfld.CanSet() {
			if db2 {
				fmt.Printf("Debug: Looking to set field %s\n", sfld.Name)
			}
			// Assign to it
			curVal := fmt.Sprintf("%s", vfld)
			if len(curVal) > 5 && curVal[0:5] == "$ENV$" {
				envVal := os.Getenv(curVal[5:])
				if db2 {
					fmt.Printf("Debug: Overwriting field %s current [%s] with [%s]\n", sfld.Name, curVal, envVal)
				}
				vfld.SetString(envVal)
			}
			if len(curVal) > 6 && curVal[0:6] == "$FILE$" {
				data, err := ioutil.ReadFile(curVal[6:])
				if db2 {
					fmt.Printf("Debug: Overwriting field %s current [%s] with [%s]\n", sfld.Name, data, data)
				}
				if err != nil {
				}
				vfld.SetString(string(data))
			}
		}
	}

	return nil
}

// Exists returns true if a directory or file exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// ExistsIsDir returns true if a direcotry exists.
func ExistsIsDir(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	if fi.IsDir() {
		return true
	}
	return false
}

var db1 = false
var db2 = false
