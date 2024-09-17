package rpcsync

type Config struct {
	Enabled                   bool    `mapstructure:"Enabled"`
	Port                      int     `mapstructure:"Port"`
	MaxRequestsPerIPAndSecond float64 `mapstructure:"MaxRequestsPerIPAndSecond"`
}
