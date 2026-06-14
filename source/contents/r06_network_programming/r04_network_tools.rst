============================================================
网络工具与调试
============================================================

网络抓包、流量分析、代理等实用工具库。

.. contents:: 目录
   :depth: 3
   :local:

pcap / pnet —— 数据包捕获与分析
==================================

pcap 提供数据包捕获功能，pnet 提供数据包构建与解析。

.. code-block:: toml

   [dependencies]
   pcap = "3"
   pnet = "0.35"

pcap 捕获数据包：

.. code-block:: rust

   use pcap::{Capture, Device};

   fn main() -> Result<(), Box<dyn std::error::Error>> {
       // 列出所有网络设备
       for device in Device::list()? {
           println!("设备: {} - {:?}", device.name, device.addresses);
       }

       // 捕获数据包
       let mut cap = Capture::from_device("en0")?
           .promisc(true)       // 混杂模式
           .snaplen(65535)      // 最大捕获长度
           .timeout(1000)       // 超时（毫秒）
           .open()?;

       println!("开始捕获数据包...");

       while let Ok(packet) = cap.next_packet() {
           println!(
               "收到 {} 字节, 时间: {}.{:06}",
               packet.len(),
               packet.header.ts.tv_sec,
               packet.header.ts.tv_usec,
           );

           // 解析以太网帧
           if let Some(eth) = pnet::packet::ethernet::EthernetPacket::new(packet.data) {
               println!("  EtherType: {:?}", eth.get_ethertype());
           }

           // 只处理前 10 个包
           // break;
       }

       Ok(())
   }

   // 保存到 pcap 文件
   fn save_pcap() -> Result<(), Box<dyn std::error::Error>> {
       let mut cap = Capture::from_device("en0")?
           .promisc(true)
           .snaplen(65535)
           .timeout(1000)
           .open()?;

       let mut savefile = cap.savefile("capture.pcap")?;

       for _ in 0..100 {
           let packet = cap.next_packet()?;
           savefile.write(&packet);
       }

       println!("已保存 100 个数据包到 capture.pcap");
       Ok(())
   }

pnet 构建数据包：

.. code-block:: rust

   use pnet::packet::{
       ip::IpNextHeaderProtocols,
       ipv4::MutableIpv4Packet,
       tcp::MutableTcpPacket,
       ethernet::MutableEthernetPacket,
       Packet, MutablePacket,
   };

   fn build_packet() {
       let mut buffer = vec![0u8; 1500];

       // 构建以太网帧
       let mut eth_pkt = MutableEthernetPacket::new(&mut buffer).unwrap();
       // 设置 MAC 地址、EtherType 等

       // 构建 IP 包
       let mut ip_pkt = MutableIpv4Packet::new(eth_pkt.payload_mut()).unwrap();
       ip_pkt.set_version(4);
       ip_pkt.set_ttl(64);
       ip_pkt.set_next_level_protocol(IpNextHeaderProtocols::Tcp);

       // 构建 TCP 包
       let mut tcp_pkt = MutableTcpPacket::new(ip_pkt.payload_mut()).unwrap();
       tcp_pkt.set_source(8080);
       tcp_pkt.set_destination(80);
       tcp_pkt.set_sequence(0);
       tcp_pkt.set_flags(0x002); // SYN

       println!("数据包构建完成: {} 字节", buffer.len());
   }

proxies (reqwest 代理)
=========================

.. code-block:: rust

   use reqwest::{Client, Proxy};

   fn proxy_examples() -> Result<(), reqwest::Error> {
       // HTTP 代理
       let http_proxy = Client::builder()
           .proxy(Proxy::http("http://127.0.0.1:8080")?)
           .build()?;

       // HTTPS 代理
       let https_proxy = Client::builder()
           .proxy(Proxy::https("https://127.0.0.1:8443")?)
           .build()?;

       // 带认证的代理
       let auth_proxy = Client::builder()
           .proxy(Proxy::all("http://user:pass@127.0.0.1:8080")?)
           .build()?;

       // 系统代理
       let system_proxy = Client::builder()
           .no_proxy()
           .build()?;

       // 自定义代理规则
       use reqwest::NoProxy;
       let custom_proxy = Client::builder()
           .proxy(Proxy::custom(|url| {
               if url.host_str() == Some("internal.example.com") {
                   None // 直连
               } else {
                   Some("http://proxy:8080".parse().unwrap())
               }
           }))
           .build()?;

       Ok(())
   }

HTTP 抓包调试
===============

使用 reqwest + 环境变量启用代理抓包：

.. code-block:: bash

   # 使用 mitmproxy / Charles / Fiddler
   export HTTPS_PROXY=http://127.0.0.1:8888
   export HTTP_PROXY=http://127.0.0.1:8888

.. code-block:: rust

   use reqwest::Client;

   #[tokio::main]
   async fn main() -> Result<(), reqwest::Error> {
       // 启用 TLS 调试日志
       std::env::set_var("RUST_LOG", "debug");
       env_logger::init();

       let client = Client::builder()
           .danger_accept_invalid_certs(true) // 允许抓包工具的证书
           .build()?;

       let resp = client.get("https://httpbin.org/get").send().await?;
       println!("{}", resp.text().await?);

       Ok(())
   }

env_logger / tracing 网络日志：

.. code-block:: rust

   #[tokio::main]
   async fn main() {
       // 启用 reqwest / hyper 的调试日志
       std::env::set_var("RUST_LOG", "reqwest=debug,hyper=debug,h2=debug");
       tracing_subscriber::fmt::init();

       let resp = reqwest::get("https://httpbin.org/ip").await.unwrap();
       println!("{}", resp.text().await.unwrap());
   }

tokio::net 工具
=================

.. code-block:: rust

   use tokio::net::{TcpListener, TcpStream, UdpSocket, lookup_host};

   #[tokio::main]
   async fn main() -> std::io::Result<()> {
       // DNS 解析
       for addr in lookup_host("rust-lang.org:80").await? {
           println!("rust-lang.org -> {}", addr);
       }

       // TcpStream 连接
       let mut stream = TcpStream::connect("example.com:80").await?;
       println!("本地: {}, 远程: {}", stream.local_addr()?, stream.peer_addr()?);

       // 获取 TTL
       println!("TTL: {}", stream.ttl()?);

       // TCP_NODELAY
       stream.set_nodelay(true)?;

       Ok(())
   }

socket2 —— 高级 Socket 配置
===============================

.. code-block:: toml

   [dependencies]
   socket2 = "0.5"

.. code-block:: rust

   use socket2::{Socket, Domain, Type, Protocol, SockAddr};
   use std::net::{SocketAddr, Ipv4Addr};

   fn advanced_socket() -> std::io::Result<()> {
       let addr: SocketAddr = (Ipv4Addr::new(127, 0, 0, 1), 8080).into();
       let sock_addr = SockAddr::from(addr);

       let socket = Socket::new(Domain::IPV4, Type::STREAM, Some(Protocol::TCP))?;

       // SO_REUSEADDR
       socket.set_reuse_address(true)?;

       // SO_REUSEPORT
       socket.set_reuse_port(true)?;

       // SO_KEEPALIVE
       socket.set_keepalive(true)?;

       // TCP_NODELAY
       socket.set_nodelay(true)?;

       // 非阻塞
       socket.set_nonblocking(true)?;

       // 缓冲区大小
       socket.set_recv_buffer_size(65536)?;
       socket.set_send_buffer_size(65536)?;

       socket.bind(&sock_addr)?;
       socket.listen(128)?;

       Ok(())
   }

socket2 常用选项：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 选项
     - 说明
   * - ``set_reuse_address``
     - SO_REUSEADDR
   * - ``set_reuse_port``
     - SO_REUSEPORT（多进程绑定同一端口）
   * - ``set_keepalive``
     - SO_KEEPALIVE
   * - ``set_nodelay``
     - TCP_NODELAY
   * - ``set_nonblocking``
     - 非阻塞模式
   * - ``set_recv_buffer_size``
     - 接收缓冲区大小
   * - ``set_send_buffer_size``
     - 发送缓冲区大小
   * - ``set_ttl``
     - IP TTL
   * - ``set_linger``
     - SO_LINGER（关闭时等待）

网络监控与工具总结
====================

.. list-table::
   :header-rows: 1
   :widths: 18 15 45 22

   * - Crate
     - 定位
     - 核心能力
     - 使用场景
   * - ``pcap``
     - 数据包捕获
     - 实时捕获、保存 pcap 文件
     - 网络分析、抓包工具
   * - ``pnet``
     - 数据包构建
     - 构建/解析各层协议数据包
     - 网络扫描、自定义协议
   * - ``socket2``
     - Socket 配置
     - 细粒度 socket 选项控制
     - 需要底层 socket 控制
   * - ``reqwest`` + 代理
     - HTTP 调试
     - 代理配置、抓包配合
     - API 调试
   * - ``env_logger`` / ``tracing``
     - 日志
     - 网络库调试日志
     - 开发调试
