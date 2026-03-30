# Frontend (Next.js)

漁港のせりシステムのフロントエンドアプリケーションです。

## 技術構成 (Tech Stack)

- **Framework**: Next.js 16 (App Router)
- **Library**: React 19
- **Language**: TypeScript
- **Styling**: [Panda CSS](https://panda-css.com/) (Styled System, Atomic CSS)
- **Data Fetching**: [TanStack Query v5](https://tanstack.com/query/latest)
- **Forms**: React Hook Form + Zod
- **I18n**: next-intl
- **Testing**: Vitest / React Testing Library

## アーキテクチャ (Architecture)

ドメイン駆動の「Feature-based」構造と、関心の分離（Separation of Concerns）を意識したレイヤード構造を採用しています。

詳細は [ARCHITECTURE.md](./docs/ARCHITECTURE.md) を参照してください。

### ディレクトリ構成

```text
frontend/
├── app/                  # Next.js App Router (Routing & Layouts)
├── src/
│   ├── features/         # ドメインごとの機能カプセル化
│   │   └── [feature]/
│   │       ├── components/  # 機能固有の UI
│   │       ├── states/      # UI ロジック・状態管理
│   │       ├── queries/     # 機能固有のデータフェッチ
│   │       └── types/       # 機能固有の型定義（ドメインモデル + UI型）
│   ├── data/             # データフェッチレイヤー
│   │   ├── api/          # プリミティブな API 呼び出し (apiClient)
│   │   ├── queries/      # ドメインごとの TanStack Query フック & キー
│   │   └── entities/     # サーバー側データ構造の型定義（@entities/*）
│   ├── components/       # 共有 UI コンポーネント (Atomic Design)
│   │   ├── atoms/        # 汎用パーツ (Button, Input, Text)
│   │   ├── molecules/    # 複数の Atom を組み合わせた塊
│   │   ├── organisms/    # 具体的かつ機能的なコンポーネント
│   │   ├── templates/    # ページレイアウト
│   │   └── functionals/  # Context Provider や初期化ロジック
│   ├── core/             # アプリケーション基盤
│   │   ├── api/          # API クライアント (fetch ラッパー)
│   │   ├── i18n/         # 国際化メッセージとリクエスト設定
│   │   └── styles/       # グローバルスタイル
│   ├── schemas/          # Zod スキーマ・入力検証
│   └── libs/             # 外部ライブラリ設定・生成物 (Panda CSS 等)
```

## 開発環境 (Development)

### 前提条件

- **Node.js**: v20+
- **npm**: CI と README の手順は npm ベースです

### コマンド

```bash
cd frontend
npm install   # 依存関係のインストール
npm run dev   # 開発サーバー起動
npm run test  # テスト実行
npm run lint  # リンター起動
```

ブラウザで `https://localhost` (Nginx 経由) または `http://localhost:3000` (直通) を開いてください。
