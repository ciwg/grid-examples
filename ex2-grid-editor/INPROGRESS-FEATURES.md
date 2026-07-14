# INPROGRESS FEATURES

Current review status for features from the older collaborative editor.

## Suggested phases

### Phase 1 - Core collaborative editor

Status: implemented slice complete

- 1, 2, 3, 4
- 6, 7, 8
- 17, 18, 19, 20, 22, 24
- 33, 35, 36, 37, 38, 39, 40
- 41, 42
- 54, 55, 60, 61, 62, 63, 64, 65, 67, 68
- 73, 74, 76, 78

Goal:
- make the shared editor feel strong, usable, and demo-ready
- finish, fix, and polish features that are already partially implemented

### Phase 2 - Document workflow and export

Status: implemented slice complete

- 9, 10, 11, 12, 13, 15, 16
- 34, 43, 44
- 45, 46, 47, 48
- 51, 52, 53
- 79, 80
- 111, 112, 113, 114, 115
- 131, 132, 133, 134, 135, 136, 137, 138
- 154, 155, 156, 157, 158, 159, 160
- 162, 163, 164, 168

Goal:
- make documents easier to create, navigate, share, export, and demo

### Phase 3 - Review, comments, and history

Status: implemented slice complete

- 81, 82, 83, 84, 85, 86, 87, 88, 90
- 102, 103, 104, 105, 106, 107, 108, 109, 110
- 116, 117, 118, 119, 120
- 121, 122, 123, 124, 125, 126, 127, 128, 129, 130
- 171, 176, 177, 180

Goal:
- add richer collaboration, review flow, visibility, and content tools

### Phase 4 - PromiseGrid-native backend features

Status: implemented slice complete

- 89
- 91, 92, 93, 94, 95, 96, 97, 98, 99, 100
- 28, 29, 30

Goal:
- handle version restore, permissions, ownership, publishing, and document exchange in PromiseGrid-native ways

### Phase 5 - Later / optional / boss review

Status: approved

- 5. Activity log
- 14. Rename document
- 21. Word count / document stats
- 23. Format document / whitespace cleanup
- 49. PromiseGrid CBOR export
- 50. Raw JSON / debug export
- 58. Shareable read-only view
- 59. Document lock / read-only mode
- 69. Browser installable / PWA behavior
- 71. Read-only spectator mode
- 72. Presenter / demo mode
- 75. Bigger cursor / higher-contrast cursor mode
- 77. Mobile-friendly layout
- 101. Presence status message
- 139. Sticky heading while scrolling
- 140. Minimap
- 141. Track changes mode
- 142. Accept / reject changes
- 143. Suggestion mode
- 144. Comment resolution
- 145. Resolved comments view
- 146. Mention notifications
- 147. Follow another user’s cursor
- 148. Bring me to where another user is editing
- 149. Presence sound / alert
- 150. Quiet / do-not-disturb mode
- 151. Session recording / replay
- 152. Demo script mode
- 153. Guided tour overlay
- 165. Deep link to heading
- 166. Deep link to comment
- 167. Deep link to selection
- 169. Personal notes on a shared doc
- 170. Private highlights on a shared doc
- 172. AI writing assistance
- 173. AI rewrite / tone change
- 174. AI explain selected text
- 175. AI action on selected text
- 178. Live translation
- 179. Multi-language UI
- 181. Plugin system
- 182. Custom extensions
- 183. Scripting hooks
- 184. Automation rules
- 185. Keyboard macro recording
- 186. Reusable command palette
- 187. Custom toolbar actions
- 188. Workspace settings per document
- 189. Workspace settings per user
- 190. Import / export settings
- 191. Better document ID naming
- 192. Human-friendly share links
- 193. Custom document slug
- 194. Alias for a document
- 195. Permanent canonical document ID
- 196. Redirect old link to new link
- 197. Merge two documents
- 198. Split one document into two
- 199. Import markdown into new doc
- 200. Import existing file into shared doc

Goal:
- reserve these for later product decisions, demo refinements, or boss review

## Confirmed

1. Remote cursors
2. Selection highlights
3. Typing indicators
4. User list with count
6. New shared document flow
7. Open existing by shared link
8. Automatic reconnection
9. Markdown preview
10. Export formats
11. Offline persistence
12. Document copy / duplicate
13. Open recent documents
15. Copy document URL
16. Share / email document link
17. Keyboard shortcuts
18. Preferences dialog
19. Line numbers toggle
20. Dark / light mode
22. Find / search
24. Bold / italic / underline / markdown formatting actions
25. Print
26. Version history
27. Compare versions / diff view
32. Presence / activity timestamps
33. Browser and Neovim interoperability
34. Document title / metadata
35. Browser menu system
36. Status indicator: online / offline / reconnecting
37. Custom user name
38. Custom user color
39. Color picker
40. Connection / session info panel
41. Remote selections in Neovim
42. Cursor color labels with names
43. Multi-document tabs or switching
44. Document registry / document list
45. Save / export as Automerge file
46. Save / export as HTML
47. Save / export as Markdown
48. Save / export as plain text
51. Document creation timestamp
52. Last edited timestamp
53. Last viewed timestamp
54. Join / leave notifications
55. Peer count in toolbar
56. Scroll sync in markdown preview
57. Markdown preview toggle shortcut
60. Presence aging / stale-offline-removal behavior
61. Reconnect banner after disconnect
62. Unsynced local changes indicator
63. Better error messages for failed sync
64. Document loaded / syncing / ready states
65. Peer color legend
66. Presence badges above editor
  Note: already implemented in the web UI; not required in the same location, and not currently needed in Neovim
67. Open doc by paste-in link
68. Simple welcome / onboarding flow
70. Service worker / offline app shell
73. Font resize / larger text controls
74. Dyslexia-friendly spacing / font options
76. Keyboard shortcut help overlay
  Note: shortcut remapping is necessary because Mac, Linux, and Windows users may need different bindings and may already have local mappings
78. Better copy/paste handling
79. File import
  Note: drag-and-drop is optional; keep plain file import even if drag-and-drop is dropped
80. Image paste / attachment support
81. Multi-user comments
82. Inline annotations
83. Per-document chat
  Note: only as inline comments, not a separate chat pane
84. Document activity feed
85. Presence history / recent participants
86. “Last edited by” display
87. “Last viewed by” display
88. Named saved versions
89. Restore old version
  Note: needs special PromiseGrid-specific backend behavior and must not be treated as a simple UI rollback feature
  Note: implement as a new current-time restore action that creates a new state from an older version, instead of deleting history or replacing the past
  Note: double-check with boss before finalizing behavior
90. Audit trail / change timeline
91. Per-user permissions
92. Document owner / admin role
93. Invite link management
94. Temporary guest access
95. Document description / summary
96. Tags or labels
97. Search across documents
98. Folder or collection grouping
99. Pin favorite documents
100. Archive document
  Note: 91-100 need PromiseGrid-native design rather than a generic app-only permission or storage model
102. Custom profile picture / avatar
103. Emoji reactions
104. Inline task checklist support
105. Table editing support
106. Code block tools
107. Link preview
108. Embedded media preview
109. Mention people with `@name`
  Note: visible `@name` in the UI should resolve to a stable id underneath
110. Notification settings
111. Document templates
112. Starter templates for notes / docs / checklists
113. Slash commands like `/todo` or `/h1`
114. Quick insert menu
115. Auto-save indicator
116. Local draft recovery after crash
  Note: implement in PromiseGrid-native ways rather than as a browser-only local draft feature
117. Conflict / merge explanation UI
118. Connection diagnostics view
119. Sync history inspector
120. Developer debug panel
121. Multi-file document set
122. Linked documents
123. Backlinks
124. Cross-document references
125. Document graph view
126. Outline / heading navigator
127. Table of contents panel
128. Jump to heading
129. Collapsible sections
130. Focus mode
131. Split view editing
132. Side-by-side two documents
133. Open same document in two panes
134. Search and replace
135. Case-sensitive search option
136. Regex search option
137. Replace all
138. Go to line
154. Built-in help panel
155. Troubleshooting panel
  Note: can be part of the built-in help panel
156. Test document generator
157. Sample collaborative demo document
158. Template gallery
159. Publish document snapshot
160. Export audit report
162. Copy as markdown
163. Copy as HTML
164. Copy share link to current section
168. Bookmark a location in doc
171. Auto-generated summary
176. Voice dictation
177. Text-to-speech readback
180. Spellcheck / grammar help

## Maybe

5. Activity log
14. Rename document
21. Word count / document stats
23. Format document / whitespace cleanup
49. PromiseGrid CBOR export
50. Raw JSON / debug export
58. Shareable read-only view
59. Document lock / read-only mode
69. Browser installable / PWA behavior
71. Read-only spectator mode
72. Presenter / demo mode
75. Bigger cursor / higher-contrast cursor mode
77. Mobile-friendly layout
101. Presence status message
139. Sticky heading while scrolling
  Note: keep the current section heading visible at the top while scrolling through long documents
140. Minimap
  Note: a tiny side overview of the whole document for quick navigation
141. Track changes mode
  Note: show edits as tracked changes instead of directly blending them into the document view
142. Accept / reject changes
  Note: allow a reviewer to accept or reject specific tracked edits
143. Suggestion mode
  Note: edits are proposed as suggestions rather than applied as normal direct edits
144. Comment resolution
  Note: mark a comment thread as resolved without deleting its history
145. Resolved comments view
  Note: separate place to review comments that were already resolved
146. Mention notifications
  Note: notify a user when they are mentioned in a comment or annotation
147. Follow another user’s cursor
  Note: keep your view centered on another person’s current cursor location
148. Bring me to where another user is editing
  Note: quick jump to the place another collaborator is currently working
149. Presence sound / alert
  Note: optional sound or alert for join, leave, or mention events
150. Quiet / do-not-disturb mode
  Note: suppress non-essential alerts, sounds, or interruptions
151. Session recording / replay
  Note: record and replay a collaborative session
152. Demo script mode
  Note: step through a planned demo flow
153. Guided tour overlay
  Note: overlay that teaches the UI
165. Deep link to heading
  Note: link straight to a heading in the document
166. Deep link to comment
  Note: link straight to a comment thread
167. Deep link to selection
  Note: link straight to a selected text range
169. Personal notes on a shared doc
170. Private highlights on a shared doc
172. AI writing assistance
173. AI rewrite / tone change
174. AI explain selected text
175. AI action on selected text
178. Live translation
179. Multi-language UI
181. Plugin system
182. Custom extensions
183. Scripting hooks
184. Automation rules
185. Keyboard macro recording
186. Reusable command palette
187. Custom toolbar actions
  Note: related to shortcuts, but not the same; lets the user choose buttons or menu actions for common tasks
188. Workspace settings per document
189. Workspace settings per user
190. Import / export settings
  Note: save settings to a file and load them later on another machine or for backup
191. Better document ID naming
192. Human-friendly share links
193. Custom document slug
194. Alias for a document
195. Permanent canonical document ID
196. Redirect old link to new link
197. Merge two documents
198. Split one document into two
199. Import markdown into new doc
200. Import existing file into shared doc

## No

161. Export selected text only

## Change

28. GitHub pull file into editor -> PromiseGrid-native document exchange
29. GitHub commit from editor -> PromiseGrid-native publish / commit flow
30. AI commit message generation -> revisit under PromiseGrid-native publish flow
31. Typing status text -> use cursor behavior instead:
- flashing cursor = typing
- solid cursor = present

## Review note

These are in-progress feature decisions for review.

Planned flow:
- confirm feature set
- implement selected features
- review with boss
- keep, change, or drop after review

## Boss review items

These need explicit boss review by feature name, not just feature number.

- Restore old version
  Note: should be a new current-time restore action, not a rollback that rewrites history
- Per-user permissions
- Document owner / admin role
- Invite link management
- Temporary guest access
- Document description / summary
- Tags or labels
- Search across documents
- Folder or collection grouping
- Pin favorite documents
- Archive document
- PromiseGrid-native document exchange
- PromiseGrid-native publish / commit flow
- AI commit message generation under PromiseGrid-native publish flow
- Shareable read-only view
- Document lock / read-only mode
- Browser installable / PWA behavior
- Read-only spectator mode
- Presenter / demo mode
- Better document ID naming
- Human-friendly share links
- Custom document slug
- Alias for a document
- Permanent canonical document ID
- Redirect old link to new link
- Merge two documents
- Split one document into two
- Import markdown into new doc
- Import existing file into shared doc

## Phase completion note

Phases 1 through 4 have been implemented as working slices.

That does not mean every single feature in those groups is fully final.
It means:
- the main Phase 1 slice is in
- the main Phase 2 slice is in
- the main Phase 3 slice is in
- the main Phase 4 slices completed so far are in

Remaining review and product-shaping work is mostly in:
- the boss review items above
- the `Maybe` section
- open TODOs in `TODO/`
