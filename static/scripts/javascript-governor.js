//////////////////////// JAVASCRIPT FOR VOTERS GUIDE /////////////////////

// This page contains variables and functions that are specific for each race

// Number of questions in the questionaire section for this specific race
// Used in two functions in javascript.js : load_all_answers (to populate answers) / questionnaire_hash (for scroll to questions)

var numberOfQuestions = 11;

// Used in the share_answer function (to share a candidates specific answer to each question)
// Associate the candidates full name with their last name in the json file
// Associate a shorthand term for each question

function generate_share_text(candidate, question) {
  var name;
  var issue;

  switch (candidate) {
    case "dickson":
      name = "Freddie Donald Dickson Jr.";
      break;
    case "edwards":
      name = "Donna F. Edwards";
      break;
    case "jaffe":
      name = "Ralph Jaffe";
      break;
    case "scaldaferri":
      name = "Theresa C. Scaldaferri";
      break;
    case "staley":
      name = "Violet Staley";
      break;
    case "taylor":
      name = "Blaine Taylor";
      break;
    case "tinus":
      name = "Ed Tinus";
      break;
    case "vanhollen":
      name = "Chris Van Hollen";
      break;
    case "young":
      name = "Lih Young";
      break;
    case "chaffee":
      name = "Chris Chaffee";
      break;
    case "connor":
      name = "Sean P. Connor";
      break;
    case "douglas":
      name = "Richard J. Douglas";
      break;
    case "graziani":
      name = "John R. Graziani";
      break;
    case "holmes":
      name = "Greg Holmes";
      break;
    case "hooe":
      name = 'Joseph D. "Joe" Hooe';
      break;
    case "kefalas":
      name = "Chrys Kefalas";
      break;
    case "mcnicholas":
      name = "Mark McNicholas";
      break;
    case "richardson":
      name = "Lynn Richardson";
      break;
    case "seda":
      name = "Anthony Seda";
      break;
    case "shawver":
      name = "Richard Shawver";
      break;
    case "szeliga":
      name = "Kathy Szeliga";
      break;
    case "wallace":
      name = "Dave Wallace";
      break;
    case "yarrington":
      name = "Garry Thomas Yarrington";
      break;
    case "flowers":
      name = "Margaret Flowers";
      break;
  }

  switch (question) {
    case 1:
      issue = "Iran";
      break;
    case 2:
      issue = "ISIS";
      break;
    case 3:
      issue = "trade";
      break;
    case 4:
      issue = "Obamacare";
      break;
    case 5:
      issue = "financial regulation";
      break;
    case 6:
      issue = "Obama's legacy";
      break;
    case 7:
      issue = "workforce issues";
      break;
    case 8:
      issue = "gun control";
      break;
    case 9:
      issue = "redistricting";
      break;
    case 10:
      issue = "Baltimore";
      break;
    case 11:
      issue = "reason for running for office";
      break;
  }

  var text = "Maryland Governor candidate " + name + " on " + issue + ":";
  var link =
    "http://data.baltimoresun.com/voter-guide-2016/mayor/" +
    candidate +
    ".html%23" +
    question;
  return [text, link];
}
