=====================
网络 & HTTP
=====================

reqwest
==================

高级 HTTP 客户端。

hyper
==================

底层 HTTP 库。

tokio
==================

异步运行时。

async-std
==================

类似标准库风格的异步运行时。

warp
==================

基于 hyper 的 Web 框架。

actix-web
==================

高性能 Web 框架。


axum
==================

基于 tower 的 Web 框架。

创建一个简单的 Web 服务器
-------------------------------

.. literalinclude:: ./code/r03_Networking_HTTP/axum_demo/01.sample_http.rs
  :caption: 01.sample_http.rs
  :language: rust
