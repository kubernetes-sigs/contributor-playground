ワークショップに参加する前に
============================

# メールアドレスについて

以下の CNCF CLA サインアップにおいて、
**企業の開発者として登録する**場合は、メールアドレスは企業ドメインのアドレスを使用する必要があります。
_企業として参画する予定だが、ワークショップまでに間に合わない場合は、個人開発者としてサインアップしておいてください。_
具体的には、以下のメールアドレスです。

* GitHub の **Primary email address**
* Linux Foundation ID のメールアドレス

Primary email address が一致していないと Pull Request を行った時に Kubernetes の CI テストで失敗し、
`cncf-cla: no`のラベルが付いてしまい、マージできないため、登録するメールアドレスには注意が必要です。

# GitHub アカウントの作成

* GitHub アカウントを持っていない人は作成してください。
* GitHub アカウントをすでに持っている場合は Settings -> Emails で、**Primary email address** が
  Linux Foundation ID に登録するアドレスに設定してください。企業の開発者の場合はとくに注意。
  
なお、Github のコミットで別の email を使用したい場合は、
そのメールアドレスを (Primary 以外の) **Emails** に登録しておけばよい。

# CNCF CLA へのサインアップ
* Linux Foundation ID の取得
  + [LF sign up](https://identity.linuxfoundation.org/) へアクセス
  + 持っていない場合は[ここ](https://identity.linuxfoundation.org/)から作成する。
    - ユーザ名、メールアドレスなどを入力して登録ボタンを押下
    - メールアドレス宛に確認メールが来るので、メール中のURLにブラウザからアクセスする。
  + すでに持っている場合は、メールアドレスが要件を満たすものになっているかを確認、企業の開発者の場合は注意。
* Linux Foundation ID と GitHub アカウントの紐付け
  + Linux Foundationn ID の [Social network logins](https://identity.linuxfoundation.org/user/me/hybridauth) に
    アクセスし github アカウントと紐付ける。

* CNCF CLA sign up
  1. 企業内で CNCF 開発参加者のリスト管理者に名前の追加を依頼する。
  2. CNCFの[当該ページ](https://identity.linuxfoundation.org/projects/cncf)の "Sign up to contribute to this project as an employee" をクリック。(この手順で "Groups:Authorized CNCF Contributors" が設定される。)
    * 個人開発者の場合は "Sign up to contribute to this project as an individual" をクリック
  3. The Linux Foudation ID を使ってログインする。
  4. [https://identity.linuxfoundation.org/user/me](https://identity.linuxfoundation.org/user/me) に下記が表示されることを確認する。
    - Groups:Authorized CNCF Contributors
    - CNCF-<企業名が判断できる名称> (← 企業の開発者として登録している場合のみ)
    - **Note**: もし "CNCF-<企業名が判断できる名称>" が付いているが、"Groups:Authorized CNCF Contributors" が付いてない場合は、再度手順2から実施する。

# Kubernetes の Slack を登録
* ここから登録します。https://slack.k8s.io/
* ログイン後に以下のチャンネルに入ってください。
  + #jp-dev
    このチャンネルの中で、事前の質疑応答、レビュー依頼の練習、フォローなど行っていきます。

# コマンドのインストール
* 以下のコマンドをインストールしておいてください。
  + git
  + golang
  + docker

# Git や GitHub の操作
* 基本的な Git や GitHub の操作には慣れておいてください。
* 以下も参考にしてください。
  + https://github.com/kubernetes/community/blob/master/contributors/guide/contributor-cheatsheet/README-ja.md#%E8%B2%A2%E7%8C%AE%E3%81%99%E3%82%8B

# リポジトリのフォークとクローン
* 以下のリポジトリを自分の GitHub アカウントにフォークしておいてください。
  + https://github.com/kubernetes-sigs/contributor-playground
  + https://github.com/kubernetes/kubernetes
* 次に自分のアカウントにフォークしたリポジトリを自分の PC にクローンしておいてください。
