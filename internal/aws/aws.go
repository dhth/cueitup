package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/dhth/cueitup/internal/types"

	"github.com/aws/aws-sdk-go-v2/config"
)

func GetAWSConfig(source types.ConfigSource) (aws.Config, error) {
	var cfg aws.Config
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch source.Kind {
	case types.Env:
		cfg, err = config.LoadDefaultConfig(ctx)
	case types.SharedProfile:
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithSharedConfigProfile(source.Value))
	}

	return cfg, err
}
