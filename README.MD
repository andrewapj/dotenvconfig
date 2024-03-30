#  DotEnvConfig

Library responsible for configuring Go applications.

## Installation

`go get github.com/andrewapj/dotenvconfig`

## Usage

### Loading a config file
```
err := Load(getFS(), "config.env", Options{})
```

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
val, err := cfg.GetKey("KEY")
```

### Get Key Must

Get a key from the config, panic if the key is missing
```
val := cfg.GetKeyMust("KEY")
```

### Get Int Key

Get an integer value from the config.
```
i, err := cfg.GetKeyAsInt("TEST_KEY")
```

### Get Int Key Must

Get an integer value from the config, panic if the key is missing
```
i := cfg.GetKeyAsIntMust("TEST_KEY")
```
