/* PULL RESULTS FROM JSON FEED AND DISPLAY ON RESULTS PAGES */



function raceResults(feed, raceID, location) {
     
        /* set no cache */
        $.ajaxSetup({ cache: false });
        


        $.getJSON(feed, function(data){ 
            var html = [];
            var html2 = [];
            /* loop through array */
            $.each(data, function(index, key){          

              if (key.raceid == raceID & key.fipscode == null) {

                /* identify percentage, multiply x 100 and round */
                var percent = Math.round((key.votepct) * 100);

                var updated = new Date(key.lastupdated);
                updated = updated.toString();
                updated = updated.replace("GMT-0400", " ");

                var votecount = key.votecount;

                /* if winner, add check */
                if (key.winner == true) {
                  won = "winner"
                } else {
                  won = ""
                };
                
                html.push("<table class=\"candidate-row\"><tr><td class=\"candidate name \"> ", key.first, " ", key.last, "</td><td class=\"party-col\">", key.party, "</td><td class=\"votes\">", votecount.toLocaleString('en'), "</td><td class=\"percent\"><div class=\"percent-bar-bg\"><div class=\"percent-bar\" style=\"width:", percent, "%;\"></div></div><div class=\"vote-percent\">", percent, "%</div></td></tr></table>");
                
                //pulls the first updated by and displays it up top. not the best solution for sure.
                html2.push("<span class=\"update\">Updated: " + updated + "</span>");                
        }

 
            $(location).html(html.join(''));
            $(".updated").html(html2);
            $( ".update" ).last().css( "display", "inline" );

            });

            // Get preceints data based on raceID
            $.getJSON("http://data.baltimoresun.com/news/elections2016/json/jsonPrecincts.json", function(data){ 
                var html3 = [];
                var precinctsReporting = data.precincts[raceID].precinctsreporting;
                var precinctsTotal = data.precincts[raceID].precinctstotal;
                var precinctsPercent = data.precincts[raceID].precinctsreportingpct;


                html3.push(" <div class=\"pre\"> " + precinctsReporting + " precincts reporting out of " + precinctsTotal +  ".</div>");                

                $(location).append(html3);

            });

        });
    };



