package options

import (
	"time"

	"github.com/beego/beego"
	"github.com/mitchellh/mapstructure"
)

// Options represents the options of service.
type Options struct {
	FEURL string `mapstructure:"fe_url"`

	DBUser string `mapstructure:"db_user"`
	DBPass string `mapstructure:"db_pass"`
	DBHost string `mapstructure:"db_host"`
	DBPort string `mapstructure:"db_port"`
	DBName string `mapstructure:"db_name"`

	AccessTokenExpiresIn time.Duration `mapstructure:"auth_access_token_expires_in"`
	AccessTokenSecret    string        `mapstructure:"auth_access_token_secret"`

	GAuthClientID     string `mapstructure:"gauth_client_id"`
	GAuthClientSecret string `mapstructure:"gauth_client_secret"`
	GAuthProfileURL   string `mapstructure:"gauth_profile_url"`
}

// Load loads Options from Viper and returns them.
func Load() (Options, error) {
	var opts Options
	cfg, err := beego.AppConfig.GetSection("default")
	if err != nil {
		return opts, err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
		Result: &opts,
	})
	if err != nil {
		return Options{}, err
	}
	if err := decoder.Decode(cfg); err != nil {
		return Options{}, err
	}

	return opts, nil
}
