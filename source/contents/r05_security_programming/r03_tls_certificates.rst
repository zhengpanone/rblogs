TLS 与证书管理 (TLS & Certificate Management)
==============================================

.. .. module:: r05_security_programming.r03_tls_certificates

TLS (Transport Layer Security) 是互联网通信安全的基石。
本章介绍 Rust 中 TLS 实现、证书管理与 HTTPS 配置。

rustls
------

``rustls`` 是纯 Rust 实现的 TLS 库，零 unsafe 依赖 ring 或 aws-lc-rs，
是 Rust 生态的标准 TLS 方案。

**TLS 客户端：**

.. code-block:: rust

    use rustls::pki_types::{ServerName, CertificateDer};
    use std::sync::Arc;
    use tokio::net::TcpStream;
    use tokio_rustls::{TlsConnector, rustls::ClientConfig};

    async fn tls_client(host: &str) -> Result<(), Box<dyn std::error::Error>> {
        // 加载系统根证书
        let mut root_store = rustls::RootCertStore::empty();
        for cert in rustls_native_certs::load_native_certs()? {
            root_store.add(cert).unwrap();
        }

        let config = ClientConfig::builder()
            .with_root_certificates(root_store)
            .with_no_client_auth();

        let connector = TlsConnector::from(Arc::new(config));
        let stream = TcpStream::connect(format!("{}:443", host)).await?;
        let domain = ServerName::try_from(host.to_string())?;

        let mut tls_stream = connector.connect(domain, stream).await?;
        // 使用 tls_stream 进行加密通信...

        Ok(())
    }

**TLS 服务端（自签名证书示例）：**

.. code-block:: rust

    use rustls::ServerConfig;
    use rustls_pemfile::{certs, pkcs8_private_keys};
    use std::sync::Arc;
    use tokio::net::TcpListener;
    use tokio_rustls::TlsAcceptor;
    use tokio::io::{AsyncReadExt, AsyncWriteExt};

    fn load_tls_config() -> Result<ServerConfig, Box<dyn std::error::Error>> {
        let cert_file = &mut std::io::BufReader::new(
            std::fs::File::open("cert.pem")?
        );
        let key_file = &mut std::io::BufReader::new(
            std::fs::File::open("key.pem")?
        );

        let cert_chain: Vec<rustls::pki_types::CertificateDer> = certs(cert_file)
            .filter_map(|c| c.ok())
            .collect();
        let key = pkcs8_private_keys(key_file)
            .filter_map(|k| k.ok())
            .next()
            .ok_or("no private key found")?;

        let config = ServerConfig::builder()
            .with_no_client_auth()
            .with_single_cert(cert_chain, rustls::pki_types::PrivateKeyDer::Pkcs8(key))?;

        Ok(config)
    }

    async fn tls_server() -> Result<(), Box<dyn std::error::Error>> {
        let config = load_tls_config()?;
        let acceptor = TlsAcceptor::from(Arc::new(config));
        let listener = TcpListener::bind("0.0.0.0:8443").await?;

        loop {
            let (stream, addr) = listener.accept().await?;
            let acceptor = acceptor.clone();

            tokio::spawn(async move {
                match acceptor.accept(stream).await {
                    Ok(mut tls_stream) => {
                        let mut buf = [0u8; 1024];
                        let n = tls_stream.read(&mut buf).await.unwrap();
                        let response = b"HTTP/1.1 200 OK\r\n\r\nHello TLS!";
                        tls_stream.write_all(response).await.unwrap();
                    }
                    Err(e) => eprintln!("TLS error: {}", e),
                }
            });
        }
    }

**与 Axum 集成：**

.. code-block:: rust

    use axum::{routing::get, Router};
    use axum_server::tls_rustls::RustlsConfig;

    async fn handler() -> &'static str {
        "Hello, HTTPS!"
    }

    #[tokio::main]
    async fn main() {
        let config = RustlsConfig::from_pem_file("cert.pem", "key.pem")
            .await
            .unwrap();

        let app = Router::new().route("/", get(handler));

        axum_server::bind_rustls("0.0.0.0:443".parse().unwrap(), config)
            .serve(app.into_make_service())
            .await
            .unwrap();
    }

.. list-table:: rustls 常用功能
   :header-rows: 1

   * - 功能
     - 说明
   * - ``ClientConfig``
     - 客户端 TLS 配置（根证书、客户端认证）
   * - ``ServerConfig``
     - 服务端 TLS 配置（证书链、私钥、客户端认证）
   * - ``RootCertStore``
     - 信任的根证书存储
   * - ``rustls-native-certs``
     - 加载操作系统原生根证书
   * - ``rustls-pemfile``
     - PEM 文件解析
   * - ``tokio-rustls``
     - Tokio 异步 I/O 适配
   * - ``hyper-rustls``
     - Hyper HTTP 适配

native-tls
----------

``native-tls`` 封装了操作系统的 TLS 实现（macOS Security.framework、Windows SChannel、Linux OpenSSL），
自动信任系统根证书。

.. code-block:: rust

    use native_tls::TlsConnector;
    use std::io::{Read, Write};
    use std::net::TcpStream;

    fn native_tls_client() -> Result<(), Box<dyn std::error::Error>> {
        let connector = TlsConnector::builder()
            .build()?;  // 自动使用系统根证书

        let stream = TcpStream::connect("example.com:443")?;
        let mut tls_stream = connector.connect("example.com", stream)?;

        tls_stream.write_all(b"GET / HTTP/1.0\r\nHost: example.com\r\n\r\n")?;
        let mut buf = Vec::new();
        tls_stream.read_to_end(&mut buf)?;
        println!("{}", String::from_utf8_lossy(&buf));

        Ok(())
    }

.. list-table:: rustls vs native-tls 对比
   :header-rows: 1

   * - 维度
     - rustls
     - native-tls
   * - 实现
     - 纯 Rust
     - 封装系统 TLS 库
   * - 依赖
     - ring 或 aws-lc-rs
     - OpenSSL / Security.framework / SChannel
   * - 编译
     - 无系统依赖，编译简单
     - 可能需要系统库
   * - 证书信任
     - 需显式加载根证书
     - 自动使用系统根证书
   * - FIPS
     - aws-lc-rs 后端支持
     - 取决于系统配置
   * - 推荐场景
     - 新项目、追求纯 Rust 栈
     - 需要系统证书集成

rcgen (证书生成)
----------------

``rcgen`` 用于生成自签名 X.509 证书，适合开发环境或内部服务间 mTLS。

.. code-block:: rust

    use rcgen::{CertificateParams, DistinguishedName, DnType, KeyPair, IsCa, BasicConstraints};
    use std::time::SystemTime;

    fn generate_self_signed_cert() -> Result<(String, String), Box<dyn std::error::Error>> {
        // CA 证书
        let mut ca_params = CertificateParams::new(vec!["My CA".to_string()])?;
        ca_params.is_ca = IsCa::Ca(BasicConstraints::Unconstrained);
        let ca_key = KeyPair::generate()?;
        let ca_cert = ca_params.self_signed(&ca_key)?;

        // 服务端证书
        let mut server_params = CertificateParams::new(vec![
            "localhost".to_string(),
            "my-server.local".to_string(),
        ])?;

        let mut dn = DistinguishedName::new();
        dn.push(DnType::CommonName, "My Server");
        dn.push(DnType::OrganizationName, "My Org");
        server_params.distinguished_name = dn;

        let server_key = KeyPair::generate()?;
        let server_cert = server_params.signed_by(&server_key, &ca_cert, &ca_key)?;

        Ok((
            server_cert.pem(),    // 证书
            server_key.serialize_pem(),  // 私钥
        ))
    }

x509-parser / x509-cert (证书解析)
--------------------------------------

.. code-block:: rust

    use x509_parser::prelude::*;

    fn parse_certificate(pem_data: &[u8]) -> Result<(), Box<dyn std::error::Error>> {
        let (_, pem) = x509_parser::pem::parse_x509_pem(pem_data)?;
        let (_, x509) = X509Certificate::from_der(&pem.contents)?;

        println!("Subject: {}", x509.subject());
        println!("Issuer: {}", x509.issuer());
        println!("Serial: {}", x509.raw_serial_as_string());
        println!("Not Before: {}", x509.validity().not_before.to_rfc2822()?);
        println!("Not After: {}", x509.validity().not_after.to_rfc2822()?);

        // SAN (Subject Alternative Names)
        if let Some((_, san)) = x509.subject_alternative_name()? {
            for name in &san.value.general_names {
                println!("SAN: {:?}", name);
            }
        }

        Ok(())
    }

mTLS (双向 TLS)
-----------------

.. code-block:: rust

    use rustls::{
        pki_types::{CertificateDer, PrivateKeyDer},
        server::WebPkiClientVerifier,
        RootCertStore, ServerConfig,
    };

    fn load_mtls_config(
        server_cert: Vec<CertificateDer<'static>>,
        server_key: PrivateKeyDer<'static>,
        client_ca_cert: CertificateDer<'static>,
    ) -> Result<ServerConfig, Box<dyn std::error::Error>> {
        // 客户端 CA 根证书
        let mut client_root_store = RootCertStore::empty();
        client_root_store.add(client_ca_cert)?;

        let client_verifier = WebPkiClientVerifier::builder(
            std::sync::Arc::new(client_root_store),
        )
        .build()?;

        let config = ServerConfig::builder()
            .with_client_cert_verifier(client_verifier)
            .with_single_cert(server_cert, server_key)?;

        Ok(config)
    }

Let's Encrypt / ACME
--------------------

``acme-micro`` 是一个轻量级的 ACME 客户端，用于自动获取 Let's Encrypt 证书。

.. code-block:: rust

    use acme_micro::{Acme, DirectoryUrl, Error};

    async fn issue_lets_encrypt_cert(
        domain: &str,
        email: &str,
    ) -> Result<(), Error> {
        let acme = Acme::new(DirectoryUrl::LetsEncryptStaging)?;
        // 创建账户
        let account = acme.new_account(email).await?;

        // 创建订单
        let order = account.new_order(domain).await?;

        // 完成 HTTP-01 验证（需在 web 服务器提供 challenge token）
        let auth = order.authorizations().await?;
        let challenge = auth[0].http_challenge().ok_or("no http challenge")?;
        // challenge.save_key_authorization() -> 写到 .well-known/acme-challenge/<token>

        // 等待验证完成
        challenge.validate(5000).await?;

        // 生成证书
        let (private_key, certificate) = order.finalize_csr(
            acme.create_csr(&[domain])?,
            5000,
        ).await?;

        std::fs::write("cert.pem", certificate)?;
        std::fs::write("key.pem", private_key)?;

        Ok(())
    }

.. note::

   Let's Encrypt 生产环境使用 ``DirectoryUrl::LetsEncrypt``；
   测试时务必先用 ``LetsEncryptStaging``，避免触发速率限制。

总结
-----

.. list-table:: TLS 与证书管理 Crate 总览
   :header-rows: 1

   * - Crate
     - 用途
     - 适用场景
   * - ``rustls``
     - 纯 Rust TLS 实现
     - HTTP 客户端/服务端、gRPC、WebSocket
   * - ``tokio-rustls``
     - rustls 的 Tokio 适配
     - 异步 TLS 通信
   * - ``native-tls``
     - 系统 TLS 封装
     - 需要系统证书信任链
   * - ``rcgen``
     - 自签名证书生成
     - 开发环境、mTLS 内部服务
   * - ``x509-parser``
     - X.509 证书解析
     - 证书信息提取、审计
   * - ``acme-micro``
     - ACME 客户端
     - 自动获取 Let's Encrypt 证书
