package boot

import (
	"go-bunrouter-example/module/article"
	"go-bunrouter-example/module/health"
	"os"

	"go-bunrouter-example/infrastructure/config"
	"go-bunrouter-example/infrastructure/database"
	"go-bunrouter-example/infrastructure/limiter"
	logger "go-bunrouter-example/infrastructure/log"
	"go-bunrouter-example/infrastructure/redis"
	"go-bunrouter-example/utils"

	log "github.com/sirupsen/logrus"
)

type HandlerSetup struct {
	Limiter     *limiter.RateLimiter
	HealthHttp  health.InterfaceHttp
	ArticleHttp article.InterfaceHttp
}

func MakeHandler() HandlerSetup {
	//initiate config
	config.Initialize()

	//initiate logger
	logger.Init(config.Conf.LogFormat, config.Conf.LogLevel)

	var err error

	//initiate a redis client
	redisClient, err := redis.NewRedisClient(&config.Conf)
	if err != nil {
		log.Fatalf("failed initiate redis: %v", err)
		os.Exit(1)
	}

	//initiate a redis library interface
	redisLibInterface, err := redis.NewRedisLibInterface(redisClient)
	if err != nil {
		log.Fatalf("failed initiate redis library: %v", err)
		os.Exit(1)
	}

	//setup infrastructure postgres
	db, err := database.NewPostgresDatabaseClient(&config.Conf)
	if err != nil {
		log.Fatalf("failed initiate database postgres: %v", err)
		os.Exit(1)
	}

	//add limiter
	interval := utils.StringUnitToDuration(config.Conf.Interval)
	middlewareWithLimiter := limiter.NewRateLimiter(int(config.Conf.Rate), interval)

	//health module
	healthRepository := health.NewRepository(db.DbConn)
	healthService := health.NewService(healthRepository, redisClient)
	healthModule := health.NewHttp(healthService)

	//article module
	articleRepository := article.NewRepository(db.DbConn)
	articleService := article.NewService(articleRepository, redisLibInterface)
	articleModule := article.NewHttp(articleService)

	return HandlerSetup{
		Limiter:     middlewareWithLimiter,
		HealthHttp:  healthModule,
		ArticleHttp: articleModule,
	}
}
