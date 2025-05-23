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

package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/flatgrassdotnet/forecaster/ingame/browser"
	"github.com/flatgrassdotnet/forecaster/ingame/viewer"
)

func main() {
	proto := flag.String("proto", "tcp", "proto for web server")
	addr := flag.String("addr", "127.0.0.1:80", "address for web server")
	flag.Parse()

	if os.Getenv("API_URL") == "" {
		os.Setenv("API_URL", "https://api.cl0udb0x.com")
	}

	// static assets
	http.Handle("GET /client/", http.StripPrefix("/client/", http.FileServer(http.Dir("data/client"))))

	// toybox.garrysmod.com
	http.HandleFunc("GET /", browser.Home)
	http.HandleFunc("GET /{type}", browser.Browser)
	http.HandleFunc("GET /search", browser.Search)

	http.HandleFunc("GET /view/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fmt.Sprintf("/view/%s/details", r.PathValue("id")), http.StatusSeeOther)
	})
	http.HandleFunc("GET /view/{id}/{type}", viewer.Viewer)

	// redundant?
	//http.HandleFunc("GET /IG/{show}", browser.Browser)

	// redirect
	http.HandleFunc("GET /ingame/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("GET /IG/maps/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/maps", http.StatusSeeOther)
	})

	// http stuff
	if *proto == "unix" {
		err := os.Remove(*addr)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("failed to delete unix socket: %s", err)
		}
	}

	l, err := net.Listen(*proto, *addr)
	if err != nil {
		log.Fatalf("failed to create web server listener: %s", err)
	}

	defer l.Close()

	if *proto == "unix" {
		err = os.Chmod(*addr, 0777)
		if err != nil {
			log.Fatalf("failed to set unix socket permissions: %s", err)
		}
	}

	http.Serve(l, nil)
}
