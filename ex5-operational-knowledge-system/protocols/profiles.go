package protocols

import (
	_ "embed"
	"fmt"

	"github.com/ipfs/go-cid"
)

//go:embed knowledge-item.md
var knowledgeItemSpec []byte

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
)

func mustProfile(name string, spec []byte) Profile {
	c, err := CIDForBytes(spec)
	if err != nil {
		panic(err)
	}
	return Profile{Name: name, Spec: spec, CID: c}
}

func ProfileByCIDText(cidText string) (Profile, error) {
	for _, profile := range []Profile{KnowledgeItemProfile} {
		if profile.CID.String() == cidText {
			return profile, nil
		}
	}
	return Profile{}, fmt.Errorf("unknown profile cid %q", cidText)
}
