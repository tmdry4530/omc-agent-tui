# OMC Agent TUI 아키텍처 v0.1

## 1) 목적
OMC(Oh My ClaudeCode) 멀티에이전트 실행 이벤트를 수집/정규화해서 TUI로 실시간 시각화한다.

## 2) 아키텍처 개요
파이프라인:
`OMC/ClaudeCode logs/hooks -> Collector -> Normalizer/Redactor -> State Store -> TUI Renderer`

## 3) 모듈 경계

### 3.1 Collector
- 입력: OMC 훅, stdout/stderr, 파일 로그(JSONL/txt)
- 출력: RawEvent
- 책임: 입력 수집, 타임스탬프 보정, 소스 장애 격리

### 3.2 Normalizer
- RawEvent -> CanonicalEvent 변환
- provider/mode/role/state/type 정규화
- 민감정보 마스킹(redaction)

### 3.3 Store
- RunState/AgentState/TaskState 유지
- 최근 이벤트 ring buffer
- 메트릭 집계(latency/tokens/error/cost)

### 3.4 TUI
- Agent Arena / Timeline / Task Graph / Inspector / Footer 렌더
- 필터/포커스/키바인딩 처리
- live/replay 모드 전환

### 3.5 Replay
- JSONL 로딩
- 가상 시계 기반 재생(1x/4x/8x/16x)

## 4) 데이터 흐름
1. Collector 수신
2. Normalizer 변환 + Redaction
3. Store 상태 갱신
4. Renderer 부분 업데이트
5. Metrics/Footer 계산

## 5) 장애/복구 전략
- 입력 소스 실패: 소스 단위 circuit-breaker
- 파싱 실패: 이벤트 drop + warning count 증가
- 렌더 실패: 패널 단위 recover(프로세스 종료 금지)
- burst 트래픽: 샘플링/축약 모드 전환

## 6) 보안
- redaction 기본 ON
- token/key/secret/password 패턴 마스킹
- raw 로그 저장은 opt-in

## 7) OMC 모드 반영
- Ralph: verify-fix 루프 강조
- Ultrawork: fan-out 병렬 작업 강조
- Ultrapilot: A→Z 연속 실행 추적 강조

## 8) 확장 포인트
- provider plugin
- exporter(JSON/CSV/OTel)
- alert sink(webhook/desktop)
