const proofRoot = process.env.EX7_PROOF_ROOT;
if (!proofRoot) throw new Error('EX7_PROOF_ROOT is required');
const pause = ms => new Promise(resolve => setTimeout(resolve, ms));
const newTab = async url => await (await fetch(`http://127.0.0.1:9229/json/new?${encodeURIComponent(url)}`, {method: 'PUT'})).json();
function cdp(url) { const ws = new WebSocket(url); let next = 0; const pending = new Map(); ws.onmessage = event => { const message = JSON.parse(event.data); if (message.id) { const request = pending.get(message.id); pending.delete(message.id); message.error ? request.reject(message.error) : request.resolve(message.result); } }; return new Promise(resolve => ws.onopen = () => resolve({call: (method, params = {}) => new Promise((resolve, reject) => { const id = ++next; pending.set(id, {resolve, reject}); ws.send(JSON.stringify({id, method, params})); })})); }
const bobTab = await newTab('about:blank'); const aliceTab = await newTab('about:blank');
const bob = await cdp(bobTab.webSocketDebuggerUrl); const alice = await cdp(aliceTab.webSocketDebuggerUrl);
const consoleEntries = [];
alice.call('Runtime.enable');
alice.call('Log.enable');
await bob.call('Page.enable'); await alice.call('Page.enable'); await bob.call('Page.navigate', {url: 'http://127.0.0.1:7038/'}); await alice.call('Page.navigate', {url: 'http://127.0.0.1:7037/'});
const evaluate = (target, expression) => target.call('Runtime.evaluate', {expression, awaitPromise: true, returnByValue: true});
for (let index = 0; index < 20; index++) { if ((await evaluate(bob, `Boolean(document.querySelector('#approval-target'))`)).result.value) break; if (index === 19) throw new Error('Bob terminal form did not load'); await pause(250); }
await evaluate(bob, `document.querySelector('#approval-target').value='http://127.0.0.1:7037';document.querySelector('#approval-observation').value='Chrome two-agent proof';document.querySelector('#approval-request-form').requestSubmit();true`);
await pause(1000); const message = await evaluate(bob, `document.querySelector('#message').textContent`); if (message.result.value) throw new Error(message.result.value);
await alice.call('Page.reload', {ignoreCache: true});
for (let index = 0; index < 20; index++) { if ((await evaluate(alice, `Boolean(document.querySelector('[data-approve]'))`)).result.value) break; if (index === 19) { const fs = await import('node:fs/promises'); const body = await evaluate(alice, `document.body.innerText`); const incoming = await evaluate(alice, `document.querySelector('#incoming-approvals')?.innerHTML || ''`); await fs.writeFile(`${proofRoot}/alice-body.txt`, body.result.value || ''); await fs.writeFile(`${proofRoot}/alice-incoming.html`, incoming.result.value || ''); await fs.writeFile(`${proofRoot}/alice-console.json`, JSON.stringify(consoleEntries)); throw new Error('Alice approval control did not render'); } await pause(250); }
const approval = await evaluate(alice, `fetch('/api/approval-requests/' + document.querySelector('[data-approve]').dataset.approve + '/approve',{method:'POST'}).then(async response => ({status:response.status,body:await response.text()}))`);
await (await import('node:fs/promises')).writeFile(`${proofRoot}/alice-approval-response.json`, JSON.stringify(approval.result.value || approval));
for (let index = 0; index < 20; index++) { await pause(750); const status = await evaluate(bob, `document.querySelector('#approval-status').textContent`); if (status.result.value.includes('independently retained')) break; if (index === 19) throw new Error(status.result.value); }
const screenshot = await bob.call('Page.captureScreenshot', {format: 'png'}); await (await import('node:fs/promises')).writeFile(`${proofRoot}/bob-final.png`, Buffer.from(screenshot.data, 'base64'));
