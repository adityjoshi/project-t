# 🧠 Project Synapse - START HERE

## ⚡ Quick Start (3 Steps)

### Step 1: Get OpenAI API Key

1. Go to: **https://platform.openai.com/api-keys**
2. Sign up or login
3. Add payment method (required): **https://platform.openai.com/account/billing**
4. Create a new API key
5. Copy the key (starts with `sk-`)

**Cost**: ~$1-5/month for personal use (pay per API call)

### Step 2: Set Your API Key

```bash
cd /Users/aditya/Desktop/project-t
echo "OPENAI_API_KEY=sk-your-key-here" > .env
```

**Replace `sk-your-key-here` with your actual API key!**

### Step 3: Run Everything

```bash
docker-compose up -d --build
```

Wait 30-60 seconds, then open: **http://localhost:3000**

---

## 📋 Complete Command List

```bash
# 1. Navigate to project
cd /Users/aditya/Desktop/project-t

# 2. Create .env file with your OpenAI API key
echo "OPENAI_API_KEY=sk-your-actual-key-here" > .env

# 3. Start all services
docker-compose up -d --build

# 4. Check status
docker-compose ps

# 5. View logs
docker-compose logs -f

# 6. Access the app
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
```

---

## 🛠️ Useful Commands

```bash
# Stop services
docker-compose down

# Restart services
docker-compose restart

# View logs
docker-compose logs -f

# Rebuild after code changes
docker-compose up -d --build
```

---

## 📚 More Information

- **Quick Start**: [QUICKSTART.md](QUICKSTART.md)
- **All Commands**: [RUN_COMMANDS.md](RUN_COMMANDS.md)
- **Docker Guide**: [DOCKER.md](DOCKER.md)
- **Full Docs**: [README.md](README.md)

---

## ❓ Troubleshooting

**Services won't start?**
```bash
docker-compose logs
```

**Backend errors?**
```bash
docker-compose logs backend
```

**Check if API key is set:**
```bash
cat .env
```

**Port already in use?**
```bash
lsof -i :3000
lsof -i :8080
```

---

## ✅ Verify It's Working

1. Open http://localhost:3000
2. Click "+ Capture"
3. Save a test item
4. Try searching for it

If it works, you're all set! 🎉

