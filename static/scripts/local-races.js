var local_race = {

	toggle_races: function(){

		$("#race-toggle ul li").click(function(){

			//Grab which race was clicked (e.g. 1 or 2)
			var race = $(this).data("list");

			//Hide all race lists
			$(".race-lists").hide();

			//Show the clicked race
			$("#race-"+race).show();

			//Update button styles
			$("#race-toggle ul li").removeClass("selected");
			$(this).addClass("selected");

		});

	}, 

	mobile_nav_action: function(){

		//jQuery's .change() will detect when something is selected in the dropdown
		$("#local-mobile-nav-dropdown").change(function(){

			//Grab the link of the selected county
			var link = $("option:selected").data("link");

			//Redirect browser to this page
			location.href = link+".html";

		});

	}

}

$(document).ready(function(){

	local_race.toggle_races();
	local_race.mobile_nav_action();

});