import { copyFile, mkdir, rm } from "node:fs/promises";
import { spawn } from "node:child_process";

const repoRoot = process.cwd();
const harnessRoot = "/tmp/ex5-browser-interaction";
const runtimeRoot = `${harnessRoot}/runtime`;
const chromeRoot = `${harnessRoot}/chrome-profile`;
const listenAddress = "127.0.0.1:7046";
const appURL = `http://${listenAddress}`;
const devToolsURL = "http://127.0.0.1:9225";

function run(command, args) {
  return new Promise((resolve, reject) => {
    const child = spawn(command, args, { cwd: repoRoot, stdio: "inherit" });
    child.once("error", reject);
    child.once("exit", (code) => code === 0 ? resolve() : reject(new Error(`${command} exited ${code}`)));
  });
}

async function waitFor(url) {
  for (let attempt = 0; attempt < 100; attempt++) {
    try {
      if ((await fetch(url)).ok) return;
    } catch {
      // The process is still starting.
    }
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
  throw new Error(`timed out waiting for ${url}`);
}

function stop(child) {
  return new Promise((resolve) => {
    if (!child || child.exitCode !== null) {
      resolve();
      return;
    }
    child.once("exit", resolve);
    child.kill("SIGTERM");
    setTimeout(() => child.kill("SIGKILL"), 2000).unref();
  });
}

async function evaluate(page, expression) {
  const socket = new WebSocket(page.webSocketDebuggerUrl);
  let nextID = 0;
  const call = (method, params) => new Promise((resolve, reject) => {
    const id = ++nextID;
    const receive = (event) => {
      const message = JSON.parse(event.data);
      if (message.id !== id) return;
      socket.removeEventListener("message", receive);
      message.error ? reject(new Error(message.error.message)) : resolve(message.result);
    };
    socket.addEventListener("message", receive);
    socket.send(JSON.stringify({ id, method, params }));
  });
  await new Promise((resolve, reject) => {
    socket.addEventListener("open", resolve, { once: true });
    socket.addEventListener("error", reject, { once: true });
  });
  try {
    const result = await call("Runtime.evaluate", { expression, awaitPromise: true, returnByValue: true });
    if (result.exceptionDetails) throw new Error(result.exceptionDetails.text);
    return result.result.value;
  } finally {
    socket.close();
  }
}

let runtime;
let installer;
try {
  // Intent: Exercise the shipped Chrome extension and its native host against
  // a disposable runtime rather than relying on an already-open browser or
  // the developer's live demo. Source: DI-sabek; DI-bulaf; DI-temur
  await rm(harnessRoot, { recursive: true, force: true });
  await run("bash", ["scripts/setup-demo-browser.sh"]);
  await run("bash", ["scripts/load-sample-data.sh", runtimeRoot]);
  await mkdir(`${chromeRoot}/NativeMessagingHosts`, { recursive: true });
  await copyFile(
    "/tmp/ex5-demo-browser/native-host/operational_browser_host.json",
    `${chromeRoot}/NativeMessagingHosts/operational_browser_host.json`,
  );

  runtime = spawn("go", ["run", "./cmd/operational-knowledge", "-listen", listenAddress, "-data-root", runtimeRoot], {
    cwd: repoRoot,
    stdio: "inherit",
  });
  await waitFor(`${appURL}/api/meta`);

  installer = spawn("node", [
    "scripts/install-demo-chrome-extension.mjs",
    chromeRoot,
    `${repoRoot}/chrome-extension`,
    "9225",
    "about:blank",
  ], { cwd: repoRoot, stdio: ["ignore", "pipe", "inherit"] });
  await new Promise((resolve, reject) => {
    const timer = setTimeout(() => reject(new Error("timed out loading the Ex5 unpacked extension")), 15000);
    installer.stdout.on("data", (chunk) => {
      const output = chunk.toString();
      process.stdout.write(output);
      if (output.includes("Ex5 Chrome extension ready: miagfmaampfgjkojhccdilogehbjijpe")) {
        clearTimeout(timer);
        resolve();
      }
    });
    installer.once("error", (error) => { clearTimeout(timer); reject(error); });
    installer.once("exit", (code) => { clearTimeout(timer); reject(new Error(`extension installer exited ${code}`)); });
  });
  await waitFor(`${devToolsURL}/json/list`);
  let targets;
  for (let attempt = 0; attempt < 50; attempt++) {
    targets = await (await fetch(`${devToolsURL}/json/list`)).json();
    if (targets.some((target) => target.type === "service_worker" && target.url.includes("chrome-extension://miagfmaampfgjkojhccdilogehbjijpe/"))) break;
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
  console.log(`Chrome targets: ${targets.map((target) => `${target.type}:${target.url}`).join(" | ")}`);
  if (!targets.some((target) => target.type === "service_worker" && target.url.includes("chrome-extension://miagfmaampfgjkojhccdilogehbjijpe/"))) {
    throw new Error("Ex5 extension worker did not start");
  }
  const blankPage = targets.find((target) => target.type === "page" && target.url === "about:blank");
  if (!blankPage) throw new Error("Chrome did not expose its initial page through DevTools");
  await evaluate(blankPage, `location.href = ${JSON.stringify(appURL)}`);
  for (let attempt = 0; attempt < 50; attempt++) {
    targets = await (await fetch(`${devToolsURL}/json/list`)).json();
    if (targets.some((target) => target.type === "page" && target.url.startsWith(appURL))) break;
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
  const page = targets.find((target) => target.type === "page" && target.url.startsWith(appURL));
  if (!page) throw new Error("Ex5 page was not exposed through Chrome DevTools");

  const result = await evaluate(page, `(async () => {
    const visible = (id) => !document.getElementById(id).hidden;
    const click = (id) => document.getElementById(id).click();
    const settle = () => new Promise((resolve) => setTimeout(resolve, 500));
    for (let attempt = 0; attempt < 20 && !document.querySelector("#draft-review button"); attempt++) await settle();
    if (!visible("review-drafts-lane")) throw new Error("Draft queue is not visible on load");
    if (!document.querySelector("#draft-review button")) {
      const status = document.getElementById("workspace-status")?.textContent.trim() || "no workspace status";
      throw new Error("Chrome did not complete the extension/native-host/runtime handshake: " + status);
    }
    click("review-lane-hotspots"); await settle();
    if (!visible("review-hotspots-lane")) throw new Error("Problem hotspots did not become visible");
    if (!document.querySelector("#problem-review button")) throw new Error("Problem hotspots has no direct-contract results");
    click("review-lane-search"); await settle();
    if (!visible("review-search-lane")) throw new Error("Known record search did not become visible");
    document.getElementById("search-preset-drafts").click(); await settle();
    const inspect = [...document.querySelectorAll("#search-results button, #review-search-lane button")]
      .find((button) => button.textContent.trim() === "Inspect");
    if (!inspect) throw new Error("Known record search exposed no Inspect action");
    inspect.click(); await settle();
    const current = document.getElementById("detail-meta").textContent.trim();
    if (!current || current.startsWith("Open a record")) throw new Error("Current Record did not visibly update");
    return "PASS";
  })()`);
  if (result !== "PASS") throw new Error("Chrome/native-host interaction assertions did not pass");
  console.log("Ex5 Chrome/native-host interaction workflow: PASS");
} finally {
  await stop(installer);
  await stop(runtime);
  await rm(harnessRoot, { recursive: true, force: true });
}
