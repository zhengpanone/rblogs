================
Trait
================

.. contents:: 目录
   :depth: 3
   :local:


Trait
===================

Trait 是 Rust 的 **多态** 机制。

Trait 它不仅用于接口抽象，还广泛参与：

- 泛型约束
- 运算符重载
- 类型转换
- 迭代器系统
- 错误处理
- 格式化输出
- 并发安全
- 资源释放
- 异步转换

基本语法
-------------------

.. code-block:: rust

    trait Greeter {
        // 只有声明，没有默认实现（必须由实现者提供）
        fn name(&self) -> &str;
        
        // 有默认实现的方法
        fn greet(&self) {
            println!("Hello, {}!", self.name());
        }
    }

    struct Person {
        name: String,
    }

    impl Greeter for Person {
        // 必须实现没有默认实现的方法
        fn name(&self) -> &str {
            &self.name
        }
        // greet 方法使用默认实现
    }

    struct Robot {
        id: u32,
    }

    impl Greeter for Robot {
        fn name(&self) -> &str {
            "Robot-42"
        }
        
        // 覆盖默认实现
        fn greet(&self) {
            println!("Beep boop, I am {}", self.name());
        }
    }

    fn main() {
        let person = Person { name: String::from("Alice") };
        let robot = Robot { id: 42 };
        
        person.greet(); // Hello, Alice!
        robot.greet();  // Beep boop, I am Robot-42
    }

Trait 作为参数
-------------------

.. code-block:: rust

    // 定义一个 trait
    trait Summary {
        fn summarize(&self) -> String;
    }

    // 为不同类型实现
    struct Article {
        headline: String,
        content: String,
    }

    impl Summary for Article {
        fn summarize(&self) -> String {
            format!("{}...", &self.headline)
        }
    }

    struct Tweet {
        username: String,
        content: String,
    }

    impl Summary for Tweet {
        fn summarize(&self) -> String {
            format!("@{}: {}", self.username, self.content)
        }
    }

使用 impl Trait 语法
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

    fn notify(item: impl Summary) {
        println!("Breaking news! {}", item.summarize());
    }

    notify(article);
    notify(tweet);

使用 Trait Bound 语法
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

    // 等价写法
    fn notify<T: Summary>(item: T) {
        println!("Breaking news! {}", item.summarize());
    }

多个 Trait 约束
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

    use std::fmt::Display;

    // impl Trait 语法
    fn notify(item: impl Summary + Display) {
        println!("Breaking news! {}", item.summarize());
    }

    // Trait Bound 语法
    fn notify<T: Summary + Display>(item: T) {
        println!("Breaking news! {}", item.summarize());
    }

    // where 子句（更清晰）
    fn notify<T>(item: T)
    where
        T: Summary + Display,
    {
        println!("Breaking news! {}", item.summarize());
    }

Trait 作为返回值
--------------------------

.. code-block:: rust

    fn returns_summarizable() -> impl Summary {
        Tweet {
            username: String::from("horse_ebooks"),
            content: String::from("of course, as you probably already know, people"),
        }
    }


限制：只能返回一种具体类型：

.. code-block:: rust

    // 编译错误：if 和 else 返回不同类型
    fn returns_summarizable(switch: bool) -> impl Summary {
        if switch {
            Article { /* ... */ }
        } else {
            Tweet { /* ... */ }
        }
    }

    // 使用 Box<dyn Trait>
    fn returns_summarizable(switch: bool) -> Box<dyn Summary> {
        if switch {
            Box::new(Article { /* ... */ })
        } else {
            Box::new(Tweet { /* ... */ })
        }
    }

Trait 对象（动态分发）
--------------------------

基本用法
>>>>>>>>>>>>>>>>>>

.. code-block:: rust

    // 使用 dyn Trait 创建 trait 对象
    let articles: Vec<Box<dyn Summary>> = vec![
        Box::new(Article { /* ... */ }),
        Box::new(Tweet { /* ... */ }),
    ];

    for article in articles {
        println!("{}", article.summarize());
    }

静态分发 vs 动态分发
>>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

    // 静态分发（编译时确定，零开销）
    fn static_dispatch<T: Summary>(item: T) {
        item.summarize();
    }

    // 动态分发（运行时确定，有虚表开销）
    fn dynamic_dispatch(item: &dyn Summary) {
        item.summarize();
    }

完全限定语法
---------------------

当多个 trait 有同名方法时：

.. code-block:: rust

    trait Pilot {
        fn fly(&self);
    }

    trait Wizard {
        fn fly(&self);
    }

    struct Human;

    impl Pilot for Human {
        fn fly(&self) {
            println!("This is your captain speaking.");
        }
    }

    impl Wizard for Human {
        fn fly(&self) {
            println!("Up!");
        }
    }

    impl Human {
        fn fly(&self) {
            println!("*waving arms furiously*");
        }
    }

    fn main() {
        let person = Human;
        person.fly();           // *waving arms furiously*
        Pilot::fly(&person);    // This is your captain speaking.
        Wizard::fly(&person);   // Up!
    }

Supertrait（Trait 继承）
-----------------------------

.. code-block:: rust

    trait Shape {
        fn area(&self) -> f64;
    }

    trait ShapeDescription: Shape {
        fn description(&self) -> String {
            format!("A shape with area {:.2}", self.area())
        }
    }

    struct Circle {
        radius: f64,
    }

    impl Shape for Circle {
        fn area(&self) -> f64 {
            std::f64::consts::PI * self.radius * self.radius
        }
    }

    impl ShapeDescription for Circle {} // 因为 Circle 已经实现了 Shape





常见的Trait和用法
==================

转换
-----------------

From<T>
>>>>>>>>>>>>>>>>>>>>>>>

把一种类型 无损 转换成另一种类型。
例子： ``let s: String = String::from("hello");``

.. code-block:: rust
  :caption: From 的例子

  struct UserId(u64);

  impl From<u64> for UserId {
      fn from(value: u64) -> Self {
          UserId(value)
      }
  }

  fn main(){
      let id = UserId::from(1234567890);
      println!("{}", id);

      let id:UserId = 100.into();
  }



Into<T>
>>>>>>>>>>>>>>>>>>>>>>>

From 的反向实现。只要实现了 From<A> for B，就自动实现 Into<B> for A。
常用在泛型函数里： ``fn foo<T: Into<String>>(s: T)``。

TryFrom<T> / TryInto<T>
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

带错误处理的转换（可能失败）。
例子： ``let n: u8 = u8::try_from(300).unwrap_err()``;

FromStr
>>>>>>>>>>>>>>>>>>>>>>>

从字符串解析出某个类型。
例子： ``let ip: IpAddr = "127.0.0.1".parse().unwrap()``;

运算符重载
-------------------

Add, Sub, Mul, Div, Rem
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

定义 + - * / % 运算。
例子：impl Add for Point { … } 可以自定义点的加法。

Neg
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

定义一元负号 -x。

Index, IndexMut
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

实现 [] 下标操作。
例子：vec[0] 背后就是调用了 Index。

Deref, DerefMut
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

Deref示例
:::::::::::::::::::::

智能指针解引用。
Box<T>、Rc<T> 等就是通过实现 Deref 来模拟指针的。

.. code-block:: rust
  :caption: Deref 的例子

  use std::ops::Deref;

  struct UserName(String);

  impl Deref for UserName {
      type Target = String;

      fn deref(&self) -> &Self::Target {
          &self.0
      }
  }

  fn main(){
      let name = UserName(String::from("Alice"));
      println!("{}", name.len());
      // 等价于
      println!("{}", name.0.len());
  }

Deref Trait + Deref Coercion(解引用强制转换)示例
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: rust

    use std::ops::Deref;

    struct MyBox<T>(T);

    impl<T> Deref for MyBox<T> {
        type Target = T;

        fn deref(&self) -> &Self::Target {
            &self.0
        }
    }

    fn hello(s: &str) {
        println!("{}", s);
    }

    fn main() {
        let x = MyBox("hello");

        hello(&x);
    }

``hello(&x);`` x 的类型是 ``MyBox<&str>``, 但是``hello()`` 需要的是 ``&str`` 。为什么能直接传？答案就是：

.. code-block:: text

    Deref Trait
    +
    Deref Coercion（自动解引用转换）

理解 MyBox<T>
""""""""""""""""""""""

这是一个 **Tuple Struct**：

.. code-block:: rust

    struct MyBox<T>(T);

例如：

.. code-block:: rust

    let x = MyBox("hello");

实际等价于：

.. code-block:: rust

    let x = MyBox::<&str>("hello");
    
内存：

.. code-block:: text
        
    x
    │
    └── "hello"

类型：

.. code-block:: rust

    MyBox<&str>

理解 Deref Trait
""""""""""""""""""""""
   
Rust 标准库：

.. code-block:: rust

    pub trait Deref {
        type Target;

        fn deref(&self) -> &Self::Target;
    }

作用：

.. code-block:: text

    告诉编译器：
    如何从当前类型得到内部引用

实现：

.. code-block:: rust

    impl<T> Deref for MyBox<T> {
        type Target = T;

        fn deref(&self) -> &T {
            &self.0
        }
    }

对于：

.. code-block:: rust

    MyBox<&str>

等价于：

.. code-block:: rust

    impl Deref for MyBox<&str> {
        type Target = &str;

        fn deref(&self) -> &&str {
            &self.0
        }
    }

注意返回值：

.. code-block:: rust
    
    &&str

为什么是 &&str
~~~~~~~~~~~~~~~~~~

假设：

.. code-block:: rust

    let x = MyBox("hello");

内部：

.. code-block:: text

    x.0

类型：

.. code-block:: text

    &str

而：

.. code-block:: text

    &self.0

是：

.. code-block:: text

    &&str

因为：

.. code-block:: text

    x.0       -> &str
    &x.0      -> &&str

所以：

.. code-block:: text

    x.deref()

返回：

.. code-block:: text

    &&str

hello(&x) 发生了什么
~~~~~~~~~~~~~~~~~~~~~~~~~

函数：

.. code-block:: rust

    fn hello(s: &str)

需要：

.. code-block:: rust

    &str

传入：

.. code-block:: rust

    hello(&x);

此时：

.. code-block:: rust

    &x

类型：

.. code-block:: rust

    &MyBox<&str>

编译器发现：

.. code-block:: text

    需要:

    &str

    实际:
    &MyBox<&str>

不匹配。

于是开始寻找：

.. code-block:: rust

    Deref

**第一次 Deref**
~~~~~~~~~~~~~~~~~~~~~~~~~

编译器自动调用：

.. code-block:: rust

    Deref::deref(&x)

即：

.. code-block:: rust

    x.deref()

得到：

.. code-block:: rust

    &&str

此时：

.. code-block:: text

    &MyBox<&str>
        ↓ Deref
    &&str

**第二次 Deref**
~~~~~~~~~~~~~~~~~~~~~~~~~

编译器继续发现：

.. code-block:: rust

    &&str

还不是：

.. code-block:: rust

    &str

于是再解一次引用：

.. code-block:: text

    &&str
    ↓
    &str

最终：

.. code-block:: rust

    hello(&x);

被编译器理解为：

.. code-block:: rust

    hello(*x.deref());

最终参数：

.. code-block:: rust

    &str

匹配成功。

编译器实际过程

.. code-block:: text

    hello(&x)

    &MyBox<&str>
        ↓
    Deref::deref()
        ↓
    &&str
        ↓
    自动解引用
        ↓
    &str
        ↓
    hello()

这就是：

.. code-block:: text

    Deref Coercion
    （解引用强制转换）


Deref 最核心用途
:::::::::::::::::::::

Rust 中以下类型都依赖 Deref：

.. code-block:: rust

    Box<T>
    Rc<T>
    Arc<T>
    String
    Vec<T>
    Cow<T>
    PathBuf
    OsString


DerefMut示例
:::::::::::::::::::::

.. code-block:: rust
  :caption: DerefMut 的例子

  use std::ops::{Deref, DerefMut};

  struct UserName(String);

  impl Deref for UserName {
      type Target = String;

      fn deref(&self) -> &String {
          &self.0
      }
  }

  impl DerefMut for UserName {
      fn deref_mut(&mut self) -> &mut String {
          &mut self.0
      }
  }

  fn main(){
      let mut name = UserName(String::from("Alice"));
      name.push_str(" Bob");
      println!("{}", name);
  }

格式化
----------------

Display
>>>>>>>>>>>>>>>>>>>>>>>

人类可读的格式化，用 {}。
例子：``println!("{}", my_struct)``

.. code-block:: rust

    use std::fmt;

    struct UserId(u64);

    // Display：用户友好的输出
    impl fmt::Display for UserId {
        fn fmt(
            &self,
            f: &mut fmt::Formatter<'_>,
        ) -> fmt::Result {
            write!(f, "{}", self.0)
        }
    }

    fn main(){
        let id = UserId(1234567890);
        println!("{}", id);
    }


Debug
>>>>>>>>>>>>>>>>>>>>>>>

调试格式化，用 {:?}。
几乎所有类型都会 #[derive(Debug)]。

.. code-block:: rust

    // Debug：调试输出（可通过 #[derive(Debug)] 自动实现）
    #[derive(Debug)]
    struct Point {
        x: i32,
        y: i32,
    }

    println!("{:?}", point); // Debug: Point { x: 3, y: 4 }

Write / Read（来自 std::io）
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

IO 写入/读取接口。
文件、网络流都实现了这些 trait。

集合
----------------------------------

Iterator
>>>>>>>>>>>>>>>>>>>>>>>

所有迭代器的核心 trait，提供 .next()。
for 循环、map、filter 都基于它。

IntoIterator
>>>>>>>>>>>>>>>>>>>>>>>

用于 for x in collection。
Vec<T> 同时实现了按值、按引用、按可变引用的 IntoIterator。

Extend
>>>>>>>>>>>>>>>>>>>>>>>

往集合里追加元素。
例子：vec.extend(&[1,2,3]);

FromIterator
>>>>>>>>>>>>>>>>>>>>>>>

把迭代器转成集合。
例子：let v: Vec<i32> = (0..5).collect();

并发 & 生命周期
----------------------

Send
>>>>>>>>>>>>>>>>>>>>>>>

类型能否安全地跨线程转移所有权。
大多数类型都是 Send，除了 Rc<T>。

Sync
>>>>>>>>>>>>>>>>>>>>>>>

类型能否安全地跨线程共享引用。
Arc<T> 是 Sync，Rc<T> 不是。

Drop
>>>>>>>>>>>>>>>>>>>>>>>

自定义资源释放逻辑，类似 C++ 的析构函数。
File、MutexGuard 都用它来自动清理资源。

Clone
>>>>>>>>>>>>>>>>>>>>>>>

显式复制一个值。和 Copy 不同，Clone 可以做深拷贝。

.. code-block:: rust

    // Clone：显式克隆
    #[derive(Clone)]
    struct Data {
        value: String,
    }

    let a = Data { value: String::from("hello") };
    let b = a.clone(); // 深拷贝


Copy
>>>>>>>>>>>>>>>>>>>>>>>

位拷贝（轻量类型）。数字、布尔、引用是 Copy，String、Vec 不是。

.. code-block:: rust

    // Copy：按位复制（隐式）
    #[derive(Copy, Clone)]
    struct Point {
        x: i32,
        y: i32,
    }

    let p1 = Point { x: 1, y: 2 };
    let p2 = p1; // 隐式复制，p1 仍然有效

比较 & 默认
---------------------

PartialEq / Eq
>>>>>>>>>>>>>>>>>>>>>>>

定义 ==、!=。Eq 代表完全等价，PartialEq 允许“部分等价”（比如浮点数 NaN）。

.. code-block:: rust

    #[derive(PartialEq, Eq)]
    struct Point {
        x: i32,
        y: i32,
    }

    let p1 = Point { x: 1, y: 2 };
    let p2 = Point { x: 1, y: 2 };

    assert_eq!(p1, p2); 

PartialOrd / Ord
>>>>>>>>>>>>>>>>>>>>>>>

定义排序比较 < > <= >=。
PartialOrd 允许不可比（NaN），Ord 代表全序。

Hash
>>>>>>>>>>>>>>>>>>>>>

定义哈希值，用于 HashMap、HashSet。

Default
>>>>>>>>>>>>>>>>>>>>>

提供一个默认值。
例子：let v: Vec<i32> = Default::default();