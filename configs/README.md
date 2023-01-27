## configs

A collection of configuration files used to test the `op-rosetta` client using [`rosetta-cli`](https://github.com/coinbase/rosetta-cli).

These config files were built using the [rosett-cli documentation](https://www.rosetta-api.org/docs/rosetta_test.html).


### Setup

First, [install `rosetta-cli`](https://github.com/coinbase/rosetta-cli#install).

If the installation script doesn't work (on mac it may not), try cloning the repo and installing from source:

```bash
gh repo clone coinbase/rosetta-cli
cd rosetta-cli
make install
```


### Testing

_NOTE: `op-rosetta` must be running on the specified host and port provided in the configuration file. For local testing, this can be done as described in the root [README.md](../README.md#running) section, which will run an instance on localhost, port 8080._

To validate with `rosetta-cli`, run one of the following commands:

* `rosetta-cli check:data --configuration-file optimism/config.json` - This command validates that the Data API implementation is correct using the ethereum `optimism` node. It also ensures that the implementation does not miss any balance-changing operations.

* `rosetta-cli check:construction --configuration-file optimism/config.json` - This command validates the Construction API implementation. It also verifies transaction construction, signing, and submissions to the `optimism` network.

* `rosetta-cli check:data --configuration-file optimism/config.json` - This command validates that the Data API implementation is correct using the ethereum `optimism` node. It also ensures that the implementation does not miss any balance-changing operations.


### Reference

The `rosetta-cli` offers a number of commands useful for checking the implementation of the `op-rosetta` client. For more information, run `rosetta-cli --help` for the following output:

```
CLI for the Rosetta API

Usage:
  rosetta-cli [command]

Available Commands:
  check:construction           Check the correctness of a Rosetta Construction API Implementation
  check:data                   Check the correctness of a Rosetta Data API Implementation
  check:perf                   Benchmark performance of time-critical endpoints of Asset Issuer's Rosetta Implementation
  check:spec                   Check that a Rosetta implementation satisfies Rosetta spec
  completion                   Generate the autocompletion script for the specified shell
  configuration:create         Create a default configuration file at the provided path
  configuration:validate       Ensure a configuration file at the provided path is formatted correctly
  help                         Help about any command
  utils:asserter-configuration Generate a static configuration file for the Asserter
  utils:train-zstd             Generate a zstd dictionary for enhanced compression performance
  version                      Print rosetta-cli version
  view:balance                 View an account balance
  view:block                   View a block
  view:networks                View all network statuses

Flags:
      --block-profile string        Save the pprof block profile in the specified file
      --configuration-file string   Configuration file that provides connection and test settings.
                                    If you would like to generate a starter configuration file (populated
                                    with the defaults), run rosetta-cli configuration:create.

                                    Any fields not populated in the configuration file will be populated with
                                    default values.
      --cpu-profile string          Save the pprof cpu profile in the specified file
  -h, --help                        help for rosetta-cli
      --mem-profile string          Save the pprof mem profile in the specified file

Use "rosetta-cli [command] --help" for more information about a command.
```