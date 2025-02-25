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
        "contact": {
            "name": "Nokib Sarkar",
            "url": "https://github.com/nokibsarkar",
            "email": "nokibsarkar@gmail.com"
        },
        "license": {
            "name": "GPL-3.0",
            "url": "https://www.gnu.org/licenses/gpl-3.0.html"
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
                            "type": "string"
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
                            "$ref": "#/definitions/services.CampaignCreateRequest"
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
        },
        "/campaign/{id}": {
            "post": {
                "description": "Update a campaign",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Campaign"
                ],
                "summary": "Update a campaign",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The campaign ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "The campaign request",
                        "name": "campaignRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.CampaignUpdateRequest"
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
        },
        "/round/": {
            "get": {
                "description": "get all rounds",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Round"
                ],
                "summary": "List all rounds",
                "parameters": [
                    {
                        "type": "string",
                        "name": "campaignId",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.ResponseList-database_CampaignRound"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new round for a campaign",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Round"
                ],
                "summary": "Create a new round",
                "parameters": [
                    {
                        "description": "The round request",
                        "name": "roundRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.RoundRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.ResponseSingle-database_CampaignRound"
                        }
                    }
                }
            }
        },
        "/round/import/{roundId}": {
            "get": {
                "description": "It would be used as a server sent event stream to broadcast on the frontend about current status of the round",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Round"
                ],
                "summary": "Get the import status about a round",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The round ID",
                        "name": "roundId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.ResponseSingle-services_RoundImportSummary"
                        }
                    }
                }
            },
            "post": {
                "description": "The user would provide a round ID and a list of commons categories and the system would import images from those categories",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Round"
                ],
                "summary": "Import images from commons",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The round ID",
                        "name": "roundId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "The import from commons request",
                        "name": "ImportFromCommons",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/services.ImportFromCommonsPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/routes.ResponseSingle-services_RoundImportSummary"
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
                "am",
                "ar",
                "as",
                "af",
                "bn",
                "zh",
                "da",
                "nl",
                "en",
                "fi",
                "fr",
                "de",
                "el",
                "gu",
                "ha",
                "hi",
                "ig",
                "it",
                "ja",
                "kn",
                "ko",
                "ml",
                "mr",
                "ne",
                "no",
                "or",
                "om",
                "fa",
                "pt",
                "pa",
                "ru",
                "sa",
                "st",
                "sn",
                "sd",
                "so",
                "es",
                "sw",
                "sv",
                "ta",
                "te",
                "ts",
                "tn",
                "tr",
                "ur",
                "-",
                "ve",
                "xh",
                "yo",
                "zu"
            ],
            "x-enum-varnames": [
                "Amharic",
                "Arabic",
                "Assamese",
                "Afrikaans",
                "Bangla",
                "Chinese",
                "Danish",
                "Dutch",
                "English",
                "Finnish",
                "French",
                "German",
                "Greek",
                "Gujarati",
                "Hausa",
                "Hindi",
                "Igbo",
                "Italian",
                "Japanese",
                "Kannada",
                "Korean",
                "Malayalam",
                "Marathi",
                "Nepali",
                "Norwegian",
                "Odia",
                "Oromo",
                "Persian",
                "Portuguese",
                "Punjabi",
                "Russian",
                "Sanskrit",
                "Sesotho",
                "Shona",
                "Sindhi",
                "Somali",
                "Spanish",
                "Swahili",
                "Swedish",
                "Tamil",
                "Telugu",
                "Tsonga",
                "Tswana",
                "Turkish",
                "Urdu",
                "Undefined",
                "Venda",
                "Xhosa",
                "Yoruba",
                "Zulu"
            ]
        },
        "database.Campaign": {
            "type": "object",
            "properties": {
                "campaignId": {
                    "type": "string"
                },
                "createdAt": {
                    "description": "read only",
                    "type": "string"
                },
                "createdById": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "endDate": {
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
                "startDate": {
                    "type": "string"
                }
            }
        },
        "database.CampaignRound": {
            "type": "object",
            "properties": {
                "allowCreations": {
                    "type": "boolean"
                },
                "allowExpansions": {
                    "type": "boolean"
                },
                "allowJuryToParticipate": {
                    "type": "boolean"
                },
                "allowMultipleJudgement": {
                    "type": "boolean"
                },
                "allowedMediaTypes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.MediaType"
                    }
                },
                "blacklist": {
                    "type": "string"
                },
                "campaignId": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "createdById": {
                    "type": "string"
                },
                "dependsOnRoundId": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "endDate": {
                    "type": "string"
                },
                "isOpen": {
                    "type": "boolean"
                },
                "isPublic": {
                    "type": "boolean"
                },
                "maximumSubmissionOfSameArticle": {
                    "type": "integer"
                },
                "minimumAddedBytes": {
                    "type": "integer"
                },
                "minimumAddedWords": {
                    "type": "integer"
                },
                "minimumDurationMilliseconds": {
                    "type": "integer"
                },
                "minimumHeight": {
                    "type": "integer"
                },
                "minimumResolution": {
                    "type": "integer"
                },
                "minimumTotalBytes": {
                    "type": "integer"
                },
                "minimumTotalWords": {
                    "type": "integer"
                },
                "minimumWidth": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "roundId": {
                    "type": "string"
                },
                "secretBallot": {
                    "type": "boolean"
                },
                "serial": {
                    "type": "integer"
                },
                "startDate": {
                    "type": "string"
                },
                "totalSubmissions": {
                    "type": "integer"
                }
            }
        },
        "database.MediaType": {
            "type": "string",
            "enum": [
                "ARTICLE",
                "BITMAP",
                "AUDIO",
                "VIDEO",
                "PDF"
            ],
            "x-enum-varnames": [
                "MediaTypeArticle",
                "MediaTypeImage",
                "MediaTypeAudio",
                "MediaTypeVideo",
                "MediaTypePDF"
            ]
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
        "routes.ResponseList-database_CampaignRound": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.CampaignRound"
                    }
                }
            }
        },
        "routes.ResponseSingle-database_CampaignRound": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/database.CampaignRound"
                }
            }
        },
        "routes.ResponseSingle-services_RoundImportSummary": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/services.RoundImportSummary"
                }
            }
        },
        "services.CampaignCreateRequest": {
            "type": "object",
            "properties": {
                "allowCreations": {
                    "type": "boolean"
                },
                "allowExpansions": {
                    "type": "boolean"
                },
                "allowJuryToParticipate": {
                    "type": "boolean"
                },
                "allowMultipleJudgement": {
                    "type": "boolean"
                },
                "allowedMediaTypes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.MediaType"
                    }
                },
                "blacklist": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "endDate": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "jury": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "language": {
                    "$ref": "#/definitions/consts.Language"
                },
                "maximumSubmissionOfSameArticle": {
                    "type": "integer"
                },
                "minimumAddedBytes": {
                    "type": "integer"
                },
                "minimumAddedWords": {
                    "type": "integer"
                },
                "minimumDurationMilliseconds": {
                    "type": "integer"
                },
                "minimumHeight": {
                    "type": "integer"
                },
                "minimumResolution": {
                    "type": "integer"
                },
                "minimumTotalBytes": {
                    "type": "integer"
                },
                "minimumTotalWords": {
                    "type": "integer"
                },
                "minimumWidth": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "rules": {
                    "type": "string"
                },
                "secretBallot": {
                    "type": "boolean"
                },
                "startDate": {
                    "type": "string"
                }
            }
        },
        "services.CampaignUpdateRequest": {
            "type": "object",
            "properties": {
                "allowCreations": {
                    "type": "boolean"
                },
                "allowExpansions": {
                    "type": "boolean"
                },
                "allowJuryToParticipate": {
                    "type": "boolean"
                },
                "allowMultipleJudgement": {
                    "type": "boolean"
                },
                "allowedMediaTypes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.MediaType"
                    }
                },
                "blacklist": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "endDate": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "jury": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "language": {
                    "$ref": "#/definitions/consts.Language"
                },
                "maximumSubmissionOfSameArticle": {
                    "type": "integer"
                },
                "minimumAddedBytes": {
                    "type": "integer"
                },
                "minimumAddedWords": {
                    "type": "integer"
                },
                "minimumDurationMilliseconds": {
                    "type": "integer"
                },
                "minimumHeight": {
                    "type": "integer"
                },
                "minimumResolution": {
                    "type": "integer"
                },
                "minimumTotalBytes": {
                    "type": "integer"
                },
                "minimumTotalWords": {
                    "type": "integer"
                },
                "minimumWidth": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "rules": {
                    "type": "string"
                },
                "secretBallot": {
                    "type": "boolean"
                },
                "startDate": {
                    "type": "string"
                }
            }
        },
        "services.ImportFromCommonsPayload": {
            "type": "object",
            "properties": {
                "categories": {
                    "description": "Categories from which images will be fetched",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "services.ImportStatus": {
            "type": "string",
            "enum": [
                "success",
                "failed",
                "pending"
            ],
            "x-enum-varnames": [
                "ImportStatusSuccess",
                "ImportStatusFailed",
                "ImportStatusPending"
            ]
        },
        "services.RoundImportSummary": {
            "type": "object",
            "properties": {
                "failedCount": {
                    "type": "integer"
                },
                "failedIds": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "status": {
                    "$ref": "#/definitions/services.ImportStatus"
                },
                "successCount": {
                    "type": "integer"
                }
            }
        },
        "services.RoundRequest": {
            "type": "object",
            "properties": {
                "allowCreations": {
                    "type": "boolean"
                },
                "allowExpansions": {
                    "type": "boolean"
                },
                "allowJuryToParticipate": {
                    "type": "boolean"
                },
                "allowMultipleJudgement": {
                    "type": "boolean"
                },
                "allowedMediaTypes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/database.MediaType"
                    }
                },
                "blacklist": {
                    "type": "string"
                },
                "campaignId": {
                    "type": "string"
                },
                "dependsOnRoundId": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "endDate": {
                    "type": "string"
                },
                "isOpen": {
                    "type": "boolean"
                },
                "isPublic": {
                    "type": "boolean"
                },
                "maximumSubmissionOfSameArticle": {
                    "type": "integer"
                },
                "minimumAddedBytes": {
                    "type": "integer"
                },
                "minimumAddedWords": {
                    "type": "integer"
                },
                "minimumDurationMilliseconds": {
                    "type": "integer"
                },
                "minimumHeight": {
                    "type": "integer"
                },
                "minimumResolution": {
                    "type": "integer"
                },
                "minimumTotalBytes": {
                    "type": "integer"
                },
                "minimumTotalWords": {
                    "type": "integer"
                },
                "minimumWidth": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "secretBallot": {
                    "type": "boolean"
                },
                "serial": {
                    "type": "integer"
                },
                "startDate": {
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
