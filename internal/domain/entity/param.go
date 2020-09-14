package entity

type FanParam struct {
	TimeOut  int   `json:"timeout"`
	PoolSize int   `json:"poolsize"`
	URLs     []URL `json:"urls"`
}

type FanParamRepo interface {
	Load() (*FanParam, error)
}
