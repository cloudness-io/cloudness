package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Kubectl struct {
	config *Config
	log    *Logger
}

func NewKubectl(cfg *Config, log *Logger) *Kubectl {
	return &Kubectl{
		config: cfg,
		log:    log,
	}
}

// =============================================================================
// kubectl-based Apply
// =============================================================================

// ApplyYAMLFile applies a YAML file using kubectl
func (k *Kubectl) ApplyYAMLFile(ctx context.Context, filepath string) error {
	// Check if file exists and has content
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) || (err == nil && info.Size() == 0) {
		k.log.Debug("Skipping empty or non-existent file: %s", filepath)
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", filepath, err)
	}

	// Use kubectl apply with retry
	return k.kubectlApplyWithRetry(ctx, filepath, 3)
}

// kubectlApplyWithRetry runs kubectl apply with exponential backoff
func (k *Kubectl) kubectlApplyWithRetry(ctx context.Context, filepath string, maxRetries int) error {
	var lastErr error
	backoff := 1 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			k.log.Debug("Retry %d/%d after %v...", attempt, maxRetries, backoff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
		}

		err := k.kubectlApply(ctx, filepath)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry on certain errors
		if strings.Contains(err.Error(), "invalid") ||
			strings.Contains(err.Error(), "forbidden") ||
			strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	return fmt.Errorf("kubectl apply failed after %d retries: %w", maxRetries, lastErr)
}

// kubectlApply runs kubectl apply -f
func (k *Kubectl) kubectlApply(ctx context.Context, filepath string) error {
	args := []string{"apply", "-f", filepath}

	if k.config.Verbose {
		k.log.Debug("kubectl %s", strings.Join(args, " "))
	}

	cmd := exec.CommandContext(ctx, "kubectl", args...)
	output, err := cmd.CombinedOutput()

	if k.config.Verbose && len(output) > 0 {
		k.log.Info("%s", strings.TrimSpace(string(output)))
	}

	if err != nil {
		return fmt.Errorf("kubectl apply failed: %s: %w", strings.TrimSpace(string(output)), err)
	}

	return nil
}

// Delete runs kubectl delete
func (k *Kubectl) Delete(ctx context.Context, resource, name, namespace string) error {
	args := []string{"delete", resource, name, "-n", namespace, "--ignore-not-found=true"}

	cmd := exec.CommandContext(ctx, "kubectl", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		k.log.Debug("kubectl delete failed: %s", strings.TrimSpace(string(output)))
		return err
	}

	return nil
}
