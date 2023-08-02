#  DotEnvConfig

## Overview

Library responsible for configuring a Go application. Makes use of [Godotenv](https://github.com/joho/godotenv).

# Prerequisites

- Go 1.20
- Make

# Building

- Run `make` which will fully test the application.

A file called `coverage.html` will be generated that shows the test coverage of the code.

# Installation

`go get github.com/andrewapj/dotenvconfig`

# Usage

### Config file

A config file should end with `.env` and be in the following format:
```
TEST_KEY=123
```

### Load config

Load the config using a default environment. 
This will look for a file called 'default.env'.

```
c, err := config.NewConfig(os.DirFS("internal/testconfig")).Load()
	if err != nil {
		panic(err.Error())
	}
```

Or load the config using a specific environment.
This will look for a file called 'custom.env'
```
c, err := config.NewConfig(os.DirFS("internal/testconfig")).WithEnvironment("custom").Load()
	if err != nil {
		panic(err.Error())
	}
```

If an environment variable of `APP_ENVIRONMENT` is set then this will take precedence over the previous options.
This will look for a file called 'envvar.env'
```
_ = os.Setenv(configEnvKey, "envvar")
c, err := config.NewConfig(os.DirFS("internal/testconfig")).WithEnvironment("custom").Load()
	if err != nil {
		panic(err.Error())
	}
```

### Get Key

Get a key from the config
```
k, err := c.GetKey("TEST_KEY")
	if err != nil {
		panic(err.Error())
	}
```

Get a key from the config but allow an os environment variable to override it.
```
_ = os.Setenv("TEST_KEY", "NEW_VALUE")
	k, err := c.GetKey("TEST_KEY")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(k) // NEW_VALUE
```

### Get int Key

Get an integer value from the config.
```
k, err := c.GetKeyAsInt("TEST_KEY")
	if err != nil {
		panic(err.Error())
	}
```

### Add a config to a context
```
ctx, err := config.ToContext(context.Background(), cfg)
	if err != nil {
		panic(err.Error())
	}
```

### Get a config from a context
```
cfg, err := config.FromContext(ctx)
	if err != nil {
		panic(err.Error())
	}
```

