============================================================
系统 & 命令行 Crate
============================================================

Rust 生态中构建 CLI 工具、错误处理和日志系统的主流 crate。

.. contents:: 目录
   :depth: 3
   :local:

clap
==========

命令行参数解析库，支持 derive（推荐）和 builder 两种模式。``structopt`` 已合并进 clap 3.x+。

.. code-block:: toml

   [dependencies]
   clap = { version = "4", features = ["derive"] }

derive 模式（推荐）：

.. code-block:: rust

   use clap::Parser;

   #[derive(Parser, Debug)]
   #[command(name = "myapp")]
   #[command(version = "1.0")]
   #[command(about = "一个示例 CLI 工具", long_about = None)]
   struct Cli {
       /// 输入文件路径
       #[arg(short, long)]
       input: String,

       /// 输出文件路径
       #[arg(short, long, default_value = "output.txt")]
       output: String,

       /// 启用详细输出
       #[arg(short, long, default_value_t = false)]
       verbose: bool,

       /// 线程数
       #[arg(short, long, default_value_t = 4)]
       threads: u32,
   }

   fn main() {
       let cli = Cli::parse();

       println!("输入: {}", cli.input);
       println!("输出: {}", cli.output);
       println!("详细: {}", cli.verbose);
       println!("线程: {}", cli.threads);
   }

子命令：

.. code-block:: rust

   use clap::{Parser, Subcommand};

   #[derive(Parser)]
   #[command(name = "git-cli")]
   struct Cli {
       #[command(subcommand)]
       command: Commands,
   }

   #[derive(Subcommand)]
   enum Commands {
       /// 克隆仓库
       Clone {
           /// 仓库 URL
           url: String,
       },
       /// 添加文件
       Add {
           /// 要添加的文件
           #[arg(short, long)]
           file: String,
       },
       /// 提交更改
       Commit {
           /// 提交信息
           #[arg(short, long)]
           message: String,
       },
   }

   fn main() {
       let cli = Cli::parse();

       match cli.command {
           Commands::Clone { url } => println!("克隆: {}", url),
           Commands::Add { file } => println!("添加: {}", file),
           Commands::Commit { message } => println!("提交: {}", message),
       }
   }

验证与参数约束：

.. code-block:: rust

   use clap::Parser;

   #[derive(Parser)]
   struct Cli {
       /// 端口号 (1-65535)
       #[arg(short, long, value_parser = clap::value_parser!(u16).range(1..=65535))]
       port: u16,

       /// 模式
       #[arg(short, long, value_parser = ["dev", "test", "prod"])]
       mode: String,

       /// 必须二选一的参数
       #[arg(long, group = "action")]
       create: bool,

       #[arg(long, group = "action")]
       delete: bool,
   }

常用属性：

.. list-table:: clap derive 常用属性
   :header-rows: 1
   :widths: 35 65

   * - 属性
     - 说明
   * - ``#[arg(short, long)]``
     - 同时支持 ``-x`` 和 ``--xxx``
   * - ``#[arg(default_value = "...")]``
     - 默认值
   * - ``#[arg(default_value_t = 0)]``
     - 类型安全的默认值
   * - ``#[arg(value_parser = ...)]``
     - 自定义值验证
   * - ``#[arg(required = true)]``
     - 必填参数
   * - ``#[arg(conflicts_with = "other")]``
     - 与其他参数互斥
   * - ``#[arg(requires = "other")]``
     - 依赖其他参数
   * - ``#[arg(group = "name")]``
     - 参数组
   * - ``#[command(subcommand)]``
     - 声明子命令

anyhow
==========

灵活的应用程序级错误处理，适合业务逻辑层。自动为错误添加上下文，无需自定义错误类型。

.. code-block:: rust

   use anyhow::{Context, Result, anyhow, bail};

   fn read_config(path: &str) -> Result<String> {
       let content = std::fs::read_to_string(path)
           .with_context(|| format!("无法读取配置文件: {}", path))?;
       Ok(content)
   }

   fn validate_config(content: &str) -> Result<()> {
       if content.trim().is_empty() {
           bail!("配置文件为空");
       }
       if !content.contains("server") {
           return Err(anyhow!("配置文件缺少 [server] 节"));
       }
       Ok(())
   }

   fn main() -> Result<()> {
       let content = read_config("config.toml")?;
       validate_config(&content)?;
       println!("配置加载成功");
       Ok(())
   }

   // 在 main 中使用 anyhow::Result，? 操作符自动转换所有错误

常用 API：

.. list-table:: anyhow 常用 API
   :header-rows: 1
   :widths: 35 65

   * - 方法
     - 说明
   * - ``anyhow::Result<T>``
     - 通用的 Result 类型别名
   * - ``anyhow!("message")``
     - 创建错误
   * - ``bail!("message")``
     - 创建错误并 return
   * - ``.context("message")``
     - 为错误添加上下文
   * - ``.with_context(|| format!(...))``
     - 延迟计算上下文
   * - ``.map_err(|e| anyhow!(...))``
     - 手动转换错误

anyhow vs thiserror：

.. list-table:: anyhow vs thiserror
   :header-rows: 1
   :widths: 25 40 35

   * - 特性
     - ``anyhow``
     - ``thiserror``
   * - 使用场景
     - 应用程序（调用方）
     - 库（定义错误类型）
   * - 错误类型
     - 抹除具体类型
     - 保留具体枚举类型
   * - 调用者匹配
     - 不能 match 具体错误
     - 可以 match 错误变体
   * - 上手难度
     - 极低
     - 低
   * - 推荐
     - CLI 工具、应用层
     - 库、需要类型化错误

thiserror
==========

自定义错误类型的派生宏，适合库作者。编译期生成 ``Display`` 和 ``Error`` 实现。

.. code-block:: rust

   use thiserror::Error;

   #[derive(Error, Debug)]
   pub enum MyError {
       #[error("IO 错误: {0}")]
       Io(#[from] std::io::Error),

       #[error("解析错误: {0}")]
       Parse(#[from] std::num::ParseIntError),

       #[error("配置无效: {message}")]
       Config { message: String },

       #[error("未找到用户: id={user_id}")]
       NotFound { user_id: u64 },

       #[error("超时: {0}ms")]
       Timeout(u64),

       #[error("未知错误")]
       Unknown,
   }

   fn parse_user_id(input: &str) -> Result<u64, MyError> {
       let id: u64 = input.parse()?; // ParseIntError 自动转为 MyError::Parse
       if id == 0 {
           return Err(MyError::NotFound { user_id: 0 });
       }
       Ok(id)
   }

   fn load_file(path: &str) -> Result<String, MyError> {
       let content = std::fs::read_to_string(path)?; // io::Error 自动转为 MyError::Io
       Ok(content)
   }

   fn main() {
       match parse_user_id("0") {
           Ok(id) => println!("用户 ID: {}", id),
           Err(e) => println!("错误: {}", e), // 未找到用户: id=0
       }
   }

常用属性：

.. list-table:: thiserror 常用属性
   :header-rows: 1
   :widths: 35 65

   * - 属性
     - 说明
   * - ``#[error("...")]``
     - 定义 Display 格式
   * - ``#[from]``
     - 自动实现 From，支持 ``?`` 传播
   * - ``#[error(transparent)]``
     - 透传源错误的 Display 和 source
   * - ``{0}`` / ``{field}``
     - 在错误消息中引用字段
   * - ``#[source]``
     - 标记错误来源（用于 ``Error::source()``）

log + env_logger
==================

``log`` 是 Rust 的日志门面（Facade），定义日志宏；``env_logger`` 是常用后端实现。

.. code-block:: toml

   [dependencies]
   log = "0.4"
   env_logger = "0.11"

基本使用：

.. code-block:: rust

   use log::{info, warn, error, debug, trace};

   fn main() {
       // 初始化 logger（通常在 main 开头）
       env_logger::init();

       trace!("跟踪级别日志");
       debug!("调试信息");
       info!("服务启动成功，监听端口: 8080");
       warn!("磁盘空间不足: 剩余 5%");
       error!("数据库连接失败");
   }

通过环境变量控制日志级别：

.. code-block:: console

   $ RUST_LOG=debug cargo run
   $ RUST_LOG=my_app=debug,other_crate=warn cargo run
   $ RUST_LOG=info cargo run

模块级日志过滤：

.. code-block:: rust

   use log::{info, debug};

   mod database {
       pub fn connect() {
           debug!("尝试连接数据库...");  // RUST_LOG=database=debug 才会显示
       }
   }

   mod server {
       pub fn start() {
           info!("服务器启动");  // 默认 info 级别显示
       }
   }

   fn main() {
       env_logger::Builder::from_env(env_logger::Env::default().default_filter_or("info"))
           .init();

       database::connect();
       server::start();
   }

自定义格式：

.. code-block:: rust

   use env_logger::Builder;
   use std::io::Write;

   fn main() {
       Builder::new()
           .format(|buf, record| {
               writeln!(
                   buf,
                   "{} [{}] {} - {}",
                   chrono::Local::now().format("%Y-%m-%d %H:%M:%S"),
                   record.level(),
                   record.target(),
                   record.args()
               )
           })
           .filter(None, log::LevelFilter::Info)
           .init();

       log::info!("自定义格式的日志");
   }

tracing
==========

结构化、异步友好的日志与诊断框架。比 ``log`` 更强大，支持 spans（操作追踪）、结构化字段和异步上下文。

.. code-block:: toml

   [dependencies]
   tracing = "0.1"
   tracing-subscriber = "0.3"

基本使用：

.. code-block:: rust

   use tracing::{info, warn, error, debug, span, Level};

   fn main() {
       // 初始化 subscriber
       tracing_subscriber::fmt::init();

       info!("服务启动");

       // 结构化字段
       info!(user_id = 42, action = "login", "用户登录");
       warn!(remaining_space = "5%", "磁盘空间不足");
       error!(error = %"connection refused", "数据库连接失败");

       // Span：追踪操作
       let span = span!(Level::INFO, "处理请求", request_id = "abc-123");
       let _guard = span.enter();

       info!("开始处理");
       process_data();
       info!("处理完成");
   }

   fn process_data() {
       info!("数据处理中...");
   }

输出示例：

.. code-block:: text

   2024-01-15T10:30:00.000Z  INFO 服务启动
   2024-01-15T10:30:00.001Z  INFO 用户登录 user_id=42 action="login"
   2024-01-15T10:30:00.002Z  INFO 处理请求{request_id="abc-123"}: 开始处理
   2024-01-15T10:30:00.003Z  INFO 处理请求{request_id="abc-123"}: 数据处理中...
   2024-01-15T10:30:00.004Z  INFO 处理请求{request_id="abc-123"}: 处理完成

与 tokio 集成（异步上下文传播）：

.. code-block:: rust

   use tracing::{info, Instrument};

   #[tokio::main]
   async fn main() {
       tracing_subscriber::fmt::init();

       let handle = tokio::spawn(
           async {
               info!("异步任务中的日志");
               tokio::time::sleep(std::time::Duration::from_millis(100)).await;
               info!("异步任务完成");
           }
           .instrument(tracing::info_span!("background-task")),
       );

       handle.await.unwrap();
   }

Subscriber 配置：

.. code-block:: rust

   use tracing_subscriber::{fmt, EnvFilter};

   fn main() {
       fmt()
           .with_env_filter(EnvFilter::from_default_env()) // RUST_LOG
           .with_target(true)       // 显示模块路径
           .with_thread_ids(true)   // 显示线程 ID
           .with_line_number(true)  // 显示行号
           .json()                  // JSON 格式输出
           .init();
   }

log vs tracing：

.. list-table:: log vs tracing
   :header-rows: 1
   :widths: 25 40 35

   * - 特性
     - ``log`` + ``env_logger``
     - ``tracing``
   * - 结构化日志
     - 不支持
     - 原生支持（键值对）
   * - Span（操作追踪）
     - 不支持
     - 核心特性
   * - 异步友好
     - 一般
     - 原生支持上下文传播
   * - 生态兼容
     - 广泛，很多库使用
     - 越来越流行
   * - 性能
     - 轻量
     - 略重但功能更强
   * - 推荐
     - 简单项目、库
     - 异步项目、分布式系统

总结
=====

.. list-table:: 系统 & 命令行 Crate 总览
   :header-rows: 1
   :widths: 20 25 25 30

   * - Crate
     - 类型
     - 核心能力
     - 典型场景
   * - ``clap``
     - CLI 参数解析
     - derive 模式、子命令、验证
     - CLI 工具开发
   * - ``anyhow``
     - 应用级错误处理
     - 灵活错误、上下文、bail!
     - CLI 工具、应用程序
   * - ``thiserror``
     - 库级错误定义
     - 派生宏、From 自动转换
     - 库开发、类型化错误
   * - ``log`` + ``env_logger``
     - 日志门面 + 后端
     - 轻量日志、环境变量控制
     - 简单日志需求
   * - ``tracing``
     - 结构化诊断
     - Span、结构化字段、异步友好
     - 异步项目、分布式追踪
