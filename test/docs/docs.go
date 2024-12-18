// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/products": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetProductsResponse"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateProductRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/UpdateProductResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "description": "request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/CreateProductsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/CreateProductsResponse"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/DeleteProductResponse"
                        }
                    }
                }
            }
        },
        "/productsById": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetProductResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "CreateProductsRequest": {
            "type": "object",
            "required": [
                "type"
            ],
            "properties": {
                "type": {
                    "type": "string"
                }
            }
        },
        "CreateProductsResponse": {
            "type": "object",
            "properties": {
                "type": {
                    "type": "string"
                }
            }
        },
        "DeleteProductDataMyDataResponse": {
            "type": "object",
            "properties": {
                "type": {
                    "type": "string"
                }
            }
        },
        "DeleteProductDataResponse": {
            "type": "object",
            "properties": {
                "my_data": {
                    "$ref": "#/definitions/DeleteProductDataMyDataResponse"
                }
            }
        },
        "DeleteProductResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/DeleteProductDataResponse"
                }
            }
        },
        "GetProductResponse": {
            "type": "object",
            "properties": {
                "type": {
                    "type": "string"
                }
            }
        },
        "GetProductsResponse": {
            "type": "object",
            "properties": {
                "type": {
                    "type": "string"
                }
            }
        },
        "UpdateProductRequest": {
            "type": "object",
            "required": [
                "type"
            ],
            "properties": {
                "type": {
                    "type": "string"
                }
            }
        },
        "UpdateProductResponse": {
            "type": "object",
            "properties": {
                "type": {
                    "type": "string"
                }
            }
        }
    },
    "tags": [
        {
            "description": "API For Frontend only. Please aware that this API keeps changing by new requirements.",
            "name": "API For Frontend"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Service Name",
	Description:      "Description of the service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
