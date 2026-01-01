package main

import "context"

type Operator interface {
	// ApplyCommon applies common artifact files
	ApplyCommon(ctx context.Context) error

	// Volumes handles PVC creation and resizing
	Volumes(ctx context.Context) error

	// Deploy deploys the actual application
	Deploy(ctx context.Context) error

	// Ingress deploys the ingress
	Ingress(ctx context.Context) error

	// Cleanup cleans up resources
	Cleanup(ctx context.Context)
}
