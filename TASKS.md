# OMC Agent TUI TASKS

## M0. 부트스트랩
- [ ] 프로젝트 문서 구조 확정
- [ ] status/* 초기화
- [ ] references/* 핵심 스키마 채우기

## M1. 이벤트 스키마/정규화
- [ ] CanonicalEvent 필드 확정
- [ ] OMC mode 매핑 규칙 확정
- [ ] redaction 규칙 확정

## M2. 컴포넌트 설계
- [ ] Agent Arena 상세 명세
- [ ] Timeline 상세 명세
- [ ] Task Graph 상세 명세
- [ ] Inspector 상세 명세

## M3. 아키텍처 고정
- [ ] 모듈 경계 확정(collector/normalizer/store/tui/replay)
- [ ] 장애/복구 전략 확정
- [ ] 성능 목표 수치 확정

## M4. UX 시나리오
- [ ] 정상 시나리오 3개
- [ ] 실패 시나리오 3개
- [ ] replay 시나리오 2개

## M5. 실행 준비
- [ ] MASTER_PROMPT 최종 고정
- [ ] PHASE 0/1/2 프롬프트 고정
- [ ] 우선순위 Top 5 재정렬

---

## 강제 운영 원칙
- [ ] OMC Ralph/Ultrawork/Ultrapilot로 A→Z 완주
- [ ] 실패 시 verify-fix 루프
- [ ] 부분완료 금지 (구현+검증+보고)
- [ ] 환경변수는 사용자 세팅 전제
