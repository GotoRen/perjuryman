## Usage
```
### RootCAとサーバ証明書を準備
> server
$ make create/certificate/server

### クライアント証明書を準備
> client
$ make create/certificate/client
```

## Certificates and Keys
| 鍵の種類 | 説明 |
| :--- | :--- |
| `ca.pem` | ルート証明書 |
| `ca_private_key.pem` | ルート証明書のRSA秘密鍵 |
| `server.pem` | サーバ証明書 |
| `server.key` | サーバ証明書のRSA秘密鍵 |
| `client_certreq.pem` | クライアント証明書 |
| `client.pem` | クライアント証明書のRSA公開鍵 |
| `client_private.key` | クライアント証明書のRSA秘密鍵 |
