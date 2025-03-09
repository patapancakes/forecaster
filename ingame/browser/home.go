package browser

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/flatgrassdotnet/forecaster/utils"
)

type HomeData struct {
	Type string
	PackageColumns
}

var th = template.Must(template.New("home.html").ParseGlob("data/templates/browser/*.html"))

func Home(w http.ResponseWriter, r *http.Request) {
	hd := HomeData{Type: "home"}

	var err error
	hd.PackageColumns, err = GetPackageColumns(r.Header.Get("MAP"), "")
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get package columns: %s", err))
		return
	}

	// execute template
	err = th.Execute(w, hd)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
