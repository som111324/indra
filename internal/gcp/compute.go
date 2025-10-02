package gcp

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

type ComputeService struct {
	client    *compute.Service
	projectID string
	zone      string
}

type VMCreateRequest struct {
	Name          string
	MachineType   string
	Zone          string
	StartupScript string
	DiskSize      int64
}

func NewComputeService() (*ComputeService, error) {
	ctx := context.Background()

	// Load service account credentials
	client, err := google.DefaultClient(ctx, compute.ComputeScope)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCP client: %v", err)
	}

	// Create compute service
	service, err := compute.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %v", err)
	}

	return &ComputeService{
		client:    service,
		projectID: "indraai",
		zone:      "us-central1-a",
	}, nil
}

func (s *ComputeService) CreateVM(req VMCreateRequest) (string, error) {
	// Generate random VM name if not provided
	if req.Name == "" {
		req.Name = fmt.Sprintf("auto-deploy-%d", time.Now().Unix())
	}

	// Set defaults
	if req.MachineType == "" {
		req.MachineType = "e2-micro"
	}
	if req.Zone == "" {
		req.Zone = s.zone
	}
	if req.DiskSize == 0 {
		req.DiskSize = 10
	}

	// Create VM instance configuration
	instance := &compute.Instance{
		Name:        req.Name,
		MachineType: fmt.Sprintf("zones/%s/machineTypes/%s", req.Zone, req.MachineType),
		Zone:        req.Zone,

		// Boot disk configuration
		Disks: []*compute.AttachedDisk{
			{
				Boot:       true,
				AutoDelete: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "projects/ubuntu-os-cloud/global/images/family/ubuntu-2204-lts",
					DiskSizeGb:  req.DiskSize,
				},
			},
		},

		// Network configuration
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				Network: fmt.Sprintf("projects/%s/global/networks/default", s.projectID),
				AccessConfigs: []*compute.AccessConfig{
					{
						Type: "ONE_TO_ONE_NAT",
						Name: "External NAT",
					},
				},
			},
		},

		// Metadata with startup script
		Metadata: &compute.Metadata{
			Items: []*compute.MetadataItems{
				{
					Key:   "startup-script",
					Value: &req.StartupScript,
				},
			},
		},

		// Service account for VM
		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: "default",
				Scopes: []string{
					"https://www.googleapis.com/auth/cloud-platform",
				},
			},
		},

		// Tags for firewall rules
		Tags: &compute.Tags{
			Items: []string{"http-server", "https-server", "auto-deployer"},
		},
	}

	// Create the VM
	op, err := s.client.Instances.Insert(s.projectID, req.Zone, instance).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create VM: %v", err)
	}

	fmt.Printf("VM creation initiated. Operation: %s\n", op.Name)
	return req.Name, nil
}

// GetVMStatus returns the current status of a VM
func (s *ComputeService) GetVMStatus(vmName, zone string) (*compute.Instance, error) {
	if zone == "" {
		zone = s.zone
	}

	instance, err := s.client.Instances.Get(s.projectID, zone, vmName).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get VM status: %v", err)
	}

	return instance, nil
}

// DeleteVM deletes a VM
func (s *ComputeService) DeleteVM(vmName, zone string) error {
	if zone == "" {
		zone = s.zone
	}

	_, err := s.client.Instances.Delete(s.projectID, zone, vmName).Do()
	if err != nil {
		return fmt.Errorf("failed to delete VM: %v", err)
	}

	return nil
}

// Helper function to generate random string
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
