# 鹦鹉代理命名系统调研报告

> **调研日期**: 2026-01-31
> **状态**: 待审核

---

## 执行摘要

DivineSense 的「鹦鹉」代理系统需要一个符合真实鹦鹉品种特性的命名系统。当前命名存在一些不匹配（如"金刚"对应 Macaw 与日程管理能力的关联度较低），且未来需要扩展新的代理（如邮件代理）。

---

## 一、当前命名分析

### 1.1 现有代理映射

| AgentType | 中文名 | 英文名 | 鹦鹉品种 | 实际能力 | 匹配度 |
|:---------|:------|:-------|:--------|:--------|:------:|
| `MEMO` | 灰灰 | Grey | 非洲灰鹦鹉 (African Grey) | 语义搜索、记忆检索 | ✅ 高 |
| `SCHEDULE` | 金刚 | King Kong | 金刚鹦鹉 (Macaw) | 日程管理、时间规划 | ⚠️ 中 |
| `AMAZING` | 惊奇 | Amazing | 亚马逊鹦鹉 (Amazon Parrot) | 综合助理、并发执行 | ⚠️ 中 |
| `GEEK` | 极客 | Geek | - | Claude Code CLI 集成 | - |
| `EVOLUTION` | 进化 | Evolution | - | 源代码修改 + PR | - |

### 1.2 品种特性分析

#### 非洲灰鹦鹉 (African Grey Parrot) ✅

**真实特性**：
- **记忆力惊人**：著名的非洲灰鹦鹉 Alex 能识别 100+ 物品、理解概念（颜色、形状、数量、相同/不同）
- **语言能力**：可记住数千个词汇，理解语义而非简单模仿
- **问题解决**：能完成复杂的认知任务

**与 MEMO 代理匹配度**：🟢 完美匹配
- 鹦鹉的卓越记忆力 → 语义搜索能力
- 词汇理解 → 笔记内容理解

#### 金刚鹦鹉 (Macaw) ⚠️

**真实特性**：
- **体型巨大**：最大的鹦鹉之一，体长可达 1 米
- **喙部强壮**：能咬开坚硬的坚果
- **寿命极长**：可活 50+ 年
- **群居社交**：以家庭为单位活动

**与 SCHEDULE 代理匹配度**：🟡 部分匹配
- 喙部力量 → 无法类比时间管理
- 长寿命 → 可隐喻长期规划
- **问题**：金刚鹦鹉与"时间管理"、"日程安排"的关联度不够直观

#### 亚马逊鹦鹉 (Amazon Parrot) ⚠️

**真实特性**：
- **语言天赋**：世界上最会说话的鹦鹉之一
- **性格活泼**：好奇心强，喜欢探索
- **社交能力**：善于与人类互动

**与 AMAZING 代理匹配度**：🟡 部分匹配
- 语言能力 → 可以类比综合表达能力
- **问题**："亚马逊"是一个属名（包含 30+ 种），而非单一品种；命名过于通用

---

## 二、品种能力映射研究

### 2.1 日程管理理想的鹦鹉品种

研究显示，以下品种更适合日程管理代理：

| 品种 | 特性 | 匹配理由 |
|:-----|:-----|:--------|
| **太平洋鹦鹉 (Pacific Parrotlet)** | 体型小但精力充沛 | 高效执行 |
| **虎皮鹦鹉 (Budgerigar)** | 极其规律的生活作息 | 时间规律性 |
| **吸蜜鹦鹉 (Lory)** | 活跃、喜欢穿梭 | 灵活调度 |
| **鸡尾鹦鹉 (Cockatiel)** | 精准的时间感知能力 | 生物钟准确 |

### 2.2 综合助理理想的鹦鹉品种

| 品种 | 特性 | 匹配理由 |
|:-----|:-----|:--------|
| **折衷鹦鹉 (Eclectus Parrot)** | 两性差异显著、高度社会化 | 多面能力 |
| **灰头情侣鹦鹉 (Agapornis)** | 合作能力强 | 并发协作 |
| **葵花凤头鹦鹉 (Sulphur-crested Cockatoo)** | 智力极高、工具使用 | 综合问题解决 |

---

## 三、建议的命名系统

### 3.1 命名原则

1. **品种-能力一致性**：鹦鹉品种的自然能力应与代理功能匹配
2. **中英文对称**：中文名应是英文名的意译或音译
3. **扩展性**：为未来新代理预留品种空间
4. **文化适配**：中文名易记且富有含义

### 3.2 方案 A：品种重映射（推荐）

| AgentType | 新中文 | 新英文 | 鹦鹉品种 | 变更说明 |
|:---------|:------|:-------|:--------|:--------|
| `MEMO` | 灰灰 | Grey | 非洲灰鹦鹉 | ✅ 保持不变 |
| `SCHEDULE` | 时巧 | Tick | 鸡尾鹦鹉 | 🔄 新品种 |
| `AMAZING` | 折衷 | Echo | 折衷鹦鹉 | 🔄 新品种 |
| `GEEK` | 极客 | Geek | - | ✅ 保持不变 |
| `EVOLUTION` | 进化 | Evolution | - | ✅ 保持不变 |

**鸡尾鹦鹉 (Cockatiel) → 时巧 (Tick)**
- 真实特性：精准的生物钟、规律作息
- 中文"时巧"：时间上的机巧，寓意时间管理专家
- 英文"Tick"：时钟的滴答声

**折衷鹦鹉 (Eclectus Parrot) → 折衷 (Echo)**
- 真实特性：两性差异显著（雄翠红）、高智商
- 中文"折衷"：谐音"折衷"，寓意综合多方能力
- 英文"Echo"：回声，寓意能回应各种需求

### 3.3 方案 B：保持现有命名（最小改动）

保持灰灰、金刚、惊奇不变，但增强品种特性的解释：

- **金刚**：强调"长寿命"→"长期规划能力"
- **惊奇**：强调"亚马逊鹦鹉的社交性"→"综合协调能力"

### 3.4 未来扩展品种

| 计划代理 | 候选品种 | 中文名 | 英文名 | 能力映射 |
|:--------|:--------|:------|:-------|:--------|
| **邮件** | 虹彩吸蜜鹦鹉 (Rainbow Lorikeet) | 彩虹 | Rainbow | 多彩表达 |
| **提醒** | 灰胸吸蜜鹦鹉 (Dusky Lory) | 唤唤 | Nudge | 持续提醒 |
| **归档** | 深林王鹦鹉 (King Parrot) | 典藏 | Archive | 整理归类 |
| **搜索** | 超级鹦鹉 (Superb Parrot) | 卓越 | Scout | 敏锐搜索 |
| **统计** | 红冠灰鹦鹉 (Galah) | 统计 | Chart | 数据可视化 |

---

## 四、命名系统技术规范

### 4.1 代码映射

```go
// plugin/ai/agent/types.go

// ParrotSpecies represents the actual parrot species for each agent.
type ParrotSpecies struct {
    CommonName   string   // 英文俗名
    ScientificName string // 学名
    ChineseName  string   // 中文名
    Emoji         string   // 表情符号
}

var AgentSpeciesMap = map[ParrotAgentType]ParrotSpecies{
    MEMO: {
        CommonName: "African Grey Parrot",
        ScientificName: "Psittacus erithacus",
        ChineseName: "非洲灰鹦鹉",
        Emoji: "🦜",
    },
    SCHEDULE: {
        CommonName: "Cockatiel",           // 方案A
        // CommonName: "Macaw",            // 方案B
        ScientificName: "Nymphicus hollandicus",
        ChineseName: "鸡尾鹦鹉",
        Emoji: "🦜",
    },
    AMAZING: {
        CommonName: "Eclectus Parrot",     // 方案A
        // CommonName: "Amazon Parrot",    // 方案B
        ScientificName: "Eclectus roratus",
        ChineseName: "折衷鹦鹉",
        Emoji: "🦜",
    },
}
```

### 4.2 前端映射

```typescript
// web/src/types/parrot.ts

export interface ParrotSpecies {
  commonName: string;
  scientificName: string;
  chineseName: string;
  emoji: string;
  origin: string;
  naturalAbilities: string[];
  symbolicMeaning: string;
}

export const PARROT_SPECIES: Record<ParrotAgentType, ParrotSpecies> = {
  [ParrotAgentType.MEMO]: {
    commonName: "African Grey Parrot",
    scientificName: "Psittacus erithacus",
    chineseName: "非洲灰鹦鹉",
    emoji: "🦜",
    origin: "非洲热带雨林",
    naturalAbilities: ["惊人的记忆力", "强大的模仿能力", "复杂问题解决"],
    symbolicMeaning: "智慧与记忆的象征",
  },
  // ...
};
```

---

## 五、实施计划

| 阶段 | 任务 | 工作量 | 优先级 |
|:----:|:-----|:-------|:-------|
| Phase 1 | 确定命名方案 | 0.5 人日 | P0 |
| Phase 2 | 更新后端 SelfDescribe | 1 人日 | P0 |
| Phase 3 | 更新前端 PARROT_AGENTS | 1 人日 | P0 |
| Phase 4 | 更新文档 | 0.5 人日 | P1 |
| Phase 5 | 添加品种科普内容 | 1 人日 | P2 |

**总计**: 3-4 人日

---

## 六、参考资源

### 鹦鹉品种权威资源

- [World Parrot Trust](https://www.parrots.org/) - 全球鹦鹉保护组织
- [BirdLife International](https://www.birdlife.org/) - 鸟类物种数据库
- [IUCN Red List](https://www.iucnredlist.org/) - 濒危物种名录

### 品种特性参考

| 品种 | 智力排名 | 语言能力 | 寿命 | 保护状态 |
|:-----|:--------|:--------|:-----|:--------|
| 非洲灰鹦鹉 | #1 | ⭐⭐⭐⭐⭐ | 40-60年 | EN |
| 金刚鹦鹉 | #3 | ⭐⭐⭐ | 50年 | LC/NT |
| 折衷鹦鹉 | #2 | ⭐⭐⭐⭐ | 30年 | LC |
| 鸡尾鹦鹉 | #5 | ⭐⭐⭐ | 15-25年 | LC |
| 吸蜜鹦鹉 | #4 | ⭐⭐⭐ | 20年 | VU |

---

## 七、附录

### A. 品种详细信息

#### 非洲灰鹦鹉 (Psittacus erithacus)
- **别名**：灰鹦鹉、刚果灰
- **分布**：西非至中非热带雨林
- **著名个体**：Alex（Irene Pepperberg 研究对象）
- **认知能力**：理解零的概念、相同/不同、基本数学

#### 鸡尾鹦鹉 (Nymphicus hollandicus)
- **别名**：玄凤、面拨凤凰
- **分布**：澳大利亚内陆
- **特点**：头上有可移动的冠羽
- **能力**：精准的时间感知，能预测日常事件

#### 折衷鹦鹉 (Eclectus roratus)
- **别名**：极乐鸟鹦鹉
- **分布**：新几内亚、澳大利亚北部
- **特点**：极度两性异形（雄绿、雌红）
- **能力**：复杂的社交结构，高智商

---

*报告版本: v1.0 | 最后更新: 2026-01-31*
