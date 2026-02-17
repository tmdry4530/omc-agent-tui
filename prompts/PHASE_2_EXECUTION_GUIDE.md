# PHASE 2 — EXECUTION GUIDE (OMC Agent TUI)

## 0) 실행 전 게이트 (필수)
구현/수정 작업 전에 반드시 Preflight를 수행한다.

### Preflight 절차
1. OMC 활성 상태 확인
2. Gemini MCP 테스트 호출 1회
3. Codex MCP 테스트 호출 1회
4. 결과 기록:
   - 성공/실패
   - latency
   - 에러 메시지(실패 시)
5. 판정:
   - 둘 다 성공 → `PRECHECK: PASS` 후 진행
   - 하나라도 실패 → `PRECHECK: FAIL` + **즉시 작업 취소**

### Fail 처리
- `status/FAILURES.md`에 상세 기록
- 사용자에게 “MCP 호출 실패로 작업 취소” 보고
- 우회 구현/대체 구현 금지

---

## 1) 핵심 실행 원칙
- OMC Ralph / Ultrawork / Ultrapilot 중심 A→Z 완주
- 실패 시 verify-fix 루프 기반 복구
- 부분완료 보고 금지 (구현+검증+결과보고까지 완료)
- 환경변수는 사용자 세팅 전제

---

## 2) 작업 흐름
1. Preflight PASS 확인
2. TASKS.md에서 현재 마일스톤 선택
3. 이번 턴 목표(1~3개) 설정
4. 작업 수행
5. 검증 증거 수집
6. status 파일 업데이트
7. 다음 액션 제시

---

## 3) 검증 규칙
- 기능: PRD/COMPONENT_SPEC 요구사항 체크
- 안정성: 실패/복구 시나리오 확인
- 보안: redaction 정책 준수 확인
- 성능: 이벤트 burst 상황 가독성 확인

---

## 4) 실패 기록 템플릿 (필수)
`status/FAILURES.md`에 아래 항목 기록:
- 일시
- 단계
- 증상
- 재현 방법
- 추정 원인
- 시도한 해결
- 결과
- 대안 1/2
- 다음 액션
- 상태(Open/Mitigating/Closed)

---

## 5) 매 턴 보고 형식
1) 이번 작업
2) 검증 결과(증거)
3) 변경 파일
4) 남은 리스크
5) 다음 액션(최대 3개)

---

## 6) 완료 정의 (DoD)
- 요구사항 충족
- 핵심 시나리오(정상/실패/replay) 검증 완료
- status/DECISIONS 최신화 완료
- 운영자가 바로 이어받을 수 있는 상태
