/*

Validator to ensure JSON data being passed back and forth is correctly formatted

*/

ALLOWED_KEYS = [
	"to",
	"from",
	"action",
	"table",
	"fields",
]

ALLOWED_ACTIONS = [
	"select",
	"update",
	"delete",
	"create",
]

// Check for required fields and allowed actions/keys
function _validate_json(json_data) {
	for (var key in json_data) {
		if (! ALLOWED_KEYS.includes(key)) {
			throw new Error(`Key ${key} not in allowed keys!`)
		}
		if (key == "action") {
			action = json_data[key]
			if (! ALLOWED_ACTIONS.includes(action)) {
				throw new Error(`Action ${action} not in allowed actions!`)
			}
		}
	}
}

// Template for websocket messages being sent to the server from the client to request data
function viz_select(ws_conn, table="hosts", fields={}) {
	json_request = {
		// Temporary text for frontend
		"from": "frontend-1",
		"to": "server",
		"action": "select",
		"table": table,
		"fields": fields,
	}

	_validate_json(json_request)

	ws_conn.send(JSON.stringify(json_request))
}