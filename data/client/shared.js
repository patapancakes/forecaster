function TogglePanelVisible(panel) {
	$(panel).toggle();
	return false;
}

function Popup(url) {
	var func = function(data) { 
		$("BODY").append(data); 

		$("#PopupDialog").dialog({
			modal: true,
			buttons: {
				Ok: function() {
					$(this).dialog('close');
					$("#PopupDialog").remove();
				}
			}
		});
	}

	$.post(url, {}, func, "html");

	return false;
}

function DoPopupAction(url) {
	$("#PopupFeedback").html("loading..");

	$.post(url, {}, function(data) {$("#PopupFeedback").html(data);}, "html");
	return false;
}

function DoAction(url, target) {
	$(target).html("loading..");

	$.post(url, {}, function(data) {$(target).html(data);}, "html");
	return false;
}
