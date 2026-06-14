日志与监控 (Logging & Monitoring)
==================================

.. .. module:: r07_deployment_operations.r03_logging_monitoring

可观测性（Observability）是生产运维的核心：日志记录应用行为，指标量化系统状态，追踪定位调用链路。
本章介绍 Rust 中的日志、指标和分布式追踪方案。

tracing (结构化日志与分布式追踪)
----------------------------------

``tracing`` 是 Rust 生态最强大的可观测性框架，支持结构化日志、Span 追踪和异步上下文传播。

**基础使用：**

.. code-block:: rust

    use tracing::{info, warn, error, debug, span, Level, instrument};
    use tracing_subscriber;

    #[tokio::main]
    async fn main() {
        // 初始化订阅器
        tracing_subscriber::fmt()
            .with_max_level(Level::DEBUG)
            .with_target(false)      // 不显示模块路径
            .with_thread_ids(true)   // 显示线程 ID
            .with_file(true)         // 显示文件名
            .with_line_number(true)  // 显示行号
            .json()                  // JSON 格式输出（生产环境推荐）
            .init();

        info!(user_id = 42, "User logged in");
        debug!("Processing request...");

        let span = span!(Level::INFO, "handle_request", request_id = "abc123");
        let _guard = span.enter();
        info!("Processing payment");

        do_work().await;
    }

    #[instrument]
    async fn do_work() {
        info!("Working...");
        tokio::time::sleep(std::time::Duration::from_millis(100)).await;
        warn!(duration_ms = 100, "Work completed slowly");
    }

**结构化字段与 Span：**

.. code-block:: rust

    use tracing::{info, instrument, Span};

    #[instrument(skip(db_pool), fields(user_id = %user.id))]
    async fn get_user(
        db_pool: &sqlx::PgPool,
        user: &User,
    ) -> Result<User, sqlx::Error> {
        info!("Fetching user from database");

        // 在 Span 中添加额外字段
        Span::current().record("db_query_time_ms", 15);

        let row = sqlx::query_as::<_, User>("SELECT * FROM users WHERE id = $1")
            .bind(user.id)
            .fetch_one(db_pool)
            .await?;

        Ok(row)
    }

**与 Axum 集成：**

.. code-block:: rust

    use axum::{routing::get, Router, middleware};
    use tower_http::trace::TraceLayer;
    use tracing::info;

    async fn handler() -> &'static str {
        info!("Handling request");
        "Hello"
    }

    #[tokio::main]
    async fn main() {
        tracing_subscriber::fmt().json().init();

        let app = Router::new()
            .route("/", get(handler))
            .layer(
                TraceLayer::new_for_http()
                    .make_span_with(|request: &axum::http::Request<_>| {
                        tracing::info_span!(
                            "http_request",
                            method = %request.method(),
                            uri = %request.uri(),
                            version = ?request.version(),
                        )
                    })
            );

        let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
        axum::serve(listener, app).await.unwrap();
    }

.. list-table:: tracing 核心概念
   :header-rows: 1

   * - 概念
     - 说明
   * - ``Event``
     - 日志事件 (info!/warn!/error!/debug!/trace!)
   * - ``Span``
     - 操作的生命周期范围，可嵌套
   * - ``Subscriber``
     - 日志收集器，决定如何处理 Event/Span
   * - ``Layer``
     - 订阅器的中间件，可组合（fmt/json/过滤等）
   * - ``#[instrument]``
     - 自动为函数创建 Span，记录参数和返回值

env_logger (轻量级日志)
--------------------------

对于简单应用，``env_logger`` 配合 ``log`` Crate 是最轻量的方案。

.. code-block:: rust

    use log::{info, warn, error, debug, LevelFilter};
    use env_logger::Env;

    fn main() {
        env_logger::Builder::from_env(Env::default().default_filter_or("info"))
            .format_timestamp_millis()
            .format_module_path(true)
            .init();

        info!("Server starting on port 8080");
        debug!("This won't show with default 'info' level");
        warn!("Disk usage at 80%");
        error!("Failed to connect to database");
    }

.. code-block:: bash

    # 运行时控制日志级别
    RUST_LOG=debug cargo run
    RUST_LOG=my_app=debug,actix_web=info cargo run

.. list-table:: tracing vs log/env_logger 对比
   :header-rows: 1

   * - 维度
     - tracing
     - log + env_logger
   * - 结构化日志
     - ✅ 原生支持
     - ❌ 仅文本
   * - Span / 追踪
     - ✅ 核心功能
     - ❌ 不支持
   * - 异步感知
     - ✅ 自动传播上下文
     - ❌ 不支持
   * - 性能
     - 极佳（无锁设计）
     - 良好
   * - 复杂度
     - 中等
     - 低
   * - 生态
     - OpenTelemetry / Jaeger / 丰富 Layer
     - 简单直接
   * - 适合场景
     - 微服务、分布式系统
     - 简单应用、CLI 工具

OpenTelemetry 集成
----------------------

.. code-block:: rust

    use opentelemetry::{
        global,
        trace::{Tracer, TracerProvider},
        KeyValue,
    };
    use opentelemetry_sdk::{
        trace::{Config, RandomIdGenerator, Sampler, TracerProvider as SdkProvider},
        Resource,
    };
    use opentelemetry_otlp::WithExportConfig;
    use tracing_opentelemetry::OpenTelemetryLayer;
    use tracing_subscriber::layer::SubscriberExt;

    fn init_tracing() -> Result<(), Box<dyn std::error::Error>> {
        // 配置 OTLP Exporter（发送到 Jaeger/Tempo/...）
        let exporter = opentelemetry_otlp::new_exporter()
            .tonic()
            .with_endpoint("http://localhost:4317");

        let tracer_provider = SdkProvider::builder()
            .with_batch_exporter(exporter)
            .with_config(
                Config::default()
                    .with_sampler(Sampler::AlwaysOn)
                    .with_id_generator(RandomIdGenerator::default())
                    .with_resource(Resource::new(vec![
                        KeyValue::new("service.name", "my-service"),
                        KeyValue::new("service.version", "1.0.0"),
                    ])),
            )
            .build();

        let tracer = tracer_provider.tracer("my-service");
        global::set_tracer_provider(tracer_provider);

        let telemetry_layer = OpenTelemetryLayer::new(tracer);

        tracing_subscriber::registry()
            .with(tracing_subscriber::fmt::layer().json())
            .with(telemetry_layer)
            .init();

        Ok(())
    }

metrics (指标收集)
---------------------

``metrics`` Crate 提供与日志类似的 facade 模式，支持多种后端。

.. code-block:: rust

    use metrics::{counter, gauge, histogram, describe_counter, describe_gauge, Unit};
    use metrics_exporter_prometheus::PrometheusBuilder;
    use axum::{routing::get, Router};
    use std::time::Instant;

    async fn handle_request() -> String {
        let start = Instant::now();

        // 请求计数
        counter!("http.requests.total", 1, "method" => "GET", "path" => "/");

        // 模拟处理
        tokio::time::sleep(std::time::Duration::from_millis(50)).await;

        // 请求延迟直方图
        let elapsed = start.elapsed().as_secs_f64();
        histogram!("http.requests.duration", elapsed, "method" => "GET", "path" => "/");

        "OK".to_string()
    }

    async fn metrics_handler() -> String {
        metrics_exporter_prometheus::PrometheusBuilder::new().build().render()
    }

    #[tokio::main]
    async fn main() {
        // 初始化 Prometheus 导出器
        PrometheusBuilder::new()
            .install_recorder()
            .unwrap();

        // 注册指标描述
        describe_counter!("http.requests.total", Unit::Count, "Total HTTP requests");
        describe_gauge!("db.connections.active", Unit::Count, "Active DB connections");
        describe_histogram!("http.requests.duration", Unit::Seconds, "Request duration");

        let app = Router::new()
            .route("/", get(handle_request))
            .route("/metrics", get(metrics_handler));

        let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
        axum::serve(listener, app).await.unwrap();
    }

.. list-table:: metrics 后端选择
   :header-rows: 1

   * - Crate
     - 后端
     - 适用场景
   * - ``metrics-exporter-prometheus``
     - Prometheus pull 模式
     - Kubernetes / 标准监控
   * - ``metrics-exporter-statsd``
     - StatsD / Datadog
     - 传统监控栈
   * - ``metrics-exporter-tcp``
     - TCP 推送
     - 自定义采集器
   * - ``metrics-observer``
     - 内存观察
     - 测试 / 调试

sentry (错误追踪)
--------------------

.. code-block:: rust

    use sentry::{ClientOptions, IntoDsn};
    use sentry::integrations::tracing::SentryLayer;
    use tracing_subscriber::layer::SubscriberExt;

    fn init_sentry() -> sentry::ClientInitGuard {
        let _guard = sentry::init((
            "https://key@o0.ingest.sentry.io/0".into_dsn().unwrap(),
            sentry::ClientOptions {
                release: sentry::release_name!(),
                traces_sample_rate: 0.2,     // 20% 追踪采样
                environment: Some(
                    std::env::var("ENV").unwrap_or("development".into()).into()
                ),
                ..Default::default()
            },
        ));

        // 与 tracing 集成
        tracing_subscriber::registry()
            .with(tracing_subscriber::fmt::layer())
            .with(SentryLayer::default())
            .init();

        _guard
    }

    #[tracing::instrument]
    async fn risky_operation() {
        if let Err(e) = do_something().await {
            sentry::capture_error(&e);
            tracing::error!(error = %e, "Operation failed");
        }
    }

.. note::

   ``sentry`` 可自动捕获 panic 并上报，配合 ``#[instrument]`` 可关联 Span 信息。
   生产环境建议设置 ``traces_sample_rate`` 控制采样率。

健康检查端点
---------------

.. code-block:: rust

    use axum::{routing::get, Router, Json, http::StatusCode};
    use serde::Serialize;
    use sqlx::PgPool;

    #[derive(Serialize)]
    struct HealthStatus {
        status: &'static str,
        version: &'static str,
        uptime_seconds: u64,
    }

    async fn health_check() -> Json<HealthStatus> {
        Json(HealthStatus {
            status: "ok",
            version: env!("CARGO_PKG_VERSION"),
            uptime_seconds: 0, // 通过全局状态跟踪
        })
    }

    async fn readiness_check(
        axum::extract::State(pool): axum::extract::State<PgPool>,
    ) -> Result<Json<HealthStatus>, StatusCode> {
        // 检查数据库连接
        sqlx::query("SELECT 1")
            .execute(&pool)
            .await
            .map_err(|_| StatusCode::SERVICE_UNAVAILABLE)?;

        Ok(Json(HealthStatus {
            status: "ready",
            version: env!("CARGO_PKG_VERSION"),
            uptime_seconds: 0,
        }))
    }

    fn health_routes(pool: PgPool) -> Router {
        Router::new()
            .route("/health", get(health_check))
            .route("/ready", get(readiness_check))
            .with_state(pool)
    }

.. list-table:: 健康检查端点
   :header-rows: 1

   * - 端点
     - 用途
     - 检查内容
   * - ``/health``
     - Liveness Probe
     - 进程是否存活
   * - ``/ready``
     - Readiness Probe
     - 依赖服务（DB/Redis）是否就绪
   * - ``/metrics``
     - Prometheus Metrics
     - 指标数据

总结
-----

.. list-table:: 日志与监控 Crate 总览
   :header-rows: 1

   * - Crate
     - 用途
     - 适用场景
   * - ``tracing``
     - 结构化日志 + Span 追踪
     - 异步应用、微服务（推荐首选）
   * - ``log`` + ``env_logger``
     - 简单文本日志
     - CLI 工具、简单应用
   * - ``opentelemetry``
     - 分布式追踪
     - 微服务调用链分析
   * - ``metrics`` + Prometheus
     - 指标收集与导出
     - 性能监控、告警
   * - ``sentry``
     - 错误追踪
     - 生产环境异常监控
   * - ``tower-http::trace``
     - HTTP 请求追踪
     - Axum/Actix-web 集成
