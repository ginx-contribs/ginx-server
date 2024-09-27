package conf

import (
	"github.com/246859/duration"
	"github.com/ginx-contribs/logx"
	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"
	"log/slog"
	"os"
	"reflect"
)

// ReadFrom read app configuration from specified file
func ReadFrom(filename string) (App, error) {
	var app App
	configFile, err := os.Open(filename)
	if err != nil {
		return App{}, err
	}
	defer configFile.Close()

	err = toml.NewDecoder(configFile).Decode(&app)
	if err != nil {
		return App{}, err
	}

	return app, nil
}

// WriteTo save the app configuration to specified file
func WriteTo(filename string, app App) error {
	configFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer configFile.Close()
	return toml.NewEncoder(configFile).Encode(app)
}

// DefaultConfig is the default configuration for application
var DefaultConfig = App{
	Server: Server{
		Address:      "127.0.0.1:8080",
		BasePath:     "/api",
		ReadTimeout:  duration.Minute,
		WriteTimeout: duration.Minute,
		IdleTimeout:  5 * duration.Minute,
		MultipartMax: 50 << 20,
		Pprof:        false,
	},
	Log: Log{
		Filename: "/etc/ginx-internal/internal.log",
		Prompt:   "[ginx-internal]",
		Level:    slog.LevelInfo,
		Format:   logx.TextFormat,
		Source:   false,
		Color:    false,
	},
	DB: DB{
		Driver:             "unset",
		Address:            "127.0.0.1:3306",
		User:               "user",
		Password:           "password",
		Database:           "database",
		Params:             "",
		MaxIdleConnections: 10,
		MaxOpenConnections: 100,
		MaxLifeTime:        duration.Hour,
		MaxIdleTime:        10 * duration.Minute,
	},
	Redis: Redis{
		Address:      "127.0.0.1:6379",
		Password:     "password",
		WriteTimeout: duration.Minute,
		ReadTimeout:  duration.Minute,
	},
	Email: Email{
		Host:     "",
		Port:     0,
		Username: "",
		Password: "",
		MQ: EmailMq{
			Topic:     "email",
			BatchSize: 20,
			Group:     "email-group",
			Consumers: []string{"consumerA"},
		},
		Code: VerifyCode{
			TTL:      5 * duration.Minute,
			RetryTTL: duration.Minute,
		},
	},
	Jwt: Jwt{
		Issuer: "lobby",
		Access: AccessToken{
			Expire: 4 * duration.Hour,
			Delay:  10 * duration.Minute,
			Key:    "01J6EA2G4FSF9ABC218VFJ2B3C",
		},
		Refresh: RefreshToken{
			Expire: 144 * duration.Hour,
			Key:    "01J6EA3FKDDHTT9Q8Z5YKWHVCE",
		},
	},
}

// Revise check the given configuration, if field value is zero then it will be overwritten by same filed value of DefaultConfig
func Revise(cfg App) (App, error) {
	src, dst := make(map[string]any), make(map[string]any)

	err := mapstructure.Decode(DefaultConfig, &src)
	if err != nil {
		return App{}, err
	}

	err = mapstructure.Decode(cfg, &dst)
	if err != nil {
		return App{}, err
	}
	reviseMap(src, dst)

	var destConf App
	err = mapstructure.Decode(dst, &destConf)
	if err != nil {
		return App{}, err
	}
	return destConf, nil
}

func reviseMap(src, dst map[string]any) {
	for srcKey, srcVal := range src {
		dstVal := dst[srcKey]
		if reflect.TypeOf(dstVal).Kind() == reflect.Map && reflect.TypeOf(srcVal).Kind() == reflect.Map {
			reviseMap(srcVal.(map[string]any), dstVal.(map[string]any))
		} else {
			if reflect.ValueOf(dstVal).IsZero() {
				dst[srcKey] = srcVal
			}
		}
	}
}
