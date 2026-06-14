=============================
异步 & 并发
=============================

Rust 生态中异步运行时、并发工具和异步编程基础设施的核心 crate。

.. contents:: 目录
   :depth: 3
   :local:

tokio
==========

异步运行时，Rust 异步编程的事实标准。提供异步 I/O、任务调度、定时器等功能。

.. code-block:: toml

   [dependencies]
   tokio = { version = "1", features = ["full"] }

基础用法：

.. code-block:: rust

   #[tokio::main]
   async fn main() {
       // 并发执行多个任务
       let (r1, r2) = tokio::join!(
           async { tokio::time::sleep(tokio::time::Duration::from_millis(10)).await; 42 },
           async { tokio::time::sleep(tokio::time::Duration::from_millis(5)).await; "hello" },
       );
       println!("{} {}", r1, r2); // 42 hello
   }

spawn 任务：

.. code-block:: rust

   use tokio::task;

   #[tokio::main]
   async fn main() {
       let handle = task::spawn(async {
           // 耗时计算
           let sum: u64 = (0..1_000_000).sum();
           sum
       });

       let result = handle.await.unwrap();
       println!("Sum: {}", result); // Sum: 499999500000
   }

select! 竞速：

.. code-block:: rust

   use tokio::select;
   use tokio::time::{sleep, Duration};

   #[tokio::main]
   async fn main() {
       let slow = async {
           sleep(Duration::from_secs(5)).await;
           "慢任务"
       };

       let fast = async {
           sleep(Duration::from_millis(100)).await;
           "快任务"
       };

       select! {
           result = slow => println!("慢任务完成: {}", result),
           result = fast => println!("快任务完成: {}", result),
       }
       // 输出: 快任务完成: 快任务
   }

tokio::sync 模块 —— 异步同步原语：

.. code-block:: rust

   use tokio::sync::{mpsc, oneshot, broadcast, watch, Mutex, RwLock, Semaphore, Notify, Barrier};

   // mpsc 通道
   #[tokio::main]
   async fn main() {
       let (tx, mut rx) = mpsc::channel(32);

       tokio::spawn(async move {
           for i in 0..10 {
               tx.send(i).await.unwrap();
           }
       });

       while let Some(val) = rx.recv().await {
           println!("收到: {}", val);
       }
   }

   // oneshot 单次通道
   async fn oneshot_example() {
       let (tx, rx) = oneshot::channel();

       tokio::spawn(async {
           let _ = tx.send("一次性消息");
       });

       let msg = rx.await.unwrap();
       println!("{}", msg);
   }

   // broadcast 广播通道
   async fn broadcast_example() {
       let (tx, mut rx1) = broadcast::channel(16);
       let mut rx2 = tx.subscribe();

       tokio::spawn(async move {
           let _ = tx.send("广播消息");
       });

       assert_eq!(rx1.recv().await.unwrap(), "广播消息");
       assert_eq!(rx2.recv().await.unwrap(), "广播消息");
   }

   // 异步 Mutex
   async fn mutex_example() {
       let counter = std::sync::Arc::new(Mutex::new(0));

       let mut handles = vec![];
       for _ in 0..100 {
           let counter = counter.clone();
           handles.push(tokio::spawn(async move {
               let mut num = counter.lock().await;
               *num += 1;
           }));
       }

       for handle in handles {
           handle.await.unwrap();
       }

       println!("Count: {}", *counter.lock().await); // 100
   }

   // Semaphore 信号量
   async fn semaphore_example() {
       let semaphore = std::sync::Arc::new(Semaphore::new(3)); // 最多 3 个并发
       let mut handles = vec![];

       for i in 0..10 {
           let sem = semaphore.clone();
           handles.push(tokio::spawn(async move {
               let _permit = sem.acquire().await.unwrap();
               println!("任务 {} 开始", i);
               tokio::time::sleep(Duration::from_millis(100)).await;
               println!("任务 {} 完成", i);
           }));
       }

       for handle in handles {
           handle.await.unwrap();
       }
   }

tokio::time 模块 —— 定时器与超时：

.. code-block:: rust

   use tokio::time::{sleep, timeout, interval, Duration, Instant};

   #[tokio::main]
   async fn main() {
       // sleep
       sleep(Duration::from_secs(1)).await;

       // timeout 超时控制
       match timeout(Duration::from_millis(100), slow_operation()).await {
           Ok(result) => println!("完成: {}", result),
           Err(_) => println!("操作超时"),
       }

       // interval 定时器
       let mut interval = interval(Duration::from_secs(1));
       for i in 0..5 {
           interval.tick().await;
           println!("第 {} 次触发", i + 1);
       }
   }

   async fn slow_operation() -> String {
       sleep(Duration::from_secs(5)).await;
       "结果".to_string()
   }

tokio::io 模块 —— 异步 I/O：

.. code-block:: rust

   use tokio::io::{self, AsyncReadExt, AsyncWriteExt, AsyncBufReadExt, BufReader};
   use tokio::fs::{self, File};
   use tokio::net::{TcpListener, TcpStream};

   // 异步文件读写
   async fn file_io() -> io::Result<()> {
       let mut file = File::create("test.txt").await?;
       file.write_all(b"Hello, async IO!").await?;

       let mut file = File::open("test.txt").await?;
       let mut contents = String::new();
       file.read_to_string(&mut contents).await?;
       println!("{}", contents); // Hello, async IO!

       // 按行读取
       let file = File::open("test.txt").await?;
       let mut lines = BufReader::new(file).lines();
       while let Some(line) = lines.next_line().await? {
           println!("行: {}", line);
       }

       Ok(())
   }

   // TCP 服务器
   async fn tcp_server() -> io::Result<()> {
       let listener = TcpListener::bind("127.0.0.1:8080").await?;

       loop {
           let (mut socket, addr) = listener.accept().await?;
           println!("新连接: {}", addr);

           tokio::spawn(async move {
               let mut buf = [0; 1024];
               loop {
                   match socket.read(&mut buf).await {
                       Ok(0) => return, // 连接关闭
                       Ok(n) => {
                           if socket.write_all(&buf[..n]).await.is_err() {
                               return; // 写入错误
                           }
                       }
                       Err(_) => return,
                   }
               }
           });
       }
   }

tokio 常用模块一览：

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - 模块
     - 功能
   * - ``tokio::sync``
     - 异步同步原语：mpsc, oneshot, broadcast, watch, Mutex, RwLock, Semaphore, Barrier, Notify
   * - ``tokio::time``
     - sleep, timeout, interval, Instant
   * - ``tokio::io``
     - AsyncReadExt, AsyncWriteExt, BufReader, copy, split
   * - ``tokio::fs``
     - 异步文件操作：read, write, rename, create_dir, remove_file
   * - ``tokio::net``
     - TcpListener, TcpStream, UdpSocket, UnixListener
   * - ``tokio::task``
     - spawn, spawn_blocking, yield_now, JoinSet, LocalSet
   * - ``tokio::signal``
     - 操作系统信号处理（unix::signal, ctrl_c）
   * - ``tokio::process``
     - 异步子进程管理

tokio 运行时配置：

.. code-block:: rust

   #[tokio::main(flavor = "multi_thread", worker_threads = 4)]
   async fn main() {
       // 多线程运行时，4 个工作线程
   }

   // 单线程运行时
   #[tokio::main(flavor = "current_thread")]
   async fn main() {
       // 适用于 I/O 不密集或需要线程本地存储的场景
   }

   // 手动创建运行时
   fn manual_runtime() {
       let rt = tokio::runtime::Runtime::new().unwrap();
       rt.block_on(async {
           println!("手动运行时");
       });
   }

   // 带配置的运行时
   fn configured_runtime() {
       let rt = tokio::runtime::Builder::new_multi_thread()
           .worker_threads(4)
           .max_blocking_threads(512)
           .thread_name("my-worker")
           .enable_all()
           .build()
           .unwrap();

       rt.block_on(async {
           // ...
       });
   }

futures
==========

Rust 异步编程核心 trait 与组合子库。提供 `Future` trait、`Stream` trait 以及丰富的异步组合器。

.. code-block:: toml

   [dependencies]
   futures = "0.3"

核心 trait：

.. code-block:: rust

   use std::future::Future;
   use std::pin::Pin;
   use std::task::{Context, Poll};

   // Future trait（标准库定义，futures 扩展了更多组合子）
   // pub trait Future {
   //     type Output;
   //     fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output>;
   // }

   // Stream trait
   use futures::stream::{self, Stream, StreamExt, TryStreamExt};

   async fn stream_example() {
       let mut stream = stream::iter(vec![1, 2, 3, 4, 5]);

       // 消费 stream
       while let Some(val) = stream.next().await {
           println!("{}", val);
       }
   }

Stream 组合器：

.. code-block:: rust

   use futures::stream::{self, StreamExt};

   async fn stream_combinators() {
       let stream = stream::iter(0..10);

       // filter: 过滤
       let evens: Vec<_> = stream.clone()
           .filter(|x| futures::future::ready(x % 2 == 0))
           .collect()
           .await;
       println!("偶数: {:?}", evens); // [0, 2, 4, 6, 8]

       // map: 映射
       let squares: Vec<_> = stream::iter(1..5)
           .map(|x| async move { x * x })
           .buffered(2) // 最多 2 个并发
           .collect()
           .await;
       println!("平方: {:?}", squares); // [1, 4, 9, 16]

       // fold: 累加
       let sum = stream::iter(1..=100)
           .fold(0u64, |acc, x| async move { acc + x })
           .await;
       println!("Sum: {}", sum); // 5050

       // take / skip / chain / enumerate
       let first_3: Vec<_> = stream::iter(0..10).take(3).collect().await;
       println!("前3: {:?}", first_3);

       // merge: 合并两个 stream
       let s1 = stream::iter(vec![1, 3, 5]);
       let s2 = stream::iter(vec![2, 4, 6]);
       let merged: Vec<_> = s1.merge(s2).collect().await;
       println!("合并: {:?}", merged);
   }

Future 组合器：

.. code-block:: rust

   use futures::future::{self, FutureExt, TryFutureExt};

   async fn future_combinators() {
       // join: 等待所有完成
       let (a, b, c) = future::join3(
           async { 1 },
           async { 2 },
           async { 3 },
       ).await;
       println!("{} {} {}", a, b, c); // 1 2 3

       // join_all: 等待所有 Future 完成
       let futures: Vec<_> = (0..10).map(|i| async move { i * i }).collect();
       let results = future::join_all(futures).await;
       println!("{:?}", results);

       // select: 竞速
       future::select(
           async { tokio::time::sleep(std::time::Duration::from_secs(5)).await; "慢" },
           async { tokio::time::sleep(std::time::Duration::from_millis(10)).await; "快" },
       ).await;

       // select_ok: 等待第一个成功的结果
       // abortable: 可取消的 Future
       let (fut, abort_handle) = future::abortable(long_running_task());
       tokio::spawn(async {
           tokio::time::sleep(std::time::Duration::from_millis(100)).await;
           abort_handle.abort();
       });
       match fut.await {
           Ok(result) => println!("完成: {}", result),
           Err(_) => println!("被取消"),
       }

       // either: 返回两种类型之一
       let result = if rand::random() {
           future::Either::Left(async { "左" })
       } else {
           future::Either::Right(async { "右" })
       };
       println!("{}", result.await);
   }

   async fn long_running_task() -> String {
       tokio::time::sleep(std::time::Duration::from_secs(10)).await;
       "完成".to_string()
   }

Sink trait —— 异步写入：

.. code-block:: rust

   use futures::sink::SinkExt;
   use futures::channel::mpsc;

   async fn sink_example() {
       let (tx, mut rx) = mpsc::channel(10);

       // tx 实现了 Sink trait
       let mut sink = tx.sink_map_err(|e| format!("发送错误: {}", e));

       sink.send("hello").await.unwrap();
       sink.send("world").await.unwrap();

       drop(sink);
       while let Some(msg) = rx.next().await {
           println!("收到: {}", msg);
       }
   }

futures 常用模块一览：

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - 模块
     - 功能
   * - ``futures::future``
     - Future trait 扩展：join, select, join_all, abortable, Either, BoxFuture, OptionFuture
   * - ``futures::stream``
     - Stream trait 扩展：iter, unfold, filter, map, fold, merge, zip, select
   * - ``futures::sink``
     - Sink trait 扩展：send, send_all, flush, close, fanout
   * - ``futures::channel``
     - mpsc (多生产者单消费者), oneshot
   * - ``futures::executor``
     - block_on, ThreadPool, LocalPool
   * - ``futures::io``
     - AsyncReadExt, AsyncWriteExt, BufReader, copy

futures 与 tokio 的关系：

.. list-table::
   :header-rows: 1
   :widths: 30 50

   * - 概念
     - 说明
   * - ``std::future::Future``
     - Rust 标准库定义的核心 trait
   * - ``futures`` crate
     - 提供 Stream / Sink trait、丰富的组合子（join, select, map 等）、channel 等
   * - ``tokio``
     - 异步运行时（reactor + executor），实现 I/O、定时器、任务调度
   * - 关系
     - tokio 提供运行时执行 futures 提供的 Future/Stream；两者互补

crossbeam
==========

高性能并发工具集。提供无锁数据结构、作用域线程、通道等，是 std::sync 的有力补充。

.. code-block:: toml

   [dependencies]
   crossbeam = "0.8"

作用域线程 —— 借用栈上数据：

.. code-block:: rust

   use crossbeam::thread;

   fn main() {
       let data = vec![1, 2, 3, 4, 5, 6, 7, 8];

       let results = thread::scope(|s| {
           let mid = data.len() / 2;

           let handle1 = s.spawn(|_| {
               data[..mid].iter().sum::<i32>()
           });

           let handle2 = s.spawn(|_| {
               data[mid..].iter().sum::<i32>()
           });

           handle1.join().unwrap() + handle2.join().unwrap()
       }).unwrap();

       println!("总和: {}", results); // 36
       // data 仍然可用，scope 保证线程已结束
       println!("原始数据: {:?}", data);
   }

channel —— 多生产者多消费者通道：

.. code-block:: rust

   use crossbeam::channel::{self, select, tick, after};
   use std::time::Duration;
   use std::thread;

   fn main() {
       let (tx, rx) = channel::unbounded();

       // 多生产者
       let tx1 = tx.clone();
       let tx2 = tx.clone();

       thread::spawn(move || {
           for i in 0..5 {
               tx1.send(format!("线程1: {}", i)).unwrap();
               thread::sleep(Duration::from_millis(100));
           }
       });

       thread::spawn(move || {
           for i in 0..5 {
               tx2.send(format!("线程2: {}", i)).unwrap();
               thread::sleep(Duration::from_millis(150));
           }
       });

       drop(tx);

       // 接收所有消息
       for msg in rx {
           println!("{}", msg);
       }
   }

   // select! 多通道选择
   fn select_example() {
       let (tx1, rx1) = channel::unbounded();
       let (tx2, rx2) = channel::unbounded();

       thread::spawn(move || {
           tx1.send("通道1").unwrap();
       });

       thread::spawn(move || {
           thread::sleep(Duration::from_millis(50));
           tx2.send("通道2").unwrap();
       });

       // select! 宏：等待任意一个通道
       select! {
           recv(rx1) -> msg => println!("收到: {}", msg.unwrap()),
           recv(rx2) -> msg => println!("收到: {}", msg.unwrap()),
           default(Duration::from_millis(500)) => println!("超时"),
       }
   }

crossbeam::queue —— 无锁队列：

.. code-block:: rust

   use crossbeam::queue::{ArrayQueue, SegQueue};
   use std::sync::Arc;
   use std::thread;

   fn array_queue_example() {
       // ArrayQueue: 有界无锁队列
       let queue = Arc::new(ArrayQueue::new(100));

       let q_producer = queue.clone();
       let producer = thread::spawn(move || {
           for i in 0..1000 {
               while q_producer.push(i).is_err() {
                   // 队列满，自旋等待
                   thread::yield_now();
               }
           }
       });

       let q_consumer = queue.clone();
       let consumer = thread::spawn(move || {
           let mut sum = 0;
           loop {
               match q_consumer.pop() {
                   Ok(val) => sum += val,
                   Err(_) => {
                       // 队列空，可以等待或重试
                       if producer.is_finished() {
                           break;
                       }
                       thread::yield_now();
                   }
               }
           }
           sum
       });

       producer.join().unwrap();
       println!("总和: {}", consumer.join().unwrap());
   }

   fn seg_queue_example() {
       // SegQueue: 无界无锁队列
       let queue = Arc::new(SegQueue::new());

       let q = queue.clone();
       thread::spawn(move || {
           for i in 0..100 {
               q.push(i);
           }
       });

       let q = queue.clone();
       thread::spawn(move || {
           for i in 100..200 {
               q.push(i);
           }
       });

       thread::sleep(Duration::from_millis(100));
       let mut count = 0;
       while queue.pop().is_ok() {
           count += 1;
       }
       println!("弹出 {} 个元素", count);
   }

crossbeam::atomic —— AtomicCell：

.. code-block:: rust

   use crossbeam::atomic::AtomicCell;
   use std::sync::Arc;
   use std::thread;

   fn main() {
       let counter = Arc::new(AtomicCell::new(0u64));

       let mut handles = vec![];
       for _ in 0..10 {
           let counter = counter.clone();
           handles.push(thread::spawn(move || {
               for _ in 0..100_000 {
                   counter.fetch_add(1);
               }
           }));
       }

       for handle in handles {
           handle.join().unwrap();
       }

       println!("计数: {}", counter.load()); // 1_000_000
   }

crossbeam::utils::sync —— WaitGroup / Parker / ShardedLock：

.. code-block:: rust

   use crossbeam::sync::{Parker, Unparker, WaitGroup, ShardedLock};
   use std::thread;

   // WaitGroup: 等待一组任务完成
   fn waitgroup_example() {
       let wg = WaitGroup::new();

       for i in 0..5 {
           let wg = wg.clone();
           thread::spawn(move || {
               println!("任务 {} 完成", i);
               drop(wg); // 任务完成时 drop
           });
       }

       wg.wait(); // 等待所有任务
       println!("所有任务完成");
   }

   // Parker: 线程挂起与唤醒
   fn parker_example() {
       let parker = Parker::new();
       let unparker = parker.unparker();

       thread::spawn(move || {
           thread::sleep(Duration::from_secs(1));
           println!("唤醒主线程");
           unparker.unpark();
       });

       println!("主线程挂起...");
       parker.park();
       println!("主线程被唤醒");
   }

   // ShardedLock: 分片读写锁（比 RwLock 有更好的读并发性能）
   fn sharded_lock_example() {
       let lock = ShardedLock::new(0);

       // 读
       {
           let guard = lock.read().unwrap();
           println!("读取: {}", *guard);
       }

       // 写
       {
           let mut guard = lock.write().unwrap();
           *guard += 1;
       }
   }

crossbeam::deque —— 工作窃取双端队列：

.. code-block:: rust

   use crossbeam::deque::{Injector, Steal, Stealer, Worker};
   use std::thread;

   fn main() {
       let injector = Injector::new();
       let mut workers = vec![];
       let mut stealers = vec![];

       for _ in 0..4 {
           let worker = Worker::new_lifo();  // LIFO: 局部任务后进先出
           stealers.push(worker.stealer());
           workers.push(worker);
       }

       // 全局队列注入任务
       for i in 0..100 {
           injector.push(i);
       }

       let stealers = std::sync::Arc::new(stealers);

       let mut handles = vec![];
       for (id, worker) in workers.into_iter().enumerate() {
           let stealers = stealers.clone();
           handles.push(thread::spawn(move || {
               let mut sum = 0;

               loop {
                   // 1. 先处理自己队列的任务
                   match worker.pop() {
                       Some(task) => { sum += task; continue; }
                       None => {}
                   }

                   // 2. 从全局队列窃取
                   match injector.steal_batch_and_pop(&worker) {
                       Steal::Success(task) => { sum += task; continue; }
                       Steal::Empty => break,
                       Steal::Retry => continue,
                   }

                   // 3. 从其他 worker 窃取
                   let mut stolen = false;
                   for stealer in stealers.iter() {
                       match stealer.steal_batch_and_pop(&worker) {
                           Steal::Success(task) => {
                               sum += task;
                               stolen = true;
                               break;
                           }
                           Steal::Empty => {}
                           Steal::Retry => {}
                       }
                   }

                   if !stolen {
                       // 尝试从全局队列窃取
                       match injector.steal_batch_and_pop(&worker) {
                           Steal::Success(task) => { sum += task; }
                           Steal::Empty => break,
                           Steal::Retry => continue,
                       }
                   }
               }

               println!("Worker {} 总和: {}", id, sum);
               sum
           }));
       }

       let total: i32 = handles.into_iter().map(|h| h.join().unwrap()).sum();
       println!("总计: {}", total);
   }

crossbeam 常用模块一览：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 模块
     - 功能
   * - ``crossbeam::thread``
     - scope: 作用域线程，允许借用栈上变量
   * - ``crossbeam::channel``
     - 高性能 MPMC 通道：unbounded, bounded, select!, tick, after
   * - ``crossbeam::queue``
     - 无锁队列：ArrayQueue (有界), SegQueue (无界)
   * - ``crossbeam::deque``
     - 工作窃取双端队列：Worker, Stealer, Injector
   * - ``crossbeam::sync``
     - WaitGroup, Parker/Unparker, ShardedLock, TreiberStack
   * - ``crossbeam::atomic``
     - AtomicCell: 原子共享类型
   * - ``crossbeam::utils``
     - Backoff (自旋等待), CachePadded (缓存行对齐)

std::sync 与 crossbeam 对比：

.. list-table::
   :header-rows: 1
   :widths: 30 30 30

   * - 功能
     - std::sync
     - crossbeam
   * - 多生产者通道
     - mpsc (多生产者单消费者)
     - channel (多生产者多消费者)
   * - 作用域线程
     - std::thread::scope (Rust 1.63+)
     - thread::scope (更早支持)
   * - 无锁队列
     - 无
     - ArrayQueue, SegQueue
   * - 工作窃取
     - 无
     - deque::Worker
   * - WaitGroup
     - 无
     - sync::WaitGroup
   * - 原子类型
     - AtomicBool 等基础类型
     - AtomicCell (泛型原子类型)
   * - 选择宏
     - 无
     - select!

总结
==========

异步 & 并发 Crate 总览：

.. list-table::
   :header-rows: 1
   :widths: 15 15 50 20

   * - Crate
     - 定位
     - 核心能力
     - 使用场景
   * - ``tokio``
     - 异步运行时
     - I/O 事件循环、任务调度、定时器、异步同步原语
     - Web 服务器、网络客户端、异步 I/O
   * - ``futures``
     - 异步抽象层
     - Stream / Sink trait、Future 组合子、异步 channel
     - 所有异步代码的通用组合与抽象
   * - ``crossbeam``
     - 并发工具集
     - 无锁数据结构、作用域线程、MPMC 通道、工作窃取
     - CPU 密集型并行计算、高性能数据结构
