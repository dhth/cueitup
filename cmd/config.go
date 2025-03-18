package cmd

import (
	"errors"
	"fmt"
	"strings"

	t "github.com/dhth/cueitup/internal/types"
	yaml "github.com/goccy/go-yaml"
)

var (
	errCouldntParseConfig   = errors.New("couldn't parse config file")
	errProfileNotFound      = errors.New("profile not found")
	errNoProfilesDefined    = errors.New("no profiles defined")
	errProfileConfigInvalid = errors.New("profile config is invalid")
)

func getProfile(configBytes []byte, profileName string) (t.Profile, error) {
	var cfg t.Config
	var zero t.Profile

	err := yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return zero, fmt.Errorf("%w: %s", errCouldntParseConfig, err.Error())
	}

	if len(cfg.Profiles) == 0 {
		return zero, errNoProfilesDefined
	}

	availableProfiles := make([]string, len(cfg.Profiles))
	for i, pc := range cfg.Profiles {
		availableProfiles[i] = pc.Name
		if pc.Name == profileName {
			profile, errors := t.ParseProfileConfig(pc)
			if len(errors) > 0 {
				if len(errors) == 1 {
					return zero, fmt.Errorf("%w: %s", errProfileConfigInvalid, errors[0].Error())
				}

				errorStrs := make([]string, len(errors))
				for i, err := range errors {
					errorStrs[i] = fmt.Sprintf("  - %s", err.Error())
				}
				return zero, fmt.Errorf("%w:\n%s", errProfileConfigInvalid, strings.Join(errorStrs, "\n"))
			}

			return profile, nil
		}
	}

	return zero, fmt.Errorf("%w; available profiles: %v", errProfileNotFound, availableProfiles)
}
