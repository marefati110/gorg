package swagger

import (
	"github.com/go-openapi/spec"
	"github.com/marefati110/gorg/types"
)

type SwaggerConfig struct {
	Swagger *spec.Swagger
}

func (c *SwaggerConfig) Init() {
	c.Swagger = &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Contact: &spec.ContactInfo{},
					License: nil,
				},
				VendorExtensible: spec.VendorExtensible{
					Extensions: spec.Extensions{},
				},
			},
			Paths: &spec.Paths{
				Paths: make(map[string]spec.PathItem),
				VendorExtensible: spec.VendorExtensible{
					Extensions: nil,
				},
			},
			Definitions:         make(map[string]spec.Schema),
			SecurityDefinitions: make(map[string]*spec.SecurityScheme),
		},
		VendorExtensible: spec.VendorExtensible{
			Extensions: nil,
		},
	}
}

func (c *SwaggerConfig) SetVersion() {
	c.Swagger.Swagger = "2.0"
}

func (c *SwaggerConfig) SetInfo(i spec.InfoProps) {
	c.Swagger.Info = &spec.Info{
		InfoProps: i,
	}
}

func (c *SwaggerConfig) SetPath(r types.Route, m types.Module, cfg types.Config) {

	// methods := gorg.MethodResolve(r)
	// url := gorg.UrlResolve(r, m, cfg)
	// fnName := utils.GetFunctionName(r.Handler)

	// for _, method := range methods {

	// }

	// pathItem := spec.PathItem{
	// 	PathItemProps: spec.PathItemProps{
	// 		Get: &spec.Operation{
	// 			OperationProps: spec.OperationProps{
	// 				Description: r.Doc.Description,
	// 				Summary:     r.Doc.Summary,
	// 				ID:          fnName,
	// 			},
	// 		},
	// 	},
	// }

	// c.Swagger.Paths.Paths[url] = pathItem

}
