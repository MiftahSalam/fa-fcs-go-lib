package logger

type (
	Context struct {
		ExternalID     string                 `json:"external_id"`
		JourneyID      string                 `json:"journey_id"`
		ChainID        string                 `json:"chain_id"`
		AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
	}

	ctxKeyLogger struct{}
)

var ctxKey = ctxKeyLogger{}
