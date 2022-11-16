package compo

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func Routes() {
	app.Route("/", &Root{})
}
