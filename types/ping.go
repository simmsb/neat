package types

type PingRequest struct {
	Sender   string  `mapstructure:"sender" json:"sender"`
	Target   string  `mapstructure:"target" json:"target"`
	Count    uint    `mapstructure:"count" json:"count"`
	Interval float64 `mapstructure:"interval" json:"interval"`
}

type PingResponse struct {
	Sent       uint    `mapstructure:"sent" json:"sent"`
	Received   uint    `mapstructure:"received" json:"received"`
	AverageRTT float64 `mapstructure:"avg_rtt" json:"avg_rtt"`
	StdDev     float64 `mapstructure:"std_dev" json:"std_dev"`
}
