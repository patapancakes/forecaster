/*
	forecaster - cloudbox frontend
	Copyright (C) 2024  patapancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package browser

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/flatgrassdotnet/forecaster/common"
	"github.com/flatgrassdotnet/forecaster/utils"
)

type Browser struct {
	InGame     bool
	HomePage   bool
	IsDarkMode bool
	Search     string
	Sort       string
	Category   string
	Packages   []common.Package
	PrevLink   string
	NextLink   string
}

const itemsPerPage = 50

var (
	categories = map[string]string{
		"mine":     "mine",
		"entities": "entity",
		"weapons":  "weapon",
		"props":    "prop",
		"saves":    "savemap",
		"maps":     "map",
	}
	t = template.Must(template.New("browser.html").Funcs(template.FuncMap{"StripHTTPS": func(url string) string { s, _ := strings.CutPrefix(url, "https:"); return s }}).ParseGlob("data/templates/browser/*.html"))
)

func Handle(w http.ResponseWriter, r *http.Request) {
	category, ok := categories[r.PathValue("category")]
	if !ok {
		http.Error(w, "unknown category", http.StatusNotFound)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	var darkmode bool
	darkCookie, err := r.Cookie("darkmode")
	if err == nil {
		darkmode = darkCookie.Value == "true"
	}

	v := make(url.Values)
	v.Set("type", category)
	v.Set("offset", strconv.Itoa((page-1)*itemsPerPage))
	v.Set("count", strconv.Itoa(itemsPerPage))

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "random"
	}

	v.Set("sort", sort)

	if r.URL.Query().Has("search") {
		v.Set("search", r.URL.Query().Get("search"))
	}

	if r.Host == "safe.cl0udb0x.com" {
		v.Set("safemode", "true")
	}

	resp, err := http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get package list: %s", err))
		return
	}

	var list []common.Package
	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to decode package list: %s", err))
		return
	}

	prev := fmt.Sprintf("?page=%d", page-1)
	if page <= 1 {
		prev = "#"
	}

	next := fmt.Sprintf("?page=%d", page+1)
	if len(list) < itemsPerPage {
		next = "#"
	}

	ingame := strings.Contains(strings.ToLower(r.UserAgent()), "gmod/") || r.Host == "toybox.garrysmod.com" || r.Host == "ingame.cl0udb0x.com" || r.Host == "safe.cl0udb0x.com"

	/*if ingame && strings.Contains(strings.ToLower(r.UserAgent()), "awesomium") {
		http.Redirect(w, r, "/assets/awesomium/awesomium.html", http.StatusSeeOther)
		return
	}*/

	err = t.Execute(w, Browser{
		InGame:     ingame,
		IsDarkMode: darkmode,
		Search:     r.URL.Query().Get("search"),
		Sort:       sort,
		Category:   category,
		Packages:   list,
		PrevLink:   prev,
		NextLink:   next,
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
