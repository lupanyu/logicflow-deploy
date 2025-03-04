package services

import "encoding/json"

type DeployRequest struct {
	FlowData json.RawMessage `json:"flow_data" binding:"required"`
	EnvType  string          `json:"env_type" binding:"required,oneof=dev test prod"`
	Operator string          `json:"operator" binding:"required"`
}

type DeployResponse struct {
	DeployID  string `json:"deploy_id"`
	Timestamp int64  `json:"timestamp"`
}

type DeployStatus struct {
	Status     string            `json:"status"`
	Progress   float32           `json:"progress"`
	NodesState map[string]string `json:"nodes_state"`
}

type DeployService struct{}

func (service *DeployService) Deploy() error {
	return nil
}
func NewDeployService() *DeployService {
	return &DeployService{}
}

func (s *DeployService) GetStatus(id string) string {
	return "success"
}

func (service *DeployService) DeployStatus(id string) error {
	return nil
}
func (s *DeployService) CancelDeploy(id string) error {
	return nil
}
