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
	CreatedAt  time.Time      `json:"created_at"`
	PushedAt   time.Time      `json:"pushed_at"`
	SystemIDs  map[string]int `json:"system_ids"`
}

var regexID = regexp.MustCompile(`^[a-z0-9\-]{36}`)

func NewFromGithubWebhook(payload github.PushPayload) *Notification {
	n := &Notification{
		Repository: payload.Repository.URL,
		Ref:        payload.Ref,
		CommitID:   payload.After,
		UpdatedAt:  payload.Repository.UpdatedAt,
		CreatedAt:  time.Unix(payload.Repository.CreatedAt, 0),
		PushedAt:   time.Unix(payload.Repository.PushedAt, 0),
		SystemIDs:  make(map[string]int),
	}

	for _, commit := range payload.Commits {
		for _, added := range commit.Added {
			if regexID.Match([]byte(added)) {
				system_id := string(regexID.Find([]byte(added)))
				n.SystemIDs[system_id]++
			}
		}
		for _, removed := range commit.Removed {
			if regexID.Match([]byte(removed)) {
				system_id := string(regexID.Find([]byte(removed)))
				n.SystemIDs[system_id]++
			}
		}
		for _, modified := range commit.Modified {
			if regexID.Match([]byte(modified)) {
				system_id := string(regexID.Find([]byte(modified)))
				n.SystemIDs[system_id]++
			}
		}
	}

	return n
}

func NewFromJson(b []byte) *Notification {
	var n Notification

	json.Unmarshal(b, &n)

	return &n
}

func (n *Notification) Contains(system_id string) bool {
	if _, ok := n.SystemIDs[system_id]; ok {
		return true
	}
	return false
}

func (n *Notification) IsSelected(system_id, ref string) bool {
	_, ok := n.SystemIDs[system_id]

	// empty system_id or a found system_id is true
	SystemID := (system_id == "") || ok

	// empty ref or a matching ref is true
	Ref := (ref == "") || (ref == n.Ref)

	return SystemID && Ref
}

func (n *Notification) ToJson() (b []byte, err error) {
	b, err = json.Marshal(n)
	return b, err
}
