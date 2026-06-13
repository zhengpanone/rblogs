==================================
主流编程风格 / 编程范式
==================================

编程语言的发展过程中，形成了多种主流编程风格（Programming Style）和编程范式（Programming Paradigm）。

可以理解为：

- **编程风格（Style）**：代码怎么写
- **编程范式（Paradigm）**：程序怎么组织

命令式编程（Imperative Programming）
----------------------------------------

核心思想：告诉计算机一步一步怎么做，通过修改状态、执行语句来改变程序状态。

特点:

- 关注“执行过程”
- 大量使用赋值、循环、条件判断
- 状态可变（mutable state）

示例：

.. code-block:: java

    int sum = 0;

    for (int i = 1; i <= 100; i++) {
        sum += i;
    }

执行过程：

.. code-block:: text

    初始化
      ↓
    循环
      ↓
    累加
      ↓
    结束

代表语言：

- C
- Pascal
- Fortran
- 早期 BASIC
- 汇编语言

过程式编程（Procedural Programming）
----------------------------------------

核心思想：将代码封装成 **过程和函数**，按步骤调用。

特点:

- 是命令式编程的重要分支
- 强调“过程 / 函数”的组织
- 比纯命令式更结构化
- 函数划分模块
- 数据和行为分离

示例：

.. code-block:: python

    def calculate_sum(n):
      total = 0
      for i in range(1, n + 1):
          total += i
      return total

代表语言：

- C
- Go
- Python
- Rust（支持）

面向对象编程（Object-Oriented Programming）
------------------------------------------------

核心思想：将数据和操作数据的行为封装成对象，通过对象之间的交互完成任务。

把：

.. code-block:: text

    数据
    +
    行为

封装到一起。

**四大特征**:

- 封装（Encapsulation）
- 继承（Inheritance）
- 多态（Polymorphism）
- 抽象（Abstraction）

特点

- 以“对象”为中心
- 强调职责划分
- 适合大型、复杂系统的建模

示例：

.. code-block:: java

    public class User {

        private String name;

        public void login() {
            System.out.println("登录");
        }
    }

使用：

.. code-block:: java

    User user = new User();
    user.login();

OOP 四大特性
~~~~~~~~~~~~

封装（Encapsulation）
^^^^^^^^^^^^^^^^^^^^^

.. code-block:: java

    private String name;

隐藏实现细节。

继承（Inheritance）
^^^^^^^^^^^^^^^^^^^

.. code-block:: java

    class Animal {}

    class Dog extends Animal {}

实现代码复用。

多态（Polymorphism）
^^^^^^^^^^^^^^^^^^^^

.. code-block:: java

    Animal animal = new Dog();

同一接口不同实现。

抽象（Abstraction）
^^^^^^^^^^^^^^^^^^^

.. code-block:: java

    interface Payment {
        void pay();
    }

关注能力，不关注实现。

代表语言：

- Java
- C#
- C++
- Kotlin
- Python（多范式）
- Ruby

函数式编程（Functional Programming）
----------------------------------------

核心思想：

    函数是一等公民。

函数可以：

- 赋值
- 传参
- 返回

示例：

.. code-block:: java

    Function<Integer, Integer> addOne =
        x -> x + 1;

纯函数
~~~~~~

相同输入永远得到相同输出。

.. code-block:: java

    int add(int a, int b) {
        return a + b;
    }

不可变数据
~~~~~~~~~~

不推荐：

.. code-block:: java

    list.add("A");

推荐：

.. code-block:: java

    List<String> newList =
        Stream.concat(
            list.stream(),
            Stream.of("A")
        ).toList();

高阶函数
~~~~~~~~

函数接收函数：

.. code-block:: java

    list.stream()
        .filter(x -> x > 10)
        .map(x -> x * 2)
        .forEach(System.out::println);

代表语言：

- Haskell
- Lisp
- Scala
- F#
- Elixir

现代语言支持：

- Java Stream
- Kotlin
- Rust
- JavaScript

声明式编程（Declarative Programming）
-----------------------------------------

核心思想： 只说明想要什么结果，不关心具体执行过程。

特点:

- 关注“做什么”，而非“怎么做”
- 代码更接近问题描述
- 通常可读性更强

命令式：

.. code-block:: java

    List<String> result = new ArrayList<>();

    for (String s : list) {
        if (s.length() > 5) {
            result.add(s);
        }
    }

声明式：

.. code-block:: java

    list.stream()
        .filter(s -> s.length() > 5)
        .toList();

SQL 示例：

.. code-block:: sql

    SELECT *
    FROM user
    WHERE age > 18;

你无需关心：

.. code-block:: text

    索引扫描
    排序
    执行计划

数据库负责执行。

响应式编程（Reactive Programming）
--------------------------------------

核心思想：

    数据变化驱动程序执行。

传统方式：

.. code-block:: java

    String data = getData();
    process(data);

响应式方式：

.. code-block:: java

    Flux.just("A", "B", "C")
        .map(String::toLowerCase)
        .subscribe(System.out::println);

代表技术：

- Reactor
- RxJava
- Vue
- React

Spring Boot：

.. code-block:: text

    Spring WebFlux

属于典型响应式框架。

事件驱动编程（Event Driven Programming）
--------------------------------------------

核心思想：程序流程由外部事件（点击、消息、信号）驱动。

特点:

- 以事件为中心
- 大量使用回调函数 / 监听器
- 常用于 GUI、服务器、前端开发

示例：

.. code-block:: java

    button.addActionListener(
        e -> System.out.println("点击")
    );

流程：

.. code-block:: text

    点击按钮
        ↓
    触发事件
        ↓
    执行回调

代表领域：

- GUI
- 浏览器
- Node.js

并发 / 并行编程（Concurrent & Parallel Programming）
----------------------------------------------------------

核心思想：处理多个任务同时执行的问题。

特点:

- 关注线程、进程、通信、同步
- 解决竞态、死锁等问题
- 常与 OOP / FP 结合

典型模型:

- 多线程（Java Thread）
- Actor 模型（Erlang、Akka）
- CSP（Go goroutine + channel）

传统方式：

.. code-block:: java

    synchronized

Actor 模型：

.. code-block:: text

    Actor A
        ↓ 消息
    Actor B

每个 Actor 包含：

.. code-block:: text

    状态
    邮箱
    处理逻辑

代表语言：

- Erlang
- Elixir

代表框架：

- Akka

数据导向编程（Data-Oriented Programming）
----------------------------------------------------------

核心思想：

    数据结构优先，而不是对象优先。

OOP 风格：

.. code-block:: java

    User.login()
    User.logout()
    User.changePassword()

DOP 风格：

.. code-block:: java

    Map<String,Object>
    Record
    Struct

逻辑：

.. code-block:: java

    process(userData);

而不是：

.. code-block:: java

    user.process();

代表语言：

- C
- Rust
- Go

典型领域：

- 游戏开发
- 高性能计算

ECS（Entity Component System）
----------------------------------

游戏开发热门范式。

结构：

.. code-block:: text

    Entity
        ↓
    Component
        ↓
    System

示例：

.. code-block:: text

    Player
     ├─ Position
     ├─ Health
     └─ Weapon

System：

.. code-block:: text

    MoveSystem
    RenderSystem
    AttackSystem

代表框架：

- Bevy（Rust）
- Unity DOTS
- Unreal Mass

多范式编程（Multi-Paradigm）
----------------------------------

现代语言通常不只支持一种风格，而是混合多种范式：

.. list-table:: 语言与支持范式
   :widths: 25 75
   :header-rows: 1

   * - 语言
     - 支持范式
   * - Python
     - OOP + 函数式 + 过程式
   * - JavaScript
     - 事件驱动 + 函数式 + OOP
   * - Scala
     - OOP + 函数式
   * - Rust
     - 命令式 + 函数式 + 泛型
   * - C++
     - 过程式 + OOP + 泛型

Rust 常见范式定位
-----------------------

Rust 既不是纯 OOP，也不是纯函数式语言。

更接近：

.. code-block:: text

    过程式
    +
    泛型编程
    +
    数据导向
    +
    函数式
    +
    所有权模型

典型 Rust 风格：

.. code-block:: rust

    struct User {
        name: String,
    }

    impl User {
        fn new(name: String) -> Self {
            Self { name }
        }
    }

配合：

.. code-block:: text

    Iterator
    Trait
    Enum
    Pattern Matching
    Ownership

形成独特的 Rust 编程模型。

现代后端开发主流范式
--------------------------

以 Spring Boot 为例：

.. code-block:: text

    OOP
    +
    DDD
    +
    事件驱动
    +
    响应式
    +
    函数式（Stream）
    +
    微服务

典型技术栈：

.. code-block:: text

    Spring Boot 3
    Spring Security
    Spring Authorization Server
    Flowable
    MyBatis-Plus
    Vue

推荐学习顺序：

.. code-block:: text

    过程式
        ↓
    OOP
        ↓
    设计模式
        ↓
    函数式编程
        ↓
    DDD
        ↓
    响应式编程（WebFlux）
        ↓
    Actor 模型
        ↓
    Rust 数据导向编程

如果重点学习 Rust，建议优先掌握：

.. code-block:: text

    1. Ownership（所有权）
    2. RAII
    3. Trait-Oriented Programming
    4. Generic Programming
    5. Functional Programming
    6. Data-Oriented Programming
    7. ECS

这条路线比直接学习复杂框架更容易建立正确的 Rust 思维模型。