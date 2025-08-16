package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	school "github.com/s21platform/school-proto/school-proto"
)

type Skill struct {
	Name   string `json:"name"`
	Points int32  `json:"points"`
}

type Badge struct {
	Name            string `json:"name"`
	ReceiptDateTime string `json:"receiptDateTime"`
	IconURL         string `json:"iconURL"`
}
type Skills []Skill
type Badges []Badge
type ParticipantDataValue struct {
	ClassName            string `json:"className"`
	ParallelName         string `json:"parallelName"`
	ExpValue             int64  `json:"expValue"`
	Level                int32  `json:"level"`
	ExpToNextLevel       int64  `json:"expToNextLevel"`
	CampusUUID           string `json:"campusUuid"`
	Status               string `json:"status"`
	Skills               Skills `json:"skills"`
	PeerReviewPoints     int64  `json:"peerReviewPoints"`
	PeerCodeReviewPoints int64  `json:"peerCodeReviewPoints"`
	Coins                int64  `json:"coins"`
	Badges               Badges `json:"badges"`
	TribeID              int64  `json:"tribeId"`
}

type ParticipantData struct {
	Login          string `json:"login"`
	CampusId       int64  `db:"campus_id" json:"campus_id"`
	ClassName      string `db:"class_name" json:"class_name"`
	ParallelName   string `db:"parallel_name" json:"parallel_name"`
	TribeID        int64  `db:"tribe_id" json:"tribe_id"`
	Status         string `json:"status"`
	CreatedAt      string `db:"created_at" json:"created_at"`
	ExpValue       int64  `db:"exp_value" json:"exp_value"`
	Level          int64  `json:"level"`
	ExpToNextLevel int64  `db:"exp_to_next_level" json:"exp_to_next_level"`
	Crp            int64  `json:"crp"`
	Skills         Skills `json:"skills"`
	Prp            int64  `json:"prp"`
	Coins          int64  `json:"coins"`
	Badges         Badges `json:"badges"`
}

type Participant struct {
	Login    string `db:"login"`
	ExpLevel int64  `db:"exp_level"`
	Level    int32  `db:"level"`
	Status   string `db:"status"`
}

const (
	ParticipantStatusActive   = "ACTIVE"
	ParticipantStatusBlocked  = "BLOCKED"
	ParticipantStatusExpelled = "EXPELLED"
	ParticipantStatusFrozen   = "FROZEN"
)

func (s *Skills) ConvertSkillsFromProto(skills []*school.Skills) {
	*s = make([]Skill, len(skills))
	for i, skill := range skills {
		(*s)[i] = Skill{
			Name:   skill.Name,
			Points: skill.Points,
		}
	}
}

func (b *Badges) ConvertBadgesFromProto(badges []*school.Badges) {
	*b = make([]Badge, len(badges))
	for i, badge := range badges {
		(*b)[i] = Badge{
			Name:            badge.Name,
			ReceiptDateTime: badge.ReceiptDateTime,
			IconURL:         badge.IconURL,
		}
	}
}

// Реализация driver.Valuer для Skills
func (s Skills) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Skills: %w", err)
	}
	return string(b), nil
}

// Реализация sql.Scanner для Skills
func (s *Skills) Scan(src interface{}) error {
	if src == nil {
		*s = Skills{}
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("cannot scan type %T into Skills", src)
	}
	return json.Unmarshal(data, s)
}

// Реализация driver.Valuer для Badges
func (b Badges) Value() (driver.Value, error) {
	if len(b) == 0 {
		return "[]", nil
	}
	data, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Badges: %w", err)
	}
	return string(data), nil
}

// Реализация sql.Scanner для Badges
func (b *Badges) Scan(src interface{}) error {
	if src == nil {
		*b = Badges{}
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("cannot scan type %T into Badges", src)
	}
	return json.Unmarshal(data, b)
}
