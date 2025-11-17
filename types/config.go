package types

import (
	"fmt"
	"time"

	"github.com/cloudness-io/cloudness/lock"
	"github.com/cloudness-io/cloudness/pubsub"
)

// Config stores the system configuration.
type Config struct {
	// InstanceID specifis the ID of the cloudness instance.
	// NOTE: If the value is not provided the hostname of the machine is used.
	InstanceID  string `envconfig:"CLOUDNESS_INSTANCE_ID"`
	Environment string `envconfig:"CLOUDNESS_ENVIRONMENT"`

	Debug bool `envconfig:"CLOUDNESS_DEBUG"`
	Trace bool `envconfig:"CLOUDNESS_TRACE"`

	// GracefulShutdownTime defines the max time we wait when shutting down a server.
	// 5min should be enough for most git clones to complete.
	GracefulShutdownTime time.Duration `envconfig:"CLOUDNESS_GRACEFUL_SHUTDOWN_TIME" default:"300s"`

	Profiler struct {
		Type        string `envconfig:"CLOUDNESS_PROFILER_TYPE"`
		ServiceName string `envconfig:"CLOUDNESS_PROFILER_SERVICE_NAME" default:"cloudness"`
	}

	Instance struct {
		AllowNewTenantCreation bool `envconfig:"CLOUDNESS_INSTANCE_ALLOW_NEW_TENANT" default:"true"`
	}

	TenantConfig     TenantConfig
	KubeServerConfig KubeServerConfig

	// Token defines token configuration parameters.
	Token struct {
		CookieName string        `envconfig:"CLOUDNESS_TOKEN_COOKIE_NAME" default:"token"`
		Expire     time.Duration `envconfig:"CLOUDNESS_TOKEN_EXPIRE" default:"720h"`
	}

	// Database defines the database configuration parameters.
	Database struct {
		Driver     string `envconfig:"CLOUDNESS_DATABASE_DRIVER"     default:"sqlite3"`
		Datasource string `envconfig:"CLOUDNESS_DATABASE_DATASOURCE" default:"database.sqlite3"`
		Host       string `envconfig:"CLOUDNESS_DATABASE_HOST"`
		Port       string `envconfig:"CLOUDNESS_DATABASE_PORT"`
		Name       string `envconfig:"CLOUDNESS_DATABASE_NAME"`
		User       string `envconfig:"CLOUDNESS_DATABASE_USER"`
		Password   string `envconfig:"CLOUDNESS_DATABASE_PASSWORD"`
		SSLMode    string `envconfig:"CLOUDNESS_DATABASE_SSL_MODE"`
	}

	PubSub struct {
		Provider         pubsub.Provider `envconfig:"CLOUDNESS_PUBSUB_PROVIDER"          default:"inmemory"`
		AppNamespace     string          `envconfig:"CLOUDNESS_PUBSUB_APP_NAMESPACE"     default:"cloudness"`
		DefaultNamespace string          `envconfig:"CLOUDNESS_PUBSUB_DEFAULT_NAMESPACE" default:"default"`
		HealthInterval   time.Duration   `envconfig:"CLOUDNESS_PUBSUB_HEALTH_INTERVAL"   default:"3s"`
		SendTimeout      time.Duration   `envconfig:"CLOUDNESS_PUBSUB_SEND_TIMEOUT"      default:"60s"`
		ChannelSize      int             `envconfig:"CLOUDNESS_PUBSUB_CHANNEL_SIZE"      default:"100"`
	}

	// Server defines the server configuration parameters.
	Server struct {
		// HTTP defines the http configuration parameters
		HTTP struct {
			Port  int    `envconfig:"CLOUDNESS_HTTP_PORT" default:"8000"`
			Proto string `envconfig:"CLOUDNESS_HTTP_PROTO" default:"http"`
		}

		// Acme defines Acme configuration parameters.
		Acme struct {
			Enabled bool   `envconfig:"CLOUDNESS_ACME_ENABLED"`
			Endpont string `envconfig:"CLOUDNESS_ACME_ENDPOINT"`
			Email   bool   `envconfig:"CLOUDNESS_ACME_EMAIL"`
			Host    string `envconfig:"CLOUDNESS_ACME_HOST"`
		}
	}

	// Cors defines http cors parameters
	Cors struct {
		AllowedOrigins   []string `envconfig:"CLOUDNESS_CORS_ALLOWED_ORIGINS"   default:"*"`
		AllowedMethods   []string `envconfig:"CLOUDNESS_CORS_ALLOWED_METHODS"   default:"GET,POST,PATCH,PUT,DELETE,OPTIONS"`
		AllowedHeaders   []string `envconfig:"CLOUDNESS_CORS_ALLOWED_HEADERS"   default:"Origin,Accept,Accept-Language,Authorization,Content-Type,Content-Language,X-Requested-With,X-Request-Id"` //nolint:lll // struct tags can't be multiline
		ExposedHeaders   []string `envconfig:"CLOUDNESS_CORS_EXPOSED_HEADERS"   default:"Link"`
		AllowCredentials bool     `envconfig:"CLOUDNESS_CORS_ALLOW_CREDENTIALS" default:"true"`
		MaxAge           int      `envconfig:"CLOUDNESS_CORS_MAX_AGE"           default:"300"`
	}

	Redis struct {
		Endpoint           string `envconfig:"CLOUDNESS_REDIS_ENDPOINT"              default:"localhost:6379"`
		MaxRetries         int    `envconfig:"CLOUDNESS_REDIS_MAX_RETRIES"           default:"3"`
		MinIdleConnections int    `envconfig:"CLOUDNESS_REDIS_MIN_IDLE_CONNECTIONS"  default:"0"`
		Username           string `envconfig:"CLOUDNESS_REDIS_USERNAME"`
		Password           string `envconfig:"CLOUDNESS_REDIS_PASSWORD"`
	}

	// Secure defines http security parameters.
	Secure struct {
		AllowedHosts          []string          `envconfig:"CLOUDNESS_HTTP_ALLOWED_HOSTS"`
		HostsProxyHeaders     []string          `envconfig:"CLOUDNESS_HTTP_PROXY_HEADERS"`
		SSLRedirect           bool              `envconfig:"CLOUDNESS_HTTP_SSL_REDIRECT"`
		SSLTemporaryRedirect  bool              `envconfig:"CLOUDNESS_HTTP_SSL_TEMPORARY_REDIRECT"`
		SSLHost               string            `envconfig:"CLOUDNESS_HTTP_SSL_HOST"`
		SSLProxyHeaders       map[string]string `envconfig:"CLOUDNESS_HTTP_SSL_PROXY_HEADERS"`
		STSSeconds            int64             `envconfig:"CLOUDNESS_HTTP_STS_SECONDS"`
		STSIncludeSubdomains  bool              `envconfig:"CLOUDNESS_HTTP_STS_INCLUDE_SUBDOMAINS"`
		STSPreload            bool              `envconfig:"CLOUDNESS_HTTP_STS_PRELOAD"`
		ForceSTSHeader        bool              `envconfig:"CLOUDNESS_HTTP_STS_FORCE_HEADER"`
		BrowserXSSFilter      bool              `envconfig:"CLOUDNESS_HTTP_BROWSER_XSS_FILTER"    default:"true"`
		FrameDeny             bool              `envconfig:"CLOUDNESS_HTTP_FRAME_DENY"            default:"true"`
		ContentTypeNosniff    bool              `envconfig:"CLOUDNESS_HTTP_CONTENT_TYPE_NO_SNIFF"`
		ContentSecurityPolicy string            `envconfig:"CLOUDNESS_HTTP_CONTENT_SECURITY_POLICY"`
		ReferrerPolicy        string            `envconfig:"CLOUDNESS_HTTP_REFERRER_POLICY"`
	}

	//
	// Source code management.
	//

	// Github provides the github client configuration.
	Github struct {
		ClientId     string   `envconfig:"CLOUDNESS_SCM_GITHUB_CLIENT_ID"`
		ClientSecret string   `envconfig:"CLOUDNESS_SCM_GITHUB_CLIENT_SECRET"`
		Scope        []string `envconfig:"CLOUDNESS_SCM_GITHUB_SCOPE" default:"repo,repo:status,user:email,read:org"`
		Debug        bool     `envconfig:"CLOUDNESS_SCM_GITHUB_DEBUG"`
	}

	Lock struct {
		// Provider is a name of distributed lock service like redis, memory, file etc...
		Provider      lock.Provider `envconfig:"CLOUDNESS_LOCK_PROVIDER"          default:"inmemory"`
		Expiry        time.Duration `envconfig:"CLOUDNESS_LOCK_EXPIRE"            default:"8s"`
		Tries         int           `envconfig:"CLOUDNESS_LOCK_TRIES"             default:"8"`
		RetryDelay    time.Duration `envconfig:"CLOUDNESS_LOCK_RETRY_DELAY"       default:"250ms"`
		DriftFactor   float64       `envconfig:"CLOUDNESS_LOCK_DRIFT_FACTOR"      default:"0.01"`
		TimeoutFactor float64       `envconfig:"CLOUDNESS_LOCK_TIMEOUT_FACTOR"    default:"0.25"`
		// AppNamespace is just service app prefix to avoid conflicts on key definition
		AppNamespace string `envconfig:"CLOUDNESS_LOCK_APP_NAMESPACE"     default:"cloudness"`
		// DefaultNamespace is when mutex doesn't specify custom namespace for their keys
		DefaultNamespace string `envconfig:"CLOUDNESS_LOCK_DEFAULT_NAMESPACE" default:"default"`
	}
}

// Process performs post-processing on the configuration after it has been loaded.
func (c *Config) Process() {
	// If the database driver is postgres and a full datasource is not already provided,
	// construct it from the individual parts.
	if c.Database.Driver == "postgres" && c.Database.Datasource == "" {
		c.Database.Datasource = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			c.Database.User,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Name,
			c.Database.SSLMode,
		)
	}
}

type TenantConfig struct {
	DefaultTenantName                 string  `envconfig:"CLOUDNESS_TENANT_DEFAULT_NAME"           default:"root"`
	DefaultAllowAdminToModify         bool    `envconfig:"CLOUDNESS_TENANT_ALLOW_ADMIN_TO_MODIFY"  default:"false"`
	DefaultMaxProjectsPerTenant       int64   `envconfig:"CLOUDNESS_TENANT_MAX_PROJECTS"           default:"10"`
	DefaultMaxApplicationsPerTenant   int64   `envconfig:"CLOUDNESS_TENANT_MAX_APPLICATIONS"       default:"10"`
	DefaultMaxInstancesPerApplication int64   `envconfig:"CLOUDNESS_TENANT_MAX_INSTANCES"          default:"2"`
	DefaultMaxCPUPerApplication       int64   `envconfig:"CLOUDNESS_TENANT_MAX_CPU"                default:"2"`
	DefaultMaxMemoryPerApplication    float64 `envconfig:"CLOUDNESS_TENANT_MAX_MEMORY"             default:"2"`
	DefaultMaxVolumeCount             int64   `envconfig:"CLOUDNESS_TENANT_MAX_VOLUME_COUNT"       default:"10"`
	DefaultMinVolumeSize              int64   `envconfig:"CLOUDNESS_TENANT_MIN_VOLUME_SIZE"        default:"1"`
	DefaultMaxVolumeSize              int64   `envconfig:"CLOUDNESS_TENANT_MAX_VOLUME_SIZE"        default:"10"`
}

type KubeServerConfig struct {
	DefaultVolumeSupportsOnlineExpansion bool `envconfig:"CLOUDNESS_KUBE_UNMOUNT_BEFORE_RESIZE"  default:"true"`
}
