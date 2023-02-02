package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
// Accepts variadic network IDs used to filter tokens.
func UnmarshalTokenConfig(contents []byte, networks ...uint64) ([]configuration.Token, error) {
	// Try to parse the file as a list of tokens.
	var payload []configuration.Token
	if err := json.Unmarshal(contents, &payload); err == nil {
		return FilterNetworks(payload, networks...), nil
	}

	// If this fails, try the backwards-compatible token json format
	var outer map[string]interface{}
	if err := json.Unmarshal(contents, &outer); err == nil {
		for k, v := range outer {
			for t, b := range v.(map[string]interface{}) {
				if b.(bool) {
					switch strings.ToLower(k) {
					case "mainnet":
						payload = append(payload, configuration.Token{
							ChainID: 1,
							Address: t,
						})
					case "testnet":
						payload = append(payload, configuration.Token{
							ChainID: 420,
							Address: t,
						})
					case "goerli":
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
		return FilterNetworks(payload, networks...), nil
	}

	return nil, fmt.Errorf("error parsing file contents as json token list")
}

// FilterNetworks filters a list of [configuration.Token] by network ID.
func FilterNetworks(tokens []configuration.Token, networks ...uint64) []configuration.Token {
	if len(networks) == 0 {
		return tokens
	}

	filtered := []configuration.Token{}
	for _, token := range tokens {
		for _, network := range networks {
			if token.ChainID == uint64(network) {
				filtered = append(filtered, token)
			}
		}
	}

	return filtered
}
