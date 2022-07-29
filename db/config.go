package db

type Config struct {
	Driver     string `yaml:"driver"`
	Addr       string `yaml:"addr"`
	Port       string `yaml:"port"`
	DB         string `yaml:"db"`
	UserEnvKey string `yaml:"user_env_key"`
	PassEnvKey string `yaml:"pass_env_key"`
}
