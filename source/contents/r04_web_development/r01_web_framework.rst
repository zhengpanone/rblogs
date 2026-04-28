============================
Web 框架
============================

- **Actix-web**：高性能、功能完整、社区活跃
- **Axum**：基于Tokio的现代化框架，类型安全
- **Rocket**：零配置、开发友好、安全
- **Warp**：组合式、函数式编程风格
- **Tide**：异步、简洁的设计

示例
----------------------

Actix-web示例
======================

.. literalinclude:: code/r01_web_framework/actix-web-demo/src/main.rs
  :caption: main.rs
  :language: rust

.. literalinclude:: code/r01_web_framework/actix-web-demo/Cargo.toml
  :caption: Cargo.toml
  :language: toml

Axum示例
======================

.. literalinclude:: code/r01_web_framework/axum-demo/src/main.rs
  :caption: main.rs
  :language: rust

.. literalinclude:: code/r01_web_framework/axum-demo/Cargo.toml
  :caption: Cargo.toml
  :language: toml


tokio features解释

.. code-block:: toml

  tokio = { version = "1.52.1", features = [
    "rt-multi-thread", 
    "macros", 
    "net", 
    "io-util"
    ] }


- rt-multi-thread — 启用多线程运行时
- macros — 启用 #[tokio::main] 宏
- net — 网络功能（用于 TcpListener）
- io-util — I/O 工具函数




