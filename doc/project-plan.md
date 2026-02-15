# VimGo 项目计划书

## 一、项目概述

**VimGo** 是一个融合 Vim 操作哲学与围棋文化的创新项目，旨在为用户提供独特的终端围棋体验。

### 1.1 项目愿景
- 让 Vim 用户能够以熟悉的快捷键方式体验围棋
- 保留 Vim 的操作习惯与界面风格
- 同时支持终端和网页双平台部署

### 1.2 核心特色
- **GBA 像素风格**：复古怀旧的视觉体验
- **Vim 快捷键操作**：hjkl 移动光标，数字+w/e/b 跳转
- **Vim 命令模式**：`:q` 退出、`:w` 保存棋谱、`:undo` 悔棋等
- **Vim 状态栏**：显示当前模式、回合、手数、坐标系

---

## 二、功能规划

### 2.1 基础功能

| 功能 | 描述 |
|------|------|
| 落子 | 在光标位置落子，自动判断气数 |
| 吃子 | 自动提子，统计提子数量 |
| 劫争检测 | 正确处理打劫规则 |
| 禁入点检测 | 禁止自杀和违反规则 |
| 终局判断 | 数子或点目判定胜负 |

### 2.2 Vim 快捷键

| 模式 | 快捷键 | 功能 |
|------|--------|------|
| Normal | `h/j/k/l` | 光标上下左右移动 |
| Normal | `w/e/b` | 跳转到下/上/前一位置 |
| Normal | `0/$` | 跳转到行首/行尾 |
| Normal | `gg/G` | 跳转到棋盘开头/末尾 |
| Normal | `x` | 当前光标位置落子 |
| Normal | `u` | 悔棋 (undo) |
| Normal | `Ctrl+r` | 重做 (redo) |
| Normal | `:` | 进入命令模式 |
| Normal | `/` | 进入搜索模式 |
| Command | `:q` | 退出游戏 |
| Command | `:q!` | 强制退出（不保存） |
| Command | `:w [file]` | 保存棋谱 |
| Command | `:wqa` | 保存并退出 |
| Command | `:undo` | 悔棋 |
| Command | `:redo` | 重做 |
| Command | `:pass` | 停一手 |
| Command | `:count` | 显示提子数 |
| Command | `:help` | 显示帮助 |
| Visual | `v` | 进入可视模式选择区域 |

### 2.3 Vim 状态栏

```
-- VimGo -- Normal -- 19x19 -- 黑方 -- 第 127 手 -- [A1] --
```

状态栏显示内容：
- 当前模式 (Normal/Insert/Visual/Command)
- 棋盘尺寸 (9x9/13x13/19x19)
- 当前执子方
- 当前手数
- 光标位置坐标
- 保存状态 (* 表示未保存)

---

## 三、技术架构

```
vimgo/
├── cmd/
│   └── vimgo/          # CLI 入口
├── internal/
│   ├── game/           # 围棋逻辑引擎
│   ├── board/          # 棋盘数据结构
│   ├── rules/          # 规则判定
│   ├── ui/
│   │   ├── terminal/   # 终端界面
│   │   └── web/        # 网页界面
│   ├── vim/            # Vim 模式与快捷键处理
│   └── sgf/            # SGF 棋谱读写
├── web/
│   ├── static/         # 前端资源
│   └── server.go       # Web 服务器
├── doc/                # 项目文档
└── go.mod              # Go 模块定义
```

### 3.1 终端版本 (TUI)

**技术选型**
- **Go + tview** 或 **bubbletea**：跨平台终端 UI 库
- **gocui**：更贴近 Vim 的体验

**特性**
- 256 色终端支持
- 响应式布局
- 鼠标支持（可选）

### 3.2 网页版本

**技术选型**
- **Go + WebAssembly**：核心逻辑编译为 WASM
- **React/Vue**：前端界面
- **Canvas/WebGL**：GBA 像素风格渲染

**特性**
- 像素风格渲染
- 响应式设计
- 键盘事件监听

---

## 四、像素风格 UI 技术方案

### 4.1 方案一：Go + Ebiten（桌面应用推荐）

**Ebiten** 是 Go 语言最成熟的 2D 游戏库，专为像素风格设计。

**项目地址**: https://github.com/hajimehoshi/ebiten

**优势**：
- 纯 Go 实现，无 C 依赖，跨平台编译简单
- 一键编译到 Windows/macOS/Linux/WebAssembly
- 内置像素绘图 API，性能优秀
- 自动处理不同分辨率和 DPI

**代码示例**：
```go
package main

import (
    "github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
    board *Board
}

func (g *Game) Update() error {
    // 游戏逻辑更新
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // 绘制像素棋盘
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Scale(3, 3)  // 像素放大 3 倍
    screen.DrawImage(g.board.Image, op)

    // 绘制棋子
    for _, stone := range g.board.Stones {
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(float64(stone.X*12), float64(stone.Y*12))
        screen.DrawImage(stone.Image, op)
    }
}

func main() {
    ebiten.SetWindowSize(720, 480)
    ebiten.SetWindowTitle("VimGo - 围棋")
    ebiten.RunGame(&Game{})
}
```

---

### 4.2 方案二：Phaser.js + Canvas（网页版推荐）

**Phaser.js** 是最流行的 HTML5 2D 游戏框架，像素效果极佳。

**项目地址**: https://phaser.io/

**优势**：
- 专为像素游戏设计，内置 `pixelArt: true` 模式
- 丰富的插件生态系统
- 易于分享，无需安装即可游玩
- 支持 Sprite 动画、特效、粒子系统

**代码示例**：
```javascript
const config = {
    type: Phaser.AUTO,
    width: 240,
    height: 160,
    parent: 'game-container',
    pixelArt: true,           // 关键：开启像素模式
    scale: {
        mode: Phaser.SCALE.FIT,
        autoCenter: Phaser.SCALE.CENTER_BOTH,
        zoom: 3                 // 放大 3 倍显示
    },
    scene: {
        preload: preload,
        create: create,
        update: update
    }
};

// GBA 风格调色板
const COLORS = {
    board: 0xE8CDA5,      // 竹色棋盘
    boardDark: 0xC9A05C,  // 深色格线
    black: 0x1A1A1A,      // 黑子
    white: 0xF5F5F5,      // 白子
    star: 0x8B4513,       // 星位
    cursor: 0xFF6B6B,     // 光标
    text: 0x2D2D2D        // 文字
};

const game = new Phaser.Game(config);
```

**像素增强效果**：
```javascript
// 扫描线效果
class ScanlineEffect extends Phaser.Scene {
    create() {
        const scanline = this.add.graphics();
        scanline.fillStyle(0x000000, 0.1);
        for (let y = 0; y < 160; y += 2) {
            scanline.fillRect(0, y, 240, 1);
        }
    }
}

// 像素抖动过渡
function drawWithDither(ctx, x, y, color) {
    ctx.imageSmoothingEnabled = false;
    // 棋盘抖动算法
}
```

---

### 4.3 方案三：终端增强 + 像素图混合（终端版）

在纯终端环境中显示像素风格画面。

**技术工具**：

| 工具 | 用途 |
|------|------|
| **chafa** | 在终端显示像素图片 |
| **icat** | kitty 终端图片协议 |
| **termpix** | Rust 写的终端图片查看器 |
| **viu** | 轻量级图片查看器 |

**安装**：
```bash
# Arch Linux
pacman -S chafa

# macOS
brew install chafa

# Ubuntu
apt install chafa
```

**使用示例**：
```bash
# 显示像素风格棋盘
chafa -s 80x40 -f sixel board.png

# 在游戏中嵌入图片
vimgo --render-image | chafa -s 80x50 -
```

**终端渲染流程**：
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  游戏逻辑    │ ──► │  渲染 Canvas │ ──► │  chafa 输出  │
│  (Go)       │     │  (240x160)  │     │  (终端显示)   │
└─────────────┘     └─────────────┘     └─────────────┘
```

---

### 4.4 技术架构（更新）

```
vimgo/
├── cmd/
│   ├── vimgo/           # CLI 入口（通用）
│   ├── terminal/       # 终端版（tview + chafa）
│   ├── desktop/         # 桌面版（Ebiten）
│   └── web/             # 网页版（Phaser.js）
├── internal/
│   ├── core/            # 核心围棋引擎（共享）
│   │   ├── board/       # 棋盘数据结构
│   │   ├── rules/       # 规则判定
│   │   └── sgf/         # SGF 棋谱读写
│   ├── vim/             # Vim 模式与快捷键处理
│   └── renderer/        # 渲染器接口
│       ├── ebiten/      # Ebiten 渲染实现
│       ├── phaser/      # Phaser 渲染实现
│       └── terminal/    # 终端渲染实现
├── assets/              # 静态资源
│   ├── sprites/         # 像素素材（.png）
│   │   ├── stones/      # 棋子精灵图
│   │   ├── ui/          # UI 元素
│   │   └── effects/     # 特效素材
│   └── fonts/           # 像素字体（.ttf）
└── web/                 # Web 服务器
    ├── static/          # 前端静态资源
    └── server.go        # HTTP 服务
```

---

### 4.5 GBA 像素风格设计规范

**分辨率标准**：
| 标准 | 分辨率 | 用途 |
|------|--------|------|
| GBA 原生 | 240x160 | 基础画布 |
| 放大 2x | 480x320 | 桌面窗口 |
| 放大 3x | 720x480 | 全屏显示 |
| 放大 4x | 960x640 | 高清展示 |

**像素字体推荐**：
| 字体 | 风格 | 下载 |
|------|------|------|
| TinyPico | 像素完美 | Google Fonts |
| Pixelicroff | 等宽像素 | GitHub |
|苹方像素| 中文像素| 方正字库 |
| Press Start 2P | 游戏风格 | Google Fonts |

**像素边框样式**：
```
// GBA 风格按钮
┌────────┐     ┌────────┐     ┌────────┐
│        │     │  Text  │     │  [OK]  │
└────────┘     └────────┘     └────────┘
 (扁平)       (斜角)        (圆角)
```

**动态效果**：
| 效果 | 描述 | 实现 |
|------|------|------|
| 落子动画 | 棋子从上方落下 | 逐帧 y 轴偏移 |
| 吃子动画 | 棋子缩小消失 | 缩放 + 透明度 |
| 光标闪烁 | 1Hz 频率闪烁 | 交替显示 |
| 胜利特效 | 粒子爆炸庆祝 | 粒子系统 |

---

## 五、界面设计

### 4.1 GBA 像素风格

- 棋盘：240x160 分辨率（GBA 原生分辨率）
- 像素大小：2x2 或 3x3 像素模拟
- 配色方案：
  - 棋盘背景：#e8cda5（竹色）
  - 黑子：#1a1a1a
  - 白子：#f5f5f5
  - 星位：#c9a05c

### 4.2 布局

```
+----------------------------------+
|  VimGo - 围棋对战                 |
+----------------------------------+
|                                  |
|     ┌──────────────────────┐     |
|     │    ●    ┌───┐        │     │
|     │  ┌───┐  │ ● │  ┌───┐  │     |
|     │  │ ● │  └───┘  │ ● │  │     |
|     │  └───┘        └───┘  │     |
|     │    ●                   │     |
|     └──────────────────────┘     |
|                                  |
+----------------------------------+
| Normal | 19x19 | 黑方 | 第 1 手 | A1|
+----------------------------------+
```

---

## 六、开发计划

### 第一阶段：核心引擎
- [ ] 棋盘数据结构设计
- [ ] 落子逻辑实现
- [ ] 气数计算
- [ ] 提子判定
- [ ] 禁入点检测
- [ ] 劫争处理

### 第二阶段：Vim 交互
- [ ] Normal 模式实现
- [ ] 光标移动 (hjkl, w/e/b, gg/G)
- [ ] 落子操作
- [ ] Command 模式实现
- [ ] 命令解析 (:q, :w, :undo, :pass)
- [ ] 状态栏更新

### 第三阶段：UI 开发

#### 3.1 终端 UI
- [ ] TUI 框架搭建（tview/bubbletea）
- [ ] 字符棋盘渲染
- [ ] 颜色主题适配
- [ ] 鼠标支持（可选）
- [ ] chafa 集成（可选）

#### 3.2 桌面 UI（Ebiten）
- [ ] Ebiten 项目搭建
- [ ] 像素素材制作（240x160）
- [ ] 棋盘绘制模块
- [ ] 棋子精灵图渲染
- [ ] 像素字体集成
- [ ] 落子/吃子动画
- [ ] 状态栏 UI

#### 3.3 网页 UI（Phaser.js）
- [ ] Phaser.js 项目搭建
- [ ] 像素配置（pixelArt: true）
- [ ] Canvas 棋盘渲染
- [ ] 键盘事件绑定
- [ ] 扫描线效果
- [ ] 响应式布局
- [ ] 移动端适配

### 第四阶段：Web 部署
- [ ] WebAssembly 编译配置
- [ ] 前后端联调
- [ ] 静态资源打包
- [ ] Web 服务器搭建
- [ ] 部署到云服务

### 第五阶段：高级功能
- [ ] SGF 棋谱读写
- [ ] AI 对战（简单算法）
- [ ] 悔棋/重做历史
- [ ] 多人联机对战
- [ ] 棋局复盘

---

## 七、命令行工具

### 6.1 安装与使用

```bash
# 安装
go install github.com/vimgo/vimgo@latest

# 运行游戏
vimgo

# 指定棋盘尺寸
vimgo --size 19    # 19路
vimgo --size 13    # 13路
vimgo --size 9     # 9路

# 加载棋谱
vimgo --file game.sgf

# 网页模式
vimgo --web        # 启动 Web 服务器
vimgo --web :8080  # 指定端口
```

### 6.2 命令行参数

| 参数 | 描述 |
|------|------|
| `-s, --size` | 棋盘尺寸 (9/13/19) |
| `-f, --file` | 加载 SGF 棋谱文件 |
| `-w, --web` | 启动网页模式 |
| `-p, --port` | Web 服务器端口 |
| `-h, --help` | 显示帮助 |
| `-v, --version` | 显示版本 |

---

## 八、部署方案

### 7.1 终端版本

```bash
# Homebrew (macOS)
brew install vimgo

# Arch Linux (AUR)
yay -S vimgo

# Go 安装
go install github.com/vimgo/vimgo/cmd/vimgo@latest

# 源码编译
git clone https://github.com/vimgo/vimgo.git
cd vimgo
make install
```

### 7.2 网页版本

**静态部署**
```bash
# 构建 WebAssembly 版本
make build-web

# 输出到 dist/ 目录
# 可部署到 GitHub Pages, Vercel, Netlify
```

**服务器部署**
```bash
# 启动 Web 服务器
vimgo --web --port 8080

# Docker 部署
docker run -p 8080:8080 vimgo/vimgo:latest
```

---

## 九、项目里程碑

| 里程碑 | 目标 | 预计时间 |
|--------|------|----------|
| M1 | 核心围棋引擎 | 2 周 |
| M2 | Vim 交互系统 | 2 周 |
| M3 | 终端 UI 完成 | 2 周 |
| M4 | 网页版本完成 | 3 周 |
| M5 | 首次发布 (v1.0) | 1 周 |

---

## 十、资源与参考

### 围棋规则参考
- 中国围棋规则（数子法）
- 日本围棋规则（点目法）
- 通用规则：打劫、禁入点

### 技术参考

**像素游戏框架**：
- Ebiten: https://github.com/hajimehoshi/ebiten
- Phaser.js: https://phaser.io/
- tview: https://github.com/rivo/tview
- bubbletea: https://github.com/charmbracelet/bubbletea
- gocui: https://github.com/jesseduffield/gocui

**像素工具**：
- chafa: https://github.com/hpjansson/chafa
- WebAssembly: https://github.com/golang/go/wiki/WebAssembly

**像素字体**：
- TinyPico: https://github.com/google/fonts/tree/main/ofl/tinypoco
- Pixelicroff: https://github.com/google/fonts/tree/main/ofl/pixelocoff

### SGF 格式
- SGF 标准: https://www.red-bean.com/sgf/

---

## 十一、贡献指南

欢迎参与 VimGo 项目的开发！

- 提交 Issue：报告 Bug 或提出建议
- 提交 PR：贡献代码
- 完善文档：帮助更多用户

**联系方式**
- GitHub: https://github.com/vimgo/vimgo

---

*文档版本: v1.1*
*最后更新: 2026-02-08*
