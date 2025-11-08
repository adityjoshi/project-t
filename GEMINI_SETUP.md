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

## Understanding Quota Exceeded Errors

### What Happens When Quota is Exceeded?

If you see errors like "You exceeded your current quota" or "429 rate limit":

1. **Item Creation Still Works**: Items are saved successfully even if summary generation fails
2. **Temporary Summary**: Items get a temporary summary (first 200 characters) until AI summary is generated
3. **Automatic Fallback**: If `OPENAI_API_KEY` is configured, the system automatically uses OpenAI when Gemini quota is exceeded
4. **Graceful Degradation**: The app continues to function normally; only AI summarization is affected

### Why Quota Gets Exceeded

- **Free Tier Limits**: The free tier has daily/monthly request limits
- **Development/Testing**: High usage during development can quickly exhaust quotas
- **Multiple Features**: Summarization, categorization, tagging, and embeddings all use the API
- **Rate Limits**: Even within daily limits, per-minute rate limits can be hit

### Solutions

1. **Wait for Reset**: Quotas usually reset daily - wait 24 hours
2. **Upgrade Plan**: Upgrade to paid Gemini API for higher limits
3. **Use OpenAI Fallback**: Add `OPENAI_API_KEY` to `.env` for automatic fallback
4. **Monitor Usage**: Check your usage at https://ai.dev/usage?tab=rate-limit
5. **Reduce Usage**: Temporarily disable some AI features if needed

### How to Explain This in Documentation/Interviews

**Short Explanation**:
"Project Synapse uses Google Gemini API for AI-powered summarization. The free tier has rate limits (60 requests/minute, 1,500/day). When quota is exceeded, the system gracefully degrades by using OpenAI as a fallback (if configured) or skipping summarization without affecting core functionality. Items are still saved successfully with temporary summaries until quota resets."

**Technical Explanation**:
"The summarization service implements a fallback mechanism: it attempts to use Gemini 2.5 Pro first, then falls back to Gemini 2.5 Flash, and finally to OpenAI if Gemini quota is exceeded. This ensures high availability and graceful degradation. The system logs quota errors but continues operating normally, allowing users to refresh summaries later when quota resets."

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

