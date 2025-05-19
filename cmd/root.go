package cmd

import (
	"context"
	"fmt"
	"gha-register-build-artifact/internal/artifacts"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:   "gha-register-build-action",
		Short: "Publish the build artifact metadata to CloudBees Build Platform",
		Long:  "Publish the build artifact metadata to CloudBees Build Platform",
		RunE:  run,
	}
	cfg artifacts.Config
)

func Execute() error {
	return cmd.Execute()
}

func init() {
	setDefaultValues(&cfg)
}

func setDefaultValues(cfg *artifacts.Config) {
	artifactType := os.Getenv(artifacts.ArtifactType)
	if artifactType != "" {
		cfg.ArtifactType = artifactType
	} else {
		cfg.ArtifactType = ""
	}

	artifactDigest := os.Getenv(artifacts.ArtifactDigest)
	if artifactDigest != "" {
		cfg.ArtifactDigest = artifactDigest
	} else {
		cfg.ArtifactDigest = ""
	}

	artifactLabel := os.Getenv(artifacts.ArtifactLabel)
	if artifactLabel != "" {
		cfg.ArtifactLabel = artifactLabel
	} else {
		cfg.ArtifactLabel = ""
	}
	cfg.ArtifactOperation = artifacts.PUBLISHED
}

func run(_ *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown arguments: %v", args)
	}
	newContext, cancel := context.WithCancel(context.Background())
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, os.Interrupt)
	go func() {
		<-osChannel
		cancel()
	}()

	return cfg.Run(newContext)
}
