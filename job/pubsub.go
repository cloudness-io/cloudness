package job

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/cloudness-io/cloudness/pubsub"
)

const (
	PubSubTopicCancelJob   = "cloudness:job:cancel_job"
	PubSubTopicStateChange = "cloudness:job:state_change"
)

func encodeStateChange(job *Job) ([]byte, error) {
	stateChange := &StateChange{
		UID:      job.UID,
		Type:     job.Type,
		State:    job.State,
		Progress: job.RunProgress,
		Result:   job.Result,
		Failure:  job.LastFailureError,
	}

	buffer := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buffer).Encode(stateChange); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DecodeStateChange(payload []byte) (*StateChange, error) {
	stateChange := &StateChange{}
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(stateChange); err != nil {
		return nil, err
	}

	return stateChange, nil
}

func publishStateChange(ctx context.Context, publisher pubsub.Publisher, job *Job) error {
	payload, err := encodeStateChange(job)
	if err != nil {
		return fmt.Errorf("failed to gob encode StateChange: %w", err)
	}

	err = publisher.Publish(ctx, PubSubTopicStateChange, payload)
	if err != nil {
		return fmt.Errorf("failed to publish StateChange: %w", err)
	}

	return nil
}
