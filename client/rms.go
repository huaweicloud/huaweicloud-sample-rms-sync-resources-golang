package client

import (
	rms_model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rms/v1/model"
)

var rmsClient = GetRmsClient()

func ListLimitedResources(provider, resourceType, regionId string, limit int32, marker *string) (*rms_model.ListResourcesResponse, error) {
	req := rms_model.ListResourcesRequest{
		Provider: provider,
		Type:     resourceType,
		RegionId: &regionId,
		Limit:    &limit,
		Marker:   marker,
	}
	resp, err := rmsClient.ListResources(&req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ListProviders() (*rms_model.ListProvidersResponse, error) {
	req := rms_model.ListProvidersRequest{}
	resp, err := rmsClient.ListProviders(&req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetResourceById(provider, resourceType, resourceId string) (*rms_model.ShowResourceByIdResponse, error) {
	req := rms_model.ShowResourceByIdRequest{
		Provider:   provider,
		Type:       resourceType,
		ResourceId: resourceId,
	}
	resp, err := rmsClient.ShowResourceById(&req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
