# 未読メール(IMAP)通知アプリ
- このプログラムは，IMAPサーバーから未読メールを取得し，LINEに通知するアプリケーションです．

## 概要
- IMAPサーバーに接続し，INBOX内の未読メールを取得します．
- 各メールの送り元と件名を表示します．
- 件名がエンコードされている場合はデコードします（ISO-2022-JPエンコーディングに対応）．
- 取得したメールの情報をLINEに通知します．
- 最新のメールシーケンス番号をファイルに保存し，次回実行時に未読メールの重複取得を防ぎます．

## 前提条件
- go version go1.22.1 darwin/arm64 (筆者の環境)
- LINEのmessaging apiが利用できること

## installation
- `mv .env.sample .env`
    - SERVER: IMAPサーバのアドレスおよびポート番号，993を利用していなければ適宜変更してください
    - ID: IMAPサーバにログインするためのID
    - PASSWORD: IMAPサーバにログインするためのパスワード
    - TOKEN: LINEのMessaging API利用のためのトークン    
- `go build`
- 以下のように `/etc/etc/systemd/system/get-mail.service` を作成
```
[Unit]
Description=Run getmail every 5 minutes

[Service]
ExecStart=/path/to/your/get-mail/get-mail
WorkingDirectory=/path/to/your/get-mail
```

- 以下のように `/etc/systemd/system/get-mail.timer` を作成
```
[Unit]
Description=Run getmail every 5 minutes

[Timer]
OnBootSec=5min
OnUnitActiveSec=5min
Unit=get-mail.service

[Install]
WantedBy=timers.target
```

- `sudo systemctl enable get-mail.timer`
- `sudo systemctl start get-mail.timer`