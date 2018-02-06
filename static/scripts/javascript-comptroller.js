//////////////////////// JAVASCRIPT FOR VOTERS GUIDE /////////////////////

// This page contains variables and functions that are specific for each race


// Number of questions in the questionaire section for this specific race
  // Used in two functions in javascript.js : load_all_answers (to populate answers) / questionnaire_hash (for scroll to questions)

var numberOfQuestions = 5;


// Used in the share_answer function (to share a candidates specific answer to each question)
  // Associate the candidates full name with their last name in the json file
  // Associate a shorthand term for each question

function generate_share_text(candidate, question) {

    var name;
    var issue;

    switch(candidate){

      case "phukan": name = "Anjali Reed Phukan"; break;
      case "franchot": name = "Peter Franchot"; break;

    }

    switch(question){

      case 1: issue = "Q1"; break;
      case 2: issue = "Q2"; break;
      case 3: issue = "Q3"; break;
      case 4: issue = "Obamacare"; break;
      case 5: issue = "financial regulation"; break;

    }

    var text = "Maryland Comptroller candidate "+name+" on "+issue+":";
    var link = "http://data.baltimoresun.com/voter-guide-2016/mayor/"+candidate+".html%23"+question;
    return [text, link];

  };
