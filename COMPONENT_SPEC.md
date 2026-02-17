# OMC Agent TUI 컴포넌트 명세 v0.1

## 1) Panel A — Agent Arena
목적: 에이전트 상태를 즉시 파악

표시:
- 마스코트(색상+역할 배지)
- agent id/name
- state badge
- progress(0~100)
- blocked_by

정렬 우선순위:
`error > blocked > running > waiting > idle/done`

상태 스타일:
- running: pulse
- waiting: dim
- blocked: yellow border
- error: red border
- done: check badge

---

## 2) Panel B — Live Timeline
목적: 이벤트 흐름 추적

행 포맷:
`HH:MM:SS | agent | type | summary`

필터:
- `/agent:`
- `/type:`
- `/provider:`
- `/state:`
- `/mode:`

동작:
- Enter: 상세 보기
- y: payload 복사(선택)

---

## 3) Panel C — Task Graph
목적: 태스크 관계/병목 시각화

표시:
- parent-child DAG
- critical path
- blocker chain

규칙:
- 기본 depth 3
- 이벤트 증가 시 그래프 축약 모드

---

## 4) Panel D — Inspector
목적: 선택 에이전트/태스크 상세 분석

섹션:
1. Summary
2. Intent(plan)
3. Action(tool calls/results)
4. Diff(intent vs action)
5. Verify/Fix history

핵심 지표:
- avg latency
- retries
- error count
- success ratio

---

## 5) Panel E — Footer Metrics
표시:
- tokens in/out
- avg latency
- error rate
- cost estimate
- live/replay 상태

---

## 6) 키바인딩 핵심
- q 종료
- tab 포커스 이동
- / 필터
- r replay 토글
- m mascot 토글
- p pause
- space (replay play/pause)
