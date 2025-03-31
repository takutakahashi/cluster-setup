# Proxmox VM Provisioning with Terraform

このディレクトリには、Proxmox VEに複数のVMをプロビジョニングするためのTerraformコードが含まれています。

## 前提条件

- Terraform v1.0.0以上
- Proxmox VE 6.x以上
- Proxmox APIアクセス権限を持つユーザー
- Cloud-Init対応のUbuntuテンプレート

## 初期設定

1. Terraformをインストールします。
2. `terraform.tfvars.example`ファイルを`terraform.tfvars`としてコピーし、環境に合わせて設定します。
3. Proxmox VEにCloud-Init対応のUbuntuテンプレートを用意します。

## 使用方法

### 1. 初期化

```bash
cd terraform
terraform init
```

### 2. 設定の検証

```bash
terraform validate
```

### 3. 実行プランの確認

```bash
terraform plan
```

### 4. リソースのデプロイ

```bash
terraform apply
```

### 5. リソースの破棄

```bash
terraform destroy
```

## カスタマイズ

### VM設定のカスタマイズ

`terraform.tfvars`ファイル内の`vms`マップを編集して、VMの構成を変更できます。

### Cloud-Initのカスタマイズ

起動時に実行するカスタムスクリプトを`terraform.tfvars`ファイルの`custom_script`変数で定義できます。

## トラブルシューティング

- **API接続エラー**: Proxmox APIのURLとユーザー認証情報を確認してください。
- **テンプレート関連のエラー**: テンプレート名とストレージの設定を確認してください。
- **Cloud-Init問題**: Proxmoxノードでcloud-init snippetsのディレクトリが適切にセットアップされていることを確認してください。

## リソース

- [Terraform Proxmoxプロバイダーのドキュメント](https://registry.terraform.io/providers/Telmate/proxmox/latest/docs)
- [Proxmox VEのドキュメント](https://pve.proxmox.com/wiki/Main_Page)