package client

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	rms "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rms/v1"
	my_config "huaweicloud-sample-rms-sync-resources/config"
	"sync"
)

var singletonRmsClient *rms.RmsClient
var once sync.Once

func GetRmsClient() *rms.RmsClient {
	once.Do(buildRmsClient)
	return singletonRmsClient
}

func buildRmsClient() {
	cfg := my_config.GetConfig()
	domain := cfg.Domain
	service := cfg.Rms
	credential := getCredential(domain.Ak, domain.Sk, domain.DomainId)
	singletonRmsClient = rms.NewRmsClient(
		rms.RmsClientBuilder().
			WithEndpoint(service.Endpoint).
			WithCredential(credential).
			WithHttpConfig(config.DefaultHttpConfig()).
			Build())
}

func getCredential(ak, sk, domainId string) auth.ICredential {
	return global.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		WithDomainId(domainId).
		Build()
}
