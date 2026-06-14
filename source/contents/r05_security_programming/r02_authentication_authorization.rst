认证与授权 (Authentication & Authorization)
=============================================

.. .. module:: r05_security_programming.r02_authentication_authorization

认证（验证"你是谁"）和授权（验证"你能做什么"）是 Web 安全的核心支柱。
本章介绍 Rust 生态中 JWT、OAuth2、会话管理等认证授权相关 Crate。

jsonwebtoken (JWT)
-------------------

``jsonwebtoken`` 是 Rust 中最常用的 JWT 库，支持 HS256 / RS256 / ES256 等算法。

**签发 JWT：**

.. code-block:: rust

    use jsonwebtoken::{encode, decode, Header, Algorithm, Validation, EncodingKey, DecodingKey};
    use serde::{Deserialize, Serialize};
    use chrono::{Utc, Duration};

    #[derive(Debug, Serialize, Deserialize)]
    struct Claims {
        sub: String,          // 用户 ID
        role: String,         // 角色
        exp: usize,           // 过期时间 (Unix timestamp)
        iat: usize,           // 签发时间
    }

    fn create_jwt(user_id: &str, role: &str, secret: &[u8]) -> Result<String, jsonwebtoken::errors::Error> {
        let now = Utc::now();
        let claims = Claims {
            sub: user_id.to_string(),
            role: role.to_string(),
            iat: now.timestamp() as usize,
            exp: (now + Duration::hours(24)).timestamp() as usize,
        };

        encode(&Header::default(), &claims, &EncodingKey::from_secret(secret))
    }

**验证 JWT：**

.. code-block:: rust

    fn validate_jwt(token: &str, secret: &[u8]) -> Result<Claims, jsonwebtoken::errors::Error> {
        let validation = Validation::new(Algorithm::HS256);
        let token_data = decode::<Claims>(
            token,
            &DecodingKey::from_secret(secret),
            &validation,
        )?;
        Ok(token_data.claims)
    }

**RS256 非对称签名：**

.. code-block:: rust

    use jsonwebtoken::EncodingKey;

    // 签发端：使用私钥
    let private_key = include_bytes!("../keys/private.pem");
    let encoding_key = EncodingKey::from_rsa_pem(private_key).unwrap();
    let token = encode(&Header::new(Algorithm::RS256), &claims, &encoding_key).unwrap();

    // 验证端：使用公钥
    let public_key = include_bytes!("../keys/public.pem");
    let decoding_key = DecodingKey::from_rsa_pem(public_key).unwrap();
    let token_data = decode::<Claims>(&token, &decoding_key, &Validation::new(Algorithm::RS256)).unwrap();

.. list-table:: jsonwebtoken 常用功能
   :header-rows: 1

   * - 功能
     - 说明
   * - ``Header``
     - JWT 头部，指定算法和类型
   * - ``Validation``
     - 验证规则：算法、过期时间 (exp)、生效时间 (nbf)、签发者 (iss)、受众 (aud)
   * - ``EncodingKey``
     - 签名密钥（HMAC 对称密钥或 RSA/EC 私钥）
   * - ``DecodingKey``
     - 验证密钥（对称密钥或公钥）
   * - ``Algorithm::HS256/384/512``
     - HMAC 对称签名算法
   * - ``Algorithm::RS256/384/512``
     - RSA 非对称签名算法
   * - ``Algorithm::ES256/384``
     - ECDSA 非对称签名算法

JWT 最佳实践
^^^^^^^^^^^^^^^^

.. list-table:: JWT 安全实践
   :header-rows: 1

   * - 实践
     - 说明
   * - 设置合理过期时间
     - Access Token 15-60 分钟，Refresh Token 7-30 天
   * - 使用 RS256 而非 HS256
     - 非对称签名更安全，私钥仅在认证服务持有
   * - 验证所有声明
     - 始终验证 ``exp`` / ``nbf`` / ``iss`` / ``aud``
   * - 不在 JWT 中存放敏感数据
     - Payload 仅 Base64 编码，非加密
   * - 使用 Token 黑名单
     - 登出时加入黑名单，防止已签发 Token 继续使用

oauth2
------

``oauth2`` Crate 提供 OAuth2 授权码流程的客户端实现，支持所有标准 OAuth2 提供商。

**授权码流程 (Authorization Code Grant)：**

.. code-block:: rust

    use oauth2::{
        AuthUrl, ClientId, ClientSecret, CsrfToken, PkceCodeChallenge,
        RedirectUrl, Scope, TokenResponse, TokenUrl,
        basic::BasicClient,
        reqwest::async_http_client,
    };
    use url::Url;

    async fn oauth2_login() -> Result<(), Box<dyn std::error::Error>> {
        let client = BasicClient::new(
            ClientId::new("your-client-id".to_string()),
            Some(ClientSecret::new("your-client-secret".to_string())),
            AuthUrl::new("https://provider.com/authorize".to_string())?,
            Some(TokenUrl::new("https://provider.com/token".to_string())?),
        )
        .set_redirect_uri(RedirectUrl::new("https://yourapp.com/callback".to_string())?);

        // Step 1: 生成授权 URL + PKCE 验证器
        let (pkce_challenge, pkce_verifier) = PkceCodeChallenge::new_random_sha256();
        let (auth_url, csrf_token) = client
            .authorize_url(CsrfToken::new_random)
            .add_scope(Scope::new("read".to_string()))
            .add_scope(Scope::new("write".to_string()))
            .set_pkce_challenge(pkce_challenge)
            .url();

        // 将 auth_url 重定向给用户...
        println!("Visit: {}", auth_url);

        // Step 2: 用户授权后回调，用 authorization code 换 token
        // let code = AuthorizationCode::new("code-from-callback".to_string());
        // let token_result = client
        //     .exchange_code(code)
        //     .set_pkce_verifier(pkce_verifier)
        //     .request_async(async_http_client)
        //     .await?;

        // println!("Access Token: {}", token_result.access_token().secret());

        Ok(())
    }

.. list-table:: OAuth2 授权流程
   :header-rows: 1

   * - 流程
     - 适用场景
     - PKCE
   * - Authorization Code
     - Web 应用（有后端）
     - 强烈推荐
   * - Authorization Code + PKCE
     - SPA / 移动端 / CLI
     - 必须
   * - Client Credentials
     - 服务间通信
     - 不适用
   * - Device Code
     - 输入受限设备 (TV/CLI)
     - 不适用

openidconnect
---------------

``openidconnect`` 基于 ``oauth2`` Crate，增加了 OpenID Connect (OIDC) 身份层支持。

.. code-block:: rust

    use openidconnect::{
        core::{
            CoreClient, CoreIdToken, CoreIdTokenVerifier, CoreProviderMetadata,
            CoreResponseType, CoreUserInfoClaims,
        },
        reqwest::async_http_client,
        AuthenticationFlow, ClientId, ClientSecret, CsrfToken, IssuerUrl,
        Nonce, RedirectUrl, Scope,
    };

    async fn oidc_login() -> Result<(), Box<dyn std::error::Error>> {
        let provider_metadata = CoreProviderMetadata::discover_async(
            IssuerUrl::new("https://accounts.google.com".to_string())?,
            async_http_client,
        )
        .await?;

        let client = CoreClient::from_provider_metadata(
            provider_metadata,
            ClientId::new("client-id".to_string()),
            Some(ClientSecret::new("client-secret".to_string())),
        )
        .set_redirect_uri(RedirectUrl::new("https://app.com/callback".to_string())?);

        let (auth_url, csrf_token, nonce) = client
            .authorize_url(
                AuthenticationFlow::<CoreResponseType>::AuthorizationCode,
                CsrfToken::new_random,
                Nonce::new_random,
            )
            .add_scope(Scope::new("email".to_string()))
            .add_scope(Scope::new("profile".to_string()))
            .url();

        println!("Visit: {}", auth_url);
        Ok(())
    }

.. note::

   ``openidconnect`` 提供了 ``IdToken`` 验证器，自动处理 JWT 签名验证、
   ``nonce`` 校验、``iss`` / ``aud`` 声明检查等安全细节。
   使用 OIDC 时务必启用 ``IdToken`` 验证。

tower-sessions
--------------

``tower-sessions`` 是 Axum 生态的会话管理中间件，支持多种存储后端。

.. code-block:: rust

    use axum::{routing::get, Router};
    use tower_sessions::{Session, SessionManagerLayer};
    use tower_sessions::cookie::time::Duration;
    use tower_sessions_sqlx_store::SqliteStore;
    use sqlx::sqlite::SqlitePoolOptions;

    async fn login(session: Session) -> String {
        session.insert("user_id", "42").await.unwrap();
        session.insert("role", "admin").await.unwrap();
        "logged in".to_string()
    }

    async fn dashboard(session: Session) -> String {
        let user_id: Option<String> = session.get("user_id").await.unwrap();
        match user_id {
            Some(id) => format!("Welcome, user {}", id),
            None => "Please login first".to_string(),
        }
    }

    async fn logout(session: Session) -> String {
        session.flush().await.unwrap();
        "logged out".to_string()
    }

    #[tokio::main]
    async fn main() {
        let pool = SqlitePoolOptions::new()
            .connect("sqlite:data.db")
            .await
            .unwrap();

        let session_store = SqliteStore::new(pool);
        let session_layer = SessionManagerLayer::new(session_store)
            .with_secure(false)  // 开发环境关闭 HTTPS Only
            .with_max_age(Duration::hours(2));

        let app = Router::new()
            .route("/login", get(login))
            .route("/dashboard", get(dashboard))
            .route("/logout", get(logout))
            .layer(session_layer);

        let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
        axum::serve(listener, app).await.unwrap();
    }

.. list-table:: tower-sessions 存储后端
   :header-rows: 1

   * - Crate
     - 后端
   * - ``tower-sessions-sqlx-store``
     - PostgreSQL / SQLite / MySQL (via sqlx)
   * - ``tower-sessions-redis-store``
     - Redis
   * - ``tower-sessions-moka-store``
     - 内存 (Moka 缓存)
   * - ``tower-sessions-mongodb-store``
     - MongoDB

paseto
------

``paseto`` (Platform-Agnostic SEcurity TOkens) 是 JWT 的现代化替代方案，
避免了 JWT 的算法混淆攻击风险。

.. code-block:: rust

    use paseto::{
        keys::SymmetricKey,
        tokens::{PasetoBuilder, TimeBackend},
        builder::Builder,
    };
    use serde::{Serialize, Deserialize};

    #[derive(Debug, Serialize, Deserialize)]
    struct UserClaims {
        sub: String,
        role: String,
    }

    fn create_paseto_token(key: &[u8]) -> String {
        let symmetric_key = SymmetricKey::<32>::from(key.try_into().unwrap());
        let claims = UserClaims {
            sub: "user-1".to_string(),
            role: "admin".to_string(),
        };

        PasetoBuilder::<SymmetricKey<32>>::default()
            .set_claim(claims)
            .set_expiration(&TimeBackend::now().add_hours(2).unwrap())
            .build(&symmetric_key)
            .unwrap()
    }

.. note::

   PASETO 提供 ``v1.local`` (对称加密+认证)、``v1.public`` (非对称签名)、
   ``v2.local``、``v2.public`` 等多种模式。与 JWT 不同，PASETO 不允许
   算法选择——每个版本只允许一种算法，从根本上杜绝了算法混淆攻击。

Argon2 密码认证流程
----------------------

.. code-block:: rust

    use argon2::{
        Argon2, PasswordHash, PasswordHasher, PasswordVerifier,
        password_hash::{rand_core::OsRng, SaltString, Error},
    };

    struct AuthService;

    impl AuthService {
        fn register(password: &str) -> Result<String, Error> {
            let salt = SaltString::generate(&mut OsRng);
            let argon2 = Argon2::default();
            let hash = argon2.hash_password(password.as_bytes(), &salt)?;
            Ok(hash.to_string())
        }

        fn login(password: &str, stored_hash: &str) -> Result<bool, Error> {
            let parsed_hash = PasswordHash::new(stored_hash)?;
            Ok(Argon2::default()
                .verify_password(password.as_bytes(), &parsed_hash)
                .is_ok())
        }
    }

总结
-----

.. list-table:: 认证与授权 Crate 总览
   :header-rows: 1

   * - Crate
     - 用途
     - 适用场景
   * - ``jsonwebtoken``
     - JWT 签发与验证
     - API 认证、微服务通信
   * - ``oauth2``
     - OAuth2 客户端
     - 第三方登录 (Google/GitHub/...)
   * - ``openidconnect``
     - OIDC 身份认证
     - SSO 单点登录、企业级身份
   * - ``tower-sessions``
     - 服务端会话管理
     - Axum Web 应用会话
   * - ``paseto``
     - 安全 Token (JWT 替代)
     - 需要更高安全保障的 Token 场景
   * - ``argon2``
     - 密码哈希与验证
     - 用户密码存储
