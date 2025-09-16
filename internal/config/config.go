package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port" json:"port"`
	} `mapstructure:"server" json:"server"`

	Output struct {
		Directory                 string `mapstructure:"directory" json:"directory"`
		FPS                       int    `mapstructure:"fps" json:"fps"`
		StartVideoRecordingHotkey string `mapstructure:"start_video_recording_hotkey" json:"start_video_recording_hotkey"`
		CaptureScreenShotHotkey   string `mapstructure:"capture_screenshot_hotkey" json:"capture_screenshot_hotkey"`
	} `mapstructure:"output" json:"output"`
}

type ConfigUpdate struct {
	Server struct {
		Port *int `json:"port"`
	} `json:"server"`

	Output struct {
		Directory                 *string `json:"directory"`
		FPS                       *int    `json:"fps"`
		StartVideoRecordingHotkey *string `json:"start_video_recording_hotkey"`
		CaptureScreenShotHotkey   *string `json:"capture_screenshot_hotkey"`
	} `json:"output"`
}

func Load() (*Config, error) {
	v := viper.New()

	// defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("output.directory", "./captures")
	v.SetDefault("output.fps", 30)
	v.SetDefault("output.start_video_recording_hotkey", "2")
	v.SetDefault("output.capture_screenshot_hotkey", "1")

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

	return &cfg, nil
}

func Save(cfg *Config) error {
	v := viper.New()
	v.Set("server.port", cfg.Server.Port)
	v.Set("output.directory", cfg.Output.Directory)
	v.Set("output.fps", cfg.Output.FPS)
	v.Set("output.start_video_recording_hotkey", cfg.Output.StartVideoRecordingHotkey)
	v.Set("output.capture_screenshot_hotkey", cfg.Output.CaptureScreenShotHotkey)
	return v.WriteConfigAs("config.yaml")
}
