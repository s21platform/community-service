package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	school "github.com/s21platform/school-proto/school-proto"
)

type Skill struct {
	Name   string `json:"name"`
	Points int32  `json:"points"`
}
type Skills []Skill

type Badge struct {
	Name            string `json:"name"`
	ReceiptDateTime string `json:"receiptDateTime"`
	IconURL         string `json:"iconURL"`
}

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
	TribeID              string `json:"tribeId,omitempty"`
}

func (s *Skills) ConvertSkillsFromProto(skills []*school.Skills) []Skill {
	result := make([]Skill, len(skills))
	for i, skill := range skills {
		result[i] = Skill{
			Name:   skill.Name,
			Points: skill.Points,
		}
	}
	return result
}

func (s *Skills) Value() (driver.Value, error) {
	return Value(s)
}

func (s *Skills) Scan(value interface{}) error {
	return Scan(value, s)
}

type Badges []Badge

func (b *Badges) ConvertBadgesFromProto(badges []*school.Badges) []Badge {
	result := make([]Badge, len(badges))
	for i, badge := range badges {
		result[i] = Badge{
			Name:            badge.Name,
			ReceiptDateTime: badge.ReceiptDateTime,
			IconURL:         badge.IconURL,
		}
	}
	return result
}

func (b Badges) Value() (driver.Value, error) {
	return Value(b)
}

func (b *Badges) Scan(value interface{}) error {
	return Scan(value, b)
}

func Value(s interface{}) (driver.Value, error) {
	j, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func Scan(value interface{}, s interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return errors.New("failed to scan Skills, not string or []byte")
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, s)
}
