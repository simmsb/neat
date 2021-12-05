package testbeds

import "github.com/segmentio/ksuid"

func generateID() string {
	id, err := ksuid.NewRandom()
	if err != nil {
		panic("neat id generation failed")
	}
	return id.String()
}
