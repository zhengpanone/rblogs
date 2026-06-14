================================
智能指针（Smart Pointers）
================================

智能指针是 Rust 中管理内存和资源的利器。它像指针一样工作，但带有额外的元数据和能力，例如引用计数、自动释放、内部可变性等。

在 Rust 中，智能指针通常通过实现 ``Deref`` 和 ``Drop`` 两个 Trait 来工作：

- ``Deref``：让智能指针可以像普通引用一样被解引用
- ``Drop``：在离开作用域时自动清理资源（RAII）

.. contents:: 目录
   :depth: 3
   :local:

什么是智能指针
==================

普通指针 vs 智能指针：

.. code-block:: rust

   // 普通引用：只指向数据，不管理生命周期
   let x = 42;
   let r: &i32 = &x;

   // 智能指针：拥有数据，管理释放，可解引用
   let b = Box::new(42);   // 在堆上分配，离开作用域自动释放

核心特征：

.. list-table:: 智能指针 vs 普通引用
   :header-rows: 1
   :widths: 30 35 35

   * - 特征
     - 普通引用 ``&T``
     - 智能指针
   * - 所有权
     - 不拥有数据
     - 拥有数据
   * - 内存位置
     - 指向栈或堆
     - 通常在堆上
   * - 生命周期管理
     - 编译器借用检查
     - RAII（Drop）
   * - 解引用
     - 自动
     - 通过 ``Deref`` Trait
   * - 典型例子
     - ``&x``
     - ``Box<T>``, ``Rc<T>``, ``Arc<T>``

Box\<T\>：堆分配
=====================

``Box<T>`` 是最简单的智能指针，将数据分配到堆上。

基本使用：

.. code-block:: rust

   fn main() {
       // 在堆上分配一个 i32
       let b = Box::new(5);
       println!("b = {}", b); // 自动解引用

       // 离开作用域，堆内存自动释放
   }

适用场景：

.. list-table:: Box\<T\> 的使用场景
   :header-rows: 1
   :widths: 30 70

   * - 场景
     - 说明
   * - 类型大小不确定
     - 递归类型（如链表）编译期无法确定大小
   * - 转移大块数据所有权
     - 避免在栈上复制大量数据
   * - Trait Object
     - ``Box<dyn Trait>`` 实现动态派发

递归类型示例（链表）：

.. code-block:: rust

   // 使用 Box 解决递归类型大小问题
   enum List {
       Cons(i32, Box<List>),
       Nil,
   }

   use List::{Cons, Nil};

   fn main() {
       let list = Cons(1, Box::new(Cons(2, Box::new(Cons(3, Box::new(Nil))))));
   }

为什么必须用 Box？

.. code-block:: text

   enum List {
       Cons(i32, List),  // ❌ 编译错误：递归类型大小无限
   }

   enum List {
       Cons(i32, Box<List>),  // ✅ Box 大小固定（一个指针大小）
   }

Trait Object：

.. code-block:: rust

   trait Draw {
       fn draw(&self);
   }

   struct Circle;
   impl Draw for Circle {
       fn draw(&self) { println!("画圆"); }
   }

   struct Square;
   impl Draw for Square {
       fn draw(&self) { println!("画方"); }
   }

   fn main() {
       let shapes: Vec<Box<dyn Draw>> = vec![
           Box::new(Circle),
           Box::new(Square),
       ];

       for shape in &shapes {
           shape.draw();
       }
   }

Deref Trait：像引用一样使用
================================

智能指针的核心能力来自 ``Deref``：

.. code-block:: rust

   use std::ops::Deref;

   struct MyBox<T>(T);

   impl<T> MyBox<T> {
       fn new(x: T) -> MyBox<T> {
           MyBox(x)
       }
   }

   impl<T> Deref for MyBox<T> {
       type Target = T;

       fn deref(&self) -> &Self::Target {
           &self.0
       }
   }

   fn main() {
       let x = MyBox::new(5);
       assert_eq!(5, *x); // *x 等价于 *(x.deref())
   }

解引用强制转换（Deref Coercion）：

.. code-block:: rust

   fn hello(name: &str) {
       println!("Hello, {}!", name);
   }

   fn main() {
       let m = MyBox::new(String::from("Rust"));
       hello(&m); // &MyBox<String> → &String → &str
   }

编译器自动执行多步解引用，直到类型匹配。

1. Drop Trait：自动清理
=========================

智能指针离开作用域时自动执行清理：

.. code-block:: rust

   struct CustomSmartPointer {
       data: String,
   }

   impl Drop for CustomSmartPointer {
       fn drop(&mut self) {
           println!("释放 CustomSmartPointer: `{}`", self.data);
       }
   }

   fn main() {
       let c = CustomSmartPointer {
           data: String::from("我的数据"),
       };
       let d = CustomSmartPointer {
           data: String::from("其他数据"),
       };
       println!("CustomSmartPointer 已创建");
   }

输出：

.. code-block:: text

   CustomSmartPointer 已创建
   释放 CustomSmartPointer: `其他数据`    ← 后创建先销毁（LIFO）
   释放 CustomSmartPointer: `我的数据`

提前释放：

.. code-block:: rust

   fn main() {
       let c = CustomSmartPointer {
           data: String::from("数据"),
       };
       drop(c); // 提前释放，不等作用域结束
       println!("c 已被提前释放");
   }

注意：必须用 ``std::mem::drop()``，不能直接调用 ``c.drop()``。

5. Rc\<T\>：引用计数（单线程）
================================

``Rc<T>`` 允许多个所有者共享同一份数据，引用计数归零时释放。

基本使用：

.. code-block:: rust

   use std::rc::Rc;

   fn main() {
       let a = Rc::new(String::from("hello"));
       println!("count after creating a = {}", Rc::strong_count(&a)); // 1

       let b = Rc::clone(&a);
       println!("count after creating b = {}", Rc::strong_count(&a)); // 2

       {
           let c = Rc::clone(&a);
           println!("count after creating c = {}", Rc::strong_count(&a)); // 3
       }
       // c 离开作用域，count 减为 2

       println!("count after c goes out = {}", Rc::strong_count(&a)); // 2
   }

内部机制：

.. code-block:: text

   Rc::new(data)
   │
   ├── 堆上分配:
   │   ├── data
   │   └── strong_count = 1
   │
   ├── Rc::clone() → strong_count + 1
   └── Drop → strong_count - 1 → 0 → 释放

``Rc<T>`` 的限制：只读共享，不能修改内部数据（需要配合 ``RefCell<T>``）。

6. RefCell\<T\>：内部可变性
=============================

``RefCell<T>`` 将借用检查从编译期推迟到运行时，允许在拥有不可变引用时修改内部数据。

基本使用：

.. code-block:: rust

   use std::cell::RefCell;

   fn main() {
       let data = RefCell::new(42);

       // 运行时借用检查
       {
           let mut borrowed = data.borrow_mut();
           *borrowed += 1;
       } // borrow_mut 在此释放

       println!("data = {}", data.borrow()); // 43
   }

``borrow()`` 和 ``borrow_mut()`` 的规则：

.. list-table:: RefCell 借用规则
   :header-rows: 1
   :widths: 30 70

   * - 方法
     - 规则
   * - ``borrow()``
     - 返回 ``Ref<T>``，多个不可变借用可同时存在
   * - ``borrow_mut()``
     - 返回 ``RefMut<T>``，同一时刻只能有一个可变借用
   * - 冲突
     - 运行时 panic（与编译期借用检查不同）

违反规则会 panic：

.. code-block:: rust

   use std::cell::RefCell;

   fn main() {
       let data = RefCell::new(42);

       let a = data.borrow();
       let b = data.borrow_mut(); // ❌ 运行时 panic：已有不可变借用
   }

7. Rc\<RefCell\<T\>\>：共享可变数据
======================================

``Rc`` 提供多所有权，``RefCell`` 提供内部可变性，两者组合使用：

.. code-block:: rust

   use std::rc::Rc;
   use std::cell::RefCell;

   #[derive(Debug)]
   struct Node {
       value: i32,
       children: Vec<Rc<RefCell<Node>>>,
   }

   fn main() {
       let leaf = Rc::new(RefCell::new(Node {
           value: 3,
           children: vec![],
       }));

       let branch = Rc::new(RefCell::new(Node {
           value: 5,
           children: vec![Rc::clone(&leaf)],
       }));

       // 修改 leaf（即使被多个所有者共享）
       leaf.borrow_mut().value = 10;

       println!("leaf: {:?}", leaf.borrow());
       println!("branch children: {:?}", branch.borrow().children[0].borrow());
   }

典型应用：树结构、图结构等需要共享可变数据的场景。

8. Arc\<T\>：原子引用计数（多线程）
=====================================

``Arc<T>`` 与 ``Rc<T>`` 类似，但使用原子操作保证线程安全：

.. code-block:: rust

   use std::sync::Arc;
   use std::thread;

   fn main() {
       let data = Arc::new(vec![1, 2, 3]);

       let mut handles = vec![];

       for i in 0..3 {
           let data = Arc::clone(&data);
           let handle = thread::spawn(move || {
               println!("线程 {}: {:?}", i, data);
           });
           handles.push(handle);
       }

       for handle in handles {
           handle.join().unwrap();
       }
   }

``Rc<T>`` vs ``Arc<T>``：

.. list-table:: Rc\<T\> vs Arc\<T\>
   :header-rows: 1
   :widths: 25 35 40

   * - 特性
     - ``Rc<T>``
     - ``Arc<T>``
   * - 线程安全
     - 否（非原子计数）
     - 是（原子计数）
   * - 性能
     - 更快（无原子开销）
     - 较慢（原子操作开销）
   * - 使用场景
     - 单线程共享所有权
     - 多线程共享所有权
   * - 实现 Send/Sync
     - 否
     - 是（T 满足时）

9. Arc\<Mutex\<T\>\>：多线程共享可变数据
===========================================

``Mutex`` 提供互斥访问，与 ``Arc`` 组合实现多线程安全共享：

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
   // MutexGuard 在此 drop，自动解锁

10. Arc\<RwLock\<T\>\>：读写锁
=================================

``RwLock`` 允许多个读者同时访问，但写者独占：

.. code-block:: rust

   use std::sync::{Arc, RwLock};
   use std::thread;

   fn main() {
       let data = Arc::new(RwLock::new(0));

       // 多个读者
       let readers: Vec<_> = (0..5).map(|i| {
           let data = Arc::clone(&data);
           thread::spawn(move || {
               let value = data.read().unwrap();
               println!("读者 {}: {}", i, *value);
           })
       }).collect();

       // 一个写者
       {
           let data = Arc::clone(&data);
           thread::spawn(move || {
               let mut value = data.write().unwrap();
               *value = 42;
           }).join().unwrap();
       }

       for reader in readers {
           reader.join().unwrap();
       }
   }

``Mutex<T>`` vs ``RwLock<T>``：

.. list-table:: Mutex vs RwLock
   :header-rows: 1
   :widths: 25 40 35

   * - 特性
     - ``Mutex<T>``
     - ``RwLock<T>``
   * - 并发读
     - 不允许
     - 允许（多个读锁）
   * - 写操作
     - 独占锁
     - 独占锁
   * - 适用场景
     - 读写频繁相当
     - 读多写少
   * - 开销
     - 较低
     - 略高

11. Weak\<T\>：弱引用，解决循环引用
======================================

``Rc<T>`` 和 ``Arc<T>`` 都支持 ``Weak<T>``，不增加强引用计数，避免内存泄漏。

循环引用问题：

.. code-block:: rust

   use std::rc::Rc;
   use std::cell::RefCell;

   #[derive(Debug)]
   struct Node {
       parent: RefCell<Option<Rc<Node>>>,
   }

   // ❌ 循环引用导致内存泄漏
   // node_a → node_b → node_a ...

使用 ``Weak<T>`` 解决：

.. code-block:: rust

   use std::rc::{Rc, Weak};
   use std::cell::RefCell;

   #[derive(Debug)]
   struct Node {
       value: i32,
       parent: RefCell<Weak<Node>>,       // 弱引用，不增加 count
       children: RefCell<Vec<Rc<Node>>>,  // 强引用
   }

   fn main() {
       let leaf = Rc::new(Node {
           value: 3,
           parent: RefCell::new(Weak::new()),
           children: RefCell::new(vec![]),
       });

       println!("leaf strong = {}, weak = {}", Rc::strong_count(&leaf), Rc::weak_count(&leaf));

       {
           let branch = Rc::new(Node {
               value: 5,
               parent: RefCell::new(Weak::new()),
               children: RefCell::new(vec![Rc::clone(&leaf)]),
           });

           // 设置 leaf 的父节点为 branch（弱引用）
           *leaf.parent.borrow_mut() = Rc::downgrade(&branch);

           println!("branch strong = {}, weak = {}", Rc::strong_count(&branch), Rc::weak_count(&branch));
       }

       // branch 离开作用域后被释放，leaf.parent 自动变为无效
       println!("leaf parent = {:?}", leaf.parent.borrow().upgrade()); // None
   }

``Weak<T>`` 核心 API：

.. list-table:: Weak\<T\> 核心 API
   :header-rows: 1
   :widths: 35 65

   * - 方法
     - 说明
   * - ``Rc::downgrade(&rc)``
     - 从强引用创建弱引用
   * - ``weak.upgrade()``
     - 尝试获取强引用，返回 ``Option<Rc<T>>``
   * - ``Rc::weak_count(&rc)``
     - 查看弱引用计数

12. Cow\<T\>：写时复制
=========================

``Cow<T>``（Clone on Write）在需要修改时才复制数据：

.. code-block:: rust

   use std::borrow::Cow;

   fn process(input: &str) -> Cow<str> {
       if input.contains(' ') {
           Cow::Owned(input.replace(' ', "_"))
       } else {
           Cow::Borrowed(input) // 不需要修改，不复制
       }
   }

   fn main() {
       let s1 = "hello";
       let s2 = "hello world";

       let r1 = process(s1);
       let r2 = process(s2);

       println!("r1: {}", r1); // "hello"  (Borrowed，未复制)
       println!("r2: {}", r2); // "hello_world" (Owned，复制并修改)
   }

适用场景：可能不需要修改的大数据，避免不必要的克隆开销。

13. Cell\<T\>：值的内部可变性
================================

``Cell<T>`` 适用于 ``Copy`` 类型，通过 ``get()`` 和 ``set()`` 读写：

.. code-block:: rust

   use std::cell::Cell;

   fn main() {
       let x = Cell::new(42);

       x.set(100);
       println!("x = {}", x.get()); // 100

       // 不需要 borrow_mut，直接替换整个值
       let old = x.replace(200);
       println!("旧值: {}, 新值: {}", old, x.get());
   }

``Cell<T>`` vs ``RefCell<T>``：

.. list-table:: Cell\<T\> vs RefCell\<T\>
   :header-rows: 1
   :widths: 25 40 35

   * - 特性
     - ``Cell<T>``
     - ``RefCell<T>``
   * - 获取值
     - 复制整个值（``get/set``）
     - 获取引用（``borrow/borrow_mut``）
   * - 类型要求
     - ``T: Copy``
     - 无限制
   * - 运行时检查
     - 无（总是安全）
     - 有（运行时借用检查）
   * - 性能
     - 更高
     - 有检查开销

14. Pin\<T\>：固定内存位置
=============================

``Pin<T>`` 保证数据不会被移动，用于自引用结构和异步编程：

.. code-block:: rust

   use std::pin::Pin;
   use std::marker::PhantomPinned;

   #[derive(Debug)]
   struct SelfReferential {
       value: String,
       pointer: *const String,
       _pin: PhantomPinned,
   }

   impl SelfReferential {
       fn new(value: String) -> Pin<Box<Self>> {
           let mut boxed = Box::pin(SelfReferential {
               value,
               pointer: std::ptr::null(),
               _pin: PhantomPinned,
           });

           let pointer = &boxed.value as *const String;

           // 安全：数据已被 Pin 固定
           unsafe {
               let mut_ref = Pin::as_mut(&mut boxed);
               Pin::get_unchecked_mut(mut_ref).pointer = pointer;
           }

           boxed
       }
   }

主要应用：``Future``（异步编程中，future 常包含自引用）。

15. 智能指针对比总览
======================

.. list-table:: Rust 智能指针总览
   :header-rows: 1
   :widths: 25 25 25 25

   * - 类型
     - 所有权模型
     - 线程安全
     - 核心用途
   * - ``Box<T>``
     - 唯一所有权
     - 是
     - 堆分配、递归类型、Trait Object
   * - ``Rc<T>``
     - 共享所有权（强引用计数）
     - 否
     - 单线程共享数据
   * - ``Arc<T>``
     - 共享所有权（原子引用计数）
     - 是
     - 多线程共享数据
   * - ``RefCell<T>``
     - 内部可变性
     - 否
     - 运行时借用检查
   * - ``Cell<T>``
     - 内部可变性（值替换）
     - 否
     - Copy 类型的内部可变性
   * - ``Mutex<T>``
     - 互斥访问
     - 是
     - 多线程互斥
   * - ``RwLock<T>``
     - 读写锁
     - 是
     - 多线程读多写少
   * - ``Weak<T>``
     - 弱引用
     - 跟随 Rc/Arc
     - 解决循环引用
   * - ``Cow<T>``
     - 写时复制
     - 否
     - 延迟克隆
   * - ``Pin<T>``
     - 固定内存
     - 跟随内部类型
     - 自引用结构、Future

16. 常见组合模式
==================

.. list-table:: 智能指针组合模式
   :header-rows: 1
   :widths: 35 30 35

   * - 组合
     - 场景
     - 能力
   * - ``Rc<RefCell<T>>``
     - 单线程共享可变数据
     - 多个所有者 + 运行时可变
   * - ``Arc<Mutex<T>>``
     - 多线程互斥共享
     - 线程安全 + 互斥访问
   * - ``Arc<RwLock<T>>``
     - 多线程读写共享
     - 线程安全 + 读写锁
   * - ``Box<dyn Trait>``
     - 动态派发
     - 堆上的 Trait Object
   * - ``Pin<Box<T>>``
     - 固定堆分配
     - 自引用 / Future
   * - ``Cow<'a, T>``
     - 条件克隆
     - 借用或拥有

17. 实现一个自定义智能指针
============================

综合练习——实现一个带引用计数的简化版智能指针：

.. code-block:: rust

   use std::ops::Deref;
   use std::alloc::{alloc, dealloc, Layout};

   struct MyRc<T> {
       ptr: *mut RcInner<T>,
   }

   struct RcInner<T> {
       count: usize,
       value: T,
   }

   impl<T> MyRc<T> {
       fn new(value: T) -> Self {
           let layout = Layout::new::<RcInner<T>>();
           unsafe {
               let ptr = alloc(layout) as *mut RcInner<T>;
               (*ptr) = RcInner { count: 1, value };
               MyRc { ptr }
           }
       }
   }

   impl<T> Clone for MyRc<T> {
       fn clone(&self) -> Self {
           unsafe {
               (*self.ptr).count += 1;
           }
           MyRc { ptr: self.ptr }
       }
   }

   impl<T> Drop for MyRc<T> {
       fn drop(&mut self) {
           unsafe {
               (*self.ptr).count -= 1;
               if (*self.ptr).count == 0 {
                   let layout = Layout::new::<RcInner<T>>();
                   dealloc(self.ptr as *mut u8, layout);
               }
           }
       }
   }

   impl<T> Deref for MyRc<T> {
       type Target = T;

       fn deref(&self) -> &T {
           unsafe { &(*self.ptr).value }
       }
   }

   fn main() {
       let a = MyRc::new(42);
       let b = a.clone();

       println!("*a = {}", *a);
       println!("*b = {}", *b);
   }

总结
=====

.. code-block:: text

   智能指针体系
   │
   ├── 基础层
   │   ├── Box<T>       堆分配
   │   ├── Deref        解引用
   │   └── Drop         自动释放
   │
   ├── 共享所有权
   │   ├── Rc<T>        单线程引用计数
   │   ├── Arc<T>       多线程原子引用计数
   │   └── Weak<T>      弱引用（防循环）
   │
   ├── 内部可变性
   │   ├── Cell<T>      值替换（Copy 类型）
   │   └── RefCell<T>   运行时借用检查
   │
   ├── 线程同步
   │   ├── Mutex<T>     互斥锁
   │   └── RwLock<T>    读写锁
   │
   ├── 高级
   │   ├── Cow<T>       写时复制
   │   └── Pin<T>       固定内存
   │
   └── 组合模式
       ├── Rc<RefCell<T>>
       ├── Arc<Mutex<T>>
       ├── Arc<RwLock<T>>
       └── Box<dyn Trait>
