import { build } from "esbuild";
import { copyFileSync, rmSync } from "node:fs";

rmSync("helper.bundle.mjs", { force: true });

await build({
  entryPoints: ["src/helper.mjs"],
  bundle: true,
  platform: "node",
  format: "cjs",
  target: "node20",
  outfile: "helper.bundle.cjs",
});

copyFileSync("node_modules/@automerge/automerge/dist/cjs/automerge_wasm_bg.wasm", "automerge_wasm_bg.wasm");
