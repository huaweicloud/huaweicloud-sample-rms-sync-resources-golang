package core

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"huaweicloud-sample-rms-sync-resources/config"
	"huaweicloud-sample-rms-sync-resources/db"
	"huaweicloud-sample-rms-sync-resources/helper"
	"huaweicloud-sample-rms-sync-resources/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const LAYOUT = "2006-01-02T15:04:05Z"

var memory_cache = cache.New(30*time.Minute, 35*time.Minute)

var resourceTypeSet map[string]bool

func init() {
	cfg := config.GetConfig()
	resourceTypeSet = make(map[string]bool)
	for _, rt := range cfg.Rms.ResourceTypes {
		resourceTypeSet[rt] = true
	}
}

func HandleSmnMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Failed to read response body when handling SMN message, err: %v", err)
		return
	}
	smnMessage := &models.SmnMessage{}
	if err := json.Unmarshal(body, smnMessage); err != nil {
		logrus.Errorf("Failed to load SMN message when handling SMN message, err: %v", err)
		return
	}
	if !isSmnMessageValid(smnMessage) {
		logrus.Error("Find that message is not valid when handling SMN message")
		return
	}
	switch smnMessage.Type {
	case "SubscriptionConfirmation":
		if err := handleSubscriptionConfirmation(smnMessage); err != nil {
			logrus.Errorf("Failed to handle subscription confirmation message, err: %v", err)
		}
	case "Notification":
		if err := handleNotification(smnMessage); err != nil {
			logrus.Errorf("Failed to handle notification message, err: %v", err)
		}
	default:
		logrus.Warnf("Received %v SMN message, but we won't deal with it", smnMessage.Type)
	}
}

func handleSubscriptionConfirmation(smnMessage *models.SmnMessage) error {
	logrus.Infof("Received SMN subscription confirmation message, full body: %v", *smnMessage)
	resp, err := http.Get(smnMessage.SubscribeUrl)
	if err != nil {
		return fmt.Errorf("Failed to invoke subscribe url, err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		logrus.Info("Handle SMN confirmation message successfully")
	} else {
		return fmt.Errorf("Failed to confirm SMN confirmation message")
	}
	return nil
}

func handleNotification(smnMessage *models.SmnMessage) error {
	logrus.Infof("Received SMN notification message, full body: %v", *smnMessage)

	message := smnMessage.Message
	notification := &models.Notification{}
	if err := json.Unmarshal([]byte(message), notification); err != nil {
		return fmt.Errorf("Failed to load notification, err: %v", err)
	}
	notificationType := notification.NotificationType
	switch notificationType {
	case "ResourceChanged":
		logrus.Infof("Start to deal with %v notification", notificationType)
		return handleResourceChangedNotification(notification)
	case "ResourceRelationChanged":

		logrus.Warnf("We are not gonna deal with %v notification", notificationType)
	case "SnapshotArchiveCompleted":
		logrus.Warnf("We are not gonna deal with %v notification", notificationType)
	case "NotificationArchiveCompleted":
		logrus.Warnf("We are not gonna deal with %v notification", notificationType)
	default:
		return fmt.Errorf("Unexpected notification type: %v", notification.NotificationType)
	}
	return nil
}

func handleResourceChangedNotification(notification *models.Notification) error {
	detail := &models.ResourceChangedDetail{}
	if err := mapstructure.Decode(notification.Detail, detail); err != nil {
		return fmt.Errorf("Failed to decode notification detail, err: %v", err)
	}
	if !resourceTypeSet[detail.ResourceType] {
		logrus.Warnf("ResourceType %v is not configured, so we are not gonna deal with it", detail.ResourceType)
		return nil
	}
	resourceFromNotification := detail.Resource
	queryAt, err := time.Parse(LAYOUT, detail.CaptureTime)
	if err != nil {
		return fmt.Errorf("Failed to parse capture time, err: %v", err)
	}
	resourceFromNotification.QueryAt = queryAt

	if detail.EventType == models.EventDelete {
		if err := db.DeleteOneResource(resourceFromNotification); err != nil {
			return err
		}
	} else if detail.EventType == models.EventCreate || detail.EventType == models.EventUpdate {
		resourceFromNotification.State = models.ResourceStateNormal
		if err := db.CreateOrUpdateOneResource(resourceFromNotification); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unexpected event type, %v", detail.EventType)
	}
	return nil
}

func isSmnMessageValid(smnMessage *models.SmnMessage) bool {
	if !helper.CheckMessageType(smnMessage.Type) {
		logrus.Error("Message type is not valid while validating SMN message")
		return false
	}
	regionId := helper.GetRegionIdFromTopicUrn(smnMessage.TopicUrn)
	if regionId == "" {
		logrus.Error("Empty regionId while validating SMN message")
		return false
	}
	endpoint := helper.GetEndpointByRegionId(regionId)
	if endpoint == "" {
		logrus.Error("Empty endpoint while validating SMN message")
		return false
	}
	certificate, err1 := getCertificate(smnMessage.SigningCertUrl)
	if err1 != nil {
		logrus.Error(err1)
		return false
	}
	block, _ := pem.Decode(certificate)
	if block == nil {
		logrus.Errorf("Failed to parse certificate PEM")
		return false
	}
	cert, err2 := x509.ParseCertificate(block.Bytes)
	if err2 != nil {
		logrus.Errorf("Failed to parse certificate, err: %v", err2)
		return false
	}
	signed := []byte(buildSignMessage(smnMessage))
	signature, err3 := base64.StdEncoding.DecodeString(smnMessage.Signature)
	if err3 != nil {
		logrus.Errorf("Failed to decode signature, err: %v", err3)
		return false
	}
	err4 := cert.CheckSignature(cert.SignatureAlgorithm, signed, signature)
	if err4 != nil {
		logrus.Errorf("Failed to check signature, err: %v", err4)
		return false
	}
	return true
}

func buildSignMessage(smnMessage *models.SmnMessage) string {
	messageType := smnMessage.Type
	var message string
	if messageType == "Notification" {
		message = buildNotificationMessage(smnMessage)
	} else if messageType == "SubscriptionConfirmation" || messageType == "UnsubscribeConfirmation" {
		message = buildSubscriptionMessage(smnMessage)
	}
	return message
}

func buildNotificationMessage(smnMessage *models.SmnMessage) string {
	var sb strings.Builder
	sb.WriteString("message\n")
	sb.WriteString(smnMessage.Message)
	sb.WriteString("\nmessage_id\n")
	sb.WriteString(smnMessage.MessageId)
	if smnMessage.Subject != "" {
		sb.WriteString("\nsubject\n")
		sb.WriteString(smnMessage.Subject)
	}
	sb.WriteString("\ntimestamp\n")
	sb.WriteString(smnMessage.Timestamp)
	sb.WriteString("\ntopic_urn\n")
	sb.WriteString(smnMessage.TopicUrn)
	sb.WriteString("\ntype\n")
	sb.WriteString(smnMessage.Type)
	sb.WriteString("\n")
	return sb.String()
}

func buildSubscriptionMessage(smnMessage *models.SmnMessage) string {
	var sb strings.Builder
	sb.WriteString("message\n")
	sb.WriteString(smnMessage.Message)
	sb.WriteString("\nmessage_id\n")
	sb.WriteString(smnMessage.MessageId)
	sb.WriteString("\nsubscribe_url\n")
	sb.WriteString(smnMessage.SubscribeUrl)
	sb.WriteString("\ntimestamp\n")
	sb.WriteString(smnMessage.Timestamp)
	sb.WriteString("\ntopic_urn\n")
	sb.WriteString(smnMessage.TopicUrn)
	sb.WriteString("\ntype\n")
	sb.WriteString(smnMessage.Type)
	sb.WriteString("\n")
	return sb.String()
}

func getCertificate(certUrl string) ([]byte, error) {
	cert, found := memory_cache.Get("certificate")
	if found {
		return cert.([]byte), nil
	}
	resp, err := http.Get(certUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke getting certificate url, err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read certificate, err: %v", err)
		}
		memory_cache.Set("certificate", body, cache.DefaultExpiration)
		return body, nil
	}
	return nil, fmt.Errorf("Failed to download certificate")
}
