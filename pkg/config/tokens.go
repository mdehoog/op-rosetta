package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/coinbase/rosetta-geth-sdk/configuration"
)

// ReadTokenConfig attempts to read file contents from a filename.
func ReadTokenConfig(filename string) (string, error) {
	contents := ""
	if file, err := os.ReadFile(filename); err == nil {
		// if the envvar points to a file, read it; otherwise the envvar contents is expected to be JSON
		contents = string(file)
	}
	if contents == "" {
		contents = "[]"
	}
	return contents, nil
}

// UnmarshalTokenConfig attempts to construct a list of [configuration.Token] from a JSON file.
func UnmarshalTokenConfig(contents string) ([]configuration.Token, error) {
	// Try to parse the file as a list of tokens.
	var payload []configuration.Token
	if err := json.Unmarshal([]byte(contents), &payload); err == nil {
		return payload, nil
	}

	// If this fails, try the backwards-compatible token json format
	var outer map[string]interface{}
	if err := json.Unmarshal([]byte(contents), &outer); err == nil {
		for k, v := range outer {
			for t, b := range v.(map[string]interface{}) {
				if b.(bool) {
					switch k {
					case "Mainnet":
						payload = append(payload, configuration.Token{
							ChainID: 1,
							Address: t,
						})
					case "Testnet":
						payload = append(payload, configuration.Token{
							ChainID: 420,
							Address: t,
						})
					case "Goerli":
						payload = append(payload, configuration.Token{
							ChainID: 420,
							Address: t,
						})
					default:
						return nil, fmt.Errorf("unknown chain %s found when parsing json token list", k)
					}
				}
			}
		}
		return payload, nil
	}

	return nil, fmt.Errorf("error parsing file contents as json token list")
}
