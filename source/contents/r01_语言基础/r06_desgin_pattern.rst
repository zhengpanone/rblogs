======================
设计模式
======================

Newtype 模式
==================

Newtype(新类型模式)是Rust中的一种设计模式。它通过创建一个新的类型来包装一个现有的类型，从而提供更强的类型安全和更清晰的代码表达。

  用一个元组结构体（Tuple Struct）包装已有类型，从而创建一个全新的类型。

.. code-block:: rust

  struct UserId(i64);
  struct OrderId(i64);

  fn main() {
    let id = UserId(1);
    let order_id = OrderId(2);

    // 编译错误：类型不匹配
    // user_id = order_id;
  }


虽然底层是 i64，但：

- 它 不等于 i64

- 不能和 OrderId 混用

- 是一个独立类型
  

NewType它可以解决：

- 原始类型污染（String 到处飞）

- ID 混用（UserId 和 OrderId 都是 i64）

- 无法在类型层表达领域语义

- 违反 DDD 的“显式建模”


这就是 DDD 想要的：领域语义通过类型表达

Newtype 应该主要出现在：

- Value Object

- Entity Id

- 强语义字段（Email / Password / Username）

Newtype 的使用场景
---------------------------

隐藏内部实现
>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

  pub struct Password(String);

  impl Password {
    pub fn verify(&self, input: &str) -> bool {
      self.0 == input
    }
  }

绕过孤儿规则(Orphan Rule)
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

  impl Display for Vec<String>{}
  // error: cannot implement foreign trait

因为 Display、Vec都不是自己的。

利用NewType 可以绕过孤儿规则。

.. code-block:: rust

  struct MyVec(Vec<String>);

  impl Display for MyVec {
    fn fmt(
        &self,
        f: &mut std::fmt::Formatter<'_>,
    ) -> std::fmt::Result {

        write!(f, "{}", self.0.join(","))
    }
  }

DDD 中的Newtype
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

  #[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize,Deserialize)]
  pub struct UserId(i64);

  #[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
  pub struct OrderId(i64);

  pub struct Phone(String);

  pub struct Email(String);

带校验的Newtype

.. code-block:: rust

  pub struct Email(String);

  impl Email {
      pub fn new(value: String) -> Result<Self, String> {
          if value.contains('@') {
              Ok(Self(value))
          } else {
              Err("invalid email".into())
          }
      }
  }

  fn main(){
      // 创建成功后一定合法。
      let email = Email::new("test@example.com".to_string());
      println!("{:?}", email);
  }

数据库中的Newtype
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

  #[derive(sqlx::Type,Debug)]
  #[sqlx(transparent)]
  struct UserId(i64);


Path 参数转换
>>>>>>>>>>>>>>>>>>>>>>>>>

.. code-block:: rust

  use axum::{extract::Path};
  use std::str::FromStr;

  async fn get_user(Path(id): Path<UserId>) {
  }

  impl FromStr for UserId {
    type Err = std::num::ParseIntError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(UserId(s.parse()?))
    }
  }

零成本抽象(Zero-Cost Abstraction)
------------------------------------

Newtype不增加运行时开销

.. code-block:: rust

  struct UserId(u64);

编译后

.. code-block:: text

  UserId
  ≈
  u64

大小一致

.. code-block:: rust

  use std::mem;

  println!(
      "{}",
      mem::size_of::<UserId>()
  );
  // 8

项目最佳实践
--------------------


.. code-block:: rust

  #[derive(
    Debug,
    Clone,
    Copy,
    PartialEq,
    Eq,
    Hash,
    PartialOrd,
    Ord
  )]
  pub struct UserId(u64);

  impl From<u64> for UserId {}
  impl From<UserId> for u64 {}

对于字符串

.. code-block:: rust

  #[derive(
    Debug,
    Clone,
    PartialEq,
    Eq,
    Hash
  )]
  pub struct Email(String);

  // 进行合法性校验
  impl TryFrom<String> for Email

NewType 设计模板
--------------------

.. code-block:: rust

  use serde::{Deserialize, Serialize};

  #[derive(
      Debug,
      Clone,
      Copy,
      PartialEq,
      Eq,
      Hash,
      PartialOrd,
      Ord,
      Serialize,
      Deserialize
  )]
  pub struct UserId(u64);

  impl UserId {
      pub fn new(id: u64) -> Self {
          Self(id)
      }

      pub fn value(self) -> u64 {
          self.0
      }
  }

  impl From<u64> for UserId {
      fn from(value: u64) -> Self {
          Self(value)
      }
  }

  impl From<UserId> for u64 {
      fn from(id: UserId) -> Self {
          id.0
      }
  }

适用于:

- DDD（领域驱动设计）
- Web API DTO
- 数据库实体
- 微服务 ID 类型
- 金额类型
- 邮箱/手机号类型
- OAuth ClientId
- Flowable ProcessInstanceId 封装
- Kubernetes ResourceName 封装

