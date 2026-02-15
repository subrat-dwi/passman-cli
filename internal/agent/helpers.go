package agent

// ok creates a Response with OK set to true, indicating a successful operation without any additional data.
func ok() Response {
	return Response{OK: true}
}

// errResp creates a Response with OK set to false and includes the provided error message.
func errResp(msg string) Response {
	return Response{OK: false, Error: msg}
}

// data creates a Response with OK set to true and includes the provided data in the Data field.
func data(v any) Response {
	return Response{OK: true, Data: v}
}
