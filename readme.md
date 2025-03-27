# 🔐 ssmtools

SSM-based CLI utilities for secure, auditable access and file transfer to EC2 instances — without using SSH.

This toolkit includes:

- `ssmssh` — Securely open an interactive SSM session to a given EC2 instance
- `ssmcp` — Copy files **to** and **from** EC2 instances using an SSM port-forward tunnel (like `scp` over SSM)

Built in Go, with support for AWS IAM roles via OneLogin SSO.

---

## ⚙️ Features

✅ No SSH or public key setup required  
✅ Role-based access via `--role` (OneLogin SSO)  
✅ DNS resolution or Name tag-based targeting  
✅ SSM port-forwarding for file uploads  
✅ Works even when the EC2 instance has **no public IP**  
✅ Fully auditable via AWS CloudTrail  


---

## 🚀 Build & Install

### 🧪 Local build

```bash
go build -o bin/ssmssh ./cmd/ssmssh
go build -o bin/ssmcp ./cmd/ssmcp
```


## 🔧 Usage
🔐 Login via OneLogin SSO
All commands support:

--role <onelogin-role>: Use OneLogin-based SSO

--profile <aws-profile>: Use a pre-configured AWS profile

### 🖥️ ssmssh — Connect to EC2 over SSM
```bash
./ssmssh --hostname my-bastion-host --region eu-west-1 --role onelogin-dev
```
Resolves host via:

Private IP (via DNS)

EC2 Name tag fallback

### 📁 ssmcp — Copy file to EC2
```bash
./ssmcp --hostname my-bastion-host \
        --from ./local.txt \
        --to /home/ubuntu/shared/remote.txt \
        --region eu-west-1 \
        --role onelogin-dev
```
Uses SSM port-forwarding to upload the file via an HTTP server running on the instance.
