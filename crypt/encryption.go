package crypt

type envelope struct {
	Key     []byte
	Message []byte
}

func Encrypt(contents []byte) ([]byte, error) {

	dataKey, err := getEncryptionKey()
	if err != nil {
		return nil, err
	}

	messageEnvelope := &envelope{
		Key: dataKey.encryptedDataKey,
	}

	return nil, nil
}
