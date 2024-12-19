package config

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/kataras/iris/v12"
	"github.com/majid-cj/go-chat-server/infrastructure/auth"
	"github.com/majid-cj/go-chat-server/infrastructure/persistence"
	"github.com/majid-cj/go-chat-server/util/fileupload"
	"github.com/olahol/melody"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AppConfig ...
type AppConfig struct {
	sync.RWMutex
	Melody      *melody.Melody
	IPInfo      *ipinfo.Client
	Log         *zap.SugaredLogger
	AppContext  context.Context
	Wg          sync.WaitGroup
	ErrChan     chan error
	App         *iris.Application
	Persistence *persistence.Repository
	Auth        *auth.DBAuth
	Token       *auth.Token
	Upload      *fileupload.UploadFile
	Session     map[string]*melody.Session
}

// NewAppConfig ...
func NewAppConfig() (*AppConfig, error) {
	Persistence, err := persistence.NewRepository()
	if err != nil {
		return nil, err
	}
	Auth := auth.NewDBAuth()
	IpInfo := ipinfo.NewClient(nil, nil, os.Getenv("IP_INFO"))

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC822)

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()

	return &AppConfig{
		Melody:      melody.New(),
		IPInfo:      IpInfo,
		Log:         sugar,
		AppContext:  context.Background(),
		ErrChan:     make(chan error),
		App:         iris.New(),
		Persistence: Persistence,
		Auth:        Auth,
		Token:       auth.NewToken(),
		Upload:      fileupload.NewUploadFile(),
		Session:     make(map[string]*melody.Session),
	}, nil
}

// Set ...
func (config *AppConfig) Set(key string, value *melody.Session) {
	config.Lock()
	defer config.Unlock()
	config.Session[key] = value
}

// Get ...
func (config *AppConfig) Get(key string) *melody.Session {
	config.RLock()
	defer config.RUnlock()
	return config.Session[key]
}

// CloseSession ...
func (config *AppConfig) CloseSession(key string) {
	config.RLock()
	defer config.RUnlock()
	delete(config.Session, key)
}

// SendNotifications ...
func (config *AppConfig) SendNotifications(message interface{}) error {
	return nil
}
