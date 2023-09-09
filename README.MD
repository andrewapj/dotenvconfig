#  DotEnvConfig

Library responsible for configuring Go applications.

## Installation

`go get github.com/andrewapj/dotenvconfig`

## Usage

### Profile Selection

The current profile can be selected in three ways, in order of precedence.

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
- Environment variables are set from the values in the config. If they already exist then they will be preserved. 

### Config file

A config file should end with `.env` and be in the following format:
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

### Get int Key

Get an integer value from the config.
```
i := cfg.GetKeyAsInt("TEST_KEY")
```
