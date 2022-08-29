package utils

type Config struct {
	AppName    string
	Env        string `env:"ENV,default=local"`
	DBURL      string `env:"DB_URL,default=./db.sqlite"`
	Port       string `env:"PORT,default=3000"`
	MsgService string `env:"MESSAGING_SERVICE,required"`
	MsgKey     string `env:"MESSAGING_KEY,required"`
}
