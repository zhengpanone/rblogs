============
基础工具
============

serde
=========

序列化/反序列化（JSON、YAML 等）。

serde_json
-------------

serde 的 JSON 支持。

regex
========
正则表达式。

chrono
=========
日期和时间处理。

uuid
=========

生成和解析 UUID。

rand
==========

随机数生成。


config、dotenvy
=========================

config 配置文件、dotenvy 环境变量

项目结构
---------------

.. code-block:: text

  src/
  ├─ config/
  │   ├─ mod.rs
  │   └─ application.toml
  ├─ main.rs
  .env

.env 文件
--------------------

.. code-block:: text
  :caption: .env

  APP_PORT=8080
  DATABASE_URL=mysql://root:123456@localhost/demo

application.toml
------------------------

.. code-block:: toml
  :caption: application.toml

  app_name = "hello-rust"

  [server]
  host = "0.0.0.0"
  port = 8080

配置结构体
-------------------

.. code-block:: rust
  :caption: src/config/mod.rs

  use serde::Deserialize;

  #[derive(Debug, Deserialize)]
  pub struct Settings {
      pub app_name: String,
      pub server: ServerConfig,
      pub database_url: String,
  }

  #[derive(Debug, Deserialize)]
  pub struct ServerConfig {
      pub host: String,
      pub port: u16,
  }


加载配置
-------------------

.. code-block:: rust
  :caption: src/main.rs

  use config::{Config, File, Environment};
  use dotenvy::dotenv;

  pub fn load_config() -> Result<Settings, config::ConfigError> {
      dotenv().ok();

      let settings = Config::builder()
          .add_source(File::with_name("src/config/application"))
          .add_source(Environment::default())
          .build()?;

      settings.try_deserialize()
  }

main.rs 
----------------

.. code-block:: rust
  :caption: src/main.rs

  mod config;

  use crate::config::load_config;

  fn main() {
      let settings = load_config().unwrap();

      println!("{:#?}", settings);
  }

运行结果
--------------------

.. code-block:: text

  Settings {
    app_name: "hello-rust",
    server: ServerConfig {
        host: "0.0.0.0",
        port: 8080,
    },
    database_url: "mysql://root:123456@localhost/demo",
  }

