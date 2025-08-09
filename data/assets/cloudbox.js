function isScrolling(element) { // adapted from https://www.geeksforgeeks.org/check-whether-html-element-has-scrollbars-using-javascript/
	var res = !! element["scrollTop"]; 
	if (!res) { 
		element["scrollTop"] = 1; 
		res = !!element["scrollTop"]; 
		element["scrollTop"] = 0; 
	} 
	return res; 
}

// Metabox & Context menu
var items = document.querySelectorAll(".item > a");
if (items.length > 0) {
	var metabox = document.querySelector(".metabox");
	var metaboxTarget = metabox.querySelector(".metabox-target");
	itemHover = function(e) {
		metabox.classList.add("active");
		if (metabox.dataset.itemid != this.parentNode.dataset.itemid) {
			metabox.dataset.itemid = this.parentNode.dataset.itemid;
			metaboxTarget.innerHTML = this.parentNode.querySelector(".metabox-template").innerHTML;
			
			var metaDesc = metaboxTarget.querySelector(".meta-desc");
			while (isScrolling(metaDesc)) { // Make sure the description fits
				var scale = Number(metaDesc.dataset.scale);
				if (scale >= 14) {
					if (metaDesc.classList.contains("wrap")) {
						break;
					}
					metaDesc.classList.add("wrap");
					scale = -1;
				}
				scale++;
				metaDesc.dataset.scale = scale;
			}
		}
	};
	itemLeave = function(e) {
		metabox.classList.remove("active");
	};
	itemMenu = function(e) {
		e.preventDefault();
		var prevMenu = document.querySelector(".item dialog.contextmenu[open]");
		if (prevMenu) {
			prevMenu.close();
		}
		var thisMenu = this.parentNode.querySelector("dialog.contextmenu");
		thisMenu.show();
		
		if (window.matchMedia("(hover: none)").matches) {
			// touch devices
			var menuSize = thisMenu.getBoundingClientRect();
			
			thisMenu.style.top = (e.clientY - (menuSize.height/3)) + "px";
			thisMenu.style.left = (e.clientX - (menuSize.width/2)) + "px";
		} else {
			// mouse devices
			thisMenu.style.top = e.clientY + "px";
			thisMenu.style.left = e.clientX + "px";
		}
		
		if (thisMenu.getBoundingClientRect().top < 0) thisMenu.style.top = "0px";
		if (thisMenu.getBoundingClientRect().left < 0) thisMenu.style.left = "0px";

		if (thisMenu.getBoundingClientRect().bottom > document.body.clientHeight) thisMenu.style.top = document.body.clientHeight-thisMenu.getBoundingClientRect().height + "px";
		if (thisMenu.getBoundingClientRect().right > document.body.clientWidth)	thisMenu.style.left = document.body.clientWidth-thisMenu.getBoundingClientRect().width + "px";

		
	};
	for (i=0; i<items.length; i++) {
		items[i].addEventListener("mouseover", itemHover);
		items[i].addEventListener("mouseout", itemLeave);
		items[i].addEventListener("contextmenu", itemMenu);
		
		if (document.documentElement.classList.contains("gm13")) {
			// in-game shouldn't open context menu links in new tab as that's not supported ingame.
			items[i].parentNode.querySelectorAll("dialog.contextmenu a").forEach( function(lnk) {
				lnk.removeAttribute("target");
			});
		}
	}
	// Close any open context menus when anything else is clicked.
	document.body.addEventListener("click", function(e) {
		if (!e.target.classList.contains("contextmenu")) {
			var prevMenu = document.querySelector(".item dialog.contextmenu[open]");
			if (prevMenu) {
				prevMenu.close();
			}
		}
	});
}


// Metabox - don't overlap scrollbar
var content = document.querySelector(".content");
var bottomclouds = document.querySelector("#bottomclouds");
function checkScrollbar() {
	var scrollbarWidth = content.offsetWidth - content.clientWidth;
	if (metabox) metabox.style.right = scrollbarWidth + "px";
}
window.addEventListener("resize", checkScrollbar);
checkScrollbar();


// cbalert() - alternative to alert()
function cbalert(title, txt) {
	if (typeof txt === "undefined") {
		txt = title;
		title = "";
	}
	var dialog = document.querySelector("#alert");
	dialog.querySelector("h3").innerText = title;
	dialog.querySelector("p").innerText = txt;
	
	dialog.showModal();
}

// Translation - Ingame only
var tPrefix = "translate_";
function performTranslation() { // Perform the translation, from cached terms
	var langStrings = {};
	var elements = document.querySelectorAll("[translate]");
	for (i=0; i<elements.length; i++) {
		var el = elements[i];
		var strResult = sessionStorage.getItem(tPrefix+el.getAttribute("translate"));
		if (strResult) {
			el.innerText = strResult;
			el.setAttribute("dir", "auto"); // Allow rtl languages to display correctly
		}
	}
}
function checkNeededTranslation() { // Get list of non-cached translation strings
	var langStrings = [];
	var elements = document.querySelectorAll("[translate]");
	for (i=0; i<elements.length; i++) {
		var strKey = elements[i].getAttribute("translate");
		if (sessionStorage.getItem(tPrefix+strKey) == undefined && !(langStrings.indexOf(strKey)> -1)) {
			langStrings.push(strKey);
		}
	}
	return langStrings;
}
function requestTranslation(terms) { // Request missing terms from game
	try {
		cloudbox.GetTranslations(terms.toString(), function(res) {
			var results = JSON.parse(res);
			// loop over the results, add to sessionStorage
			for (i=0; i<terms.length; i++) {
				sessionStorage.setItem(tPrefix+terms[i], results[terms[i]]);
			}
			performTranslation(); // Now go ahead with translation
		});
	} catch(ex) {}
}

if (document.documentElement.classList.contains("gm13")) {
	var needs = checkNeededTranslation(); // Get list of non-cached translation strings 
	if (needs.length == 0) { // Any terms are cached, go ahead with translation
		performTranslation();
	} else {
		requestTranslation(needs); // Request missing terms from game
	}
}

// Show search button on Linux and Mac, due to Gmod not accepting 'enter' event on those platforms.
if ( (navigator.userAgent.toLowerCase().indexOf("linux") > -1) || (navigator.userAgent.toLowerCase().indexOf("macintosh") > -1) ) {
	document.documentElement.classList.add("show-searchbtn");
}

// Enhance news items
function DualTypeLink(str, match) {
	if (document.documentElement.classList.contains("gm13")) { // ingame
		return `<a class="textlink tt_top" onclick="cloudbox.OpenLink('https://${match}')" tooltip="Open in Steam Overlay" style="position: relative;">https://${match} <span class="linkicon external"></span></a>`;
	} else { // Desktop
		return `<a class="textlink" href="https://${match}" target="_blank">${match} <span class="linkicon external"></span></a>`;
	}
}
function DateFormatter(datestring) {
	var date = new Date(datestring);
	
	var dfOffset = (date.getTimezoneOffset()/60)*-1;
	if (dfOffset >= 0) dfOffset = "+" + dfOffset;
	var dfDay = Intl.DateTimeFormat("en-GB", {year:"numeric", month:"long", day:"2-digit"}).format(date);
	var dfTime = Intl.DateTimeFormat("en-GB", {hour:"numeric", minute:"2-digit", hour12:true}).format(date);
	
	var dateFormatted = `${dfDay}, ${dfTime} (UTC ${dfOffset})`;
	var dateTooltip = dateFormatted;
	
	var timeDif = (nowMS - date.getTime()) / 1000 / 60; // in minutes
	
	if (timeDif >= 60*24*2) { // More than 48 hours
		dateTooltip = "";
	} else if (timeDif >= 60*24) { // 1 day ~ 1d 23h 59min
		dateFormatted = "Yesterday";
	} else if (timeDif >= 60*2) { // 2 hours ~ 23h 59min
		dateFormatted = Math.floor(timeDif/60) + " hours ago";
	} else if (timeDif >= 60) { // 1 hour ~ 1h 59mins
		dateFormatted = "1 hour ago";
	} else if (timeDif >= 2) { // 2 mins ~ 59 minutes
		dateFormatted = Math.floor(timeDif) + " minutes ago";
	} else { // now ~ 1 min 59 seconds
		dateFormatted = "A minute ago";
	}
	
	return {formatted:dateFormatted, tooltip:dateTooltip};
}
var newslist = document.querySelectorAll(".newsitem");
var nowMS = (new Date()).getTime();
for (i=0; i<newslist.length; i++) {
	var item = newslist[i].querySelector(".newsformat");
	
	// Flatgrass links
	item.innerHTML = item.innerHTML.replace(/\[(flatgrass\.net[\w\/\.\#\-]*)\]/gi, DualTypeLink);
	
	// Time formatting
	var timeEl = newslist[i].querySelector("time[datetime]");

	var date = DateFormatter(timeEl.getAttribute("datetime"));

	timeEl.innerText = date.formatted;
	if (date.tooltip != "") {
		timeEl.classList.add("tt_bottomright");
		timeEl.setAttribute("tooltip", date.tooltip);
	}
}

// In-game on low-performance devices shouldn't have animations
function useReducedMotion() {
	document.documentElement.classList.add("reduced-motion");
	sessionStorage.setItem("reduced-motion", true);
}
if (sessionStorage.getItem("reduced-motion")) useReducedMotion();