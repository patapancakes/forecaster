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

	"github.com/flatgrassdotnet/cloudbox/common"
	"github.com/flatgrassdotnet/forecaster/utils"
)

type BrowserData struct {
	Map        string
	Type       string
	Category   int
	Categories []string
	Sort       string
	Packages   []common.Package
	PageNav
}

const itemsPerPage = 50

var (
	t         = template.Must(template.New("browser.html").ParseGlob("data/templates/browser/*.html"))
	pagetypes = map[string]string{
		"home": "home",

		"mine":       "mine",
		"favourites": "favorites",

		"entities": "entity",
		"weapons":  "weapon",
		"props":    "prop",
		"savemap":  "savemap",

		"maps": "map",
	}
)

func Browser(w http.ResponseWriter, r *http.Request) {
	bd := BrowserData{
		Map:  r.Header.Get("MAP"),
		Sort: r.URL.Query().Get("sort"),
	}

	// get page type
	var ok bool
	bd.Type, ok = pagetypes[r.PathValue("type")]
	if !ok {
		bd.Type = "entity"
	}

	// get page number
	bd.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if bd.Page < 1 {
		bd.Page = 1
	}

	// handle login ticket
	steamid, _ := utils.SteamIDFromTicket(r.Header.Get("TICKET"))

	// fill api request values
	v := make(url.Values)
	v.Set("type", bd.Type)
	v.Set("offset", strconv.Itoa((bd.Page-1)*itemsPerPage))
	v.Set("count", strconv.Itoa(itemsPerPage))
	if bd.Type == "mine" {
		v.Del("type")
		v.Set("author", string(steamid))
	}

	switch r.URL.Query().Get("sort") {
	case "popmap":
		v.Set("dataname", bd.Map)
		fallthrough
	case "popular":
		v.Set("sort", "popular")
	case "newmap":
		v.Set("dataname", bd.Map)
		fallthrough
	case "newest":
		v.Set("sort", "newest")
	default:
		v.Set("sort", "random")
	}

	// set tier2nav categories
	switch bd.Type {
	case "prop":
		bd.Categories = []string{"All", "Fun", "Building", "Other Games"}
	case "savemap":
		bd.Categories = []string{"All", "Fun", "Contraption", "Scene", "Assault Course", "Gun Fight"}
	case "entity":
		bd.Categories = []string{"All", "Fun", "Weapons", "Showcase", "Tools", "NPCs", "Vehicles"}
	case "weapon":
		bd.Categories = []string{"All", "Fun", "Violent", "Showcase", "Realistic", "Tools", "Other Games"}
	case "map":
		bd.Categories = []string{"Construct", "Roleplaying", "Eye Candy", "Puzzle", "Physics", "WIP", "Rats", "Contraption", "Gamemode"}
	}

	// get tier2nav category
	bd.Category, _ = strconv.Atoi(r.URL.Query().Get("category"))
	if bd.Category < 0 {
		bd.Category = 0
	}

	// make api request
	resp, err := http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get package list: %s", err))
		return
	}

	// decode api request
	err = json.NewDecoder(resp.Body).Decode(&bd.Packages)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to decode package list: %s", err))
		return
	}

	// set prev page
	bd.Prev = max(0, bd.Page-1)

	// set next page
	bd.Next = bd.Page + 1
	if len(bd.Packages) < itemsPerPage {
		bd.Next = 0
	}

	// execute template
	err = t.Execute(w, bd)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to execute template: %s", err))
		return
	}
}
