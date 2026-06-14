=======================
其他实用 Crate
=======================

不属于前面分类但在实际开发中非常有用的实用 crate 集合。

.. contents:: 目录
   :depth: 3
   :local:

rust-embed
============

将静态资源（HTML、CSS、JS、图片等）直接嵌入到可执行文件中，常用于 Web 服务或 CLI 工具，避免部署时携带静态文件。

.. code-block:: toml

   [dependencies]
   rust-embed = "8"

基础用法：

.. code-block:: rust

   use rust_embed::RustEmbed;

   #[derive(RustEmbed)]
   #[folder = "static/"]
   struct Asset;

   fn main() {
       // 获取嵌入的文件
       if let Some(file) = Asset::get("index.html") {
           let content = std::str::from_utf8(&file.data).unwrap();
           println!("index.html:\n{}", content);
       }

       // 遍历所有嵌入文件
       for file in Asset::iter() {
           println!("嵌入文件: {}", file);
       }
   }

与 Web 框架集成（axum）：

.. code-block:: rust

   use axum::{
       response::Html,
       routing::get,
       Router,
   };
   use rust_embed::RustEmbed;

   #[derive(RustEmbed)]
   #[folder = "static/"]
   struct Asset;

   async fn index() -> Html<&'static [u8]> {
       let file = Asset::get("index.html").unwrap();
       Html(file.data)
   }

   async fn serve_static(path: String) -> Option<(HeaderMap, &'static [u8])> {
       let file = Asset::get(&path)?;
       let mut headers = HeaderMap::new();
       // 根据扩展名设置 Content-Type
       if path.ends_with(".css") {
           headers.insert("Content-Type", "text/css".parse().unwrap());
       } else if path.ends_with(".js") {
           headers.insert("Content-Type", "application/javascript".parse().unwrap());
       }
       Some((headers, file.data))
   }

rust-embed 常用功能：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - 功能
     - 说明
   * - ``#[folder = "path/"]``
     - 指定要嵌入的目录
   * - ``Asset::get(name)``
     - 获取单个文件，返回 ``Option<EmbeddedFile>``
   * - ``Asset::iter()``
     - 遍历所有嵌入文件路径
   * - ``file.data``
     - 文件内容（``Cow<'static, [u8]>``）
   * - ``file.metadata``
     - 文件元数据（last_modified 等）
   * - ``file.metadata.sha256_hash()``
     - 文件的 SHA-256 哈希
   * - ``#[include = "*.html"]``
     - 只嵌入匹配模式的文件
   * - ``#[exclude = "*.map"]``
     - 排除特定文件

lazy_static / once_cell
=========================

延迟初始化的静态变量。``once_cell`` 已部分进入标准库（``std::sync::OnceLock`` / ``std::sync::LazyLock``，Rust 1.80+），但 ``lazy_static`` 和 ``once_cell`` 仍在大量项目中广泛使用。

.. code-block:: toml

   [dependencies]
   lazy_static = "1"
   once_cell = "1"

lazy_static：

.. code-block:: rust

   use lazy_static::lazy_static;
   use std::collections::HashMap;
   use std::sync::Mutex;

   lazy_static! {
       // 编译期不可计算的常量
       static ref CONFIG: HashMap<String, String> = {
           let mut m = HashMap::new();
           m.insert("host".to_string(), "localhost".to_string());
           m.insert("port".to_string(), "8080".to_string());
           m
       };

       // 全局可变状态
       static ref COUNTER: Mutex<u64> = Mutex::new(0);

       // 正则表达式（编译开销大，只做一次）
       static ref RE_EMAIL: regex::Regex =
           regex::Regex::new(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$").unwrap();
   }

   fn main() {
       println!("Host: {}", CONFIG.get("host").unwrap());

       let mut count = COUNTER.lock().unwrap();
       *count += 1;
       println!("Count: {}", count);

       assert!(RE_EMAIL.is_match("user@example.com"));
       assert!(!RE_EMAIL.is_match("invalid-email"));
   }

once_cell（推荐，更接近标准库 API）：

.. code-block:: rust

   use once_cell::sync::{Lazy, OnceCell};
   use once_cell::unsync::OnceCell as UnsyncOnceCell;
   use std::collections::HashMap;

   // Lazy: 惰性初始化（类似 lazy_static）
   static CONFIG: Lazy<HashMap<String, String>> = Lazy::new(|| {
       let mut m = HashMap::new();
       m.insert("host".to_string(), "localhost".to_string());
       m.insert("port".to_string(), "8080".to_string());
       m
   });

   // OnceCell: 一次性赋值（可非惰性）
   static INSTANCE: OnceCell<String> = OnceCell::new();

   fn main() {
       // 使用 Lazy
       println!("{}", CONFIG.get("host").unwrap());

       // 使用 OnceCell
       INSTANCE.set("initialized".to_string()).unwrap();
       println!("{}", INSTANCE.get().unwrap());

       // 只能设置一次
       assert!(INSTANCE.set("another".to_string()).is_err());
   }

   // 单线程 OnceCell
   fn unsync_example() {
       let cell = UnsyncOnceCell::new();
       cell.set("value").unwrap();
       assert_eq!(cell.get(), Some(&"value"));
   }

标准库替代（Rust 1.80+）：

.. code-block:: rust

   use std::sync::{OnceLock, LazyLock};

   static CELL: OnceLock<String> = OnceLock::new();
   static LAZY: LazyLock<String> = LazyLock::new(|| "computed".to_string());

   fn main() {
       CELL.set("hello".to_string()).unwrap();
       println!("{} {}", CELL.get().unwrap(), &*LAZY);
   }

三种方案对比：

.. list-table::
   :header-rows: 1
   :widths: 20 30 30

   * - 特性
     - lazy_static
     - once_cell / std
   * - 语法
     - ``lazy_static! { static ref X: T = expr; }``
     - ``static X: Lazy<T> = Lazy::new(|| expr)``
   * - 宏依赖
     - 需要宏
     - 无需宏（Lazy::new 是普通函数）
   * - 是否标准库
     - 第三方
     - once_cell 接近标准库 API，std 直接内置
   * - 一次性赋值
     - 不支持
     - OnceCell::set()
   * - 运行时检查
     - Deref 时 panic（若初始化失败）
     - get() 返回 Option
   * - 推荐程度
     - 遗留项目仍可用
     - 新项目优先使用 std::sync::LazyLock / OnceLock

num_cpus
=============

获取系统 CPU 核心数，用于配置线程池大小。

.. code-block:: toml

   [dependencies]
   num_cpus = "1"

.. code-block:: rust

   fn main() {
       // 物理核心数
       let physical = num_cpus::get();
       println!("物理核心: {}", physical);

       // 逻辑核心数（含超线程）
       let logical = num_cpus::get_physical();
       println!("逻辑核心: {}", logical);

       // 用于线程池配置
       let pool = rayon::ThreadPoolBuilder::new()
           .num_threads(num_cpus::get())
           .build()
           .unwrap();
   }

bytes
======

高性能字节缓冲区，零拷贝切片。是网络编程和序列化的基础类型。

.. code-block:: toml

   [dependencies]
   bytes = "1"

.. code-block:: rust

   use bytes::{Bytes, BytesMut, Buf, BufMut};

   fn main() {
       // Bytes: 不可变字节缓冲区（共享所有权，零拷贝切片）
       let mut b1 = Bytes::from("hello ");
       let b2 = Bytes::from("world");

       // 拼接（未达到内联阈值时零拷贝）
       let b3 = [b1.clone(), b2.clone()].concat();
       println!("{:?}", b3); // b"hello world"

       // 切片不复制数据
       let slice = b3.slice(0..5);
       println!("{:?}", slice); // b"hello"

       // BytesMut: 可变字节缓冲区
       let mut buf = BytesMut::with_capacity(64);
       buf.put_slice(b"hello");
       buf.put_u8(b' ');
       buf.put_slice(b"world");
       println!("{:?}", buf); // b"hello world"

       // 转换为 Bytes
       let frozen = buf.freeze();
       println!("{:?}", frozen);
   }

常用 API：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - API
     - 说明
   * - ``Bytes::from(&str)`` / ``Bytes::from_static(b"..")``
     - 创建 Bytes
   * - ``Bytes::slice(range)``
     - 零拷贝切片
   * - ``Bytes::copy_from_slice(&[u8])``
     - 拷贝创建
   * - ``[b1, b2].concat()``
     - 拼接多个 Bytes
   * - ``BytesMut::put_slice / put_u8 / put_i32``
     - 写入数据
   * - ``BytesMut::freeze()``
     - 转为不可变 Bytes
   * - ``BytesMut::split_to(n)`` / ``split_off(n)``
     - 分割缓冲区

derive_more / derive_builder
=============================

减少样板代码的 derive 宏。

.. code-block:: toml

   [dependencies]
   derive_more = "1"
   derive_builder = "0.20"

derive_more —— 自动实现常见 trait：

.. code-block:: rust

   use derive_more::{
       Display, Debug, From, Into, Add, Mul, Constructor,
       Deref, DerefMut, AsRef, AsMut, Index, IndexMut,
   };

   // Display + Debug
   #[derive(Display, Debug)]
   #[display("Error: {_0}")]
   struct MyError(String);

   // From 自动转换
   #[derive(From, Debug)]
   struct MyInt(i32);

   #[derive(From, Debug)]
   struct MyFloat(f64);

   // Into
   #[derive(Into, Debug)]
   struct Wrapper(String);

   // 算术运算
   #[derive(Add, Mul, Debug, Clone, Copy)]
   struct Point {
       x: i32,
       y: i32,
   }

   fn main() {
       // Display
       let err = MyError("something went wrong".to_string());
       println!("{}", err); // Error: something went wrong

       // From
       let n: MyInt = 42.into();
       let f: MyFloat = 3.14.into();

       // Into
       let w = Wrapper("hello".to_string());
       let s: String = w.into();

       // 算术
       let p1 = Point { x: 1, y: 2 };
       let p2 = Point { x: 3, y: 4 };
       let sum = p1 + p2;
       println!("Sum: {:?}, Mul: {:?}", sum, p1 * p2);
   }

derive_builder —— 自动生成 Builder 模式：

.. code-block:: rust

   use derive_builder::Builder;

   #[derive(Builder, Debug)]
   struct ServerConfig {
       host: String,
       port: u16,
       #[builder(default = "4")]
       workers: usize,
       #[builder(default = "false")]
       tls: bool,
       #[builder(setter(into, strip_option), default)]
       cert_path: Option<String>,
   }

   fn main() {
       let config = ServerConfigBuilder::default()
           .host("localhost")
           .port(8080)
           .workers(8)
           .tls(true)
           .cert_path(Some("/etc/ssl/cert.pem"))
           .build()
           .unwrap();

       println!("{:?}", config);
       // ServerConfig { host: "localhost", port: 8080, workers: 8, tls: true, cert_path: Some("/etc/ssl/cert.pem") }
   }

derive_more 常用 derive：

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - Derive
     - 实现
   * - ``Display``
     - ``std::fmt::Display``，支持 ``#[display("fmt", arg)]``
   * - ``Debug``
     - ``std::fmt::Debug``
   * - ``From`` / ``Into``
     - 类型转换
   * - ``Add`` / ``Sub`` / ``Mul`` / ``Div``
     - 算术运算符
   * - ``Deref`` / ``DerefMut``
     - 解引用（newtype 模式）
   * - ``AsRef`` / ``AsMut``
     - 引用转换
   * - ``Constructor``
     - 生成 ``new()`` 方法
   * - ``Index`` / ``IndexMut``
     - 索引操作
   * - ``Error``
     - ``std::error::Error``（配合 Display）

derive_builder 常用属性：

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - 属性
     - 说明
   * - ``#[builder(default)]``
     - 使用 Default::default()
   * - ``#[builder(default = "expr")]``
     - 指定默认值表达式
   * - ``#[builder(setter(into))]``
     - setter 自动调用 .into()，接受多种类型
   * - ``#[builder(setter(strip_option))]``
     - ``Option<T>`` 字段的 setter 直接接受 T
   * - ``#[builder(setter(skip))]``
     - 跳过该字段的 setter
   * - ``#[builder(try_setter)]``
     - setter 返回 Result
   * - ``#[builder(pattern = "owned")]``
     - setter 接受 owned 值而非引用

strum
======

枚举工具集，提供枚举的字符串转换、迭代等功能。

.. code-block:: toml

   [dependencies]
   strum = { version = "0.26", features = ["derive"] }

.. code-block:: rust

   use strum::{EnumString, EnumIter, EnumCount, EnumVariantNames, Display, IntoStaticStr, EnumProperty};

   #[derive(Debug, PartialEq, EnumString, Display, IntoStaticStr, EnumIter, EnumCount)]
   enum Color {
       #[strum(serialize = "red", to_string = "红色")]
       Red,
       #[strum(serialize = "green", to_string = "绿色")]
       Green,
       #[strum(serialize = "blue", to_string = "蓝色")]
       Blue,
   }

   fn main() {
       // 字符串 -> 枚举
       let color: Color = "red".parse().unwrap();
       assert_eq!(color, Color::Red);

       // 枚举 -> 字符串（Display）
       println!("{}", Color::Green); // 绿色

       // 枚举 -> &'static str
       let s: &'static str = Color::Blue.into();
       println!("{}", s); // blue

       // 遍历所有变体
       for variant in Color::iter() {
           println!("{:?}", variant);
       }

       // 变体数量
       println!("共 {} 种颜色", Color::COUNT); // 3
   }

   // EnumMessage: 为每个变体关联消息
   #[derive(strum::EnumMessage)]
   enum HttpStatus {
       #[strum(message = "OK")]
       Ok,
       #[strum(message = "Not Found")]
       NotFound,
       #[strum(message = "Internal Server Error")]
       InternalServerError,
   }

   fn main() {
       println!("{}", HttpStatus::NotFound.get_message().unwrap()); // Not Found
   }

strum 常用 derive：

.. list-table::
   :header-rows: 1
   :widths: 25 55

   * - Derive
     - 功能
   * - ``Display`` / ``ToString``
     - 枚举 -> 字符串
   * - ``EnumString``
     - 字符串 -> 枚举（parse）
   * - ``EnumIter``
     - 遍历所有变体
   * - ``EnumCount``
     - 变体总数（``::COUNT``）
   * - ``IntoStaticStr``
     - 转为 ``&'static str``
   * - ``EnumMessage``
     - 为变体关联消息文本
   * - ``EnumProperty``
     - 为变体关联键值属性
   * - ``AsRefStr``
     - ``as_ref()`` 返回 ``&str``

num
====

数值 trait 和类型扩展。

.. code-block:: toml

   [dependencies]
   num = "0.4"

.. code-block:: rust

   use num::{Integer, Num, BigInt, rational::Ratio, complex::Complex};
   use num::traits::{Zero, One, Signed, Bounded, NumCast, ToPrimitive, FromPrimitive, Saturating, WrappingAdd};

   fn main() {
       // BigInt: 任意精度整数
       let a: BigInt = "12345678901234567890".parse().unwrap();
       let b: BigInt = "98765432109876543210".parse().unwrap();
       println!("a + b = {}", &a + &b);

       // Ratio: 有理数
       let half = Ratio::new(1, 2);
       let third = Ratio::new(1, 3);
       println!("1/2 + 1/3 = {}", half + third); // 5/6

       // Complex: 复数
       let c1 = Complex::new(1.0, 2.0);
       let c2 = Complex::new(3.0, 4.0);
       println!("({}) + ({}) = {}", c1, c2, c1 + c2);

       // 通用数值函数
       fn generic_sum<T: Num + Copy>(a: T, b: T) -> T {
           a + b
       }
       println!("{}", generic_sum(1, 2));
       println!("{}", generic_sum(1.5, 2.5));
   }

   // 带约束的泛型
   fn factorial<T: Integer + Clone + MulAssign>(n: T) -> T {
       let mut result = T::one();
       let mut i = T::one();
       while i <= n {
           result *= i.clone();
           i = i + T::one();
       }
       result
   }

num 主要模块：

.. list-table::
   :header-rows: 1
   :widths: 20 60

   * - 模块
     - 功能
   * - ``num::BigInt`` / ``num::BigUint``
     - 任意精度整数
   * - ``num::rational::Ratio<T>``
     - 有理数（分数）
   * - ``num::complex::Complex<T>``
     - 复数
   * - ``num::traits``
     - Zero, One, Num, Integer, Signed, Bounded, NumCast 等
   * - ``num::iter``
     - range_step 等数值迭代器

总结
==========

其他实用 Crate 总览：

.. list-table::
   :header-rows: 1
   :widths: 18 15 45 22

   * - Crate
     - 定位
     - 核心能力
     - 使用场景
   * - ``rust-embed``
     - 静态资源嵌入
     - 编译期将文件嵌入二进制
     - Web 服务、CLI 单文件部署
   * - ``lazy_static`` / ``once_cell``
     - 延迟初始化
     - 运行时一次性初始化的全局变量
     - 全局配置、正则、缓存
   * - ``num_cpus``
     - CPU 信息
     - 获取物理/逻辑核心数
     - 线程池大小配置
   * - ``bytes``
     - 字节缓冲区
     - 零拷贝切片、共享所有权
     - 网络编程、序列化框架
   * - ``derive_more``
     - 样板代码减少
     - 自动实现 Display/From/Add/Deref 等
     - newtype 模式、错误类型
   * - ``derive_builder``
     - Builder 模式
     - 自动生成 Builder
     - 复杂配置对象的构建
   * - ``strum``
     - 枚举工具
     - 字符串互转、迭代、计数
     - CLI 参数解析、配置序列化
   * - ``num``
     - 数值扩展
     - 大整数、有理数、复数、数值 trait
     - 数学计算、科学计算
