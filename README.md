# MCP (Model Chat Protocol)

MCP는 대화형 AI 모델과의 통신을 위한 프로토콜과 구현체를 제공하는 프로젝트입니다.

## Branch

- <a href="https://github.com/zkfmapf123/mcp-project-template/tree/mcp/go"> MCP-Golang </a>
- <a href="https://github.com/zkfmapf123/mcp-project-template/tree/mcp/ts"> MCP-Typescript </a>

## 프로젝트 구조

```sh
├── cmd/                   # 실행 파일 디렉토리
│   └── server/            # 서버 구현체
├── internal/              # 내부 패키지
│   ├── context/           # 대화 컨텍스트 관리
│   ├── model/             # AI 모델 구현체
│   └── utils/             # 유틸리티 함수
├── pkg/                   # 공개 패키지
│   └── protocol/          # MCP 프로토콜 정의
├── web/                   # 웹 인터페이스
│   └── static/            # 정적 파일
```

## Protocol

### JSON

```go
type Message struct {
    ID        string    `json:"id"`
    Content   string    `json:"content"`
    TimeStamp time.Time `json:"timestamp"`
    ContextID string    `json:"context_id"`
}

|- contextID
    |- messageId
    |- messageId
    |- messageId
    |- messageId
    |- messageId
    ...
```

### Protocol Buffers

```sh

```
