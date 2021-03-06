package models

import (
	"time"
)

type Profile struct {
	UID        int              `json:"uid,string" bson:"_id"`
	Name       string           `json:"name" bson:"name"`
	Premium    bool             `json:"premium" bson:"premium"`
	Champion   bool             `json:"champion" bson:"champion"`
	Picture    string           `json:"picture" bson:"picture"`
	Ranks      map[string]*Rank `json:"ranks" bson:"ranks"`
	LastUpdate int64            `json:"last_update,omitempty" bson:"last_update"` //local
}

//Last time data was downloaded from remote servers
func LastRemoteUpdate() time.Time {
	currentTime := time.Now().UTC()
	return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond(), currentTime.Location())
}

//Indicates whether new data may be available on remote servers
func (profile Profile) NeedsRefresh() bool {
	lastUpdate := time.Unix(profile.LastUpdate, 0).UTC()
	now := time.Now().UTC()
	if time.Since(lastUpdate).Hours() > 24 {
		return true
	} else if (lastUpdate.Hour() < 6 && lastUpdate.Day() == now.Day()) || (lastUpdate.Day() != now.Day() && now.Hour() >= 6) {
		return true
	}
	return false
}
