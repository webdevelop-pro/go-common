package error

type SourceLocation struct {
	FilePath     string `json:"filePath"`
	LineNumber   int    `json:"lineNumber"`
	FunctionName string `json:"functionName"`
}

type SourceReference struct {
	Repository string `json:"repository"`
	RevisionID string `json:"revisionId"`
}
