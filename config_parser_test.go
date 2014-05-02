package chef

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	_, err := ParseConfig("test/support/knife.rb")
	if err != nil {
		t.Error(err)
	}
}

func TestKeyFromString(t *testing.T) {
	config := testConfig()
	_, err := keyFromString([]byte(config.KeyString))
	if err != nil {
		t.Error(err)
	}
}

func TestKeyFromFile(t *testing.T) {
	config := testConfig()
	_, err := keyFromFile(config.KeyPath)
	if err != nil {
		t.Error(err)
	}
}
