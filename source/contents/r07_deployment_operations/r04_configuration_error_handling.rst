配置管理与错误处理 (Configuration & Error Handling)
=============================================================

.. .. module:: r07_deployment_operations.r04_configuration_error_handling

稳健的配置管理和错误处理是生产级应用的基础。
本章介绍配置加载方案、错误处理模式和优雅关闭策略。

config (分层配置)
-------------------

``config`` Crate 支持从多种来源加载配置，按优先级合并。

.. code-block:: rust

    use config::{Config, File, Environment};
    use serde::Deserialize;
    use secrecy::SecretString;

    #[derive(Debug, Deserialize, Clone)]
    struct AppConfig {
        server: ServerConfig,
        database: DatabaseConfig,
        redis: RedisConfig,
        logging: LoggingConfig,
    }

    #[derive(Debug, Deserialize, Clone)]
    struct ServerConfig {
        host: String,
        port: u16,
        workers: usize,
    }

    #[derive(Debug, Deserialize, Clone)]
    struct DatabaseConfig {
        url: SecretString,
        max_connections: u32,
        timeout_seconds: u64,
    }

    #[derive(Debug, Deserialize, Clone)]
    struct RedisConfig {
        url: String,
    }

    #[derive(Debug, Deserialize, Clone)]
    struct LoggingConfig {
        level: String,
        format: String,
    }

    fn load_config() -> Result<AppConfig, config::ConfigError> {
        let run_mode = std::env::var("RUN_MODE").unwrap_or_else(|_| "development".into());

        let config = Config::builder()
            // 1. 默认配置（最低优先级）
            .set_default("server.host", "0.0.0.0")?
            .set_default("server.port", 8080)?
            .set_default("server.workers", 4)?
            // 2. 基础配置文件
            .add_source(File::with_name("config/default").required(false))
            // 3. 环境特定配置文件（覆盖默认）
            .add_source(File::with_name(&format!("config/{}", run_mode)).required(false))
            // 4. 本地覆盖（不提交到 Git）
            .add_source(File::with_name("config/local").required(false))
            // 5. 环境变量（最高优先级），如 APP_DATABASE__URL
            .add_source(Environment::with_prefix("APP").separator("__"))
            .build()?;

        config.try_deserialize()
    }

.. code-block:: toml

    # config/default.toml
    [server]
    host = "0.0.0.0"
    port = 8080
    workers = 4

    [database]
    max_connections = 10
    timeout_seconds = 30

    [redis]
    url = "redis://localhost:6379"

    [logging]
    level = "info"
    format = "json"

    # config/production.toml
    [server]
    workers = 8

    [database]
    max_connections = 50

.. code-block:: toml

    # config/local.toml (加入 .gitignore)
    [database]
    url = "postgres://localhost:5432/myapp_dev"

    [logging]
    level = "debug"
    format = "pretty"

.. list-table:: config 配置来源优先级
   :header-rows: 1

   * - 优先级
     - 来源
     - 用途
   * - 1 (最低)
     - ``set_default``
     - 硬编码默认值
   * - 2
     - ``config/default.toml``
     - 通用默认配置
   * - 3
     - ``config/{env}.toml``
     - 环境特定配置
   * - 4
     - ``config/local.toml``
     - 开发者本地覆盖
   * - 5 (最高)
     - 环境变量 ``APP_*``
     - 运行时注入（密钥、动态配置）

dotenvy (.env 文件)
----------------------

``dotenvy`` 是 ``dotenv`` 的维护分支，用于从 ``.env`` 文件加载环境变量。

.. code-block:: rust

    use dotenvy;

    fn main() {
        // 加载 .env 文件（开发环境）
        let _ = dotenvy::dotenv();

        let database_url = std::env::var("DATABASE_URL")
            .expect("DATABASE_URL must be set");

        println!("Connecting to database...");
    }

.. code-block:: text

    # .env (不提交到 Git)
    DATABASE_URL=postgres://user:pass@localhost:5432/mydb
    REDIS_URL=redis://localhost:6379
    JWT_SECRET=super-secret-key-change-in-production
    RUST_LOG=debug

    # .env.example (提交到 Git 作为模板)
    DATABASE_URL=postgres://user:password@localhost:5432/mydb
    REDIS_URL=redis://localhost:6379
    JWT_SECRET=change-me
    RUST_LOG=info

.. warning::

   1. ``.env`` 文件必须加入 ``.gitignore``，防止凭据泄露。
   2. 生产环境应通过 Kubernetes Secrets / 云平台环境变量注入凭据。
   3. ``dotenvy::dotenv()`` 在已设置环境变量时不会覆盖，确保生产安全。

figment (声明式配置)
---------------------

``figment`` 提供更声明式的配置构建方式。

.. code-block:: rust

    use figment::{Figment, providers::{Env, Format, Toml, Serialized}};
    use serde::Deserialize;

    #[derive(Debug, Deserialize)]
    struct Config {
        port: u16,
        database_url: String,
    }

    fn load_config() -> Result<Config, figment::Error> {
        Figment::new()
            .merge(Serialized::defaults(Config {
                port: 8080,
                database_url: String::new(),
            }))
            .merge(Toml::file("config/default.toml"))
            .merge(Toml::file("config/local.toml"))
            .merge(Env::prefixed("APP_"))
            .extract()
    }

clap (CLI 参数解析)
---------------------

.. code-block:: rust

    use clap::{Parser, Subcommand, ValueEnum};

    #[derive(Parser)]
    #[command(name = "myapp")]
    #[command(version, about, long_about = None)]
    struct Cli {
        /// 配置文件路径
        #[arg(short, long, default_value = "config/default.toml")]
        config: String,

        /// 运行模式
        #[arg(short, long, value_enum, default_value_t = RunMode::Server)]
        mode: RunMode,

        /// 日志级别
        #[arg(short, long, default_value = "info")]
        log_level: String,

        #[command(subcommand)]
        command: Option<Commands>,
    }

    #[derive(ValueEnum, Clone, Debug)]
    enum RunMode {
        Server,
        Worker,
        Migration,
    }

    #[derive(Subcommand)]
    enum Commands {
        /// 运行数据库迁移
        Migrate {
            /// 迁移方向
            #[arg(short, long, default_value = "up")]
            direction: String,
        },
        /// 生成配置文件模板
        Init,
    }

    fn main() {
        let cli = Cli::parse();

        match cli.command {
            Some(Commands::Migrate { direction }) => {
                println!("Running migration: {}", direction);
            }
            Some(Commands::Init) => {
                println!("Generating config template...");
            }
            None => {
                println!("Starting server in {:?} mode", cli.mode);
            }
        }
    }

错误处理模式
---------------

**thiserror (库错误类型)：**

.. code-block:: rust

    use thiserror::Error;

    #[derive(Error, Debug)]
    pub enum AppError {
        #[error("Database error: {0}")]
        Database(#[from] sqlx::Error),

        #[error("IO error: {0}")]
        Io(#[from] std::io::Error),

        #[error("Validation error: {0}")]
        Validation(String),

        #[error("Not found: {0}")]
        NotFound(String),

        #[error("Unauthorized: {0}")]
        Unauthorized(String),

        #[error("Internal error: {0}")]
        Internal(#[from] anyhow::Error),
    }

    impl axum::response::IntoResponse for AppError {
        fn into_response(self) -> axum::response::Response {
            let (status, message) = match &self {
                AppError::NotFound(msg) => (StatusCode::NOT_FOUND, msg.clone()),
                AppError::Unauthorized(msg) => (StatusCode::UNAUTHORIZED, msg.clone()),
                AppError::Validation(msg) => (StatusCode::BAD_REQUEST, msg.clone()),
                _ => (StatusCode::INTERNAL_SERVER_ERROR, "Internal server error".into()),
            };

            (status, Json(serde_json::json!({ "error": message }))).into_response()
        }
    }

**anyhow (应用错误处理)：**

.. code-block:: rust

    use anyhow::{Context, Result, anyhow, bail};

    fn read_config(path: &str) -> Result<String> {
        let content = std::fs::read_to_string(path)
            .with_context(|| format!("Failed to read config file: {}", path))?;
        Ok(content)
    }

    fn validate_port(port: u16) -> Result<()> {
        if port < 1024 {
            bail!("Port {} is privileged, use port >= 1024", port);
        }
        Ok(())
    }

    fn main() -> Result<()> {
        let config = read_config("config.toml")?;
        validate_port(8080)?;
        Ok(())
    }

.. list-table:: thiserror vs anyhow 选择指南
   :header-rows: 1

   * - 维度
     - thiserror
     - anyhow
   * - 用途
     - 定义库的错误类型
     - 应用层的便捷错误处理
   * - 特点
     - 显式枚举，可匹配
     - 类型擦除，快速原型
   * - 适用
     - Library Crate
     - Binary Crate (main.rs)
   * - 组合
     - 可与 anyhow 互转
     - 可包装 thiserror 类型

eyre (彩色错误报告)
---------------------

``eyre`` 是 ``anyhow`` 的增强版，提供彩色错误报告和更多上下文。

.. code-block:: rust

    use color_eyre::{eyre::eyre, Result, Section};
    use color_eyre::owo_colors::OwoColorize;

    fn main() -> Result<()> {
        // 安装全局错误 Hook（彩色 panic + 回溯）
        color_eyre::install()?;

        let result: Result<()> = Err(eyre!("Something went wrong")
            .suggestion("Try running with --verbose for more details")
            .warning("This may cause data loss"));

        result?;

        Ok(())
    }

优雅关闭 (Graceful Shutdown)
-------------------------------

.. code-block:: rust

    use axum::{routing::get, Router};
    use tokio::signal;
    use tower_http::trace::TraceLayer;
    use sqlx::PgPool;

    async fn handler() -> &'static str {
        "Hello"
    }

    async fn shutdown_signal() {
        let ctrl_c = async {
            signal::ctrl_c()
                .await
                .expect("Failed to install Ctrl+C handler");
        };

        #[cfg(unix)]
        let terminate = async {
            signal::unix::signal(signal::unix::SignalKind::terminate())
                .expect("Failed to install signal handler")
                .recv()
                .await;
        };

        #[cfg(not(unix))]
        let terminate = std::future::pending::<()>();

        tokio::select! {
            _ = ctrl_c => {},
            _ = terminate => {},
        }

        tracing::info!("Shutdown signal received, starting graceful shutdown");
    }

    #[tokio::main]
    async fn main() {
        tracing_subscriber::fmt::init();

        let pool = PgPool::connect("postgres://localhost/mydb").await.unwrap();

        let app = Router::new()
            .route("/", get(handler))
            .layer(TraceLayer::new_for_http());

        let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();

        tracing::info!("Server listening on port 3000");

        axum::serve(listener, app)
            .with_graceful_shutdown(shutdown_signal())
            .await
            .unwrap();

        // 关闭后清理
        pool.close().await;
        tracing::info!("Server shut down gracefully");
    }

.. note::

   ``with_graceful_shutdown`` 会在收到信号后：
   1. 停止接受新连接
   2. 等待现有请求处理完成（默认无限等待，可通过 timeout 限制）
   3. 关闭服务器

.. code-block:: rust

    // 带超时的优雅关闭
    use tokio::time::{timeout, Duration};

    async fn shutdown_with_timeout() {
        // ... signal handling ...

        // 最多等待 30 秒完成现有请求
        match timeout(Duration::from_secs(30), async {
            // 等待关闭完成...
        }).await {
            Ok(_) => tracing::info!("Graceful shutdown complete"),
            Err(_) => tracing::warn!("Shutdown timed out, forcing exit"),
        }
    }

生产环境检查清单
---------------------

.. list-table:: 生产环境上线检查清单
   :header-rows: 1

   * - 类别
     - 检查项
   * - 构建
     - Release profile 已优化 (lto + strip + codegen-units=1)
   * - 配置
     - 敏感信息通过环境变量注入，不硬编码
   * - 日志
     - 使用 JSON 结构化日志，设置合适的日志级别
   * - 监控
     - 暴露 /health、/ready、/metrics 端点
   * - 错误
     - 错误不泄露内部细节，统一错误响应格式
   * - 关闭
     - 实现优雅关闭，处理 SIGTERM / SIGINT
   * - 安全
     - TLS 已配置，安全 Header 已设置
   * - 容器
     - 使用多阶段构建，非 root 用户运行
   * - 资源
     - 设置 CPU/Memory limits，配置健康检查
   * - CI/CD
     - 自动化测试、lint、安全审计、镜像构建

总结
-----

.. list-table:: 配置管理与错误处理 Crate 总览
   :header-rows: 1

   * - Crate
     - 用途
     - 适用场景
   * - ``config``
     - 分层配置管理
     - 多环境配置
   * - ``dotenvy``
     - .env 文件加载
     - 开发环境
   * - ``figment``
     - 声明式配置
     - 复杂配置需求
   * - ``clap``
     - CLI 参数解析
     - 命令行工具
   * - ``thiserror``
     - 库错误类型定义
     - Library Crate
   * - ``anyhow``
     - 应用错误处理
     - Binary Crate
   * - ``color-eyre``
     - 彩色错误报告
     - 开发调试
   * - ``axum::serve`` (graceful)
     - 优雅关闭
     - 所有 HTTP 服务
