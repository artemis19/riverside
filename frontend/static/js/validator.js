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

// Generate UUID for frontend(s)
function uuidv4() {
    return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}

FRONTEND_ID = uuidv4()

// Check for required fields and allowed actions/keys
function _validate_json(json_data) {
    for (var key in json_data) {
        if (!ALLOWED_KEYS.includes(key)) {
            throw new Error(`Key ${key} not in allowed keys!`)
        }
        if (key == "action") {
            action = json_data[key]
            if (!ALLOWED_ACTIONS.includes(action)) {
                throw new Error(`Action ${action} not in allowed actions!`)
            }
        }
    }
}

// Template for websocket messages being sent to the server from the client to request data
async function viz_select(ws_conn, table = "hosts", fields = {}) {
    json_request = {
        "from": FRONTEND_ID,
        "to": "server",
        "action": "select",
        "table": table,
        "fields": fields,
    }

    _validate_json(json_request)

    var newPromiseResolution
    var ws_result

    var promise = new Promise(function(resolve, reject) {
        newPromiseResolution = resolve
        ws_message_queue.push(newPromiseResolution)
        ws_conn.send(JSON.stringify(json_request))
    });
    return promise
}