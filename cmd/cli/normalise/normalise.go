package main

import (
	"pj/internal/metricise"
	"regexp"
	"strings"
)

type TagRemoval struct {
	From string
	To   string
}

var (
	tagCleaning = []TagRemoval{
		// Tags to delete.
		{From: `\r`, To: ""}, // newlines
		{From: `\n`, To: ""}, // newlines
		{From: `\t`, To: ""}, // tabs
		//{From: `<br *?/>`, To: ""},                // breaks
		{From: `<em>`, To: " "},                   // emphasis
		{From: `</em>`, To: " "},                  // closing emphasis
		{From: `•`, To: " "},                      /* weird utf thing */
		{From: `<a.*?>`, To: " "},                 // anchors
		{From: `</a>`, To: " "},                   // closing anchor
		{From: `<div.*?>`, To: " "},               // divs
		{From: `</div>`, To: " "},                 // closing div
		{From: `<header.*?>`, To: " "},            // header tag
		{From: `</header>`, To: " "},              // closing header tag
		{From: `<footer.*>.*?</footer>`, To: " "}, // footer tag
		{From: `<span.*?>`, To: " "},              // spans
		{From: `</span>`, To: " "},                // closing spans
		{From: `<time.*?>`, To: " "},              // time tag
		{From: `</time>`, To: " "},                // closing time tag
		{From: `<h2.*?>`, To: " "},                // 2nd header
		{From: `</h2>`, To: " "},                  // closing 2nd header
		{From: `<p.*?>`, To: " "},                 // paragraphs
		{From: `</p>`, To: " "},                   // closing paragraphs
		{From: `<em> </em>`, To: " "},             // empty em tags

		// Dedupe spaces.
		{From: ` +`, To: " "},

		/*
			// Uncommon strings inside <strong> tags.
			{From: `<p><strong>CROSSFIT GAMES WEEK</strong></p>`, To: ""},
			{From: `<p><strong>HAPPY.*?</strong></p>`, To: ""},
			{From: `<p><strong>–.*OPTION.*?</strong></p>`, To: ""},
			{From: `<p><em>Week [56]/16`, To: `<strong>Strength</strong><p><em>Week 5/16`},
			{From: `</em></strong><em>21-15-9 reps:</em><br/>`, To: " "},
			{From: `<p><strong><em>ADJUST.*?</em></strong>.*?</p>`, To: " "},
			{From: `<strong><br/> *</strong>`, To: " "},
			{From: `<strong><br/><em>4 rounds:</em><br/></strong>`, To: " "},

			// Identifiers.
			{From: `</strong>`, To: " "},          // Delete ending strong tags.
			{From: "<strong> *?", To: "<strong>"}, //
		*/
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

	cleaned := cleanTags(metricise.Metricise(src))

	// sections identified with <strong>.*</strong>
	// split at that tag
	// then split second string at <br/> for individual movements

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
