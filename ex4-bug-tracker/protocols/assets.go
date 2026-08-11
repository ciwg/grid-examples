package protocols

import _ "embed"

const (
	IssueReportSpec              = "issue-report.md"
	IssueLifecycleUpdateSpec     = "issue-lifecycle-update.md"
	IssueAttachmentReferenceSpec = "issue-attachment-reference.md"
)

//go:embed issue-report.md
var issueReport []byte

//go:embed issue-lifecycle-update.md
var issueLifecycleUpdate []byte

//go:embed issue-attachment-reference.md
var issueAttachmentReference []byte

// MustRead returns the exact profile source used to derive its pCID.
func MustRead(name string) []byte {
	switch name {
	case IssueReportSpec:
		return append([]byte(nil), issueReport...)
	case IssueLifecycleUpdateSpec:
		return append([]byte(nil), issueLifecycleUpdate...)
	case IssueAttachmentReferenceSpec:
		return append([]byte(nil), issueAttachmentReference...)
	default:
		panic("unknown protocol source " + name)
	}
}
