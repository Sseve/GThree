package models

type Zone struct {
	Targt      string `json:"target" binding:"required"`
	ZId        string `json:"zid" binding:"required"`
	Ip         string `json:"ip" binding:"required"`
	Name       string `json:"name" binding:"required"`
	SvnVersion uint   `json:"svnversion"`
}

type ZoneOpt struct {
	Zone []Zone `json:"zone" binding:"required"`
}
