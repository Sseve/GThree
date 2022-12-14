package models

type Zone struct {
	Targt      string `json:"target" binding:"required"`
	Zid        string `json:"zid" binding:"required"`
	Ip         string `json:"ip" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Closed     bool   `json:"closed"`
	SvnVersion uint64 `json:"svnversion"` // 进行区服svn操作时,此字段是必须字段
}

type ZoneOpt struct {
	Zone []Zone `json:"zone" binding:"required"`
}
