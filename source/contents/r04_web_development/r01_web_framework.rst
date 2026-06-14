=======================
Web 框架
=======================

Rust Web 开发的主流框架介绍、对比与实践。

.. contents:: 目录
   :depth: 3
   :local:

框架概览
==========

Rust Web 框架生态丰富，以下为最主流的五个框架：

.. list-table::
   :header-rows: 1
   :widths: 15 20 20 25

   * - 框架
     - 异步运行时
     - 风格
     - 特点
   * - Actix-web
     - actix-rt (tokio)
     - 类 MVC
     - 高性能、功能完整、社区活跃
   * - Axum
     - tokio
     - 模块化
     - 类型安全、基于 Tower、与 tokio 深度集成
   * - Rocket
     - tokio
     - 声明式
     - 零配置、开发友好、宏驱动
   * - Warp
     - tokio
     - 函数式
     - Filter 组合、类型安全
   * - Tide
     - async-std
     - 中间件
     - 简洁、渐进式

Actix-web
===========

Actix-web 是 Rust 社区历史最久、性能最高的 Web 框架之一。基于 actor 模型构建，但在 4.x 中 actor 已变为可选。

特性：

- 高性能：常年位于 TechEmpower 基准测试前列
- 类型安全：请求/响应处理类型明确
- 丰富的中间件生态：日志、压缩、会话、CORS
- 支持 HTTP/1.x、HTTP/2、WebSocket、TLS

基础示例：

.. code-block:: rust

   use actix_web::{App, HttpResponse, HttpServer, Responder, web, middleware};

   async fn index() -> impl Responder {
       HttpResponse::Ok().body("Hello World!")
   }

   async fn greet(path: web::Path<String>) -> impl Responder {
       HttpResponse::Ok().body(format!("Hello {}!", path))
   }

   // JSON 请求处理
   #[derive(serde::Deserialize)]
   struct CreateUser {
       name: String,
       email: String,
   }

   async fn create_user(user: web::Json<CreateUser>) -> impl Responder {
       HttpResponse::Created().json(serde_json::json!({
           "name": &user.name,
           "email": &user.email,
       }))
   }

   #[actix_web::main]
   async fn main() -> std::io::Result<()> {
       HttpServer::new(|| {
           App::new()
               .route("/", web::get().to(index))
               .route("/{name}", web::get().to(greet))
               .route("/users", web::post().to(create_user))
       })
       .bind("127.0.0.1:8080")?
       .run()
       .await
   }

状态共享（Application State）：

.. code-block:: rust

   use actix_web::{web, App, HttpServer, HttpResponse};
   use std::sync::Mutex;

   struct AppState {
       counter: Mutex<i64>,
   }

   async fn count(data: web::Data<AppState>) -> HttpResponse {
       let mut counter = data.counter.lock().unwrap();
       *counter += 1;
       HttpResponse::Ok().body(format!("访问次数: {}", counter))
   }

   #[actix_web::main]
   async fn main() -> std::io::Result<()> {
       let data = web::Data::new(AppState {
           counter: Mutex::new(0),
       });

       HttpServer::new(move || {
           App::new()
               .app_data(data.clone())
               .route("/count", web::get().to(count))
       })
       .bind("127.0.0.1:8080")?
       .run()
       .await
   }

中间件：

.. code-block:: rust

   use actix_web::{App, HttpServer, middleware};

   #[actix_web::main]
   async fn main() -> std::io::Result<()> {
       HttpServer::new(|| {
           App::new()
               .wrap(middleware::Logger::default())       // 日志
               .wrap(middleware::Compress::default())     // 压缩
               .wrap(middleware::NormalizePath::trim())   // 路径规范化
               // ... routes
       })
       .bind("127.0.0.1:8080")?
       .run()
       .await
   }

提取器（Extractors）：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 提取器
     - 说明
   * - ``web::Path<T>``
     - 路径参数（/users/{id}）
   * - ``web::Query<T>``
     - URL 查询参数
   * - ``web::Json<T>``
     - JSON 请求体
   * - ``web::Form<T>``
     - 表单数据
   * - ``web::Data<T>``
     - 应用状态
   * - ``HttpRequest``
     - 原始请求
   * - ``web::Payload``
     - 请求体流（用于文件上传等）

Axum
======

Axum 是 tokio 团队开发的模块化 Web 框架，基于 Tower 和 Hyper 构建，强调类型安全和可组合性。

特性：

- 类型安全的路由和提取器
- 基于 Tower 的中间件生态（tower-http）
- 原生 async/await 支持
- 共享状态、依赖注入
- 与 tokio 生态无缝集成

基础示例：

.. code-block:: rust

   use axum::{
       Router,
       extract::{Path, Query, State},
       response::Json,
       routing::get,
   };
   use serde::Deserialize;
   use std::sync::Arc;

   #[derive(Clone)]
   struct AppState {
       db_pool: String, // 模拟
   }

   async fn root() -> Json<serde_json::Value> {
       Json(serde_json::json!({ "message": "Hello, World!" }))
   }

   async fn greet(Path(name): Path<String>) -> Json<serde_json::Value> {
       Json(serde_json::json!({ "message": format!("Hello, {}!", name) }))
   }

   #[derive(Deserialize)]
   struct SearchParams {
       q: String,
       page: Option<u32>,
   }

   async fn search(Query(params): Query<SearchParams>) -> Json<serde_json::Value> {
       Json(serde_json::json!({
           "query": params.q,
           "page": params.page.unwrap_or(1),
       }))
   }

   async fn get_state(State(state): State<AppState>) -> Json<serde_json::Value> {
       Json(serde_json::json!({ "db": state.db_pool }))
   }

   #[tokio::main]
   async fn main() {
       let state = AppState {
           db_pool: "postgres://localhost/mydb".to_string(),
       };

       let app = Router::new()
           .route("/", get(root))
           .route("/{name}", get(greet))
           .route("/search", get(search))
           .route("/state", get(get_state))
           .with_state(state);

       let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
           .await
           .unwrap();
       axum::serve(listener, app).await.unwrap();
   }

嵌套路由：

.. code-block:: rust

   use axum::{Router, routing::get};

   async fn list_users() -> &'static str { "用户列表" }
   async fn get_user() -> &'static str { "用户详情" }
   async fn create_user() -> &'static str { "创建用户" }

   #[tokio::main]
   async fn main() {
       let user_routes = Router::new()
           .route("/", get(list_users).post(create_user))
           .route("/{id}", get(get_user));

       let app = Router::new()
           .route("/", get(|| async { "首页" }))
           .nest("/users", user_routes)
           .nest("/api", api_routes());

       // ...
   }

   fn api_routes() -> Router {
       Router::new()
           .route("/health", get(|| async { "OK" }))
   }

中间件（tower-http）：

.. code-block:: rust

   use tower_http::{
       cors::CorsLayer,
       trace::TraceLayer,
       compression::CompressionLayer,
       limit::RequestBodyLimitLayer,
   };
   use axum::Router;

   fn app() -> Router {
       Router::new()
           // 请求日志
           .layer(TraceLayer::new_for_http())
           // CORS
           .layer(CorsLayer::permissive())
           // 响应压缩
           .layer(CompressionLayer::new())
           // 请求体大小限制 (5MB)
           .layer(RequestBodyLimitLayer::new(5 * 1024 * 1024))
           // ... routes
   }

Axum 提取器：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 提取器
     - 说明
   * - ``Path<T>``
     - 路径参数
   * - ``Query<T>``
     - URL 查询参数
   * - ``Json<T>``
     - JSON 请求体
   * - ``Form<T>``
     - 表单数据
   * - ``State<T>``
     - 共享应用状态
   * - ``Extension<T>``
     - 请求扩展（中间件注入）
   * - ``Headers`` / ``HeaderMap``
     - 请求头
   * - ``Method`` / ``Uri``
     - HTTP 方法和 URI
   * - ``Bytes`` / ``String``
     - 原始请求体
   * - ``Request<Body>``
     - 完整请求对象

Rocket
========

Rocket 以开发体验著称，宏驱动、零配置即可运行。

.. code-block:: toml

   [dependencies]
   rocket = "0.5"

.. code-block:: rust

   #[macro_use] extern crate rocket;

   #[get("/")]
   fn index() -> &'static str {
       "Hello, World!"
   }

   #[get("/<name>")]
   fn greet(name: &str) -> String {
       format!("Hello, {}!", name)
   }

   #[get("/search?<q>&<page>")]
   fn search(q: String, page: Option<u32>) -> String {
       format!("搜索: {}, 页码: {}", q, page.unwrap_or(1))
   }

   #[derive(serde::Deserialize, rocket::form::FromForm)]
   struct LoginForm<'r> {
       username: &'r str,
       password: &'r str,
   }

   #[post("/login", data = "<form>")]
   fn login(form: rocket::form::Form<LoginForm<'_>>) -> String {
       format!("登录: {}", form.username)
   }

   #[launch]
   fn rocket() -> _ {
       rocket::build()
           .mount("/", routes![index, greet, search, login])
   }

Rocket 路由属性：

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - 属性
     - 说明
   * - ``#[get("/path")]``
     - GET 请求
   * - ``#[post("/path", data = "<param>")]``
     - POST 请求
   * - ``#[put("/path", data = "<param>")]``
     - PUT 请求
   * - ``#[delete("/path")]``
     - DELETE 请求
   * - ``#[patch("/path", data = "<param>")]``
     - PATCH 请求
   * - ``<param>``
     - 动态路径段
   * - ``<param..>``
     - 多个路径段
   * - ``?<param>``
     - 可选查询参数

Warp
======

Warp 是函数式组合风格的 Web 框架，通过 Filter 组合构建路由。

.. code-block:: toml

   [dependencies]
   warp = "0.3"
   tokio = { version = "1", features = ["full"] }

.. code-block:: rust

   use warp::Filter;

   #[tokio::main]
   async fn main() {
       // GET /
       let hello = warp::path::end()
           .map(|| warp::reply::html("Hello, World!"));

       // GET /hello/<name>
       let greet = warp::path!("hello" / String)
           .map(|name| format!("Hello, {}!", name));

       // GET /search?q=<query>
       let search = warp::path("search")
           .and(warp::query::<std::collections::HashMap<String, String>>())
           .map(|params: std::collections::HashMap<String, String>| {
               let q = params.get("q").map(|s| s.as_str()).unwrap_or("");
               format!("搜索: {}", q)
           });

       // POST /users (JSON)
       #[derive(serde::Deserialize)]
       struct CreateUser {
           name: String,
           email: String,
       }

       let create_user = warp::path("users")
           .and(warp::post())
           .and(warp::body::json::<CreateUser>())
           .map(|user: CreateUser| {
               warp::reply::json(&serde_json::json!({
                   "created": user.name,
               }))
           });

       let routes = hello
           .or(greet)
           .or(search)
           .or(create_user)
           .with(warp::log("api"));

       warp::serve(routes).run(([127, 0, 0, 1], 8080)).await;
   }

Filter 组合模式：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - Filter
     - 说明
   * - ``warp::path("...")``
     - 匹配路径段
   * - ``warp::path::end()``
     - 匹配路径结束（即 "/"）
   * - ``warp::path::param::<T>()``
     - 提取路径参数
   * - ``warp::query::<T>()``
     - 提取查询参数
   * - ``warp::body::json::<T>()``
     - 提取 JSON 请求体
   * - ``warp::body::form::<T>()``
     - 提取表单数据
   * - ``warp::header::<T>("name")``
     - 提取请求头
   * - ``filter.and(other)``
     - 组合两个 filter（同时满足）
   * - ``filter.or(other)``
     - 选择 filter（任一满足）
   * - ``filter.with(wrapper)``
     - 包裹中间件
   * - ``filter.recover(handler)``
     - 错误恢复

Tide
======

Tide 是 async-std 生态的 Web 框架，设计简洁。

.. code-block:: toml

   [dependencies]
   tide = "0.16"
   async-std = { version = "1", features = ["attributes"] }

.. code-block:: rust

   use tide::{Request, Response, StatusCode};

   async fn index(_req: Request<()>) -> tide::Result {
       Ok(Response::builder(200)
           .body("Hello, World!")
           .build())
   }

   async fn greet(req: Request<()>) -> tide::Result {
       let name: String = req.param("name")?.into();
       Ok(format!("Hello, {}!", name).into())
   }

   #[async_std::main]
   async fn main() -> tide::Result<()> {
       let mut app = tide::new();

       app.at("/").get(index);
       app.at("/hello/:name").get(greet);

       app.listen("127.0.0.1:8080").await?;
       Ok(())
   }

框架对比
==========

.. list-table::
   :header-rows: 1
   :widths: 15 15 15 20 15

   * - 特性
     - Actix-web
     - Axum
     - Rocket
     - Warp
   * - 异步运行时
     - actix-rt
     - tokio
     - tokio
     - tokio
   * - HTTP 版本
     - 1.x / 2
     - 1.x / 2
     - 1.x
     - 1.x
   * - WebSocket
     - 原生支持
     - 通过 axum::extract::ws
     - 原生支持
     - 原生支持
   * - 中间件
     - 内置 + 自定义
     - Tower 生态
     - Fairing
     - Filter 包装
   * - 提取器
     - 类型化 Extractors
     - 类型化 Extractors
     - FromRequest
     - Filter 组合
   * - 共享状态
     - web::Data<T>
     - State<T>
     - rocket::State<T>
     - filter.and(with_state)
   * - 测试工具
     - actix_web::test
     - axum_test / tower::ServiceExt
     - rocket::local
     - warp::test
   * - 编译速度
     - 较慢
     - 较快
     - 较慢（宏多）
     - 中等
   * - 学习曲线
     - 中等
     - 较平缓
     - 平缓
     - 较陡（函数式）
   * - 适用场景
     - 高性能 API、全栈
     - API 服务、微服务
     - 全栈应用、快速原型
     - API 网关、组合式服务

选型建议：

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - 场景
     - 推荐
   * - 追求极致性能、大型项目
     - Actix-web
   * - 与 tokio 深度集成、类型安全
     - Axum
   * - 快速原型、开发体验优先
     - Rocket
   * - 函数式编程偏好、Filter 组合
     - Warp
   * - async-std 生态
     - Tide

代码示例
==========

Actix-web 示例
----------------

.. literalinclude:: code/r01_web_framework/actix-web-demo/src/main.rs
  :caption: main.rs
  :language: rust

.. literalinclude:: code/r01_web_framework/actix-web-demo/Cargo.toml
  :caption: Cargo.toml
  :language: toml

Axum 示例
-----------

.. literalinclude:: code/r01_web_framework/axum-demo/src/main.rs
  :caption: main.rs
  :language: rust

.. literalinclude:: code/r01_web_framework/axum-demo/Cargo.toml
  :caption: Cargo.toml
  :language: toml

Tokio Features 解释
---------------------

.. code-block:: toml

   tokio = { version = "1.52.1", features = [
       "rt-multi-thread",
       "macros",
       "net",
       "io-util"
   ] }

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - Feature
     - 作用
   * - ``rt-multi-thread``
     - 启用多线程运行时（生产环境推荐）
   * - ``macros``
     - 启用 ``#[tokio::main]`` 宏
   * - ``net``
     - 网络功能（TcpListener / TcpStream）
   * - ``io-util``
     - I/O 工具函数（copy / split 等）
