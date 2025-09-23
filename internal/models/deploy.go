package models

import "time"

type DeployRequest struct {
	RepoURL      string   `json:"repo_url" binding:"required"`
	ProjectType  string   `json:"project_type"`
	StartCommand string   `json:"start_command"`
	VMConfig     VMConfig `json:"vm_config"`
}

// VMConfig holds VM-specific configuration
type VMConfig struct {
	MachineType string `json:"machine_type"` // e.g., "e2-micro"
	Zone        string `json:"zone"`         // e.g., "us-central1-a"
	DiskSize    int64  `json:"disk_size"`    // GB, default 10
}

// DeployResponse is returned after initiating deployment
type DeployResponse struct {
	DeploymentID string    `json:"deployment_id"`
	Status       string    `json:"status"`
	VMID         string    `json:"vm_id"`
	CreatedAt    time.Time `json:"created_at"`
	Message      string    `json:"message"`
}

// VMStatus represents current state of the VM
type VMStatus struct {
	VMID      string    `json:"vm_id"`
	Status    string    `json:"status"` // "PROVISIONING", "RUNNING", "TERMINATED"
	PublicIP  string    `json:"public_ip"`
	UpdatedAt time.Time `json:"updated_at"`
}
