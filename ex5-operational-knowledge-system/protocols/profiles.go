package protocols

import (
	_ "embed"
	"fmt"

	"github.com/ipfs/go-cid"
)

//go:embed knowledge-item.md
var knowledgeItemSpec []byte

//go:embed knowledge-approval.md
var knowledgeApprovalSpec []byte

//go:embed knowledge-evidence.md
var knowledgeEvidenceSpec []byte

//go:embed knowledge-link.md
var knowledgeLinkSpec []byte

//go:embed knowledge-responsibility.md
var knowledgeResponsibilitySpec []byte

//go:embed operational-run.md
var operationalRunSpec []byte

//go:embed operational-place.md
var operationalPlaceSpec []byte

//go:embed operational-resource.md
var operationalResourceSpec []byte

type Profile struct {
	Name string
	Spec []byte
	CID  cid.Cid
}

var (
	// Intent: Bind the first ex5 PromiseGrid-native runtime slice to the exact
	// shipped knowledge-item spec bytes rather than a handwritten constant.
	// Source: DI-mibor
	KnowledgeItemProfile = mustProfile("knowledge-item", knowledgeItemSpec)
	// Intent: Bind the second ex5 PromiseGrid-native runtime slice to the exact
	// shipped knowledge-approval spec bytes rather than a handwritten constant.
	// Source: DI-vosul
	KnowledgeApprovalProfile = mustProfile("knowledge-approval", knowledgeApprovalSpec)
	// Intent: Bind the third ex5 PromiseGrid-native runtime slice to the exact
	// shipped knowledge-evidence spec bytes rather than a handwritten constant.
	// Source: DI-kavup
	KnowledgeEvidenceProfile = mustProfile("knowledge-evidence", knowledgeEvidenceSpec)
	// Intent: Bind the fourth ex5 PromiseGrid-native runtime slice to the exact
	// shipped knowledge-link spec bytes rather than a handwritten constant.
	// Source: DI-votek
	KnowledgeLinkProfile = mustProfile("knowledge-link", knowledgeLinkSpec)
	// Intent: Bind the fifth ex5 PromiseGrid-native runtime slice to the exact
	// shipped knowledge-responsibility spec bytes rather than a handwritten
	// constant. Source: DI-sarib
	KnowledgeResponsibilityProfile = mustProfile("knowledge-responsibility", knowledgeResponsibilitySpec)
	// Intent: Bind the sixth ex5 PromiseGrid-native runtime slice to the exact
	// shipped operational-run spec bytes rather than a handwritten constant.
	// Source: DI-vamok
	OperationalRunProfile = mustProfile("operational-run", operationalRunSpec)
	// Intent: Bind the seventh ex5 PromiseGrid-native runtime slice to the
	// exact shipped operational-place spec bytes rather than a handwritten
	// constant. Source: DI-pivul
	OperationalPlaceProfile = mustProfile("operational-place", operationalPlaceSpec)
	// Intent: Bind the eighth ex5 PromiseGrid-native runtime slice to the exact
	// shipped operational-resource spec bytes rather than a handwritten
	// constant. Source: DI-pivul
	OperationalResourceProfile = mustProfile("operational-resource", operationalResourceSpec)
)

func mustProfile(name string, spec []byte) Profile {
	c, err := CIDForBytes(spec)
	if err != nil {
		panic(err)
	}
	return Profile{Name: name, Spec: spec, CID: c}
}

func ProfileByCIDText(cidText string) (Profile, error) {
	for _, profile := range []Profile{KnowledgeItemProfile, KnowledgeApprovalProfile, KnowledgeEvidenceProfile, KnowledgeLinkProfile, KnowledgeResponsibilityProfile, OperationalRunProfile, OperationalPlaceProfile, OperationalResourceProfile} {
		if profile.CID.String() == cidText {
			return profile, nil
		}
	}
	return Profile{}, fmt.Errorf("unknown profile cid %q", cidText)
}
