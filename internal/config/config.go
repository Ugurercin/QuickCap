package config

import (
	"fmt"

	"example.com/internal/common"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port" json:"port"`
	} `mapstructure:"server" json:"server"`

	Output struct {
		Directory string                 `mapstructure:"directory" json:"directory"`
		FPS       int                    `mapstructure:"fps" json:"fps"`
		Mode      common.CaptureMode     `mapstructure:"mode" json:"mode"`
		Position  common.OverlayPosition `mapstructure:"position" json:"position"`
		WebcamW   int                    `mapstructure:"webcam_width" json:"webcam_width"`
		WebcamH   int                    `mapstructure:"webcam_height" json:"webcam_height"`
	} `mapstructure:"output" json:"output"`
}

type ConfigUpdate struct {
	Server *struct {
		Port *int `json:"port"`
	} `json:"server"`

	Output *struct {
		Directory *string                 `json:"directory"`
		FPS       *int                    `json:"fps"`
		Mode      *common.CaptureMode     `json:"mode"`     // enum
		Position  *common.OverlayPosition `json:"position"` // enum
		WebcamW   *int                    `json:"webcam_width"`
		WebcamH   *int                    `json:"webcam_height"`
	} `json:"output"`
}

func Load() (*Config, error) {
	v := viper.New()

	// defaults
	v.SetDefault("output.directory", "./captures")
	v.SetDefault("output.fps", 30)
	v.SetDefault("output.mode", string(common.CaptureScreen))
	v.SetDefault("output.position", string(common.PositionBottomLeft))
	v.SetDefault("output.webcam_width", 320)
	v.SetDefault("output.webcam_height", 240)

	// config file settings
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./backend")

	// env support
	v.SetEnvPrefix("SCREENRECORDER")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		fmt.Println("No config file found, will create one with defaults")

		if err := v.WriteConfigAs("config.yaml"); err != nil {
			fmt.Println("Failed to write default config:", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	ValidateEnums(&cfg)
	return &cfg, nil
}

func ValidateEnums(cfg *Config) {

	if cfg.Output.Mode != common.CaptureScreen &&
		cfg.Output.Mode != common.CaptureCamera &&
		cfg.Output.Mode != common.CaptureAsOverlay {
		fmt.Printf("Invalid mode %q in config, falling back to 'screen'\n", cfg.Output.Mode)
		cfg.Output.Mode = common.CaptureScreen
	}

	if cfg.Output.Position != common.PositionTopLeft &&
		cfg.Output.Position != common.PositionTopRight &&
		cfg.Output.Position != common.PositionBottomLeft &&
		cfg.Output.Position != common.PositionBottomRight {
		fmt.Printf("Invalid position %q in config, falling back to 'bottom-right'\n", cfg.Output.Position)
		cfg.Output.Position = common.PositionBottomRight
	}
}

func Save(cfg *Config) error {
	v := viper.New()
	v.Set("server.port", cfg.Server.Port)

	v.Set("output.directory", cfg.Output.Directory)
	v.Set("output.fps", cfg.Output.FPS)
	v.Set("output.mode", string(cfg.Output.Mode))
	v.Set("output.position", string(cfg.Output.Position))
	v.Set("output.webcam_width", cfg.Output.WebcamW)
	v.Set("output.webcam_height", cfg.Output.WebcamH)

	return v.WriteConfigAs("config.yaml")
}
