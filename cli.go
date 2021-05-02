package weather

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// CurrentWeatherCLI accepts a slice of command line flags and arguments,
// determines the location of interest and the measurement units to use
//  (e.g. imperial, standard, metric) and prints the current weather conditions
// for that location using the given measurement units. An error is returned if
// the OPENWEATHER_API_KEY environment variable is not set, if the command line
// flags and arguments are invalid, or if the call to get the weather conditions
// has a problem.
func CurrentWeatherCLI(args []string) error {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return errors.New("environment variable OPENWEATHER_API_KEY must be set")
	}

	var cfg cliEnv
	if err := cfg.fromArgs(args[1:]); err != nil {
		return err
	}

	c, err := Conditions(cfg.location, cfg.units, apiKey)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", c)
	return nil
}

// cliEnv represents command line arguments and flags.
type cliEnv struct {
	units    string
	location string
}

// fromArgs accepts a slice of strings representing command line flags and
// positional arguments and tries to parse them into a cliEnv struct. An
// error is returned if the units flag cannot be parsed correctly or if the
// location positional parameter is not provided.
func (c *cliEnv) fromArgs(args []string) error {
	fs := flag.NewFlagSet("weather", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fs.Output().Write([]byte("USAGE: weather [-units={standard|metric|imperial}] <location>\n\n"))
		fs.PrintDefaults()
	}
	fs.StringVar(&c.units, "units", "imperial", "the units to use, one of: standard, metric, imperial")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if c.units != "imperial" && c.units != "standard" && c.units != "metric" {
		return errors.New("units flag must be set to one of: imperial, metric, standard")
	}
	loc := fs.Arg(0)
	if loc == "" {
		return errors.New("positional argument for location must be given (e.g. 'london', 'tampa,us', etc.)")
	}
	c.location = loc

	return nil
}
