package queue

import (
	"encoding/json"
	"fmt"

	"github.com/betterme-dev/go-server-core/pkg/env"
)

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) PublishMessage(queue string, msg Message) error {
	item, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("error while encoding to json: %w", err)
	}
	err = env.Queue().Publish(queue, item)
	if err != nil {
		return fmt.Errorf("error while publishing message into queue(%s): %w", queue, err)
	}

	return nil
}
