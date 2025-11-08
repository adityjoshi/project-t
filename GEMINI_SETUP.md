# Google Gemini API Setup (Free)

Project Synapse now uses **Google Gemini** as the default AI provider, which is **completely free**!

## Get Your Free Gemini API Key

1. **Go to Google AI Studio**: https://makersuite.google.com/app/apikey
2. **Sign in** with your Google account
3. **Click "Create API Key"**
4. **Copy the API key** (starts with `AIza...`)

## Set Your API Key

### Option 1: Using .env file

```bash
# In project root
echo "GEMINI_API_KEY=AIza-your-key-here" >> .env
```

### Option 2: Environment variable

```bash
export GEMINI_API_KEY=AIza-your-key-here
```

## Restart Services

```bash
docker-compose restart backend
```

## Verify It's Working

Try creating a new item in the app. It should:
- ✅ Generate summaries
- ✅ Create tags automatically
- ✅ Enable semantic search

## Free Tier Limits

Google Gemini free tier includes:
- **60 requests per minute**
- **1,500 requests per day**
- More than enough for personal use!

## Switch Back to OpenAI (Optional)

If you want to use OpenAI instead:

1. Set `AI_PROVIDER=openai` in `.env`
2. Set `OPENAI_API_KEY=sk-...` in `.env`
3. Restart backend: `docker-compose restart backend`

## Troubleshooting

**"GEMINI_API_KEY not set" warning?**
- Make sure you've added the key to `.env` file
- Restart the backend container

**API errors?**
- Verify your API key is correct
- Check you haven't exceeded rate limits
- Make sure the key starts with `AIza`

