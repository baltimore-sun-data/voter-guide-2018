function filter_party(party) {
  if (party == "democrats") {
    $("#democrats").show();
    $("#republicans").hide();
    $("#greens").hide();
    $("#choosePartyText").empty();
    $("#choosePartyText").text("Democrats");
  } else if (party == "republicans") {
    $("#democrats").hide();
    $("#republicans").show();
    $("#greens").hide();
    $("#choosePartyText").empty();
    $("#choosePartyText").text("Republicans");
  } else if (party == "greens") {
    $("#democrats").hide();
    $("#republicans").hide();
    $("#greens").show();
    $("#choosePartyText").empty();
    $("#choosePartyText").text("Greens");
  } else if (party == "all") {
    $("#democrats").show();
    $("#republicans").show();
    $("#choosePartyText").empty();
    $("#choosePartyText").text("Choose Party");
  }
}

function swapRace() {
  if ($("#local-races").is(":visible")) {
    $("#local-races").hide();
    $("#local-button").removeClass("selected");
    $("#state-races").show();
    $("#state-button").addClass("selected");
  } else {
    $("#state-races").hide();
    $("#state-button").removeClass("selected");
    $("#local-races").show();
    $("#local-button").addClass("selected");
  }
}

function share(socialMessage) {
  $(".icon-twitter.top").on("click", function() {
    var tweet = socialMessage;
    var url = window.location.href; //Interactive URL

    var twitter_url =
      "https://twitter.com/intent/tweet?text=" +
      tweet +
      "&url=" +
      url +
      "&tw_p=tweetbutton";
    window.open(
      twitter_url,
      "mywin",
      "left=200,top=200,width=500,height=300,toolbar=1,resizable=0"
    );
    return false;
  });

  $(".icon-facebook.top").on("click", function() {
    var picture =
      "http://www.trbimg.com/img-53fdf16a/turbine/bal-baltimore-default-facebook-icon"; //Picture URL
    var title = "Baltimore Sun Voter Guide 2016"; //Post title
    var description = socialMessage; //Post description

    var url = window.location.href; //Interactive URL

    var facebook_url =
      "https://www.facebook.com/dialog/feed?display=popup&app_id=310302989040998&link=" +
      url +
      "&picture=" +
      picture +
      "&name=" +
      title +
      "&description=" +
      description +
      "&redirect_uri=http://www.facebook.com";
    window.open(
      facebook_url,
      "mywin",
      "left=200,top=200,width=500,height=300,toolbar=1,resizable=0"
    );
    return false;
  });
}

var app = {
  init: function() {
    app.news_animation();
    app.questionnaire_nav();
    app.all_candidates_toggle();
    app.mobile_nav();
  },

  //these are the toggle buttons for the cadidate navigation toggle buttons
  all_candidates_toggle: function() {
    $("#all-candidates-toggle-button").click(function() {
      $("#all-candidates-toggle").slideToggle("fast");

      if ($("#all-candidates-toggle-button").text() == "+ SHOW MORE") {
        $("#all-candidates-toggle-button").text("- SHOW LESS");
      } else {
        $("#all-candidates-toggle-button").text("+ SHOW MORE");
      }
    });

    $("#all-candidates-toggle-button-rep").click(function() {
      $("#all-candidates-toggle-rep").slideToggle("fast");

      if ($("#all-candidates-toggle-button-rep").text() == "+ SHOW MORE") {
        $("#all-candidates-toggle-button-rep").text("- SHOW LESS");
      } else {
        $("#all-candidates-toggle-button-rep").text("+ SHOW MORE");
      }
    });
  },

  news_animation: function() {
    $(document)
      .on("mouseenter", "#news ul li", function() {
        $("#news ul li").addClass("faded");
        $(this).removeClass("faded");
      })
      .on("mouseleave", "#news ul li", function() {
        $("#news ul li").removeClass("faded");
      });
  },

  questionnaire_nav: function() {
    //Define a variable to house the setTimeout
    var default_text;

    $("#questionnaire-nav ul li").hover(
      function() {
        clearTimeout(default_text);
        $("#questionnaire-nav div").html($(this).attr("data-subject"));
      },
      function() {
        default_text = setTimeout(function() {
          $("#questionnaire-nav div").html("Jump to:");
        }, 1000);
      }
    );

    $("#questionnaire-nav ul li").click(function() {
      //We have to do some math because of the fixed nav

      //Find vertical displacement of the question we want to scroll to
      var q_position = $("#question-" + $(this).html()).offset();
      $.scrollTo(q_position.top - 85, 800);
    });
  },

  load_answer: function(qnum, candidate) {
    //Load headshot
    //$("#question-"+qnum).find(".answer-headshot").attr("src","../images/candidates-no-bg/"+(window[candidate]["candidateLastName"].replace(/\s+/g, '-').toLowerCase())+".jpg");

    //Change name
    $("#question-" + qnum)
      .find(".speaker")
      .html(window[candidate]["candidateLastName"]);

    //Insert the first paragraph (done so separately bc of the headshot)
    $("#question-" + qnum)
      .find(".first-para")
      .html(window[candidate]["q" + qnum + "p1"]);

    // If q1 is blank, hide and questionnaire section and say there is no questionaire.
    if (window[candidate]["q" + qnum + "p1"] == "") {
      $("#questionnaire-candidate").hide();
      $("#message").show();
    }

    //Append the rest of the answer
    $("#question-" + qnum)
      .find(".other-paras")
      .append(window[candidate]["q" + qnum + "p2"]);
  },

  share_answer: function() {
    console.log("TODO");
  },

  toggle_fixed_nav: function() {
    //Establish at what vertical spot we want the nav to be revealed

    //We will add to where the inline-nav begins. We don't want that
    //to be visible and we want some cushion before the other drops down

    var reveal_point = $("#inline-nav").offset().top + 120;

    $(document).scroll(function() {
      var current_position = $(this).scrollTop();

      if (current_position >= reveal_point) {
        $("nav").addClass("revealed");
      } else {
        $("nav").removeClass("revealed");
      }
    });

    //Add click/scroll functionality to the nav buttons
    $(".nav-item").click(function() {
      var section = $(this).data("section");
      var position = $("#" + section).offset().top;
      $.scrollTo(position - 100, 800);
    });
  },

  mobile_nav: function() {
    $("#mobile-nav").click(function() {
      if (!$(this).hasClass("clicked")) {
        $("#mobile-nav-drop").show();
        $(this).addClass("clicked");
      } else {
        $("#mobile-nav-drop").hide();
        $(this).removeClass("clicked");
      }
    });
  },

  bill_reference_nav: function() {
    $("#bill-nav ul li").click(function() {
      //Grab the year clicked
      var year = $(this).html();
      if (year === "Feedback") year = "feedback";

      //Jump to that year
      $.scrollTo($("#year-" + year), 0);
    });
  }
};

function numberWithCommas(x) {
  return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

$(document).ready(function() {
  app.init();
});
