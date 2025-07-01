package model

import (
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
	TribeID              string `json:"tribeId,omitempty"`
}

func (s *Skills) ConvertSkillsFromProto(skills []*school.Skills)  {
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
