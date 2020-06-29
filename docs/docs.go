// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/block/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "block"
                ],
                "summary": "分页获取区块列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/controllers.Block"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/block/transaction/{hash}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "block"
                ],
                "summary": "交易详情",
                "parameters": [
                    {
                        "type": "string",
                        "description": "tx hash",
                        "name": "hash",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/controllers.TransactionReceipt"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/block/transaction/{number}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "block"
                ],
                "summary": "分页列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "blockNumber",
                        "name": "number",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/controllers.Transaction"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/chain/{chain_id}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chain"
                ],
                "summary": "删除链信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "chainId",
                        "name": "chain_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功删除链信息",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/node": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "node"
                ],
                "summary": "添加node",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "关联链Id",
                        "name": "chain_id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "节点地址",
                        "name": "address",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "端口",
                        "name": "port",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "名称",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "NodeId",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/node/change": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "node"
                ],
                "summary": "切换node",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "useId",
                        "name": "user_id",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "nodeId",
                        "name": "node_id",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/node/current": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "node"
                ],
                "summary": "获取当前登录账户的节点",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/controllers.Node"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/node/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "node"
                ],
                "summary": "节点列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/controllers.Node"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "用户名",
                        "name": "username",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "密码",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"token\": token, \"user_id\": loadUser.ID}",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/user/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "description": "用户名",
                        "name": "username",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "密码",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\":0,\"data\":\"用户Id\"}",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/wallet": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "添加钱包",
                "parameters": [
                    {
                        "type": "string",
                        "description": "钱包昵称",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "keystore string",
                        "name": "content",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "walletId",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/wallet/list": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallet"
                ],
                "summary": "钱包列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.JSONResult"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dao.Wallet"
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
        "controllers.Block": {
            "type": "object",
            "properties": {
                "difficulty": {
                    "type": "string"
                },
                "gasLimit": {
                    "type": "integer"
                },
                "gasUsed": {
                    "type": "integer"
                },
                "miner": {
                    "type": "string"
                },
                "nonce": {
                    "type": "integer"
                },
                "number": {
                    "type": "string"
                },
                "parentHash": {
                    "type": "string"
                },
                "sha3Uncles": {
                    "type": "string"
                },
                "timestamp": {
                    "type": "integer"
                },
                "transactions": {
                    "type": "integer"
                }
            }
        },
        "controllers.JSONResult": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "controllers.Node": {
            "type": "object",
            "required": [
                "address",
                "chain_id",
                "name",
                "port"
            ],
            "properties": {
                "address": {
                    "description": "地址",
                    "type": "string"
                },
                "chain_id": {
                    "description": "链id",
                    "type": "integer"
                },
                "chain_name": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "is_https": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "network_id": {
                    "type": "integer"
                },
                "port": {
                    "description": "端口",
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "controllers.Transaction": {
            "type": "object",
            "properties": {
                "from": {
                    "type": "string"
                },
                "gas": {
                    "type": "integer"
                },
                "gasPrice": {
                    "type": "string"
                },
                "hash": {
                    "type": "string"
                },
                "input": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "nonce": {
                    "type": "integer"
                },
                "to": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "controllers.TransactionReceipt": {
            "type": "object",
            "properties": {
                "blockHash": {
                    "type": "string"
                },
                "blockNumber": {
                    "type": "string"
                },
                "contractAddress": {
                    "type": "string"
                },
                "cumulativeGasUsed": {
                    "type": "integer"
                },
                "gas": {
                    "type": "integer"
                },
                "gasPrice": {
                    "type": "string"
                },
                "gasUsed": {
                    "type": "integer"
                },
                "input": {
                    "type": "string"
                },
                "nonce": {
                    "type": "integer"
                },
                "root": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "status": {
                    "type": "string"
                },
                "to": {
                    "type": "string"
                },
                "transactionHash": {
                    "type": "string"
                },
                "transactionIndex": {
                    "type": "integer"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "dao.Wallet": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "content": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "name": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
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
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
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
	swag.Register(swag.Name, &s{})
}
