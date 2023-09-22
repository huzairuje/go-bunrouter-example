package health

import (
	"context"
	"go-bunrouter-example/infrastructure/config"

	logger "go-bunrouter-example/infrastructure/log"
	"go-bunrouter-example/module/primitive"

	"github.com/go-redis/redis"
)

type InterfaceService interface {
	CheckUpTime(ctx context.Context) (resp primitive.HealthResp, err error)
}

type Service struct {
	repository  RepositoryInterface
	redisClient *redis.Client
}

func NewService(repository RepositoryInterface, redisClient *redis.Client) InterfaceService {
	return &Service{
		repository:  repository,
		redisClient: redisClient,
	}
}

func (u *Service) CheckUpTime(ctx context.Context) (primitive.HealthResp, error) {
	ctxName := "CheckUpTime"
	errCheckDb := u.repository.CheckUpTimeDB()
	if errCheckDb != nil {
		logger.Error(ctx, ctxName, "got error when %s : %v", ctxName, errCheckDb)
		return primitive.HealthResp{}, errCheckDb
	}

	if config.Conf.Redis.EnableRedis && u.redisClient != nil {
		errCheckRedis := u.redisClient.Ping().Err()
		if errCheckRedis != nil {
			logger.Error(ctx, ctxName, "got error when %s : %v", ctxName, errCheckRedis)
			return primitive.HealthResp{}, errCheckRedis
		}
	}

	return primitive.HealthResp{
		Db:    "healthy",
		Redis: "healthy",
	}, nil
}
