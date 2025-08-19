# Nakama Quick Reference Guide

## 🚀 **How to Run and Validate Nakama**

### **1. Start Services**
```bash
# Start Nakama with CockroachDB
docker compose up -d

# Check service status
docker ps | grep -E "(nakama|cockroach)"
```

### **2. Quick Validation**
```bash
# Basic test
node test-nakama.js

# Comprehensive validation
node validate-nakama.js
```

### **3. Manual Testing**

#### **API Endpoints**
```bash
# Authenticate user
curl -s -X POST "http://127.0.0.1:7350/v2/account/authenticate/custom?create=true" \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic ZGVmYXVsdGtleTo=" \
  -d '{"id":"testuser123"}' | jq .

# Get account info (replace TOKEN)
curl -s "http://127.0.0.1:7350/v2/account" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq .
```

#### **Web Interfaces**
- **Nakama Console**: http://localhost:7351 (admin/password)
- **CockroachDB Admin**: http://localhost:8080

### **4. Service Information**

| Service | Port | Status | URL |
|---------|------|--------|-----|
| Nakama API | 7350 | ✅ Running | http://localhost:7350 |
| Nakama Console | 7351 | ✅ Running | http://localhost:7351 |
| CockroachDB | 26257 | ✅ Running | - |
| CockroachDB Admin | 8080 | ✅ Running | http://localhost:8080 |

### **5. Common Commands**

```bash
# View logs
docker compose logs nakama
docker compose logs cockroachdb

# Restart services
docker compose restart nakama
docker compose restart cockroachdb

# Stop all services
docker compose down

# Start all services
docker compose up -d

# Check health
curl http://localhost:7350/v2/health
curl http://localhost:8080/health
```

### **6. Development Workflow**

#### **With Custom Go Module**
```bash
# 1. Edit your module
vim modules/main.go

# 2. Build for Linux
docker run --rm -v $(pwd)/modules:/app -w /app golang:1.21 \
  sh -c "go mod tidy && go build -buildmode=plugin -o main.so main.go"

# 3. Deploy module
cp modules/main.so data/modules/

# 4. Restart Nakama
docker compose restart nakama

# 5. Test
node test-nakama.js
```

#### **Without Custom Module (Current Setup)**
```bash
# Just test basic functionality
node validate-nakama.js
```

### **7. Troubleshooting**

#### **Service Won't Start**
```bash
# Check if ports are in use
lsof -i :7350
lsof -i :7351
lsof -i :8080

# Clean up containers
docker compose down --remove-orphans
docker compose up -d
```

#### **Database Connection Issues**
```bash
# Check CockroachDB logs
docker compose logs cockroachdb

# Check Nakama logs
docker compose logs nakama

# Verify database is running
docker exec -it effective-golang-cockroachdb-1 cockroach sql --insecure
```

#### **Module Loading Issues**
```bash
# Check if module exists
ls -la data/modules/

# Rebuild module
docker run --rm -v $(pwd)/modules:/app -w /app golang:1.21 \
  sh -c "go build -buildmode=plugin -o main.so main.go"

# Copy and restart
cp modules/main.so data/modules/ && docker compose restart nakama
```

### **8. Validation Results**

✅ **All Core Services Working:**
- Authentication system
- Account management
- Friends system
- Storage system
- Leaderboards
- HTTP API
- Web console
- Database health

### **9. Next Steps**

1. **Read the setup guide**: `docs/nakama-setup-guide.md`
2. **Learn Nakama concepts**: `docs/nakama-learning-guide.md`
3. **Try examples**: `tutorials/examples/`
4. **Build your game**: Using the patterns you've learned

### **10. Useful Links**

- **Documentation**: https://heroiclabs.com/docs/
- **GitHub**: https://github.com/heroiclabs/nakama
- **Discord**: https://discord.gg/heroiclabs
- **Forums**: https://forum.heroiclabs.com/

---

**🎉 Your Nakama server is fully operational and ready for development!**
