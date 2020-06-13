package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shuymn/nijisanji-db-collector/src/interfaces/collector"
)

var cdep *collector.Dependency

func init() {
	cdep = &collector.Dependency{}
}

func main() {
	lambda.Start(cdep.CollectLiversHandler)
}
