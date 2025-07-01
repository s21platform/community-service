package config

type Key string

const (
	KeyMetrics                = Key("metrics")
	KeyUUID                   = Key("uuid")
	KeyLogger                 = Key("logger")
	KeyParticipantLastUpdated = Key("participant_last_updated")
	KeyCampusesLastUpdated    = Key("campuses_last_update")
	KeyParticipantDataLastUpdated = Key("participant_data_last_updated")
)
