package gcp

import (
	"fmt"
)

type ComputeService struct {
	// Placeholder for GCP Compute Service client
}

func NewComputeService() (*ComputeService, error) {
	//todo: implement authentication and client creation

	return &ComputeService{}, nil
}

func (cs *ComputeService) CreateVm(startupScript string) (string, error) {
	vmID := "vm-placeholder-123"
	fmt.Printf("Creating VM with script: %s\n", startupScript[:100]+"...")
	return vmID, nil
}
