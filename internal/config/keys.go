package config

type key string

const (
	KeyMetrics             = key("metrics")
	KeyUUID                = key("uuid")
	KeyLogger              = key("logger")
	KeyLoginsLastUpdated   = "logins_last_updated"
	KeyCampusesLastUpdated = "campuses_last_update"
)
