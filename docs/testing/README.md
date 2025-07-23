# URL-DB 테스트 문서

## 개요

이 디렉토리는 URL-DB 프로젝트의 테스트 전략과 설계 원칙을 정의하는 문서들을 포함합니다. URL-DB는 수만 URL 그래프를 실시간으로 전파·분석하는 MCP(Management Control Plane) 서버로, **6-단계 테스트 피라미드**를 기반으로 한 체계적인 테스트 전략을 제공합니다.

## 문서 목록

### 1. [TEST_STRATEGY.md](./TEST_STRATEGY.md)
**Clean Architecture 기반 테스트 전략**

- **6-단계 테스트 피라미드**: Unit 40%, Property-Based & Fuzz 15%, Integration 20%, Contract 10%, E2E 10%, Performance 5%
- 레이어별 테스트 전략 (Domain, Application, Infrastructure, Interface)
- 고급 테스트 기법 (속성 기반, 퍼즈, 계약, 동시성 테스트)
- CI/CD 자동화 및 커버리지 거버넌스
- 성능 및 부하 테스트 전략

**주요 내용:**
- Domain Layer: 엔티티, 값 객체, 도메인 서비스 테스트
- Application Layer: 유스케이스, DTO 테스트
- Infrastructure Layer: 리포지토리, 데이터베이스 테스트
- Interface Layer: HTTP/MCP 핸들러 테스트
- 고급 기법: Rapid, bufconn, Pact, goleak, Testcontainers

### 2. [TEST_DESIGN_PRINCIPLES.md](./TEST_DESIGN_PRINCIPLES.md)
**테스트 친화적 설계 원칙**

- 의존성 주입 (Dependency Injection)
- 순수 함수 (Pure Functions)
- 값 객체 (Value Objects)
- 팩토리 패턴 (Factory Pattern)
- 인터페이스 분리 (Interface Segregation)
- 이벤트 기반 설계 (Event-Driven Design)
- 설정 주입 (Configuration Injection)
- 시간 의존성 분리 (Time Dependency Separation)
- 에러 처리 표준화 (Error Handling Standardization)
- 테스트 헬퍼 함수 (Test Helper Functions)

**주요 내용:**
- 각 원칙의 문제점과 해결책
- Before/After 코드 예시
- 테스트 작성 가이드라인
- 모킹 전략

### 3. [PRACTICAL_EXAMPLES.md](./PRACTICAL_EXAMPLES.md)
**실제 적용 예시**

- 도메인 엔티티 설계 개선
- 값 객체 설계 개선
- 서비스 레이어 설계 개선
- 팩토리 패턴 적용
- 이벤트 기반 설계
- 설정 주입 패턴

**주요 내용:**
- 구체적인 코드 예시
- 테스트 코드 예시
- 실제 URL-DB 프로젝트 적용 사례

### 4. [ADVANCED_TESTING_GUIDE.md](./ADVANCED_TESTING_GUIDE.md)
**고성능 Go MCP 서버 테스트 가이드**

- **6-단계 피라미드 재구성**: 경계·동시성·계약·성능·보안까지 체계적 검증
- **Unit Layer 강화**: 속성 기반·퍼즈 테스트, 경쟁 조건 탐지, 고루틴 누수 방지
- **Integration Layer 고도화**: gRPC bufconn, Testcontainers, Pact 계약 테스트
- **전-시스템 검증**: 시나리오 드리븐 E2E, 벤치마크·프로파일, 부하·스트레스 테스트
- **CI/CD 자동화**: 커버리지 거버넌스, 동시성·컨텍스트·시간 의존성 테스트

**주요 내용:**
- Rapid, gopter, Go 1.18+ fuzzing을 활용한 속성 기반 테스트
- bufconn, testcontainers-go를 활용한 경량화된 통합 테스트
- Pact를 활용한 소비자 주도 계약 테스트
- 동시성 안전을 위한 레이스 감지 및 고루틴 누수 방지

## 사용법

### 1. 테스트 전략 수립
```bash
# 테스트 전략 문서 참조
cat docs/testing/TEST_STRATEGY.md
```

### 2. 설계 원칙 적용
```bash
# 설계 원칙 문서 참조
cat docs/testing/TEST_DESIGN_PRINCIPLES.md
```

### 3. 실제 적용 예시 참조
```bash
# 실제 적용 예시 문서 참조
cat docs/testing/PRACTICAL_EXAMPLES.md
```

### 4. 고급 테스트 기법 학습
```bash
# 고급 테스트 가이드 참조
cat docs/testing/ADVANCED_TESTING_GUIDE.md
```

## 테스트 작성 워크플로우

### 1. 설계 단계
1. **TEST_DESIGN_PRINCIPLES.md** 참조하여 테스트 친화적 설계 적용
2. 의존성 주입과 인터페이스 분리 고려
3. 순수 함수와 값 객체 사용

### 2. 구현 단계
1. **PRACTICAL_EXAMPLES.md** 참조하여 구체적인 구현 방법 적용
2. 팩토리 패턴과 헬퍼 함수 사용
3. 표준화된 에러 처리 적용

### 3. 테스트 단계
1. **TEST_STRATEGY.md** 참조하여 적절한 테스트 레벨 선택
2. **6-단계 피라미드**에 따른 우선순위 설정
3. 레이어별 테스트 전략 적용

### 4. 고급 테스트 단계
1. **ADVANCED_TESTING_GUIDE.md** 참조하여 고급 테스트 기법 적용
2. 속성 기반, 퍼즈, 계약 테스트 구현
3. 동시성 안전 및 성능 테스트 추가

## 6-단계 테스트 피라미드

```
    Performance & Load Tests (5%)
           ▲
    End-to-End Scenario Tests (10%)
           ▲
    Contract Tests (10%)
           ▲
    Integration Tests (20%)
           ▲
    Property-Based & Fuzz Tests (15%)
           ▲
    Unit Tests (40%)
```

### 각 단계별 목적과 도구

| 단계 | 비율 | 목적 | 도구 | 대상 |
|------|------|------|------|------|
| **Unit Tests** | 40% | 순수 함수, 값 객체, DTO | `testing`, `testify` | Domain, Application Layer |
| **Property-Based & Fuzz Tests** | 15% | 경계·인코딩 취약점 | `rapid`, `gopter`, Go 1.18+ fuzzing | Value Objects, Parsers |
| **Integration Tests** | 20% | gRPC bufconn, Testcontainers | `bufconn`, `testcontainers-go` | Infrastructure Layer |
| **Contract Tests** | 10% | Pact (Producer·Consumer) | `pact-go` | Interface Layer |
| **End-to-End Scenario Tests** | 10% | 실제 MCP workflow | `httptest`, `grpctest` | 전체 시스템 |
| **Performance & Load Tests** | 5% | benchmark + pprof | `testing.B`, `pprof` | 핫패스, 부하 |

## 핵심 원칙 요약

### 테스트 친화적 설계
- **의존성 주입**: 인터페이스를 통한 느슨한 결합
- **순수 함수**: 부수 효과 없는 비즈니스 로직
- **값 객체**: 불변 객체로 상태 관리 단순화
- **팩토리 패턴**: 복잡한 객체 생성 단순화

### 테스트 작성 원칙
- **AAA 패턴**: Arrange-Act-Assert
- **명확한 테스트 이름**: `Test[Function]_[Scenario]_[ExpectedResult]`
- **테스트 격리**: 독립적인 테스트 실행
- **헬퍼 함수**: 재사용 가능한 테스트 유틸리티

### 고급 테스트 기법
- **속성 기반 테스트**: Rapid, gopter를 활용한 자동 테스트 케이스 생성
- **퍼즈 테스트**: Go 1.18+ 내장 fuzzing으로 경계값 취약점 탐지
- **계약 테스트**: Pact를 활용한 서비스 간 호환성 보장
- **동시성 테스트**: 레이스 감지, 고루틴 누수 방지

### 테스트 전략
- **6-단계 피라미드**: 기능과 신뢰성을 동시에 담보하는 체계적 구조
- **레이어별 접근**: 각 레이어의 책임에 맞는 테스트
- **비즈니스 가치**: 핵심 비즈니스 로직 우선 테스트

## Clean Architecture와의 연관성

### Domain Layer
- **엔티티**: 비즈니스 규칙 검증 테스트
- **값 객체**: 불변성과 유효성 검증 테스트
- **도메인 서비스**: 복잡한 비즈니스 로직 테스트

### Application Layer
- **유스케이스**: 애플리케이션 비즈니스 규칙 테스트
- **DTO**: 데이터 변환 로직 테스트
- **애플리케이션 서비스**: 오케스트레이션 로직 테스트

### Infrastructure Layer
- **리포지토리**: 데이터 접근 로직 테스트
- **외부 서비스**: 통합 테스트
- **설정**: 환경별 설정 테스트

### Interface Layer
- **HTTP 핸들러**: REST API 테스트
- **MCP 핸들러**: 프로토콜 구현 테스트
- **CLI**: 명령줄 인터페이스 테스트

## 성능 고려사항

### 테스트 실행 속도
- **단위 테스트**: < 1초
- **통합 테스트**: < 10초
- **E2E 테스트**: < 60초

### 메모리 사용량
- **모킹 사용**: 실제 객체 대신 가벼운 모의 객체
- **테스트 격리**: 각 테스트 후 정리
- **리소스 관리**: 데이터베이스 연결 풀 관리

### 병렬 실행
- **독립적인 테스트**: 순서에 의존하지 않는 테스트
- **공유 상태 최소화**: 전역 변수 사용 금지
- **테스트 데이터 격리**: 각 테스트별 독립적인 데이터

## 고급 도구 스택

### 속성 기반 테스트
- **Rapid**: 현대적인 Go 속성 기반 테스트 라이브러리
- **Gopter**: Go 속성 기반 테스트 프레임워크
- **Go 1.18+ Fuzzing**: 내장 퍼즈 테스트 기능

### 통합 테스트
- **bufconn**: gRPC 인-메모리 테스트
- **testcontainers-go**: 컨테이너 기반 통합 테스트
- **httptest**: HTTP 서버 테스트

### 계약 테스트
- **Pact**: 소비자 주도 계약 테스트
- **Pact Broker**: 계약 관리 및 검증

### 동시성 테스트
- **Race Detector**: Go 내장 레이스 감지
- **goleak**: 고루틴 누수 감지
- **Chronos**: 정적 레이스 스캐너

### 성능 테스트
- **testing.B**: 벤치마크 테스트
- **pprof**: 성능 프로파일링
- **vegeta**: HTTP 부하 테스트 도구

## 결론

이 문서들은 URL-DB 프로젝트에서 테스트 가능한 코드를 작성하고, 효과적인 테스트를 구현하기 위한 가이드를 제공합니다.

**핵심 메시지:**
1. **테스트 친화적 설계**가 먼저다
2. **Clean Architecture**와 테스트는 자연스럽게 맞는다
3. **6-단계 피라미드**로 기능과 신뢰성을 동시에 담보한다
4. **고급 테스트 기법**으로 예상치 못한 버그를 탐지한다
5. **동시성 안전**과 **성능 보장**으로 안정성을 확보한다
6. **지속적인 개선**이 필요하다

이 원칙들을 따르면 테스트 작성이 쉬워지고, 코드 품질이 향상되며, 유지보수가 용이해집니다. 특히 MCP 서버의 고병렬 환경에서도 데이터 정합성과 자원 누수를 예방할 수 있습니다. 