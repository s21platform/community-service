package model

import 	school "github.com/s21platform/school-proto/school-proto"

type Skill struct {
	Name   string `json:"name"`
	Points int32  `json:"points"`
}

type Badge struct {
	Name            string `json:"name"`
	ReceiptDateTime string `json:"receiptDateTime"`
	IconURL         string `json:"iconURL"`
}

type ParticipantData struct {
	ClassName            string  `json:"className"`
	ParallelName         string  `json:"parallelName"`
	ExpValue             int64   `json:"expValue"`
	Level                int32   `json:"level"`
	ExpToNextLevel       int64   `json:"expToNextLevel"`
	CampusUUID           string  `json:"campusUuid"`
	Status               string  `json:"status"`
	Skills               []Skill `json:"skills"`
	PeerReviewPoints     int64   `json:"peerReviewPoints"`
	PeerCodeReviewPoints int64   `json:"peerCodeReviewPoints"`
	Coins                int64   `json:"coins"`
	Badges               []Badge `json:"badges"`
}



func (p *ParticipantData) ToDTO(in *school.GetParticipantDataOut) (ParticipantData, error) {
	result := ParticipantData{
		AttributeId: in.AttributeId,
		Value:       in.Value,
		ParentId:    in.ParentId,
	}
	return result, nil
}
func (p *ParticipantData) ToDTO(in *school.GetParticipantDataOut) (ParticipantData, error){
	result := ParticipantData{
		ClassName:            in.ClassName,
		ParallelName:         in.ParallelName,
		ExpValue:             in.ExpValue,
		Level:                in.Level,
		ExpToNextLevel:       in.ExpToNextLevel,
		CampusUUID:           in.CampusUuid,
		Status:               in.Status,
		Skills:               convertSkillsToProto(p.Skills),
		PeerReviewPoints:     p.PeerReviewPoints,
		PeerCodeReviewPoints: p.PeerCodeReviewPoints,
		Coins:                p.Coins,
		Badges:               convertBadgesToProto(p.Badges),
	}
}
