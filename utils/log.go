package utils

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
)

// SetupDefaults sets up the default logging level.
func SetupDefaults() {
	log.Root().SetHandler(
		log.LvlFilterHandler(
			log.LvlInfo,
			log.StreamHandler(os.Stdout, log.TerminalFormat(true)),
		),
	)
}
