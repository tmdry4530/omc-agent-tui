# Keybindings (OMC Agent TUI)

## 1) 기본 원칙
- 모든 단축키는 한 손 조작 가능하게 설계
- 패널 포커스 기반 동작 + 전역 단축키 분리
- 충돌 키는 금지(중복 매핑 금지)

---

## 2) 전역 키

| 키 | 동작 | 비고 |
|---|---|---|
| `q` | 앱 종료 | 즉시 종료 |
| `tab` | 다음 패널 포커스 | 순환 |
| `shift+tab` | 이전 패널 포커스 | 순환 |
| `/` | 필터/검색 입력 열기 | 공통 검색바 |
| `esc` | 입력/필터 취소 | 기본 복귀 |
| `r` | Live ↔️ Replay 전환 | 모드 토글 |
| `p` | 일시정지/재개 | live soft pause |
| `m` | 마스코트 렌더 on/off | 성능 저하시 off 권장 |
| `?` | 단축키 도움말 | 오버레이 |

---

## 3) Timeline 패널

| 키 | 동작 |
|---|---|
| `j` / `k` | 아래/위 이벤트 이동 |
| `g` | 맨 위로 이동 |
| `G` | 맨 아래로 이동 |
| `enter` | 이벤트 상세 보기 |
| `y` | 이벤트 payload 복사 |
| `f` | 빠른 필터(타입/레벨) |

---

## 4) Agent Arena 패널

| 키 | 동작 |
|---|---|
| `h/j/k/l` | 카드 이동 |
| `enter` | 선택 에이전트 Inspector 열기 |
| `s` | 상태 기준 정렬 토글 |
| `o` | 역할 기준 정렬 토글 |
| `b` | blocker만 보기 토글 |
| `e` | error만 보기 토글 |

---

## 5) Task Graph 패널

| 키 | 동작 |
|---|---|
| `g` | 그래프 밀도 토글(요약/상세) |
| `[` / `]` | depth 감소/증가 |
| `c` | critical path 하이라이트 토글 |
| `x` | blocker chain 하이라이트 토글 |
| `enter` | 노드 상세(태스크 정보) |

---

## 6) Inspector 패널

| 키 | 동작 |
|---|---|
| `1` | Summary 탭 |
| `2` | Intent 탭 |
| `3` | Action 탭 |
| `4` | Diff 탭 |
| `5` | Verify/Fix 히스토리 탭 |
| `u` | 선택 에이전트 최신 상태로 리프레시 |

---

## 7) Replay 모드 전용

| 키 | 동작 |
|---|---|
| `space` | 재생/일시정지 |
| `←` / `→` | 한 스텝 뒤/앞 |
| `1` / `4` / `8` / `16` | 재생 배속 변경 |
| `t` | 특정 시각 점프 |
| `n` | 다음 에러 이벤트로 점프 |
| `N` | 이전 에러 이벤트로 점프 |

---

## 8) 필터 문법

입력창(`/`)에서 사용:

- `agent:<id|*>`
- `type:<eventType|*>`
- `state:<running|blocked|error|*>`
- `provider:<claude|gemini|codex|*>`
- `mode:<ralph|ultrawork|ultrapilot|*>`
- `text:<keyword>`

예시:
- `agent:coder-auth type:error`
- `provider:codex state:blocked`
- `mode:ralph type:verify`

---

## 9) 접근성/대체 키(선택)

- 화살표 불편 시 `hjkl` 완전 지원
- 컬러 비활성 모드에서 상태 아이콘 강화:
  - RUNNING `▶️`
  - WAITING `…`
  - BLOCKED `!`
  - ERROR `✖️`
  - DONE `✓`

---

## 10) 변경 관리

- 단축키 변경 시 `DECISIONS.md`에 이유 기록
- 사용자 피드백으로 충돌 키 수정 시 version bump와 함께 공지
