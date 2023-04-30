package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
			{
				"id": "jks0sxuwq9zzuo3",
				"created": "2022-10-16 01:55:23.183Z",
				"updated": "2023-01-31 21:38:29.622Z",
				"name": "cards",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "bd1pp6wy",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "375zbr9q",
						"name": "cost",
						"type": "number",
						"required": false,
						"unique": false,
						"options": {
							"min": 0,
							"max": 20
						}
					},
					{
						"system": false,
						"id": "esu5fw5o",
						"name": "power",
						"type": "number",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null
						}
					},
					{
						"system": false,
						"id": "22rxfrjc",
						"name": "health",
						"type": "number",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null
						}
					},
					{
						"system": false,
						"id": "eaybcxcu",
						"name": "ability",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "icicf0zt",
						"name": "chance",
						"type": "number",
						"required": false,
						"unique": false,
						"options": {
							"min": 0,
							"max": 100
						}
					},
					{
						"system": false,
						"id": "9eiaqhu0",
						"name": "artist",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "ovbsmcol",
						"name": "legendary",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "vrucqktd",
						"name": "set",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "38c2su06fbtuxso",
							"cascadeDelete": false,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "qouanbbo",
						"name": "rarity",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "ogq7g8jxn44s1n6",
							"cascadeDelete": false,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "qrz6secv",
						"name": "type",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "dx2a9n8cdft8dp4",
							"cascadeDelete": false,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "hcqer4ku",
						"name": "colour",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "7isgvhleqxzqhj5",
							"cascadeDelete": false,
							"maxSelect": 2,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "hnvlycsv",
						"name": "subtype",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "pf1nki03sxulklg",
							"cascadeDelete": false,
							"maxSelect": 2,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "syblbsvx",
						"name": "charges",
						"type": "number",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null
						}
					},
					{
						"system": false,
						"id": "dnqw3diz",
						"name": "generates",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "jks0sxuwq9zzuo3",
							"cascadeDelete": false,
							"maxSelect": 3,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "7byfpqtu",
						"name": "reminders",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "3jb9so7x4a0jr7v",
							"cascadeDelete": false,
							"maxSelect": 3,
							"displayFields": null
						}
					}
				],
				"listRule": "type.name != \"Special\" && type.name != \"Deckbuilding\"",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "38c2su06fbtuxso",
				"created": "2022-10-16 01:55:38.341Z",
				"updated": "2023-01-31 21:38:29.626Z",
				"name": "sets",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "tm4c1qa1",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "bv1sog3h",
						"name": "sort_order",
						"type": "number",
						"required": true,
						"unique": true,
						"options": {
							"min": 1,
							"max": null
						}
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "ogq7g8jxn44s1n6",
				"created": "2022-10-16 02:00:48.679Z",
				"updated": "2023-01-31 21:38:29.627Z",
				"name": "rarities",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "hraykhhl",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "fnredoni",
						"name": "sort_order",
						"type": "number",
						"required": true,
						"unique": true,
						"options": {
							"min": 1,
							"max": 7
						}
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "dx2a9n8cdft8dp4",
				"created": "2022-10-16 02:02:39.263Z",
				"updated": "2023-01-31 21:38:29.628Z",
				"name": "types",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "93ac7j6a",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"listRule": "name != \"Special\" && name != \"Deckbuilding\"",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "7isgvhleqxzqhj5",
				"created": "2022-10-16 02:05:59.167Z",
				"updated": "2023-01-31 21:38:29.629Z",
				"name": "colours",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "4ncvoop3",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "xi3ufqn9",
						"name": "sort_order",
						"type": "number",
						"required": true,
						"unique": true,
						"options": {
							"min": 1,
							"max": 6
						}
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "wapd61gcde5iim4",
				"created": "2022-10-16 02:09:29.595Z",
				"updated": "2023-01-31 21:38:29.630Z",
				"name": "spellslingers",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "9a2f7ej4",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "jvbevtaj",
						"name": "colour",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "7isgvhleqxzqhj5",
							"cascadeDelete": false,
							"maxSelect": 2,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "eshdsakv",
						"name": "signatures",
						"type": "relation",
						"required": true,
						"unique": true,
						"options": {
							"collectionId": "jks0sxuwq9zzuo3",
							"cascadeDelete": false,
							"maxSelect": 4,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "els1ujvc",
						"name": "abilities",
						"type": "relation",
						"required": true,
						"unique": true,
						"options": {
							"collectionId": "jks0sxuwq9zzuo3",
							"cascadeDelete": false,
							"maxSelect": 5,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "pskcgvv0",
						"name": "health",
						"type": "number",
						"required": true,
						"unique": false,
						"options": {
							"min": 15,
							"max": 35
						}
					},
					{
						"system": false,
						"id": "cnexadfa",
						"name": "starter_deck",
						"type": "relation",
						"required": true,
						"unique": true,
						"options": {
							"collectionId": "hmotoqp66vj93u6",
							"cascadeDelete": false,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				],
				"listRule": "",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "hmotoqp66vj93u6",
				"created": "2022-10-16 02:11:42.187Z",
				"updated": "2023-02-04 18:54:11.204Z",
				"name": "decks",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "2kgqcs8j",
						"name": "spellslinger",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "wapd61gcde5iim4",
							"cascadeDelete": false,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "ejcsbysf",
						"name": "land",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "jks0sxuwq9zzuo3",
							"cascadeDelete": false,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "2rj9bjqq",
						"name": "cards",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "jks0sxuwq9zzuo3",
							"cascadeDelete": false,
							"maxSelect": 30,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "uqekbszj",
						"name": "splash",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "7isgvhleqxzqhj5",
							"cascadeDelete": false,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "tskxssly",
						"name": "is_public",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "dlsmcqrt",
						"name": "description",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "egedzv8r",
						"name": "code",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "e6unvake",
						"name": "image",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "2ngr12v2",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "ozzc7mon",
						"name": "owner",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "bg5qils1rdz6h0f",
							"cascadeDelete": true,
							"maxSelect": 1,
							"displayFields": [
								"display_name"
							]
						}
					},
					{
						"system": false,
						"id": "znrd8fa7",
						"name": "card_details",
						"type": "json",
						"required": true,
						"unique": false,
						"options": {}
					}
				],
				"listRule": "is_public = true || @request.auth.id = owner",
				"viewRule": "is_public = true || @request.auth.id = owner",
				"createRule": "@request.auth.id = owner",
				"updateRule": "@request.auth.id = owner",
				"deleteRule": "@request.auth.id = owner",
				"options": {}
			},
			{
				"id": "pf1nki03sxulklg",
				"created": "2022-11-17 04:10:38.859Z",
				"updated": "2023-01-31 21:38:29.634Z",
				"name": "subtypes",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "2evljhpq",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "3jb9so7x4a0jr7v",
				"created": "2022-12-28 02:40:42.660Z",
				"updated": "2023-01-31 21:38:29.636Z",
				"name": "reminders",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "qvmyavpd",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "tr8wa2gk",
						"name": "text",
						"type": "text",
						"required": true,
						"unique": true,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "bg5qils1rdz6h0f",
				"created": "2023-01-06 21:36:46.463Z",
				"updated": "2023-01-31 21:38:29.637Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "y1cn9izu",
						"name": "display_name",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": 50,
							"pattern": ""
						}
					}
				],
				"listRule": null,
				"viewRule": "",
				"createRule": "",
				"updateRule": "@request.auth.id = id",
				"deleteRule": null,
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": false,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"requireEmail": true
				}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
