package knowledgeapi

// WriteContentRequest overwrites a node's whole content body
// (`PUT /nodes/:id/content`). An empty content string is a legitimate value
// (it clears the body), so `content` is required-as-a-field but may be "".
type WriteContentRequest struct {
	Content string `json:"content"`
}

// EditContentRequest is a surgical, anchored edit of a node's content
// (`POST /nodes/:id/content/actions/edit`) — the server-side analog of Claude
// Code's `Edit` tool. `old_string` must occur (uniquely, unless `replace_all`)
// in the current body.
type EditContentRequest struct {
	OldString  string `json:"old_string" validate:"required"`
	NewString  string `json:"new_string"`
	ReplaceAll bool   `json:"replace_all"`
}

// ContentResponse is the raw body of a node plus its line count
// (`GET /nodes/:id/content`). `cat -n` style framing is an MCP-layer concern
// (LUM-99); this returns the body verbatim.
type ContentResponse struct {
	Id         string `json:"id"`
	Content    string `json:"content"`
	LinesTotal int    `json:"lines_total"`
}

// WriteContentResponse summarizes a content overwrite.
type WriteContentResponse struct {
	Id         string `json:"id"`
	LinesTotal int    `json:"lines_total"`
}

// EditContentResponse summarizes a surgical edit: how many occurrences were
// replaced and the resulting line count.
type EditContentResponse struct {
	Id           string `json:"id"`
	Replacements int    `json:"replacements"`
	LinesTotal   int    `json:"lines_total"`
}

// ── Within-node navigation: outline + grep ─────────────────────────────────
// Large nodes (mostly ingested file/url docs) shouldn't be read blindly. These
// two multi-node operations let an agent map a body's structure or grep it for
// an exact token, then read just the relevant line window with read_node_content.
// Line numbers are 1-based and align with the cat -n read framing.

// OutlineRequest asks for the structural outline of one or several nodes' bodies
// (`POST /nodes/content/actions/outline`).
type OutlineRequest struct {
	NodeIds []string `json:"node_ids" validate:"required,min=1"`
}

// OutlineHeading is one markdown heading in a node's outline plus the line range
// of the section it opens. LineEnd runs to the line before the next heading of
// the same-or-higher level (nested subsections stay inside it), or to the end of
// the body for the last section at its level.
type OutlineHeading struct {
	Level     int    `json:"level"`
	Title     string `json:"title"`
	LineStart int    `json:"line_start"`
	LineEnd   int    `json:"line_end"`
}

// NodeOutline is one node's structural map. Headings is empty for a body with no
// markdown headings (a flat doc — grep it or read it paged instead).
type NodeOutline struct {
	NodeId     string           `json:"node_id"`
	Title      string           `json:"title"`
	LinesTotal int              `json:"lines_total"`
	Headings   []OutlineHeading `json:"headings"`
}

// OutlineResponse returns one NodeOutline per existing requested node, in request
// order; ids absent from the caller's org are silently dropped.
type OutlineResponse struct {
	Results []NodeOutline `json:"results"`
}

// SearchContentRequest greps one or several nodes' bodies for a pattern
// (`POST /nodes/content/actions/search`). Substring by default; set `is_regex`
// for a Go regular expression. `context_lines` and `max_matches_per_node` are
// pointers so an explicit value is distinguishable from omitted (→ defaults).
type SearchContentRequest struct {
	NodeIds           []string `json:"node_ids" validate:"required,min=1"`
	Pattern           string   `json:"pattern" validate:"required"`
	IsRegex           bool     `json:"is_regex,omitempty"`
	IsCaseSensitive   bool     `json:"is_case_sensitive,omitempty"`
	ContextLines      *int     `json:"context_lines,omitempty"`
	MaxMatchesPerNode *int     `json:"max_matches_per_node,omitempty"`
}

// ContentMatch is one grep hit: the 1-based matching line, the line range of the
// returned snippet (the match plus its context window), and the snippet in cat -n format.
type ContentMatch struct {
	Line      int    `json:"line"`
	LineStart int    `json:"line_start"`
	LineEnd   int    `json:"line_end"`
	Snippet   string `json:"snippet"`
}

// NodeContentMatches is one node's grep result. MatchCount is the true total;
// Matches is capped at the request's max_matches_per_node and Truncated reports
// whether the cap dropped any.
type NodeContentMatches struct {
	NodeId     string         `json:"node_id"`
	Title      string         `json:"title"`
	LinesTotal int            `json:"lines_total"`
	MatchCount int            `json:"match_count"`
	Truncated  bool           `json:"truncated"`
	Matches    []ContentMatch `json:"matches"`
}

// SearchContentResponse returns one NodeContentMatches per existing requested
// node, in request order; ids absent from the caller's org are dropped.
type SearchContentResponse struct {
	Results []NodeContentMatches `json:"results"`
}
