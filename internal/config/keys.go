package config

type key string

const (
	KeyMetrics                = key("metrics")
	KeyUUID                   = key("uuid")
	KeyLogger                 = key("logger")
	KeyParticipantLastUpdated = key("participant_last_updated")
	KeyCampusesLastUpdated    = key("campuses_last_update")
	KeyParticipantDataLastUpdated = key("participant_data_last_updated")
)
