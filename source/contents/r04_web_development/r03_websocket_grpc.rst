===========================
WebSocket 与 gRPC
===========================

Rust 中实现 WebSocket 实时通信和 gRPC 微服务通信。

.. contents:: 目录
   :depth: 3
   :local:

WebSocket
===========

WebSocket 提供全双工、低延迟的实时通信通道，适用于聊天、实时通知、协作编辑等场景。

Axum WebSocket
----------------

.. code-block:: toml

   [dependencies]
   axum = { version = "0.8", features = ["ws"] }
   tokio = { version = "1", features = ["full"] }
   futures = "0.3"

基础 WebSocket 服务端：

.. code-block:: rust

   use axum::{
       Router, routing::get,
       extract::ws::{WebSocket, WebSocketUpgrade, Message},
       response::IntoResponse,
   };
   use futures::{sink::SinkExt, stream::StreamExt};

   async fn ws_handler(ws: WebSocketUpgrade) -> impl IntoResponse {
       ws.on_upgrade(handle_socket)
   }

   async fn handle_socket(mut socket: WebSocket) {
       // 发送欢迎消息
       let _ = socket
           .send(Message::Text("连接成功！".to_string()))
           .await;

       // 接收并回显消息
       while let Some(Ok(msg)) = socket.recv().await {
           match msg {
               Message::Text(text) => {
                   println!("收到: {}", text);
                   let response = format!("回显: {}", text);
                   if socket.send(Message::Text(response)).await.is_err() {
                       break; // 发送失败，断开连接
                   }
               }
               Message::Binary(data) => {
                   if socket.send(Message::Binary(data)).await.is_err() {
                       break;
                   }
               }
               Message::Ping(data) => {
                   let _ = socket.send(Message::Pong(data)).await;
               }
               Message::Close(_) => {
                   let _ = socket.send(Message::Close(None)).await;
                   break;
               }
               _ => {}
           }
       }

       println!("WebSocket 连接关闭");
   }

   #[tokio::main]
   async fn main() {
       let app = Router::new()
           .route("/ws", get(ws_handler));

       let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
           .await.unwrap();
       axum::serve(listener, app).await.unwrap();
   }

WebSocket 聊天室：

.. code-block:: rust

   use axum::{
       Router, routing::get,
       extract::ws::{WebSocket, WebSocketUpgrade, Message},
       response::IntoResponse,
   };
   use futures::{SinkExt, StreamExt};
   use std::sync::Arc;
   use tokio::sync::{broadcast, Mutex};
   use std::collections::HashSet;

   struct ChatState {
       users: Mutex<HashSet<String>>,
       tx: broadcast::Sender<String>,
   }

   async fn ws_handler(
       ws: WebSocketUpgrade,
       state: axum::extract::State<Arc<ChatState>>,
   ) -> impl IntoResponse {
       ws.on_upgrade(|socket| handle_chat(socket, state))
   }

   async fn handle_chat(socket: WebSocket, state: Arc<ChatState>) {
       let (mut sender, mut receiver) = socket.split();
       let user_id = format!("user_{}", rand::random::<u32>());

       // 加入聊天室
       {
           let mut users = state.users.lock().await;
           users.insert(user_id.clone());
       }

       let mut rx = state.tx.subscribe();

       // 广播加入消息
       let join_msg = format!("{} 加入了聊天室", user_id);
       let _ = state.tx.send(join_msg);

       // 发送消息任务
       let send_state = state.clone();
       let send_user = user_id.clone();
       let mut send_task = tokio::spawn(async move {
           while let Some(Ok(msg)) = receiver.next().await {
               if let Message::Text(text) = msg {
                   let full_msg = format!("{}: {}", send_user, text);
                   let _ = send_state.tx.send(full_msg);
               }
           }
       });

       // 接收广播任务
       let mut recv_task = tokio::spawn(async move {
           while let Ok(msg) = rx.recv().await {
               if sender.send(Message::Text(msg)).await.is_err() {
                   break;
               }
           }
       });

       // 等待任一任务结束
       tokio::select! {
           _ = &mut send_task => recv_task.abort(),
           _ = &mut recv_task => send_task.abort(),
       }

       // 离开聊天室
       let mut users = state.users.lock().await;
       users.remove(&user_id);
       let _ = state.tx.send(format!("{} 离开了聊天室", user_id));
   }

   #[tokio::main]
   async fn main() {
       let (tx, _) = broadcast::channel(100);
       let state = Arc::new(ChatState {
           users: Mutex::new(HashSet::new()),
           tx,
       });

       let app = Router::new()
           .route("/chat", get(ws_handler))
           .with_state(state);

       let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
           .await.unwrap();
       axum::serve(listener, app).await.unwrap();
   }

Actix-web WebSocket
---------------------

.. code-block:: rust

   use actix_web::{web, App, HttpRequest, HttpResponse, HttpServer, Error};
   use actix_ws;

   async fn ws_handler(req: HttpRequest, stream: web::Payload) -> Result<HttpResponse, Error> {
       let (response, mut session, mut msg_stream) = actix_ws::handle(&req, stream)?;

       actix_web::rt::spawn(async move {
           while let Some(Ok(msg)) = msg_stream.next().await {
               match msg {
                   actix_ws::Message::Text(text) => {
                       println!("收到: {}", text);
                       let _ = session.text(format!("回显: {}", text)).await;
                   }
                   actix_ws::Message::Binary(bin) => {
                       let _ = session.binary(bin).await;
                   }
                   actix_ws::Message::Ping(bytes) => {
                       let _ = session.pong(&bytes).await;
                   }
                   actix_ws::Message::Close(reason) => {
                       let _ = session.close(reason).await;
                       break;
                   }
                   _ => {}
               }
           }
       });

       Ok(response)
   }

   #[actix_web::main]
   async fn main() -> std::io::Result<()> {
       HttpServer::new(|| {
           App::new()
               .route("/ws", web::get().to(ws_handler))
       })
       .bind("127.0.0.1:8080")?
       .run()
       .await
   }

WebSocket 消息类型：

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - 消息类型
     - 说明
   * - ``Text(String)``
     - UTF-8 文本消息
   * - ``Binary(Vec<u8>)``
     - 二进制消息
   * - ``Ping(Vec<u8>)``
     - 心跳 ping，应回复 Pong
   * - ``Pong(Vec<u8>)``
     - 心跳 pong 响应
   * - ``Close(Option<CloseFrame>)``
     - 关闭连接请求

gRPC
======

gRPC 是高性能的 RPC 框架，使用 Protocol Buffers 作为接口定义语言。

tonic —— Rust gRPC 框架
--------------------------

.. code-block:: toml

   [dependencies]
   tonic = "0.12"
   prost = "0.13"
   tokio = { version = "1", features = ["full"] }

   [build-dependencies]
   tonic-build = "0.12"

Proto 定义 ``proto/hello.proto``：

.. code-block:: protobuf

   syntax = "proto3";

   package hello;

   service Greeter {
       rpc SayHello (HelloRequest) returns (HelloReply);
       rpc SayHelloStream (HelloRequest) returns (stream HelloReply);
       rpc ChatStream (stream ChatMessage) returns (stream ChatMessage);
   }

   message HelloRequest {
       string name = 1;
   }

   message HelloReply {
       string message = 1;
   }

   message ChatMessage {
       string user = 1;
       string text = 2;
   }

构建脚本 ``build.rs``：

.. code-block:: rust

   fn main() -> Result<(), Box<dyn std::error::Error>> {
       tonic_build::compile_protos("proto/hello.proto")?;
       Ok(())
   }

服务端实现：

.. code-block:: rust

   use tonic::{transport::Server, Request, Response, Status, Streaming};
   use hello::greeter_server::{Greeter, GreeterServer};
   use hello::{HelloRequest, HelloReply, ChatMessage};
   use tokio::sync::mpsc;
   use tokio_stream::wrappers::ReceiverStream;

   pub mod hello {
       tonic::include_proto!("hello");
   }

   #[derive(Default)]
   struct MyGreeter;

   #[tonic::async_trait]
   impl Greeter for MyGreeter {
       // 一元 RPC
       async fn say_hello(
           &self,
           request: Request<HelloRequest>,
       ) -> Result<Response<HelloReply>, Status> {
           let name = request.into_inner().name;
           Ok(Response::new(HelloReply {
               message: format!("Hello, {}!", name),
           }))
       }

       // 服务端流式 RPC
       type SayHelloStreamStream = ReceiverStream<Result<HelloReply, Status>>;

       async fn say_hello_stream(
           &self,
           request: Request<HelloRequest>,
       ) -> Result<Response<Self::SayHelloStreamStream>, Status> {
           let name = request.into_inner().name;
           let (tx, rx) = mpsc::channel(4);

           tokio::spawn(async move {
               for i in 1..=5 {
                   tx.send(Ok(HelloReply {
                       message: format!("Hello, {}! (第 {} 次)", name, i),
                   }))
                   .await
                   .unwrap();
                   tokio::time::sleep(std::time::Duration::from_secs(1)).await;
               }
           });

           Ok(Response::new(ReceiverStream::new(rx)))
       }

       // 双向流式 RPC
       type ChatStreamStream = ReceiverStream<Result<ChatMessage, Status>>;

       async fn chat_stream(
           &self,
           request: Request<Streaming<ChatMessage>>,
       ) -> Result<Response<Self::ChatStreamStream>, Status> {
           let mut in_stream = request.into_inner();
           let (tx, rx) = mpsc::channel(4);

           tokio::spawn(async move {
               while let Some(Ok(msg)) = in_stream.message().await {
                   let reply = ChatMessage {
                       user: "Server".to_string(),
                       text: format!("回显: {}", msg.text),
                   };
                   tx.send(Ok(reply)).await.unwrap();
               }
           });

           Ok(Response::new(ReceiverStream::new(rx)))
       }
   }

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let addr = "[::1]:50051".parse()?;
       let greeter = MyGreeter::default();

       println!("gRPC 服务启动于 {}", addr);

       Server::builder()
           .add_service(GreeterServer::new(greeter))
           .serve(addr)
           .await?;

       Ok(())
   }

客户端实现：

.. code-block:: rust

   use hello::greeter_client::GreeterClient;
   use hello::{HelloRequest, ChatMessage};

   pub mod hello {
       tonic::include_proto!("hello");
   }

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let mut client = GreeterClient::connect("http://[::1]:50051").await?;

       // 一元 RPC
       let request = tonic::Request::new(HelloRequest {
           name: "Rust".to_string(),
       });
       let response = client.say_hello(request).await?;
       println!("响应: {}", response.into_inner().message);

       // 服务端流式
       let request = tonic::Request::new(HelloRequest {
           name: "Stream".to_string(),
       });
       let mut stream = client.say_hello_stream(request).await?.into_inner();
       while let Some(reply) = stream.message().await? {
           println!("流式响应: {}", reply.message);
       }

       // 双向流式
       let outbound = async_stream::stream! {
           for i in 1..=3 {
               yield ChatMessage {
                   user: "Client".to_string(),
                   text: format!("消息 {}", i),
               };
               tokio::time::sleep(std::time::Duration::from_secs(1)).await;
           }
       };

       let response = client.chat_stream(tonic::Request::new(outbound)).await?;
       let mut in_stream = response.into_inner();
       while let Some(msg) = in_stream.message().await? {
           println!("服务端回复: {}", msg.text);
       }

       Ok(())
   }

gRPC 拦截器（Interceptor）：

.. code-block:: rust

   use tonic::{Request, Status, metadata::MetadataValue};

   // 服务端拦截器：认证
   fn auth_interceptor(mut req: Request<()>) -> Result<Request<()>, Status> {
       let token: MetadataValue<_> = "Bearer my-secret-token".parse().unwrap();

       match req.metadata().get("authorization") {
           Some(t) if t == token => Ok(req),
           _ => Err(Status::unauthenticated("无效的认证令牌")),
       }
   }

   // 使用拦截器
   // Server::builder()
   //     .add_service(GreeterServer::with_interceptor(greeter, auth_interceptor))
   //     .serve(addr)
   //     .await?;

RPC 类型对比：

.. list-table::
   :header-rows: 1
   :widths: 20 25 35

   * - RPC 类型
     - 签名
     - 适用场景
   * - 一元 RPC
     - ``rpc Call(Req) returns (Resp)``
     - 简单请求-响应
   * - 服务端流式
     - ``rpc Call(Req) returns (stream Resp)``
     - 大量数据分页推送、日志流
   * - 客户端流式
     - ``rpc Call(stream Req) returns (Resp)``
     - 文件上传、批量处理
   * - 双向流式
     - ``rpc Call(stream Req) returns (stream Resp)``
     - 实时聊天、协作编辑

gRPC vs REST：

.. list-table::
   :header-rows: 1
   :widths: 20 40 40

   * - 特性
     - gRPC
     - REST
   * - 协议
     - HTTP/2
     - HTTP/1.x / HTTP/2
   * - 数据格式
     - Protocol Buffers（二进制）
     - JSON / XML（文本）
   * - 接口定义
     - .proto 文件
     - OpenAPI / 手动定义
   * - 代码生成
     - 强类型客户端/服务端
     - 需要额外工具
   * - 流式传输
     - 原生支持（双向流）
     - SSE / WebSocket 补充
   * - 浏览器支持
     - 需 grpc-web
     - 原生支持
   * - 性能
     - 更高（二进制 + HTTP/2）
     - 较低
   * - 调试
     - 需专用工具（grpcurl）
     - curl / 浏览器直接调试

WebSocket vs gRPC Stream 选择：

.. list-table::
   :header-rows: 1
   :widths: 20 40 40

   * - 特性
     - WebSocket
     - gRPC Stream
   * - 通信模型
     - 全双工消息
     - 请求-响应流
   * - 协议
     - WebSocket (upgrade from HTTP)
     - HTTP/2
   * - 浏览器客户端
     - 原生支持
     - 需 grpc-web 代理
   * - 类型安全
     - 手动序列化
     - 自动生成（Protobuf）
   * - 适用场景
     - 浏览器实时通信、聊天
     - 微服务间通信、流处理
