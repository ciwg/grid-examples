// Intent: Keep browser private signing keys non-extractable and local to the
// browser profile while the service owns canonical CBOR construction. Source: DI-zubot; DI-kolaf
const databaseName = "ex4-promise-keys";
const storeName = "keys";

function encodeBytes(bytes) {
  let text = "";
  for (const byte of bytes) text += String.fromCharCode(byte);
  return btoa(text);
}

function decodeBytes(value) {
  const text = atob(value);
  return Uint8Array.from(text, (character) => character.charCodeAt(0));
}

function requestJSON(path, body) {
  return fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  }).then(async (response) => {
    if (!response.ok) throw new Error((await response.text()) || `${path} failed`);
    return response.json();
  });
}

function openDatabase() {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(databaseName, 1);
    request.onupgradeneeded = () => request.result.createObjectStore(storeName);
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

async function loadKey(role) {
  const database = await openDatabase();
  return new Promise((resolve, reject) => {
    const request = database.transaction(storeName).objectStore(storeName).get(role);
    request.onsuccess = () => resolve(request.result || null);
    request.onerror = () => reject(request.error);
  });
}

async function saveKey(role, keyPair) {
  const database = await openDatabase();
  await new Promise((resolve, reject) => {
    const request = database.transaction(storeName, "readwrite").objectStore(storeName).put(keyPair, role);
    request.onsuccess = () => resolve();
    request.onerror = () => reject(request.error);
  });
}

async function keyForRole(role) {
  let keyPair = await loadKey(role);
  if (!keyPair) {
    keyPair = await crypto.subtle.generateKey({ name: "Ed25519" }, false, ["sign", "verify"]);
    await saveKey(role, keyPair);
  }
  return keyPair;
}

async function enrollmentForRole(role) {
  const keyPair = await keyForRole(role);
  const publicKey = new Uint8Array(await crypto.subtle.exportKey("raw", keyPair.publicKey));
  const prepared = await requestJSON("/api/agents/enroll/prepare", { public_key: encodeBytes(publicKey), role });
  const signature = new Uint8Array(await crypto.subtle.sign("Ed25519", keyPair.privateKey, decodeBytes(prepared.signable_bytes)));
  const enrollment = await requestJSON("/api/agents/enroll/finalize", { enrollment: prepared.enrollment, proof: { signature: encodeBytes(signature) } });
  return { keyPair, agentID: enrollment.agent_id, publicKey };
}

export async function submitPromise(role, profile, payload) {
  const signer = await enrollmentForRole(role);
  const prepared = await requestJSON("/api/promises/prepare", { profile, agent_id: signer.agentID, payload });
  const signature = new Uint8Array(await crypto.subtle.sign("Ed25519", signer.keyPair.privateKey, decodeBytes(prepared.signable_bytes)));
  const finalized = await requestJSON("/api/promises/finalize", { draft_id: prepared.draft_id, public_key: encodeBytes(signer.publicKey), signature: encodeBytes(signature) });
  const response = await fetch("/api/promises", { method: "POST", headers: { "Content-Type": "application/cbor" }, body: decodeBytes(finalized.envelope) });
  if (!response.ok) throw new Error((await response.text()) || "promise submission failed");
  return response.json();
}

export async function agentIDForRole(role) {
  return (await enrollmentForRole(role)).agentID;
}
