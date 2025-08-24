FROM golang:1.23-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git

# 作業ディレクトリを設定
WORKDIR /app

# go.modとgo.sumをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# 実行用の軽量イメージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk --no-cache add ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドされたバイナリをコピー
COPY --from=builder /app/main .

# ポートを公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./main"]