package browser

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"

	"github.com/flatgrassdotnet/forecaster/common"
	"github.com/flatgrassdotnet/forecaster/utils"
)

var ts = template.Must(template.New("search.html").ParseGlob("data/templates/browser/*.html"))

type SearchData struct {
	Map      string
	Search   string
	Packages []common.Package

	Type string
	PackageColumns
}

func Search(w http.ResponseWriter, r *http.Request) {
	sd := SearchData{
		Map:    r.Header.Get("MAP"),
		Search: r.URL.Query().Get("search"),

		Type: r.URL.Query().Get("type"),
	}

	steamid, _ := utils.SteamIDFromTicket(r.Header.Get("TICKET"))

	if sd.Type == "home" {
		var err error
		sd.PackageColumns, err = GetPackageColumns(r.Header.Get("MAP"), r.URL.Query().Get("search"))
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get package columns: %s", err))
			return
		}
	} else {
		v := make(url.Values)

		v.Set("type", sd.Type)
		v.Set("search", r.URL.Query().Get("search"))

		if sd.Type == "mine" {
			v.Del("type")
			v.Set("author", string(steamid))
		}

		resp, err := http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get package list: %s", err))
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&sd.Packages)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to decode package list: %s", err))
			return
		}
	}

	err := ts.Execute(w, sd)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
