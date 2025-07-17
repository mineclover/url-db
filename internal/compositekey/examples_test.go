package compositekey_test

import (
	"fmt"
	"log"

	"url-db/internal/compositekey"
)

func ExampleService_Create() {
	service := compositekey.NewService("url-db")
	
	// 정규화된 합성키 생성
	key, err := service.Create("Tech Articles", 123)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(key)
	// Output: url-db:tech-articles:123
}

func ExampleService_CreateWithTool() {
	service := compositekey.NewService("url-db")
	
	// 사용자 정의 도구명으로 합성키 생성
	key, err := service.CreateWithTool("Bookmark Manager", "Personal Links", 456)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(key)
	// Output: bookmark-manager:personal-links:456
}

func ExampleService_Parse() {
	service := compositekey.NewService("url-db")
	
	// 합성키 파싱
	ck, err := service.Parse("url-db:tech-articles:123")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Tool: %s, Domain: %s, ID: %d\n", ck.ToolName, ck.DomainName, ck.ID)
	// Output: Tool: url-db, Domain: tech-articles, ID: 123
}

func ExampleService_ParseComponents() {
	service := compositekey.NewService("url-db")
	
	// 합성키를 구성 요소로 분해
	toolName, domainName, id, err := service.ParseComponents("url-db:tech-articles:123")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Tool: %s, Domain: %s, ID: %d\n", toolName, domainName, id)
	// Output: Tool: url-db, Domain: tech-articles, ID: 123
}

func ExampleService_Validate() {
	service := compositekey.NewService("url-db")
	
	// 합성키 유효성 검사
	isValid := service.Validate("url-db:tech-articles:123")
	fmt.Printf("Valid: %t\n", isValid)
	
	isValid = service.Validate("invalid-format")
	fmt.Printf("Valid: %t\n", isValid)
	
	// Output: Valid: true
	// Valid: false
}

func ExampleService_GetToolName() {
	service := compositekey.NewService("url-db")
	
	// 합성키에서 도구명 추출
	toolName, err := service.GetToolName("url-db:tech-articles:123")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(toolName)
	// Output: url-db
}

func ExampleService_GetDomainName() {
	service := compositekey.NewService("url-db")
	
	// 합성키에서 도메인명 추출
	domainName, err := service.GetDomainName("url-db:tech-articles:123")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(domainName)
	// Output: tech-articles
}

func ExampleService_GetID() {
	service := compositekey.NewService("url-db")
	
	// 합성키에서 ID 추출
	id, err := service.GetID("url-db:tech-articles:123")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(id)
	// Output: 123
}

func ExampleService_NormalizeComponents() {
	service := compositekey.NewService("url-db")
	
	// 구성 요소 정규화
	toolName, domainName, err := service.NormalizeComponents("URL Database", "Tech Articles")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Tool: %s, Domain: %s\n", toolName, domainName)
	// Output: Tool: url-database, Domain: tech-articles
}

func ExampleCreateNormalized() {
	// 정규화된 합성키 생성
	ck, err := compositekey.CreateNormalized("URL Database", "Tech Articles", 123)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(ck.String())
	// Output: url-database:tech-articles:123
}

func ExampleParse() {
	// 합성키 파싱
	ck, err := compositekey.Parse("url-db:tech-articles:123")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Tool: %s, Domain: %s, ID: %d\n", ck.ToolName, ck.DomainName, ck.ID)
	// Output: Tool: url-db, Domain: tech-articles, ID: 123
}

func ExampleIsValid() {
	// 합성키 유효성 검사
	isValid := compositekey.IsValid("url-db:tech-articles:123")
	fmt.Printf("Valid: %t\n", isValid)
	
	isValid = compositekey.IsValid("invalid-format")
	fmt.Printf("Valid: %t\n", isValid)
	
	// Output: Valid: true
	// Valid: false
}

func ExampleNormalizeToolName() {
	// 도구명 정규화
	normalized, err := compositekey.NormalizeToolName("URL Database")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(normalized)
	// Output: url-database
}

func ExampleNormalizeDomainName() {
	// 도메인명 정규화
	normalized, err := compositekey.NormalizeDomainName("Tech Articles")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(normalized)
	// Output: tech-articles
}