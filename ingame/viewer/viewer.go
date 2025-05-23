package viewer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"text/template"

	"github.com/flatgrassdotnet/cloudbox/common"
	"github.com/flatgrassdotnet/forecaster/utils"
	"github.com/xeonx/timeago"
)

type ViewerData struct {
	PageType string
	Item     common.Package
}

var t = template.Must(template.New("viewer.html").Funcs(template.FuncMap{"timeago": timeago.English.Format}).ParseGlob("data/templates/viewer/*.html"))

func Viewer(w http.ResponseWriter, r *http.Request) {
	var vd ViewerData

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode id: %s", err), http.StatusBadRequest)
		return
	}

	vd.PageType = r.PathValue("type")

	v := make(url.Values)
	v.Set("id", strconv.Itoa(id))

	// make api request
	resp, err := http.Get(fmt.Sprintf("%s/packages/get?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get package info: %s", err))
		return
	}

	// decode api request
	err = json.NewDecoder(resp.Body).Decode(&vd.Item)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to decode package list: %s", err))
		return
	}

	// execute template
	err = t.Execute(w, vd)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
