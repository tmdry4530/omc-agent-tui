# PRD — OMC Agent TUI

## 1. 문제 정의
OMC(Oh My ClaudeCode) 기반 멀티에이전트 실행은 강력하지만, 실행 중 “누가/무엇을/왜” 하는지 즉시 파악하기 어렵다.  
현재는 로그 중심 관찰이라 병목, 실패 원인, verify-fix 루프의 상태를 빠르게 추적하기 힘들다.

### 현재 Pain Point
- 에이전트별 현재 상태(실행/대기/블록/오류) 가시성 부족
- 계획(Intent)과 실제 실행(Action)의 괴리 추적 어려움
- 장애/재시도/병목이 텍스트 로그에 묻힘
- 세션 종료 후 재현(replay) 난이도 높음

---

## 2. 목표
### 제품 목표
1. 사용자(개발자)가 **2초 이내**에 현재 멀티에이전트 상태를 이해한다.
2. 사용자(개발자)가 **1분 이내**에 병목 원인을 특정한다.
3. 실행 세션을 replay하여 실패 원인을 **재현 가능**하게 만든다.

### 운영 목표
- OMC Ralph/Ultrawork/Ultrapilot 모드에서 A→Z 실행 관찰 지원
- 실패 발생 시 verify-fix 루프 상태를 명확히 시각화

### 비목표(초기)
- TUI에서 코드 직접 편집/실행 제어
- 브라우저 UI 우선 개발
- 벤더별 API 심층 통합(초기에는 이벤트 레벨 통합)

---

## 3. 핵심 사용자
- Claude Code + OMC로 실제 개발하는 개인/팀 개발자
- 멀티에이전트 자동화 결과를 검증/관제해야 하는 Tech Lead

---

## 4. 핵심 기능 (Functional Requirements)
### FR-1 Agent Arena
- 에이전트별 마스코트 카드 표시
- 상태: idle/running/waiting/blocked/error/done
- 역할/진행률/blocked-by 표시

### FR-2 Live Timeline
- 실시간 이벤트 스트림 표시
- 필터: agent/type/provider/level
- 상세 이벤트 inspect 지원

### FR-3 Task Graph
- parent-child 태스크 DAG 시각화
- critical path 표시
- blocker chain 표시

### FR-4 Inspector (Intent vs Action)
- 선택 에이전트의 계획(Intent)과 실제 행동(Action) 비교
- 이탈(diff), 재시도, verify-fix 히스토리 표시

### FR-5 Replay
- JSONL 기반 세션 재생
- 속도 조절(1x/4x/8x/16x), step 탐색

### FR-6 Security Redaction
- token/key/secret/password 자동 마스킹 기본 ON

---

## 5. 비기능 요구사항 (NFR)
- 성능: 50 events/s 처리 시 UI 사용 가능 수준 유지
- 안정성: 파서 오류가 발생해도 프로세스 전체 종료 금지
- 가독성: 80컬럼 터미널 fallback 제공
- 보안: 민감정보 노출 방지 기본값 유지
- 이식성: 로컬 터미널 환경에서 동작

---

## 6. 성공 지표 (Success Metrics)
- 병목 식별 평균 시간: 1분 이내
- 장애 원인 회귀 분석 성공률: 80% 이상
- 사용자가 체감하는 디버깅 시간: 기존 대비 30% 단축(정성+로그 지표)

---

## 7. 제약사항
- 환경변수는 사용자 직접 세팅(본 프로젝트는 ENV 존재를 가정)
- OMC 및 외부 MCP 이벤트 형식 변동 가능성 존재
- 초기 단계는 “관찰/분석” 중심, 실행 제어는 제외

---

## 8. 운영 원칙
- OMC Ralph/Ultrawork/Ultrapilot 기반 A→Z 완주 추적
- 실패 시 중단하지 않고 verify-fix 루프 중심으로 복구 상태 표시
- 부분완료 보고 금지: 구현+검증+결과 보고까지 완료 기준

---

## 9. 릴리즈 범위
### v0.1 (MVP)
- Agent Arena + Timeline + Graph + Inspector + Replay(기본)
- Redaction 기본 적용

### v0.2
- Intent/Action diff 고도화
- 비용/토큰 히트맵
- 알림(선택)

### v0.3
- 플러그인형 provider 확장
- 운영 리포트 자동 요약
