package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/api"
	"github.com/khaivutri/bookmark-service/pkg/jwtutils"
	"github.com/khaivutri/bookmark-service/pkg/logger"
	"github.com/khaivutri/bookmark-service/pkg/validation"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	privateKeyEnv = "JWT_PRIVATE_KEY"
	publicKeyEnv  = "JWT_PUBLIC_KEY"
)

func CreateApp() api.Engine {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	if err := validation.RegisterValidation(); err != nil {
		panic(err)
	}

	logger.SetLogLevel(cfg.LogLevel)

	redisClient := createRedis()

	dbClient := createDB()
	migrate(dbClient)

	privateKey, err := loadKeyFromEnv(privateKeyEnv)
	if err != nil {
		panic(err)
	}
	publicKey, err := loadKeyFromEnv(publicKeyEnv)
	if err != nil {
		panic(err)
	}

	a := createAPIApp(cfg, redisClient, dbClient, privateKey, publicKey)
	return a
}

func createAPIApp(cfg *api.Config, redis *redis.Client, db *gorm.DB, privateKey, publicKey []byte) api.Engine {
	app := gin.New()

	jwtGen, err := jwtutils.NewJWTGeneratorFromPEM(privateKey)
	if err != nil {
		panic(err)
	}
	jwtVal, err := jwtutils.NewJWTValidatorFromPEM(publicKey)
	if err != nil {
		panic(err)
	}
	a := api.NewEngine(api.EngineOpts{
		App:    app,
		Cfg:    cfg,
		Redis:  redis,
		DB:     db,
		JWTGen: jwtGen,
		JWTVal: jwtVal,
	})
	return a
}
