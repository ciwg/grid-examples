#!/usr/bin/env node

import { spawn } from "node:child_process";
import { resolve } from "node:path";

const [profileRoot, extensionRoot, debugPort, pageURL] = process.argv.slice(2);
const expectedExtensionID = "miagfmaampfgjkojhccdilogehbjijpe";

if (!profileRoot || !extensionRoot || !debugPort || !pageURL) {
  console.error("usage: install-demo-chrome-extension.mjs <profile-root> <extension-root> <debug-port> <page-url>");
  process.exit(2);
}

// Intent: Keep Chrome's sanctioned DevTools-pipe owner alive for the one
// disposable demo session, while the interaction suite attaches separately
// through the localhost DevTools port. Source: DI-bilim.
const chrome = spawn("google-chrome", [
  `--user-data-dir=${resolve(profileRoot)}`,
  "--enable-unsafe-extension-debugging",
  "--remote-debugging-pipe",
  "--remote-debugging-address=127.0.0.1",
  `--remote-debugging-port=${debugPort}`,
  "--no-first-run",
  "--no-default-browser-check",
  pageURL,
], { stdio: ["ignore", "inherit", "inherit", "pipe", "pipe"] });

let nextID = 0;
let settled = false;
let stopping = false;
let receiveBuffer = Buffer.alloc(0);
const pending = new Map();

function finish(code) {
  if (settled) return;
  settled = true;
  clearTimeout(timeout);
  process.exitCode = code;
}

function send(method, params = {}) {
  return new Promise((resolvePromise, rejectPromise) => {
    const id = ++nextID;
    // Chrome's remote-debugging pipe frames each CDP JSON message with NUL.
    const frame = Buffer.from(`${JSON.stringify({ id, method, params })}\0`);
    pending.set(id, { resolve: resolvePromise, reject: rejectPromise });
    chrome.stdio[3].write(frame);
  });
}

function rejectPending(error) {
  for (const { reject } of pending.values()) reject(error);
  pending.clear();
}

chrome.stdio[4].on("data", (chunk) => {
  receiveBuffer = Buffer.concat([receiveBuffer, chunk]);
  while (true) {
    const frameEnd = receiveBuffer.indexOf(0);
    if (frameEnd < 0) return;
    const message = JSON.parse(receiveBuffer.subarray(0, frameEnd).toString("utf8"));
    receiveBuffer = receiveBuffer.subarray(frameEnd + 1);
    const request = pending.get(message.id);
    if (!request) continue;
    pending.delete(message.id);
    if (message.error) request.reject(new Error(message.error.message));
    else request.resolve(message.result);
  }
});

chrome.once("error", (error) => {
  console.error(`could not start Google Chrome: ${error.message}`);
  rejectPending(error);
  finish(1);
});

chrome.once("exit", (code, signal) => {
  if (!settled && !stopping && (code !== 0 || signal !== null)) {
    const error = new Error(`Chrome installer exited (${signal || code})`);
    rejectPending(error);
    console.error(error.message);
    finish(1);
  }
});

const timeout = setTimeout(() => {
  const error = new Error("Chrome DevTools-pipe extension installation timed out");
  rejectPending(error);
  console.error(error.message);
  chrome.kill("SIGTERM");
  finish(1);
}, 15000);

for (const signal of ["SIGINT", "SIGTERM"]) {
  process.on(signal, () => {
    stopping = true;
    chrome.kill(signal);
  });
}

try {
  const result = await send("Extensions.loadUnpacked", { path: resolve(extensionRoot) });
  if (result.id !== expectedExtensionID) {
    throw new Error(`Chrome installed unexpected extension ID: ${result.id}`);
  }
  clearTimeout(timeout);
  console.log(`Ex5 Chrome extension ready: ${result.id}`);
} catch (error) {
  console.error(`could not install Ex5 Chrome extension: ${error.message}`);
  chrome.kill("SIGTERM");
  finish(1);
}
