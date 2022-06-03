# Read in a JSON Configuration File Withe Defaults

Read a JSON configuration file into a Go structure with
default values specified in tags.

It would be nice to have the default configuration specified in the 
code and just read in the values that are chained.

For example:

```
type GlobalConfigData struct {
	ExampeWithDefault string `default:"dflt-1"`
	SomePassword      string `default:"dflt-2"`
}
```

Specifies the defaults for `SomePassword` as `dflt-2`.  This can also be used with the tags for JSON.

```
type GlobalConfigData struct {
	ExampeWithDefault string `default:"dflt-1"`
	SomePassword      string `json:"some_password" default:"dflt-2"`
}
```

So if we need a value that is different from the default we can specify it in the JSON configuration file.
Let's say the file is `./cfg.json` and it gets read in at the start.

```
{
	"some_password": "bob's-bad!password"
}
```

When it is read in the JSON file will override the default.  The configuration can also come from the
process environment.  This is useful for not putting things like passwords and authentication tokens
into files that get checked into Git.  (The ideal is to use a configuration sever that will authenticate
and fetch a particular set of authentication from a single remote system, place the auth-token or
password into an environment variable and then this code can pick it up)

To do this specify the default for the field as "$ENV$name" where `name` is the name of the environment
variable.

```
type S3Config struct {
	S3_bucket string `json:"s3_bucket" default:"s3://documents"`
	S3_region string `json:"s3_region" default:"$ENV$AWS_REGION"`
}

type GlobalConfigData struct {
	S3Config

	ExampeWithDefault string `default:"dflt-1"`
	SomePassword      string `default:"dflt-2"`
}
```

In this case `s3_region` is pulled from the environment as a default but could be specified in the
JSON configuration file.




## Example

```
package main

import (
	"fmt"
	"os"

	"github.com/pschlump/godebug"
	"git.q8s.co/pschlump//ReadConfig"
)

// GlobalConfigData is the gloal configuration data.
// It holds all the data from the cfg.json file.
type GlobalConfigData struct {
	ExampeWithDefault string `default:"dflt-1"`
	SomePassword      string `default:"dflt-2"`
}

var gCfg GlobalConfigData // global configuration data.

func main() {
	err := ReadConfig.ReadFile("./testdata/a.json", &gCfg)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS: read %s\n", godebug.SVarI(gCfg))
}

```

