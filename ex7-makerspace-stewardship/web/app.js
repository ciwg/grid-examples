let state;
const byId = (id) => document.getElementById(id);
const message = (text, error = false) => { const node = byId('message'); node.textContent = text; node.style.color = error ? '#a63b32' : '#166245'; };
const escape = (value) => String(value).replace(/[&<>'"]/g, character => ({'&':'&amp;','<':'&lt;','>':'&gt;',"'":'&#39;','"':'&quot;'}[character]));
const option = (value, label) => `<option value="${escape(value)}">${escape(label)}</option>`;
async function load() { state = await fetch('/api/state').then(response => response.json()); render(); }
function member(id) { return state.members.find(item => item.id === id); }
function area(id) { return state.areas.find(item => item.id === id).name; }
function render() {
  byId('tools').innerHTML = state.tools.map(tool => `<article class="tool ${tool.safetyHold ? 'hold' : ''}"><h3>${escape(tool.name)}</h3><p>${escape(area(tool.areaId))}</p><p class="status">${escape(tool.safetyHold ? 'Safety hold — unavailable' : tool.activeLoan ? `Loaned to ${member(tool.activeLoan.memberId).initials} until ${new Date(tool.activeLoan.dueAt).toLocaleString()}${tool.activeLoan.termsComplete ? '' : ' — accepted terms unavailable'}` : tool.condition)}</p><p>${tool.offSiteLoan ? 'Off-site loan eligible under area rules.' : 'In-space use only.'}</p><p>${tool.observations.length} recorded observation(s).</p></article>`).join('');
  byId('policies').innerHTML = state.areas.map(item => `<article class="tool"><h3>${escape(item.name)}</h3><p class="status">Current policy ${escape(item.policyVersion)}</p><p>${escape(item.policy)}</p><p><small>Authority delegated by ${escape(item.delegatedBy)}.</small></p></article>`).join('');
  byId('authorities').innerHTML = state.authorities.map(authority => `<article class="tool"><h3>${escape(member(authority.memberId).name)}</h3><p>${escape(area(authority.areaId))} steward</p><p>${authority.scopes.map(escape).join(' · ')}</p><p><small>Recognized by ${escape(authority.recognizedBy)}; review by ${escape(authority.reviewAt)}.</small></p></article>`).join('');
  const members = state.members.map(item => option(item.id, `${item.name} (${item.initials})`)).join('');
  byId('access-member').innerHTML = members;
  renderEligibility();
  const evidence = state.tools.flatMap(tool => tool.observations.map(observation => ({...observation, toolName: tool.name}))).sort((left, right) => new Date(right.createdAt) - new Date(left.createdAt));
  byId('timeline').innerHTML = evidence.length ? evidence.map(event => `<article class="event"><strong>${escape(event.toolName)}</strong> — ${escape(member(event.reporterId)?.name || event.reporterId)}<br>${escape(event.text)}<br><small>${new Date(event.createdAt).toLocaleString()}</small></article>`).join('') : '<p>No projected observations recorded yet.</p>';
}
function renderEligibility() {
  const memberID = byId('access-member').value;
  const viewer = member(memberID);
  const qualifications = state.qualifications.filter(item => item.memberId === memberID);
  byId('qualification-summary').textContent = qualifications.length ? `${viewer.name} has ${qualifications.length} currently accepted area qualification(s).` : `${viewer.name} has no currently accepted area qualifications.`;
  byId('eligibility').innerHTML = state.tools.map(tool => {
    const qualification = qualifications.find(item => item.areaId === tool.areaId && item.scope === tool.requiredQualification && item.status === 'accepted');
    const eligible = Boolean(qualification) && !tool.safetyHold;
    const required = tool.requiredQualification.replaceAll('-', ' ');
    const reason = tool.safetyHold ? 'This tool is on a safety hold pending area inspection.' : qualification ? `Eligible through the ${required} qualification issued by ${member(qualification.issuedBy).name}.` : `Not currently eligible: this requires an accepted ${required} qualification.`;
    return `<article class="tool ${eligible ? '' : 'hold'}"><h3>${escape(tool.name)}</h3><p class="status">${eligible ? 'Currently eligible' : 'Not currently eligible'}</p><p>${escape(reason)}</p><p><small>${tool.offSiteLoan ? 'Off-site lending is separately subject to current area terms.' : 'In-space use only.'}</small></p></article>`;
  }).join('');
}
byId('record-ingress-form').addEventListener('submit', async (event) => { event.preventDefault(); try { const records = byId('signed-records').value.split(/\s+/).filter(Boolean); const response = await fetch('/api/records', {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({records})}); const payload = await response.json(); if (!response.ok) throw new Error(payload.error); event.target.reset(); await load(); message('Exact signed record bytes retained; projection changed only when local recognition permits it.'); } catch (error) { message(error.message, true); } });
byId('access-member').addEventListener('change', renderEligibility);
load();
