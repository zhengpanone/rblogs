============================================================
基础网络编程
============================================================

Rust 标准库和底层网络编程基础，涵盖 TCP/UDP、DNS 解析、地址处理等。

.. contents:: 目录
   :depth: 3
   :local:

TCP 编程
==========

Rust 标准库提供了 ``std::net::TcpListener`` 和 ``std::net::TcpStream`` 用于 TCP 通信。

TCP 服务器：

.. code-block:: rust

   use std::io::{Read, Write};
   use std::net::{TcpListener, TcpStream};
   use std::thread;

   fn handle_client(mut stream: TcpStream) {
       let mut buf = [0; 1024];

       loop {
           match stream.read(&mut buf) {
               Ok(0) => {
                   println!("客户端断开连接");
                   break;
               }
               Ok(n) => {
                   // 回显数据
                   if stream.write_all(&buf[..n]).is_err() {
                       break;
                   }
               }
               Err(e) => {
                   eprintln!("读取错误: {}", e);
                   break;
               }
           }
       }
   }

   fn main() -> std::io::Result<()> {
       let listener = TcpListener::bind("127.0.0.1:8080")?;
       println!("TCP 服务器监听于 127.0.0.1:8080");

       for stream in listener.incoming() {
           match stream {
               Ok(stream) => {
                   println!("新连接: {}", stream.peer_addr()?);
                   thread::spawn(move || handle_client(stream));
               }
               Err(e) => eprintln!("连接错误: {}", e),
           }
       }

       Ok(())
   }

TCP 客户端：

.. code-block:: rust

   use std::io::{Read, Write};
   use std::net::TcpStream;

   fn main() -> std::io::Result<()> {
       let mut stream = TcpStream::connect("127.0.0.1:8080")?;
       println!("已连接到服务器");

       // 发送数据
       stream.write_all(b"Hello, Server!")?;

       // 读取响应
       let mut buf = [0; 1024];
       let n = stream.read(&mut buf)?;
       println!("收到: {}", String::from_utf8_lossy(&buf[..n]));

       // 关闭写入端（发送 FIN）
       stream.shutdown(std::net::Shutdown::Write)?;

       Ok(())
   }

TCP 选项配置：

.. code-block:: rust

   use std::net::{TcpStream, TcpListener};
   use std::time::Duration;

   fn configure_socket() -> std::io::Result<()> {
       let listener = TcpListener::bind("127.0.0.1:0")?;

       // 设置 SO_REUSEADDR
       listener.set_nonblocking(false)?;

       let stream = TcpStream::connect("example.com:80")?;

       // 设置读超时
       stream.set_read_timeout(Some(Duration::from_secs(5)))?;

       // 设置写超时
       stream.set_write_timeout(Some(Duration::from_secs(5)))?;

       // 启用 TCP_NODELAY（禁用 Nagle 算法）
       stream.set_nodelay(true)?;

       // 设置 SO_KEEPALIVE
       stream.set_keepalive(Some(Duration::from_secs(60)))?;

       // 获取本地和远程地址
       println!("本地: {}", stream.local_addr()?);
       println!("远程: {}", stream.peer_addr()?);

       // 获取 TTL
       println!("TTL: {}", stream.ttl()?);

       Ok(())
   }

常用 TCP 选项：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 选项
     - 说明
   * - ``set_read_timeout``
     - 读超时，避免阻塞
   * - ``set_write_timeout``
     - 写超时
   * - ``set_nodelay(true)``
     - 禁用 Nagle 算法，减少延迟
   * - ``set_keepalive``
     - TCP Keep-Alive 保活
   * - ``set_ttl``
     - 设置 IP TTL
   * - ``set_nonblocking``
     - 非阻塞模式

UDP 编程
==========

.. code-block:: rust

   use std::net::UdpSocket;

   fn udp_server() -> std::io::Result<()> {
       let socket = UdpSocket::bind("127.0.0.1:8081")?;
       println!("UDP 服务器监听于 127.0.0.1:8081");

       let mut buf = [0; 1024];
       loop {
           let (n, src) = socket.recv_from(&mut buf)?;
           println!("收到 {} 字节来自 {}", n, src);

           // 回显
           socket.send_to(&buf[..n], src)?;
       }
   }

   fn udp_client() -> std::io::Result<()> {
       let socket = UdpSocket::bind("127.0.0.1:0")?;

       socket.send_to(b"Hello, UDP!", "127.0.0.1:8081")?;

       let mut buf = [0; 1024];
       let (n, src) = socket.recv_from(&mut buf)?;
       println!("收到 {} 字节来自 {}: {}", n, src,
                String::from_utf8_lossy(&buf[..n]));

       Ok(())
   }

   // UDP 组播
   fn multicast() -> std::io::Result<()> {
       use std::net::{IpAddr, Ipv4Addr};

       let multicast_addr = Ipv4Addr::new(224, 0, 0, 1);
       let socket = UdpSocket::bind("0.0.0.0:0")?;

       // 加入组播组
       socket.join_multicast_v4(&multicast_addr, &Ipv4Addr::UNSPECIFIED)?;

       // 设置组播 TTL
       socket.set_multicast_ttl_v4(32)?;

       // 发送组播数据
       socket.send_to(b"multicast message", (multicast_addr, 9000))?;

       Ok(())
   }

   // UDP 广播
   fn broadcast() -> std::io::Result<()> {
       let socket = UdpSocket::bind("0.0.0.0:0")?;
       socket.set_broadcast(true)?;

       socket.send_to(b"broadcast message", "255.255.255.255:9000")?;

       Ok(())
   }

地址解析与 DNS
================

.. code-block:: rust

   use std::net::{IpAddr, Ipv4Addr, Ipv6Addr, SocketAddr, ToSocketAddrs};

   fn main() -> std::io::Result<()> {
       // 解析域名
       for addr in "rust-lang.org:443".to_socket_addrs()? {
           println!("rust-lang.org: {}", addr);
       }

       // 解析 IP
       let ip: IpAddr = "192.168.1.1".parse()?;
       println!("IPv4: {}", ip);

       let ip: IpAddr = "::1".parse()?;
       println!("IPv6: {}", ip);

       // SocketAddr 构建
       let addr = SocketAddr::new(IpAddr::V4(Ipv4Addr::new(127, 0, 0, 1)), 8080);
       println!("地址: {}", addr);

       // IP 地址判断
       let ip = Ipv4Addr::new(10, 0, 0, 1);
       println!("私有地址: {}", ip.is_private());
       println!("环回地址: {}", ip.is_loopback());
       println!("链路本地: {}", ip.is_link_local());

       // 网络地址
       let ip = Ipv4Addr::new(192, 168, 1, 100);
       println!("{}", ip.to_string());

       Ok(())
   }

非阻塞与异步 I/O 基础
========================

使用 mio 进行事件驱动编程：

.. code-block:: toml

   [dependencies]
   mio = { version = "1", features = ["net", "os-poll"] }

.. code-block:: rust

   use mio::{Events, Interest, Poll, Token};
   use mio::net::{TcpListener, TcpStream};
   use std::io::{Read, Write};
   use std::collections::HashMap;

   const SERVER: Token = Token(0);

   fn main() -> std::io::Result<()> {
       let mut poll = Poll::new()?;
       let mut events = Events::with_capacity(128);

       // 创建监听器
       let mut listener = TcpListener::bind("127.0.0.1:8080".parse().unwrap())?;
       poll.registry()
           .register(&mut listener, SERVER, Interest::READABLE)?;

       let mut connections: HashMap<Token, TcpStream> = HashMap::new();
       let mut next_token = Token(1);

       println!("mio 服务器监听于 127.0.0.1:8080");

       loop {
           poll.poll(&mut events, None)?;

           for event in events.iter() {
               match event.token() {
                   SERVER => loop {
                       match listener.accept() {
                           Ok((mut stream, addr)) => {
                               println!("新连接: {}", addr);
                               let token = next_token;
                               next_token.0 += 1;

                               poll.registry().register(
                                   &mut stream,
                                   token,
                                   Interest::READABLE,
                               )?;
                               connections.insert(token, stream);
                           }
                           Err(ref e) if e.kind() == std::io::ErrorKind::WouldBlock => break,
                           Err(e) => eprintln!("接受连接错误: {}", e),
                       }
                   },
                   token => {
                       if let Some(mut stream) = connections.remove(&token) {
                           let mut buf = [0; 1024];
                           match stream.read(&mut buf) {
                               Ok(0) => {
                                   println!("连接关闭: {:?}", token);
                               }
                               Ok(n) => {
                                   let msg = String::from_utf8_lossy(&buf[..n]);
                                   println!("收到: {}", msg);
                                   let _ = stream.write_all(b"HTTP/1.1 200 OK\r\n\r\nHello!");
                                   connections.insert(token, stream);
                               }
                               Err(ref e) if e.kind() == std::io::ErrorKind::WouldBlock => {
                                   connections.insert(token, stream);
                               }
                               Err(e) => eprintln!("读取错误: {}", e),
                           }
                       }
                   }
               }
           }
       }
   }

TcpStream 常用方法：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 方法
     - 说明
   * - ``connect(addr)``
     - 连接到远程地址
   * - ``read(&mut buf)``
     - 读取数据
   * - ``write(&buf)`` / ``write_all(&buf)``
     - 写入数据
   * - ``flush()``
     - 刷新写入缓冲区
   * - ``shutdown(Shutdown)``
     - 关闭读/写/双向
   * - ``peer_addr()``
     - 获取远程地址
   * - ``local_addr()``
     - 获取本地地址
   * - ``set_read_timeout``
     - 设置读超时
   * - ``set_nodelay``
     - TCP_NODELAY

UdpSocket 常用方法：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 方法
     - 说明
   * - ``bind(addr)``
     - 绑定本地地址
   * - ``send_to(buf, addr)``
     - 发送到指定地址
   * - ``recv_from(&mut buf)``
     - 接收数据和来源地址
   * - ``connect(addr)``
     - 关联默认远程地址（之后可用 send/recv）
   * - ``set_broadcast(true)``
     - 允许广播
   * - ``join_multicast_v4``
     - 加入 IPv4 组播组
   * - ``set_multicast_ttl_v4``
     - 设置组播 TTL

总结
======

.. list-table::
   :header-rows: 1
   :widths: 20 40 20

   * - 协议
     - 特点
     - 适用场景
   * - TCP
     - 可靠、有序、面向连接
     - HTTP、数据库、文件传输
   * - UDP
     - 不可靠、无序、无连接
     - DNS、视频流、游戏
   * - 组播
     - 一对多传输
     - 视频会议、服务发现
   * - 广播
     - 局域网内所有主机
     - DHCP、服务发现
