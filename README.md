# neis-gym-scraper

[![Check Open Classes](https://github.com/iku/neis-tools/actions/workflows/check-open-classes.yml/badge.svg)](https://github.com/iku/neis-tools/actions/workflows/check-open-classes.yml)

このツールは、NEIS体操教室の予約サイトを自動的にスクレイピングして空きクラスを確認し、空きが見つかった場合にSlackへ通知を送信します。人気のクラスに空きが出た際、すぐに把握できるように定期的に実行することを想定して設計されています。

## 機能

- 向こう4週間分のNEIS-GYMスケジュールをスクレイピングします。
- 週末（土曜日と日曜日）のクラスを特に対象としています。
- 「満員（Full）」と表示されていないクラスを特定します。
- 空き状況の詳細を含むフォーマットされた通知をSlackチャンネルに送信します。
- 堅牢な監視のために、エラー通知を別のSlackチャンネルに送信します。
- 環境変数または `.env` ファイルで設定可能です。
- 自動定期チェック用のGitHub Actionsワークフローが含まれています。

## 仕組み

1.  **設定の読み込み**: 環境変数からSlack Webhook URLを読み込みます。
2.  **ウェブサイトのスクレイピング**: NEIS-GYMのスケジュールページにHTTPリクエストを送信します。より広い期間を確認するために、向こう4週間分を繰り返し処理します。
3.  **HTMLの解析**: HTMLレスポンスを解析してスケジュール表を見つけます。土曜日と日曜日に対応する行を探します。
4.  **空き状況の確認**: 週末の行の中で、各クラスの枠を検査します。枠が有効で、かつ「Full（満員）」というテキストが含まれていない場合、空きがあるとみなされます。
5.  **通知の送信**: 1つ以上の空き枠が見つかった場合、各空き枠の日付、時間、クラス名をリストアップしたメッセージを作成し、設定されたSlack Webhook URLに送信します。処理中にエラーが発生した場合は、別のWebhookを使用してアラートを送信します。

## 前提条件

- **Go**: バージョン 1.19 以上。
- **Slack Webhook URLs**: Slackワークスペースからの有効なIncoming Webhook URLが必要です。

## インストール

1.  **リポジトリをクローンします:**
    ```bash
    git clone https://github.com/iku/neis-tools.git
    cd neis-tools
    ```

2.  **依存関係をインストールします:**
    ```bash
    go mod download
    ```

## 設定

アプリケーションは、Slack通知を送信するために2つの環境変数を設定する必要があります。

プロジェクトのルートディレクトリに `.env` ファイルを作成してください：

```env
# .env file
SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/SUCCESS/WEBHOOK"
SLACK_WEBHOOK_URL_FOR_ERROR="https://hooks.slack.com/services/YOUR/ERROR/WEBHOOK"
```

Alternatively, you can export these variables in your shell environment. The application will fail to start if these are not set.

## Usage

To run the scraper locally:

```bash
go run .
```

The tool will log its progress to the console and send notifications to Slack if any openings are found.

## Automated Execution with GitHub Actions

This repository includes a GitHub Actions workflow (`.github/workflows/check-open-classes.yml`) to run the scraper on a schedule.

-   **Trigger**: The workflow is set to be triggered manually (`workflow_dispatch`). A cron schedule to run every 10 minutes is commented out and can be enabled if desired.
-   **Secrets**: The workflow requires the `SLACK_WEBHOOK_URL` and `SLACK_WEBHOOK_URL_FOR_ERROR` to be configured as secrets in your GitHub repository settings. Navigate to `Settings` > `Secrets and variables` > `Actions` to add them.
-   **Execution**: The workflow checks out the code, downloads a pre-compiled binary from a GitHub Release (it expects a release with the binary `neis-tools`), and executes it. For it to work correctly, you must create a release and attach the compiled binary.

### Building the binary for releases:

```bash
# For Linux (as used in GitHub Actions)
GOOS=linux GOARCH=amd64 go build -o neis-tools .
```

## License

This project is licensed under the MIT License.
