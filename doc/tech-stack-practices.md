# VimGo 技术栈与实践说明

## 1. 项目定位
VimGo 是一个基于 Go 实现的围棋应用，核心逻辑统一在后端，提供两种交互入口：
- 终端 UI（Bubble Tea + Lipgloss）
- Web UI（xterm.js + WebSocket + PTY，复用终端渲染效果）

## 2. 当前技术栈

### 2.1 语言与运行时
- Go `1.24.2`

### 2.2 终端 UI
- `github.com/charmbracelet/bubbletea`：TUI 状态机与事件循环
- `github.com/charmbracelet/lipgloss`：终端样式与布局

### 2.3 Web UI
- 原生 HTML/CSS/JavaScript
- `xterm.js`：浏览器终端模拟
- `xterm-addon-fit`：终端自适应尺寸
- `github.com/gorilla/websocket`：WebSocket 通道
- `github.com/creack/pty`：启动/管理伪终端子进程

### 2.4 核心业务模块
- `internal/board`：棋盘数据结构
- `internal/rules`：合法性、提子、计分
- `internal/game`：对局状态与回合管理
- `internal/vim`：Vim 模式与按键处理
- `internal/sgf`：SGF 读写

## 3. 目录与分层实践
- `cmd/vimgo`：终端程序入口
- `cmd/web`：Web 服务入口（静态资源 + `/ws`）
- `internal/*`：领域与应用逻辑，不暴露为公共库
- `web/static`：前端静态资源
- `doc`：项目文档

分层原则：
- 游戏规则和状态在 `internal` 中实现，UI 层不直接承载规则。
- Web 端通过 PTY 复用终端版本，避免维护两套渲染/交互逻辑。

## 4. 工程实践

### 4.1 代码组织
- 以小模块拆分核心能力（board/rules/game/vim/sgf）。
- 使用 Go 原生包管理与模块化，不引入复杂框架。

### 4.2 测试策略
- 以单元测试为主，覆盖核心规则与状态流转。
- 现有测试集中在：
  - `internal/board/board_test.go`
  - `internal/game/game_test.go`
  - `internal/rules/score_test.go`
- 重点验证内容：提子、打劫、悔棋恢复、计分规则与边界。

### 4.3 一致性与可维护性
- 使用 `gofmt` 保持格式一致。
- 优先修复“行为正确性”问题，再补回归测试防止回归。
- Web 端与终端端共享同一业务逻辑，减少功能漂移。

## 5. 运行与开发流程

### 5.1 常用命令
- 终端版运行：
```bash
go run ./cmd/vimgo
```
- Web 版运行（需先有 `vimgo` 可执行文件）：
```bash
go build -o vimgo ./cmd/vimgo
go run ./cmd/web
```
- 测试：
```bash
go test ./...
```

### 5.2 Web 版工作方式
- 浏览器连接 `/ws`。
- 服务端为每个连接启动一个 `vimgo` PTY 子进程。
- PTY 输出流实时转发到 xterm。
- 键盘输入与窗口 resize 事件通过 WebSocket 回传 PTY。

## 6. 当前规则实现边界
- 计分模块默认在“盘上棋子视为活棋”前提下计算。
- 中日计分都支持，非法 method 会回退到 Chinese。
- 悔棋已支持恢复提子数、当前行棋方、最后一手与记录长度。

## 7. 后续可演进方向
- 增加死活判定/终局确认流程，提升实战计分准确性。
- 增加集成测试，覆盖 `cmd/web` 的端到端行为。
- 将 Web 配置（端口、默认棋盘尺寸、二进制路径）参数化。
