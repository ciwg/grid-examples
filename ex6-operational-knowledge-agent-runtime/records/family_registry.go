package records

// packageFamilyPCIDs is the fixed built-in family registry. The specification
// files are immutable and the verification test binds every literal here to
// its exact specification bytes. Source: DI-solan.
var packageFamilyPCIDs = map[string]string{
	"moks.context.place.v1":          "bafkreibczv3ytaah2bgoelclio275nblrps4lrtwjvsa5rwo6r3aybndim",
	"moks.context.resource.v1":       "bafkreidrnhbzqmhxstqq2ev325hwws6qeucaunumk7p7tbiwesnbs7rwma",
	"moks.context.responsibility.v1": "bafkreibwa7cc65rmie3z5tltubvz3uz34dinpsvm3nroczcb6exvaitidq",
	"moks.correctiveaction.event.v1": "bafkreicnlp64sd7obleleg5exm7c2yysg2wjv6wnmzx5nvx4uhpcftwvya",
	"moks.inventory.count.v1":        "bafkreiftbdmcewooqz73ise7klyy7zpzrlqa2q5p5hpf4mcy5sesex33a4",
	"moks.inventory.item.v1":         "bafkreiguud5rv5ybe7dpwmqq4xf2lkjmh5wajvpl2rfew2q5iu4ewathwm",
	"moks.inventory.reconcile.v1":    "bafkreidxt63iq3zazcj6jdtuzkuwyyaht56fuustik6um7vyielpor7kiy",
	"moks.knowledge.item.v1":         "bafkreigcuuwxlm3c2bvj73mvftzo4qbkyq5pqecx7zgtf4bnmdafpwjk64",
	"moks.knowledge.lifecycle.v1":    "bafkreihffm26lnfpfzdn466pomobsmh2lyhl4kwfuuik3pzf4axm2ke2eu",
	"moks.knowledge.revision.v1":     "bafkreidnyfrg5eifflrbjpb7kckaoxlwca77swygm5ae2cdo6wweccgs6e",
	"moks.links.typed.v1":            "bafkreihhg4mxoh37frbzw22ujyuigu546qkn5pr5iyy4naxd6dwr3udsyu",
	"moks.maintenance.finding.v1":    "bafkreif2pt5vcyijajwrurebyk2eoa4h77gxtj7iq7suwuw2egaogksw7u",
	"moks.maintenance.item.v1":       "bafkreies7a4w3bpnl7beri3ofky3lewiztwkmalzkb3bu3mgkfh5vfklyi",
	"moks.maintenance.service.v1":    "bafkreibsdabewazp4ibxcjp2uwyddczv2qim6dqyg7rckrzlzgpadvppi4",
	"moks.ops.note.v1":               "bafkreia3ctr6uzj3ygnpdastrn4fydsvknnonkdpgxsydyj5qvz2qcs5gy",
	"moks.procedures.item.v1":        "bafkreid3xltxowg4uftabtchtmqcdnr6vcmhxd7ziqmrxxe3v5ew6vu3yi",
	"moks.procedures.use.v1":         "bafkreic3x6pdeeazo66unr637j5nib3lkp7y26kwrocltiwkjejxds6uve",
	"moks.quarantine.event.v1":       "bafkreicqudjaew2khqckm7yasfsthrdwldps26wcsm7rjg6aswagjp7md4",
	"moks.receiving.disposition.v1":  "bafkreifvlesjb2sqwvdholdrhpvttsmrrjka7as2ig3ijqm3ypnywnhfp4",
	"moks.receiving.item.v1":         "bafkreihtq33ntlmhgu2jknbpehuq3tjcqh6y5ruaekhao3yl6ew32nsqbq",
	"moks.receiving.receipt.v1":      "bafkreiar4u5irmomc2lkptxhrsdsk5podezhrji55247hhkh2ngzkcba7i",
	"moks.runs.approval.v1":          "bafkreidbpikyfswoqv6uja2liswlao5t2jjpilklmkqqbgrl3fhtkpvpba",
	"moks.runs.evidence.v1":          "bafkreiflfm2hhldipxdzmjzmnlufisf4tihzn7wjrdbdduasawqesvdnly",
	"moks.runs.run.v1":               "bafkreia75oenisiarp2psot5tumzsw3snu5ofqzk5lzijirgs6r5jidox4",
	"moks.training.completion.v1":    "bafkreigm66qcalyb3bedzjmnybbqwvn6cmpcjw32lsjertk6i6kdxolpcm",
	"moks.training.item.v1":          "bafkreihsft37k2fefo7h3viuhzgmtdxy3k64lgypgnsvd63fbxtfyccq7y",
	"moks.training.session.v1":       "bafkreic5pd5ijoxdwdqnbj6yiy2ufkm4sl6p4jqnfftyqezvcfo5updnee",
}

// PackageProtocolPCID returns the fixed CIDv1 for a built-in package family.
// An empty result means the family is external to this binary's built-in
// registry and must supply its separately frozen pCID. Source: DI-solan.
func PackageProtocolPCID(family string) string {
	return packageFamilyPCIDs[family]
}
