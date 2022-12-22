// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/auth/nonce": {
            "post": {
                "description": "Generate a message for user's wallet",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Generate a message",
                "parameters": [
                    {
                        "description": "Generate message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.GenerateMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.GeneratedMessage"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/auth/nonce/verify": {
            "post": {
                "description": "Verified the generated message",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Verified the generated message",
                "parameters": [
                    {
                        "description": "Verify message request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.VerifyMessageRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.VerifyResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/configs": {
            "get": {
                "description": "Get configs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Configs"
                ],
                "summary": "Get configs",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.ConfigResp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "description": "create config",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Configs"
                ],
                "summary": "create config",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.ConfigResp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/configs/{key}": {
            "get": {
                "description": "get one config",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Configs"
                ],
                "summary": "get one config",
                "parameters": [
                    {
                        "type": "string",
                        "description": "config key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.ConfigResp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "description": "delete config",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Configs"
                ],
                "summary": "delete config",
                "parameters": [
                    {
                        "type": "string",
                        "description": "config key",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.ConfigResp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/files": {
            "post": {
                "security": [
                    {
                        "Authorization": []
                    }
                ],
                "description": "Upload file",
                "produces": [
                    "multipart/form-data"
                ],
                "tags": [
                    "Files"
                ],
                "summary": "Upload file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.FileRes"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/profile": {
            "get": {
                "security": [
                    {
                        "Authorization": []
                    }
                ],
                "description": "User profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Profile"
                ],
                "summary": "User profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.ProfileResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "Authorization": []
                    }
                ],
                "description": "Edit User profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Profile"
                ],
                "summary": "Edit User profile",
                "parameters": [
                    {
                        "description": "Update profile request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdateProfileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.ProfileResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/profile/logout": {
            "post": {
                "security": [
                    {
                        "Authorization": []
                    }
                ],
                "description": "Logout",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Profile"
                ],
                "summary": "Logout",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.LogoutResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/project": {
            "get": {
                "description": "get projects",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "get projects",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The cursor returned in the previous response (used for getting the next page).",
                        "name": "cursor",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "contract address",
                        "name": "contractAddress",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JsonResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create projects",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Create project",
                "parameters": [
                    {
                        "description": "Create profile request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateProjectReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JsonResponse"
                        }
                    }
                }
            }
        },
        "/project/{contractAddress}/tokens": {
            "get": {
                "description": "get tokens by project address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "get project's tokens",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "The cursor returned in the previous response (used for getting the next page).",
                        "name": "cursor",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "contract address",
                        "name": "contractAddress",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JsonResponse"
                        }
                    }
                }
            }
        },
        "/project/{contractAddress}/tokens/{projectID}": {
            "get": {
                "description": "get project's detail",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "get project's detail",
                "parameters": [
                    {
                        "type": "string",
                        "description": "contract address",
                        "name": "contractAddress",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token ID",
                        "name": "projectID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.JsonResponse"
                        }
                    }
                }
            }
        },
        "/token/{contractAddress}/{tokenID}": {
            "get": {
                "description": "get token uri data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "token_uri"
                ],
                "summary": "get token uri data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "contract address",
                        "name": "contractAddress",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token ID",
                        "name": "tokenID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.TokenURIResp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/trait/{contractAddress}/{tokenID}": {
            "get": {
                "description": "get token's traits",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "token_uri"
                ],
                "summary": "get token's traits",
                "parameters": [
                    {
                        "type": "string",
                        "description": "contract address",
                        "name": "contractAddress",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token ID",
                        "name": "tokenID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.JsonResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/response.TokenTraitsResp"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "request.CreateProjectReq": {
            "type": "object"
        },
        "request.GenerateMessageRequest": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                }
            }
        },
        "request.UpdateProfileRequest": {
            "type": "object",
            "properties": {
                "bio": {
                    "type": "string"
                },
                "display_name": {
                    "type": "string"
                }
            }
        },
        "request.VerifyMessageRequest": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "signature": {
                    "type": "string"
                }
            }
        },
        "response.ConfigResp": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "response.FileRes": {
            "type": "object",
            "properties": {
                "file_name": {
                    "type": "string"
                },
                "file_size": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "mime_type": {
                    "type": "string"
                },
                "uploaded_by": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "response.GeneratedMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "response.JsonResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "$ref": "#/definitions/response.RespondErr"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "response.LogoutResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "response.ProfileResponse": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "bio": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "displayName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "walletAddress": {
                    "type": "string"
                }
            }
        },
        "response.RespondErr": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "response.TokenTraitsResp": {
            "type": "object",
            "properties": {
                "attributes": {}
            }
        },
        "response.TokenURIResp": {
            "type": "object",
            "properties": {
                "animation_url": {
                    "type": "string"
                },
                "attributes": {},
                "description": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "response.VerifyResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "isVerified": {
                    "type": "boolean"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Api-Key": {
            "type": "apiKey",
            "name": "Api-Key",
            "in": "header"
        },
        "Authorization": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "",
	BasePath:    "/rederinghub.io/v1",
	Schemes:     []string{},
	Title:       "Generative.xyz APIs",
	Description: "This is a sample server Autonomous devices management server.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
