package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shuymn/nijisanji-db-collector/src/application"
	"github.com/shuymn/nijisanji-db-collector/src/infrastructure/external"
	"github.com/shuymn/nijisanji-db-collector/src/infrastructure/persistence"
	"github.com/shuymn/nijisanji-db-collector/src/interfaces/collector"
)

var cdep *collector.Dependency

func init() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg := new(aws.Config)
	level := zerolog.WarnLevel
	if os.Getenv("STAGE") != "production" {
		level = zerolog.DebugLevel
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		if os.Getenv("AWS_ENDPOINT") == "" {
			log.Panic().Msg("AWS_ENDPOINT is empty")
		}
		cfg.Endpoint = aws.String(os.Getenv("AWS_ENDPOINT"))
	}
	zerolog.SetGlobalLevel(level)

	irepo, err := external.NewIchikaraRepository()
	if err != nil {
		log.Panic().Err(err).Send()
	}

	cfg.Region = aws.String(os.Getenv("AWS_REGION"))
	db := dynamo.New(session.Must(session.NewSession()), cfg)
	lrepo := persistence.NewLiveRepository(db)
	csvc := application.NewCollectService(irepo, lrepo)
	cdep = &collector.Dependency{
		CollectService: csvc,
	}
}

func main() {
	lambda.Start(cdep.CollectLiversHandler)
}
