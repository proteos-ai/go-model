package knowledgemodel

// KnowledgeNodeSearchResult is a node returned by retrieval: its metadata plus a
// relevance Score. The body is excluded (retrieval ranks over it but returning
// it would move far more bytes than a result list needs). Score's scale depends
// on the match mode — for keyword search it is the Postgres `ts_rank` of the
// content against the query; a higher score is more relevant.
type KnowledgeNodeSearchResult struct {
	KnowledgeNodeMetadata
	Score float64 `json:"score"`
}
