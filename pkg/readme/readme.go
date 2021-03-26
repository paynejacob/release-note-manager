package readme

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/paynejacob/release-note-manager/pkg/configuration"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Readme struct {
	notes map[int]Note
}

func NewReadme() Readme {
	return Readme{
		notes: make(map[int]Note, 0),
	}
}

func (r *Readme) AddNote(n Note) {
	r.notes[n.IssueId] = n
}

var issueNumberRegex = regexp.MustCompile("\\[#(\\d+)\\]")

const startOfReleaseNotes = "<!-- DO NOT DELETE THIS (start release notes) DO NOT DELETE THIS -->"
const endOfReleaseNotes = "<!-- DO NOT DELETE THIS (end release notes) DO NOT DELETE THIS -->"
const noteTemplate = "- %s [%d](https://github.com/%s/%s/issues/%d)\n"
const sectionTemplate = "## %s\n"

func ReadmeFromReleaseDraft(draft string) Readme {
	var line []byte
	var lineStr string
	var afterStart bool
	var err error
	var section string

	readme := NewReadme()

	for _, lineStr = range strings.Split(draft, "\n") {
		line = []byte(lineStr)
		// check if we are in the release notes section
		if !afterStart {
			if lineStr == startOfReleaseNotes {
				afterStart = true
			}

			continue
		}

		// break if we are at the end of our section
		if lineStr == endOfReleaseNotes {
			break
		}

		// check if this is a section starting
		if strings.HasPrefix(lineStr, "## ") {
			section = strings.TrimPrefix(lineStr, "## ")
			continue
		}

		// check if release note line
		if strings.HasPrefix(lineStr, "- ") {

			// get the issue number
			matches := issueNumberRegex.FindIndex(line)

			// if this line does not match the pattern we ignore it
			if len(matches) != 2 {
				continue
			}

			// get the note from the readme line
			note := Note{}
			note.Message = string(line[2:matches[0]])
			note.IssueId, err = strconv.Atoi(string(line[matches[0]:matches[1]]))
			note.Section = section

			// if we failed to parse the issue number correctly we skip this line
			if err != nil {
				continue
			}

			// write the note from the line
			readme.AddNote(note)
		}
	}

	return readme
}

func ReadmeFromIssue(issueChan chan *github.Issue, sections []configuration.Section) Readme {
	readme := NewReadme()
	var issue *github.Issue
	var note Note
	sectionMap := make(map[string]map[string]bool)

	// construct section map for fast lookups
	for i := range sections {
		sectionMap[sections[i].Name] = make(map[string]bool)
		for ii := 0; ii < len(sections[i].Labels); ii++ {
			sectionMap[sections[i].Name][sections[i].Labels[ii]] = true
		}
	}

	for issue = range issueChan {
		note = Note{Message: issue.GetTitle(), IssueId: issue.GetNumber()}

		// find the note section
		var found bool
		for section := range sectionMap {
			for i := range issue.Labels {
				if _, found = sectionMap[section][issue.Labels[i].GetName()]; found {
					note.Section = section
					break
				}
			}

			if found {
				break
			}
		}

		readme.AddNote(note)
	}

	return readme
}

func GenerateMarkdown(readme Readme, owner, repo string, sections []configuration.Section) string {
	var out string
	sectionNoteMap := make(map[string][]Note, 0)

	// make buckets for each section
	sectionNoteMap[""] = make([]Note, 0)
	for i := 0; i < len(sections); i++ {
		sectionNoteMap[sections[i].Name] = make([]Note, 0)
	}

	// group notes by section
	for _, n := range readme.notes {
		sectionNoteMap[n.Section] = append(sectionNoteMap[n.Section], n)
	}

	// for each section print the header and all notes
	for _, section := range sortedMapKeys(sectionNoteMap) {
		if len(sectionNoteMap[section]) == 0 {
			continue
		}

		if section != "" {
			out += fmt.Sprintf(sectionTemplate, section)
		}

		for i := 0; i < len(sectionNoteMap[section]); i++ {
			out += fmt.Sprintf(
				noteTemplate,
				sectionNoteMap[section][i].Message,
				sectionNoteMap[section][i].IssueId,
				owner,
				repo,
				sectionNoteMap[section][i].IssueId)
		}

		out += "\n"
	}

	return out
}

func AddLeft(right, left *Readme) {
	for k := range left.notes {
		if _, ok := right.notes[k]; !ok {
			right.notes[k] = left.notes[k]
		}
	}
}

func sortedMapKeys(in map[string][]Note) []string {
	keys := make([]string, len(in))

	var i int
	for k := range in {
		keys[i] = k
		i ++
	}

	sort.Strings(keys)

	return keys
}
