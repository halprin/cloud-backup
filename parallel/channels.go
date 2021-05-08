package parallel

import "github.com/aws/aws-sdk-go/service/s3"

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
		currentError := <- currentErrorChannel
		errors = append(errors, currentError)
	}

	return errors
}

func ConvertChannelsOfCompletedPartsToSlice(channels []chan *s3.CompletedPart) []*s3.CompletedPart {
	var completedParts []*s3.CompletedPart

	for _, currentErrorChannel := range channels {
		currentCompletedPart := <- currentErrorChannel
		completedParts = append(completedParts, currentCompletedPart)
	}

	return completedParts
}
