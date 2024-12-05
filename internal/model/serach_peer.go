package model

type SearchPeers struct {
	Login string `db:"login"`
}

type PeerSchoolData struct {
	ClassName    string `db:"class_name"`
	ParallelName string `db:"parallel_name"`
}
