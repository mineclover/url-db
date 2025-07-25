# URL-DB 기능 가이드

URL-DB는 AI 어시스턴트와 함께 사용할 수 있는 URL 관리 시스템입니다. 링크를 체계적으로 정리하고 태그를 붙여 쉽게 찾을 수 있게 도와줍니다.

## 🏠 도메인 관리 - URL 그룹 만들기

### 도메인이란?
URL을 주제별로 분류하는 폴더 같은 개념입니다. 예: 'tech', 'recipes', 'shopping'

### 할 수 있는 일
- **도메인 만들기**: 새로운 카테고리 생성 (예: "레시피" 도메인)
- **도메인 목록 보기**: 만든 모든 카테고리 확인
- **도메인 정보 수정**: 이름이나 설명 변경
- **도메인 삭제**: 더 이상 필요 없는 카테고리 제거

## 🔗 URL 관리 - 링크 저장하고 정리하기

### URL이란?
도메인 안에 저장되는 개별 웹사이트 링크입니다.

### 할 수 있는 일
- **URL 추가**: 좋은 웹사이트를 도메인에 저장
- **URL 목록 보기**: 도메인 안의 모든 링크 확인
- **URL 검색**: 제목이나 내용으로 링크 찾기
- **URL 정보 수정**: 제목이나 설명 변경
- **URL 삭제**: 필요 없는 링크 제거
- **URL로 찾기**: 정확한 웹주소로 링크 검색

## 🏷️ 태그 시스템 - 더 세밀한 분류

### 태그 타입
1. **일반 태그**: 색깔, 카테고리 등 (예: "빨강", "중요")
2. **순서 태그**: 우선순위가 있는 태그 (예: "1순위", "2순위")
3. **숫자**: 가격, 평점 등 (예: "29900", "4.5")
4. **텍스트**: 간단한 메모 (예: "나중에 읽기")
5. **마크다운**: 자세한 설명 (서식 포함)
6. **이미지**: 스크린샷이나 관련 이미지

### 할 수 있는 일
- **태그 타입 정의**: 도메인에서 사용할 태그 종류 설정
- **URL에 태그 붙이기**: 링크에 분류 정보 추가
- **태그로 검색**: 특정 태그가 붙은 링크만 찾기
- **태그 정보 보기**: URL에 붙은 모든 태그 확인

## 🔍 고급 검색 - 원하는 링크 쉽게 찾기

### 검색 방법
- **태그 조합 검색**: 여러 태그 조건으로 정확한 링크 찾기
- **검색 옵션**: 
  - 정확히 일치
  - 부분 포함
  - 시작 단어
  - 끝 단어

### 할 수 있는 일
- **복합 검색**: "가격이 3만원 이하이면서 리뷰가 좋은 제품" 같은 조건 검색
- **한 번에 모든 정보 보기**: URL과 태그를 동시에 확인

## 📮 이벤트 알림 - 변경사항 추적

### 이벤트란?
URL이나 태그가 변경될 때 발생하는 알림입니다.

### 할 수 있는 일
- **알림 구독**: 특정 URL의 변경사항 알림 받기
- **이벤트 종류**:
  - URL 추가됨
  - URL 수정됨
  - URL 삭제됨
  - 태그 변경됨
- **구독 관리**: 알림 설정 변경 또는 취소

## 🔗 연결 관리 - URL 간의 관계 설정

### 연결 타입
1. **강한 연결**: 한 URL이 삭제되면 연결된 URL도 영향받음
2. **약한 연결**: 참고용 연결, 삭제해도 영향 없음
3. **참조 연결**: 단순히 관련 있다는 표시

### 할 수 있는 일
- **URL 연결**: 관련 있는 링크들 연결
- **연결 관계 보기**: 어떤 URL이 어떻게 연결되어 있는지 확인
- **연결 해제**: 더 이상 관련 없는 링크들 분리

## 📊 시스템 정보 - 현재 상태 확인

### 확인할 수 있는 정보
- **서버 상태**: 시스템이 정상 작동하는지 확인
- **이벤트 통계**: 얼마나 많은 변경이 일어났는지 확인
- **처리 대기 목록**: 아직 처리되지 않은 작업들

## 🤖 AI 어시스턴트와 함께 사용하기

URL-DB는 Claude Desktop이나 Cursor 같은 AI 도구와 연결해서 사용할 수 있습니다.

### 자연어로 요청하기
- "기술 관련 사이트들을 보여줘"
- "가격이 10만원 이하인 제품 링크 찾아줘"
- "이 URL에 '중요' 태그 붙여줘"
- "레시피 도메인 만들어줘"

### AI가 도와주는 일
- 복잡한 검색 조건을 간단한 말로 요청
- 여러 단계 작업을 한 번에 처리
- 관련 링크들을 자동으로 찾아서 연결
- 태그를 분석해서 유용한 정보 제공

---

## 실제 사용 예시

### 예시 1: 요리 레시피 정리
1. "recipes" 도메인 만들기
2. 태그 타입 설정: "요리시간"(숫자), "난이도"(태그), "재료"(텍스트)
3. 레시피 사이트 추가하면서 태그 붙이기
4. "30분 이하로 만들 수 있는 쉬운 요리" 검색

### 예시 2: 기술 자료 수집
1. "tech" 도메인 만들기  
2. 태그 타입 설정: "프로그래밍언어"(태그), "난이도"(순서태그), "즐겨찾기"(태그)
3. 유용한 기술 문서들 저장하면서 분류
4. AI에게 "Python 초급 자료만 보여줘" 요청

### 예시 3: 쇼핑 정보 관리
1. "shopping" 도메인 만들기
2. 태그 타입 설정: "가격"(숫자), "브랜드"(태그), "리뷰점수"(숫자)
3. 관심 상품 링크들 저장
4. "10만원 이하이면서 리뷰 4점 이상인 제품" 검색

이렇게 URL-DB를 사용하면 인터넷에서 찾은 유용한 정보들을 체계적으로 정리하고, 나중에 쉽게 찾아볼 수 있습니다!