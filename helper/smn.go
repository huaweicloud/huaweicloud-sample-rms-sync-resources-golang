package helper

import (
	"huaweicloud-sample-rms-sync-resources/config"
	"strings"
)

var regionToEndpoint map[string]string
var messageTypeSet map[string]bool

func init() {
	regionToEndpoint = make(map[string]string)
	cfg := config.GetConfig()
	for _, item := range cfg.Smn.Items {
		regionToEndpoint[item.RegionId] = item.Endpoint
	}

	messageTypeSet = map[string]bool{
		"Notification":             true,
		"SubscriptionConfirmation": true,
		"UnsubscribeConfirmation":  true,
	}
}

func CheckMessageType(messageType string) bool {
	return messageTypeSet[messageType]
}

func GetRegionIdFromTopicUrn(topicUrn string) string {
	if topicUrn == "" {
		return ""
	}
	topicContents := strings.Split(topicUrn, ":")
	if len(topicContents) == 5 {
		return topicContents[2]
	}
	return ""
}

func GetEndpointByRegionId(regionId string) string {
	return regionToEndpoint[regionId]
}
