"use strict";

var app = {
  init: function() {
    app.activate_social_buttons();
    app.news_animation();
    app.questionnaire_nav();
    app.all_candidates_toggle();
    app.find_district();
    app.mobile_nav();
    app.homepage_toggle();
    app.candidate_table_filter();
    app.fetch_coverage();
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
        "&redirect_uri=https://www.facebook.com";
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
  /* global L, leafletPip */
  find_district: function() {
    var $map = $("#map");
    if ($map.length === 0) {
      return;
    }
    var map = L.map("map", {scrollWheelZoom: false}).setView([39.000419, -76.7591], 8);

    var info = L.control();

    var address;

    $.getJSON($map.data("map-layer"), function(districtData) {
      L.tileLayer(
        "https://cartodb-basemaps-{s}.global.ssl.fastly.net/light_all/{z}/{x}/{y}.png",
        {
          attribution:
            '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="http://cartodb.com/attributions">CartoDB</a>'
        }
      ).addTo(map);

      function onEachFeature(feature, layer) {
        if (feature.properties && feature.properties.name && feature.properties.county && feature.properties.city) {
          layer.bindPopup("<div><b>District:</b> <a href='"+window.location.href+'district-'+feature.properties.name+"'>" + feature.properties.name + "</a></div><div><b>Counties in district:</b> " + feature.properties.county + "</div><div><b>Cities/neighborhoods in district:</b> " + feature.properties.city + "</div>");
        }
        else {
          layer.bindPopup("<b>District:</b> <a href='"+window.location.href+'district-'+feature.properties.name+"'>" + feature.properties.name + "</a>");
        }
      }

      var districtLayer = L.geoJSON(districtData, {
        onEachFeature: onEachFeature
      });

      districtLayer.addTo(map);

      L.Control.geocoder({
        defaultMarkGeocode: false
      })
        .on("markgeocode", function(e) {
          var polygon = leafletPip.pointInLayer(
            [e.geocode.center.lng, e.geocode.center.lat],
            districtLayer
          );
          if (polygon.length === 0) {
            address =
              "Your address cannot be found in Maryland. Re-enter your full address and try again.";
            info.update(address);
          } else {
            var district = polygon[0].feature.properties.name;
            address = e.geocode.html;
            var bbox = e.geocode.bbox;
            var poly = L.polygon([
              bbox.getSouthEast(),
              bbox.getNorthEast(),
              bbox.getNorthWest(),
              bbox.getSouthWest()
            ]).addTo(map);
            map.fitBounds(poly.getBounds());
            districtLayer.eachLayer(function(feature) {
              if (feature.feature.properties.name === district) {
                info.update();
                info.update(address, district);
              }
            });
          }
        })
        .addTo(map);

      info.onAdd = function(map) {
        this._div = L.DomUtil.create("div", "info"); // create a div with a class "info"
        return this._div;
      };

      // method that we will use to update the control based on feature properties passed
      info.update = function(address, district) {
        this._div.innerHTML =
          '<div class="result">' +
          address +
          '<div class="district-result">District: <a href="'+ window.location.href+'district-'+district+'">' +
          district +
          "</a></div></div>";
      };

      info.addTo(map);
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
  },

  fetch_coverage: function() {
    var els = document.querySelectorAll("[data-feed-url]");
    [].slice.call(els).forEach(function(el) {
      var feedURL = el.getAttribute("data-feed-url");
      // Newsfeed script based on https://rss2json.com/
      var jsonFeedURL =
        "https://api.rss2json.com/v1/api.json?&api_key=q3gkae8uetnoaynpco9iwje8fpuqcibubkxfr5g8&count=5&rss_url=" +
        encodeURIComponent(feedURL);

      var xhr = new XMLHttpRequest();
      xhr.addEventListener("load", function() {
        var data;
        try {
          data = JSON.parse(xhr.responseText);
        } catch (ex) {
          console.warn("JSON parse error", ex);
          xhr.dispatchEvent(new Event("error"));
          return;
        }
        if (data.status !== "ok") {
          xhr.dispatchEvent(new Event("error"));
          return;
        }

        data.items.forEach(function(item) {
          var li = document.createElement("li");
          li.classList.add("article-title");
          var a = document.createElement("a");
          a.href = item.link;
          a.innerText = item.title;
          li.appendChild(a);
          el.appendChild(li);
        });
      });

      xhr.addEventListener("error", function(e) {
        el.innerHTML = "<li class='article-title'>Could not load feed.</li>";
        console.warn("JSON feed error");
      });
      xhr.open("GET", jsonFeedURL, true);
      xhr.send();
    });
  },

  homepage_toggle: function() {
    $(".js-swap-race").on("click", function() {
      $("#swap-race").toggleClass("hide-fed-state");
      $("#swap-race").toggleClass("hide-local");
      $("#fed-state-button").toggleClass("selected");
      $("#local-button").toggleClass("selected");
    });
  },

  candidate_table_filter: function() {
    // Source: https://codepen.io/chriscoyier/pen/tIuBL
    var Arr = Array.prototype;
    var _input;

    function _onInputEvent(e) {
      _input = e.target;
      var tables = document.getElementsByClassName(
        _input.getAttribute("data-table")
      );
      Arr.forEach.call(tables, function(table) {
        Arr.forEach.call(table.tBodies, function(tbody) {
          Arr.forEach.call(tbody.rows, _filter);
        });
      });
    }

    function _filter(row) {
      var text = row.textContent.toLowerCase();
      var val = _input.value.toLowerCase();
      row.style.display = text.indexOf(val) === -1 ? "none" : "table-row";
    }

    var inputs = document.getElementsByClassName("light-table-filter");
    Arr.forEach.call(inputs, function(input) {
      input.oninput = _onInputEvent;
    });
  }
};

$(document).ready(function() {
  app.init();
});
