package notify

func buildPrompt(data, today, tomorrow string) string {
	return `You are a school schedule classifier for an Israeli school. Analyze the schedule changes below and return a structured JSON result.

Schedule changes (JSON):
` + data + `

Today's date: ` + today + `
Tomorrow's date: ` + tomorrow + `

Instructions:

Step 1 — identify RESCHEDULES (highest priority): pair a deleted event with an added event where "type" and "group" match, dates differ, and the texts are identical or very similar (same core content — minor detail differences like a room number or time slot are allowed). These matched pairs are rescheduled events. Remove both events from the pool before proceeding.

Step 2 — identify MODIFICATIONS: from the remaining (unmatched) events, pair up a deleted event with an added event that share the same "date", "type", and "group" but differ in "text". Each matched pair is a modification.

Step 3 — classify each modification:
  • CLARIFICATION: the new text is essentially the same event with added detail — a time slot was added, a room number was specified, a teacher was named, or class numbers were clarified. The old text is fully contained in or superseded by the new text without losing any core information.
  • CHANGE: the content changed significantly — different subject, major restructuring, an activity was replaced or removed.

Step 3.5 — for each CLARIFICATION, identify a short note (1-2 Hebrew words) describing the specific detail that was added or clarified. Examples: "חדר", "כיתות", "שעה", "מורה", "חדר, שעה". Use the most specific and concise description.

Step 4 — filter trivial modifications. Discard a matched pair entirely if it differs ONLY in:
  • Blank lines added or removed
  • Whitespace-only differences (spaces/tabs around dashes, colons, etc.)
  • Trailing/leading whitespace or newline normalised
  • Punctuation style only (comma → dash, period added, etc.)
  • Duplicate content rearranged without adding or removing information

Apply the same filter to unmatched added/deleted events that are empty or whitespace-only.

Step 5 — output a single JSON object with this exact schema:
{
  "rescheduled":    [{"old_date":"...","new_date":"...","group":"...","type":"...","text":"...","source_url":"..."}],
  "clarifications": [{"date":"...","group":"...","type":"...","text":"...","note":"...","source_url":"..."}],
  "changes":        [{"date":"...","group":"...","type":"...","old_text":"...","new_text":"...","source_url":"..."}],
  "added":          [{"date":"...","group":"...","type":"...","text":"...","source_url":"..."}],
  "deleted":        [{"date":"...","group":"...","type":"...","text":"..."}]
}

Rules:
- All five arrays are required; use [] for empty ones.
- Copy all date and text values verbatim from the input — do NOT modify, translate, or summarise them.
- source_url handling — this is CRITICAL, follow exactly:
  • rescheduled: copy source_url from the matching added event (the one with new_date)
  • clarifications: copy source_url from the matching added event
  • changes: copy source_url from the matching added event
  • added: copy source_url EXACTLY from the corresponding input added event — if the input event has a "source_url" field, it MUST appear in the output
  • deleted: no source_url
- If an event has no source_url in the input, omit the field in the output.
- Output ONLY the JSON object. No markdown fences, no prose, no comments.`
}
