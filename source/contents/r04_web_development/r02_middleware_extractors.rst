==========================
中间件与请求处理
==========================

Web 框架中的中间件、提取器、请求处理与响应构造。

.. contents:: 目录
   :depth: 3
   :local:

中间件概述
============

中间件是 Web 请求处理管道中的一个环节，在请求到达路由处理器之前和响应返回客户端之前执行逻辑。

.. code-block::

   请求 → [中间件1] → [中间件2] → [路由处理器] → [中间件2] → [中间件1] → 响应

中间件常见用途：

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - 用途
     - 说明
   * - 日志
     - 记录请求方法、路径、响应时间
   * - 认证/授权
     - 验证用户身份、检查权限
   * - CORS
     - 跨域资源共享策略
   * - 压缩
     - gzip / brotli 响应压缩
   * - 限流
     - 控制请求频率
   * - 请求体大小限制
     - 防止大请求攻击
   * - 超时
     - 请求处理超时控制
   * - 追踪
     - 分布式链路追踪（tracing）

Actix-web 中间件
==================

内置中间件：

.. code-block:: rust

   use actix_web::{App, HttpServer, middleware};

   #[actix_web::main]
   async fn main() -> std::io::Result<()> {
       HttpServer::new(|| {
           App::new()
               // 请求日志
               .wrap(middleware::Logger::default())
               // 响应压缩
               .wrap(middleware::Compress::default())
               // 路径规范化（/path/ → /path）
               .wrap(middleware::NormalizePath::trim())
               // 默认响应头
               .wrap(middleware::DefaultHeaders::new()
                   .add(("X-Version", "1.0")))
               // ... routes
       })
       .bind("127.0.0.1:8080")?
       .run()
       .await
   }

自定义中间件：

.. code-block:: rust

   use actix_web::{
       dev::{Service, ServiceRequest, ServiceResponse, Transform},
       Error, HttpMessage,
   };
   use std::future::{ready, Ready, Future};
   use std::pin::Pin;
   use std::time::Instant;

   // Step 1: 定义中间件工厂
   pub struct RequestTimer;

   // Step 2: 实现 Transform trait
   impl<S, B> Transform<S, ServiceRequest> for RequestTimer
   where
       S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error> + 'static,
       S::Future: 'static,
       B: 'static,
   {
       type Response = ServiceResponse<B>;
       type Error = Error;
       type Transform = RequestTimerMiddleware<S>;
       type InitError = ();
       type Future = Ready<Result<Self::Transform, Self::InitError>>;

       fn new_transform(&self, service: S) -> Self::Future {
           ready(Ok(RequestTimerMiddleware { service }))
       }
   }

   // Step 3: 中间件服务
   pub struct RequestTimerMiddleware<S> {
       service: S,
   }

   impl<S, B> Service<ServiceRequest> for RequestTimerMiddleware<S>
   where
       S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error> + 'static,
       S::Future: 'static,
       B: 'static,
   {
       type Response = ServiceResponse<B>;
       type Error = Error;
       type Future = Pin<Box<dyn Future<Output = Result<Self::Response, Self::Error>>>>;

       fn poll_ready(&self, cx: &mut std::task::Context<'_>) -> std::task::Poll<Result<(), Self::Error>> {
           self.service.poll_ready(cx)
       }

       fn call(&self, req: ServiceRequest) -> Self::Future {
           let start = Instant::now();
           let path = req.path().to_string();
           let method = req.method().to_string();

           let fut = self.service.call(req);

           Box::pin(async move {
               let res = fut.await?;
               let elapsed = start.elapsed();
               println!("{} {} - {:?}", method, path, elapsed);
               Ok(res)
           })
       }
   }

   // 使用
   #[actix_web::main]
   async fn main() -> std::io::Result<()> {
       HttpServer::new(|| {
           App::new()
               .wrap(RequestTimer)
               .route("/", actix_web::web::get().to(|| async { "OK" }))
       })
       .bind("127.0.0.1:8080")?
       .run()
       .await
   }

Axum 中间件（Tower）
======================

Axum 使用 Tower 的 ``Layer`` / ``Service`` trait 体系，与 tower-http 生态无缝集成。

常用 tower-http 中间件：

.. code-block:: rust

   use axum::{Router, routing::get};
   use tower_http::{
       cors::CorsLayer,
       trace::TraceLayer,
       compression::CompressionLayer,
       limit::RequestBodyLimitLayer,
       timeout::TimeoutLayer,
       validate_request::ValidateRequestHeaderLayer,
   };
   use std::time::Duration;

   fn app() -> Router {
       Router::new()
           .route("/", get(|| async { "Hello" }))
           // CORS: 允许所有来源
           .layer(CorsLayer::permissive())
           // 请求追踪日志
           .layer(TraceLayer::new_for_http())
           // 响应压缩
           .layer(CompressionLayer::new())
           // 请求体限制 5MB
           .layer(RequestBodyLimitLayer::new(5 * 1024 * 1024))
           // 请求超时 30s
           .layer(TimeoutLayer::new(Duration::from_secs(30)))
   }

CORS 配置：

.. code-block:: rust

   use tower_http::cors::{CorsLayer, AllowOrigin, AllowMethods, AllowHeaders};
   use http::{Method, header};

   fn cors_layer() -> CorsLayer {
       CorsLayer::new()
           .allow_origin(AllowOrigin::exact("https://example.com".parse().unwrap()))
           .allow_methods(AllowMethods::list([
               Method::GET,
               Method::POST,
               Method::PUT,
               Method::DELETE,
           ]))
           .allow_headers(AllowHeaders::list([
               header::AUTHORIZATION,
               header::CONTENT_TYPE,
           ]))
           .allow_credentials(true)
           .max_age(Duration::from_secs(3600))
   }

自定义 Tower 中间件：

.. code-block:: rust

   use axum::{Router, routing::get};
   use tower::{ServiceBuilder, ServiceExt};
   use tower::layer::layer_fn;
   use std::time::Instant;

   async fn handler() -> &'static str {
       "Hello"
   }

   #[tokio::main]
   async fn main() {
       let app = Router::new()
           .route("/", get(handler))
           .layer(
               ServiceBuilder::new()
                   .layer(TraceLayer::new_for_http())
                   .layer(CompressionLayer::new())
                   .into_inner()
           );

       // ...
   }

   // 使用 layer_fn 创建简单中间件
   fn my_layer<S>(inner: S) -> impl tower::Service<
       http::Request<axum::body::Body>,
       Response = S::Response,
       Error = S::Error,
   >
   where
       S: tower::Service<http::Request<axum::body::Body>>,
   {
       // 在请求前后添加逻辑
       layer_fn(|inner: S| {
           move |mut req: http::Request<axum::body::Body>| {
               let start = Instant::now();
               let uri = req.uri().to_string();

               let fut = inner.call(req);
               async move {
                   let res = fut.await?;
                   let elapsed = start.elapsed();
                   println!("{} - {:?}", uri, elapsed);
                   Ok(res)
               }
           }
       })
   }

提取器（Extractors）
======================

提取器是 Web 框架中从 HTTP 请求中提取数据的机制。

Axum 提取器详解
----------------

.. code-block:: rust

   use axum::{
       extract::{Path, Query, Json, Form, State, Extension, Host, MatchedPath},
       http::{HeaderMap, Method, Uri, header},
       response::Json as JsonResponse,
   };
   use serde::Deserialize;

   // 1. Path: 路径参数
   async fn get_user(Path(user_id): Path<u64>) -> JsonResponse<serde_json::Value> {
       JsonResponse(serde_json::json!({ "user_id": user_id }))
   }

   // 2. Query: 查询参数
   #[derive(Deserialize)]
   struct Pagination {
       page: Option<u32>,
       limit: Option<u32>,
   }

   async fn list_users(Query(pagination): Query<Pagination>) -> JsonResponse<serde_json::Value> {
       let page = pagination.page.unwrap_or(1);
       let limit = pagination.limit.unwrap_or(10).min(100);
       JsonResponse(serde_json::json!({ "page": page, "limit": limit }))
   }

   // 3. Json: JSON 请求体
   #[derive(Deserialize)]
   struct CreateUser {
       name: String,
       email: String,
       #[serde(default)]
       age: Option<u8>,
   }

   async fn create_user(Json(payload): Json<CreateUser>) -> JsonResponse<serde_json::Value> {
       JsonResponse(serde_json::json!({
           "created": payload.name,
           "email": payload.email,
       }))
   }

   // 4. Form: 表单数据
   #[derive(Deserialize)]
   struct LoginForm {
       username: String,
       password: String,
   }

   async fn login(Form(form): Form<LoginForm>) -> JsonResponse<serde_json::Value> {
       // 验证逻辑...
       JsonResponse(serde_json::json!({ "user": form.username }))
   }

   // 5. State: 共享应用状态
   #[derive(Clone)]
   struct AppState {
       db_pool: String,
   }

   async fn health(State(state): State<AppState>) -> JsonResponse<serde_json::Value> {
       JsonResponse(serde_json::json!({ "db": state.db_pool }))
   }

   // 6. Headers: 请求头
   async fn headers_handler(headers: HeaderMap) -> JsonResponse<serde_json::Value> {
       let user_agent = headers
           .get(header::USER_AGENT)
           .and_then(|v| v.to_str().ok())
           .unwrap_or("unknown");

       JsonResponse(serde_json::json!({ "user_agent": user_agent }))
   }

   // 7. 组合提取器
   async fn complex_handler(
       Path(user_id): Path<u64>,
       Query(pagination): Query<Pagination>,
       headers: HeaderMap,
   ) -> JsonResponse<serde_json::Value> {
       JsonResponse(serde_json::json!({
           "user_id": user_id,
           "page": pagination.page,
           "user_agent": headers.get("user-agent").map(|v| v.to_str().unwrap_or("")),
       }))
   }

自定义提取器：

.. code-block:: rust

   use axum::{
       extract::{FromRequest, FromRequestParts},
       http::{request::Parts, StatusCode},
       response::{IntoResponse, Response},
       async_trait,
   };

   // 从请求部分提取
   struct ExtractUserAgent(String);

   #[async_trait]
   impl<S> FromRequestParts<S> for ExtractUserAgent
   where
       S: Send + Sync,
   {
       type Rejection = (StatusCode, &'static str);

       async fn from_request_parts(parts: &mut Parts, _state: &S) -> Result<Self, Self::Rejection> {
           let user_agent = parts
               .headers
               .get(http::header::USER_AGENT)
               .and_then(|v| v.to_str().ok())
               .ok_or((StatusCode::BAD_REQUEST, "缺少 User-Agent 头"))?;

           Ok(ExtractUserAgent(user_agent.to_string()))
       }
   }

   // 从完整请求提取
   struct ExtractAuth(String);

   #[async_trait]
   impl<S> FromRequest<S> for ExtractAuth
   where
       S: Send + Sync,
   {
       type Rejection = Response;

       async fn from_request(req: http::Request<axum::body::Body>, state: &S) -> Result<Self, Self::Rejection> {
           let auth_header = req
               .headers()
               .get(http::header::AUTHORIZATION)
               .and_then(|v| v.to_str().ok())
               .and_then(|v| v.strip_prefix("Bearer "))
               .ok_or_else(|| {
                   (StatusCode::UNAUTHORIZED, "需要认证").into_response()
               })?;

           Ok(ExtractAuth(auth_header.to_string()))
       }
   }

响应构造
==========

Axum 响应类型：

.. code-block:: rust

   use axum::{
       response::{Html, Json, IntoResponse, Response, Redirect, sse::Sse},
       http::StatusCode,
   };
   use serde_json::{json, Value};

   // 多种响应类型
   async fn html_response() -> Html<&'static str> {
       Html("<h1>Hello, World!</h1>")
   }

   async fn json_response() -> Json<Value> {
       Json(json!({ "status": "ok" }))
   }

   async fn text_response() -> &'static str {
       "plain text"
   }

   async fn status_response() -> (StatusCode, &'static str) {
       (StatusCode::NOT_FOUND, "Not Found")
   }

   async fn redirect_response() -> Redirect {
       Redirect::permanent("/new-location")
   }

   async fn custom_response() -> impl IntoResponse {
       (StatusCode::CREATED, [("X-Custom", "value")], "created")
   }

   // 自定义响应类型
   struct AppError {
       code: StatusCode,
       message: String,
   }

   impl IntoResponse for AppError {
       fn into_response(self) -> Response {
           let body = Json(json!({
               "error": self.message,
           }));
           (self.code, body).into_response()
       }
   }

   async fn error_example() -> Result<&'static str, AppError> {
       Err(AppError {
           code: StatusCode::NOT_FOUND,
           message: "资源未找到".to_string(),
       })
   }

错误处理
==========

全局错误处理：

.. code-block:: rust

   use axum::{
       Router, routing::get,
       http::StatusCode,
       response::{IntoResponse, Response},
   };
   use std::convert::Infallible;

   // 错误类型
   #[derive(Debug)]
   enum ApiError {
       NotFound(String),
       BadRequest(String),
       Internal(String),
   }

   impl IntoResponse for ApiError {
       fn into_response(self) -> Response {
           let (status, message) = match self {
               ApiError::NotFound(msg) => (StatusCode::NOT_FOUND, msg),
               ApiError::BadRequest(msg) => (StatusCode::BAD_REQUEST, msg),
               ApiError::Internal(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
           };
           (status, message).into_response()
       }
   }

   // 从 std::error::Error 转换
   impl<E: std::error::Error> From<E> for ApiError {
       fn from(err: E) -> Self {
           ApiError::Internal(err.to_string())
       }
   }

   // 在路由中使用
   async fn get_item() -> Result<String, ApiError> {
       // Ok(item)
       Err(ApiError::NotFound("项目不存在".to_string()))
   }

Actix-web 错误处理：

.. code-block:: rust

   use actix_web::{HttpResponse, ResponseError, http::StatusCode};
   use std::fmt;

   #[derive(Debug)]
   enum AppError {
       NotFound(String),
       BadRequest(String),
       Internal(String),
   }

   impl fmt::Display for AppError {
       fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
           match self {
               AppError::NotFound(msg) => write!(f, "NotFound: {}", msg),
               AppError::BadRequest(msg) => write!(f, "BadRequest: {}", msg),
               AppError::Internal(msg) => write!(f, "Internal: {}", msg),
           }
       }
   }

   impl ResponseError for AppError {
       fn error_response(&self) -> HttpResponse {
           let (status, msg) = match self {
               AppError::NotFound(msg) => (StatusCode::NOT_FOUND, msg.clone()),
               AppError::BadRequest(msg) => (StatusCode::BAD_REQUEST, msg.clone()),
               AppError::Internal(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg.clone()),
           };
           HttpResponse::build(status).json(serde_json::json!({ "error": msg }))
       }
   }

请求验证
==========

使用 validator 进行请求验证：

.. code-block:: toml

   [dependencies]
   validator = { version = "0.18", features = ["derive"] }

.. code-block:: rust

   use validator::Validate;
   use serde::Deserialize;

   #[derive(Debug, Deserialize, Validate)]
   struct CreateUserRequest {
       #[validate(length(min = 1, max = 100, message = "名称长度需在 1-100 之间"))]
       name: String,

       #[validate(email(message = "邮箱格式不正确"))]
       email: String,

       #[validate(range(min = 0, max = 150, message = "年龄需在 0-150 之间"))]
       age: Option<u8>,

       #[validate(url(message = "URL 格式不正确"))]
       website: Option<String>,

       #[validate(length(min = 6, message = "密码至少 6 位"))]
       password: String,
   }

   async fn create_user(
       Json(payload): Json<CreateUserRequest>,
   ) -> Result<Json<serde_json::Value>, (StatusCode, Json<serde_json::Value>)> {
       payload.validate().map_err(|e| {
           (StatusCode::BAD_REQUEST, Json(json!({
               "error": "验证失败",
               "details": e.to_string(),
           })))
       })?;

       Ok(Json(json!({ "created": payload.name })))
   }

请求处理流程总览
==================

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - 阶段
     - 处理内容
   * - 中间件（Before）
     - CORS 检查、认证令牌验证、请求日志记录
   * - 提取器
     - 解析 Path/Query/Json/Form/Headers 到 Rust 类型
   * - 验证
     - 业务规则校验（长度、格式、范围）
   * - 处理器
     - 核心业务逻辑
   * - 响应构造
     - 构建 JSON/HTML/Redirect 等响应
   * - 中间件（After）
     - 响应头注入、压缩、响应日志
