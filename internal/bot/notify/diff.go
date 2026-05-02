package notify

import "strings"

type segment struct {
	text     string
	removed  bool
	ellipsis bool
}

func computeDiff(oldText, newText string) []segment {
	oldLines := splitNonEmpty(oldText)
	if len(oldLines) == 0 {
		return nil
	}

	newSet := make(map[string]bool, len(oldLines))
	for _, l := range splitNonEmpty(newText) {
		newSet[strings.TrimSpace(l)] = true
	}

	type marked struct {
		text    string
		removed bool
	}
	marks := make([]marked, 0, len(oldLines))
	for _, l := range oldLines {
		marks = append(marks, marked{
			text:    l,
			removed: !newSet[strings.TrimSpace(l)],
		})
	}

	var result []segment
	i := 0
	for i < len(marks) {
		if !marks[i].removed {
			j := i
			for j < len(marks) && !marks[j].removed {
				j++
			}
			result = append(result, segment{ellipsis: true})
			i = j
		} else {
			result = append(result, segment{text: marks[i].text, removed: true})
			i++
		}
	}

	// If there are no removed lines, the diff is empty (trivial change).
	hasRemoved := false
	for _, seg := range result {
		if seg.removed {
			hasRemoved = true
			break
		}
	}
	if !hasRemoved {
		return nil
	}

	return result
}

func splitNonEmpty(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if strings.TrimSpace(l) != "" {
			out = append(out, l)
		}
	}
	return out
}
