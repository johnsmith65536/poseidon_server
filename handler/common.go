package handler

type Status struct {
	StatusCode    int
	StatusMessage string
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
