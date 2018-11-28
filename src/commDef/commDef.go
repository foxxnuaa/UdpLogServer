package commDef

type RecourceItem struct {
	TimeStamp int64
	ServerId int32
	Uid int64
	RecourceType int32
	RecourceId int32
	RecourceCount int32
	StrReason string
}

type ActionLogUnit struct {
	ServerId int32
	Uid int64
	RecourceType int32
	RecourceId int32
	RecourceCount int32
}

type ItemLogUnit struct {
	ServerId int32
	Uid int64
	RecourceType int32
	RecourceId int32
	RecourceCount int32
}

type ResourceDef struct {
	Type                int32
	Id                  int32
	Name                string
	FiveMinuteThreshold int32
	OneHourThreshold    int32
	OneDayThreshold     int32
	Filter string//过滤来源,用|隔离开
	FilterMap map[string]bool
}

type  AlarmDef struct{
	AlarmDay int
	AlarmTimeStamp int64
	ServerId int
	Uid int64
	ResourceType int
	ResourceId int
	ExpireTime int
	ReportCount int
	CurrentCount int
	Threshold int64
	ResourceName string
	Reason string
}