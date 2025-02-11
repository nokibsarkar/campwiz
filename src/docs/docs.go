// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "LGPL-3.0",
            "url": "http://www.gnu.org/licenses/lgpl-3.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/campaign/": {
            "get": {
                "description": "get all campaigns",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Campaign"
                ],
                "summary": "List all campaigns",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "collectionFormat": "csv",
                        "name": "ids",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "nextToken",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.ResponseList-database_Campaign"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new campaign",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Campaign"
                ],
                "summary": "Create a new campaign",
                "parameters": [
                    {
                        "description": "The campaign request",
                        "name": "campaignRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.CampaignRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/database.Campaign"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "consts.Language": {
            "type": "string",
            "enum": [
                "en",
                "bn"
            ],
            "x-enum-varnames": [
                "English",
                "Bangla"
            ]
        },
        "database.Campaign": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "read only",
                    "type": "string"
                },
                "created_by": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "language": {
                    "$ref": "#/definitions/consts.Language"
                },
                "name": {
                    "type": "string"
                },
                "rules": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                }
            }
        },
        "routes.ResponseList-database_Campaign": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.Campaign"
                    }
                }
            }
        },
        "services.CampaignRequest": {
            "type": "object",
            "properties": {
                "created_by": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "jury": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "language": {
                    "$ref": "#/definitions/consts.Language"
                },
                "name": {
                    "type": "string"
                },
                "rules": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1",
	Host:             "localhost:8080",
	BasePath:         "/api/v2",
	Schemes:          []string{"http", "https"},
	Title:            "Campwiz API",
	Description:      "This is the API documentation for Campwiz",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
