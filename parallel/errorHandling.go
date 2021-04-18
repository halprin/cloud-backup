package parallel

func InvokeErrorReturnFunction(errorFunction func() error) chan error {
	errorChannel := make(chan error, 1)
	go func() {
		errorChannel <- errorFunction()
		close(errorChannel)
	}()

	return errorChannel
}

func ConvertChannelsOfErrorToErrorSlice(errorChannels []chan error) []error {
	var errors  []error

	for _, currentErrorChannel := range errorChannels {
		currentError := <- currentErrorChannel
		errors = append(errors, currentError)
	}

	return errors
}
