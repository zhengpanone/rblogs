============================================================
RAII（Resource Acquisition Is Initialization） in Rust
============================================================

RAII（资源获取即初始化）是 Rust 内存安全和资源管理的核心思想之一。

在 Rust 中，RAII 通过：

- 所有权（Ownership）
- 生命周期（Lifetime）
- ``Drop`` Trait

实现自动资源管理，不需要 GC，也不需要手动 ``free()``。

.. contents:: 目录
   :depth: 3
   :local:

1. 什么是 RAII
===============

C 语言：

.. code-block:: c

   FILE* fp = fopen("test.txt", "r");

   /* 使用文件 */

   fclose(fp);

问题：

.. code-block:: c

   if (error) {
       return;
   }

可能忘记：

.. code-block:: c

   fclose(fp);

导致资源泄露。

Rust：

.. code-block:: rust

   {
       let file = File::open("test.txt")?;

       // 使用文件
   }

离开作用域：

.. code-block:: text

   file
   ↓
   Drop
   ↓
   close()

自动释放。

这就是 RAII。

2. Rust 中的 RAII
==================

示例：

.. code-block:: rust

   struct Resource;

   impl Drop for Resource {
       fn drop(&mut self) {
           println!("释放资源");
       }
   }

使用：

.. code-block:: rust

   fn main() {
       let r = Resource;

       println!("业务逻辑");
   }

输出：

.. code-block:: text

   业务逻辑
   释放资源

实际执行顺序：

.. code-block:: text

   创建 Resource
   ↓
   使用 Resource
   ↓
   离开作用域
   ↓
   调用 Drop

3. Drop Trait
==============

Rust 的析构函数：

.. code-block:: rust

   pub trait Drop {
       fn drop(&mut self);
   }

例如：

.. code-block:: rust

   struct DatabaseConnection;

实现：

.. code-block:: rust

   impl Drop for DatabaseConnection {
       fn drop(&mut self) {
           println!("关闭数据库连接");
       }
   }

使用：

.. code-block:: rust

   {
       let conn = DatabaseConnection;
   }

自动：

.. code-block:: text

   关闭数据库连接

4. 多个对象的 Drop 顺序
========================

.. code-block:: rust

   struct A;
   struct B;

.. code-block:: rust

   impl Drop for A {
       fn drop(&mut self) {
           println!("drop A");
       }
   }

   impl Drop for B {
       fn drop(&mut self) {
           println!("drop B");
       }
   }

.. code-block:: rust

   fn main() {
       let a = A;
       let b = B;
   }

输出：

.. code-block:: text

   drop B
   drop A

规则：

**后创建先销毁（LIFO）**

类似栈：

.. code-block:: text

   push A
   push B

   pop B
   pop A

5. 提前释放资源
================

默认：

.. code-block:: rust

   let file = File::open("a.txt")?;

   // ...

   // 作用域结束才释放

提前释放：

.. code-block:: rust

   std::mem::drop(file);

例如：

.. code-block:: rust

   let mutex_guard = mutex.lock().unwrap();

   do_work();

   drop(mutex_guard);

   do_other_work();

锁提前释放。

6. Mutex 为什么依赖 RAII
=========================

看标准库：

.. code-block:: rust

   let guard = mutex.lock().unwrap();

类型：

.. code-block:: rust

   MutexGuard<T>

内部：

.. code-block:: rust

   impl Drop for MutexGuard<'_, T> {
       fn drop(&mut self) {
           unlock_mutex();
       }
   }

所以：

.. code-block:: rust

   {
       let guard = mutex.lock().unwrap();

       // 自动加锁
   }

离开作用域：

.. code-block:: text

   guard Drop
   ↓
   unlock

自动解锁。

7. 文件为什么依赖 RAII
=======================

.. code-block:: rust

   use std::fs::File;

   let file = File::open("a.txt")?;

底层：

.. code-block:: text

   File
   ↓
   fd
   ↓
   Drop
   ↓
   close(fd)

因此：

.. code-block:: rust

   {
       let file = File::open("a.txt")?;
   }

自动关闭文件。

8. Box 为什么依赖 RAII
=======================

.. code-block:: rust

   let b = Box::new(100);

离开作用域：

.. code-block:: text

   Box Drop
   ↓
   释放堆内存

等价于：

.. code-block:: c

   free(ptr);

但自动执行。

9. Rc 为什么依赖 RAII
======================

.. code-block:: rust

   use std::rc::Rc;

   let a = Rc::new(String::from("hello"));

   let b = Rc::clone(&a);

引用计数：

.. code-block:: text

   count = 2

离开作用域：

.. code-block:: text

   drop(b)
   count = 1

   drop(a)
   count = 0
   释放内存

``Rc`` 的核心也是 ``Drop``。

10. RAII 管理事务
==================

数据库事务是经典案例。

.. code-block:: rust

   struct Transaction {
       committed: bool,
   }

.. code-block:: rust

   impl Drop for Transaction {
       fn drop(&mut self) {
           if !self.committed {
               println!("rollback");
           }
       }
   }

提交：

.. code-block:: rust

   impl Transaction {
       fn commit(mut self) {
           self.committed = true;
           println!("commit");
       }
   }

使用：

.. code-block:: rust

   {
       let tx = Transaction {
           committed: false,
       };

       // 出错直接 return
   }

自动：

.. code-block:: text

   rollback

这就是很多数据库库的实现思路。

11. Scope Guard（高级 RAII）
=============================

例如：

.. code-block:: rust

   let temp_file = create_temp_file();

希望退出时删除：

.. code-block:: rust

   struct TempFile {
       path: PathBuf,
   }

.. code-block:: rust

   impl Drop for TempFile {
       fn drop(&mut self) {
           std::fs::remove_file(&self.path).ok();
       }
   }

使用：

.. code-block:: rust

   {
       let temp = TempFile::new();

       // 使用临时文件
   }

自动：

.. code-block:: text

   删除临时文件

类似 Go：

.. code-block:: go

   defer os.Remove(...)

但 Rust 是 RAII。

12. 企业级案例：连接池
=======================

.. code-block:: rust

   let conn = pool.get()?;

类型：

.. code-block:: rust

   PooledConnection

业务完成：

.. code-block:: rust

   {
       let conn = pool.get()?;
   }

自动：

.. code-block:: text

   Drop
   ↓
   归还连接池

而不是：

.. code-block:: java

   finally {
       conn.close();
   }

例如：

- ``sqlx``
- ``deadpool``
- ``bb8``
- ``r2d2``

都大量使用 RAII。

13. RAII + Builder + Deref
===========================

Rust 常见设计组合：

.. code-block:: text

   Builder
       ↓
   创建对象

   Deref
       ↓
   像引用一样使用

   Drop(RAII)
       ↓
   自动释放资源

例如：

.. code-block:: rust

   let client = Client::builder()
       .timeout(...)
       .build()?;

内部：

.. code-block:: text

   Client
       ├─ Builder
       ├─ Deref
       └─ Drop

14. Rust 核心智能指针与 RAII
=============================

.. list-table:: Rust 核心智能指针与 RAII 行为
   :header-rows: 1
   :widths: 30 50

   * - 类型
     - RAII 行为
   * - ``Box<T>``
     - 释放堆内存
   * - ``Vec<T>``
     - 释放动态数组
   * - ``String``
     - 释放字符串缓冲区
   * - ``File``
     - 关闭文件
   * - ``TcpStream``
     - 关闭 Socket
   * - ``MutexGuard<T>``
     - 自动解锁
   * - ``RwLockGuard<T>``
     - 自动解锁
   * - ``Rc<T>``
     - 减少引用计数
   * - ``Arc<T>``
     - 原子减少引用计数
   * - ``Transaction``
     - 自动回滚
   * - ``PooledConnection``
     - 自动归还连接

必学的下一步
=============

理解 RAII 后，建议按顺序继续：

.. code-block:: text

   Ownership（所有权）
       ↓
   Borrowing（借用）
       ↓
   Drop（析构）
       ↓
   RAII
       ↓
   Box
       ↓
   Rc
       ↓
   RefCell
       ↓
   Arc
       ↓
   Mutex
       ↓
   Pin
       ↓
   Async Runtime

其中：

- ``RAII + Drop`` 是资源管理基础
- ``Rc + RefCell`` 是单线程共享所有权
- ``Arc + Mutex`` 是多线程共享所有权

这条线基本贯穿整个 Rust 智能指针生态。