package protocol

import (
	"fmt"

	"github.com/ipfs/go-cid"

	"github.com/computerscienceiscool/grid-examples/ex1-order-flow/specdocs"
)

type Profile struct {
	Name string
	Spec []byte
	CID  cid.Cid
}

var (
	OrderProfile          = mustProfile("order", specdocs.MustRead("order.md"))
	PickPackProfile       = mustProfile("pick_pack", specdocs.MustRead("pick-pack.md"))
	AccountingProfile     = mustProfile("accounting", specdocs.MustRead("accounting.md"))
	ShipmentProfile       = mustProfile("shipment", specdocs.MustRead("shipment.md"))
	KernelRegisterProfile = mustProfile("kernel_register", specdocs.MustRead("kernel-register.md"))
)

func mustProfile(name string, spec []byte) Profile {
	c, err := CIDForBytes(spec)
	if err != nil {
		panic(err)
	}
	return Profile{Name: name, Spec: spec, CID: c}
}

func ProfileByCIDText(cidText string) (Profile, error) {
	for _, profile := range []Profile{
		OrderProfile,
		PickPackProfile,
		AccountingProfile,
		ShipmentProfile,
		KernelRegisterProfile,
	} {
		if profile.CID.String() == cidText {
			return profile, nil
		}
	}
	return Profile{}, fmt.Errorf("unknown profile cid %q", cidText)
}
