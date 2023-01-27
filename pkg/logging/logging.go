package logging

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
)

// SetupLogging sets up the default logging level.
func SetupLogging() {
	log.Root().SetHandler(
		log.LvlFilterHandler(
			log.LvlInfo,
			log.StreamHandler(os.Stdout, log.TerminalFormat(true)),
		),
	)
}
