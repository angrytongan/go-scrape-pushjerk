package main

import (
	"regexp"
	"strings"
)

type TagRemoval struct {
	From string
	To   string
}

var (
	tagCleaning = []TagRemoval{
		{From: `\r`, To: " "}, // newlines
		{From: `\n`, To: " "}, // newlines
		{From: `\t`, To: " "}, // tabs
		{From: `<p><strong>CROSSFIT GAMES WEEK</strong></p>`, To: " "},
		{From: `<p><strong>HAPPY APRIL FOOLS DAY!</strong></p>`, To: " "},
		{From: `<p><strong>APRIL FOOLS’ DAY!</strong></p>`, To: " "},
		{From: `<p><strong>HAPPY NEW YEAR</strong></p>`, To: " "},
		{From: `<p><strong>HAPPY NEW YEAR!</strong></p>`, To: " "},
		{From: `<p><strong>MERRY CHRISTMAS!</strong></p>`, To: " "},
		{From: `<p><strong>HAPPY VALENTINES DAY!</strong></p>`, To: " "},
		{From: `<p><strong>HAPPY HALLOWEEN!</strong></p>`, To: " "},
		{From: `<p><strong>HAPPY 4<sup>TH</sup> OF JULY!</strong></p>`, To: " "},
		{From: `<p><strong>HAPPY ST. PATRICK’S DAY</strong></p>`, To: " "},
		{From: `<p><strong>–.*OPTION.*?</strong></p>`, To: " "},
		{From: `<p><strong><em>Rest.*?</em></strong></p>`, To: " "},
		{From: `<p><em>Week [56]/16`, To: `<strong>Strength</strong><p><em>Week 6/16`},
		{From: `<strong><em> </em></strong>`, To: " "},
		{From: `</em></strong><em>21-15-9 reps:</em><br/>`, To: " "},
		{From: `<p><strong><em>ADJUST.*?</em></strong>.*?</p>`, To: " "},
		{From: `<strong><br/> *</strong>`, To: " "},
		{From: `<strong>•</strong>`, To: " "}, // bold dot things
		{From: `<strong><br/><em>4 rounds:</em><br/></strong>`, To: " "},
		{From: `<a.*?>`, To: " "},      // anchors
		{From: `</a>`, To: " "},        // closing anchor
		{From: `<div.*?>`, To: " "},    // divs
		{From: `</div>`, To: " "},      // closing div
		{From: `<header.*?>`, To: " "}, // header tag
		{From: `</header>`, To: " "},   // closing header tag
		{From: `<footer.*?>`, To: " "}, // footer tag
		{From: `</footer>`, To: " "},   // closing footer tag
		{From: `<span.*?>`, To: " "},   // spans
		{From: `</span>`, To: " "},     // closing spans
		{From: `<time.*?>`, To: " "},   // time tag
		{From: `</time>`, To: " "},     // closing time tag
		{From: `<h2.*?>`, To: " "},     // 2nd header
		{From: `</h2>`, To: " "},       // closing 2nd header
		{From: `<p.*?>`, To: " "},      // paragraphs
		{From: `</p>`, To: " "},        // closing paragraphs
		{From: `<br *?/>`, To: " "},    // breaks
		{From: `<em>`, To: " "},        // emphasis
		{From: `</em>`, To: " "},       // closing emphasis
		{From: `</strong>`, To: " "},   // ending strong tag
		{From: `•`, To: " "},           // Uh, some UTF thing
		{From: ` +`, To: " "},          // dedupe spaces

		{From: "<strong> *", To: "<strong>"}, //
	}
)

type WorkoutSection struct {
	Section string
	Content string
}

func cleanTags(src string) string {
	for _, pattern := range tagCleaning {
		re := regexp.MustCompile(pattern.From)
		src = re.ReplaceAllLiteralString(src, pattern.To)
	}

	return strings.Trim(src, " ")
}

func normalise(src string) []WorkoutSection {
	workoutSections := []WorkoutSection{}

	cleaned := cleanTags(src)

	sections := strings.Split(cleaned, "<strong>")
	if len(sections) == 0 {
		return workoutSections
	}

	for _, section := range sections[1:] {
		subsections := strings.SplitN(section, " ", 2)

		ws := WorkoutSection{
			Section: subsections[0],
			Content: subsections[1],
		}

		workoutSections = append(workoutSections, ws)
	}

	return workoutSections
}
