#  DotEnvConfig

Library responsible for configuring Go applications.

## Installation

`go get github.com/andrewapj/dotenvconfig`

## Usage

### Profile Selection

A profile is used to determine which .env config file should be loaded.
The current profile can be selected in three ways, in order of precedence:

- Lookup profile via an environment variable.
```
err := dotenvconfig.Load(fs, Options{
	ProfileKey: "ENV_VAR",
})
```
- Specify a profile
```
err := dotenvconfig.Load(fs, Options{
	Profile: "development",
})
```

- If no profile is specified then `default` is used.
- The library will look for a file called `{profile}.env`
- Environment variables are set from the values in the config. 
If an existing environment variable already exist then it will be preserved. 

### Config file

A config file should end with `.env` and its content should be in the following format:
```
TEST_KEY=123

# This is a comment

TEST_KEY2 = 456
```

### Get Key

Get a key from the config
```
val := cfg.GetKey("KEY")
```

### Get Key Must

Get a key from the config, panic if the key is missing
```
val := cfg.GetKeyMust("KEY")
```

### Get Int Key

Get an integer value from the config.
```
i := cfg.GetKeyAsInt("TEST_KEY")
```

### Get Int Key Must

Get an integer value from the config, panic if the key is missing
```
i := cfg.GetKeyAsIntMust("TEST_KEY")
```
