# Event Schema (OMC Agent TUI)

## 1) 목적

OMC / Claude Code / Gemini / Codex에서 발생하는 이벤트를
하나의 공통 포맷(**CanonicalEvent**)으로 정규화하여
TUI가 일관된 방식으로 소비하도록 한다.

---

## 2) CanonicalEvent 정의

```json
{
  "ts": "2026-02-17T22:27:00Z",
  "run_id": "run-20260217-2227",
  "provider": "claude|gemini|codex|system",
  "mode": "ralph|ultrawork|ultrapilot|team|autopilot|unknown",
  "agent_id": "coder-auth",
  "parent_agent_id": "planner-main",
  "role": "planner|coder|reviewer|guard|tester|writer|custom",
  "state": "idle|running|waiting|blocked|error|done",
  "type": "task_spawn|task_update|task_done|tool_call|tool_result|message|error|replan|verify|fix",
  "task_id": "task-42",
  "intent_ref": "plan-7",
  "payload": {},
  "metrics": {
    "latency_ms": 420,
    "tokens_in": 210,
    "tokens_out": 95,
    "cost_usd": 0.0021
  },
  "raw_ref": "optional://source-pointer"
}
```

---

## 3) 필드 설명

* **ts (required)**
  이벤트 발생 시간 (ISO 8601, UTC 권장)

* **run_id (required)**
  실행 세션 식별자

* **provider (required)**
  이벤트 유래 모델/시스템

* **mode (optional)**
  OMC 실행 모드

* **agent_id (required)**
  이벤트를 발생시킨 에이전트

* **parent_agent_id (optional)**
  상위 에이전트 ID

* **role (required)**
  에이전트 역할

* **state (required)**
  에이전트 상태

* **type (required)**
  이벤트 타입

* **task_id (optional)**
  관련 태스크 ID

* **intent_ref (optional)**
  계획 단계 참조 ID

* **payload (optional)**
  원본 이벤트 세부 데이터

* **metrics (optional)**
  성능/비용 메타데이터

* **raw_ref (optional)**
  원본 로그 포인터 (디버깅용)

---

## 4) 상태 전이 규칙

* `idle → running` : 작업 시작
* `running → waiting` : 외부 결과 대기
* `waiting → running` : 결과 수신 후 재개
* `running → blocked` : 의존성 미충족
* `blocked → running` : 의존성 해소
* `running → error` : 실패
* `error → running` : 재시도
* `running → done` : 완료

---

## 5) 이벤트 타입 규약

* **task_spawn** : 자식 태스크 생성
* **task_update** : 진행/상태 업데이트
* **task_done** : 태스크 종료
* **tool_call** : 도구 호출 시작
* **tool_result** : 도구 호출 결과
* **message** : 에이전트 메시지
* **error** : 예외/오류
* **replan** : 계획 재수립
* **verify** : 검증 단계 진입/결과
* **fix** : 수정 단계 진입/결과

---

## 6) OMC 모드별 기대 패턴

### Ralph

* `verify → fix → verify` 반복 이벤트 빈도 높음
* `done` 전까지 loop 지속

### Ultrawork

* `task_spawn` fan-out 빈도 높음
* 병렬 agent 다수 `running` 상태

### Ultrapilot

* 단일 목표 중심의 긴 연속 실행
* 중간 `replan` + `recover` 이벤트 중요

---

## 7) Provider 매핑 가이드 (초안)

* Claude Code 내부 이벤트 → `provider = claude`
* Gemini MCP 호출 결과 → `provider = gemini`
* Codex MCP 호출 결과 → `provider = codex`
* 시스템 / 파서 / 스토어 이벤트 → `provider = system`

---

## 8) Redaction 규칙 (기본 ON)

### 마스킹 대상 키/패턴

**키명:**

* `api_key`
* `token`
* `secret`
* `password`
* `authorization`

**값 패턴:**

* 긴 base64/hex 토큰
* `sk-...` 형태 API 키
* `Bearer` 토큰

### 처리 규칙

* 원문 보존 금지 (기본)
* 출력은 `***REDACTED***`

---

## 9) 유효성 검증 규칙

### Required 필드

* `ts`
* `run_id`
* `provider`
* `agent_id`
* `role`
* `state`
* `type`

### Enum 검증

* `provider / mode / role / state / type`는 정의된 enum만 허용
* 미등록 값은 `unknown`으로 강등 + warning 기록

### 오류 처리

* 파싱 실패 이벤트는 drop하되 카운터 증가
* store / TUI는 파싱 실패로 종료되면 안 됨

---

## 10) 샘플 이벤트

### Sample: `task_spawn`

```json
{
  "ts": "2026-02-17T22:28:10Z",
  "run_id": "run-1",
  "provider": "claude",
  "mode": "ultrawork",
  "agent_id": "planner-main",
  "role": "planner",
  "state": "running",
  "type": "task_spawn",
  "task_id": "task-100",
  "payload": {
    "title": "Fix auth flow",
    "child_agent": "coder-auth"
  }
}
```

---

### Sample: `verify_fail_then_fix`

```json
{
  "ts": "2026-02-17T22:29:00Z",
  "run_id": "run-1",
  "provider": "claude",
  "mode": "ralph",
  "agent_id": "reviewer-1",
  "role": "reviewer",
  "state": "error",
  "type": "verify",
  "task_id": "task-100",
  "payload": {
    "result": "fail",
    "reason": "test regression"
  }
}
```

---


