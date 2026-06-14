============================================================
HTTP 客户端
============================================================

Rust 中发送 HTTP 请求的主流客户端库。

.. contents:: 目录
   :depth: 3
   :local:

reqwest
=========

Rust 最流行的 HTTP 客户端，支持同步和异步、TLS、连接池、代理等。

.. code-block:: toml

   [dependencies]
   reqwest = { version = "0.12", features = ["json", "rustls-tls"] }
   tokio = { version = "1", features = ["full"] }

基础请求：

.. code-block:: rust

   use reqwest;

   #[tokio::main]
   async fn main() -> Result<(), reqwest::Error> {
       // GET 请求
       let body = reqwest::get("https://httpbin.org/ip")
           .await?
           .text()
           .await?;
       println!("IP: {}", body);

       // POST 请求
       let client = reqwest::Client::new();
       let resp = client
           .post("https://httpbin.org/post")
           .body("hello world")
           .send()
           .await?;
       println!("状态: {}", resp.status());

       Ok(())
   }

JSON 请求：

.. code-block:: rust

   use reqwest;
   use serde::{Deserialize, Serialize};

   #[derive(Debug, Serialize, Deserialize)]
   struct User {
       name: String,
       email: String,
   }

   #[derive(Debug, Deserialize)]
   struct PostResponse {
       json: User,
   }

   #[tokio::main]
   async fn main() -> Result<(), reqwest::Error> {
       let client = reqwest::Client::new();

       // 发送 JSON
       let user = User {
           name: "Alice".to_string(),
           email: "alice@example.com".to_string(),
       };

       let resp: PostResponse = client
           .post("https://httpbin.org/post")
           .json(&user)
           .send()
           .await?
           .json()
           .await?;

       println!("返回: {:?}", resp.json);

       Ok(())
   }

请求定制：

.. code-block:: rust

   use reqwest::{Client, header};
   use std::time::Duration;

   #[tokio::main]
   async fn main() -> Result<(), reqwest::Error> {
       let client = Client::builder()
           .timeout(Duration::from_secs(30))           // 超时
           .connect_timeout(Duration::from_secs(10))    // 连接超时
           .pool_max_idle_per_host(5)                   // 每主机最大空闲连接
           .user_agent("MyApp/1.0")                     // User-Agent
           .default_headers({
               let mut headers = header::HeaderMap::new();
               headers.insert("X-Custom", "value".parse().unwrap());
               headers
           })
           .cookie_store(true)                          // Cookie 存储
           .gzip(true)                                  // 自动解压
           .brotli(true)
           .build()?;

       let resp = client
           .get("https://httpbin.org/headers")
           .header("Authorization", "Bearer token123")
           .query(&[("page", "1"), ("limit", "10")])
           .send()
           .await?;

       println!("{}", resp.text().await?);

       Ok(())
   }

代理设置：

.. code-block:: rust

   use reqwest::{Client, Proxy};

   #[tokio::main]
   async fn main() -> Result<(), reqwest::Error> {
       // HTTP 代理
       let client = Client::builder()
           .proxy(Proxy::http("http://proxy.example.com:8080")?)
           .build()?;

       // HTTPS 代理
       let client = Client::builder()
           .proxy(Proxy::https("https://proxy.example.com:8443")?)
           .build()?;

       // 代理认证
       let client = Client::builder()
           .proxy(Proxy::all("http://user:pass@proxy.example.com:8080")?)
           .build()?;

       // 环境变量自动代理 (HTTP_PROXY, HTTPS_PROXY, NO_PROXY)
       let client = Client::builder()
           .no_proxy()
           .build()?;

       Ok(())
   }

文件上传：

.. code-block:: rust

   use reqwest::{Client, multipart};

   #[tokio::main]
   async fn main() -> Result<(), reqwest::Error> {
       let client = Client::new();

       // multipart 表单上传
       let form = multipart::Form::new()
           .text("username", "alice")
           .file("avatar", "/path/to/avatar.png")?;

       let resp = client
           .post("https://httpbin.org/post")
           .multipart(form)
           .send()
           .await?;

       // 流式上传
       let file = tokio::fs::File::open("large_file.bin").await.unwrap();
       let resp = client
           .put("https://httpbin.org/put")
           .body(reqwest::Body::wrap_stream(
               tokio_util::codec::FramedRead::new(file, tokio_util::codec::BytesCodec::new())
           ))
           .send()
           .await?;

       Ok(())
   }

响应处理：

.. code-block:: rust

   use reqwest::Client;
   use tokio::io::AsyncWriteExt;

   #[tokio::main]
   async fn main() -> Result<(), reqwest::Error> {
       let client = Client::new();

       // 流式下载（大文件）
       let mut resp = client
           .get("https://example.com/large_file.bin")
           .send()
           .await?;

       let mut file = tokio::fs::File::create("download.bin").await.unwrap();
       while let Some(chunk) = resp.chunk().await? {
           file.write_all(&chunk).await.unwrap();
       }

       // 文本响应
       let text = client.get("https://example.com").send().await?.text().await?;

       // JSON 响应
       #[derive(serde::Deserialize)]
       struct IpResponse { origin: String }
       let ip: IpResponse = client.get("https://httpbin.org/ip").send().await?.json().await?;

       // 字节响应
       let bytes = client.get("https://example.com").send().await?.bytes().await?;

       // 检查状态码
       let resp = client.get("https://example.com").send().await?;
       if resp.status().is_success() {
           println!("请求成功");
       } else if resp.status().is_client_error() {
           println!("客户端错误: {}", resp.status());
       } else if resp.status().is_server_error() {
           println!("服务端错误: {}", resp.status());
       }

       // 获取响应头
       let content_type = resp.headers().get("content-type")
           .and_then(|v| v.to_str().ok());
       println!("Content-Type: {:?}", content_type);

       Ok(())
   }

错误处理：

.. code-block:: rust

   use reqwest::{Client, StatusCode};

   #[tokio::main]
   async fn main() {
       let client = Client::new();

       let result = client
           .get("https://httpbin.org/status/404")
           .send()
           .await;

       match result {
           Ok(resp) => {
               match resp.error_for_status() {
                   Ok(resp) => println!("成功: {}", resp.status()),
                   Err(err) => {
                       if err.status() == Some(StatusCode::NOT_FOUND) {
                           println!("资源不存在");
                       } else {
                           eprintln!("请求错误: {}", err);
                       }
                   }
               }
           }
           Err(err) => {
               if err.is_timeout() {
                   eprintln!("请求超时");
               } else if err.is_connect() {
                   eprintln!("连接失败");
               } else {
                   eprintln!("其他错误: {}", err);
               }
           }
       }
   }

reqwest 常用 API：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - API
     - 说明
   * - ``Client::new()``
     - 创建默认客户端
   * - ``Client::builder()...build()``
     - 自定义构建客户端
   * - ``.get(url)`` / ``.post(url)`` / ``.put(url)`` / ``.delete(url)``
     - HTTP 方法
   * - ``.json(&data)``
     - 发送 JSON 请求体
   * - ``.form(&data)``
     - 发送表单数据
   * - ``.multipart(form)``
     - 发送 multipart 表单
   * - ``.query(&params)``
     - URL 查询参数
   * - ``.header(key, val)``
     - 添加请求头
   * - ``.timeout(dur)``
     - 设置超时
   * - ``.basic_auth(user, pass)``
     - HTTP Basic 认证
   * - ``.bearer_auth(token)``
     - Bearer Token 认证
   * - ``resp.text()`` / ``resp.json::<T>()`` / ``resp.bytes()``
     - 解析响应体

ureq
======

轻量级同步 HTTP 客户端，零异步依赖，适合简单场景。

.. code-block:: toml

   [dependencies]
   ureq = { version = "3", features = ["json"] }

.. code-block:: rust

   use ureq;

   fn main() -> Result<(), ureq::Error> {
       // GET 请求
       let body: String = ureq::get("https://httpbin.org/ip")
           .call()?
           .into_string()?;
       println!("{}", body);

       // POST JSON
       let resp: serde_json::Value = ureq::post("https://httpbin.org/post")
           .set("Content-Type", "application/json")
           .send_json(serde_json::json!({
               "name": "Alice",
               "email": "alice@example.com",
           }))?
           .into_json()?;
       println!("{:?}", resp);

       // 自定义配置
       let agent = ureq::AgentBuilder::new()
           .timeout(std::time::Duration::from_secs(30))
           .user_agent("MyApp/1.0")
           .build();

       let resp = agent
           .get("https://httpbin.org/get")
           .query("page", "1")
           .query("limit", "10")
           .call()?;

       println!("{}", resp.into_string()?);

       Ok(())
   }

reqwest vs ureq：

.. list-table::
   :header-rows: 1
   :widths: 20 40 40

   * - 特性
     - reqwest
     - ureq
   * - 异步支持
     - 原生异步（tokio）
     - 仅同步
   * - 依赖
     - tokio + hyper + TLS
     - 少量依赖
   * - 连接池
     - 内置
     - 内置（Agent）
   * - 代理
     - 内置
     - 内置
   * - HTTP/2
     - 支持
     - 不支持
   * - 编译速度
     - 较慢
     - 较快
   * - 适用场景
     - 异步服务、高并发
     - CLI 工具、简单脚本

hyper
=======

底层 HTTP 实现，reqwest 和 axum 都基于 hyper。适合需要完全控制的场景。

.. code-block:: toml

   [dependencies]
   hyper = { version = "1", features = ["full"] }
   tokio = { version = "1", features = ["full"] }
   http-body-util = "0.1"
   hyper-util = { version = "0.1", features = ["full"] }

客户端：

.. code-block:: rust

   use hyper::{Request, Method};
   use hyper_util::client::legacy::Client;
   use hyper_util::rt::TokioExecutor;
   use http_body_util::BodyExt;

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let client: Client<_, String> = Client::builder(TokioExecutor::new())
           .build_http();

       let req = Request::builder()
           .method(Method::GET)
           .uri("https://httpbin.org/ip")
           .body(String::new())?;

       let resp = client.request(req).await?;
       println!("状态: {}", resp.status());

       let body = resp.collect().await?.to_bytes();
       println!("{}", String::from_utf8_lossy(&body));

       Ok(())
   }

HTTP 客户端选择：

.. list-table::
   :header-rows: 1
   :widths: 15 20 45

   * - 库
     - 层级
     - 适用场景
   * - reqwest
     - 高级
     - 绝大多数场景，功能最全
   * - ureq
     - 高级（同步）
     - 简单脚本、CLI 工具
   * - hyper
     - 低级
     - 构建自定义 HTTP 客户端/服务端
