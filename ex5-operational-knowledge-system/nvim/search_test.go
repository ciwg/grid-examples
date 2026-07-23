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

func TestNeovimSearchRendersGroupedBrowseResults(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/search" {
			http.NotFound(response, request)
			return
		}
		if request.URL.Query().Get("q") != "receiving audit" {
			http.Error(response, fmt.Sprintf("unexpected query %q", request.URL.RawQuery), http.StatusBadRequest)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprint(response, `{
			"filters":{"q":"receiving audit"},
			"places":[{"id":"PLACE-0001","kind":"area","name":"Receiving","summary":"Inbound receiving area"}],
			"resources":[{"id":"RES-0001","kind":"container","name":"RJ45 Bin","summary":"Connector bin","place_id":"PLACE-0001"}],
			"responsibilities":[{"id":"RESP-0001","title":"Receiving lead","summary":"Owns intake review"}],
			"items":[{"id":"ITEM-0001","kind":"receiving_check","status":"approved","title":"Inspect inbound pallet","summary":"Inbound pallet receiving check"}],
			"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"ITEM-0001","outcome":"accepted_with_notes","notes":"Outer wrap torn"}]
		}`); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "search.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
vim.env.OKS_SOCKET = "off"
local oks = require("oks")
oks.setup()

vim.cmd("OksSearch receiving audit")

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "## Places", 1, true) then
  error("missing places section")
end
if not string.find(body, "PLACE-0001 kind=area name=Receiving", 1, true) then
  error("missing place summary")
end
if not string.find(body, "inspect: :OksInspectEntity place PLACE-0001", 1, true) then
  error("missing place inspect hint")
end
if not string.find(body, "inspect: :OksInspect ITEM-0001", 1, true) then
  error("missing item inspect hint")
end
if not string.find(body, "inspect: :OksInspectRun RUN-0001", 1, true) then
  error("missing run inspect hint")
end
if not string.find(vim.api.nvim_buf_get_name(0), "oks-search://receiving%%20audit", 1, true) then
  error("unexpected search buffer name: " .. vim.api.nvim_buf_get_name(0))
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
		t.Fatalf("nvim search regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing ") || strings.Contains(string(output), "unexpected search buffer name") {
		t.Fatalf("unexpected search output: %s", string(output))
	}
}

func TestNeovimSearchAddsStructuredFilters(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/search" {
			http.NotFound(response, request)
			return
		}
		values := request.URL.Query()
		if values.Get("q") != "receiving audit" {
			http.Error(response, fmt.Sprintf("unexpected q %q", values.Get("q")), http.StatusBadRequest)
			return
		}
		if values.Get("kind") != "receiving_check" || values.Get("problem") != "true" || values.Get("place_id") != "PLACE-0001" {
			http.Error(response, fmt.Sprintf("unexpected filters %q", request.URL.RawQuery), http.StatusBadRequest)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprint(response, `{
			"filters":{"q":"receiving audit","kind":"receiving_check","problem":true,"place_id":"PLACE-0001"},
			"places":[],
			"resources":[],
			"responsibilities":[],
			"items":[{"id":"ITEM-0001","kind":"receiving_check","status":"approved","title":"Inspect inbound pallet","summary":"Inbound pallet receiving check"}],
			"runs":[{"id":"RUN-0001","kind":"receiving_check","item_id":"ITEM-0001","outcome":"accepted_with_notes","notes":"Outer wrap torn"}]
		}`); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer server.Close()

	script := filepath.Join(t.TempDir(), "search_filters.lua")
	scriptBody := fmt.Sprintf(`
vim.env.OKS_BASE_URL = %q
vim.env.OKS_SOCKET = "off"
local oks = require("oks")
oks.setup()

vim.cmd("OksSearch receiving audit kind=receiving_check problem=true place_id=PLACE-0001")

local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
local body = table.concat(lines, "\n")
if not string.find(body, "## Filters", 1, true) then
  error("missing filters section")
end
if not string.find(body, "- kind: receiving_check", 1, true) then
  error("missing kind filter")
end
if not string.find(body, "- problem: true", 1, true) then
  error("missing problem filter")
end
if not string.find(body, "- place_id: PLACE-0001", 1, true) then
  error("missing place filter")
end
if not string.find(vim.api.nvim_buf_get_name(0), "oks-search://receiving%%20audit", 1, true) then
  error("unexpected filtered search buffer name: " .. vim.api.nvim_buf_get_name(0))
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
		t.Fatalf("nvim filtered search regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing ") || strings.Contains(string(output), "unexpected filtered search buffer name") {
		t.Fatalf("unexpected filtered search output: %s", string(output))
	}
}

func TestNeovimSearchRejectsUnsupportedFilter(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	script := filepath.Join(t.TempDir(), "search_bad_filter.lua")
	scriptBody := `
local notices = {}
vim.notify = function(message, level)
  table.insert(notices, message)
end
local oks = require("oks")
oks.setup()

vim.cmd("OksSearch receiving audit owner=alice")

local joined = table.concat(notices, "\n")
if not string.find(joined, 'unsupported search filter "owner"', 1, true) then
  error("missing unsupported filter warning: " .. joined)
end
if string.find(vim.api.nvim_buf_get_name(0), "oks-search://", 1, true) then
  error("search buffer should not open for unsupported filter")
end
vim.cmd("qa!")
`
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
		t.Fatalf("nvim unsupported filter regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing unsupported filter warning") || strings.Contains(string(output), "search buffer should not open") {
		t.Fatalf("unexpected unsupported filter output: %s", string(output))
	}
}

func TestNeovimSearchRejectsTrailingWordsAfterFilters(t *testing.T) {
	nvimPath, err := exec.LookPath("nvim")
	if err != nil {
		t.Skip("nvim not available")
	}

	script := filepath.Join(t.TempDir(), "search_bad_order.lua")
	scriptBody := `
local notices = {}
vim.notify = function(message, level)
  table.insert(notices, message)
end
local oks = require("oks")
oks.setup()

vim.cmd("OksSearch receiving kind=procedure extra")

local joined = table.concat(notices, "\n")
if not string.find(joined, "free-text query words must come before key=value filters", 1, true) then
  error("missing ordering warning: " .. joined)
end
if string.find(vim.api.nvim_buf_get_name(0), "oks-search://", 1, true) then
  error("search buffer should not open for invalid query/filter ordering")
end
vim.cmd("qa!")
`
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
		t.Fatalf("nvim invalid search ordering regression: %v\n%s", err, string(output))
	}
	if strings.Contains(string(output), "missing ordering warning") || strings.Contains(string(output), "search buffer should not open") {
		t.Fatalf("unexpected invalid search ordering output: %s", string(output))
	}
}
