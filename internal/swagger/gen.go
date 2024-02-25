package swagger

import (
	"encoding/json"
	"os"
)

func WriteJSON(c SwaggerConfig) error {

	b, e := json.MarshalIndent(c.Swagger, "", "	")

	if e != nil {
		return e
	}

	f, e := os.Create("swagger.json")
	if e != nil {
		return e
	}
	defer f.Close()

	_, e = f.Write(b)
	if e != nil {
		return e
	}

	return nil
}
