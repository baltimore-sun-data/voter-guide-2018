 function filter(party) {
  if (party == "democrats") {

        $("#democrats").show();
        $("#republicans").hide();
        $("#greens").hide();
        $('#choosePartyText').empty();
        $('#choosePartyText').text('Democrats');

    }  else if (party == "republicans") {

        $("#democrats").hide();
        $("#republicans").show();
        $("#greens").hide();
        $('#choosePartyText').empty();
        $('#choosePartyText').text('Republicans');

  }  else if (party == "greens") {

        $("#democrats").hide();
        $("#republicans").hide();
        $("#greens").show();
        $('#choosePartyText').empty();
        $('#choosePartyText').text('Greens');

  } else if (party == "all") {

        $("#democrats").show();
        $("#republicans").show();
        $('#choosePartyText').empty();
        $('#choosePartyText').text('Choose Party');

  }
}
