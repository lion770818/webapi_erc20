package config

type Config struct {
	Mode  Web   `yaml:"web"`
	Web   Web   `yaml:"web"`
	Mysql Mysql `yaml:"mysql"`
	Auth  Auth  `yaml:"auth"`
	Redis Redis `yaml:"redis"`
	Node  Node  `yaml:"node"`
	Log   Log   `yaml:"log"`
}
type Web struct {
	Mode string `yaml:"mode"`
	Port string `yaml:"port"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
type Auth struct {
	Active     string `yaml:"active"`
	ExpireTime string `yaml:"expireTime"`
	PrivateKey string `yaml:"privateKey"`
}
type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}
type Node struct {
	Url     []string `yaml:"url"`
	Confirm int64    `yaml:"confirm"`
}
type Log struct {
	Env        string `yaml:"env"`
	Path       string `yaml:"path"`
	Encoding   string `yaml:"encoding"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
}
