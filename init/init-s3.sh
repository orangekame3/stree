#!/bin/bash

# AWS CLIの設定
export AWS_ACCESS_KEY_ID=dummy
export AWS_SECRET_ACCESS_KEY=dummy
export AWS_DEFAULT_REGION=us-east-1
export AWS_ENDPOINT=http://localstack:4566

# バケットの作成
aws --endpoint-url=$AWS_ENDPOINT s3 mb s3://my-bucket

# 階層の深いディレクトリ構造を作成し、CSVファイルを生成
mkdir -p data/test/dir1/dir1_1/dir1_1_1
mkdir -p data/test/dir1/dir1_2
mkdir -p data/test/dir2/dir2_1/dir2_1_1
mkdir -p data/test/dir3

# 各ディレクトリに異なる数のCSVファイルを生成
for i in {1..30}
do
  echo "name,age" > "data/test/dir1/dir1_1/dir1_1_1/file${i}.csv"
  echo "Alice,30" >> "data/test/dir1/dir1_1/dir1_1_1/file${i}.csv"
done

for i in {1..20}
do
  echo "name,age" > "data/test/dir1/dir1_2/file${i}.csv"
  echo "Bob,40" >> "data/test/dir1/dir1_2/file${i}.csv"
done

for i in {1..50}
do
  echo "name,age" > "data/test/dir2/dir2_1/dir2_1_1/file${i}.csv"
  echo "Charlie,35" >> "data/test/dir2/dir2_1/dir2_1_1/file${i}.csv"
done

for i in {1..10}
do
  echo "name,age" > "data/test/dir3/file${i}.csv"
  echo "Dave,50" >> "data/test/dir3/file${i}.csv"
done

# CSVファイルをS3バケットにアップロード
aws --endpoint-url=$AWS_ENDPOINT s3 sync data/test s3://my-bucket/test

# データディレクトリの削除
rm -rf data
