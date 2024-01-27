package parallel

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func InvokeErrorReturnFunction(errorFunction func() error) chan error {
	errorChannel := make(chan error, 1)
	go func() {
		errorChannel <- errorFunction()
		close(errorChannel)
	}()

	return errorChannel
}

func ConvertChannelsOfErrorToErrorSlice(errorChannels []chan error) []error {
	var errors []error

	for _, currentErrorChannel := range errorChannels {
		currentError := <-currentErrorChannel
		errors = append(errors, currentError)
	}

	return errors
}

func ConvertChannelsOfCompletedPartsToSlice(channels []chan types.CompletedPart) []types.CompletedPart {
	var completedParts []types.CompletedPart

	for _, currentErrorChannel := range channels {
		currentCompletedPart := <-currentErrorChannel
		completedParts = append(completedParts, currentCompletedPart)
	}

	return completedParts
}
