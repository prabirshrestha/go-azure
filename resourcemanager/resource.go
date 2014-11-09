package resourcemanager

import (
	"errors"
)

func NewResourceService(s *Service) *ResourceService {
	rs := &ResourceService{s: s}
	return rs
}

type ResourceService struct {
	s *Service
}

type ResourceListOpts struct {
	SubscriptionId string
	Top            int
	SkipToken      string
	Filter         string
	ApiVersion     string
}

func (rs *ResourceService) List(options *ResourceListOpts) (*ResourceService, error) {
	if options == nil {
		return nil, errors.New("options is nil")
	}

	if options.SubscriptionId == "" {
		return nil, errors.New("options.SubscriptionId is empty")
	}

	return nil, nil
}
