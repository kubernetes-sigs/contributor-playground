ワークショップに参加する前に
============================

# GitHub アカウントの作成
GitHub アカウント持っていない人は作成しておいてください。

# CNCF CLA へのサインアップ
* [LF sign up](https://identity.linuxfoundation.org/)
  The Linux Foundation ID の取得
  + 企業の開発者として登録する場合は、メールアドレスは企業ドメインのアドレスを登録する必要がある。
    _企業として参画する予定だが、ワークショップまでにに間に合わない場合は、個人開発者としてサインアップしておいてください。_
  + 持っていない場合は[ここ](https://identity.linuxfoundation.org/)から作成する
    - ユーザ名、メールアドレスなどを入力して登録ボタンを押下
    - メールアドレス宛に確認メールが来るので、メール中のURLにブラウザからアクセスする。
* CNCF CLA sign up
  1. 企業内で CNCF 開発参加者のリスト管理者に名前の追加を依頼する。
  2. CNCFの[当該ページ](https://identity.linuxfoundation.org/projects/cncf)の "Sign up to contribute to this project as an employee" をクリック。(この手順で "Groups:Authorized CNCF Contributors" が設定される。)
  3. The Linux Foudation ID を使ってログインする。
  4. "To sign up as a contributor to this project you must associate a GitHub account."が表示されることを確認する。
  5. [https://identity.linuxfoundation.org/user/me/hybridauth](https://identity.linuxfoundation.org/user/me/hybridauth)  で github アカウントと紐付ける。
  6. [https://identity.linuxfoundation.org/user/me](https://identity.linuxfoundation.org/user/me)  に下記が表示されることを確認する。
    - CNCF-<企業名が判断できる名称>
    - Groups:Authorized CNCF Contributors
    - **Note**: もし "CNCF-<企業名が判断できる名称>" が付いているが、"Groups:Authorized CNCF Contributors" が付いてない場合は、再度手順2から実施する。
  7. LF に登録しているメールアドレスを github のユーザ情報に登録する。
    - これをしないと PullRequest した時に Kubernetes の CI テストで失敗し、`cncf-cla: no`のラベルが付いてしまい、マージできない。
    - Github の Primary email address を LF に登録したアドレスとして登録すること。
    - Github 上で別の email を公開したい場合は、そのメールアドレスを Public email address として登録すればよい。

# Kubernetes の Slack を登録
* ここから登録します。https://slack.k8s.io/
* ログイン後に以下のチャンネルに入ってください。
  + #jp-dev

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

