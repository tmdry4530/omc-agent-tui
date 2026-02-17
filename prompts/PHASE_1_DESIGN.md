# PHASE 1 — DESIGN (OMC Agent TUI)

## 0) Preflight 게이트 (필수)
설계 작업 시작 전 Gemini/Codex MCP 호출 상태를 확인한다.

절차:
1. OMC 활성 확인
2. Gemini 테스트 호출 1회
3. Codex 테스트 호출 1회
4. 판정:
   - 둘 다 성공: `PRECHECK: PASS`
   - 하나라도 실패: `PRECHECK: FAIL` + 작업 취소

실패 시:
- `status/FAILURES.md` 기록
- 사용자에게 취소 사유 보고
- 우회 진행 금지

---

## 1) 목표
구현 가능한 수준의 컴포넌트/데이터/상태머신 설계를 확정한다.
(코드 작성 금지)

## 2) 입력 문서
- PRD.md
- COMPONENT_SPEC.md
- ARCHITECTURE.md
- references/event-schema.md
- references/keybindings.md
- references/mascot-guidelines.md
- assets/mascot/palette.md
- status/DECISIONS.md

## 3) 수행 지침
1. 패널별 책임/입출력/의존성 정의
2. CanonicalEvent 기반 상태 전이 검증
3. Intent vs Action diff 규칙 설계
4. replay/live 공용 경로 및 예외처리 정의
5. OMC 모드별 UI 차이 반영
6. 접근성/성능/보안 요건 반영

## 4) 출력 형식
1) Design Decisions
2) Component Contracts
3) State Machine Summary
4) Risk & Trade-offs
5) Next Actions(최대 5개)
