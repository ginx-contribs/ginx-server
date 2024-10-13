package conf

import (
	"github.com/246859/duration"
	"log/slog"
)

// App is configuration for the whole application
type App struct {
	Server Server   `toml:"server" comment:"http server configuration"`
	Log    Log      `toml:"log" comment:"server log configuration"`
	DB     DB       `toml:"db" comment:"database connection configuration"`
	Redis  Redis    `toml:"redis" comment:"redis connection configuration"`
	Email  Email    `toml:"email" comment:"email smtp client configuration"`
	Jwt    Jwt      `toml:"jwt" comment:"jwt secret configuration"`
	Meta   MetaInfo `toml:"-"`
}

// MetaInfo for program
type MetaInfo struct {
	AppName   string
	Author    string
	Version   string
	BuildTime string
}

// Server is configuration for the http server
type Server struct {
	Address      string            `toml:"address" comment:"server bind address"`
	BasePath     string            `toml:"basepath" comment:"base path for api"`
	ReadTimeout  duration.Duration `toml:"readTimeout" comment:"the maximum duration for reading the entire request"`
	WriteTimeout duration.Duration `toml:"writeTimeout" comment:"the maximum duration before timing out writes of the response"`
	IdleTimeout  duration.Duration `toml:"idleTimeout" comment:"the maximum amount of time to wait for the next request when keep-alives are enabled"`
	MultipartMax int64             `toml:"multipartMax" comment:"value of 'maxMemory' param that is given to http.Request's ParseMultipartForm"`
	Pprof        bool              `toml:"pprof" comment:"enabled pprof program profiling"`
	Swagger      bool              `toml:"swagger" comment:"enable swagger documentation"`
	TLS          TLS               `toml:"tls" comment:"tls certificate"`
}

type TLS struct {
	Cert string `toml:"cert" comment:"tls certificate"`
	Key  string `toml:"key" comment:"tls key"`
}

// Log is configuration for logging
type Log struct {
	Filename string     `toml:"filename" comment:"log output file"`
	Prompt   string     `toml:"-"`
	Level    slog.Level `toml:"level" comment:"support levels: DEBUG | INFO | WARN | ERROR"`
	Format   string     `toml:"format" comment:"TEXT or JSON"`
	Source   bool       `toml:"source" comment:"whether to show source file in logs"`
	Color    bool       `toml:"color" comment:"enable color log"`
}

// DB is configuration for database
type DB struct {
	Driver             string            `toml:"driver" comment:"sqlite | mysql | postgresql"`
	Address            string            `toml:"address" comment:"db internal host"`
	User               string            `toml:"user" comment:"db username"`
	Password           string            `toml:"password" comment:"db password"`
	Database           string            `toml:"database" comment:"database name"`
	Params             string            `toml:"param" comment:"connection params"`
	MaxIdleConnections int               `toml:"maxIdleConnections" comment:"max idle connections limit"`
	MaxOpenConnections int               `toml:"maxOpenConnections" comment:"max opening connections limit"`
	MaxLifeTime        duration.Duration `toml:"maxLifeTime" comment:"max connection lifetime"`
	MaxIdleTime        duration.Duration `toml:"maxIdleTime" comment:"max connection idle time"`
}

// Redis is configuration for redis internal
type Redis struct {
	Address      string            `toml:"address" comment:"host address"`
	Password     string            `toml:"password" comment:"redis auth"`
	WriteTimeout duration.Duration `toml:"writeTimeout" comment:"Timeout for socket writes."`
	ReadTimeout  duration.Duration `toml:"readTimeout" comment:"Timeout for socket reads."`
}

// Jwt is configuration for jwt signing
type Jwt struct {
	Issuer  string       `toml:"issuer" comment:"jwt issuer"`
	Access  AccessToken  `toml:"access" comment:"access token configuration"`
	Refresh RefreshToken `toml:"refresh" comment:"refresh token configuration"`
}

type AccessToken struct {
	Expire duration.Duration `toml:"expire" comment:"duration to expire access token"`
	Delay  duration.Duration `toml:"delay" comment:"delay duration after expiration"`
	Key    string            `toml:"key" comment:"access token signing key"`
}

type RefreshToken struct {
	Expire duration.Duration `toml:"expire" comment:"duration to expire refresh token"`
	Key    string            `toml:"key" comment:"refresh token signing key"`
}

type RateLimit struct {
	Public struct {
		Limit  int               `toml:"limit"`
		Window duration.Duration `toml:"window"`
	} `toml:"public"`
}

type Email struct {
	Host     string     `toml:"host" comment:"smtp internal host"`
	SSL      bool       `toml:"ssl" comment:"use ssl port"`
	Port     int        `toml:"port" comment:"smtp internal port"`
	Username string     `toml:"username" comment:"smtp user name"`
	Password string     `toml:"password" comment:"password to authenticate"`
	Template string     `toml:"template" comment:"custom email template dir"`
	MQ       EmailMq    `toml:"-"`
	Code     VerifyCode `toml:"code" comment:"email verification code configuration"`
}

type EmailMq struct {
	Topic     string   `toml:"topic" comment:"email mq topic"`
	BatchSize int64    `toml:"batchSize" comment:"max batch size of per reading"`
	Group     string   `toml:"group" comment:"consumer group"`
	Consumers []string `toml:"consumers" comment:"how many consumer in groups, must >=1."`
}

type VerifyCode struct {
	TTL      duration.Duration `toml:"ttl" comment:"lifetime for verification code"`
	RetryTTL duration.Duration `toml:"retry" comment:"max wait time before asking for another new verification code"`
}
