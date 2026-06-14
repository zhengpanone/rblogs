密码学与哈希 (Cryptography & Hashing)
=====================================

.. .. module:: r05_security_programming.r01_cryptography_hashing

Rust 密码学生态涵盖哈希算法、对称/非对称加密、数字签名和密钥派生等核心能力。
本章介绍最常用的密码学 Crate，帮助你在项目中安全地实现数据加密和完整性保护。

ring
----

``ring`` 是 Rust 密码学生态的基石之一，由 BoringSSL 核心代码移植而来，提供最基础的密码学原语。
其设计哲学是"少即是多"——只暴露经过安全审计的、难以误用的 API。

**哈希计算：**

.. code-block:: rust

    use ring::digest::{digest, SHA256, SHA512};

    let data = b"hello, world";
    let hash = digest(&SHA256, data);
    println!("SHA256: {:x?}", hash.as_ref());

    let hash = digest(&SHA512, data);
    println!("SHA512: {:x?}", hash.as_ref());

**HMAC 消息认证码：**

.. code-block:: rust

    use ring::hmac;

    let key = hmac::Key::new(hmac::HMAC_SHA256, b"my-secret-key");
    let tag = hmac::sign(&key, b"authenticated message");
    hmac::verify(&key, b"authenticated message", tag.as_ref()).unwrap();

**PBKDF2 密钥派生：**

.. code-block:: rust

    use ring::pbkdf2;

    const CREDENTIAL_LEN: usize = digest::SHA256_OUTPUT_LEN;
    let mut to_store = [0u8; CREDENTIAL_LEN];

    pbkdf2::derive(
        pbkdf2::PBKDF2_HMAC_SHA256,
        std::num::NonZeroU32::new(100_000).unwrap(),
        b"salt",
        b"user-password",
        &mut to_store,
    );

**AEAD 认证加密 (AES-256-GCM)：**

.. code-block:: rust

    use ring::aead::{Aad, BoundKey, LessSafeKey, Nonce, NonceSequence, UnboundKey, AES_256_GCM};

    struct CounterNonce(u32);

    impl NonceSequence for CounterNonce {
        fn advance(&mut self) -> Result<Nonce, ring::error::Unspecified> {
            let mut nonce_bytes = vec![0u8; 12];
            nonce_bytes[..4].copy_from_slice(&self.0.to_be_bytes());
            self.0 += 1;
            Nonce::try_assume_unique_for_key(&nonce_bytes)
        }
    }

    let unbound_key = UnboundKey::new(&AES_256_GCM, &[0u8; 32]).unwrap();
    let mut key = LessSafeKey::new(unbound_key);

    let nonce = CounterNonce(1);
    let aad = Aad::from(b"additional data");
    let mut in_out = b"secret message".to_vec();

    key.seal_in_place_append_tag(nonce, aad, &mut in_out).unwrap();
    // in_out now contains ciphertext + 16-byte tag

    let nonce = CounterNonce(1);
    let plaintext = key.open_in_place(nonce, aad, &mut in_out).unwrap();
    // plaintext == b"secret message"

**随机数生成：**

.. code-block:: rust

    use ring::rand::{SecureRandom, SystemRandom};

    let rng = SystemRandom::new();
    let mut key = [0u8; 32];
    rng.fill(&mut key).unwrap();

.. list-table:: ring 常用模块
   :header-rows: 1

   * - 模块
     - 用途
   * - ``digest``
     - SHA-256 / SHA-512 / SHA-1 哈希
   * - ``hmac``
     - HMAC 消息认证码
   * - ``pbkdf2``
     - PBKDF2 密钥派生
   * - ``aead``
     - AES-GCM / ChaCha20-Poly1305 认证加密
   * - ``rand``
     - 安全随机数生成
   * - ``signature``
     - Ed25519 / ECDSA 数字签名
   * - ``agreement``
     - ECDH 密钥协商

.. note::

   ``ring`` 的 AEAD API 分为 ``LessSafeKey``（单一密钥+外部 Nonce 序列）和
   ``SealingKey`` / ``OpeningKey``（编译期保证 Nonce 唯一性）。
   生产代码中应优先使用 ``SealingKey`` / ``OpeningKey``。


rust-crypto (RustCrypto) 生态
--------------------------------

`RustCrypto <https://github.com/RustCrypto>`_ 是一个模块化的纯 Rust 密码学项目，
每个算法对应一个独立 Crate，按需引入。

**SHA-2 / SHA-3 哈希：**

.. code-block:: rust

    use sha2::{Sha256, Sha512, Digest};

    let mut hasher = Sha256::new();
    hasher.update(b"hello ");
    hasher.update(b"world");
    let result = hasher.finalize();
    println!("SHA256: {:x}", result);

**HMAC：**

.. code-block:: rust

    use hmac::{Hmac, Mac};
    use sha2::Sha256;

    type HmacSha256 = Hmac<Sha256>;

    let mut mac = HmacSha256::new_from_slice(b"secret-key").unwrap();
    mac.update(b"message");
    let code = mac.finalize();
    let code_bytes = code.into_bytes();

**AES-GCM 加密：**

.. code-block:: rust

    use aes_gcm::{Aes256Gcm, KeyInit, Nonce};
    use aes_gcm::aead::{Aead, OsRng};

    let key = Aes256Gcm::generate_key(&mut OsRng);
    let cipher = Aes256Gcm::new(&key);
    let nonce = Nonce::from_slice(b"unique-12bytes"); // 96-bit (12 bytes)

    let ciphertext = cipher.encrypt(nonce, b"plaintext message".as_ref()).unwrap();
    let plaintext = cipher.decrypt(nonce, ciphertext.as_ref()).unwrap();
    assert_eq!(plaintext, b"plaintext message");

**Ed25519 数字签名：**

.. code-block:: rust

    use ed25519_dalek::{SigningKey, VerifyingKey, Signature, Signer, Verifier};
    use rand::rngs::OsRng;

    let signing_key = SigningKey::generate(&mut OsRng);
    let verifying_key = VerifyingKey::from(&signing_key);

    let signature: Signature = signing_key.sign(b"important document");
    verifying_key.verify(b"important document", &signature).unwrap();

**Argon2 密码哈希：**

.. code-block:: rust

    use argon2::{Argon2, PasswordHash, PasswordHasher, PasswordVerifier};
    use argon2::password_hash::{SaltString, rand_core::OsRng};

    let password = b"hunter2";
    let salt = SaltString::generate(&mut OsRng);

    let argon2 = Argon2::default();
    let hash = argon2.hash_password(password, &salt).unwrap().to_string();

    let parsed_hash = PasswordHash::new(&hash).unwrap();
    argon2.verify_password(password, &parsed_hash).unwrap();

.. list-table:: RustCrypto 常用 Crate
   :header-rows: 1

   * - Crate
     - 功能
     - 类型
   * - ``sha2``
     - SHA-224 / SHA-256 / SHA-384 / SHA-512
     - 哈希
   * - ``sha3``
     - SHA3-224 / SHA3-256 / SHA3-384 / SHA3-512 / SHAKE
     - 哈希
   * - ``md-5``
     - MD5（仅用于遗留兼容，不可用于安全）
     - 哈希
   * - ``hmac``
     - HMAC 消息认证码
     - MAC
   * - ``aes-gcm``
     - AES-GCM 认证加密
     - AEAD
   * - ``chacha20poly1305``
     - ChaCha20-Poly1305 认证加密
     - AEAD
   * - ``ed25519-dalek``
     - Ed25519 数字签名
     - 签名
   * - ``rsa``
     - RSA 加密/签名
     - 非对称
   * - ``argon2``
     - Argon2 密码哈希（推荐）
     - 密钥派生
   * - ``bcrypt``
     - bcrypt 密码哈希
     - 密钥派生
   * - ``x25519-dalek``
     - X25519 密钥协商
     - 密钥协商

ring vs RustCrypto 对比
------------------------

.. list-table:: ring vs RustCrypto 选择指南
   :header-rows: 1

   * - 维度
     - ring
     - RustCrypto
   * - 实现语言
     - C + 汇编（BoringSSL 子集）
     - 纯 Rust
   * - 算法覆盖
     - 精选核心算法，API 简洁
     - 算法丰富，模块化
   * - 性能
     - 极高（手写汇编优化）
     - 良好
   * - FIPS 认证
     - 部分支持
     - 不直接支持
   * - 编译速度
     - 较快（单 Crate）
     - 较慢（多 Crate 依赖链）
   * - 学习曲线
     - API 独特，有学习成本
     - 标准 trait 设计，易上手
   * - 适合场景
     - 只需核心原语、追求极致性能
     - 需要丰富算法、纯 Rust 栈

orion
-----

``orion`` 是一个纯 Rust 密码学库，强调可用性和默认安全。

.. code-block:: rust

    use orion::{
        hash::{digest, sha512::Sha512},
        auth::Authenticator,
        aead,
        pwhash::{self, PasswordHash},
    };

    // 哈希
    let hash = digest(&Sha512, b"hello").unwrap();

    // HMAC
    let key = Authenticator::generate_key().unwrap();
    let tag = Authenticator::authenticate(&key, b"message").unwrap();
    Authenticator::authenticate_verify(&tag, &key, b"message").unwrap();

    // AEAD (XChaCha20-Poly1305)
    let key = aead::SecretKey::generate().unwrap();
    let ciphertext = aead::seal(&key, b"secret data").unwrap();
    let plaintext = aead::open(&key, &ciphertext).unwrap();

    // 密码哈希（Argon2i）
    let hash = PasswordHash::hash_password(b"user-password", 1 << 16).unwrap();
    hash.verify(b"user-password").unwrap();

.. note::

   ``orion`` 不暴露任何 ``unsafe`` 接口，所有 API 默认使用安全的参数配置，
   非常适合对密码学不熟悉但需要安全实现的开发者。

password-hash 与密码存储最佳实践
------------------------------------

.. code-block:: rust

    use argon2::{
        Argon2, Algorithm, Params, Version,
        password_hash::{rand_core::OsRng, PasswordHash, PasswordHasher, SaltString},
    };

    /// 推荐参数：Argon2id, 64MB 内存, 3 次迭代, 1 并行度
    fn hash_password(password: &str) -> Result<String, argon2::password_hash::Error> {
        let salt = SaltString::generate(&mut OsRng);
        let params = Params::new(65536, 3, 1, None).unwrap();
        let argon2 = Argon2::new(Algorithm::Argon2id, Version::V0x13, params);

        let hash = argon2.hash_password(password.as_bytes(), &salt)?;
        Ok(hash.to_string())
    }

    fn verify_password(password: &str, hash_str: &str) -> Result<bool, argon2::password_hash::Error> {
        let parsed_hash = PasswordHash::new(hash_str)?;
        Ok(parsed_hash.hash_password
            .as_ref()
            .map(|_| Argon2::default().verify_password(password.as_bytes(), &parsed_hash).is_ok())
            .unwrap_or(false))
    }

.. list-table:: 密码存储最佳实践
   :header-rows: 1

   * - 实践
     - 说明
   * - 使用 Argon2id
     - 2024 年 OWASP 推荐首选，抗 GPU/ASIC/侧信道
   * - 最少 64MB 内存
     - 增加暴力破解成本
   * - 每次生成独立随机盐
     - 防止彩虹表攻击，防止相同密码产生相同哈希
   * - 使用恒定时间比较
     - ``argon2`` Crate 内置，无需额外处理
   * - 绝不自行实现
     - 使用经过安全审计的成熟 Crate

常量时间比较
---------------

.. code-block:: rust

    use subtle::ConstantTimeEq;

    let a: &[u8] = b"secret-value";
    let b: &[u8] = b"secret-value";
    let c: &[u8] = b"wrong-value";

    // 恒定时间比较，防止时序攻击
    if a.ct_eq(b).into() {
        println!("equal");
    }

    // 避免：标准 == 比较可能泄露时序信息
    // if a == b { ... }  // 不推荐用于安全敏感场景

总结
-----

.. list-table:: 密码学与哈希 Crate 总览
   :header-rows: 1

   * - Crate
     - 定位
     - 选择场景
   * - ``ring``
     - 核心原语，C + 汇编优化
     - 极致性能、TLS 底层、嵌入式
   * - ``sha2`` / ``aes-gcm`` / ``ed25519-dalek`` ...
     - RustCrypto 生态，纯 Rust 模块化
     - 需要丰富算法、纯 Rust 栈
   * - ``orion``
     - 纯 Rust，默认安全
     - 追求易用性、零 unsafe
   * - ``argon2``
     - 密码哈希首选
     - 用户密码存储
   * - ``subtle``
     - 恒定时间操作
     - 防止时序攻击
