==========================================
Async / Await 异步编程 
==========================================

Rust 的异步编程基于 ``async`` / ``await`` 语法，配合 ``Future`` Trait 和异步运行时（如 ``tokio``），实现高效的协作式并发。

核心概念：

- ``async fn`` 返回一个 ``Future``
- ``.await`` 等待 ``Future`` 完成，不阻塞线程
- 异步运行时（tokio / async-std）负责调度执行

.. contents:: 目录
   :depth: 3
   :local:

为什么需要异步
==================

同步 I/O vs 异步 I/O：

.. list-table:: 同步 vs 异步对比
   :header-rows: 1
   :widths: 25 35 40

   * - 特性
     - 同步（多线程）
     - 异步（async/await）
   * - 线程模型
     - 每个连接一个线程
     - 少量线程处理大量连接
   * - 内存开销
     - 高（每线程约 2-8MB 栈）
     - 低（每个任务约几 KB）
   * - 上下文切换
     - 系统级，开销大
     - 用户级协作式，开销小
   * - 适用场景
     - CPU 密集型
     - I/O 密集型
   * - 10000 连接
     - 10000 线程（不现实）
     - 几个线程即可

第一个 async 程序
======================

需要添加依赖（以 tokio 为例）：

.. code-block:: toml

   [dependencies]
   tokio = { version = "1", features = ["full"] }

基本示例：

.. code-block:: rust

   async fn say_hello() {
       println!("你好，异步世界！");
   }

   #[tokio::main]
   async fn main() {
       say_hello().await;
   }

``#[tokio::main]`` 宏将 ``async fn main`` 转换为普通 ``fn main`` 并启动 tokio 运行时。

async fn 的本质：Future
============================

``async fn`` 的返回值是一个实现了 ``Future`` Trait 的类型：

.. code-block:: rust

   use std::future::Future;

   // 这两者是等价的
   async fn foo() -> i32 { 42 }

   fn foo() -> impl Future<Output = i32> {
       async { 42 }
   }

Future Trait 的定义（简化）：

.. code-block:: rust

   pub trait Future {
       type Output;
       fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output>;
   }

   pub enum Poll<T> {
       Ready(T),   // 完成，返回结果
       Pending,    // 未完成，等待下次 poll
   }

执行流程：

.. code-block:: text

   运行时
   │
   ├── poll(future)
   │   ├── Ready(result)  → 返回结果
   │   └── Pending        → 注册 waker，等待事件
   │
   ├── 事件就绪 → waker.wake()
   │
   └── 重新 poll(future) → Ready(result)

await：不阻塞的等待
=========================

``.await`` 是 Rust 异步的核心操作符：

.. code-block:: rust

   async fn fetch_data() -> String {
       // 模拟异步 I/O
       tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
       String::from("数据获取完成")
   }

   async fn process() {
       println!("开始获取数据...");
       let data = fetch_data().await; // 在此暂停，让出线程
       println!("{}", data);
   }

   #[tokio::main]
   async fn main() {
       process().await;
   }

``.await`` 的行为：

.. code-block:: text

   fetch_data().await
   │
   ├── 如果 Future 完成 → 返回结果，继续执行
   │
   └── 如果 Future 未完成 → 暂停当前函数，让出线程给其他任务
        │
        └── 事件就绪 → 恢复执行，获取结果

关键：``.await`` 不阻塞操作系统线程，而是协作式让出 CPU。

并发执行多个 Future
========================

使用 ``tokio::join!`` 并发执行：

.. code-block:: rust

   async fn task_a() -> i32 {
       tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
       1
   }

   async fn task_b() -> i32 {
       tokio::time::sleep(tokio::time::Duration::from_millis(200)).await;
       2
   }

   #[tokio::main]
   async fn main() {
       let (a, b) = tokio::join!(task_a(), task_b());
       println!("a = {}, b = {}", a, b);
       // 总耗时 ≈ 200ms，而非 300ms
   }

串行 vs 并发：

.. code-block:: rust

   #[tokio::main]
   async fn main() {
       // 串行：总耗时 ≈ 300ms
       let a = task_a().await;
       let b = task_b().await;

       // 并发：总耗时 ≈ 200ms（取最长的）
       let (a, b) = tokio::join!(task_a(), task_b());
   }

``tokio::try_join!`` 处理可能失败的任务：

.. code-block:: rust

   async fn fetch_a() -> Result<i32, &'static str> { Ok(1) }
   async fn fetch_b() -> Result<i32, &'static str> { Err("失败") }

   #[tokio::main]
   async fn main() -> Result<(), &'static str> {
       let (a, b) = tokio::try_join!(fetch_a(), fetch_b())?;
       println!("a = {}, b = {}", a, b);
       Ok(())
   }

tokio::select!：竞速执行
=============================

同时等待多个 Future，哪个先完成就用哪个：

.. code-block:: rust

   use tokio::time::{sleep, Duration};

   #[tokio::main]
   async fn main() {
       tokio::select! {
           _ = sleep(Duration::from_secs(1)) => {
               println!("1 秒到了");
           }
           _ = sleep(Duration::from_secs(2)) => {
               println!("2 秒到了");
           }
       }
       // 只输出 "1 秒到了"
   }

取消模式——select + 循环：

.. code-block:: rust

   use tokio::sync::mpsc;
   use tokio::time::{sleep, Duration};

   #[tokio::main]
   async fn main() {
       let (tx, mut rx) = mpsc::channel(32);
       let (shutdown_tx, mut shutdown_rx) = mpsc::channel(1);

       // 工作循环
       let worker = tokio::spawn(async move {
           loop {
               tokio::select! {
                   Some(msg) = rx.recv() => {
                       println!("处理消息: {}", msg);
                   }
                   _ = shutdown_rx.recv() => {
                       println!("收到关闭信号，退出");
                       break;
                   }
               }
           }
       });

       tx.send("hello".to_string()).await.unwrap();
       sleep(Duration::from_millis(100)).await;
       shutdown_tx.send(()).await.unwrap();

       worker.await.unwrap();
   }

tokio::spawn：异步任务
============================

``spawn`` 在运行时中创建新的并发任务：

.. code-block:: rust

   #[tokio::main]
   async fn main() {
       let handle = tokio::spawn(async {
           // 这是一个独立的并发任务
           tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;
           42
       });

       // 主任务继续执行
       println!("等待子任务...");
       let result = handle.await.unwrap(); // 等待子任务完成
       println!("结果: {}", result);
   }

多个 spawn 并发：

.. code-block:: rust

   #[tokio::main]
   async fn main() {
       let mut handles = vec![];

       for i in 0..10 {
           handles.push(tokio::spawn(async move {
               tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
               i * 2
           }));
       }

       for handle in handles {
           let result = handle.await.unwrap();
           println!("结果: {}", result);
       }
   }

异步 Channel：任务间通信
=============================

tokio 提供异步版本的 channel：

.. code-block:: rust

   use tokio::sync::mpsc;

   #[tokio::main]
   async fn main() {
       let (tx, mut rx) = mpsc::channel(32);

       // 生产者
       let producer = tokio::spawn(async move {
           for i in 0..5 {
               tx.send(format!("消息 {}", i)).await.unwrap();
               tokio::time::sleep(tokio::time::Duration::from_millis(50)).await;
           }
       });

       // 消费者
       let consumer = tokio::spawn(async move {
           while let Some(msg) = rx.recv().await {
               println!("收到: {}", msg);
           }
       });

       producer.await.unwrap();
       consumer.await.unwrap();
   }

异步 Mutex：tokio::sync::Mutex
=====================================

异步代码中应使用 tokio 的 Mutex 而非标准库的：

.. code-block:: rust

   use std::sync::Arc;
   use tokio::sync::Mutex;

   #[tokio::main]
   async fn main() {
       let counter = Arc::new(Mutex::new(0));
       let mut handles = vec![];

       for _ in 0..10 {
           let counter = Arc::clone(&counter);
           handles.push(tokio::spawn(async move {
               let mut num = counter.lock().await; // 异步获取锁
               *num += 1;
               // 离开作用域自动释放
           }));
       }

       for handle in handles {
           handle.await.unwrap();
       }

       println!("计数: {}", *counter.lock().await); // 10
   }

为什么不能用 ``std::sync::Mutex``？

.. list-table:: std::Mutex vs tokio::Mutex
   :header-rows: 1
   :widths: 25 40 35

   * - 特性
     - ``std::sync::Mutex``
     - ``tokio::sync::Mutex``
   * - 锁等待
     - 阻塞线程
     - 异步等待（不阻塞线程）
   * - 使用方式
     - ``.lock().unwrap()``
     - ``.lock().await``
   * - 适用场景
     - 同步代码 / 锁持有时间极短
     - 异步代码
   * - 跨 .await 持有
     - 不推荐（阻塞运行时）
     - 可以安全使用

异步 I/O：文件与网络
===========================

文件读写：

.. code-block:: rust

   use tokio::fs::File;
   use tokio::io::AsyncReadExt;

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let mut file = File::open("hello.txt").await?;
       let mut contents = String::new();
       file.read_to_string(&mut contents).await?;
       println!("文件内容: {}", contents);
       Ok(())
   }

文件写入：

.. code-block:: rust

   use tokio::fs::File;
   use tokio::io::AsyncWriteExt;

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let mut file = File::create("output.txt").await?;
       file.write_all(b"Hello, async world!").await?;
       Ok(())
   }

TCP 服务器：

.. code-block:: rust

   use tokio::net::{TcpListener, TcpStream};
   use tokio::io::{AsyncReadExt, AsyncWriteExt};

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let listener = TcpListener::bind("127.0.0.1:8080").await?;
       println!("服务器监听在 8080 端口");

       loop {
           let (mut socket, addr) = listener.accept().await?;
           println!("新连接: {}", addr);

           tokio::spawn(async move {
               let mut buf = [0; 1024];
               loop {
                   match socket.read(&mut buf).await {
                       Ok(0) => break, // 连接关闭
                       Ok(n) => {
                           if socket.write_all(&buf[..n]).await.is_err() {
                               break;
                           }
                       }
                       Err(_) => break,
                   }
               }
           });
       }
   }

HTTP 请求（需要 reqwest）：

.. code-block:: toml

   [dependencies]
   reqwest = { version = "0.12", features = ["json"] }
   tokio = { version = "1", features = ["full"] }

.. code-block:: rust

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let body = reqwest::get("https://httpbin.org/ip")
           .await?
           .text()
           .await?;

       println!("响应: {}", body);
       Ok(())
   }

Stream：异步迭代器
========================

``Stream`` 是异步版本的 ``Iterator``：

.. code-block:: rust

   use tokio_stream::StreamExt;

   #[tokio::main]
   async fn main() {
       // 创建一个间隔流
       let mut interval = tokio::time::interval(tokio::time::Duration::from_millis(100));

       for _ in 0..5 {
           interval.tick().await;
           println!("滴!");
       }
   }

使用 ``tokio_stream`` 处理流：

.. code-block:: rust

   use tokio_stream::StreamExt;

   #[tokio::main]
   async fn main() {
       // 把迭代器转为 Stream
       let stream = tokio_stream::iter(0..10);

       // 异步 map + filter
       let results: Vec<i32> = stream
           .filter(|x| {
               let x = *x;
               async move { x % 2 == 0 }
           })
           .map(|x| async move { x * 2 })
           .collect()
           .await;

       println!("{:?}", results); // [0, 4, 8, 12, 16]
   }

异步运行时对比
====================

.. list-table:: 主流异步运行时
   :header-rows: 1
   :widths: 20 40 40

   * - 运行时
     - 特点
     - 适用场景
   * - ``tokio``
     - 最流行，功能最全，多线程 work-stealing
     - 通用异步编程首选
   * - ``async-std``
     - API 贴近标准库，易上手
     - 学习 / 简单项目
   * - ``smol``
     - 轻量级，模块化
     - 嵌入式 / 资源受限
   * - ``monoio``
     - 基于 io_uring 的高性能 I/O
     - Linux 高性能网络服务
   * - ``embassy``
     - 嵌入式异步运行时
     - 嵌入式设备

错误处理
===============

``Result`` 与 async：

.. code-block:: rust

   async fn might_fail(input: i32) -> Result<i32, &'static str> {
       if input < 0 {
           Err("输入不能为负数")
       } else {
           Ok(input * 2)
       }
   }

   #[tokio::main]
   async fn main() -> Result<(), Box<dyn std::error::Error>> {
       let result = might_fail(42).await?; // ? 操作符正常使用
       println!("结果: {}", result);

       match might_fail(-1).await {
           Ok(v) => println!("{}", v),
           Err(e) => println!("错误: {}", e),
       }

       Ok(())
   }

Pin 与 Future 的关系
============================

``Future`` 被 poll 后可能包含自引用，因此需要 ``Pin`` 保证不移动：

.. code-block:: rust

   use std::pin::Pin;
   use std::future::Future;

   // async fn 编译后的状态机大致如下（简化）
   enum MyFutureState {
       Start,
       Waiting { /* 可能自引用 */ },
       Done,
   }

   // Pin<Box<dyn Future>> 保证 Future 不会在内存中移动
   fn spawn_future(future: impl Future<Output = ()> + Send + 'static) {
       tokio::spawn(future); // tokio::spawn 内部处理了 Pin
   }

日常开发中很少直接接触 ``Pin``，但理解它有助于深入掌握 async。

async 块与闭包
====================

async 块：

.. code-block:: rust

   #[tokio::main]
   async fn main() {
       let data = String::from("hello");

       // async 块
       let future = async {
           println!("{}", data);
           42
       };

       let result = future.await;
       println!("{}", result);
   }

async 闭包（Rust 2024 edition / nightly）：

.. code-block:: rust

   // 当前通常用普通闭包返回 async 块
   let async_closure = || async {
       tokio::time::sleep(tokio::time::Duration::from_millis(10)).await;
       42
   };

   let result = async_closure().await;
   println!("{}", result);

常见异步模式
==================

.. list-table:: 异步设计模式
   :header-rows: 1
   :widths: 25 75

   * - 模式
     - 说明
   * - 请求-响应
     - 客户端发请求，服务端响应（HTTP、RPC）
   * - 发布-订阅
     - 广播消息到多个消费者（broadcast channel）
   * - 扇出-扇入
     - 分发任务到多个 worker，收集结果
   * - 流水线
     - 多个阶段串行处理，阶段间用 channel 连接
   * - 背压（Backpressure）
     - 通过有界 channel 限制生产者速度
   * - 优雅关闭
     - select! + shutdown channel 实现安全退出
   * - 超时重试
     - tokio::time::timeout + 循环重试

超时模式示例：

.. code-block:: rust

   use tokio::time::{timeout, Duration};

   async fn slow_operation() -> i32 {
       tokio::time::sleep(Duration::from_secs(5)).await;
       42
   }

   #[tokio::main]
   async fn main() {
       match timeout(Duration::from_secs(1), slow_operation()).await {
           Ok(result) => println!("完成: {}", result),
           Err(_) => println!("操作超时"),
       }
   }

性能与最佳实践
====================

.. list-table:: 异步最佳实践
   :header-rows: 1
   :widths: 30 70

   * - 建议
     - 说明
   * - 避免阻塞操作
     - 不要在异步代码中调用 ``std::thread::sleep`` 或同步 I/O
   * - 用 ``tokio::task::spawn_blocking``
     - CPU 密集型或同步阻塞操作用此函数隔离
   * - 避免持异步锁跨 .await
     - 可能导致死锁或性能下降
   * - 合理设置 channel 容量
     - 用有界 channel 实现背压
   * - 注意取消安全
     - select! 分支被取消时，确保资源正确释放
   * - 使用 ``Arc<Mutex<T>>`` 而非全局变量
     - 异步代码中共享状态首选

spawn_blocking 示例：

.. code-block:: rust

   #[tokio::main]
   async fn main() {
       // CPU 密集型计算放到 blocking 线程池
       let result = tokio::task::spawn_blocking(|| {
           // 这里是同步代码，运行在独立的线程池
           let mut sum = 0u64;
           for i in 0..1_000_000 {
               sum += i;
           }
           sum
       })
       .await
       .unwrap();

       println!("计算结果: {}", result);
   }

总结
==========

.. code-block:: text

   Async/Await 体系
   │
   ├── 核心概念
   │   ├── async fn       返回 Future
   │   ├── .await         不阻塞等待
   │   └── Future Trait   poll / Ready / Pending
   │
   ├── 运行时 (tokio)
   │   ├── #[tokio::main] 启动运行时
   │   ├── tokio::spawn   创建异步任务
   │   └── 多线程 work-stealing 调度
   │
   ├── 并发原语
   │   ├── tokio::join!   并发等待
   │   ├── tokio::select! 竞速等待
   │   ├── mpsc::channel  异步消息传递
   │   └── tokio::sync::Mutex  异步互斥锁
   │
   ├── I/O
   │   ├── tokio::net     异步网络
   │   ├── tokio::fs      异步文件
   │   └── Stream         异步迭代器
   │
   └── 生态
       ├── reqwest        异步 HTTP 客户端
       ├── axum           异步 Web 框架
       ├── tonic          异步 gRPC
       └── sqlx           异步数据库
