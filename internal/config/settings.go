package config

type Settings struct {
	Db struct {
		Dsn  string `mapstructure:"dsn"`
		Type string `mapstructure:"type"`
	} `mapstructure:"db"`
	Redis struct {
		Dsn string `mapstructure:"dsn"`
	} `mapstructure:"redis"`
	Limiter struct {
		Size     int    `mapstructure:"size"`
		Interval int    `mapstructure:"interval"`
		Dsn      string `mapstructure:"dsn"`
		Name     string `mapstructure:"name"`
		Mode     int    `mapstructure:"mode"`
		Rate     int    `mapstructure:"rate"`
	} `mapstructure:"limiter"`
	Application struct {
		Name         string `mapstructure:"name"`
		ReadTimeout  int    `mapstructure:"read_timeout"`
		WriteTimeout int    `mapstructure:"write_timeout"`
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"port"`
		Mode         string `mapstructure:"mode"`
	} `mapstructure:"application"`
	Middleware []string `mapstructure:"middleware"`
	Log        struct {
		File struct {
			Enable     bool   `mapstructure:"enable"`
			Filename   string `mapstructure:"filename"`
			MaxSize    int    `mapstructure:"max_size"`
			MaxAge     int    `mapstructure:"max_age"`
			MaxBackups int    `mapstructure:"max_backups"`
			LocalTime  bool   `mapstructure:"local_time"`
			Compress   bool   `mapstructure:"compress"`
		} `mapstructure:"file"`
		Kafka struct {
			Enable       bool   `mapstructure:"enable"`
			Addr         string `mapstructure:"addr"`
			BatchTimeout int    `mapstructure:"batch_timeout"`
			WriteTimeout int    `mapstructure:"write_timeout"`
			Async        bool   `mapstructure:"async"`
			BatchSize    int    `mapstructure:"batch_size"`
			BatchBytes   int    `mapstructure:"batch_bytes"`
			RequiredAcks int    `mapstructure:"required_acks"`
			Topic        string `mapstructure:"topic"`
		} `mapstructure:"kafka"`
		Console struct {
			Enable bool `mapstructure:"enable"`
		} `mapstructure:"console"`
	} `mapstructure:"log"`
}
