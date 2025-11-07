package httperrormapper

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
