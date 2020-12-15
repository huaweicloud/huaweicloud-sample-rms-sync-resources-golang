package models

type SmnMessage struct {
	Type             string `json:"type"`
	TopicUrn         string `json:"topic_urn"`
	MessageId        string `json:"message_id"`
	Message          string `json:"message"`
	UnsubscribeUrl   string `json:"unsubscribe_url"`
	SubscribeUrl     string `json:"subscribe_url"`
	Signature        string `json:"signature"`
	SignatureVersion string `json:"signature_version"`
	SigningCertUrl   string `json:"signing_cert_url"`
	Timestamp        string `json:"timestamp"`
	Subject          string `json:"subject"`
}

type ConfirmBody struct {
	Token    string `json:"token"`
	TopicUrn string `json:"topic_urn"`
}
