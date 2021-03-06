package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func testPrivValidatorKey() string {
	return `{
"address": "5BCD69E0178E0E6C6F96F541B265CAE3178611AE",
"pub_key": {
  "type": "tendermint/PubKeyEd25519",
  "value": "KwddNyi18Ta7tPs6xwfM79O3waMn1+aJuB6GyGQjYuY="
},
"priv_key": {
  "type": "tendermint/PrivKeyEd25519",
  "value": "XQpf+QIrfT/3v0yLquLhfJ5dUaQfJ+ScLYoPzjpUuTkrB103KLXxNru0+zrHB8zv07fBoyfX5om4HobIZCNi5g=="
  }
}`
}

func testPrivValidatorState() string {
	return `{
  "height": "0",
  "round": 0,
  "step": 0
}`
}

func testValidConfig() *Config {
	return &Config{
		Init: InitConfig{
			LogLevel:               "INFO",
			SetSize:                2,
			Threshold:              10,
			Rank:                   1,
			ValidatorListenAddr:    "127.0.0.1:4000",
			ValidatorListenAddrRPC: "127.0.0.1:26657",
		},
		FilePV: FilePVConfig{
			ChainID:       "testchain",
			KeyFilePath:   "./priv_validator_key.json",
			StateFilePath: "./priv_validator_state.json",
		},
	}
}

func TestInitDir(t *testing.T) {
	configDir := "./.test"
	defer os.RemoveAll(configDir)

	if err := InitDir(configDir); err != nil {
		t.Errorf("Expected err to be nil, instead got: %v", err)
	}
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Expected %v to have been created, instead it doesn't exist", configDir)
	}
}

func TestGetDir(t *testing.T) {
	defaultDir := os.Getenv("HOME") + "/.signctrl"

	os.Unsetenv("SIGNCTRL_CONFIG_DIR")
	if GetDir() != defaultDir {
		t.Errorf("Expected SIGNCTRL_CONFIG_DIR to be \"%v\", instead got: %v", defaultDir, os.Getenv("SIGNCTRL_CONFIG_DIR"))
	}

	os.Setenv("SIGNCTRL_CONFIG_DIR", "/some/random/dir")
	if GetDir() != "/some/random/dir" {
		t.Errorf("Expected SIGNCTRL_CONFIG_DIR to be \"/some/random/dir\", instead got: %v", os.Getenv("SIGNCTRL_CONFIG_DIR"))
	}
}

func TestValidate(t *testing.T) {
	ioutil.WriteFile("./priv_validator_key.json", []byte(testPrivValidatorKey()), 0644)
	ioutil.WriteFile("./priv_validator_state.json", []byte(testPrivValidatorState()), 0644)
	defer os.Remove("./priv_validator_key.json")
	defer os.Remove("./priv_validator_state.json")

	config := testValidConfig()

	// Valid config.
	if err := config.validate(); err != nil {
		t.Errorf("Expected err to be nil, instead got: %v", err)
	}

	// Invalid loglevel.
	config.Init.LogLevel = "INVALID"
	if err := config.validate(); err == nil {
		t.Errorf("Expected err, instead got nil")
	}
	config.Init.LogLevel = "INFO"

	// Invalid setsize.
	config.Init.SetSize = 0
	if err := config.validate(); err == nil {
		t.Errorf("Expected err, instead got nil")
	}
	config.Init.SetSize = 2

	// Invalid threshold.
	config.Init.Threshold = 0
	if err := config.validate(); err == nil {
		t.Errorf("Expected err, instead got nil")
	}
	config.Init.Threshold = 10

	// Invalid rank.
	config.Init.Rank = 0
	if err := config.validate(); err == nil {
		t.Errorf("Expected err, instead got nil")
	}
	config.Init.Rank = 1

	// Invalid validator listen address.
	config.Init.ValidatorListenAddr = "127.0.0.1"
	if err := config.validate(); err == nil {
		t.Errorf("Expected err, instead got nil")
	}
	config.Init.ValidatorListenAddr = "127.0.0.1:4000"

	// Invalid validator rpc listen address.
	config.Init.ValidatorListenAddrRPC = "127.0.0.1"
	if err := config.validate(); err == nil {
		t.Errorf("Expected err, instead got nil")
	}
	config.Init.ValidatorListenAddrRPC = "127.0.0.1:26657"

	// Invalid chainid.
	config.FilePV.ChainID = ""
	if err := config.validate(); err == nil {
		t.Errorf("Expected err, instead got nil")
	}
	config.FilePV.ChainID = "testchain"

	// Non-existent path to keyfile.
	config.FilePV.KeyFilePath = "/this/path/does/not/exist"
	if err := config.validate(); err == nil {
		t.Errorf("Expected err, instead got nil")
	}
	config.FilePV.KeyFilePath = "./priv_validator_key.json"

	// Non-existent path to statefile.
	config.FilePV.StateFilePath = "/this/path/does/not/exist"
	if err := config.validate(); err == nil {
		t.Errorf("Expected err, instead got nil")
	}
	config.FilePV.StateFilePath = "./priv_validator_state.json"
}
