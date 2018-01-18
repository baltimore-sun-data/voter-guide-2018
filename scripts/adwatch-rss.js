
    
google.load("feeds", "1");

// Our callback function, for when a feed is loaded.
function feedLoaded(result) {
  if (!result.error) {

    $("#lead-story").html("");
    $("#reg-stories").html("");

    // Loop through the feeds, putting the titles onto the page.
    // Check out the result object for a list of properties returned in each entry.
    // http://code.google.com/apis/ajaxfeeds/documentation/reference.html#JSON
    for (var i = 0; i < result.feed.entries.length; i++) {

      //Select the story
      var entry = result.feed.entries[i];

      //Grab thumbnail
      var thumb = entry.mediaGroups[0].contents[0].thumbnails[0].url;

      //Format date

        //Convert the RSS date to a JavaScript date object
        var date = new Date(entry.publishedDate);

        //Using my own array of months for AP Style
        var months = Array("Jan.", "Feb.", "March", "April", "May", "June", "July", "Aug.", "Sept.", "Oct.", "Nov.", "Dec.");

        //Construct our formatted date
        var date_string = months[date.getMonth()] + " " + date.getDate() + ", " + date.getFullYear() + ", " + date.format("h:MM tt");

      //Format lede

        var lede = entry.content;

        //If there was a subhead (separated by some <br> tags), we need to find the cutoff
        var num = lede.indexOf("<br><br>");

        if (num > 0){
          lede = lede.substring(num+8, lede.length);
        }

        //console.log(entry.mediaGroups[0].contents[0].url);
        //console.log(entry.mediaGroups[0].contents[0].thumbnails[0].url);

        //Special treatment for most recent story
        if (i === 0){

          //Use larger version of thumb
          thumb = entry.mediaGroups[0].contents[0].url;

          $("#lead-story").append('<div class="group"><a href="'+entry.link+'"><img class="lead-image" border="0" src="'+thumb+'" /></a><div class="lead-head"><a href="'+entry.link+'"> '+entry.title+'</a></div><div class="story-date">'+date_string+'</div><div class="lead-lede">'+lede+'</div></div>');

        } else {

          $("#reg-stories").append('<div class="reg-story"><div class="group"><a href="'+entry.link+'"><img class="reg-image" border="0" src="'+thumb+'" /></a><div class="reg-head"><a href="'+entry.link+'"> '+entry.title+'</a></div><div class="story-date">'+date_string+'</div><div class="reg-lede">'+lede+'</div></div></div>');
       
        }

    }
  }
}

function on_load(url) {

  // Create a feed instance that will grab Digg's feed.
  var feed = new google.feeds.Feed(url);
  feed.setNumEntries(12);

  // Calling load sends the request off.  It requires a callback function.
  feed.load(feedLoaded);
}

