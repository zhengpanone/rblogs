================================
模板引擎与前端集成
================================

Rust 服务端渲染（SSR）模板引擎及前后端分离集成方案。

.. contents:: 目录
   :depth: 3
   :local:

模板引擎概述
==============

Rust 生态中主要有两类模板引擎：编译时模板（类型安全，编译期检查）和运行时模板（灵活，动态加载）。

.. list-table::
   :header-rows: 1
   :widths: 15 20 25 20

   * - 引擎
     - 检查时机
     - 语法风格
     - 特点
   * - Askama
     - 编译时
     - Jinja-like
     - 类型安全、零运行时开销
   * - Tera
     - 运行时
     - Jinja2-like
     - 功能丰富、模板继承
   * - MiniJinja
     - 运行时
     - Jinja2 兼容
     - 与 Python Jinja2 高度兼容
   * - Maud
     - 编译时
     - Rust 宏
     - 宏生成 HTML，类型安全

Askama（编译时模板）
======================

.. code-block:: toml

   [dependencies]
   askama = "0.12"
   axum = "0.8"

定义模板文件 ``templates/hello.html``：

.. code-block:: html

   <!DOCTYPE html>
   <html>
   <head>
       <meta charset="utf-8">
       <title>{{ title }}</title>
   </head>
   <body>
       <h1>Hello, {{ name }}!</h1>
       <p>你有 {{ messages.len() }} 条消息</p>
       <ul>
       {% for msg in messages %}
           <li>{{ msg }}</li>
       {% endfor %}
       </ul>

       {% if is_admin %}
       <p><a href="/admin">管理后台</a></p>
       {% endif %}
   </body>
   </html>

Rust 代码：

.. code-block:: rust

   use askama::Template;
   use axum::{Router, routing::get, response::Html};

   #[derive(Template)]
   #[template(path = "hello.html")]
   struct HelloTemplate {
       title: String,
       name: String,
       messages: Vec<String>,
       is_admin: bool,
   }

   async fn hello() -> Html<String> {
       let template = HelloTemplate {
           title: "欢迎".to_string(),
           name: "Rust".to_string(),
           messages: vec![
               "欢迎回来！".to_string(),
               "你有新的通知".to_string(),
           ],
           is_admin: true,
       };

       Html(template.render().unwrap())
   }

   #[tokio::main]
   async fn main() {
       let app = Router::new()
           .route("/", get(hello));

       let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
           .await.unwrap();
       axum::serve(listener, app).await.unwrap();
   }

模板继承 ``templates/base.html``：

.. code-block:: html

   <!DOCTYPE html>
   <html>
   <head>
       <meta charset="utf-8">
       <title>{% block title %}默认标题{% endblock %}</title>
   </head>
   <body>
       <nav>
           <a href="/">首页</a>
           <a href="/about">关于</a>
       </nav>
       <main>
           {% block content %}{% endblock %}
       </main>
       <footer>© 2024 My App</footer>
   </body>
   </html>

子模板 ``templates/page.html``：

.. code-block:: html

   {% extends "base.html" %}

   {% block title %}{{ page_title }}{% endblock %}

   {% block content %}
   <h1>{{ page_title }}</h1>
   <p>{{ content }}</p>
   {% endblock %}

Rust 代码：

.. code-block:: rust

   use askama::Template;

   #[derive(Template)]
   #[template(path = "page.html")]
   struct PageTemplate {
       page_title: String,
       content: String,
   }

Askama 语法速查：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 语法
     - 说明
   * - ``{{ var }}``
     - 输出变量（自动 HTML 转义）
   * - ``{{ var|safe }}``
     - 不转义输出（小心 XSS）
   * - ``{% if cond %}...{% endif %}``
     - 条件判断
   * - ``{% for item in list %}...{% endfor %}``
     - 循环
   * - ``{% block name %}...{% endblock %}``
     - 模板继承块
   * - ``{% extends "base.html" %}``
     - 继承父模板
   * - ``{% include "partial.html" %}``
     - 包含子模板
   * - ``{% match val %}{% when ... %}...{% endmatch %}``
     - 模式匹配
   * - ``{{ list|join(", ") }}``
     - 过滤器

Tera（运行时模板）
====================

.. code-block:: toml

   [dependencies]
   tera = "1"
   axum = "0.8"

模板文件 ``templates/index.html``：

.. code-block:: html

   <!DOCTYPE html>
   <html>
   <head><title>{{ title }}</title></head>
   <body>
       <h1>{{ title }}</h1>
       <ul>
       {% for user in users %}
           <li>{{ user.name }} - {{ user.email }}</li>
       {% endfor %}
       </ul>
   </body>
   </html>

Rust 代码：

.. code-block:: rust

   use axum::{Router, routing::get, response::Html, extract::State};
   use tera::{Tera, Context};
   use std::sync::Arc;
   use serde::Serialize;

   #[derive(Clone)]
   struct AppState {
       tera: Arc<Tera>,
   }

   #[derive(Serialize)]
   struct User {
       name: String,
       email: String,
   }

   async fn index(State(state): State<AppState>) -> Html<String> {
       let mut ctx = Context::new();
       ctx.insert("title", "用户列表");
       ctx.insert("users", &vec![
           User { name: "Alice".into(), email: "alice@example.com".into() },
           User { name: "Bob".into(), email: "bob@example.com".into() },
       ]);

       let rendered = state.tera.render("index.html", &ctx).unwrap();
       Html(rendered)
   }

   #[tokio::main]
   async fn main() {
       let tera = Tera::new("templates/**/*.html").unwrap();
       let state = AppState { tera: Arc::new(tera) };

       let app = Router::new()
           .route("/", get(index))
           .with_state(state);

       let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
           .await.unwrap();
       axum::serve(listener, app).await.unwrap();
   }

Tera 常用功能：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 功能
     - 说明
   * - ``{% if %}`` / ``{% elif %}`` / ``{% else %}``
     - 条件分支
   * - ``{% for %}`` / ``{% endfor %}``
     - 循环
   * - ``{% set var = value %}``
     - 变量赋值
   * - ``{% block %}`` / ``{% extends %}``
     - 模板继承
   * - ``{% include "file" %}``
     - 包含子模板
   * - ``{% macro name(args) %}``
     - 宏定义
   * - ``{{ var | filter }}``
     - 过滤器（lower, upper, date, json_encode 等）
   * - ``{{ loop.index }}`` / ``{{ loop.first }}``
     - 循环内置变量

Askama vs Tera 对比：

.. list-table::
   :header-rows: 1
   :widths: 20 40 40

   * - 特性
     - Askama
     - Tera
   * - 类型检查
     - 编译时（变量缺失=编译错误）
     - 运行时（变量缺失=渲染错误）
   * - 模板变更
     - 需重新编译
     - 无需重新编译
   * - 性能
     - 零运行时开销
     - 运行时解析
   * - 热重载
     - 不支持（编译时绑定）
     - 支持
   * - 学习曲线
     - 中等
     - 较低
   * - 适用场景
     - 模板稳定的项目
     - 需要热重载、动态模板

Maud（宏模板）
===============

Maud 使用 Rust 宏生成 HTML，类型安全且高性能。

.. code-block:: toml

   [dependencies]
   maud = "0.26"
   axum = "0.8"

.. code-block:: rust

   use maud::{html, Markup, DOCTYPE};

   fn layout(title: &str, content: Markup) -> Markup {
       html! {
           (DOCTYPE)
           html lang="zh" {
               head {
                   meta charset="utf-8";
                   title { (title) }
                   link rel="stylesheet" href="/style.css";
               }
               body {
                   header {
                       h1 { "My App" }
                       nav {
                           a href="/" { "首页" }
                           a href="/about" { "关于" }
                       }
                   }
                   main { (content) }
                   footer { "© 2024" }
               }
           }
       }
   }

   fn index() -> Markup {
       layout("首页", html! {
           h2 { "欢迎" }
           p { "这是一段内容" }
           ul {
               @for i in 1..=3 {
                   li { "项目 " (i) }
               }
           }
       })
   }

   fn about() -> Markup {
       layout("关于", html! {
           h2 { "关于我们" }
           p { "我们使用 Rust 构建 Web 应用" }
       })
   }

   // 与 Axum 集成
   use axum::{Router, routing::get, response::Html};

   async fn index_handler() -> Html<String> {
       Html(index().into_string())
   }

   #[tokio::main]
   async fn main() {
       let app = Router::new()
           .route("/", get(index_handler))
           .route("/about", get(|| async { Html(about().into_string()) }));

       let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
           .await.unwrap();
       axum::serve(listener, app).await.unwrap();
   }

前后端分离
============

Rust 作为 API 后端，前端使用 React / Vue / Svelte 等框架。

静态文件服务（Axum）：

.. code-block:: rust

   use axum::{Router, routing::get, response::Html};
   use tower_http::services::{ServeDir, ServeFile};

   #[tokio::main]
   async fn main() {
       let app = Router::new()
           // API 路由
           .nest("/api", api_routes())
           // 静态文件
           .nest_service("/assets", ServeDir::new("dist/assets"))
           // SPA fallback：所有非 API 路径返回 index.html
           .fallback_service(ServeFile::new("dist/index.html"));

       let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
           .await.unwrap();
       axum::serve(listener, app).await.unwrap();
   }

   fn api_routes() -> Router {
       Router::new()
           .route("/users", get(list_users).post(create_user))
           .route("/health", get(|| async { "OK" }))
   }

CORS 配置（前后端分离必备）：

.. code-block:: rust

   use tower_http::cors::{CorsLayer, AllowOrigin, AllowMethods, AllowHeaders};
   use http::{Method, header};

   fn cors_config() -> CorsLayer {
       CorsLayer::new()
           // 开发环境允许 localhost
           .allow_origin(AllowOrigin::list([
               "http://localhost:5173".parse().unwrap(),   // Vite
               "http://localhost:3000".parse().unwrap(),   // React
               "http://localhost:8080".parse().unwrap(),   // Vue
           ]))
           .allow_methods(AllowMethods::list([
               Method::GET, Method::POST, Method::PUT,
               Method::DELETE, Method::OPTIONS,
           ]))
           .allow_headers(AllowHeaders::list([
               header::AUTHORIZATION,
               header::CONTENT_TYPE,
               header::ACCEPT,
           ]))
           .allow_credentials(true)
   }

使用 rust-embed 嵌入前端构建产物：

.. code-block:: rust

   use rust_embed::RustEmbed;
   use axum::{
       Router, routing::get,
       response::{Html, IntoResponse},
       http::{header, StatusCode},
   };

   #[derive(RustEmbed)]
   #[folder = "dist/"]
   struct Assets;

   async fn serve_static(path: &str) -> impl IntoResponse {
       match Assets::get(path) {
           Some(file) => {
               let mime = mime_guess::from_path(path).first_or_octet_stream();
               ([(header::CONTENT_TYPE, mime.as_ref())], file.data).into_response()
           }
           None => (StatusCode::NOT_FOUND, "Not Found").into_response(),
       }
   }

   async fn spa_fallback() -> impl IntoResponse {
       match Assets::get("index.html") {
           Some(file) => Html(file.data).into_response(),
           None => (StatusCode::NOT_FOUND, "index.html not found").into_response(),
       }
   }

API 文档生成
==============

使用 utoipa 自动生成 OpenAPI/Swagger 文档：

.. code-block:: toml

   [dependencies]
   utoipa = { version = "5", features = ["axum_extras"] }
   utoipa-swagger-ui = { version = "8", features = ["axum"] }

.. code-block:: rust

   use axum::{Router, routing::get, Json};
   use utoipa::{OpenApi, ToSchema};
   use utoipa_swagger_ui::SwaggerUi;
   use serde::Serialize;

   #[derive(OpenApi)]
   #[openapi(
       paths(list_users, get_user, create_user),
       components(schemas(User, CreateUserPayload)),
   )]
   struct ApiDoc;

   #[derive(Serialize, ToSchema)]
   struct User {
       id: u64,
       name: String,
       email: String,
   }

   #[derive(serde::Deserialize, ToSchema)]
   struct CreateUserPayload {
       name: String,
       email: String,
   }

   /// 获取用户列表
   #[utoipa::path(
       get,
       path = "/users",
       responses(
           (status = 200, description = "用户列表", body = [User])
       )
   )]
   async fn list_users() -> Json<Vec<User>> {
       Json(vec![])
   }

   /// 根据 ID 获取用户
   #[utoipa::path(
       get,
       path = "/users/{id}",
       params(
           ("id" = u64, Path, description = "用户 ID")
       ),
       responses(
           (status = 200, description = "用户详情", body = User),
           (status = 404, description = "用户不存在")
       )
   )]
   async fn get_user(
       axum::extract::Path(id): axum::extract::Path<u64>,
   ) -> Json<User> {
       Json(User { id, name: "Alice".into(), email: "alice@example.com".into() })
   }

   /// 创建用户
   #[utoipa::path(
       post,
       path = "/users",
       request_body = CreateUserPayload,
       responses(
           (status = 201, description = "创建成功", body = User)
       )
   )]
   async fn create_user(
       Json(payload): Json<CreateUserPayload>,
   ) -> (axum::http::StatusCode, Json<User>) {
       (axum::http::StatusCode::CREATED, Json(User {
           id: 1,
           name: payload.name,
           email: payload.email,
       }))
   }

   #[tokio::main]
   async fn main() {
       let app = Router::new()
           .route("/users", get(list_users).post(create_user))
           .route("/users/{id}", get(get_user))
           .merge(SwaggerUi::new("/swagger-ui")
               .url("/api-docs/openapi.json", ApiDoc::openapi()));

       let listener = tokio::net::TcpListener::bind("127.0.0.1:3000")
           .await.unwrap();
       axum::serve(listener, app).await.unwrap();
   }

静态资源嵌入方案对比：

.. list-table::
   :header-rows: 1
   :widths: 20 40 20

   * - 方案
     - 说明
     - 适用场景
   * - ``tower_http::services::ServeDir``
     - 从文件系统提供静态文件
     - 开发环境
   * - ``rust-embed``
     - 编译期嵌入文件到二进制
     - 单文件部署
   * - ``include_str!`` / ``include_bytes!``
     - 编译期嵌入单个文件
     - 小文件、配置
   * - CDN
     - 外部 CDN 托管静态资源
     - 生产环境

总结
======

模板引擎与前端方案选择：

.. list-table::
   :header-rows: 1
   :widths: 20 40 20

   * - 方案
     - 适用场景
     - 推荐 Crate
   * - 服务端渲染 (SSR)
     - 内容网站、管理后台、SEO 敏感
     - Askama / Tera
   * - 前后端分离
     - SPA、复杂交互、移动端 API
     - Axum + React/Vue + CORS
   * - 嵌入式单文件部署
     - CLI 工具、小应用
     - rust-embed + Askama
   * - API 文档
     - 所有 API 项目
     - utoipa + Swagger UI
   * - 微服务通信
     - 服务间调用
     - tonic (gRPC)
   * - 实时通信
     - 聊天、通知、协作
     - WebSocket (axum ws / actix-ws)
