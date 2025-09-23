package api

import (
	"cloud/internal/models"
	"cloud/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeployHandler struct {
	deployService *services.DeployService
}

// ADD THIS CONSTRUCTOR FUNCTION:
func NewDeployHandler(deployService *services.DeployService) *DeployHandler {
	return &DeployHandler{
		deployService: deployService,
	}
}

func (h *DeployHandler) Deploy(c *gin.Context) {
	var req models.DeployRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "Invalid request payload",
			"error":   err.Error()})
		return

	}

	if req.RepoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"details": "Repository URL is required",
			"error":   "Missing repo_url field"})
		return
	}

	if req.RepoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "repo_url is required",
		})
		return
	}

	// Set defaults for VM config if not provided
	if req.VMConfig.MachineType == "" {
		req.VMConfig.MachineType = "e2-micro"
	}
	if req.VMConfig.Zone == "" {
		req.VMConfig.Zone = "us-central1-a"
	}
	if req.VMConfig.DiskSize == 0 {
		req.VMConfig.DiskSize = 10
	}

	resp, err := h.deployService.Deploy(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": "Failed to initiate deployment",
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)

}

func (h *DeployHandler) GetStatus(c *gin.Context) {
	vmID := c.Param("vm_id")

	if vmID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "vm_id is required",
		})
		return
	}

	status, err := h.deployService.GetVMStatus(vmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get VM status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, status)
}
