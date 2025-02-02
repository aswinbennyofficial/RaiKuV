package configs

// Parent for the yaml config schema 
type YamlConfig struct{
	LogConfig LogConfig `mapstructure:"log"`
}


// Config for the logs in the config.yaml
type LogConfig struct {
    Level  string  `mapstructure:"level"`
    Output string  `mapstructure:"output"`
    File   fileLog `mapstructure:"file"`
}

type fileLog struct {
    Path       string `mapstructure:"path"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxAge     int    `mapstructure:"max_age"`
    MaxBackups int    `mapstructure:"max_backups"`
}

