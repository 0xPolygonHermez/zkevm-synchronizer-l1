package sqlstorage

import "fmt"

type Config struct {
	// Name of the database
	DriverName string `mapstructure:"DriverName"`
	DataSource string `mapstructure:"DataSource"`
}

func (c *Config) String() string {
	return fmt.Sprintf("DriverName=%s DataSource=%s", c.DriverName, c.DataSource)
}

func (c *Config) SanityCheck() error {
	if c.DriverName == "" {
		return fmt.Errorf("DriverName is required")
	}
	if c.DataSource == "" {
		return fmt.Errorf("DataSource is required")
	}
	if c.DriverName != SqliteDriverName {
		return fmt.Errorf("DriverName not supported: %s, only supports: %s", c.DriverName, SqliteDriverName)
	}
	return nil
}
