function isScrolling(element) { // adapted from https://www.geeksforgeeks.org/check-whether-html-element-has-scrollbars-using-javascript/
	var res = !! element["scrollTop"]; 
	if (!res) { 
		element["scrollTop"] = 1; 
		res = !!element["scrollTop"]; 
		element["scrollTop"] = 0; 
	} 
	return res; 
}

// Metabox
var items = document.querySelectorAll(".item");
var metabox = document.querySelector(".metabox");
var metaboxTarget = metabox.querySelector(".metabox-target");
itemHover = function(e) {
	metabox.classList.add("active");
	if (metabox.dataset.itemid != this.dataset.itemid) {
		metabox.dataset.itemid = this.dataset.itemid;
		metaboxTarget.innerHTML = this.querySelector(".metabox-template").innerHTML;
		
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
for (i=0; i<items.length; i++) {
	items[i].addEventListener("mouseover", itemHover);
	items[i].addEventListener("mouseout", itemLeave);
}


// Clouds - don't overlap scrollbar
var content = document.querySelector(".content");
var bottomclouds = document.querySelector("#bottomclouds");
function checkScrollbar() {
	var scrollbarWidth = content.offsetWidth - content.clientWidth;
	bottomclouds.style.right = scrollbarWidth + "px";
	metabox.style.right = scrollbarWidth + "px";
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
	
	if (document.documentElement.classList.contains("no-dialog")) {
		cbalertfallbackModal();
	} else {
		dialog.showModal();
	}
}

function cbalertfallbackModal() {
	var dialog = document.querySelector("#alert");
	var backdrop = document.querySelector("#backdrop");
	dialog.setAttribute("open", "open");
	backdrop.setAttribute("open", "open");
}

function cbalertClose() {
	var dialog = document.querySelector("#alert");
	if (document.documentElement.classList.contains("no-dialog")) {
		var backdrop = document.querySelector("#backdrop");
		dialog.removeAttribute("open");
		backdrop.removeAttribute("open");
	} else {
		dialog.close();
	}
}


if (typeof HTMLDialogElement !== 'function') {
	// Fallback handling for Awesomium and GM12
	document.documentElement.classList.add("no-dialog");
	
	var backdrop = document.createElement("div");
	backdrop.setAttribute("id","backdrop");
	backdrop.addEventListener("click", cbalertClose);
	document.body.insertBefore(backdrop, document.querySelector("#alert"));
	
	document.addEventListener("keydown", function(e){
		if (event.keyCode == 27) {
			cbalertClose();
		}
	});
}


// Translation - GM13 only
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
		if (navigator.userAgent.indexOf("Chrome/18.") > -1) { // Awesomium, thanks garry
		
			var results = JSON.parse(cloudbox.GetTranslations(terms.toString()));
			// loop over the results, add to sessionStorage
			for (i=0; i<terms.length; i++) {
				sessionStorage.setItem(tPrefix+terms[i], results[terms[i]]);
			}
			performTranslation(); // Now go ahead with translation
		
		} else { // x86-64
			cloudbox.GetTranslations(terms.toString(), function(res) {
				var results = JSON.parse(res);
				// loop over the results, add to sessionStorage
				for (i=0; i<terms.length; i++) {
					sessionStorage.setItem(tPrefix+terms[i], results[terms[i]]);
				}
				performTranslation(); // Now go ahead with translation
			});
		}
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
var newslist = document.querySelectorAll(".newsitem");
var nowMS = (new Date()).getTime();
for (i=0; i<newslist.length; i++) {
	var item = newslist[i].querySelector(".newsformat");
	
	// Flatgrass link
	if (document.documentElement.classList.contains("gm13")) { // ingame
		var flatgrassLinkUnique = "class='textlink tt_top' onclick='cloudbox.OpenLink(\"flatgrass\")' tooltip='Open in Steam Overlay' style='position: relative;'>https://flatgrass.net";
	} else { // out of game
		var flatgrassLinkUnique = "class='textlink' href='https://flatgrass.net' target='_blank'>flatgrass.net";
	}
	
	item.innerHTML = item.innerHTML.replace("https://flatgrass.net", "<a " + flatgrassLinkUnique + " <span class='linkicon' style='background-image:url(\"/assets/rustmb/combined/link.png\")'></span></a>");
	
	
	// Time formatting
	if (typeof Intl === 'object' && typeof Intl.DateTimeFormat === 'function') {
		var timeEl = newslist[i].querySelector("time[datetime]");
		
		if (timeEl != undefined) {
		
			var date = new Date(timeEl.getAttribute("datetime"));
			
			var offset = (date.getTimezoneOffset()/60)*-1;
			if (offset >= 0) offset = "+" + offset;
			var dateFormatted = Intl.DateTimeFormat("en-GB", {year:"numeric", month:"long", day:"2-digit"}).format(date) + ", " + Intl.DateTimeFormat("en-GB", {hour:"numeric", minute:"2-digit", hour12:true}).format(date) + " (UTC "+ offset + ")";
			
			var showTT = true;
			
			var timeDif = (nowMS - date.getTime()) / 1000 / 60; // in minutes
			if (timeDif < 2) {
				timeEl.innerText = "A minute ago";
			} else if (timeDif < 60) {
				timeEl.innerText = Math.floor(timeDif) + " minutes ago";
			} else if (timeDif < 60*2) {
				timeEl.innerText = "1 hour ago";
			} else if (timeDif < (60*24)) {
				timeEl.innerText = Math.floor(timeDif/60) + " hours ago";
			} else if (timeDif < (60*24*2)) {
				timeEl.innerText = "Yesterday";
			} else {
				timeEl.innerText = dateFormatted;
				showTT = false;
			}
			
			if (showTT) {
				if (document.documentElement.classList.contains("gm13")) {
					timeEl.classList.add("tt_top");
					timeEl.setAttribute("tooltip", dateFormatted);
				} else {
					timeEl.setAttribute("title", dateFormatted);
				}
			}
			
		}
	}
}