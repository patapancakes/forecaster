package browser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/flatgrassdotnet/forecaster/common"
)

type PageNav struct {
	Page int
	Prev int
	Next int
}

type PackageColumns struct {
	Entities []common.Package
	Weapons  []common.Package
	Props    []common.Package
	Saves    []common.Package
}

func GetPackageColumns(mapname string, search string) (PackageColumns, error) {
	var pc PackageColumns

	v := make(url.Values)

	v.Set("count", "10")
	v.Set("sort", "newest")

	if search != "" {
		v.Set("search", search)
	}

	// entities
	v.Set("type", "entity")
	resp, err := http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		return PackageColumns{}, fmt.Errorf("failed to get entity list: %s", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&pc.Entities)
	if err != nil {
		return PackageColumns{}, fmt.Errorf("failed to decode entity list: %s", err)
	}

	// weapons
	v.Set("type", "weapon")
	resp, err = http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		return PackageColumns{}, fmt.Errorf("failed to get entity list: %s", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&pc.Weapons)
	if err != nil {
		return PackageColumns{}, fmt.Errorf("failed to decode entity list: %s", err)
	}

	// props
	v.Set("type", "prop")
	resp, err = http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		return PackageColumns{}, fmt.Errorf("failed to get entity list: %s", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&pc.Props)
	if err != nil {
		return PackageColumns{}, fmt.Errorf("failed to decode entity list: %s", err)
	}

	// saves
	v.Set("type", "savemap")
	if mapname != "" {
		v.Set("dataname", mapname)
	}
	resp, err = http.Get(fmt.Sprintf("%s/packages/list?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		return PackageColumns{}, fmt.Errorf("failed to get entity list: %s", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&pc.Saves)
	if err != nil {
		return PackageColumns{}, fmt.Errorf("failed to decode entity list: %s", err)
	}

	return pc, nil
}

func SteamIDFromTicket(ticket string) ([]byte, error) {
	if ticket == "" {
		return nil, fmt.Errorf("missing ticket")
	}

	v := make(url.Values)
	v.Set("ticket", ticket)

	resp, err := http.Get(fmt.Sprintf("%s/auth/getid?%s", os.Getenv("API_URL"), v.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to get steamid: %s", err)
	}

	steamid, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read steamid: %s", err)
	}
	
	return steamid, nil
}
