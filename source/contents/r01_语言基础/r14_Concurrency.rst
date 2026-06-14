===============================
并发（Concurrency） 
===============================

Rust 的并发模型以**所有权**和**类型系统**为基础，在编译期消除数据竞争，实现"无畏并发（Fearless Concurrency）"。

核心思想：

- ``Send`` Trait：类型可以安全地跨线程转移所有权
- ``Sync`` Trait：类型可以安全地跨线程共享引用
- 编译器在编译期检查所有并发安全问题

.. contents:: 目录
   :depth: 3
   :local:

Rust 并发哲学
==================

传统语言的并发痛点：

.. list-table:: 并发问题对比
   :header-rows: 1
   :widths: 30 35 35

   * - 问题
     - 传统语言
     - Rust
   * - 数据竞争
     - 运行时检测或未检测
     - 编译期杜绝
   * - 死锁
     - 可能发生
     - 可能发生（不保证）
   * - 悬垂指针
     - 常见
     - 编译期杜绝
   * - use-after-free
     - 运行时 bug
     - 编译期杜绝
   * - 忘记解锁
     - 容易发生
     - RAII 自动解锁

Rust 的口号：**"Fearless Concurrency"**——在编译期就帮你消除数据竞争。

线程基础：std::thread
==========================

创建线程：

.. code-block:: rust

   use std::thread;
   use std::time::Duration;

   fn main() {
       let handle = thread::spawn(|| {
           for i in 1..5 {
               println!("子线程: {}", i);
               thread::sleep(Duration::from_millis(10));
           }
       });

       for i in 1..3 {
           println!("主线程: {}", i);
           thread::sleep(Duration::from_millis(10));
       }

       handle.join().unwrap(); // 等待子线程结束
   }

``join()`` 等待线程结束并获取返回值：

.. code-block:: rust

   use std::thread;

   fn main() {
       let handle = thread::spawn(|| {
           // 大量计算
           42
       });

       let result = handle.join().unwrap();
       println!("线程返回值: {}", result); // 42
   }

move 闭包：转移所有权到线程
================================

线程闭包默认借用外部变量，需要用 ``move`` 转移所有权：

.. code-block:: rust

   use std::thread;

   fn main() {
       let v = vec![1, 2, 3];

       let handle = thread::spawn(move || {
           println!("子线程获得 v: {:?}", v); // v 的所有权移入线程
       });

       // println!("{:?}", v); // ❌ 编译错误：v 已被移走
       handle.join().unwrap();
   }

为什么需要 ``move``？

.. code-block:: text

   v 在栈上
   │
   ├── 主线程可能先结束 → v 被释放
   └── 子线程还在用 v    → use-after-free！

   move → 所有权转入子线程 → 安全

Send Trait：跨线程转移所有权
=================================

``Send`` 是标记 Trait，表示类型可以安全地转移所有权到另一个线程：

.. code-block:: rust

   // 大多数类型都是 Send
   let x: i32 = 42;        // i32: Send ✓
   let s: String = ...;    // String: Send ✓
   let v: Vec<i32> = ...;  // Vec<i32>: Send ✓

   // 反例
   let rc: Rc<i32> = ...;  // Rc<i32>: !Send（非原子引用计数）

编译器自动检查：

.. code-block:: rust

   use std::rc::Rc;
   use std::thread;

   fn main() {
       let rc = Rc::new(42);

       let handle = thread::spawn(move || {
           println!("{}", rc); // ❌ 编译错误：Rc<i32> 不是 Send
       });

       handle.join().unwrap();
   }

错误信息会明确指出 ``Rc<i32>`` 不能在线程间安全发送。

Sync Trait：跨线程共享引用
================================

``Sync`` 表示类型的不可变引用可以安全地在线程间共享：

.. code-block:: text

   // 基本类型都是 Sync
   &i32: Sync ✓
   &String: Sync ✓

   // 反例
   &RefCell<i32>: !Sync    // 内部可变性，运行时检查非线程安全
   &Cell<i32>: !Sync       // 同上
   &Rc<i32>: !Sync         // 非原子计数

Send 和 Sync 的关系：

.. list-table:: Send vs Sync
   :header-rows: 1
   :widths: 20 40 40

   * - 概念
     - 含义
     - 典型例子
   * - ``Send``
     - 所有权可跨线程转移
     - ``i32``, ``String``, ``Vec<T>``, ``Arc<T>``
   * - ``!Send``
     - 不能跨线程转移所有权
     - ``Rc<T>``, ``*const T``, ``*mut T``
   * - ``Sync``
     - 不可变引用可跨线程共享
     - ``i32``, ``Mutex<T>``, ``Arc<T>``
   * - ``!Sync``
     - 不可变引用不能跨线程共享
     - ``RefCell<T>``, ``Cell<T>``, ``Rc<T>``

Arc\<T\>：多线程共享所有权
================================

``Arc<T>`` 是线程安全的引用计数智能指针（详见智能指针章节）：

.. code-block:: rust

   use std::sync::Arc;
   use std::thread;

   fn main() {
       let data = Arc::new(vec![1, 2, 3, 4, 5]);
       let mut handles = vec![];

       for i in 0..3 {
           let data = Arc::clone(&data); // 增加引用计数
           let handle = thread::spawn(move || {
               println!("线程 {}: 数据长度 = {}", i, data.len());
           });
           handles.push(handle);
       }

       for handle in handles {
           handle.join().unwrap();
       }

       println!("所有线程完成，data 仍可用: {:?}", data);
   }

Mutex\<T\>：互斥锁
========================

``Mutex<T>`` 保证同一时刻只有一个线程访问数据：

.. code-block:: rust

   use std::sync::{Arc, Mutex};
   use std::thread;

   fn main() {
       let counter = Arc::new(Mutex::new(0));
       let mut handles = vec![];

       for _ in 0..10 {
           let counter = Arc::clone(&counter);
           let handle = thread::spawn(move || {
               let mut num = counter.lock().unwrap();
               *num += 1;
               // MutexGuard 离开作用域，自动解锁
           });
           handles.push(handle);
       }

       for handle in handles {
           handle.join().unwrap();
       }

       println!("最终计数: {}", *counter.lock().unwrap()); // 10
   }

``MutexGuard`` 的 RAII 行为：

.. code-block:: rust

   let mut num = counter.lock().unwrap();
   *num += 1;
   // num (MutexGuard) 在此 drop → 自动 unlock

即使发生 panic，MutexGuard 也会被 drop，锁不会"毒化"（但 Mutex 会被标记为 poisoned）。

Poisoned Mutex：

.. code-block:: rust

   let mutex = Arc::new(Mutex::new(0));

   let handle = thread::spawn({
       let mutex = Arc::clone(&mutex);
       move || {
           let _guard = mutex.lock().unwrap();
           panic!("线程 panic！");
       }
   });

   handle.join().unwrap_err(); // 线程 panic

   // 此时 mutex 是 poisoned 状态
   let result = mutex.lock(); // Err(PoisonError)
   match result {
       Ok(_) => println!("未中毒"),
       Err(poisoned) => {
           let data = poisoned.into_inner();
           println!("恢复数据: {}", data);
       }
   }

RwLock\<T\>：读写锁
=========================

允许多个读者同时访问，写者独占：

.. code-block:: rust

   use std::sync::{Arc, RwLock};
   use std::thread;

   fn main() {
       let data = Arc::new(RwLock::new(0));
       let mut handles = vec![];

       // 多个读者
       for i in 0..5 {
           let data = Arc::clone(&data);
           handles.push(thread::spawn(move || {
               let value = data.read().unwrap();
               println!("读者 {}: 值 = {}", i, *value);
           }));
       }

       // 一个写者
       let data = Arc::clone(&data);
       handles.push(thread::spawn(move || {
           let mut value = data.write().unwrap();
           *value = 100;
           println!("写者: 修改值为 {}", *value);
       }));

       for handle in handles {
           handle.join().unwrap();
       }

       println!("最终值: {}", *data.read().unwrap());
   }

Condvar：条件变量
========================

用于线程间的通知和等待：

.. code-block:: rust

   use std::sync::{Arc, Mutex, Condvar};
   use std::thread;

   fn main() {
       let pair = Arc::new((Mutex::new(false), Condvar::new()));
       let pair2 = Arc::clone(&pair);

       // 等待线程
       let waiter = thread::spawn(move || {
           let (lock, cvar) = &*pair2;
           let mut started = lock.lock().unwrap();
           while !*started {
               started = cvar.wait(started).unwrap();
           }
           println!("收到通知，开始执行！");
       });

       // 通知线程
       thread::sleep(std::time::Duration::from_millis(100));
       let (lock, cvar) = &*pair;
       let mut started = lock.lock().unwrap();
       *started = true;
       cvar.notify_one();

       waiter.join().unwrap();
   }

Barrier：屏障同步
========================

让多个线程在某个点汇合：

.. code-block:: rust

   use std::sync::{Arc, Barrier};
   use std::thread;

   fn main() {
       let n = 5;
       let barrier = Arc::new(Barrier::new(n));
       let mut handles = vec![];

       for i in 0..n {
           let barrier = Arc::clone(&barrier);
           handles.push(thread::spawn(move || {
               println!("线程 {} 到达屏障前", i);
               // 做一些前期工作
               barrier.wait(); // 等待所有线程到达
               println!("线程 {} 通过屏障，继续执行", i);
           }));
       }

       for handle in handles {
           handle.join().unwrap();
       }
   }

Channel（mpsc）：消息传递
===============================

Rust 标准库提供多生产者单消费者（mpsc）通道：

基本使用：

.. code-block:: rust

   use std::sync::mpsc;
   use std::thread;

   fn main() {
       let (tx, rx) = mpsc::channel();

       thread::spawn(move || {
           tx.send("你好").unwrap();
           tx.send("世界").unwrap();
       });

       // rx 迭代接收
       for received in rx {
           println!("收到: {}", received);
       }
       // 当 tx 被 drop 时，rx 自动结束迭代
   }

多个生产者：

.. code-block:: rust

   use std::sync::mpsc;
   use std::thread;

   fn main() {
       let (tx, rx) = mpsc::channel();
       let tx2 = tx.clone(); // 第二个生产者

       thread::spawn(move || {
           for i in 0..5 {
               tx.send(format!("线程1-消息{}", i)).unwrap();
               thread::sleep(std::time::Duration::from_millis(10));
           }
       });

       thread::spawn(move || {
           for i in 0..5 {
               tx2.send(format!("线程2-消息{}", i)).unwrap();
               thread::sleep(std::time::Duration::from_millis(10));
           }
       });

       for received in rx {
           println!("{}", received);
       }
   }

同步通道（有界通道）：

.. code-block:: rust

   use std::sync::mpsc;

   fn main() {
       let (tx, rx) = mpsc::sync_channel(2); // 容量为 2

       tx.send(1).unwrap();
       tx.send(2).unwrap();
       // tx.send(3).unwrap(); // 阻塞，直到 rx 取走一个消息

       println!("{}", rx.recv().unwrap());
       println!("{}", rx.recv().unwrap());
   }

``recv()`` vs ``try_recv()``：

.. list-table:: Channel 接收方法
   :header-rows: 1
   :widths: 25 75

   * - 方法
     - 行为
   * - ``rx.recv()``
     - 阻塞等待，直到收到消息（或 tx 关闭）
   * - ``rx.try_recv()``
     - 不阻塞，立即返回 ``Result<T, TryRecvError>``
   * - ``rx.iter()``
     - 迭代器，阻塞接收直到 tx 关闭

Atomic：原子类型
========================

对于简单的共享数据，原子操作比 Mutex 更高效：

.. code-block:: rust

   use std::sync::atomic::{AtomicBool, AtomicI32, Ordering};
   use std::sync::Arc;
   use std::thread;

   fn main() {
       // 原子布尔：线程标志
       let running = Arc::new(AtomicBool::new(true));
       let running_clone = Arc::clone(&running);

       let handle = thread::spawn(move || {
           while running_clone.load(Ordering::SeqCst) {
               println!("工作中...");
               thread::sleep(std::time::Duration::from_millis(50));
           }
           println!("收到停止信号");
       });

       thread::sleep(std::time::Duration::from_millis(200));
       running.store(false, Ordering::SeqCst);

       handle.join().unwrap();
   }

原子计数器：

.. code-block:: rust

   use std::sync::atomic::{AtomicI32, Ordering};
   use std::sync::Arc;
   use std::thread;

   fn main() {
       let counter = Arc::new(AtomicI32::new(0));
       let mut handles = vec![];

       for _ in 0..10 {
           let counter = Arc::clone(&counter);
           handles.push(thread::spawn(move || {
               for _ in 0..100 {
                   counter.fetch_add(1, Ordering::SeqCst);
               }
           }));
       }

       for handle in handles {
           handle.join().unwrap();
       }

       println!("最终计数: {}", counter.load(Ordering::SeqCst)); // 1000
   }

内存排序（Ordering）：

.. list-table:: 常用内存排序
   :header-rows: 1
   :widths: 30 70

   * - Ordering
     - 含义
   * - ``Relaxed``
     - 只保证原子性，不保证顺序。最高性能
   * - ``Acquire``
     - 读操作，后续读写不能重排到此操作之前
   * - ``Release``
     - 写操作，之前的读写不能重排到此操作之后
   * - ``AcqRel``
     - 兼具 Acquire 和 Release
   * - ``SeqCst``
     - 顺序一致性，最严格的保证，默认推荐

Send + Sync 自动推导
============================

编译器为大多数类型自动实现 Send 和 Sync：

.. list-table:: Send/Sync 自动推导规则
   :header-rows: 1
   :widths: 35 65

   * - 类型
     - Send/Sync 状态
   * - 基本类型 (i32, bool, f64 等)
     - Send + Sync
   * - ``&T``
     - T: Sync → &T: Send
   * - ``&mut T``
     - T: Send → &mut T: Send
   * - ``Box<T>``
     - T: Send → Box\<T\>: Send
   * - ``Vec<T>``
     - T: Send → Vec\<T\>: Send
   * - ``Arc<T>``
     - T: Send + Sync → Arc\<T\>: Send + Sync
   * - ``Mutex<T>``
     - T: Send → Mutex\<T\>: Send + Sync
   * - ``Rc<T>``
     - 永远不是 Send 或 Sync

手动标记 unsafe impl（极少使用）：

.. code-block:: rust

   struct MyType {
       // 内部使用了裸指针，但保证线程安全
       ptr: *mut i32,
   }

   unsafe impl Send for MyType {}
   unsafe impl Sync for MyType {}

thread_local：线程局部存储
=================================

每个线程拥有独立的数据副本：

.. code-block:: rust

   use std::cell::RefCell;
   use std::thread;

   thread_local! {
       static COUNTER: RefCell<u32> = RefCell::new(0);
   }

   fn main() {
       COUNTER.with(|c| {
           *c.borrow_mut() = 42;
       });

       let handle = thread::spawn(|| {
           COUNTER.with(|c| {
               println!("子线程 COUNTER: {}", *c.borrow()); // 0，独立副本
           });
       });

       COUNTER.with(|c| {
           println!("主线程 COUNTER: {}", *c.borrow()); // 42
       });

       handle.join().unwrap();
   }

Rayon：数据并行
======================

Rayon 是 Rust 生态中最流行的并行计算库，通过迭代器接口实现数据并行：

.. code-block:: rust

   // Cargo.toml 添加: rayon = "1"
   use rayon::prelude::*;

   fn main() {
       let numbers: Vec<i32> = (0..1_000_000).collect();

       // 并行求和
       let sum: i32 = numbers.par_iter().sum();
       println!("并行求和: {}", sum);

       // 并行 map
       let squares: Vec<i32> = numbers.par_iter().map(|x| x * x).collect();

       // 并行过滤
       let evens: Vec<i32> = numbers.par_iter().filter(|&&x| x % 2 == 0).cloned().collect();
       println!("偶数个数: {}", evens.len());
   }

并发模式与最佳实践
========================

.. list-table:: 并发模式选择指南
   :header-rows: 1
   :widths: 30 70

   * - 场景
     - 推荐方案
   * - 简单计数器 / 标志位
     - ``Atomic`` 类型
   * - 多线程共享只读数据
     - ``Arc<T>``
   * - 多线程共享可变数据
     - ``Arc<Mutex<T>>`` 或 ``Arc<RwLock<T>>``
   * - 生产者-消费者
     - ``mpsc::channel``
   * - 线程间通知
     - ``Condvar`` 或 ``tokio::sync::Notify``
   * - 多阶段同步
     - ``Barrier``
   * - 数据并行计算
     - ``rayon`` 库
   * - 大规模并发 I/O
     - ``tokio`` / ``async-std`` （异步）


避免常见陷阱：

.. list-table:: 并发常见陷阱
   :header-rows: 1
   :widths: 30 70

   * - 陷阱
     - 解决方案
   * - 持锁时间过长
     - 缩小锁的作用域，尽快 drop MutexGuard
   * - 死锁
     - 统一加锁顺序；用 ``try_lock()``
   * - ``Rc<T>`` 跨线程
     - 改用 ``Arc<T>``
   * - ``RefCell<T>`` 跨线程
     - 改用 ``Mutex<T>`` 或 ``RwLock<T>``
   * - 忘记 join
     - 确保所有 ``JoinHandle`` 都被处理
   * - Channel 发送端提前 drop
     - 确保 tx 的生命周期足够长

死锁示例与避免
====================

经典死锁：

.. code-block:: rust

   use std::sync::{Arc, Mutex};
   use std::thread;

   fn main() {
       let a = Arc::new(Mutex::new(0));
       let b = Arc::new(Mutex::new(0));

       let a1 = Arc::clone(&a);
       let b1 = Arc::clone(&b);
       let a2 = Arc::clone(&a);
       let b2 = Arc::clone(&b);

       let t1 = thread::spawn(move || {
           let _guard_a = a1.lock().unwrap();
           thread::sleep(std::time::Duration::from_millis(100));
           let _guard_b = b1.lock().unwrap(); // 等待 b
       });

       let t2 = thread::spawn(move || {
           let _guard_b = b2.lock().unwrap();
           thread::sleep(std::time::Duration::from_millis(100));
           let _guard_a = a2.lock().unwrap(); // 等待 a
       });

       t1.join().unwrap();
       t2.join().unwrap();
       // ❌ 死锁：t1 持有 a 等 b，t2 持有 b 等 a
   }

避免方案——统一加锁顺序：

.. code-block:: rust

   // 始终按相同顺序加锁：先 a 后 b
   // t1: lock(a) → lock(b) ✓
   // t2: lock(a) → lock(b) ✓ (先释放了 b 再锁 a，不会死锁)

总结
=====

.. code-block:: text

   并发体系
   │
   ├── 线程基础
   │   ├── thread::spawn      创建线程
   │   ├── move 闭包          所有权转移
   │   └── JoinHandle::join   等待线程
   │
   ├── 类型安全（编译期）
   │   ├── Send Trait         可跨线程转移所有权
   │   └── Sync Trait         可跨线程共享引用
   │
   ├── 共享状态
   │   ├── Arc<T>             多线程共享所有权
   │   ├── Mutex<T>           互斥访问
   │   ├── RwLock<T>          读写锁
   │   ├── Atomic             原子操作（无锁）
   │   └── Condvar            条件变量
   │
   ├── 消息传递
   │   └── mpsc::channel      "不要共享内存，而是通信"
   │
   ├── 同步原语
   │   ├── Barrier            屏障同步
   │   └── thread_local       线程局部存储
   │
   └── 生态
       └── rayon              数据并行
