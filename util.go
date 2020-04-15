package acmeserverless

import (
	"github.com/fatih/structs"
)

// ToSentryMap converts the given struct to a map[string]interface{} so it can be sent to Sentry.
// The keys of the map are the same as the JSON element names. It panics if s's kind is not struct.
func ToSentryMap(i interface{}) map[string]interface{} {
	return structs.Map(i)
}
