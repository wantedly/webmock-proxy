package main

func getReqStruct(file File) Request {
	return parseReqStruct(convertJSONToStruct(readFile(file.Path)))
}

func parseReqStruct(httpInt HttpInteractions) Request {
	return httpInt.Connection[0].Request
}
