#  DotEnvConfig

Library responsible for configuring Go applications.

## Installation

`go get github.com/andrewapj/dotenvconfig`

## Usage

### Environment Selection

The current environment can be selected in three ways, in order of precedence.

- Lookup environment via an environment variable.
```
cfg, err := dotenvconfig.NewConfig(fs, Options{
		environmentKey: "ENV_VAR",
	})
```
- Specify a specific environment
```
cfg, err := dotenvconfig.NewConfig(fs, Options{
		environment: "development",
	})
```

- If no environment is specified then `default` is used.
- The library will look for a file called `{environment}.env`

### Config file

A config file should end with `.env` and be in the following format:
```
TEST_KEY=123
```

### Get Key

Get a key from the config
```
val := cfg.GetKey("KEY")
```

Get a key from the config but allow an os environment variable to override it.
```
_ = os.Setenv("TEST_KEY", "NEW_VALUE")
k := cfg.GetKey("TEST_KEY")

fmt.Println(k) // NEW_VALUE
```

### Get int Key

Get an integer value from the config.
```
k := cfg.GetKeyAsInt("TEST_KEY")
```

### Add a config to a context
```
ctx, err := dotenvconfig.ToContext(context.Background(), cfg)
if err != nil {
	panic(err.Error())
}
```

### Get a config from a context
```
cfg, err := dotenvconfig.FromContext(ctx)
if err != nil {
	panic(err.Error())
}
```

