var app = {
  init: function() {
    app.activate_social_buttons();
    app.news_animation();
    app.questionnaire_nav();
    app.all_candidates_toggle();
    app.mobile_nav();
  },

  activate_social_buttons: function(socialMessage) {
    $(".icon-twitter.js-click").on("click", function(e) {
      var tweet = $(e.target).data("share-text");
      var url = window.location.href; // Interactive URL

      var twitterURL =
        "https://twitter.com/intent/tweet?text=" +
        encodeURIComponent(tweet) +
        "&url=" +
        encodeURIComponent(url) +
        "&tw_p=tweetbutton";
      window.open(
        twitterURL,
        "mywin",
        "left=200,top=200,width=500,height=300,toolbar=1,resizable=0"
      );
      return false;
    });

    $(".icon-facebook.js-click").on("click", function(e) {
      // FaceBook has deprecated all the options in the pop-up,
      // so it all needs to be controlled by the meta tags on the page.
      // See https://developers.facebook.com/docs/sharing/reference/feed-dialog
      var url = window.location.href;
      var facebookURL =
        "https://www.facebook.com/dialog/feed?display=popup&app_id=310302989040998&link=" +
        encodeURIComponent(url) +
        "&redirect_uri=http://www.facebook.com";
      window.open(
        facebookURL,
        "mywin",
        "left=200,top=200,width=500,height=300,toolbar=1,resizable=0"
      );
      return false;
    });
  },

  // these are the toggle buttons for the cadidate navigation toggle buttons
  all_candidates_toggle: function() {
    $("#all-candidates-toggle-button").click(function() {
      $("#all-candidates-toggle").slideToggle("fast");

      if ($("#all-candidates-toggle-button").text() === "+ SHOW MORE") {
        $("#all-candidates-toggle-button").text("- SHOW LESS");
      } else {
        $("#all-candidates-toggle-button").text("+ SHOW MORE");
      }
    });

    $("#all-candidates-toggle-button-rep").click(function() {
      $("#all-candidates-toggle-rep").slideToggle("fast");

      if ($("#all-candidates-toggle-button-rep").text() === "+ SHOW MORE") {
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
    // Define a variable to house the setTimeout
    var defaultText;

    $("#questionnaire-nav ul li").hover(
      function() {
        clearTimeout(defaultText);
        $("#questionnaire-nav div").html($(this).attr("data-subject"));
      },
      function() {
        defaultText = setTimeout(function() {
          $("#questionnaire-nav div").html("Jump to:");
        }, 1000);
      }
    );

    $("#questionnaire-nav ul li a").click(function(e) {
      // Find vertical displacement of the question we want to scroll to
      // We have to do some math because of the fixed nav
      var goal = /#.*?$/.exec(e.target.href)[0];
      var qPosition = $(goal).offset();
      $.scrollTo(qPosition.top - 85, 800);
      window.location = goal;
      return false;
    });
  },

  load_answer: function(qnum, candidate) {
    // Load headshot
    // $("#question-"+qnum).find(".answer-headshot").attr("src","../images/candidates-no-bg/"+(window[candidate]["candidateLastName"].replace(/\s+/g, '-').toLowerCase())+".jpg");

    // Change name
    $("#question-" + qnum)
      .find(".speaker")
      .html(window[candidate]["candidateLastName"]);

    // Insert the first paragraph (done so separately bc of the headshot)
    $("#question-" + qnum)
      .find(".first-para")
      .html(window[candidate]["q" + qnum + "p1"]);

    // If q1 is blank, hide and questionnaire section and say there is no questionaire.
    if (window[candidate]["q" + qnum + "p1"] === "") {
      $("#questionnaire-candidate").hide();
      $("#message").show();
    }

    // Append the rest of the answer
    $("#question-" + qnum)
      .find(".other-paras")
      .append(window[candidate]["q" + qnum + "p2"]);
  },

  share_answer: function() {
    console.log("TODO");
  },

  toggle_fixed_nav: function() {
    // Establish at what vertical spot we want the nav to be revealed

    // We will add to where the inline-nav begins. We don't want that
    // to be visible and we want some cushion before the other drops down

    var revealPoint = $("#inline-nav").offset().top + 120;

    $(document).scroll(function() {
      var currentPosition = $(this).scrollTop();

      if (currentPosition >= revealPoint) {
        $("nav").addClass("revealed");
      } else {
        $("nav").removeClass("revealed");
      }
    });

    // Add click/scroll functionality to the nav buttons
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
      // Grab the year clicked
      var year = $(this).html();
      if (year === "Feedback") year = "feedback";

      // Jump to that year
      $.scrollTo($("#year-" + year), 0);
    });
  }
};

$(document).ready(function() {
  app.init();
});
