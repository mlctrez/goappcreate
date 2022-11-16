package compo

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goappcreate/goapp"
)

var _ app.AppUpdater = (*Root)(nil)
var _ app.Mounter = (*Root)(nil)

type Root struct {
	app.Compo
}

func (r *Root) Render() app.UI {
	return app.Div().Text(goapp.RuntimeVersion())
}

func (r *Root) OnAppUpdate(ctx app.Context) {
	if goapp.IsDevelopment() && ctx.AppUpdateAvailable() {
		ctx.Reload()
	}
}

func (r *Root) OnMount(ctx app.Context) {

}
