package config

type key string

const (
	KeyMetrics                = key("metrics")
	KeyUUID                   = key("uuid")
	KeyLogger                 = key("logger")
	KeyParticipantLastUpdated = "participant_last_updated"
)
