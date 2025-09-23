package services

import (
	"cloud/internal/gcp"
	"cloud/internal/models"
	"fmt"
	"time"
)

type DeployService struct {
	gcpService *gcp.ComputeService
}

func NewDeployService(gcpService *gcp.ComputeService) *DeployService {
	return &DeployService{gcpService: gcpService}
}

func (s *DeployService) Deploy(req models.DeployRequest) (*models.DeployResponse, error) {
	// TODO: This is a stub - we'll implement this next
	deploymentID := fmt.Sprintf("deploy-%d", time.Now().Unix())

	return &models.DeployResponse{
		DeploymentID: deploymentID,
		Status:       "pending",
		VMID:         "vm-placeholder",
		CreatedAt:    time.Now(),
		Message:      "Deployment initiated successfully",
	}, nil
}

func (s *DeployService) GetVMStatus(vmID string) (*models.VMStatus, error) {
	// TODO: This is a stub - we'll implement this next
	return &models.VMStatus{
		VMID:      vmID,
		Status:    "RUNNING",
		PublicIP:  "35.123.456.789",
		UpdatedAt: time.Now(),
	}, nil
}
