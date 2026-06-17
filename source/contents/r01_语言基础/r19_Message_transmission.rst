===================
消息传递
===================

消息传递的核心思想是：线程之间不直接共享数据，而是通过发送和接收消息来通信。

channel（通道）
====================

channel 是消息传递的基本单元，它是线程之间通信的管道。

- 生产者调用 send() 把数据放进管道
- 消费者调用 recv() 从管道取出数据
- 数据的所有权被转移，不是复制

mpsc（多生产者单消费者）
==========================

mpsc 是 multiple producer, single consumer（多生产者单消费者）的缩写。

- 多生产者：可以有多个线程发送消息
- 单消费者：只有一个线程接收消息

为什么是单消费者？因为如果多个消费者同时从同一个 channel 收消息，可能会导致消息被重复处理或者丢失。


单生产者单消费者
>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust
  :caption: 单生产者单消费者
  :name: single-producer-single-consumer

  use std::sync::mpsc;
  use std::thread;

  fn main() {
      // 1. 创建 channel
      // tx = transmitter (发送端)
      // rx = receiver (接收端)
      let (tx, rx) = mpsc::channel();
      
      // 2. 创建生产者线程
      thread::spawn(move || {
          let msg = String::from("你好，消费者！");
          tx.send(msg).unwrap();
          // msg 的所有权已经转移给 channel 了
      });
      
      // 3. 消费者接收消息
      // recv 会阻塞，直到有消息到达
      let received = rx.recv().unwrap();
      println!("收到：{}", received);
  }

- mpsc::channel() 返回一对端点：发送端 tx 和接收端 rx
- send() 发送消息，返回 Result<(), SendError<T>>
- recv() 接收消息，会阻塞直到有消息，返回 Result<T, RecvError>
- 消息的所有权被转移，不是复制


多个消息
>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust
  :caption: 多个消息
  :name: multiple-messages

  use std::sync::mpsc;
  use std::thread;
  use std::time::Duration;

  fn main() {
      let (tx, rx) = mpsc::channel();
      
      thread::spawn(move || {
          let messages = vec![
              String::from("消息 1"),
              String::from("消息 2"),
              String::from("消息 3"),
          ];
          
          for msg in messages {
              tx.send(msg).unwrap();
              thread::sleep(Duration::from_millis(200));
          }
      });
      
      // 方式 1：阻塞接收
      // for received in rx {
      //     println!("收到：{}", received);
      // }
      
      // 方式 2：非阻塞接收
      loop {
          match rx.recv_timeout(Duration::from_millis(500)) {
              Ok(msg) => println!("收到：{}", msg),
              Err(_) => {
                  println!("超时，没有新消息");
                  break;
              }
          }
      }
  }

多生产者
>>>>>>>>>>>>>>>>>>

.. code-block:: rust
  :caption: 多生产者
  :name: multiple-producers

  use std::sync::mpsc;
  use std::thread;
  use std::time::Duration;

  fn main() {
      let (tx, rx) = mpsc::channel();
      
      // 生产者 1
      let tx1 = tx.clone();  // 克隆发送端
      thread::spawn(move || {
          for i in 1..=3 {
              tx1.send(format!("生产者 1: 消息 {}", i)).unwrap();
              thread::sleep(Duration::from_millis(100));
          }
      });
      
      // 生产者 2
      let tx2 = tx.clone();
      thread::spawn(move || {
          for i in 1..=3 {
              tx2.send(format!("生产者 2: 消息 {}", i)).unwrap();
              thread::sleep(Duration::from_millis(150));
          }
      });
      
      // 原始的 tx 在这里会 drop，不影响其他克隆
      
      // 消费者接收所有消息
      for received in rx {
          println!("收到：{}", received);
      }
  }

关键点：

- 用 tx.clone() 创建多个发送端
- 所有发送端共享同一个 channel
- 当所有发送端都 drop 后，recv() 会返回错误

不同类型的消息
>>>>>>>>>>>>>>>>>>>>>>

channel 可以发送任何类型的数据，包括枚举：

.. code-block:: rust
  :caption: 不同类型的消息
  :name: different-types-of-messages

  use std::sync::mpsc;
  use std::thread;

  // 定义消息类型
  enum Message {
      Text(String),
      Number(i32),
      Quit,
  }

  fn main() {
      let (tx, rx) = mpsc::channel();
      
      thread::spawn(move || {
          tx.send(Message::Text(String::from("你好"))).unwrap();
          tx.send(Message::Number(42)).unwrap();
          tx.send(Message::Quit).unwrap();
      });
      
      for msg in rx {
          match msg {
              Message::Text(text) => println!("文本：{}", text),
              Message::Number(num) => println!("数字：{}", num),
              Message::Quit => {
                  println!("收到退出信号");
                  break;
              }
          }
      }
  }

实战案例
>>>>>>>>>>>>>>>>

日志收集器

多个工作线程产生日志，一个专门的线程负责写入文件：

.. code-block:: rust
  :caption: 日志收集器
  :name: log-collector

  use std::sync::mpsc;
  use std::thread;
  use std::time::Duration;
  use std::fs::OpenOptions;
  use std::io::Write;

  enum LogMessage {
      Info(String),
      Warning(String),
      Error(String),
  }

  fn main() {
      let (tx, rx) = mpsc::channel();
      
      // 日志写入线程（消费者）
      let log_handle = thread::spawn(move || {
          let mut file = OpenOptions::new()
              .create(true)
              .append(true)
              .open("app.log")
              .unwrap();
          
          for msg in rx {
              let log_line = match msg {
                  LogMessage::Info(s) => format!("[INFO] {}\n", s),
                  LogMessage::Warning(s) => format!("[WARN] {}\n", s),
                  LogMessage::Error(s) => format!("[ERROR] {}\n", s),
              };
              
              writeln!(file, "{}", log_line).unwrap();
              print!("{}", log_line);  // 同时也输出到控制台
          }
      });
      
      // 模拟多个工作线程
      for i in 1..=3 {
          let tx_clone = tx.clone();
          thread::spawn(move || {
              for j in 1..=3 {
                  tx_clone.send(LogMessage::Info(
                      format!("工作线程 {} - 任务 {}", i, j)
                  )).unwrap();
                  thread::sleep(Duration::from_millis(100));
              }
          });
      }
      
      // 等待所有工作完成
      thread::sleep(Duration::from_secs(2));
      
      // drop 所有发送端，日志线程会退出
      drop(tx);
      log_handle.join().unwrap();
  }

Actor 模式

Actor 模式是一种并发设计模式：每个 Actor 是一个独立的实体，有自己的状态，通过消息与其他 Actor 通信。

.. code-block:: rust
  :caption: Actor 模式
  :name: actor-pattern

  use std::sync::mpsc;
  use std::thread;
  use std::collections::HashMap;

  // Actor 消息
  enum ActorMessage {
      Get(String),
      Set(String, String),
      Delete(String),
      Stop,
  }

  // Actor 响应
  enum ActorResponse {
      Value(Option<String>),
      Done,
  }

  // 简单的键值存储 Actor
  fn kv_store_actor(rx: mpsc::Receiver<ActorMessage>) {
      let mut store = HashMap::new();
      
      loop {
          match rx.recv() {
              Ok(ActorMessage::Get(key)) => {
                  let response = store.get(&key).cloned();
                  println!("Get {:?} -> {:?}", key, response);
              }
              Ok(ActorMessage::Set(key, value)) => {
                  store.insert(key, value);
                  println!("Set 完成");
              }
              Ok(ActorMessage::Delete(key)) => {
                  store.remove(&key);
                  println!("Delete {:?} 完成", key);
              }
              Ok(ActorMessage::Stop) => {
                  println!("Actor 停止");
                  break;
              }
              Err(_) => break,
          }
      }
  }

  fn main() {
      let (tx, rx) = mpsc::channel();
      
      // 启动 Actor
      let actor_handle = thread::spawn(move || {
          kv_store_actor(rx);
      });
      
      // 与 Actor 交互
      tx.send(ActorMessage::Set("name".to_string(), "Larry".to_string())).unwrap();
      tx.send(ActorMessage::Set("age".to_string(), "25".to_string())).unwrap();
      tx.send(ActorMessage::Get("name".to_string())).unwrap();
      tx.send(ActorMessage::Delete("age".to_string())).unwrap();
      tx.send(ActorMessage::Stop).unwrap();
      
      actor_handle.join().unwrap();
  }

工作池模式
:::::::::::::::::

.. code-block:: rust
  :caption: 工作池模式
  :name: worker-pool-pattern

  use std::sync::mpsc;
  use std::thread;
  use std::time::Duration;

  enum Job {
      Process(i32),
      Shutdown,
  }

  fn worker(id: usize, rx: std::sync::Arc<std::sync::Mutex<mpsc::Receiver<Job>>>) {
      loop {
          // 需要 Mutex 因为多个线程共享同一个 receiver
          let job = {
              let rx = rx.lock().unwrap();
              rx.recv().unwrap()
          };
          
          match job {
              Job::Process(n) => {
                  println!("工人 {} 处理任务 {}", id, n);
                  thread::sleep(Duration::from_millis(100));
                  println!("工人 {} 完成任务 {}", id, n);
              }
              Job::Shutdown => {
                  println!("工人 {} 收到 shutdown", id);
                  break;
              }
          }
      }
  }

  fn main() {
      let (tx, rx) = mpsc::channel();
      let rx = std::sync::Arc::new(std::sync::Mutex::new(rx));
      
      // 创建 3 个工作线程
      let mut workers = vec![];
      for i in 0..3 {
          let rx_clone = std::sync::Arc::clone(&rx);
          let handle = thread::spawn(move || {
              worker(i, rx_clone);
          });
          workers.push(handle);
      }
      
      // 发送任务
      for i in 0..6 {
          tx.send(Job::Process(i)).unwrap();
      }
      
      // 发送 shutdown 信号
      for _ in 0..3 {
          tx.send(Job::Shutdown).unwrap();
      }
      
      // 等待所有工人完成
      for handle in workers {
          handle.join().unwrap();
      }
  }

Actor 模式
====================