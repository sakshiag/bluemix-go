package cfv2

import (
	"fmt"

	"github.com/IBM-Bluemix/bluemix-go/client"
	"github.com/IBM-Bluemix/bluemix-go/rest"
)

//ServiceInstanceCreateRequest ...
type ServiceInstanceCreateRequest struct {
	Name      string                 `json:"name"`
	SpaceGUID string                 `json:"space_guid"`
	PlanGUID  string                 `json:"service_plan_guid"`
	Params    map[string]interface{} `json:"parameters,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
}

//ServiceInstanceUpdateRequest ...
type ServiceInstanceUpdateRequest struct {
	Name     *string                `json:"name,omitempty"`
	PlanGUID *string                `json:"service_plan_guid,omitempty"`
	Params   map[string]interface{} `json:"parameters,omitempty"`
	Tags     *[]string              `json:"tags,omitempty"`
}

//ServiceInstance ...
type ServiceInstance struct {
	GUID              string
	Name              string                 `json:"name"`
	Credentials       map[string]interface{} `json:"credentials"`
	ServicePlanGUID   string                 `json:"service_plan_guid"`
	SpaceGUID         string                 `json:"space_guid"`
	GatewayData       string                 `json:"gateway_data"`
	Type              string                 `json:"type"`
	DashboardURL      string                 `json:"dashboard_url"`
	LastOperation     LastOperationFields    `json:"last_operation"`
	RouteServiceURL   string                 `json:"routes_url"`
	Tags              []string               `json:"tags"`
	SpaceURL          string                 `json:"space_url"`
	ServicePlanURL    string                 `json:"service_plan_url"`
	ServiceBindingURL string                 `json:"service_bindings_url"`
	ServiceKeysURL    string                 `json:"service_keys_url"`
}

//ServiceInstanceFields ...
type ServiceInstanceFields struct {
	Metadata ServiceInstanceMetadata
	Entity   ServiceInstance
}

//ServiceInstanceMetadata ...
type ServiceInstanceMetadata struct {
	GUID string `json:"guid"`
	URL  string `json:"url"`
}

//LastOperationFields ...
type LastOperationFields struct {
	Type        string `json:"type"`
	State       string `json:"state"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

//ServiceInstanceResource ...
type ServiceInstanceResource struct {
	Resource
	Entity ServiceInstanceEntity
}

//ServiceInstanceEntity ...
type ServiceInstanceEntity struct {
	Name              string                 `json:"name"`
	Credentials       map[string]interface{} `json:"credentials"`
	ServicePlanGUID   string                 `json:"service_plan_guid"`
	SpaceGUID         string                 `json:"space_guid"`
	GatewayData       string                 `json:"gateway_data"`
	Type              string                 `json:"type"`
	DashboardURL      string                 `json:"dashboard_url"`
	LastOperation     LastOperationFields    `json:"last_operation"`
	RouteServiceURL   string                 `json:"routes_url"`
	Tags              []string               `json:"tags"`
	SpaceURL          string                 `json:"space_url"`
	ServicePlanURL    string                 `json:"service_plan_url"`
	ServiceBindingURL string                 `json:"service_bindings_url"`
	ServiceKeysURL    string                 `json:"service_keys_url"`
}

//ToModel ...
func (resource ServiceInstanceResource) ToModel() ServiceInstance {

	entity := resource.Entity

	return ServiceInstance{
		GUID:              resource.Metadata.GUID,
		Name:              entity.Name,
		Credentials:       entity.Credentials,
		ServicePlanGUID:   entity.ServicePlanGUID,
		SpaceGUID:         entity.SpaceGUID,
		GatewayData:       entity.GatewayData,
		Type:              entity.Type,
		LastOperation:     entity.LastOperation,
		RouteServiceURL:   entity.RouteServiceURL,
		DashboardURL:      entity.DashboardURL,
		Tags:              entity.Tags,
		SpaceURL:          entity.SpaceURL,
		ServicePlanURL:    entity.ServicePlanURL,
		ServiceBindingURL: entity.ServiceBindingURL,
		ServiceKeysURL:    entity.ServiceKeysURL,
	}
}

//ServiceInstances ...
type ServiceInstances interface {
	Create(req ServiceInstanceCreateRequest) (*ServiceInstanceFields, error)
	Update(instanceGUID string, req ServiceInstanceUpdateRequest) (*ServiceInstanceFields, error)
	Delete(instanceGUID string) error
	FindByName(instanceName string) (*ServiceInstance, error)
	Get(instanceGUID string) (*ServiceInstanceFields, error)
}

type serviceInstance struct {
	client *client.Client
}

func newServiceInstanceAPI(c *client.Client) ServiceInstances {
	return &serviceInstance{
		client: c,
	}
}

func (s *serviceInstance) Create(req ServiceInstanceCreateRequest) (*ServiceInstanceFields, error) {
	rawURL := "/v2/service_instances?accepts_incomplete=true&async=true"
	serviceFields := ServiceInstanceFields{}
	_, err := s.client.Post(rawURL, req, &serviceFields)
	if err != nil {
		return nil, err
	}
	return &serviceFields, nil
}

func (s *serviceInstance) Get(instanceGUID string) (*ServiceInstanceFields, error) {
	rawURL := fmt.Sprintf("/v2/service_instances/%s", instanceGUID)
	serviceFields := ServiceInstanceFields{}
	_, err := s.client.Get(rawURL, &serviceFields)
	if err != nil {
		return nil, err
	}

	return &serviceFields, err
}

func (s *serviceInstance) FindByName(instanceName string) (*ServiceInstance, error) {
	req := rest.GetRequest("/v2/service_instances")
	req.Query("return_user_provided_service_instances", "true")
	if instanceName != "" {
		req.Query("q", "name:"+instanceName)
	}
	httpReq, err := req.Build()
	if err != nil {
		return nil, err
	}
	path := httpReq.URL.String()
	services, err := s.listServicesWithPath(path)
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("Service instance:  %q doesn't exist", instanceName)
	}
	return &services[0], nil
}

func (s *serviceInstance) Delete(instanceGUID string) error {
	rawURL := fmt.Sprintf("/v2/service_instances/%s", instanceGUID)
	_, err := s.client.Delete(rawURL)
	return err
}

func (s *serviceInstance) Update(instanceGUID string, req ServiceInstanceUpdateRequest) (*ServiceInstanceFields, error) {
	rawURL := fmt.Sprintf("/v2/service_instances/%s?accepts_incomplete=true&async=true", instanceGUID)
	serviceFields := ServiceInstanceFields{}
	_, err := s.client.Put(rawURL, req, &serviceFields)
	if err != nil {
		return nil, err
	}
	return &serviceFields, nil
}

func (s *serviceInstance) listServicesWithPath(path string) ([]ServiceInstance, error) {
	var services []ServiceInstance
	_, err := s.client.GetPaginated(path, ServiceInstanceResource{}, func(resource interface{}) bool {
		if serviceInstanceResource, ok := resource.(ServiceInstanceResource); ok {
			services = append(services, serviceInstanceResource.ToModel())
			return true
		}
		return false
	})
	return services, err
}
