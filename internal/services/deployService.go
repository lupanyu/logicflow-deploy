package services

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
