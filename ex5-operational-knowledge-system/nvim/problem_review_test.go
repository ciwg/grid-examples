package nvim

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNeovimProblemReviewRendersGroupedHotspots(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/problem-review" {
			http.NotFound(response, request)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprint(response, `{
			"problem_runs":2,
			"place_groups":[
				{
					"group_id":"PLACE-0001",
					"kind":"area",
					"name":"Receiving",
					"problem_count":2,
					"receiving_problems":1,
					"inventory_problems":1,
					"highlights":["supplier mismatch","count variance"],
					"runs":[
						{"id":"RUN-0001","kind":"receiving_check","item_id":"ITEM-0001","outcome":"accepted_with_notes","notes":"Outer wrap torn","resource_ids":["RES-0001"]}
					]
				}
			],
			"resource_groups":[
				{
					"group_id":"RES-0001",
					"kind":"container",
					"name":"RJ45 Bin",
					"problem_count":1,
					"receiving_problems":0,
					"inventory_problems":1,
					"highlights":["count variance"],
					"runs":[
						{"id":"RUN-0002","kind":"inventory_audit","item_id":"ITEM-0002","outcome":"completed","notes":"Counted receiving bin","resource_ids":["RES-0001"]}
					]
				}
			]
		}`); err != nil {
			t.Fatalf("write problem review response: %v", err)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "problem_review.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
local oks = require("oks")
oks.setup()

vim.cmd("OksProblemReview")

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "# Problem review", 1, true) then
  error("missing problem review header")
end
if not string.find(body, "## Place groups", 1, true) then
  error("missing place groups section")
end
if not string.find(body, "inspect: :OksInspectEntity place PLACE-0001", 1, true) then
  error("missing place inspect hint")
end
if not string.find(body, "inspect: :OksInspectRun RUN-0001", 1, true) then
  error("missing place run inspect hint")
end
if not string.find(body, "## Resource groups", 1, true) then
  error("missing resource groups section")
end
if not string.find(body, "inspect: :OksInspectEntity resource RES-0001", 1, true) then
  error("missing resource inspect hint")
end
if not string.find(vim.api.nvim_buf_get_name(0), "oks-problem-review://review", 1, true) then
  error("unexpected problem review buffer name: " .. vim.api.nvim_buf_get_name(0))
end
vim.cmd("qa!")
`, server.URL)
	if err := os.WriteFile(script, []byte(scriptBody), 0o644); err != nil {
		t.Fatalf("write script: %v", err)
	}

	command := exec.Command(
		nvimPath,
		"--headless",
		"-u", "NONE",
		"-c", "set runtimepath+=.",
		"-l", script,
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("nvim problem review regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing ") || strings.Contains(string(output), "unexpected problem review buffer name") {
		t.Fatalf("unexpected problem review output: %s", string(output))
	}
}

func TestNeovimProblemReviewRejectsMalformedGroups(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/problem-review" {
			http.NotFound(response, request)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprint(response, `{
			"problem_runs":1,
			"place_groups":[{"group_id":"PLACE-0001","kind":"area","name":"Receiving","problem_count":1,"receiving_problems":1,"inventory_problems":0}],
			"resource_groups":[]
		}`); err != nil {
			t.Fatalf("write malformed problem review response: %v", err)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "problem_review_malformed.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
local notices = {}
vim.notify = function(message, level)
  table.insert(notices, message)
end
local oks = require("oks")
oks.setup()

vim.cmd("OksProblemReview")

local joined = table.concat(notices, "\n")
if not string.find(joined, '/api/problem-review place_groups[1] missing "runs" array', 1, true) then
  error("missing problem group contract warning: " .. joined)
end
if string.find(vim.api.nvim_buf_get_name(0), "oks-problem-review://review", 1, true) then
  error("problem review buffer should not open on malformed payload")
end
vim.cmd("qa!")
`, server.URL)
	if err := os.WriteFile(script, []byte(scriptBody), 0o644); err != nil {
		t.Fatalf("write script: %v", err)
	}

	command := exec.Command(
		nvimPath,
		"--headless",
		"-u", "NONE",
		"-c", "set runtimepath+=.",
		"-l", script,
	)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("nvim malformed problem review regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing problem group contract warning") || strings.Contains(string(output), "problem review buffer should not open") {
		t.Fatalf("unexpected problem review output: %s", string(output))
	}
}
