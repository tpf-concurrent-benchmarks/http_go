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
        "/login": {
            "post": {
                "parameters": [
                    {
                        "description": "Username and hashed password",
                        "name": "userInDB",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserInDB"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "access_token",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/poll": {
            "post": {
                "parameters": [
                    {
                        "type": "string",
                        "name": "access_token",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "name": "token_type",
                        "in": "header"
                    },
                    {
                        "description": "Poll object",
                        "name": "poll",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Poll"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Poll created successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/poll/{id}": {
            "get": {
                "parameters": [
                    {
                        "type": "string",
                        "description": "Poll ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "PollWithVotes object",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Poll not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/poll/{id}/vote": {
            "post": {
                "parameters": [
                    {
                        "type": "string",
                        "name": "access_token",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "name": "token_type",
                        "in": "header"
                    },
                    {
                        "description": "Vote object",
                        "name": "vote",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Vote"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Voted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/polls": {
            "get": {
                "responses": {
                    "200": {
                        "description": "Polls object",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Polls not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/polls/{id}": {
            "delete": {
                "parameters": [
                    {
                        "type": "string",
                        "name": "access_token",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "name": "token_type",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "Poll ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Poll deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Poll not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Add a new user to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add a new user",
                "parameters": [
                    {
                        "description": "Username and hashed password",
                        "name": "userInDB",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserInDB"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User added successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/:name": {
            "get": {
                "parameters": [
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "models.Poll": {
            "type": "object",
            "properties": {
                "options": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.UserInDB": {
            "type": "object",
            "properties": {
                "hashed_password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "models.Vote": {
            "type": "object",
            "properties": {
                "option": {
                    "type": "integer"
                },
                "poll_id": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
