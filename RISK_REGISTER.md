# Risk Register

| ID | 리스크 | 영향도 | 가능성 | 대응 전략 | 상태 |
|---|---|---:|---:|---|---|
| R-01 | OMC 이벤트 포맷 변경 | 높음 | 중간 | Normalizer 버전 분리, unknown fallback | Open |
| R-02 | 로그 누락/지연 | 중간 | 중간 | 다중 소스 수집 + timestamp 보정 | Open |
| R-03 | 이벤트 폭주로 TUI 프레임 저하 | 높음 | 중간 | 샘플링/축약 모드, 부분 렌더 | Open |
| R-04 | 민감정보 노출 | 매우높음 | 낮음 | redaction 기본 ON, 검증 체크리스트 | Open |
| R-05 | 모드별 상태 해석 불일치 | 중간 | 중간 | mode별 규칙 문서화, DECISIONS 고정 | Open |
| R-06 | replay 재현 불일치 | 중간 | 중간 | canonical schema + deterministic clock | Open |
| R-07 | 팀 내 용어 불일치 | 낮음 | 중간 | glossary/ADR(DECISIONS.md) 운영 | Open |

## 운영 규칙
- 새 리스크 발견 시 즉시 추가
- High 이상은 NEXT_ACTIONS에 대응 태스크 등록
- 상태: Open / Mitigating / Closed
