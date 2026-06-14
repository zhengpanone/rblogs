==================================
数据库 & 持久化 Crate
==================================

Rust 生态中主流的数据库和持久化方案，覆盖关系型数据库、ORM、NoSQL 和缓存。

.. contents:: 目录
   :depth: 3
   :local:

sqlx
==========

异步 SQL 工具，支持 PostgreSQL、MySQL、SQLite。编译期校验 SQL 语法，无需 DSL 或 ORM。

特点：

- 编译期检查 SQL 语法（通过 ``query!`` / ``query_as!`` 宏）
- 异步原生（基于 tokio / async-std）
- 支持连接池（内置 ``PgPool`` / ``MySqlPool`` / ``SqlitePool``）
- 支持 Migration

添加依赖：

.. code-block:: toml

   [dependencies]
   sqlx = { version = "0.8", features = ["runtime-tokio", "postgres", "uuid", "chrono"] }
   tokio = { version = "1", features = ["full"] }

连接与查询：

.. code-block:: rust

   use sqlx::postgres::PgPoolOptions;
   use sqlx::FromRow;

   #[derive(Debug, FromRow)]
   struct User {
       id: i32,
       name: String,
       email: String,
   }

   #[tokio::main]
   async fn main() -> Result<(), sqlx::Error> {
       // 创建连接池
       let pool = PgPoolOptions::new()
           .max_connections(5)
           .connect("postgres://user:password@localhost/mydb")
           .await?;

       // 查询多行
       let users = sqlx::query_as::<_, User>("SELECT id, name, email FROM users")
           .fetch_all(&pool)
           .await?;

       for user in &users {
           println!("{}: {} <{}>", user.id, user.name, user.email);
       }

       // 查询单行
       let user = sqlx::query_as::<_, User>("SELECT id, name, email FROM users WHERE id = $1")
           .bind(1)
           .fetch_one(&pool)
           .await?;

       // 插入
       let result = sqlx::query("INSERT INTO users (name, email) VALUES ($1, $2)")
           .bind("Alice")
           .bind("alice@example.com")
           .execute(&pool)
           .await?;
       println!("插入行数: {}, 新 ID: {}", result.rows_affected(), result.last_insert_id());

       // 更新
       sqlx::query("UPDATE users SET name = $1 WHERE id = $2")
           .bind("Alice Updated")
           .bind(1)
           .execute(&pool)
           .await?;

       // 删除
       sqlx::query("DELETE FROM users WHERE id = $1")
           .bind(1)
           .execute(&pool)
           .await?;

       Ok(())
   }

编译期检查（``query!`` 宏）：

.. code-block:: rust

   // 编译时检查 SQL 语法和返回列类型
   let users = sqlx::query!("SELECT id, name, email FROM users WHERE active = $1", true)
       .fetch_all(&pool)
       .await?;

   // 如果 users 表中没有 active 列，编译直接报错

事务：

.. code-block:: rust

   let mut tx = pool.begin().await?;

   sqlx::query("INSERT INTO users (name, email) VALUES ($1, $2)")
       .bind("Bob")
       .bind("bob@example.com")
       .execute(&mut *tx)
       .await?;

   sqlx::query("UPDATE accounts SET balance = balance - 100 WHERE user_id = $1")
       .bind(1)
       .execute(&mut *tx)
       .await?;

   tx.commit().await?;

Migration：

.. code-block:: console

   $ cargo install sqlx-cli
   $ sqlx migrate add create_users_table

   # 生成 migrations/20240101000000_create_users_table.sql
   $ sqlx migrate run --database-url postgres://...

.. code-block:: sql

   -- migrations/20240101000000_create_users_table.sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       name VARCHAR NOT NULL,
       email VARCHAR NOT NULL UNIQUE,
       created_at TIMESTAMP DEFAULT NOW()
   );

常用 API：

.. list-table:: sqlx 常用 API
   :header-rows: 1
   :widths: 30 70

   * - 方法
     - 说明
   * - ``sqlx::query(sql)``
     - 原始查询，返回影响行数
   * - ``sqlx::query_as::<_, T>(sql)``
     - 查询并映射到结构体（运行时）
   * - ``sqlx::query!(sql, ...)``
     - 编译期检查 SQL 和列类型
   * - ``sqlx::query_as!(T, sql, ...)``
     - 编译期检查 + 映射到结构体
   * - ``.bind(value)``
     - 绑定参数
   * - ``.fetch_all(&pool)``
     - 获取所有行
   * - ``.fetch_one(&pool)``
     - 获取单行
   * - ``.fetch_optional(&pool)``
     - 获取可选行（返回 ``Option``）
   * - ``.execute(&pool)``
     - 执行 INSERT/UPDATE/DELETE

diesel
==========

Rust 的同步 ORM 框架，编译期类型安全的查询构建器。支持 PostgreSQL、MySQL、SQLite。

特点：

- 完全类型安全的查询 DSL
- 编译期验证 schema
- 强大的关联和 join 支持
- 同步 API（适合非异步场景）

安装 CLI：

.. code-block:: console

   $ cargo install diesel_cli --no-default-features --features postgres
   $ diesel setup --database-url postgres://user:password@localhost/mydb
   $ diesel migration generate create_users

定义 Schema（自动生成）：

.. code-block:: rust

   // src/schema.rs (由 diesel CLI 自动生成)
   diesel::table! {
       users (id) {
           id -> Int4,
           name -> Varchar,
           email -> Varchar,
           created_at -> Timestamp,
       }
   }

定义 Model：

.. code-block:: rust

   // src/models.rs
   use diesel::prelude::*;
   use chrono::NaiveDateTime;

   #[derive(Queryable, Selectable)]
   #[diesel(table_name = crate::schema::users)]
   pub struct User {
       pub id: i32,
       pub name: String,
       pub email: String,
       pub created_at: NaiveDateTime,
   }

   #[derive(Insertable)]
   #[diesel(table_name = crate::schema::users)]
   pub struct NewUser<'a> {
       pub name: &'a str,
       pub email: &'a str,
   }

CRUD 操作：

.. code-block:: rust

   use diesel::prelude::*;
   use diesel::pg::PgConnection;

   fn establish_connection() -> PgConnection {
       let database_url = std::env::var("DATABASE_URL")
           .expect("DATABASE_URL must be set");
       PgConnection::establish(&database_url)
           .expect("Error connecting to database")
   }

   fn main() {
       use crate::schema::users::dsl::*;

       let mut conn = establish_connection();

       // 创建
       let new_user = NewUser { name: "Alice", email: "alice@example.com" };
       diesel::insert_into(users)
           .values(&new_user)
           .execute(&mut conn)
           .unwrap();

       // 查询所有
       let all_users: Vec<User> = users.load(&mut conn).unwrap();

       // 条件查询
       let alice: User = users
           .filter(name.eq("Alice"))
           .first(&mut conn)
           .unwrap();

       // 更新
       diesel::update(users.find(1))
           .set(name.eq("Alice Updated"))
           .execute(&mut conn)
           .unwrap();

       // 删除
       diesel::delete(users.filter(email.like("%@example.com")))
           .execute(&mut conn)
           .unwrap();
   }

sqlx vs diesel：

.. list-table:: sqlx vs diesel
   :header-rows: 1
   :widths: 25 40 35

   * - 特性
     - ``sqlx``
     - ``diesel``
   * - 编程模型
     - 手写 SQL
     - Query Builder DSL
   * - 异步
     - 原生异步
     - 同步（可配合 ``spawn_blocking``）
   * - SQL 校验
     - 编译期（``query!`` 宏）
     - 编译期（DSL 类型检查）
   * - 学习曲线
     - 较低（会 SQL 即可）
     - 较高（需学 DSL + 概念）
   * - ORM 特性
     - 轻量（偏 SQL 工具）
     - 完整 ORM（关联、join 等）
   * - 适用场景
     - 复杂 SQL、异步架构
     - 偏好 ORM、同步架构

mongodb
==========

MongoDB 官方 Rust 驱动，支持异步操作和 BSON。

.. code-block:: toml

   [dependencies]
   mongodb = "3"
   tokio = { version = "1", features = ["full"] }
   serde = { version = "1", features = ["derive"] }

基本使用：

.. code-block:: rust

   use mongodb::{Client, Collection, bson::doc};
   use serde::{Serialize, Deserialize};
   use mongodb::options::ClientOptions;

   #[derive(Debug, Serialize, Deserialize)]
   struct User {
       #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
       id: Option<mongodb::bson::oid::ObjectId>,
       name: String,
       email: String,
       age: i32,
   }

   #[tokio::main]
   async fn main() -> mongodb::error::Result<()> {
       // 连接
       let client = Client::with_uri_str("mongodb://localhost:27017").await?;
       let db = client.database("mydb");
       let collection: Collection<User> = db.collection("users");

       // 插入
       let new_user = User {
           id: None,
           name: String::from("Alice"),
           email: String::from("alice@example.com"),
           age: 28,
       };
       let result = collection.insert_one(new_user).await?;
       println!("插入 ID: {:?}", result.inserted_id);

       // 查询单条
       let user = collection
           .find_one(doc! { "name": "Alice" })
           .await?
           .expect("未找到");
       println!("找到: {} <{}>", user.name, user.email);

       // 查询多条（Cursor）
       use mongodb::bson::doc;
       use futures::TryStreamExt;

       let mut cursor = collection.find(doc! { "age": { "$gte": 18 } }).await?;
       while let Some(user) = cursor.try_next().await? {
           println!("{} - {}", user.name, user.age);
       }

       // 更新
       collection
           .update_one(
               doc! { "name": "Alice" },
               doc! { "$set": { "age": 29 } },
           )
           .await?;

       // 删除
       collection
           .delete_one(doc! { "name": "Alice" })
           .await?;

       Ok(())
   }

常用 API：

.. list-table:: mongodb 常用 API
   :header-rows: 1
   :widths: 35 65

   * - 方法
     - 说明
   * - ``collection.insert_one(doc)``
     - 插入单条文档
   * - ``collection.insert_many(docs)``
     - 批量插入
   * - ``collection.find_one(filter)``
     - 查询单条
   * - ``collection.find(filter)``
     - 查询多条（返回 Cursor）
   * - ``collection.update_one(filter, update)``
     - 更新单条
   * - ``collection.update_many(filter, update)``
     - 批量更新
   * - ``collection.delete_one(filter)``
     - 删除单条
   * - ``collection.delete_many(filter)``
     - 批量删除
   * - ``collection.count_documents(filter)``
     - 计数

redis
==========

Redis 客户端，支持异步和连接池。

.. code-block:: toml

   [dependencies]
   redis = { version = "0.25", features = ["tokio-comp", "connection-manager"] }
   tokio = { version = "1", features = ["full"] }

基本使用：

.. code-block:: rust

   use redis::{AsyncCommands, Client};

   #[tokio::main]
   async fn main() -> redis::RedisResult<()> {
       let client = Client::open("redis://127.0.0.1/")?;
       let mut conn = client.get_multiplexed_async_connection().await?;

       // String 操作
       conn.set("key", "value").await?;
       let value: String = conn.get("key").await?;
       println!("key = {}", value); // key = value

       // 设置过期时间
       conn.set_ex("temp_key", "temp_value", 60).await?; // 60 秒后过期

       // 自增
       conn.set("counter", 0).await?;
       let count: i32 = conn.incr("counter", 1).await?;
       println!("counter = {}", count); // 1

       // Hash 操作
       conn.hset("user:1", "name", "Alice").await?;
       conn.hset("user:1", "age", 28).await?;
       let name: String = conn.hget("user:1", "name").await?;
       let age: i32 = conn.hget("user:1", "age").await?;
       println!("user:1 -> name={}, age={}", name, age);

       // List 操作
       conn.rpush("queue", vec!["task1", "task2", "task3"]).await?;
       let task: String = conn.lpop("queue", None).await?;
       println!("弹出: {}", task); // task1

       // Set 操作
       conn.sadd("tags", "rust").await?;
       conn.sadd("tags", "redis").await?;
       let members: Vec<String> = conn.smembers("tags").await?;
       println!("tags: {:?}", members);

       // Pub/Sub
       // 订阅在单独连接上进行

       // Pipeline
       let (v1, v2): (String, String) = redis::pipe()
           .set("a", "1").ignore()
           .set("b", "2").ignore()
           .get("a")
           .get("b")
           .query_async(&mut conn)
           .await?;
       println!("a={}, b={}", v1, v2);

       Ok(())
   }

连接池（使用 ``bb8`` + ``redis``）：

.. code-block:: toml

   [dependencies]
   redis = { version = "0.25", features = ["tokio-comp"] }
   bb8 = "0.8"
   bb8-redis = "0.15"

常用数据类型操作：

.. list-table:: redis 常用操作
   :header-rows: 1
   :widths: 25 75

   * - 数据类型
     - 常用命令
   * - String
     - ``get`` / ``set`` / ``set_ex`` / ``incr`` / ``decr`` / ``del``
   * - Hash
     - ``hget`` / ``hset`` / ``hgetall`` / ``hdel``
   * - List
     - ``lpush`` / ``rpush`` / ``lpop`` / ``rpop`` / ``lrange``
   * - Set
     - ``sadd`` / ``smembers`` / ``sismember`` / ``srem``
   * - Sorted Set
     - ``zadd`` / ``zrange`` / ``zrangebyscore`` / ``zrem``
   * - Pub/Sub
     - ``publish`` / ``subscribe``

总结
=====

.. list-table:: 数据库 & 持久化 Crate 总览
   :header-rows: 1
   :widths: 20 25 25 30

   * - Crate
     - 类型
     - 数据库
     - 典型场景
   * - ``sqlx``
     - 异步 SQL 工具
     - PostgreSQL / MySQL / SQLite
     - 异步 Web 服务、编译期 SQL 检查
   * - ``diesel``
     - 同步 ORM
     - PostgreSQL / MySQL / SQLite
     - 传统 Web 应用、偏好 ORM 的项目
   * - ``mongodb``
     - NoSQL 驱动
     - MongoDB
     - 文档型存储、灵活 Schema
   * - ``redis``
     - 缓存客户端
     - Redis
     - 缓存、会话存储、消息队列、排行榜
