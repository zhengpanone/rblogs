==========================
安全 & 加密 Crate
==========================

Rust 生态中的安全与加密相关 crate，覆盖密码哈希、JWT 令牌、加密库和 TLS/SSL。

.. contents:: 目录
   :depth: 3
   :local:

ring
==========

高性能加密原语库，纯 Rust（含汇编）实现，零依赖（不含 OpenSSL），被 ``rustls`` 广泛使用。

.. code-block:: toml

   [dependencies]
   ring = "0.17"

哈希摘要（SHA-256）：

.. code-block:: rust

   use ring::digest;

   fn sha256_hash(data: &[u8]) -> String {
       let digest = digest::digest(&digest::SHA256, data);
       // 转为十六进制字符串
       digest.as_ref().iter().map(|b| format!("{:02x}", b)).collect()
   }

   fn main() {
       let hash = sha256_hash(b"hello world");
       println!("SHA-256: {}", hash);
   }

HMAC 消息认证码：

.. code-block:: rust

   use ring::hmac;

   fn main() {
       let key = hmac::Key::new(hmac::HMAC_SHA256, b"my-secret-key");
       let tag = hmac::sign(&key, b"important message");

       // 验证
       hmac::verify(&key, b"important message", tag.as_ref()).unwrap();
       println!("HMAC 验证通过");
   }

随机数生成：

.. code-block:: rust

   use ring::rand;

   fn main() {
       let rng = rand::SystemRandom::new();

       // 生成随机字节
       let mut bytes = [0u8; 32];
       rng.fill(&mut bytes).unwrap();
       println!("随机字节: {:?}", bytes);
   }

PBKDF2 密钥派生：

.. code-block:: rust

   use ring::pbkdf2;

   fn hash_password(password: &str) -> ([u8; 32], [u8; 16]) {
       let mut salt = [0u8; 16];
       ring::rand::SystemRandom::new().fill(&mut salt).unwrap();

       let mut hash = [0u8; 32];
       pbkdf2::derive(
           pbkdf2::PBKDF2_HMAC_SHA256,
           std::num::NonZeroU32::new(100_000).unwrap(),
           &salt,
           password.as_bytes(),
           &mut hash,
       );
       (hash, salt)
   }

   fn verify_password(password: &str, hash: &[u8], salt: &[u8]) -> bool {
       pbkdf2::verify(
           pbkdf2::PBKDF2_HMAC_SHA256,
           std::num::NonZeroU32::new(100_000).unwrap(),
           salt,
           password.as_bytes(),
           hash,
       )
       .is_ok()
   }

常用模块：

.. list-table:: ring 常用模块
   :header-rows: 1
   :widths: 30 70

   * - 模块
     - 说明
   * - ``ring::digest``
     - SHA-256 / SHA-384 / SHA-512 哈希摘要
   * - ``ring::hmac``
     - HMAC 消息认证码（签名与验证）
   * - ``ring::rand``
     - 安全随机数生成
   * - ``ring::pbkdf2``
     - PBKDF2 密钥派生（密码哈希）
   * - ``ring::aead``
     - AES-GCM / ChaCha20-Poly1305 认证加密
   * - ``ring::signature``
     - Ed25519 等数字签名

argon2
==========

Argon2 密码哈希算法，2015 年密码哈希竞赛冠军，抗 GPU/ASIC 暴力破解。

.. code-block:: toml

   [dependencies]
   argon2 = "0.5"

密码哈希与验证：

.. code-block:: rust

   use argon2::{
       password_hash::{rand_core::OsRng, PasswordHash, PasswordHasher, PasswordVerifier, SaltString},
       Argon2,
   };

   fn hash_password(password: &str) -> Result<String, argon2::password_hash::Error> {
       let salt = SaltString::generate(&mut OsRng);
       let argon2 = Argon2::default();
       let hash = argon2.hash_password(password.as_bytes(), &salt)?;
       Ok(hash.to_string())
   }

   fn verify_password(password: &str, hash_str: &str) -> Result<bool, argon2::password_hash::Error> {
       let parsed_hash = PasswordHash::new(hash_str)?;
       let argon2 = Argon2::default();
       Ok(argon2.verify_password(password.as_bytes(), &parsed_hash).is_ok())
   }

   fn main() -> Result<(), argon2::password_hash::Error> {
       let password = "my-secure-password";

       let hash = hash_password(password)?;
       println!("哈希: {}", hash);

       let is_valid = verify_password(password, &hash)?;
       println!("验证: {}", is_valid); // true

       let is_wrong = verify_password("wrong-password", &hash)?;
       println!("错误密码: {}", is_wrong); // false

       Ok(())
   }

Argon2 参数配置：

.. code-block:: rust

   use argon2::{Argon2, Algorithm, Version, Params};

   let argon2 = Argon2::new(
       Algorithm::Argon2id,      // 混合模式，最安全
       Version::V0x13,           // 最新版本
       Params::new(
           65536,                 // 内存开销 (KB)，64 MB
           3,                     // 迭代次数
           4,                     // 并行度
           Some(32),              // 输出长度
       ).unwrap(),
   );

``bcrypt`` vs ``argon2``：

.. list-table:: bcrypt vs argon2
   :header-rows: 1
   :widths: 25 40 35

   * - 特性
     - ``bcrypt``
     - ``argon2``
   * - 标准
     - 1999 年
     - 2015 年（PHC 冠军）
   * - 抗 GPU 攻击
     - 一般
     - 强（内存密集型）
   * - 抗 ASIC 攻击
     - 弱
     - 强
   * - 推荐程度
     - 旧项目兼容
     - 新项目首选

jsonwebtoken
===============

JWT（JSON Web Token）编码、解码和验证。

.. code-block:: toml

   [dependencies]
   jsonwebtoken = "9"
   serde = { version = "1", features = ["derive"] }

签发 JWT：

.. code-block:: rust

   use jsonwebtoken::{encode, decode, Header, Algorithm, Validation, EncodingKey, DecodingKey};
   use serde::{Serialize, Deserialize};
   use std::time::{SystemTime, UNIX_EPOCH};

   #[derive(Debug, Serialize, Deserialize)]
   struct Claims {
       sub: String,       // 用户标识
       exp: usize,        // 过期时间
       iat: usize,        // 签发时间
       role: String,      // 自定义字段
   }

   fn create_token(user_id: &str, secret: &str) -> Result<String, jsonwebtoken::errors::Error> {
       let now = SystemTime::now()
           .duration_since(UNIX_EPOCH)
           .unwrap()
           .as_secs() as usize;

       let claims = Claims {
           sub: user_id.to_string(),
           exp: now + 3600,          // 1 小时后过期
           iat: now,
           role: "user".to_string(),
       };

       encode(
           &Header::new(Algorithm::HS256),
           &claims,
           &EncodingKey::from_secret(secret.as_bytes()),
       )
   }

   fn verify_token(token: &str, secret: &str) -> Result<Claims, jsonwebtoken::errors::Error> {
       let validation = Validation::new(Algorithm::HS256);
       let token_data = decode::<Claims>(
           token,
           &DecodingKey::from_secret(secret.as_bytes()),
           &validation,
       )?;
       Ok(token_data.claims)
   }

   fn main() -> Result<(), jsonwebtoken::errors::Error> {
       let secret = "my-super-secret-key-change-in-production";

       // 签发
       let token = create_token("user-123", secret)?;
       println!("Token: {}", token);

       // 验证
       match verify_token(&token, secret) {
           Ok(claims) => {
               println!("用户: {}, 角色: {}", claims.sub, claims.role);
           }
           Err(e) => {
               println!("Token 无效: {}", e);
           }
       }

       Ok(())
   }

使用 RSA 密钥：

.. code-block:: rust

   // 签发（私钥）
   let encoding_key = EncodingKey::from_rsa_pem(include_bytes!("private.pem"))?;

   // 验证（公钥）
   let decoding_key = DecodingKey::from_rsa_pem(include_bytes!("public.pem"))?;
   let validation = Validation::new(Algorithm::RS256);

常用配置：

.. list-table:: JWT 常用配置
   :header-rows: 1
   :widths: 30 70

   * - 配置
     - 说明
   * - ``Algorithm::HS256``
     - HMAC-SHA256（对称密钥）
   * - ``Algorithm::RS256``
     - RSA-SHA256（非对称密钥）
   * - ``Validation { leeway: 60, .. }``
     - 允许 60 秒的时钟偏差
   * - ``Validation { aud: Some(...), .. }``
     - 验证 audience
   * - ``Validation { iss: Some(...), .. }``
     - 验证 issuer
   * - ``exp``
     - 过期时间（标准 claim）
   * - ``iat``
     - 签发时间（标准 claim）
   * - ``nbf``
     - 生效时间（标准 claim）

openssl
==========

OpenSSL 的 Rust 绑定，提供完整的 SSL/TLS、加密、证书管理功能。

.. code-block:: toml

   [dependencies]
   openssl = "0.10"

对称加密（AES-256-CBC）：

.. code-block:: rust

   use openssl::symm::{encrypt, decrypt, Cipher};

   fn main() {
       let cipher = Cipher::aes_256_cbc();
       let key = b"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F\
                   \x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F";
       let iv = b"\x00\x01\x02\x03\x04\x05\x06\x07\x00\x01\x02\x03\x04\x05\x06\x07";

       let plaintext = b"Hello, World! This is a secret message.";

       // 加密
       let ciphertext = encrypt(cipher, key, Some(iv), plaintext).unwrap();
       println!("密文长度: {} 字节", ciphertext.len());

       // 解密
       let decrypted = decrypt(cipher, key, Some(iv), &ciphertext).unwrap();
       println!("解密: {}", String::from_utf8_lossy(&decrypted));
       assert_eq!(plaintext, &decrypted[..]);
   }

RSA 密钥生成与签名：

.. code-block:: rust

   use openssl::rsa::{Rsa, Padding};
   use openssl::sign::{Signer, Verifier};
   use openssl::hash::MessageDigest;
   use openssl::pkey::PKey;

   fn main() -> Result<(), Box<dyn std::error::Error>> {
       // 生成 RSA 密钥对
       let rsa = Rsa::generate(2048)?;
       let pkey = PKey::from_rsa(rsa)?;

       // 签名
       let mut signer = Signer::new(MessageDigest::sha256(), &pkey)?;
       signer.update(b"hello world")?;
       let signature = signer.sign_to_vec()?;
       println!("签名长度: {} 字节", signature.len());

       // 验证
       let mut verifier = Verifier::new(MessageDigest::sha256(), &pkey)?;
       verifier.update(b"hello world")?;
       let valid = verifier.verify(&signature)?;
       println!("签名验证: {}", valid); // true

       Ok(())
   }

X.509 证书解析：

.. code-block:: rust

   use openssl::x509::X509;

   fn main() -> Result<(), Box<dyn std::error::Error>> {
       let cert_pem = std::fs::read("cert.pem")?;
       let cert = X509::from_pem(&cert_pem)?;

       println!("主题: {}", cert.subject_name());
       println!("颁发者: {}", cert.issuer_name());
       println!("有效期: {} - {}", cert.not_before(), cert.not_after());

       Ok(())
   }

常用模块：

.. list-table:: openssl 常用模块
   :header-rows: 1
   :widths: 25 75

   * - 模块
     - 说明
   * - ``openssl::symm``
     - 对称加密（AES、DES 等）
   * - ``openssl::rsa``
     - RSA 密钥生成、加解密、签名
   * - ``openssl::sign``
     - 签名与验证
   * - ``openssl::hash``
     - 消息摘要（SHA、MD5 等）
   * - ``openssl::x509``
     - X.509 证书处理
   * - ``openssl::pkcs12``
     - PKCS#12 证书包处理
   * - ``openssl::ssl``
     - SSL/TLS 连接
   * - ``openssl::pkey``
     - 通用密钥封装

ring vs openssl：

.. list-table:: ring vs openssl
   :header-rows: 1
   :widths: 25 40 35

   * - 特性
     - ``ring``
     - ``openssl``
   * - 依赖
     - 零系统依赖
     - 需要系统安装 OpenSSL
   * - 编译
     - 简单
     - 可能遇到版本/路径问题
   * - 功能范围
     - 精简（核心加密原语）
     - 全面（TLS、证书、多种算法）
   * - 典型用途
     - rustls 底层、轻量加密
     - 传统 TLS、证书管理、全功能加密
   * - 推荐
     - 新项目首选
     - 需要 OpenSSL 特性时

总结
=====

.. list-table:: 安全 & 加密 Crate 总览
   :header-rows: 1
   :widths: 20 25 25 30

   * - Crate
     - 类型
     - 核心能力
     - 典型场景
   * - ``ring``
     - 加密原语库
     - SHA-256、HMAC、随机数、PBKDF2、AEAD
     - 轻量加密、rustls 底层
   * - ``argon2``
     - 密码哈希
     - Argon2id 哈希与验证
     - 用户密码存储、认证
   * - ``jsonwebtoken``
     - JWT 令牌
     - HS256 / RS256 签发与验证
     - API 认证、SSO、无状态会话
   * - ``openssl``
     - 全功能加密
     - AES、RSA、签名、证书、TLS
     - 传统加密、证书管理、兼容性需求
