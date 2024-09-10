package sqlstorage

type Config struct {
	// Name of the database
	DriverName string `mapstructure:"DriverName"`
	DataSource string `mapstructure:"DataSource"`
}
