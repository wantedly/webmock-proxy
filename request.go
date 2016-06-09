package main

func getReqStruct(f file) request {
	return parseReqStruct(convertJSONToStruct(readFile(f.Path)))
}

func parseReqStruct(con connection) request {
	return con.Request
}
