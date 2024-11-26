module main

import vweb
import json

struct App {
	vweb.Context
}

fn main() {
	http_port := 8080
	app := &App{}
	vweb.run(app, http_port)
}

@['/login'; post]
fn (mut app App) login() vweb.Result {
	body := app.req.data

	data := json.encode(body)

	return app.json(data)
}
