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



创建项目
-----------------

.. code-block:: console

  $ cargo new config_demo --vcs none
  $ cd config_demo
  $ cargo add config dotenvy
  $ cargo add serde -F derive

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
  use config::{Config, File, Environment};
  use dotenvy::dotenv;

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


  pub fn load_config() -> Result<Settings, config::ConfigError> {
      dotenv().ok();

      let settings = Config::builder()
            .set_default("database_url", "postgres://localhost:5432/mydb")
            .expect("Failed to set default value")
            .add_source(File::with_name("src/config/application"))
            .add_source(Environment::default())
            .build() 
            .expect("Failed to load configuration");

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

Environment 覆盖规则
------------------------

``Environment::default()：``

会自动读取：

.. code-block:: text

  APP_NAME=xxx
  SERVER__PORT=9090

双下划线：``__`` 表示嵌套。

推荐完整写法（生产）

.. code-block:: rust

  let settings = Config::builder()
      .add_source(File::with_name("config/default"))
      .add_source(File::with_name("config/local").required(false))
      .add_source(Environment::default().separator("__"))
      .build()?;


推荐目录结构
------------------------

.. code-block:: text

  config/
  ├─ default.toml
  ├─ dev.toml
  ├─ test.toml
  ├─ prod.toml


按环境加载
-------------------

.. code-block:: rust

  let profile = std::env::var("APP_ENV")
      .unwrap_or_else(|_| "dev".into());

  let settings = Config::builder()
      .add_source(File::with_name("config/default"))
      .add_source(File::with_name(&format!("config/{}", profile)))
      .build()?;