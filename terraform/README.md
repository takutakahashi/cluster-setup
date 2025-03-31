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

### ノード固有のAPI URL設定

各Proxmoxノードに対して異なるAPI URLを指定できるようになりました。これは複数のProxmoxサーバーが異なるURLを持つ環境で有用です。

#### 使用方法

`terraform.tfvars`ファイルで以下のように各ノードのAPI URLを設定します：

```hcl
# ノードごとにAPI URLとユーザー情報を設定
proxmox_nodes_config = {
  "node-1" = {
    api_url      = "https://proxmox-node1.example.com:8006/api2/json"
    user         = "terraform@pve"
    password     = "node1-password"
    tls_insecure = true
  },
  "node-2" = {
    api_url      = "https://proxmox-node2.example.com:8006/api2/json"
    # ユーザーとパスワードを省略すると、デフォルト値が使用されます
  }
}

# デフォルト設定（ノード固有の設定がない場合に使用）
proxmox_api_url = "https://proxmox-default.example.com:8006/api2/json"
proxmox_user = "terraform@pve"
proxmox_password = "your-default-password"
```

各ノードの設定で省略されたパラメータ（ユーザー、パスワードなど）には、グローバルなデフォルト値が使用されます。

#### 技術的な実装

この実装では、各ノードに対して個別のプロバイダーエイリアスを設定しています：

- `proxmox.node-1` - node-1用のプロバイダー
- `proxmox.node-2` - node-2用のプロバイダー
- `proxmox.node-3` - node-3用のプロバイダー

デフォルトの`proxmox`プロバイダー（エイリアスなし）も引き続き利用可能で、後方互換性を確保しています。

### Cloud-Initのカスタマイズ

起動時に実行するカスタムスクリプトを`terraform.tfvars`ファイルの`custom_script`変数で定義できます。

## トラブルシューティング

- **API接続エラー**: Proxmox APIのURLとユーザー認証情報を確認してください。
- **テンプレート関連のエラー**: テンプレート名とストレージの設定を確認してください。
- **Cloud-Init問題**: Proxmoxノードでcloud-init snippetsのディレクトリが適切にセットアップされていることを確認してください。

## リソース

- [Terraform Proxmoxプロバイダーのドキュメント](https://registry.terraform.io/providers/Telmate/proxmox/latest/docs)
- [Proxmox VEのドキュメント](https://pve.proxmox.com/wiki/Main_Page)