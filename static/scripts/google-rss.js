
    
google.load("feeds", "1");

// Our callback function, for when a feed is loaded.
function feedLoaded(result) {
  if (!result.error) {

    $("#news ul").html("");

    // Loop through the feeds, putting the titles onto the page.
    // Check out the result object for a list of properties returned in each entry.
    // http://code.google.com/apis/ajaxfeeds/documentation/reference.html#JSON
    for (var i = 0; i < 4; i++) {
      var entry = result.feed.entries[i];

      var date = new Date(entry.publishedDate);

      var months = Array("JAN.", "FEB.", "MARCH", "APRIL", "MAY", "JUNE", "JULY", "AUG.", "SEPT.", "OCT.", "NOV.", "DEC.");
      var date_string = months[date.getMonth()] + " " + date.getDate();

      $("#news ul").append('<li><div class="article-date"> UPDATED '+date_string+'</div><div class="article-title"> <a href="'+entry.link+'" target="_blank"> '+entry.title+' </a> </div></li>');

    }
  }
}

function on_load(url) {

  // Create a feed instance that will grab Digg's feed.
  var feed = new google.feeds.Feed(url);

  // Calling load sends the request off.  It requires a callback function.
  feed.load(feedLoaded);
}

