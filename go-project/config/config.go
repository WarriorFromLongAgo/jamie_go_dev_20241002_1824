package config

type DbConfig struct {
	DbHost     string `yaml:"db_host"`
	DbPort     int    `yaml:"db_port"`
	DbName     string `yaml:"db_name"`
	DbUser     string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
}
