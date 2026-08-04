let state;
const byId = (id) => document.getElementById(id);
const message = (text, error = false) => { const node = byId('message'); node.textContent = text; node.style.color = error ? '#a63b32' : '#166245'; };
const escape = (value) => String(value).replace(/[&<>'"]/g, character => ({'&':'&amp;','<':'&lt;','>':'&gt;',"'":'&#39;','"':'&quot;'}[character]));
const option = (value, label) => `<option value="${escape(value)}">${escape(label)}</option>`;
async function load() { state = await fetch('/api/state').then(response => response.json()); render(); }
function member(id) { return state.members.find(item => item.id === id); }
function area(id) { return state.areas.find(item => item.id === id).name; }
function render() {
  byId('tools').innerHTML = state.tools.map(tool => `<article class="tool ${tool.safetyHold ? 'hold' : ''}"><h3>${escape(tool.name)}</h3><p>${escape(area(tool.areaId))}</p><p class="status">${escape(tool.safetyHold ? 'Safety hold — unavailable' : tool.activeLoan ? `Loaned to ${member(tool.activeLoan.memberId).initials} until ${new Date(tool.activeLoan.dueAt).toLocaleString()}` : tool.condition)}</p><p>${tool.offSiteLoan ? 'Off-site loan eligible under area rules.' : 'In-space use only.'}</p><p>${tool.observations.length} recorded observation(s).</p></article>`).join('');
  const tools = state.tools.map(tool => option(tool.id, tool.name)).join('');
  byId('observation-tool').innerHTML = tools; byId('clear-tool').innerHTML = tools;
  const members = state.members.map(item => option(item.id, `${item.name} (${item.initials})`)).join('');
  byId('reporter').innerHTML = members; byId('borrower').innerHTML = members; byId('returner').innerHTML = members;
  byId('steward').innerHTML = state.authorities.map(authority => option(authority.memberId, `${member(authority.memberId).name} — ${area(authority.areaId)} steward`)).join('');
  byId('loan-tool').innerHTML = state.tools.filter(tool => tool.offSiteLoan && !tool.activeLoan).map(tool => option(tool.id, tool.name)).join('');
  byId('return-tool').innerHTML = state.tools.filter(tool => tool.activeLoan).map(tool => option(tool.id, tool.name)).join('');
  const evidence = state.tools.flatMap(tool => tool.observations.map(observation => ({...observation, toolName: tool.name}))).sort((left, right) => new Date(right.createdAt) - new Date(left.createdAt));
  byId('timeline').innerHTML = evidence.length ? evidence.map(event => `<article class="event"><strong>${escape(event.toolName)}</strong> — ${escape(member(event.reporterId).name)}<br>${escape(event.text)}<br>${(event.photos || []).map(photo => `<img class="photo" src="${escape(photo.dataUrl)}" alt="${escape(photo.name)}">`).join('')}<br><small>${new Date(event.createdAt).toLocaleString()}</small></article>`).join('') : '<p>No observations recorded yet.</p>';
}
async function post(path, body) { const response = await fetch(path, {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body)}); const payload = await response.json(); if (!response.ok) throw new Error(payload.error); await load(); return payload; }
async function readPhotos(input) { return Promise.all([...input.files].map(file => new Promise((resolve, reject) => { const reader = new FileReader(); reader.onload = () => resolve({name: file.name, dataUrl: reader.result}); reader.onerror = () => reject(new Error(`Could not read ${file.name}`)); reader.readAsDataURL(file); }))); }
function submit(formID, action) { byId(formID).addEventListener('submit', async (event) => { event.preventDefault(); try { await action(); event.target.reset(); } catch (error) { message(error.message, true); } }); }
submit('observation-form', async () => { await post(`/api/tools/${byId('observation-tool').value}/observations`, {reporterId: byId('reporter').value, text: byId('observation-text').value, safetyHold: byId('safety-hold').checked, photos: await readPhotos(byId('observation-photos'))}); message('Observation recorded.'); });
submit('clear-form', async () => { await post(`/api/tools/${byId('clear-tool').value}/clear-safety-hold`, {stewardId: byId('steward').value, assessment: byId('assessment').value}); message('Safety hold cleared with a recorded inspection.'); });
submit('loan-form', async () => { await post(`/api/tools/${byId('loan-tool').value}/loans`, {memberId: byId('borrower').value, dueAt: new Date(byId('due-at').value).toISOString()}); message('Off-site loan recorded with its accepted terms.'); });
submit('return-form', async () => { await post(`/api/tools/${byId('return-tool').value}/returns`, {memberId: byId('returner').value, condition: byId('return-condition').value}); message('Return condition recorded.'); });
load();
