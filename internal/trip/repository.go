package trip

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type TripRepository struct {
	Client *dynamodb.Client
}
