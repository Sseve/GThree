package models

type Zone struct {
	Targt string `json:"target" binding:"required"`
	ZId   uint32 `json:"zid" binding:"required"`
	Ip    string `json:"ip" binding:"required"`
	Name  string `json:"name" binding:"required"`
}

type ZoneOpt struct {
	Zone []Zone `json:"zone" binding:"required"`
}
