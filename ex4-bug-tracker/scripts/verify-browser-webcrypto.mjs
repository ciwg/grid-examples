import { mkdtemp, rm } from "node:fs/promises";
import { spawn } from "node:child_process";
import { tmpdir } from "node:os";
import { join } from "node:path";

const runtimeRoot = await mkdtemp(join(tmpdir(), "ex4-browser-e2e-"));
const chromeRoot = await mkdtemp(join(tmpdir(), "ex4-browser-chrome-"));
const server = spawn("go", ["run", "./cmd/bug-tracker", "--listen", "127.0.0.1:7044", "--data-root", runtimeRoot], { env: { ...process.env, GOCACHE: "/tmp/ex4-go-cache" } });
const chrome = spawn("google-chrome", ["--headless=new", "--no-sandbox", "--disable-gpu", "--remote-debugging-port=9224", `--user-data-dir=${chromeRoot}`, "http://127.0.0.1:7044/"]);
const stopProcess = (process) => new Promise((resolve) => {
  if (process.exitCode !== null) return resolve();
  process.once("exit", resolve);
  process.kill();
});
const stop = async () => { await stopProcess(server); await stopProcess(chrome); await rm(runtimeRoot, { recursive: true, force: true }); await rm(chromeRoot, { recursive: true, force: true }); };

try {
  for (let i = 0; i < 40; i++) { try { await fetch("http://127.0.0.1:7044/api/meta"); break; } catch { await new Promise((resolve) => setTimeout(resolve, 100)); } }
  await new Promise((resolve) => setTimeout(resolve, 700));
  const targets = await (await fetch("http://127.0.0.1:9224/json")).json();
  const target = targets.find((item) => item.type === "page");
  const socket = new WebSocket(target.webSocketDebuggerUrl);
  let nextID = 0;
  const command = (method, params) => new Promise((resolve, reject) => { const id = ++nextID; socket.send(JSON.stringify({ id, method, params })); const receive = (event) => { const message = JSON.parse(event.data); if (message.id === id) { socket.removeEventListener("message", receive); message.error ? reject(message.error) : resolve(message.result); } }; socket.addEventListener("message", receive); });
  await new Promise((resolve) => socket.addEventListener("open", resolve, { once: true }));
  const expression = `(async()=>{document.querySelector('#new-issue-button').click();document.querySelector('#issue-title').value='WebCrypto regression';document.querySelector('#issue-description').value='headless browser';document.querySelector('#new-issue-form').dispatchEvent(new Event('submit',{bubbles:true,cancelable:true}));await new Promise(r=>setTimeout(r,1500));return await (await fetch('/api/issues')).text()})()`;
  const result = await command("Runtime.evaluate", { expression, awaitPromise: true, returnByValue: true });
  socket.close();
  if (!result.result.value.includes("WebCrypto regression")) throw new Error("browser did not create a signed issue");
  console.log("browser WebCrypto workflow: PASS");
} finally { await stop(); }
