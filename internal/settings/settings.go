package settings

type Settings struct {
	Log struct {
		Level string `yaml:"level"`
		View  string `yaml:"view"`
	} `yaml:"log"`
	Database struct {
		Host string `yaml:"host"`
		Name string `yaml:"name"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		Port string `yaml:"port"`
	} `yaml:"database"`
	ServiceSettings struct {
		DebugMode   bool  `yaml:"debug_mode"`
		WeekDays    []int `yaml:"week_days"`
		PackageSize int   `yaml:"package_size"`
	} `yaml:"service_settings"`
	Service struct {
		WorkersCount int `yaml:"workers_count"`
	} `yaml:"service"`
}
