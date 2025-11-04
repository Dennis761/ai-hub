package httpErrorMapper

type HttpError struct {
	Status int
	Body   map[string]string
}

func newHttpError(status int, err error, key string) HttpError {
	return HttpError{
		Status: status,
		Body:   map[string]string{"error": err.Error(), "key": key},
	}
}

func extractKey(err error) string {
	if err == nil {
		return ""
	}

	// 1) Code()
	if c, ok := err.(interface{ Code() string }); ok {
		return c.Code()
	}

	// 2) KeyName()
	if k, ok := err.(interface{ KeyName() string }); ok {
		return k.KeyName()
	}

	// 3) Name()
	if n, ok := err.(interface{ Name() string }); ok {
		return n.Name()
	}

	return ""
}
