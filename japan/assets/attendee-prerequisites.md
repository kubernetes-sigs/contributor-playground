ワークショップに参加する前に
============================

# コミュニティ行動規範

まずは、CNCF コミュニティ行動規範をご一読ください。
https://github.com/cncf/foundation/blob/master/code-of-conduct-languages/jp.md

# オープンソース参加規約の確認

所属企業によっては、企業として参加、個人として参加に関わらずオープンソースに貢献するために規約を定めていることがあります。

# メールアドレスについて

以下の CNCF CLA サインアップにおいて、
**企業の開発者として登録する**場合は、メールアドレスは企業ドメインのアドレスを使用する必要があります。
_企業として参画する予定だが、ワークショップまでに間に合わない場合は、個人開発者としてサインアップしておいてください。_
具体的には、以下のメールアドレスです。

* GitHub の **Primary email address**
* Linux Foundation ID のメールアドレス

GitHub の Primary email address と Linux Foundation ID のメールアドレス が一致していないと Pull Request を行った時に Kubernetes の CI テストで失敗し、
`cncf-cla: no`のラベルが付いてしまい、マージできないため、登録するメールアドレスには注意が必要です。

# GitHub アカウントの作成

* GitHub アカウントを持っていない人は作成してください。
* GitHub アカウントをすでに持っている場合は Settings -> Emails で、**Primary email address** を
  Linux Foundation ID に登録するアドレスに設定してください。企業の開発者の場合はとくに注意。
  
なお、GitHub のコミットで別の email を使用したい場合は、
GitHub でそのメールアドレスを (Primary 以外の) **Emails** に登録してください。

# CNCF CLA へのサインアップ
* Linux Foundation ID の取得
  + [LF sign up](https://identity.linuxfoundation.org/) へアクセス
  + 持っていない場合は[ここ](https://identity.linuxfoundation.org/)から作成します。
    - ユーザ名、メールアドレスなどを入力して登録ボタンを押下
    - メールアドレス宛に確認メールが来るので、メール中のURLにブラウザからアクセス。
  + すでに持っている場合は、メールアドレスが要件を満たすものになっているかを確認。企業の開発者の場合は注意。
* Linux Foundation ID と GitHub アカウントの紐付け
  + Linux Foundation ID の [Social network logins](https://identity.linuxfoundation.org/user/me/hybridauth) に
    アクセスして GitHub アカウントと紐付けます。

* CNCF CLA sign up  
企業の開発者として登録する場合と個人の開発者として登録する場合で手順と確認項目が異なります。
  + **企業の開発者として登録する場合**
    1. 企業内で CNCF 開発参加者のリスト管理者に名前の追加を依頼します。
    2. CNCFの[当該ページ](https://identity.linuxfoundation.org/projects/cncf)の "Sign up to contribute to this project as an employee" をクリック。(この手順で "Groups:Authorized CNCF Contributors" が設定されます。)
    3. The Linux Foudation ID を使ってログインします。
    4. [https://identity.linuxfoundation.org/user/me](https://identity.linuxfoundation.org/user/me) に下記が表示されることを確認します。
        - Groups:Authorized CNCF Contributors
        - CNCF-<企業名が判断できる名称>
        - **Note**: もし "CNCF-<企業名が判断できる名称>" が付いているが、"Groups:Authorized CNCF Contributors" が付いてない場合は、再度「CNCF CLA sign up」のはじめから実施してください。
  + **個人の開発者として登録する場合**
    1. CNCFの[当該ページ](https://identity.linuxfoundation.org/projects/cncf)の "Sign up to contribute to this project as an individual" をクリック。(この手順で "Groups:Authorized CNCF Contributors" が設定されます。)
    2. The Linux Foudation ID を使ってログインします。
    3. [https://identity.linuxfoundation.org/user/me](https://identity.linuxfoundation.org/user/me) に下記が表示されることを確認します。
        - Groups:Authorized CNCF Contributors
        - **Note**: もし "Groups:Authorized CNCF Contributors" が付いてない場合は、再度「CNCF CLA sign up」のはじめから実施してください。

# Kubernetes の Slack に参加
* ここから参加します。
  + https://slack.k8s.io/
* ログイン後に以下のチャンネルに入ってください。
  + #jp-dev
    - レビュー依頼の練習に使用します。そのほか、実際のコントリビューション (Issue や PR) のそれぞれについて議論する場です。
  + #jp-mentoring
    - 事前の質疑応答、講義中のフォローに利用します。そのほか、コミュニティとの関わり方について議論する場です。
    - 質疑やフォローでは発言に対してスレッドを作って回答していく予定です。複数の質問が混じらないようにするためご協力ください。
    - Slackのスレッドについては [こちら](https://slack.com/intl/ja-jp/help/articles/115000769927-%E3%82%B9%E3%83%AC%E3%83%83%E3%83%89%E3%82%92%E4%BD%BF%E7%94%A8%E3%81%97%E3%81%A6%E4%BC%9A%E8%A9%B1%E3%82%92%E6%95%B4%E7%90%86%E3%81%99%E3%82%8B) を参照してください。

# コマンドのインストール
* 以下のコマンドをインストールしておいてください。
  + git

# Git や GitHub の操作
* 基本的な Git や GitHub の操作には慣れておいてください。
* 以下も参考にしてください。
  + https://github.com/kubernetes/community/blob/master/contributors/guide/contributor-cheatsheet/README-ja.md#%E8%B2%A2%E7%8C%AE%E3%81%99%E3%82%8B

# リポジトリのフォークとクローン
* 以下のリポジトリを自分の GitHub アカウントにフォークしておいてください。
  + https://github.com/kubernetes-sigs/contributor-playground
* 次に自分のアカウントにフォークしたリポジトリを自分の PC にクローンしておいてください。
