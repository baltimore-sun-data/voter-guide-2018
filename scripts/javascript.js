//////////////////////// JAVASCRIPT FOR VOTERS GUIDE /////////////////////

//Placeholder global variable, populated by other functions and used for sharing answers
var shown_answers = {"q1":"","q2":"","q3":"","q4":"","q5":"","q6":"","q7":""};

// SPLASH PAGE FUNCTION

function swapRace() {

	if ( $("#local-races").is(":visible")  ) {

	  $('#local-races').hide();
	  $('#local-button').removeClass("selected");
	  $('#state-races').show();
	  $('#state-button').addClass("selected");

    } else {

	  $('#state-races').hide();
	  $('#state-button').removeClass("selected");
	  $('#local-races').show();
	  $('#local-button').addClass("selected");

     }
};



function share(socialMessage){


              $(".icon-twitter.top").on("click", function(){

                var tweet = socialMessage; 
                var url = window.location.href ; //Interactive URL


                var twitter_url = "https://twitter.com/intent/tweet?text="+tweet+"&url="+url+"&tw_p=tweetbutton";
                window.open(twitter_url, 'mywin','left=200,top=200,width=500,height=300,toolbar=1,resizable=0'); return false;

              });


              $(".icon-facebook.top").on("click", function(){

                var picture = "http://www.trbimg.com/img-53fdf16a/turbine/bal-baltimore-default-facebook-icon"; //Picture URL
                var title = "Baltimore Sun Voter Guide 2016"; //Post title
                var description = socialMessage;  //Post description
              

                var url = window.location.href; //Interactive URL

                  var facebook_url = "https://www.facebook.com/dialog/feed?display=popup&app_id=310302989040998&link="+url+"&picture="+picture+"&name="+title+"&description="+description+"&redirect_uri=http://www.facebook.com";        
                window.open(facebook_url, 'mywin','left=200,top=200,width=500,height=300,toolbar=1,resizable=0'); return false;

              });
              
            };



var app = {

	init: function(){

		app.news_animation();
		app.questionnaire_nav();
		app.toggle_answers();
		app.all_candidates_toggle();
		app.questionnaire_hash();
		app.mobile_nav();

	},


   

// This populates the candidate page bios based on json data

load_candidate_data: function(candidateData) {
		$("#party").html(candidateData["party"]);
		$("#job").html(candidateData["job"]);
		$("#cand-pic").find(".cand-pic").attr("src", "../images/candidates/"+candidateData["photo"].replace(/\s+/g, '-').toLowerCase()+".jpg");
		$("#nav-photo").find(".nav-photo").attr("src", "../images/candidates-no-bg/"+candidateData["photo"].replace(/\s+/g, '-').toLowerCase()+".jpg");
		$("#cand-name").html(candidateData["candidateName"]);
		$("#bio-text").html(candidateData["bio"]);
		$("#bio").find(".section-header").html("About " + candidateData["candidateLastName"]);
		$("#cand-pic").find(".cand-pic").addClass(candidateData["party"].toLowerCase());
		$("#bio").find(".section-header").addClass(candidateData["party"].toLowerCase());
		$("#background").find(".section-header").addClass(candidateData["party"].toLowerCase());
		$("#news").find(".section-header").addClass(candidateData["party"].toLowerCase());
		$("#questionnaire").find(".section-header").addClass(candidateData["party"].toLowerCase());
		$("#all-candidates").find("."+candidateData["candidateLastName"].replace(/\s+/g, '-').replace(/\./g,'').toLowerCase()).addClass("no-display");




		//if these values are NOT empty, populate

		if (candidateData["twitter"] !== "") {
			$("#social").find(".twitter-link").attr("href", "https://www.twitter.com/"+candidateData["twitter"]);
			$("#social").find(".twitter-link").html("&nbsp;&nbsp;"+candidateData["twitter"]);

			}

		if (candidateData["website"] !== "") {
			$("#social").find(".website-link").attr("href", "http://www."+candidateData["website"]);
			$("#social").find(".website-link").html("&nbsp;&nbsp;"+candidateData["website"]);

			}

		
		if (candidateData["facebook"] !== "") {
			$("#social").find(".facebook-link").attr("href", "https://www."+candidateData["facebook"]);
			$("#social").find(".facebook-link").html("&nbsp;&nbsp;"+candidateData["facebook"]);	

			}
		

		if (candidateData["age"] !== "") {
				$("#age").html("Age " + candidateData["age"]);	

			}
		
		if (candidateData["livesIn"] !== "") {
				$("#lives").html("Lives in " + candidateData["livesIn"]);	

			}

		if (candidateData["graduatedFrom"] !== "") {
				$("#graduated").html("Graduated from " + candidateData["graduatedFrom"]);

			}
		


		if (candidateData["party"] == "Democrat") {
				$("#questionnaire").addClass("dem");
					$("#social").addClass("dem");
					$("#party").addClass("blue");
					$("#cand-name").addClass("blue");
					$("#all-candidates-toggle-button").addClass("blue");
			}

		if (candidateData["party"] == "Green") {
				$("#questionnaire").addClass("gre");
					$("#social").addClass("greencolor");
					$("#party").addClass("greencolor");
					$("#cand-name").addClass("greencolor");
					$("#all-candidates-toggle-button").addClass("green");
			}


		if (candidateData["party"] == "Republican") {
				$("#social").addClass("rep");
				$("#questionnaire").addClass("rep");
				$("#party").addClass("red");
				$("#cand-name").addClass("red");
				$("#all-candidates-toggle-button-rep").addClass("red");
			}
	// If website, twitter or FB are empty, don't display the icons
		if (candidateData["website"] == "") {
			$(".icon-globe.side").addClass("no-display");
			}


		if (candidateData["twitter"] == "") {
			$(".icon-twitter.side").addClass("no-display");
			}

		if (candidateData["facebook"] == "") {
			$(".icon-facebook.side").addClass("no-display");
			}


		//Show or hide the candidate navigation based on party. 
		if (candidateData["party"] == "Democrat") {
			$(".candidate-rep").addClass("no-display");
			}

		if (candidateData["party"] == "Republican") {
			$(".candidate-dem").addClass("no-display");
			}
		if (candidateData["party"] == "Green") {
			$(".candidate-dem").addClass("no-display");
			$(".candidate-rep").addClass("no-display");

			}

	},


	//these are the toggle buttons for the cadidate navigation toggle buttons
	all_candidates_toggle: function(){

		$("#all-candidates-toggle-button").click(function(){
		    $("#all-candidates-toggle").slideToggle("fast");

 if ($("#all-candidates-toggle-button").text() == "+ SHOW MORE") 
		  { 
	     		$("#all-candidates-toggle-button").text("- SHOW LESS"); 
	  	  } 
	  	else  { 
	     		$("#all-candidates-toggle-button").text("+ SHOW MORE"); 
		};

		});


		$("#all-candidates-toggle-button-rep").click(function(){
		    $("#all-candidates-toggle-rep").slideToggle("fast");

 if ($("#all-candidates-toggle-button-rep").text() == "+ SHOW MORE") 
		  { 
	     		$("#all-candidates-toggle-button-rep").text("- SHOW LESS"); 
	  	  } 
	  	else  { 
	     		$("#all-candidates-toggle-button-rep").text("+ SHOW MORE"); 
		};

		});



	},



	news_animation: function(){

		$(document).on("mouseenter","#news ul li", function(){
			$("#news ul li").addClass("faded");
			$(this).removeClass("faded");
		}).on("mouseleave","#news ul li", function(){
			$("#news ul li").removeClass("faded");
		});
	}, 

	questionnaire_nav: function(){

		//Define a variable to house the setTimeout
		var default_text;

		$("#questionnaire-nav ul li").hover(function(){
			clearTimeout(default_text);
			$("#questionnaire-nav div").html($(this).attr("data-subject"));
		}, function(){
			default_text = setTimeout(function(){
				$("#questionnaire-nav div").html("Jump to:");
			}, 1000);
		});

		$("#questionnaire-nav ul li").click(function(){

			//We have to do some math because of the fixed nav

			//Find vertical displacement of the question we want to scroll to
			var q_position = $("#question-"+$(this).html()).offset();
			$.scrollTo(q_position.top-85, 800);
		});



	},

	questionnaire_hash: function(){

		var q = Number(window.location.hash.slice(1));

		if (q >= 1 & q <= numberOfQuestions){

			console.log("evaluated");
			var q_position = $("#question-"+q).offset();

			if ($(".container").css("width") === "1000px"){ //If desktop version...

				//Scroll to that position minus 85px to accomodate the menu				
				$.scrollTo(q_position.top-85, 0);

			} else { //If mobile...

				$.scrollTo(q_position.top, 0);

			}

		}

	},

	load_answer: function(qnum, candidate){

		//Load headshot
		//$("#question-"+qnum).find(".answer-headshot").attr("src","../images/candidates-no-bg/"+(window[candidate]["candidateLastName"].replace(/\s+/g, '-').toLowerCase())+".jpg");

		//Change name
		$("#question-"+qnum).find(".speaker").html(window[candidate]["candidateLastName"]);

		//Insert the first paragraph (done so separately bc of the headshot)
		$("#question-"+qnum).find(".first-para").html(window[candidate]["q"+qnum+"p1"]);



		// If q1 is blank, hide and questionnaire section and say there is no questionaire. 
		if(window[candidate]["q"+qnum+"p1"] == "") {
			$("#questionnaire-candidate").hide();			
			$("#message").show();
			};


		//Append the rest of the answer
		$("#question-"+qnum).find(".other-paras").append(window[candidate]["q"+qnum+"p2"]);

	},

	load_all_answers: function(candidate){

		for (var i = 1; i <=numberOfQuestions; i++){

			//Load the answer text and share texy
			app.load_answer(i,candidate);
			app.share_answer();

			//Add the "selected" class to all seven tabs of this candidate
			$("#question-"+i+" .answer ul").find('[data-cand="'+candidate+'"]').addClass("selected");

			//Add this candidates name to the "shown_answers" array, which notes which candidate's answer is currently visible (used for social puposes)
			shown_answers["q"+i] = candidate;

		}

	},

	toggle_answers: function(){

		$(".question-candidates ul li").click(function(){

			//Grab which question and which candidate
			var qnum = $(this).attr("data-q");
			var candidate = $(this).attr("data-cand");

			//Clear the .other-paras container
			$("#question-"+qnum).find(".other-paras").html("");

			//Change the text
			app.load_answer(qnum, candidate);

			//Update list style
			$(this).parent().find("li").removeClass("selected");
			$(this).addClass("selected");

			//Update the "shown_answers" array, which notes which candidate's answer is currently visible (used for social puposes)
			shown_answers["q"+qnum] = candidate;

		});

	},


share_answer: function(){

		$(".answer-social .icon-twitter").click(function(){

			var qnum = $(this).parent().data("num");
			var current_candidate = shown_answers["q"+qnum]

			var share_text = generate_share_text(current_candidate, qnum);

			var twitter_url = "https://twitter.com/intent/tweet?text="+share_text[0]+"&url="+share_text[1]+"&tw_p=tweetbutton";
			window.open(twitter_url, 'mywin','left=200,top=200,width=500,height=300,toolbar=1,resizable=0'); return false;


		});

		$(".answer-social .icon-facebook").click(function(){

			var qnum = $(this).parent().data("num");
			var current_candidate = shown_answers["q"+qnum]

			var share_text = generate_share_text(current_candidate, qnum);

			var facebook_url = "https://www.facebook.com/dialog/feed?display=popup&app_id=310302989040998&link="+share_text[1]+"&picture=http://data.baltimoresun.com/voter-guide/images/candidates/"+current_candidate+".jpg&name=Baltimore Sun Voter Guide Q%26A&description="+share_text[0]+"&redirect_uri=http://data.baltimoresun.com/voter-guide";
			window.open(facebook_url, 'mywin','left=200,top=200,width=500,height=300,toolbar=1,resizable=0'); return false;


		});		


	},

	toggle_fixed_nav: function(){

		//Establish at what vertical spot we want the nav to be revealed

		//We will add to where the inline-nav begins. We don't want that
		//to be visible and we want some cushion before the other drops down

		var reveal_point = $("#inline-nav").offset().top + 120;
		
		$(document).scroll(function(){

			var current_position = $(this).scrollTop();
			
			if (current_position >= reveal_point){
				$("nav").addClass("revealed");
			} else {
				$("nav").removeClass("revealed");
			}

		});

		//Add click/scroll functionality to the nav buttons
		$(".nav-item").click(function(){
			var section = $(this).data("section");
			var position = $("#"+section).offset().top;
			$.scrollTo(position-100, 800);
	

		});

	},

	
	mobile_nav: function(){

		$("#mobile-nav").click(function(){

			if (!$(this).hasClass("clicked")){

				$("#mobile-nav-drop").show();
				$(this).addClass("clicked");

			} else {

				$("#mobile-nav-drop").hide();
				$(this).removeClass("clicked");

			}

		});

	},

	bill_reference_nav: function(){

		$("#bill-nav ul li").click(function(){

			//Grab the year clicked
			var year = $(this).html();
			if (year === "Feedback")
				year = "feedback"

			//Jump to that year
			$.scrollTo($("#year-"+year),0);


		});

	}

}


function numberWithCommas(x) {
    return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}




$(document).ready(function(){

	//share();

	app.init();

});
