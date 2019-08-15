package notification

import (
	"encoding/json"
	"regexp"
	"time"

	"gopkg.in/go-playground/webhooks.v5/github"
)

type Notification struct {
	Repository string         `json:"repository"`
	Ref        string         `json:"ref"`
	CommitID   string         `json:"commit_id"`
	UpdatedAt  time.Time      `json:"updated_at"`
	SystemIDs  map[string]int `json:"system_ids"`
}

var regexID = regexp.MustCompile(`^[a-z0-9\-]{36}`)

func NewFromGithubWebhook(payload github.PushPayload) *Notification {
	notification := &Notification{
		Repository: payload.Repository.URL,
		Ref:        payload.Ref,
		CommitID:   payload.After,
		UpdatedAt:  payload.Repository.UpdatedAt,
		SystemIDs:  make(map[string]int),
	}

	for _, commit := range payload.Commits {
		for _, added := range commit.Added {
			if regexID.Match([]byte(added)) {
				system_id := string(regexID.Find([]byte(added)))
				notification.SystemIDs[system_id]++
			}
		}
		for _, removed := range commit.Removed {
			if regexID.Match([]byte(removed)) {
				system_id := string(regexID.Find([]byte(removed)))
				notification.SystemIDs[system_id]++
			}
		}
		for _, modified := range commit.Modified {
			if regexID.Match([]byte(modified)) {
				system_id := string(regexID.Find([]byte(modified)))
				notification.SystemIDs[system_id]++
			}
		}
	}

	return notification
}

func NewFromByteArray(b []byte) *Notification {
	var notification Notification

	json.Unmarshal(b, &notification)

	return &notification
}

func (notification *Notification) Contains(system_id string) bool {
	if _, ok := notification.SystemIDs[system_id]; ok {
		return true
	}
	return false
}
