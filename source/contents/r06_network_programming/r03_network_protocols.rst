============================================================
网络协议
============================================================

Rust 中常用网络协议的实现库，包括 HTTP/2、QUIC、WebSocket 客户端、MQTT 等。

.. contents:: 目录
   :depth: 3
   :local:

HTTP/2 与 h2
===============

h2 是 Rust 的 HTTP/2 协议实现，基于 tokio 异步 I/O。

.. code-block:: toml

   [dependencies]
   h2 = "0.4"
   tokio = { version = "1", features = ["full"] }
   bytes = "1"

HTTP/2 服务端：

.. code-block:: rust

   use h2::server;
   use tokio::net::TcpListener;
   use bytes::Bytes;

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let listener = TcpListener::bind("127.0.0.1:8080").await?;

       loop {
           let (socket, _) = listener.accept().await?;

           tokio::spawn(async move {
               let mut conn = server::handshake(socket).await.unwrap();

               while let Some(result) = conn.accept().await {
                   let (request, mut respond) = result.unwrap();
                   println!("收到请求: {:?}", request.uri());

                   let response = http::Response::builder()
                       .status(200)
                       .body(()).unwrap();

                   let mut send = respond.send_response(response, false).unwrap();
                   send.send_data(Bytes::from("Hello, HTTP/2!"), true).unwrap();
               }
           });
       }
   }

HTTP/2 客户端：

.. code-block:: rust

   use h2::client;
   use tokio::net::TcpStream;

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let stream = TcpStream::connect("127.0.0.1:8080").await?;
       let (mut client, conn) = client::handshake(stream).await?;

       tokio::spawn(async move {
           conn.await.unwrap();
       });

       let request = http::Request::builder()
           .uri("https://localhost/")
           .body(()).unwrap();

       let (response, mut stream) = client.ready().await.unwrap().send_request(request, false).unwrap();

       let mut body = response.await.unwrap();
       println!("状态: {}", body.status());

       while let Some(chunk) = body.body_mut().data().await {
           println!("数据: {:?}", chunk.unwrap());
       }

       Ok(())
   }

QUIC (quinn)
==============

QUIC 是 Google 开发的传输层协议，基于 UDP，提供 TLS 1.3 加密和 HTTP/3 基础。

.. code-block:: toml

   [dependencies]
   quinn = "0.11"
   tokio = { version = "1", features = ["full"] }
   rustls = "0.23"
   anyhow = "1"

服务端：

.. code-block:: rust

   use quinn::{Endpoint, Connection};
   use std::sync::Arc;

   #[tokio::main]
   async fn main() -> anyhow::Result<()> {
       // 生成自签名证书（生产环境使用真实证书）
       let cert = rcgen::generate_simple_self_signed(vec!["localhost".into()])?;
       let cert_chain = vec![rustls::Certificate(cert.cert.der().to_vec())];
       let key_der = rustls::PrivateKey(cert.key_pair.serialize_der());

       let mut server_crypto = rustls::ServerConfig::builder()
           .with_safe_defaults()
           .with_no_client_auth()
           .with_single_cert(cert_chain, key_der)?;

       server_crypto.alpn_protocols = vec![b"h3".to_vec()];

       let endpoint = Endpoint::server(
           quinn::ServerConfig::with_crypto(Arc::new(server_crypto)),
           "127.0.0.1:4433".parse()?,
       )?;

       println!("QUIC 服务器监听于 127.0.0.1:4433");

       while let Some(conn) = endpoint.accept().await {
           tokio::spawn(handle_connection(conn));
       }

       Ok(())
   }

   async fn handle_connection(conn: quinn::Connecting) {
       let connection = conn.await.unwrap();
       println!("新连接: {}", connection.remote_address());

       while let Ok(stream) = connection.accept_bi().await {
           tokio::spawn(handle_stream(stream));
       }
   }

   async fn handle_stream(stream: (quinn::SendStream, quinn::RecvStream)) {
       let (mut send, mut recv) = stream;

       let mut buf = vec![0; 1024];
       match recv.read(&mut buf).await {
           Ok(Some(n)) => {
               let msg = String::from_utf8_lossy(&buf[..n]);
               println!("收到: {}", msg);
               send.write_all(b"Hello from QUIC!").await.unwrap();
               send.finish().await.unwrap();
           }
           Ok(None) => println!("流关闭"),
           Err(e) => eprintln!("读取错误: {}", e),
       }
   }

客户端：

.. code-block:: rust

   use quinn::Endpoint;
   use std::sync::Arc;

   #[tokio::main]
   async fn main() -> anyhow::Result<()> {
       let mut roots = rustls::RootCertStore::empty();

       let mut client_crypto = rustls::ClientConfig::builder()
           .with_safe_defaults()
           .with_root_certificates(roots)
           .with_no_client_auth();

       client_crypto.alpn_protocols = vec![b"h3".to_vec()];

       let mut endpoint = Endpoint::client("0.0.0.0:0".parse()?)?;
       endpoint.set_default_client_config(quinn::ClientConfig::new(Arc::new(client_crypto)));

       let connection = endpoint
           .connect("127.0.0.1:4433".parse()?, "localhost")?
           .await?;

       println!("已连接到 QUIC 服务器");

       let (mut send, mut recv) = connection.open_bi().await?;
       send.write_all(b"Hello, QUIC!").await?;
       send.finish().await?;

       let mut buf = vec![0; 1024];
       let n = recv.read(&mut buf).await?.unwrap();
       println!("响应: {}", String::from_utf8_lossy(&buf[..n]));

       connection.close(0u32.into(), b"done");
       endpoint.wait_idle().await;

       Ok(())
   }

QUIC vs TCP + TLS：

.. list-table::
   :header-rows: 1
   :widths: 25 35 40

   * - 特性
     - QUIC
     - TCP + TLS
   * - 传输层
     - UDP
     - TCP
   * - 握手延迟
     - 0-RTT / 1-RTT
     - 1.5-RTT (TCP) + 1-RTT (TLS)
   * - 队头阻塞
     - 无（多路独立流）
     - 有（TCP 层面）
   * - 连接迁移
     - 支持（Connection ID）
     - 不支持（IP:Port 绑定）
   * - 加密
     - 强制 TLS 1.3
     - 可选
   * - 部署难度
     - 较高（UDP 可能被阻断）
     - 低
   * - 适用场景
     - 移动端、弱网、实时通信
     - 通用场景

WebSocket 客户端
===================

tokio-tungstenite —— 异步 WebSocket 客户端：

.. code-block:: toml

   [dependencies]
   tokio-tungstenite = "0.24"
   futures-util = "0.3"
   tokio = { version = "1", features = ["full"] }

.. code-block:: rust

   use tokio_tungstenite::{connect_async, tungstenite::Message};
   use futures_util::{SinkExt, StreamExt};

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let url = "wss://echo.websocket.org";

       let (mut ws_stream, _) = connect_async(url).await?;
       println!("WebSocket 已连接");

       // 发送消息
       ws_stream.send(Message::Text("Hello!".into())).await?;

       // 接收消息
       while let Some(msg) = ws_stream.next().await {
           match msg? {
               Message::Text(text) => println!("收到: {}", text),
               Message::Binary(data) => println!("收到二进制: {} bytes", data.len()),
               Message::Ping(data) => {
                   ws_stream.send(Message::Pong(data)).await?;
               }
               Message::Close(_) => {
                   println!("连接关闭");
                   break;
               }
               _ => {}
           }
       }

       Ok(())
   }

   // 带自定义头的连接
   async fn connect_with_headers() -> Result<(), Box<dyn std::error::Error>> {
       use tokio_tungstenite::tungstenite::http::Request;

       let request = Request::builder()
           .uri("wss://echo.websocket.org")
           .header("Authorization", "Bearer token123")
           .body(())?;

       let (mut ws_stream, _) = connect_async(request).await?;

       ws_stream.send(Message::Text("authenticated".into())).await?;

       Ok(())
   }

MQTT (rumqttc)
================

MQTT 是物联网领域最流行的轻量级消息协议。

.. code-block:: toml

   [dependencies]
   rumqttc = "0.24"
   tokio = { version = "1", features = ["full"] }

.. code-block:: rust

   use rumqttc::{MqttOptions, Client, QoS, Event, Packet};
   use std::time::Duration;

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let mut mqttoptions = MqttOptions::new("rust-client", "localhost", 1883);
       mqttoptions.set_keep_alive(Duration::from_secs(5));
       mqttoptions.set_credentials("user", "password");

       let (mut client, mut connection) = Client::new(mqttoptions, 10);

       // 订阅主题
       client.subscribe("sensors/temperature", QoS::AtLeastOnce)?;
       client.subscribe("sensors/+/status", QoS::AtLeastOnce)?;

       // 发布消息
       client.publish(
           "sensors/temperature",
           QoS::AtLeastOnce,
           false,
           "25.5",
       )?;

       // 处理消息
       tokio::spawn(async move {
           for i in 0.. {
               client.publish(
                   "sensors/temperature",
                   QoS::AtLeastOnce,
                   false,
                   format!("{:.1}", 20.0 + (i as f64 * 0.5)),
               ).unwrap();
               tokio::time::sleep(Duration::from_secs(5)).await;
           }
       });

       // 事件循环
       loop {
           let notification = connection.poll().await?;
           match notification {
               Event::Incoming(Packet::Publish(publish)) => {
                   let payload = String::from_utf8_lossy(&publish.payload);
                   println!(
                       "收到: topic={}, payload={}, qos={:?}",
                       publish.topic, payload, publish.qos,
                   );
               }
               Event::Incoming(Packet::ConnAck(_)) => {
                   println!("已连接到 MQTT Broker");
               }
               Event::Incoming(Packet::PingResp) => {}
               _ => {}
           }
       }
   }

MQTT QoS 级别：

.. list-table::
   :header-rows: 1
   :widths: 15 55

   * - QoS
     - 说明
   * - QoS 0 (AtMostOnce)
     - 最多一次，可能丢失，最快
   * - QoS 1 (AtLeastOnce)
     - 至少一次，可能重复
   * - QoS 2 (ExactlyOnce)
     - 恰好一次，最可靠但最慢

DNS (trust-dns / hickory)
===========================

hickory (原 trust-dns) 是 Rust 实现的 DNS 客户端和服务器。

.. code-block:: toml

   [dependencies]
   hickory-resolver = "0.24"

.. code-block:: rust

   use hickory_resolver::{Resolver, TokioAsyncResolver, config::*};
   use std::net::*;

   // 同步解析
   fn sync_dns() -> Result<(), Box<dyn std::error::Error>> {
       let resolver = Resolver::new(ResolverConfig::default(), ResolverOpts::default())?;

       // A 记录
       let ips = resolver.lookup_ip("rust-lang.org")?;
       for ip in ips {
           println!("rust-lang.org -> {}", ip);
       }

       // 反向解析
       let ip: IpAddr = "104.18.25.175".parse()?;
       let names = resolver.reverse_lookup(ip)?;
       for name in names {
           println!("{} -> {}", ip, name);
       }

       Ok(())
   }

   // 异步解析
   async fn async_dns() -> Result<(), Box<dyn std::error::Error>> {
       let resolver = TokioAsyncResolver::tokio(
           ResolverConfig::default(),
           ResolverOpts::default(),
       );

       let response = resolver.lookup_ip("github.com").await?;
       for addr in response {
           println!("github.com -> {}", addr);
       }

       Ok(())
   }

   // 自定义 DNS 服务器
   fn custom_dns() -> Result<(), Box<dyn std::error::Error>> {
       let mut config = ResolverConfig::new();
       config.add_name_server(NameServerConfig {
           socket_addr: "8.8.8.8:53".parse()?,  // Google DNS
           protocol: Protocol::Udp,
           tls_dns_name: None,
           trust_negative_responses: true,
           bind_addr: None,
       });

       let resolver = Resolver::new(config, ResolverOpts::default())?;
       let ips = resolver.lookup_ip("example.com")?;
       for ip in ips {
           println!("example.com -> {}", ip);
       }

       Ok(())
   }

TLS/SSL (rustls)
==================

rustls 是纯 Rust 实现的 TLS 库，无需 OpenSSL 依赖。

.. code-block:: toml

   [dependencies]
   rustls = "0.23"
   tokio = { version = "1", features = ["full"] }
   tokio-rustls = "0.26"
   webpki-roots = "0.26"

TLS 客户端：

.. code-block:: rust

   use tokio::net::TcpStream;
   use tokio_rustls::TlsConnector;
   use rustls::ClientConfig;
   use std::sync::Arc;

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       // 使用系统根证书
       let mut root_store = rustls::RootCertStore::empty();
       root_store.extend(webpki_roots::TLS_SERVER_ROOTS.iter().cloned());

       let config = ClientConfig::builder()
           .with_root_certificates(root_store)
           .with_no_client_auth();

       let connector = TlsConnector::from(Arc::new(config));
       let domain = "example.com".try_into()?;

       let stream = TcpStream::connect("example.com:443").await?;
       let mut tls_stream = connector.connect(domain, stream).await?;

       // 通过 TLS 流读写
       use tokio::io::{AsyncReadExt, AsyncWriteExt};

       tls_stream.write_all(b"GET / HTTP/1.1\r\nHost: example.com\r\n\r\n").await?;

       let mut buf = vec![0; 4096];
       let n = tls_stream.read(&mut buf).await?;
       println!("{}", String::from_utf8_lossy(&buf[..n]));

       Ok(())
   }

TLS 服务端：

.. code-block:: rust

   use tokio::net::TcpListener;
   use tokio_rustls::TlsAcceptor;
   use rustls::ServerConfig;
   use std::sync::Arc;

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       // 加载证书和密钥
       let certs = rustls_pemfile::certs(&mut std::io::BufReader::new(
           std::fs::File::open("cert.pem")?
       ))
       .collect::<Result<Vec<_>, _>>()?;

       let key = rustls_pemfile::private_key(&mut std::io::BufReader::new(
           std::fs::File::open("key.pem")?
       ))?
       .unwrap();

       let config = ServerConfig::builder()
           .with_no_client_auth()
           .with_single_cert(certs, key)?;

       let acceptor = TlsAcceptor::from(Arc::new(config));

       let listener = TcpListener::bind("127.0.0.1:8443").await?;
       println!("TLS 服务器监听于 127.0.0.1:8443");

       loop {
           let (stream, addr) = listener.accept().await?;
           let acceptor = acceptor.clone();

           tokio::spawn(async move {
               match acceptor.accept(stream).await {
                   Ok(mut tls_stream) => {
                       use tokio::io::AsyncWriteExt;
                       let _ = tls_stream.write_all(b"Hello, TLS!").await;
                       println!("TLS 连接: {}", addr);
                   }
                   Err(e) => eprintln!("TLS 握手失败: {}", e),
               }
           });
       }
   }

协议 Crate 总览：

.. list-table::
   :header-rows: 1
   :widths: 18 15 45

   * - Crate
     - 协议
     - 用途
   * - h2
     - HTTP/2
     - 底层 HTTP/2 实现
   * - quinn
     - QUIC
     - HTTP/3 基础、低延迟传输
   * - tokio-tungstenite
     - WebSocket
     - 异步 WebSocket 客户端
   * - rumqttc
     - MQTT
     - IoT 消息通信
   * - hickory-resolver
     - DNS
     - DNS 解析客户端
   * - rustls
     - TLS
     - 纯 Rust TLS 实现
