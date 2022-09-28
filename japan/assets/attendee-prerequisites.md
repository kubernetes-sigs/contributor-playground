ワークショップに参加する前に
============================

ワークショップ参加予定の皆さんに事前に実施いただく作業を説明します。

# コミュニティ行動規範

まずは[CNCF コミュニティ行動規範](https://github.com/cncf/foundation/blob/master/code-of-conduct-languages/jp.md)をご一読ください。

# オープンソース参加規約の確認

所属企業によっては、**企業として参加**、**個人として参加**に関わらずオープンソースに貢献するために規約を定めていることがあります。

# メールアドレスについて

以下の [CNCF CLA サインアップ](#cncf-cla-へのサインアップ)において**企業の開発者として登録する**場合は、メールアドレスは企業ドメインのアドレス等を使用する必要があります。  
詳しくは企業の[CLA Manager](https://docs.linuxfoundation.org/lfx/easycla/v2-current/corporate-cla-managers)に確認してください。  
企業として参画する予定だがワークショップまでに確認が間に合わない場合は、個人開発者としてサインアップすることを検討してください。

# GitHub アカウントの作成

* GitHub アカウントを持っていない人は作成してください。
* 企業の開発者として参加する人で、GitHub アカウントをすでに持っている場合は GitHub の Settings -> Emails で、**email address** の項目に企業で使う[メールアドレスを追加](https://docs.github.com/ja/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/adding-an-email-address-to-your-github-account)してください。

# CNCF CLA へのサインアップ
Kubernetesコミュニティでは2022年3月頃にCLAの手順がEasyCLAに移行しており、以前とは手順が変更されています。

事前作業としては、以下の[サインアップ事前作業](#サインアップ事前作業)を実施してください。  
実際のサインアップはワークショップの中でPullRequestを作成した際に実施します。  

## サインアップ事前作業
まず以下を実施してください。
- CNCF CLA サインアップのときに使うメールアドレスと[git コマンドに設定するコミットメールアドレスは、同じものを設定してください。](https://docs.github.com/ja/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/setting-your-commit-email-address#setting-your-commit-email-address-in-git)


次に、以下は企業の開発者として登録する場合と個人の開発者として登録する場合で作業が異なります。

- 
  **企業の開発者として登録する場合**
  - 企業内で CNCF 開発参加者のリスト管理者(CLA Manager)に自分の名前の追加を依頼してください。
  
  **個人の開発者として登録する場合**
  - 特に作業は必要ありません。

## サインアップ作業
**実際のサインアップはワークショップ中に実施します。**  
**[次のステップのSlackに参加](#kubernetes-の-slack-に参加)以降を実施してください。**  
サインアップを一度実施すると、以降のKubernetesコミュニティでの活動で再度実施する必要はありません。  
1. PullRequestを作成すると linux-foundation-easycla から以下のようにコメントされます。  
  `Please here to be authorized` をクリックします。  
  ![](https://user-images.githubusercontent.com/69111235/152226443-f6fe61ee-0e92-46c5-b6ea-c0deb718a585.png)  

1. `Authorize LF-Engineering` の緑のボタンをクリックします。  
   ![](https://user-images.githubusercontent.com/69111235/152228712-7d22f9d0-9f3c-4226-9ee0-bacba4b47725.png)  
 
1. - 企業の開発者としてサインアップするときは画面左の `Proceed as a Corporate Contributer` の青のボタンをクリックします。  
   - 個人の開発者としてサインアップするときは画面右の `Proceed as an Indivisual Contributer` の緑のボタンクリックします。  
  ![](https://user-images.githubusercontent.com/69111235/152224818-1246453a-b086-4a57-9d14-c10d62ad438f.png)  

1. 以降は画面の指示に従ってください。

その他の注意事項等は[コミュニティのドキュメント](https://github.com/kubernetes/community/blob/master/CLA.md)を参照してください。

# Kubernetes の Slack に参加
* [こちら](https://slack.k8s.io/)から参加します。
* ログイン後に以下のチャンネルに入ってください。
  + #jp-dev
    - レビュー依頼の練習に使用します。そのほか、実際のコントリビューション (Issue や PullRequest) のそれぞれについて議論する場です。
  + #jp-mentoring
    - 事前の質疑応答、講義中のフォローに利用します。そのほか、コミュニティとの関わり方について議論する場です。
    - 質疑やフォローでは発言に対してスレッドを作って回答していく予定です。複数の質問が混じらないようにするためご協力ください。
    - Slackのスレッド機能については [こちら](https://slack.com/intl/ja-jp/help/articles/115000769927-%E3%82%B9%E3%83%AC%E3%83%83%E3%83%89%E3%82%92%E4%BD%BF%E7%94%A8%E3%81%97%E3%81%A6%E4%BC%9A%E8%A9%B1%E3%82%92%E6%95%B4%E7%90%86%E3%81%99%E3%82%8B)を参照してください。

# コマンドのインストール
* 以下のコマンドをインストールしておいてください。
  + git

# Git や GitHub の操作
* 基本的な Git や GitHub の操作に慣れておいてください。
* [Kubernetesコントリビューターチートシートのこちらの節](https://github.com/kubernetes/community/blob/master/contributors/guide/contributor-cheatsheet/README-ja.md#%E8%B2%A2%E7%8C%AE%E3%81%99%E3%82%8B)から、GitHub での開発の概要を把握しておいてください。


# リポジトリのフォークとクローン
* 以下のリポジトリを自分の GitHub アカウントにフォークしておいてください。
  + https://github.com/kubernetes-sigs/contributor-playground
* 次に自分のアカウントにフォークしたリポジトリを自分の PC にクローンしておいてください。
