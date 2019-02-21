Sample read config
===

Goals:

1. Read in a config file into a Go structure.
2. Substitute in environment variables for strings staring with "$ENV$"
3. Look for the configuration in ~/local first - so that servers can be configured.


Note - test in ./ReadConfig will create a `~/local` directory with `b.json` in it when `go test` is run.

Based on open source MIT licensed.

If you use a string with `$ENV$Name` then it will pull that value from the 
environment using `Name` for the environment variable.   A similar basename
file in `~/local` will override the one specified.

You can create default values with a structure tag.  Look at ./sample-main.go for
an example.

So...

```
	err := config.ReadConfig ( "../mysetup.json", &someStruct )
```

to use.

An example in ./sample-main.go

