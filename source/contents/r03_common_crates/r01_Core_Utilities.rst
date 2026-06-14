=============================
基础工具 Crate
=============================

Rust 生态中最常用的基础工具类 crate，覆盖序列化、正则、时间、UUID、随机数、配置管理等日常开发需求。

.. contents:: 目录
   :depth: 3
   :local:

serde
=========

序列化/反序列化框架，是 Rust 生态中数据交换的事实标准。

核心概念：

- ``Serialize`` Trait：将 Rust 类型序列化为某种格式
- ``Deserialize`` Trait：从某种格式反序列化为 Rust 类型
- serde 本身不绑定具体格式，由后端 crate（如 ``serde_json``）负责

基本使用：

.. code-block:: rust

   use serde::{Serialize, Deserialize};

   #[derive(Serialize, Deserialize, Debug)]
   struct User {
       id: u64,
       name: String,
       email: String,
   }

serde 支持的属性：

.. list-table:: serde 常用属性
   :header-rows: 1
   :widths: 35 65

   * - 属性
     - 说明
   * - ``#[serde(rename = "xxx")]``
     - 重命名字段
   * - ``#[serde(rename_all = "camelCase")]``
     - 整体重命名风格（camelCase / snake_case / UPPERCASE 等）
   * - ``#[serde(skip)]``
     - 跳过该字段（不序列化/反序列化）
   * - ``#[serde(skip_serializing)]``
     - 序列化时跳过
   * - ``#[serde(skip_deserializing)]``
     - 反序列化时跳过
   * - ``#[serde(default)]``
     - 字段缺失时使用默认值
   * - ``#[serde(flatten)]``
     - 将嵌套结构展开到父级

属性示例：

.. code-block:: rust

   use serde::{Serialize, Deserialize};

   #[derive(Serialize, Deserialize, Debug)]
   #[serde(rename_all = "camelCase")]
   struct User {
       user_id: u64,          // JSON: "userId"
       #[serde(rename = "username")]
       name: String,          // JSON: "username"
       #[serde(skip)]
       internal_id: u64,      // 不参与序列化
       #[serde(default)]
       age: u8,               // 缺失时默认为 0
   }

serde_json
-------------

serde 的 JSON 后端实现：

.. code-block:: rust

   use serde::{Serialize, Deserialize};
   use serde_json;

   #[derive(Serialize, Deserialize, Debug, PartialEq)]
   struct User {
       id: u64,
       name: String,
   }

   fn main() -> Result<(), serde_json::Error> {
       // 序列化
       let user = User { id: 1, name: String::from("Alice") };
       let json = serde_json::to_string(&user)?;
       println!("序列化: {}", json);
       // {"id":1,"name":"Alice"}

       // 美化输出
       let pretty = serde_json::to_string_pretty(&user)?;
       println!("美化:\n{}", pretty);

       // 反序列化
       let parsed: User = serde_json::from_str(&json)?;
       assert_eq!(parsed, user);

       // 动态 JSON（serde_json::Value）
       let v: serde_json::Value = serde_json::from_str(r#"{"key": "value"}"#)?;
       println!("key = {}", v["key"]);

       Ok(())
   }

常用 API：

.. list-table:: serde_json 常用 API
   :header-rows: 1
   :widths: 35 65

   * - 方法
     - 说明
   * - ``serde_json::to_string(&v)``
     - 序列化为 JSON 字符串
   * - ``serde_json::to_string_pretty(&v)``
     - 美化输出
   * - ``serde_json::from_str(s)``
     - 从字符串反序列化
   * - ``serde_json::from_reader(reader)``
     - 从 Reader 反序列化
   * - ``serde_json::to_writer(writer, &v)``
     - 序列化到 Writer
   * - ``serde_json::Value``
     - 动态 JSON 值（不预先定义结构体）

regex
========

正则表达式库，编译时检查正则语法。

.. code-block:: rust

   use regex::Regex;

   fn main() {
       // 编译正则（失败时 panic）
       let re = Regex::new(r"^\d{4}-\d{2}-\d{2}$").unwrap();

       assert!(re.is_match("2024-01-15"));
       assert!(!re.is_match("2024-1-15"));

       // 捕获分组
       let re = Regex::new(r"(\w+)@(\w+)\.(\w+)").unwrap();
       let caps = re.captures("hello@example.com").unwrap();
       println!("完整: {}", &caps[0]);   // hello@example.com
       println!("用户: {}", &caps[1]);   // hello
       println!("域名: {}", &caps[2]);   // example
       println!("后缀: {}", &caps[3]);   // com

       // 命名分组
       let re = Regex::new(r"(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})").unwrap();
       let caps = re.captures("2024-03-15").unwrap();
       println!("年: {}, 月: {}, 日: {}", &caps["year"], &caps["month"], &caps["day"]);

       // 替换
       let re = Regex::new(r"\s+").unwrap();
       let result = re.replace_all("a   b  c", "_");
       println!("{}", result); // a_b_c

       // 查找所有匹配
       let re = Regex::new(r"\d+").unwrap();
       let matches: Vec<&str> = re.find_iter("a1 b22 c333").map(|m| m.as_str()).collect();
       println!("{:?}", matches); // ["1", "22", "333"]
   }

常用方法：

.. list-table:: regex 常用方法
   :header-rows: 1
   :widths: 35 65

   * - 方法
     - 说明
   * - ``re.is_match(text)``
     - 是否匹配
   * - ``re.find(text)``
     - 查找第一个匹配
   * - ``re.captures(text)``
     - 捕获分组
   * - ``re.replace_all(text, replacement)``
     - 替换所有匹配
   * - ``re.find_iter(text)``
     - 迭代所有匹配
   * - ``Regex::new(pattern)``
     - 编译正则表达式

chrono
=========

日期和时间处理的标准库。

.. code-block:: rust

   use chrono::{Local, Utc, NaiveDate, NaiveDateTime, Duration, Datelike, Timelike};

   fn main() {
       // 当前时间
       let now = Utc::now();
       println!("UTC 时间: {}", now);
       println!("本地时间: {}", Local::now());

       // 格式化输出
       println!("格式化: {}", now.format("%Y-%m-%d %H:%M:%S"));

       // 从字符串解析
       let dt = NaiveDate::parse_from_str("2024-03-15", "%Y-%m-%d")
           .unwrap()
           .and_hms_opt(12, 0, 0)
           .unwrap();
       println!("解析: {}", dt);

       // 日期运算
       let today = Utc::now().date_naive();
       let tomorrow = today + Duration::days(1);
       let last_week = today - Duration::weeks(1);
       println!("明天: {}, 上周: {}", tomorrow, last_week);

       // 日期属性
       println!("年: {}, 月: {}, 日: {}", today.year(), today.month(), today.day());

       // 时间戳
       let timestamp = now.timestamp();
       println!("时间戳: {}", timestamp);
   }

uuid
=========

生成和解析 UUID（Universally Unique Identifier）。

.. code-block:: rust

   use uuid::Uuid;

   fn main() {
       // 生成 v4 UUID（随机）
       let id = Uuid::new_v4();
       println!("v4 UUID: {}", id);

       // 简单格式（无连字符）
       println!("简单格式: {}", id.simple());

       // 从字符串解析
       let parsed: Uuid = "550e8400-e29b-41d4-a716-446655440000".parse().unwrap();
       println!("解析: {}", parsed);

       // 生成 nil UUID
       let nil = Uuid::nil();
       println!("nil UUID: {}", nil);
   }

rand
==========

随机数生成库。

.. code-block:: rust

   use rand::prelude::*;
   use rand::distributions::{Alphanumeric, Uniform};

   fn main() {
       let mut rng = rand::thread_rng();

       // 随机整数
       let n: i32 = rng.gen();
       println!("随机 i32: {}", n);

       // 指定范围的随机数
       let n: u32 = rng.gen_range(1..=100);
       println!("1-100 随机数: {}", n);

       // 随机布尔
       let b: bool = rng.gen();
       println!("随机布尔: {}", b);

       // 均匀分布
       let between = Uniform::from(10..=20);
       let nums: Vec<i32> = (0..5).map(|_| rng.sample(&between)).collect();
       println!("均匀分布: {:?}", nums);

       // 从切片随机选择
       let choices = ["苹果", "香蕉", "橘子"];
       let pick = choices.choose(&mut rng).unwrap();
       println!("随机选择: {}", pick);

       // 随机字符串
       let random_string: String = (&mut rng)
           .sample_iter(&Alphanumeric)
           .take(8)
           .map(char::from)
           .collect();
       println!("随机字符串: {}", random_string);

       // 打乱顺序
       let mut vec = vec![1, 2, 3, 4, 5];
       vec.shuffle(&mut rng);
       println!("打乱后: {:?}", vec);
   }

常用 API：

.. list-table:: rand 常用功能
   :header-rows: 1
   :widths: 30 70

   * - 功能
     - 示例
   * - 随机整数
     - ``rng.gen::<i32>()``
   * - 范围随机数
     - ``rng.gen_range(1..=100)``
   * - 随机布尔
     - ``rng.gen::<bool>()``
   * - 随机浮点
     - ``rng.gen::<f64>()``
   * - 均匀分布
     - ``Uniform::from(1..=10)`` + ``rng.sample()``
   * - 切片随机选择
     - ``slice.choose(&mut rng)``
   * - 随机字符串
     - ``Alphanumeric`` 分布
   * - 打乱顺序
     - ``vec.shuffle(&mut rng)``

config + dotenvy
=================

config 处理配置文件，dotenvy 加载 ``.env`` 环境变量。

创建项目
-----------------

.. code-block:: console

  $ cargo new config_demo --vcs none
  $ cd config_demo
  $ cargo add config dotenvy
  $ cargo add serde -F derive

项目结构
---------------

.. code-block:: text

  src/
  ├─ config/
  │   ├─ mod.rs
  │   └─ application.toml
  ├─ main.rs
  .env

.env 文件
--------------------

.. code-block:: text
  :caption: .env

  APP_PORT=8080
  DATABASE_URL=mysql://root:123456@localhost/demo

application.toml
------------------------

.. code-block:: toml
  :caption: application.toml

  app_name = "hello-rust"

  [server]
  host = "0.0.0.0"
  port = 8080

配置结构体
-------------------

.. code-block:: rust
  :caption: src/config/mod.rs

  use serde::Deserialize;
  use config::{Config, File, Environment};
  use dotenvy::dotenv;

  #[derive(Debug, Deserialize)]
  pub struct Settings {
      pub app_name: String,
      pub server: ServerConfig,
      pub database_url: String,
  }

  #[derive(Debug, Deserialize)]
  pub struct ServerConfig {
      pub host: String,
      pub port: u16,
  }


  pub fn load_config() -> Result<Settings, config::ConfigError> {
      dotenv().ok();

      let settings = Config::builder()
            .set_default("database_url", "postgres://localhost:5432/mydb")
            .expect("Failed to set default value")
            .add_source(File::with_name("src/config/application"))
            .add_source(Environment::default())
            .build()
            .expect("Failed to load configuration");

      settings.try_deserialize()
  }

main.rs
----------------

.. code-block:: rust
  :caption: src/main.rs

  mod config;

  use crate::config::load_config;

  fn main() {
      let settings = load_config().unwrap();

      println!("{:#?}", settings);
  }

运行结果
--------------------

.. code-block:: rust

  Settings {
    app_name: "hello-rust",
    server: ServerConfig {
        host: "0.0.0.0",
        port: 8080,
    },
    database_url: "mysql://root:123456@localhost/demo",
  }

Environment 覆盖规则
------------------------

``Environment::default()`` 会自动读取环境变量：

.. code-block:: text

  APP_NAME=xxx
  SERVER__PORT=9090

双下划线 ``__`` 表示嵌套层级。

配置加载优先级（从低到高）：

.. list-table:: 配置优先级
   :header-rows: 1
   :widths: 20 40 40

   * - 优先级
     - 来源
     - 说明
   * - 1（最低）
     - ``set_default``
     - 代码中设置默认值
   * - 2
     - ``File``
     - 配置文件（TOML / YAML / JSON）
   * - 3（最高）
     - ``Environment``
     - 环境变量覆盖

推荐完整写法（生产）
----------------------------

.. code-block:: rust

  let settings = Config::builder()
      .add_source(File::with_name("config/default"))
      .add_source(File::with_name("config/local").required(false))
      .add_source(Environment::default().separator("__"))
      .build()?;

推荐目录结构
------------------------

.. code-block:: text

  config/
  ├─ default.toml
  ├─ dev.toml
  ├─ test.toml
  ├─ prod.toml

按环境加载
-------------------

.. code-block:: rust

  let profile = std::env::var("APP_ENV")
      .unwrap_or_else(|_| "dev".into());

  let settings = Config::builder()
      .add_source(File::with_name("config/default"))
      .add_source(File::with_name(&format!("config/{}", profile)))
      .build()?;

总结
=====

.. list-table:: 基础工具 Crate 总览
   :header-rows: 1
   :widths: 25 25 50

   * - Crate
     - 用途
     - 典型场景
   * - ``serde``
     - 序列化/反序列化框架
     - API 数据交换、配置文件解析
   * - ``serde_json``
     - JSON 支持
     - REST API、前端通信
   * - ``regex``
     - 正则表达式
     - 输入验证、文本解析、日志分析
   * - ``chrono``
     - 日期和时间
     - 时间戳处理、日志时间、定时任务
   * - ``uuid``
     - UUID 生成
     - 数据库主键、请求追踪 ID
   * - ``rand``
     - 随机数
     - 测试数据生成、游戏、加密
   * - ``config``
     - 配置管理
     - 多环境配置、层级覆盖
   * - ``dotenvy``
     - 环境变量加载
     - 本地开发、敏感配置分离
