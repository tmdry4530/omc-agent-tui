package bridge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chamdom/omc-agent-tui/pkg/schema"
)

func TestMapAgentTypeToRole(t *testing.T) {
	tests := []struct {
		agentType string
		want      schema.Role
	}{
		{"oh-my-claudecode:executor", schema.RoleExecutor},
		{"oh-my-claudecode:architect", schema.RoleArchitect},
		{"oh-my-claudecode:analyst", schema.RolePlanner},
		{"oh-my-claudecode:security-reviewer", schema.RoleGuard},
		{"oh-my-claudecode:test-engineer", schema.RoleTester},
		{"oh-my-claudecode:writer", schema.RoleWriter},
		{"oh-my-claudecode:unknown-agent", schema.RoleCustom},
		{"executor", schema.RoleExecutor},
	}
	for _, tt := range tests {
		t.Run(tt.agentType, func(t *testing.T) {
			got := mapAgentTypeToRole(tt.agentType)
			if got != tt.want {
				t.Errorf("mapAgentTypeToRole(%q) = %q, want %q", tt.agentType, got, tt.want)
			}
		})
	}
}

func TestMapParentMode(t *testing.T) {
	tests := []struct {
		input string
		want  schema.Mode
	}{
		{"ultrawork", schema.ModeUltrawork},
		{"ralph", schema.ModeRalph},
		{"none", ""},
		{"", ""},
		{"invalid-mode", schema.ModeUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := mapParentMode(tt.input)
			if got != tt.want {
				t.Errorf("mapParentMode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMapStatus(t *testing.T) {
	tests := []struct {
		status    string
		wantState schema.AgentState
		wantType  schema.EventType
	}{
		{"completed", schema.StateDone, schema.TypeTaskDone},
		{"failed", schema.StateFailed, schema.TypeError},
		{"cancelled", schema.StateCancelled, schema.TypeStateChange},
		{"unknown", schema.StateDone, schema.TypeTaskDone},
	}
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			state, eventType := mapStatus(tt.status)
			if state != tt.wantState {
				t.Errorf("state = %q, want %q", state, tt.wantState)
			}
			if eventType != tt.wantType {
				t.Errorf("type = %q, want %q", eventType, tt.wantType)
			}
		})
	}
}

func TestConvertAgents_SpawnAndDone(t *testing.T) {
	agents := []TrackedAgent{
		{
			AgentID:     "abc123",
			AgentType:   "oh-my-claudecode:executor",
			StartedAt:   "2026-02-17T20:33:41.402Z",
			ParentMode:  "ultrawork",
			Status:      "completed",
			CompletedAt: "2026-02-17T20:41:04.584Z",
			DurationMs:  443182,
		},
	}

	events, err := ConvertAgents(agents)
	if err != nil {
		t.Fatalf("ConvertAgents error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events (spawn+done), got %d", len(events))
	}

	// Spawn event
	spawn := events[0]
	if spawn.Type != schema.TypeTaskSpawn {
		t.Errorf("spawn type = %q, want task_spawn", spawn.Type)
	}
	if spawn.State != schema.StateRunning {
		t.Errorf("spawn state = %q, want running", spawn.State)
	}
	if spawn.Role != schema.RoleExecutor {
		t.Errorf("spawn role = %q, want executor", spawn.Role)
	}
	if spawn.Mode != schema.ModeUltrawork {
		t.Errorf("spawn mode = %q, want ultrawork", spawn.Mode)
	}
	if spawn.AgentID != "abc123" {
		t.Errorf("spawn agent_id = %q, want abc123", spawn.AgentID)
	}
	if spawn.Provider != schema.ProviderClaude {
		t.Errorf("spawn provider = %q, want claude", spawn.Provider)
	}

	// Done event
	done := events[1]
	if done.Type != schema.TypeTaskDone {
		t.Errorf("done type = %q, want task_done", done.Type)
	}
	if done.State != schema.StateDone {
		t.Errorf("done state = %q, want done", done.State)
	}
	if done.Payload == nil {
		t.Error("done payload should contain duration_ms")
	}

	// Validate events
	for i, evt := range events {
		if err := evt.Validate(); err != nil {
			t.Errorf("event[%d] validation failed: %v", i, err)
		}
	}
}

func TestConvertAgents_RunningAgent(t *testing.T) {
	agents := []TrackedAgent{
		{
			AgentID:    "running1",
			AgentType:  "oh-my-claudecode:planner",
			StartedAt:  "2026-02-17T20:00:00.000Z",
			ParentMode: "none",
			Status:     "running",
		},
	}

	events, err := ConvertAgents(agents)
	if err != nil {
		t.Fatalf("ConvertAgents error: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 event (spawn only), got %d", len(events))
	}
	if events[0].Type != schema.TypeTaskSpawn {
		t.Errorf("expected task_spawn, got %q", events[0].Type)
	}
}

func TestConvertAgents_FailedAgent(t *testing.T) {
	agents := []TrackedAgent{
		{
			AgentID:     "fail1",
			AgentType:   "oh-my-claudecode:executor",
			StartedAt:   "2026-02-17T20:00:00.000Z",
			ParentMode:  "ultrawork",
			Status:      "failed",
			CompletedAt: "2026-02-17T20:05:00.000Z",
			DurationMs:  300000,
		},
	}

	events, err := ConvertAgents(agents)
	if err != nil {
		t.Fatalf("ConvertAgents error: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	errEvt := events[1]
	if errEvt.State != schema.StateFailed {
		t.Errorf("expected failed state, got %q", errEvt.State)
	}
	if errEvt.Type != schema.TypeError {
		t.Errorf("expected error type, got %q", errEvt.Type)
	}
}

func TestConvertTracking_FileRoundtrip(t *testing.T) {
	tracking := TrackingFile{
		Agents: []TrackedAgent{
			{
				AgentID:     "a1",
				AgentType:   "oh-my-claudecode:executor",
				StartedAt:   "2026-02-17T10:00:00.000Z",
				ParentMode:  "ultrawork",
				Status:      "completed",
				CompletedAt: "2026-02-17T10:05:00.000Z",
				DurationMs:  300000,
			},
			{
				AgentID:     "a2",
				AgentType:   "oh-my-claudecode:architect",
				StartedAt:   "2026-02-17T10:01:00.000Z",
				ParentMode:  "none",
				Status:      "completed",
				CompletedAt: "2026-02-17T10:03:00.000Z",
				DurationMs:  120000,
			},
		},
		TotalSpawned:   2,
		TotalCompleted: 2,
	}

	dir := t.TempDir()
	trackingPath := filepath.Join(dir, "subagent-tracking.json")
	data, _ := json.Marshal(tracking)
	os.WriteFile(trackingPath, data, 0644)

	events, err := ConvertTracking(trackingPath)
	if err != nil {
		t.Fatalf("ConvertTracking error: %v", err)
	}

	// 2 agents * 2 events each = 4
	if len(events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(events))
	}

	// All should validate
	for i, evt := range events {
		if err := evt.Validate(); err != nil {
			t.Errorf("event[%d] invalid: %v", i, err)
		}
	}
}

func TestEmitEvent(t *testing.T) {
	dir := t.TempDir()
	eventDir := filepath.Join(dir, "events")

	evt := NewSpawnEvent("test-1", "oh-my-claudecode:executor", "ultrawork")
	if err := EmitEvent(eventDir, "session-001", evt); err != nil {
		t.Fatalf("EmitEvent error: %v", err)
	}

	// Verify file exists and contains valid JSONL
	path := filepath.Join(eventDir, "session-001.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read event file: %v", err)
	}

	var decoded schema.CanonicalEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal event: %v", err)
	}
	if decoded.AgentID != "test-1" {
		t.Errorf("agent_id = %q, want test-1", decoded.AgentID)
	}
	if decoded.Type != schema.TypeTaskSpawn {
		t.Errorf("type = %q, want task_spawn", decoded.Type)
	}
}

func TestEmitEvent_Append(t *testing.T) {
	dir := t.TempDir()
	eventDir := filepath.Join(dir, "events")

	// Emit 3 events
	for i, typ := range []string{"spawn", "update", "done"} {
		var evt schema.CanonicalEvent
		switch typ {
		case "spawn":
			evt = NewSpawnEvent("agent-1", "oh-my-claudecode:executor", "ultrawork")
		case "update":
			evt = NewUpdateEvent("agent-1", "oh-my-claudecode:executor", "ultrawork", schema.StateWaiting)
		case "done":
			evt = NewDoneEvent("agent-1", "oh-my-claudecode:executor", "ultrawork")
		}
		if err := EmitEvent(eventDir, "sess", evt); err != nil {
			t.Fatalf("emit[%d] error: %v", i, err)
		}
	}

	// Read and count lines
	data, _ := os.ReadFile(filepath.Join(eventDir, "sess.jsonl"))
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	// Verify each line is valid JSON with correct type
	expectedTypes := []schema.EventType{schema.TypeTaskSpawn, schema.TypeTaskUpdate, schema.TypeTaskDone}
	for i, line := range lines {
		var evt schema.CanonicalEvent
		if err := json.Unmarshal([]byte(line), &evt); err != nil {
			t.Fatalf("line[%d] unmarshal: %v", i, err)
		}
		if evt.Type != expectedTypes[i] {
			t.Errorf("line[%d] type = %q, want %q", i, evt.Type, expectedTypes[i])
		}
	}
}

func TestNewErrorEvent(t *testing.T) {
	evt := NewErrorEvent("err-1", "oh-my-claudecode:executor", "ralph", "something broke")
	if evt.State != schema.StateError {
		t.Errorf("state = %q, want error", evt.State)
	}
	if evt.Type != schema.TypeError {
		t.Errorf("type = %q, want error", evt.Type)
	}
	if evt.Mode != schema.ModeRalph {
		t.Errorf("mode = %q, want ralph", evt.Mode)
	}

	var payload map[string]string
	json.Unmarshal(evt.Payload, &payload)
	if payload["error"] != "something broke" {
		t.Errorf("payload error = %q, want 'something broke'", payload["error"])
	}
}

func TestWriteEventsFile(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "output", "converted.jsonl")

	events := []schema.CanonicalEvent{
		NewSpawnEvent("a1", "oh-my-claudecode:executor", "ultrawork"),
		NewDoneEvent("a1", "oh-my-claudecode:executor", "ultrawork"),
	}

	if err := WriteEventsFile(outPath, events); err != nil {
		t.Fatalf("WriteEventsFile error: %v", err)
	}

	data, _ := os.ReadFile(outPath)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestConvertAndWriteRoundtrip(t *testing.T) {
	dir := t.TempDir()

	// Create tracking file
	tracking := TrackingFile{
		Agents: []TrackedAgent{
			{
				AgentID:     "rt1",
				AgentType:   "oh-my-claudecode:executor",
				StartedAt:   "2026-02-17T10:00:00.000Z",
				ParentMode:  "ultrawork",
				Status:      "completed",
				CompletedAt: "2026-02-17T10:05:00.000Z",
				DurationMs:  300000,
			},
		},
	}
	trackingPath := filepath.Join(dir, "tracking.json")
	data, _ := json.Marshal(tracking)
	os.WriteFile(trackingPath, data, 0644)

	// Convert
	events, err := ConvertTracking(trackingPath)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}

	// Write
	outPath := filepath.Join(dir, "events.jsonl")
	if err := WriteEventsFile(outPath, events); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Read back and verify
	outData, _ := os.ReadFile(outPath)
	lines := strings.Split(strings.TrimSpace(string(outData)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}

	// Verify chronological order
	var first, second schema.CanonicalEvent
	json.Unmarshal([]byte(lines[0]), &first)
	json.Unmarshal([]byte(lines[1]), &second)
	if !first.Ts.Before(second.Ts) {
		t.Error("events should be in chronological order")
	}
}

func TestEventFactories_Validation(t *testing.T) {
	factories := []struct {
		name string
		evt  schema.CanonicalEvent
	}{
		{"spawn", NewSpawnEvent("v1", "oh-my-claudecode:executor", "ultrawork")},
		{"update", NewUpdateEvent("v1", "oh-my-claudecode:executor", "ultrawork", schema.StateWaiting)},
		{"done", NewDoneEvent("v1", "oh-my-claudecode:executor", "ultrawork")},
		{"error", NewErrorEvent("v1", "oh-my-claudecode:executor", "ultrawork", "fail")},
	}
	for _, tt := range factories {
		t.Run(tt.name, func(t *testing.T) {
			if tt.evt.Ts.IsZero() {
				t.Error("ts should not be zero")
			}
			if time.Since(tt.evt.Ts) > 5*time.Second {
				t.Error("ts should be recent")
			}
			if err := tt.evt.Validate(); err != nil {
				t.Errorf("validation failed: %v", err)
			}
		})
	}
}
