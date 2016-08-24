package webmock

import (
	"path/filepath"
	"time"
)

func cacheImport(s *Server) error {
	files, err := readFilePaths(s.config.importCache)
	if err != nil {
		return err
	}
	for _, file := range files {
		dst := filepath.Join(s.config.importCache, file)
		b, err := readFile(dst)
		if err != nil {
			return err
		}
		var conn Connection
		err = jsonToStruct(b, &conn)
		if err != nil {
			return err
		}
		var conns []Connection
		conns = append(conns, conn)
		endpoint := &Endpoint{
			URL:         file,
			Connections: conns,
			Update:      time.Now(),
		}
		ce := readEndpoint(file, s.db)
		if len(ce.Connections) != 0 {
			for _, v := range ce.Connections {
				deleteConnection(&v, s.db)
				if v.Request.Method == conn.Request.Method {
					continue
				}
				conns = append(conns, v)
			}
			endpoint.Connections = conns
			updateEndpoint(ce, endpoint, s.db)
		}
		if err := insertEndpoint(endpoint, s.db); err != nil {
			return err
		}
	}
	return nil
}
