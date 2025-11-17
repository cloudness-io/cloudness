package server

import (
	"fmt"
	"os"
	"unicode"

	"github.com/cloudness-io/cloudness/lock"
	"github.com/cloudness-io/cloudness/pubsub"
	"github.com/cloudness-io/cloudness/store/database"
	"github.com/cloudness-io/cloudness/types"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func LoadConfig() (*types.Config, error) {
	config := new(types.Config)
	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	config.InstanceID, err = getSanitizedMachineName()
	if err != nil {
		return nil, fmt.Errorf("unable to ensure that instance ID is set in config: %w", err)
	}

	config.Process()

	return config, nil
}

// getSanitizedMachineName gets the name of the machine and returns it in sanitized format.
func getSanitizedMachineName() (string, error) {
	// use the hostname as default id of the instance
	hostName, err := os.Hostname()
	if err != nil {
		return "", err
	}

	// Always cast to lower and remove all unwanted chars
	// NOTE: this could theoretically lead to overlaps, then it should be passed explicitly
	// NOTE: for k8s names/ids below modifications are all noops
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/

	// The following code will:
	// * remove invalid runes
	// * remove diacritical marks (ie "smörgåsbord" to "smorgasbord")
	// * lowercase A-Z to a-z
	// * leave only a-z, 0-9, '-', '.' and replace everything else with '_'
	hostName, _, err = transform.String(
		transform.Chain(
			norm.NFD,
			runes.ReplaceIllFormed(),
			runes.Remove(runes.In(unicode.Mn)),
			runes.Map(func(r rune) rune {
				switch {
				case 'A' <= r && r <= 'Z':
					return r + 32
				case 'a' <= r && r <= 'z':
					return r
				case '0' <= r && r <= '9':
					return r
				case r == '-', r == '.':
					return r
				default:
					return '_'
				}
			}),
			norm.NFC),
		hostName)
	if err != nil {
		return "", err
	}

	return hostName, nil
}

// ProvideDatabaseConfig loads the database config from the main config.
func ProvideDatabaseConfig(config *types.Config) database.Config {
	return database.Config{
		Driver:     config.Database.Driver,
		Datasource: config.Database.Datasource,
	}
}

// ProvideLockConfig generates the `lock` package config from the config.
func ProvideLockConfig(config *types.Config) lock.Config {
	return lock.Config{
		App:           config.Lock.AppNamespace,
		Namespace:     config.Lock.DefaultNamespace,
		Provider:      config.Lock.Provider,
		Expiry:        config.Lock.Expiry,
		Tries:         config.Lock.Tries,
		RetryDelay:    config.Lock.RetryDelay,
		DriftFactor:   config.Lock.DriftFactor,
		TimeoutFactor: config.Lock.TimeoutFactor,
	}
}

func ProvidePubSubConfig(config *types.Config) pubsub.Config {
	return pubsub.Config{
		App:            config.PubSub.AppNamespace,
		Namespace:      config.PubSub.DefaultNamespace,
		Provider:       config.PubSub.Provider,
		HealthInterval: config.PubSub.HealthInterval,
		SendTimeout:    config.PubSub.SendTimeout,
		ChannelSize:    config.PubSub.ChannelSize,
	}
}
