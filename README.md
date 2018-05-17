# README

## Installation and usage

* Project requires [Hugo](https://gohugo.io) and [Yarn](https://yarnpkg.com/). The simplest way to install this on Mac is `brew install yarn hugo`.

* Run `yarn` to install JavaScript dependencies

* Run `yarn make-content` to create content files out of CSV files.

* Run `yarn css` to create CSS out of SCSS code.

* If you need test candidates, make files in content/(category)/(name).md with contents like

  ```json
  {
    "age": 57,
    "bio": "Edwards, 57, is in her fourth full term representing the 4th Congressional District in Congress. The Oxon Hill woman is a lawyer and former progressive advocate.",
    "candidateLastName": "Edwards",
    "candidateName": "Donna Edwards",
    "facebook": "donnaedwardsforsenate",
    "graduatedFrom": "Wake Forest University (BA); New Hampshire School of Law (formerly Franklin Pierce Law Center) (JD)",
    "image": "images/candidates/edwards.jpg",
    "imageNoBackground": "images/candidates-no-bg/edwards.jpg",
    "job": "U.S. Representative, 4th District",
    "livesIn": "Oxon Hill",
    "party": "Green",
    "photo": "edwards",
    "questions": [
      {
        "answer": "And Iran.\n\nIran, so far away.",
        "question": "What would you do about Iran?",
        "shortname": "Iran"
      },
      {
        "answer": "Not anymore.",
        "question": "Best name for a cat?",
        "shortname": "ISIS"
      }
    ],
    "title": "Donna Edwards",
    "twitter": "DonnaFEdwards",
    "website": "http://donnaedwardsforsenate.com"
  }
  ```

* Run `yarn server` to run a development server at [http://localhost:1313/](http://localhost:1313/).

* Run `yarn build` to create files for distribution in `public/`. Builds additionally require "scattered" and "minify". Install in current working directory with:

  - `GOBIN="$(pwd)" GOPATH="$(mktemp -d)" go get github.com/tdewolff/minify/cmd/minify`
  - `GOBIN="$(pwd)" GOPATH="$(mktemp -d)" go get github.com/carlmjohnson/scattered/cmd/scattered`


## Robocopy for primary:

This gets basic robocopy functionality working for local development.

- Use `yarn run build:robocopy` to create a robocopy executable in your project directory. You can then move it to `~/bin` or something to keep it out of the way but still runnable.
- Make the file `content/results.md` with contents like

```
+++
title = "2018 Primary Results"
type = "results-page"
+++

lorem ipsum
```

This will make the page http://localhost:1313/results/

- That page is based on `layouts/results/single.html`. Change that template to change the basic HTML for the results container page.
- The results page has a list of races for 2018 that it gets from `data/results.json`. If you want, you can change that file manually to change the races listed or regenerate the file by running `robocopy -local -results -output-dir data`. You can use the 2016 races by running `robocopy -local -results -output-dir data -metadata-src cmd/robocopy/test/Metadata.js`.
- The results page has a magic div that the JS looks for called `.js-results-container`. It uses that to figure out what to put into the `#info` box using AJAX. It reports errors downloading into the `#errors` div. It updates console.log when it downloads the page again. As written now, it redownloads the page every 30,000 milliseconds (30s).
- `robocopy` has different options you can change, but the defaults should be fine. Run `robocopy -h` or `robocopy --help` to see the options.
- As written now, you'll always want to use `robocopy` in local mode, so run `robocopy -local`.
- When you run `robocopy -local`, it will download the metadata and results from the board of elections for 2018 and make a bunch of contest pages in `dist/results/contests/{CONTEST-NUMBER}.html`.
- The individual contests are templated by `layouts-robocopy/contests.html`. This template can't use the extra functions normally used by Hugo without more work. Please ask me if you need one of those functions to be added.
- If you want to use the 2016 contest results, run `robocopy -local -metadata-src cmd/robocopy/test/Metadata.js -results-src cmd/robocopy/test/Results.js`.
- To use remote server mode, run `robocopy -dev-server` to start the dev server. (You can also change the data source with the same options as in local mode.) Change `results.md` to tell it to use the new server:

```
+++
title = "2018 Primary Results"
type = "results-page"
results-base-url = "http://localhost:9191/results/contests/"
+++

lorem ipsum
```

Remove `results-base-url` from the front-matter if you want to go back to testing in local-mode.


## Old readme

0. 2018 ELECTION SITE:

   * Requested changes: Add filtering by party
   * Timeline / Dates: end of april launch / June 14 is early voting
   * Races / Candidates: All fed & State, local for city, county, Howard, Harford, Carroll
   * Add Key Dates Page
   * Results? Tronc or Carl

POSSIBLE NEXT STEPS: - Add filtering by party - Make more pages auto populate from json file?

0. FUNCTION NOTES:

Each candidates_page.html has 2 main js files: - javascript.js has generic scripts for whole site - javascript-(race).js has race specific functions

The candidates.html pages are populated through a few functions:
A. app.toggle_fixed_nav() is called at the top of candidate pages, function is in javascript.js
Tells at what points to show / hide drop down nav

    B. init:function() pulls together a bunch of smaller functions related to questionaire. It is called on load at the bottom of javascript.js. Really only applies to candidate pages

    C. app.load_all_answers("dixon") is called at the bottom of a candidate page. The function itself is in javascript.js. It pulls together several functions but mostly populates the canadidate question answers based on json file called at top of file.

    D.  app.load_candidate_data(dixonData); populates the candidate data based on a json file called at top. Function itself is in javascript.js

    E. There are two JSON files called at the top that function as a database for these pages. One json file for data, one for the queston answers. Different races have different json files. EXAMPLE: scripts/governor-candidate-data.js and scripts/governor-questionnaire-data.js

    F. Several functions in javascript-(race).js control the question answer sharing. Need to be customized depending on the race, candidate and questions.

A Jquery load() pullS in the footer and tophat from html files in the root.

0. CREATE NEW RACE

1) Copy folder of race that is similar and rename

2. On Candidates.html:

   * Change race name in page title, social language in meta tags and script tags, omniture, page title,
     Hardcode in candidates for that race

3) On Candidate-page.html:

   * Change race name in page title, social language in meta tags and script tags, omniture, page title,

   * Hardcode in candidate names and links in the "Running Against" section

   * Change url for news feed in "var feed_url" if you choose

   * Hardcode in questions

   * Change javascript refrences for specific race in three places. Example: <script language="javascript" src="../scripts/governor-candidate-data.js"> would become <script language="javascript" src="../scripts/comptroller-candidate-data.js">.

4. Make three javascript pages

   javascript-(race).js: Swap info for canidiates and questions. This is for the SHARE QUESTION function

   (race)-candidate-data.js: JSON file for candidate data
   Complete google doc (\*see spreadsheet notes below)
   Convert to JSON (http://www.convertcsv.com/csv-to-json.htm)
   Add JSON to this js file, follow formatting

   (race)-questionaire-data.js: JSON file for candidate questionaire data
   Complete google doc (\*see spreadsheet notes below)
   Convert to JSON (http://www.convertcsv.com/csv-to-json.htm)
   Add JSON to this js file, follow formatting

* SPREADSHEET NOTES:

      	A. Cadidate Data input notes:

      	PARTY: Democrat/Republican should be capatalized
      	WEBSITE: xxx.com (no http://www.)
      	TWITTER: @XXX
      	FACEBOOK: facebook.com/xxx (no http://www.)
      	BIO:  Wrap all paragraphs in <p></p>. MUST be one long string with no returns
      	BACKGROUND:  Wrap all paragraphs in <p></p>. MUST be one long string with no returns


    B. Questionaires

    - Entries must be one long string
    - It is ok if some p2 fields are blank
    - If a p2 entry has more that one paragraph, need to wrap all graphs in  <p></p>
    - p1 entires should not be wrapped in <p>

0. NEWS FEED ON CANDIDATE PAGE:

Uses rss2json

STEP 1: Create the feed

    - Go to this page: https://rss2json.com/

    - Run xml page through their converter (Example: http://www.baltimoresun.com/news/maryland/politics/rss2.0.xml)

    - Choose advanced options and make count = 5

    - Need to have an API key to do this, just log in for one (it is free)

STEP 2: Add the feed to candidates-page.html

    Add the url provided in the converter to the variable "var feed_url" at the bottom of candidate-page.html

SOURCE CODE:
I used the AJAX code on this page, just swapped out the url provided from the converter in step 1:
https://rss2json.com/rss-to-json-api-javascript-example

0. Dev Pass
   (baltsun / data)
