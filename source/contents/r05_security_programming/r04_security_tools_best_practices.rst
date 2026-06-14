安全工具与最佳实践 (Security Tools & Best Practices)
===================================================

.. .. module:: r05_security_programming.r04_security_tools_best_practices

本章涵盖安全工具、输入验证、CSRF 防护、安全 Header、密钥管理和安全编码最佳实践。

secrecy (密钥安全)
--------------------

``secrecy`` 通过封装敏感数据防止其被意外打印、序列化或泄露到日志中。

.. code-block:: rust

    use secrecy::{Secret, ExposeSecret, SecretString};

    #[derive(Clone)]
    struct ApiClient {
        api_key: SecretString,
    }

    impl ApiClient {
        fn new(api_key: String) -> Self {
            Self {
                api_key: SecretString::from(api_key),
            }
        }

        fn authenticate(&self) {
            // 仅在需要时暴露密钥
            let key = self.api_key.expose_secret();
            // 使用 key 进行 API 调用...
        }
    }

    // Secret 不实现 Display / Debug
    let client = ApiClient::new("sk-abc123".to_string());
    // println!("{:?}", client.api_key);  // 编译错误！

    // 零化：Drop 时用 0 覆盖内存
    use secrecy::Zeroize;
    let mut secret = Secret::new([1u8, 2, 3, 4]);
    secret.zeroize();  // 立即从内存中擦除

.. note::

   ``secrecy`` 的 ``Secret<T>`` 在 Drop 时自动零化底层数据，
   确保敏感信息（密码、API Key、加密密钥）不会残留在内存中。

validator (输入验证)
--------------------

``validator`` 提供声明式输入验证，通过 derive 宏定义验证规则。

.. code-block:: rust

    use validator::{Validate, ValidationError};
    use serde::Deserialize;

    fn validate_not_blank(value: &str) -> Result<(), ValidationError> {
        if value.trim().is_empty() {
            return Err(ValidationError::new("must_not_be_blank"));
        }
        Ok(())
    }

    #[derive(Debug, Deserialize, Validate)]
    struct SignupRequest {
        #[validate(email(message = "Invalid email format"))]
        email: String,

        #[validate(length(min = 8, max = 128, message = "Password must be 8-128 chars"))]
        #[validate(custom(function = "validate_password_strength"))]
        password: String,

        #[validate(must_match(other = "password", message = "Passwords do not match"))]
        password_confirm: String,

        #[validate(range(min = 0, max = 150, message = "Invalid age"))]
        age: Option<u8>,
    }

    fn validate_password_strength(password: &str) -> Result<(), ValidationError> {
        let has_upper = password.chars().any(|c| c.is_uppercase());
        let has_lower = password.chars().any(|c| c.is_lowercase());
        let has_digit = password.chars().any(|c| c.is_ascii_digit());
        let has_special = password.chars().any(|c| !c.is_alphanumeric());

        if has_upper && has_lower && has_digit && has_special {
            Ok(())
        } else {
            Err(ValidationError::new("password_must_contain_upper_lower_digit_special"))
        }
    }

    // 在 Axum handler 中使用
    async fn signup(Json(payload): Json<SignupRequest>) -> impl axum::response::IntoResponse {
        match payload.validate() {
            Ok(_) => "Signup successful".into(),
            Err(e) => format!("Validation errors: {:?}", e).into(),
        }
    }

.. list-table:: validator 常用验证规则
   :header-rows: 1

   * - 规则
     - 说明
   * - ``email``
     - 验证邮箱格式
   * - ``url``
     - 验证 URL 格式
   * - ``length(min, max)``
     - 字符串长度范围
   * - ``range(min, max)``
     - 数值范围
   * - ``must_match``
     - 字段匹配验证
   * - ``contains``
     - 包含指定子串
   * - ``regex``
     - 正则表达式匹配
   * - ``credit_card``
     - 信用卡号格式
   * - ``custom``
     - 自定义验证函数
   * - ``nested``
     - 嵌套结构体验证

CSRF 防护
----------

Axum 中通过 ``axum-csrf`` 中间件防止跨站请求伪造。

.. code-block:: rust

    use axum::{routing::{get, post}, Router};
    use axum_csrf::{CsrfConfig, CsrfLayer, CsrfToken};

    async fn show_form(token: CsrfToken) -> String {
        // 表单中嵌入 CSRF token
        format!(
            r#"<form method="post" action="/submit">
                <input type="hidden" name="csrf_token" value="{}">
                <input type="text" name="data">
                <button type="submit">Submit</button>
            </form>"#,
            token.get()
        )
    }

    async fn submit_form() -> &'static str {
        "Form submitted successfully"
    }

    #[tokio::main]
    async fn main() {
        let csrf_config = CsrfConfig::default()
            .with_key(Some(*b"0123456789abcdef0123456789abcdef")); // 32 bytes AES key

        let app = Router::new()
            .route("/form", get(show_form))
            .route("/submit", post(submit_form))
            .layer(CsrfLayer::new(csrf_config));

        let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
        axum::serve(listener, app).await.unwrap();
    }

安全 HTTP Header
------------------

通过 ``tower-http`` 设置安全相关的 HTTP 响应头。

.. code-block:: rust

    use axum::{routing::get, Router};
    use tower_http::{
        set_header::SetResponseHeaderLayer,
        compression::CompressionLayer,
    };
    use hyper::header::{
        CONTENT_SECURITY_POLICY, STRICT_TRANSPORT_SECURITY,
        X_CONTENT_TYPE_OPTIONS, X_FRAME_OPTIONS,
    };

    async fn handler() -> &'static str {
        "Secure response"
    }

    #[tokio::main]
    async fn main() {
        let app = Router::new()
            .route("/", get(handler))
            .layer(SetResponseHeaderLayer::if_not_present(
                X_CONTENT_TYPE_OPTIONS,
                "nosniff".parse().unwrap(),
            ))
            .layer(SetResponseHeaderLayer::if_not_present(
                X_FRAME_OPTIONS,
                "DENY".parse().unwrap(),
            ))
            .layer(SetResponseHeaderLayer::if_not_present(
                STRICT_TRANSPORT_SECURITY,
                "max-age=31536000; includeSubDomains".parse().unwrap(),
            ))
            .layer(SetResponseHeaderLayer::if_not_present(
                CONTENT_SECURITY_POLICY,
                "default-src 'self'".parse().unwrap(),
            ));

        let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
        axum::serve(listener, app).await.unwrap();
    }

.. list-table:: 推荐的安全 Header
   :header-rows: 1

   * - Header
     - 值
     - 作用
   * - ``Strict-Transport-Security``
     - ``max-age=31536000; includeSubDomains``
     - 强制 HTTPS
   * - ``X-Content-Type-Options``
     - ``nosniff``
     - 禁止 MIME 类型嗅探
   * - ``X-Frame-Options``
     - ``DENY`` 或 ``SAMEORIGIN``
     - 防点击劫持
   * - ``Content-Security-Policy``
     - ``default-src 'self'``
     - 防 XSS / 数据注入
   * - ``X-XSS-Protection``
     - ``1; mode=block``
     - 启用浏览器 XSS 过滤器
   * - ``Referrer-Policy``
     - ``strict-origin-when-cross-origin``
     - 控制 Referer 信息泄露

CORS 安全配置
---------------

.. code-block:: rust

    use axum::{routing::get, Router};
    use tower_http::cors::{CorsLayer, Any, AllowOrigin};
    use hyper::Method;

    async fn handler() -> &'static str {
        "Hello"
    }

    #[tokio::main]
    async fn main() {
        // 严格 CORS 配置（推荐）
        let cors = CorsLayer::new()
            .allow_origin(AllowOrigin::exact(
                "https://myapp.com".parse().unwrap()
            ))
            .allow_methods([Method::GET, Method::POST])
            .allow_headers([hyper::header::CONTENT_TYPE])
            .allow_credentials(true)
            .max_age(std::time::Duration::from_secs(3600));

        let app = Router::new()
            .route("/", get(handler))
            .layer(cors);

        let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
        axum::serve(listener, app).await.unwrap();
    }

.. warning::

   不要在生产环境使用 ``Any`` 作为 ``allow_origin``。始终明确指定允许的源。

Rate Limiting (速率限制)
-------------------------

``tower_governor`` 基于 Token Bucket 算法提供速率限制中间件。

.. code-block:: rust

    use axum::{routing::get, Router, response::IntoResponse};
    use tower_governor::{GovernorLayer, governor::GovernorConfigBuilder};
    use std::time::Duration;

    async fn handler() -> impl IntoResponse {
        "Rate limited endpoint"
    }

    #[tokio::main]
    async fn main() {
        // 每秒 5 个请求，突发容量 10
        let governor_config = Box::new(
            GovernorConfigBuilder::default()
                .per_second(5)
                .burst_size(10)
                .finish()
                .unwrap(),
        );

        let app = Router::new()
            .route("/api", get(handler))
            .layer(GovernorLayer {
                config: Box::leak(governor_config),
            });

        let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
        axum::serve(listener, app).await.unwrap();
    }

SQL 注入防护
---------------

Rust 的类型系统和参数化查询天然防止 SQL 注入。

.. code-block:: rust

    use sqlx::postgres::PgPool;

    // ✅ 安全：参数化查询
    async fn safe_query(pool: &PgPool, user_id: i32) -> Result<String, sqlx::Error> {
        let row: (String,) = sqlx::query_as(
            "SELECT name FROM users WHERE id = $1"
        )
        .bind(user_id)  // 参数绑定，不会 SQL 注入
        .fetch_one(pool)
        .await?;
        Ok(row.0)
    }

    // ❌ 危险：字符串拼接（仅作警示）
    // async fn dangerous_query(pool: &PgPool, user_id: &str) -> Result<(), sqlx::Error> {
    //     let query = format!("SELECT * FROM users WHERE id = {}", user_id);
    //     sqlx::query(&query).execute(pool).await?;
    //     Ok(())
    // }

    // ✅ 动态表名/列名：使用白名单
    fn safe_dynamic_column(column: &str) -> Result<String, &'static str> {
        const ALLOWED: &[&str] = &["name", "email", "created_at"];
        if ALLOWED.contains(&column) {
            Ok(format!("SELECT {} FROM users", column))
        } else {
            Err("Invalid column name")
        }
    }

安全配置管理
---------------

.. code-block:: rust

    use secrecy::SecretString;
    use serde::Deserialize;
    use std::env;

    #[derive(Deserialize)]
    struct AppConfig {
        database_url: SecretString,
        jwt_secret: SecretString,
        api_key: SecretString,
        aws_secret_access_key: SecretString,
    }

    impl AppConfig {
        fn from_env() -> Result<Self, Box<dyn std::error::Error>> {
            // 从环境变量读取（不硬编码）
            Ok(envy::from_env::<Self>()?)
        }
    }

.. list-table:: 安全配置最佳实践
   :header-rows: 1

   * - 实践
     - 说明
   * - 使用环境变量
     - 敏感配置不写入代码或版本控制
   * - 使用 ``secrecy::SecretString``
     - 防止配置被意外打印或序列化
   * - ``.env`` 文件加入 ``.gitignore``
     - 防止凭据泄露到 Git 仓库
   * - 使用密钥管理服务 (KMS)
     - 生产环境通过 HashiCorp Vault / AWS KMS 获取密钥
   * - 最小权限原则
     - 数据库账号、API Key 仅授予必要权限

SQLx 与数据库安全
--------------------

.. code-block:: rust

    use sqlx::postgres::{PgPool, PgPoolOptions};

    async fn create_secure_pool() -> Result<PgPool, sqlx::Error> {
        let pool = PgPoolOptions::new()
            .max_connections(10)
            .connect_timeout(std::time::Duration::from_secs(5))
            // 使用 sslmode=require 强制 TLS
            .connect("postgres://user:pass@host/db?sslmode=require")
            .await?;

        Ok(pool)
    }

    // 预编译语句防止注入
    async fn insert_user(pool: &PgPool, name: &str, email: &str) -> Result<(), sqlx::Error> {
        sqlx::query("INSERT INTO users (name, email) VALUES ($1, $2)")
            .bind(name)
            .bind(email)
            .execute(pool)
            .await?;
        Ok(())
    }

cargo-audit (依赖安全审计)
----------------------------

.. code-block:: bash

    # 安装
    cargo install cargo-audit

    # 审计项目依赖中的已知漏洞
    cargo audit

    # 输出示例：
    # Crate:     time
    # Version:   0.1.44
    # Title:     Potential segfault in the time crate
    # Date:      2020-11-18
    # ID:        RUSTSEC-2020-0071
    # URL:       https://rustsec.org/advisories/RUSTSEC-2020-0071
    # Solution:  Upgrade to >=0.2.23

.. note::

   建议在 CI 流程中集成 ``cargo audit``，定期扫描项目依赖中的已知 CVE。
   可与 ``cargo-deny`` 配合使用，获得更全面的依赖策略管理。

安全编码检查清单
--------------------

.. list-table:: Rust 安全编码检查清单
   :header-rows: 1

   * - 类别
     - 检查项
   * - 输入验证
     - 所有外部输入必须验证（validator / 自定义验证）
   * - 密码存储
     - 使用 Argon2id，带独立随机盐，最少 64MB 内存
   * - Token 安全
     - JWT 使用 RS256/ES256，设置合理过期时间
   * - 传输安全
     - 所有通信使用 TLS 1.3，禁止明文传输
   * - 密钥管理
     - 使用 ``secrecy`` 包装，环境变量注入，禁止硬编码
   * - SQL 查询
     - 使用参数化查询，动态列名/表名使用白名单
   * - CSRF
     - 状态变更请求必须验证 CSRF Token
   * - CORS
     - 明确指定 Allow-Origin，不使用 ``*``
   * - 安全 Header
     - HSTS / CSP / X-Frame-Options / X-Content-Type-Options
   * - 依赖审计
     - CI 中集成 ``cargo audit``，定期更新依赖
   * - 日志安全
     - 不记录密码/Token/PII；使用结构化日志
   * - 错误处理
     - 不向客户端暴露内部错误细节

总结
-----

.. list-table:: 安全工具与最佳实践 Crate 总览
   :header-rows: 1

   * - Crate / 工具
     - 用途
     - 适用场景
   * - ``secrecy``
     - 敏感数据封装与零化
     - 密钥、密码、Token 管理
   * - ``validator``
     - 声明式输入验证
     - Web API 请求体验证
   * - ``axum-csrf``
     - CSRF 防护
     - Axum Web 应用
   * - ``tower-http`` (CORS / Header)
     - CORS 配置、安全 Header
     - 所有 HTTP 服务
   * - ``tower_governor``
     - 速率限制
     - API 防护、防暴力破解
   * - ``sqlx``
     - 参数化查询
     - 数据库操作
   * - ``cargo-audit``
     - 依赖漏洞扫描
     - CI/CD 流程
   * - ``cargo-deny``
     - 依赖许可证/安全策略
     - CI/CD 流程
