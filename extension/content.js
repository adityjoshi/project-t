// Content script to extract page content from various sites

let selectedText = '';

// Extract Amazon product information
function extractAmazonProduct() {
  const product = {
    title: '',
    price: '',
    rating: '',
    image: '',
    description: '',
    asin: '',
  };

  // Product title
  const titleSelectors = [
    '#productTitle',
    'h1.a-size-large',
    '[data-automation-id="title"]',
    'h1[data-automation-id="title"]',
  ];
  for (const selector of titleSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      product.title = el.innerText.trim();
      break;
    }
  }

  // Price
  const priceSelectors = [
    '.a-price .a-offscreen',
    '#priceblock_ourprice',
    '#priceblock_dealprice',
    '.a-price-whole',
    '[data-automation-id="price"]',
  ];
  for (const selector of priceSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      product.price = el.innerText.trim();
      break;
    }
  }

  // Rating
  const ratingSelectors = [
    '#acrPopover',
    '.a-icon-alt',
    '[data-automation-id="star-rating"]',
  ];
  for (const selector of ratingSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      const text = el.innerText || el.getAttribute('aria-label') || '';
      const match = text.match(/(\d+\.?\d*)\s*(out of|stars?)/i);
      if (match) {
        product.rating = match[1];
        break;
      }
    }
  }

  // Product image
  const imageSelectors = [
    '#landingImage',
    '#imgBlkFront',
    '#main-image',
    '[data-automation-id="product-image"] img',
    '.a-dynamic-image',
  ];
  for (const selector of imageSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      product.image = el.src || el.getAttribute('data-src') || el.getAttribute('data-old-src') || '';
      if (product.image) break;
    }
  }

  // Description
  const descSelectors = [
    '#feature-bullets',
    '#productDescription',
    '#productDescription_feature_div',
    '[data-automation-id="product-description"]',
  ];
  for (const selector of descSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      product.description = el.innerText.trim();
      if (product.description.length > 50) break;
    }
  }

  // ASIN
  const asinMatch = window.location.href.match(/\/dp\/([A-Z0-9]{10})/);
  if (asinMatch) {
    product.asin = asinMatch[1];
  }

  return product;
}

// Extract blog post content
function extractBlogPost() {
  const blog = {
    title: '',
    author: '',
    date: '',
    content: '',
    image: '',
  };

  // Title
  blog.title = document.title;
  const titleSelectors = [
    'article h1',
    '.post-title',
    '.entry-title',
    'h1.entry-title',
    '[itemprop="headline"]',
    'h1.post-title',
  ];
  for (const selector of titleSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      blog.title = el.innerText.trim();
      break;
    }
  }

  // Author
  const authorSelectors = [
    '[rel="author"]',
    '.author',
    '.post-author',
    '[itemprop="author"]',
    '.byline',
  ];
  for (const selector of authorSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      blog.author = el.innerText.trim();
      break;
    }
  }

  // Date
  const dateSelectors = [
    'time[datetime]',
    '.post-date',
    '.entry-date',
    '[itemprop="datePublished"]',
    '.published',
  ];
  for (const selector of dateSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      blog.date = el.innerText.trim() || el.getAttribute('datetime');
      break;
    }
  }

  // Featured image
  const imageSelectors = [
    'article img',
    '.post-thumbnail img',
    '.featured-image img',
    '[itemprop="image"]',
    'meta[property="og:image"]',
  ];
  for (const selector of imageSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      blog.image = el.src || el.getAttribute('content') || '';
      if (blog.image && !blog.image.includes('avatar')) break;
    }
  }

  // Content
  const contentSelectors = [
    'article',
    '.post-content',
    '.entry-content',
    '[itemprop="articleBody"]',
    '.post-body',
    'main article',
  ];
  for (const selector of contentSelectors) {
    const el = document.querySelector(selector);
    if (el) {
      const clone = el.cloneNode(true);
      clone.querySelectorAll('script, style, nav, aside, .ad, .advertisement').forEach(n => n.remove());
      blog.content = clone.innerText.trim();
      if (blog.content.length > 200) break;
    }
  }

  return blog;
}

// Detect todo list in selected text
function detectTodoList(text) {
  if (!text) return false;
  // Check for common todo patterns
  const todoPatterns = [
    /^[-*•]\s/m,           // Bullet points
    /^\d+\.\s/m,           // Numbered list
    /^\[[\sx]\]/im,        // Checkboxes [ ] or [x]
    /todo|to-do|task|item/i, // Contains todo keywords
  ];
  return todoPatterns.some(pattern => pattern.test(text));
}

// Extract YouTube video ID from URL
function getYouTubeVideoID(url) {
  const patterns = [
    /[?&]v=([a-zA-Z0-9_-]+)/,
    /youtu\.be\/([a-zA-Z0-9_-]+)/,
    /embed\/([a-zA-Z0-9_-]+)/,
  ];
  
  for (const pattern of patterns) {
    const match = url.match(pattern);
    if (match && match[1]) {
      return match[1];
    }
  }
  return null;
}

// Extract video information
async function extractVideoInfo() {
  const video = {
    title: '',
    channel: '',
    platform: '',
    thumbnail: '',
    description: '',
  };

  const url = window.location.href;
  const currentVideoID = getYouTubeVideoID(url);

  // YouTube
  if (url.includes('youtube.com') || url.includes('youtu.be')) {
    video.platform = 'YouTube';
    video.title = document.querySelector('h1.ytd-watch-metadata yt-formatted-string, h1.ytd-video-primary-info-renderer')?.innerText || document.title;
    video.channel = document.querySelector('#channel-name a, .ytd-channel-name a')?.innerText || '';
    const thumb = document.querySelector('meta[property="og:image"]');
    if (thumb) video.thumbnail = thumb.getAttribute('content');
    
    // PRIORITY 1: Get from YouTube's internal data (most reliable for FULL description)
    // IMPORTANT: Verify the video ID matches the current URL to avoid stale data
    try {
      // Try window.ytInitialPlayerResponse first (most reliable source)
      // But verify it's for the current video
      if (window.ytInitialPlayerResponse?.videoDetails) {
        const dataVideoID = window.ytInitialPlayerResponse.videoDetails.videoId;
        // Only use if video ID matches current URL
        if (currentVideoID && dataVideoID === currentVideoID) {
          const shortDesc = window.ytInitialPlayerResponse.videoDetails.shortDescription;
          if (shortDesc && shortDesc.length > 0) {
            video.description = shortDesc;
          }
        } else if (!currentVideoID) {
          // If we can't extract video ID from URL, use it anyway (fallback)
          const shortDesc = window.ytInitialPlayerResponse.videoDetails.shortDescription;
          if (shortDesc && shortDesc.length > 0) {
            video.description = shortDesc;
          }
        }
      }
      
      // Try window.ytInitialData (verify it's for current video)
      if (window.ytInitialData) {
        // Try to verify video ID from ytInitialData
        let isValidData = true;
        if (currentVideoID) {
          // Check if we can find the video ID in the data structure
          const videoPrimaryInfo = window.ytInitialData?.contents?.twoColumnWatchNextResults?.results?.results?.contents?.find(
            c => c.videoPrimaryInfoRenderer
          );
          if (videoPrimaryInfo?.videoPrimaryInfoRenderer?.videoActions?.menuRenderer?.topLevelButtons) {
            // Try to find video ID in the data
            const dataStr = JSON.stringify(window.ytInitialData);
            if (!dataStr.includes(currentVideoID)) {
              isValidData = false;
            }
          }
        }
        
        if (isValidData) {
          const videoDetails = window.ytInitialData?.contents?.twoColumnWatchNextResults?.results?.results?.contents?.find(
            c => c.videoSecondaryInfoRenderer
          )?.videoSecondaryInfoRenderer?.description?.runs;
          
          if (videoDetails && Array.isArray(videoDetails)) {
            const fullDesc = videoDetails.map(run => run.text || '').join('');
            if (fullDesc && fullDesc.length > (video.description?.length || 0)) {
              video.description = fullDesc;
            }
          }
        }
      }
      
      // Try to find description in any script tag (for cases where window objects aren't available)
      if (!video.description || video.description.length < 100) {
        const scripts = document.querySelectorAll('script');
        for (const script of scripts) {
          const text = script.textContent || '';
          if (text.includes('shortDescription') || text.includes('ytInitialPlayerResponse')) {
            try {
              // Try to extract from ytInitialPlayerResponse pattern (most reliable)
              const playerResponseMatch = text.match(/ytInitialPlayerResponse\s*=\s*({.+?});/s);
              if (playerResponseMatch) {
                try {
                  const playerData = JSON.parse(playerResponseMatch[1]);
                  // Verify video ID matches
                  const dataVideoID = playerData?.videoDetails?.videoId;
                  if (currentVideoID && dataVideoID && dataVideoID !== currentVideoID) {
                    // Skip this data - it's for a different video
                    continue;
                  }
                  const shortDesc = playerData?.videoDetails?.shortDescription;
                  if (shortDesc && shortDesc.length > (video.description?.length || 0)) {
                    video.description = shortDesc;
                    break; // Found it, no need to continue
                  }
                } catch (e) {
                  // Try regex extraction as fallback
                  const shortDescMatch = text.match(/"shortDescription"\s*:\s*"((?:[^"\\]|\\.)*)"/);
                  if (shortDescMatch && shortDescMatch[1]) {
                    const decoded = shortDescMatch[1]
                      .replace(/\\n/g, '\n')
                      .replace(/\\r/g, '\r')
                      .replace(/\\t/g, '\t')
                      .replace(/\\"/g, '"')
                      .replace(/\\\\/g, '\\')
                      .replace(/\\u([0-9a-fA-F]{4})/g, (match, code) => String.fromCharCode(parseInt(code, 16)));
                    if (decoded.length > (video.description?.length || 0)) {
                      video.description = decoded;
                      break;
                    }
                  }
                }
              }
              
              // Try to extract from var ytInitialData = {...}
              if (!video.description || video.description.length < 100) {
                const ytDataMatch = text.match(/var\s+ytInitialData\s*=\s*({.+?});/s);
                if (ytDataMatch) {
                  try {
                    const ytData = JSON.parse(ytDataMatch[1]);
                    // Verify video ID matches current URL
                    if (currentVideoID) {
                      const dataStr = JSON.stringify(ytData);
                      if (!dataStr.includes(currentVideoID)) {
                        // Skip - data is for a different video
                        continue;
                      }
                    }
                    const videoDetails = ytData?.contents?.twoColumnWatchNextResults?.results?.results?.contents?.find(
                      c => c.videoSecondaryInfoRenderer
                    )?.videoSecondaryInfoRenderer?.description?.runs;
                    
                    if (videoDetails && Array.isArray(videoDetails)) {
                      const fullDesc = videoDetails.map(run => run.text || '').join('');
                      if (fullDesc && fullDesc.length > (video.description?.length || 0)) {
                        video.description = fullDesc;
                      }
                    }
                  } catch (e) {
                    // Continue
                  }
                }
              }
            } catch (e) {
              // Continue to next script
            }
          }
        }
      }
    } catch (e) {
      console.log('Could not extract from YouTube internal data:', e);
    }
    
    // PRIORITY 2: Try DOM extraction (if internal data didn't work or was incomplete)
    if (!video.description || video.description.length < 100) {
      // First, try to expand the description if "Show more" button exists
      const showMoreButtons = [
        document.querySelector('ytd-expander #more'),
        document.querySelector('ytd-video-secondary-info-renderer #more'),
        document.querySelector('tp-yt-paper-button[id="more"]'),
        document.querySelector('button[aria-label*="more"]'),
        ...Array.from(document.querySelectorAll('button')).filter(btn => 
          btn.textContent.toLowerCase().includes('show more') || 
          btn.getAttribute('aria-label')?.toLowerCase().includes('more')
        ),
      ].filter(btn => btn !== null);
      
      for (const showMoreButton of showMoreButtons) {
        if (showMoreButton && (showMoreButton.textContent.toLowerCase().includes('show more') || 
            showMoreButton.getAttribute('aria-label')?.toLowerCase().includes('more'))) {
          try {
            showMoreButton.click();
            await new Promise(resolve => setTimeout(resolve, 500));
            break;
          } catch (e) {
            console.log('Could not expand description:', e);
          }
        }
      }
      
      // Try multiple selectors for YouTube description
      const descriptionSelectors = [
        'ytd-expander #content',
        'ytd-video-secondary-info-renderer #description',
        '#description-text',
        '#description',
        '.ytd-video-secondary-info-renderer #description',
        'yt-formatted-string#content-text',
        'yt-formatted-string.style-scope.ytd-video-secondary-info-renderer',
        'ytd-video-secondary-info-renderer yt-formatted-string',
      ];
      
      for (const selector of descriptionSelectors) {
        const descEl = document.querySelector(selector);
        if (descEl) {
          let descText = '';
          
          if (descEl.tagName === 'YT-FORMATTED-STRING') {
            descText = descEl.innerText || descEl.textContent || '';
            const ariaLabel = descEl.getAttribute('aria-label');
            if (ariaLabel && ariaLabel.length > descText.length) {
              descText = ariaLabel;
            }
          } else {
            descText = descEl.innerText || descEl.textContent || '';
            const ytFormattedStrings = descEl.querySelectorAll('yt-formatted-string');
            if (ytFormattedStrings.length > 0) {
              const parts = Array.from(ytFormattedStrings).map(el => {
                let text = el.innerText || el.textContent || '';
                const ariaLabel = el.getAttribute('aria-label');
                if (ariaLabel && ariaLabel.length > text.length) {
                  text = ariaLabel;
                }
                return text;
              }).filter(t => t.trim());
              if (parts.length > 0) {
                const combined = parts.join('\n');
                if (combined.length > descText.length) {
                  descText = combined;
                }
              }
            }
          }
          
          descText = descText
            .replace(/\s*(Show more|Show less)\s*/gi, '')
            .replace(/\n{3,}/g, '\n\n')
            .trim();
          
          if (descText && descText.length > (video.description?.length || 0)) {
            video.description = descText;
            break;
          }
        }
      }
    }
    
    // PRIORITY 3: Fallback to meta description
    if (!video.description || video.description.length < 50) {
      const metaDesc = document.querySelector('meta[name="description"]');
      if (metaDesc) {
        const metaText = metaDesc.getAttribute('content') || '';
        if (metaText && metaText.length > (video.description?.length || 0)) {
          video.description = metaText;
        }
      }
    }
    
    // PRIORITY 4: Try structured data as last resort
    if (!video.description || video.description.length < 100) {
      const structuredData = document.querySelector('script[type="application/ld+json"]');
      if (structuredData) {
        try {
          const data = JSON.parse(structuredData.textContent);
          if (data.description && data.description.length > (video.description?.length || 0)) {
            video.description = data.description;
          } else if (data.videoDetails && data.videoDetails.shortDescription) {
            if (data.videoDetails.shortDescription.length > (video.description?.length || 0)) {
              video.description = data.videoDetails.shortDescription;
            }
          }
        } catch (e) {
          console.log('Could not parse structured data:', e);
        }
      }
    }
  }
  // Vimeo
  else if (url.includes('vimeo.com')) {
    video.platform = 'Vimeo';
    video.title = document.querySelector('h1')?.innerText || document.title;
    const thumb = document.querySelector('meta[property="og:image"]');
    if (thumb) video.thumbnail = thumb.getAttribute('content');
    const descEl = document.querySelector('.description');
    if (descEl) video.description = descEl.innerText || '';
  }
  // Generic video detection
  else {
    const videoEl = document.querySelector('video');
    if (videoEl) {
      video.platform = 'Video';
      video.title = document.title;
      video.thumbnail = videoEl.poster || '';
    }
  }

  return video;
}

// Detect content type
function detectContentType() {
  const url = window.location.href.toLowerCase();
  const hostname = window.location.hostname.toLowerCase();

  if (hostname.includes('amazon.')) return 'amazon';
  if (hostname.includes('youtube.com') || hostname.includes('youtu.be')) return 'video';
  if (hostname.includes('vimeo.com')) return 'video';
  if (document.querySelector('article, .post, .blog-post, [itemprop="blogPost"]')) return 'blog';
  if (document.querySelector('video')) return 'video';
  
  return 'url';
}

// Extract general page content
function extractPageContent() {
  const title = document.title;
  
  let content = '';
  const contentSelectors = [
    'article',
    'main',
    '[role="main"]',
    '.content',
    '.post',
    '.entry-content',
    '#content',
  ];
  
  for (const selector of contentSelectors) {
    const element = document.querySelector(selector);
    if (element) {
      const clone = element.cloneNode(true);
      clone.querySelectorAll('script, style, nav, header, footer, aside, .ad').forEach(el => el.remove());
      content = clone.innerText.trim();
      if (content.length > 200) {
        break;
      }
    }
  }
  
  if (!content || content.length < 100) {
    const body = document.body.cloneNode(true);
    body.querySelectorAll('script, style, nav, header, footer, aside, .ad, .advertisement').forEach(el => el.remove());
    content = body.innerText.trim();
  }
  
  if (content.length > 5000) {
    content = content.substring(0, 5000) + '...';
  }
  
  return { title, content };
}

// Listen for messages from popup
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === 'getSelectedText') {
    const text = window.getSelection().toString().trim();
    selectedText = text;
    sendResponse({ selectedText: text });
    return true;
  } else if (request.action === 'extractContent') {
    // Use async IIFE to handle async extractVideoInfo
    (async () => {
      try {
        // For YouTube videos, wait a bit to ensure page has loaded current video data
        const url = window.location.href;
        if ((url.includes('youtube.com') || url.includes('youtu.be')) && url.includes('/watch')) {
          // Wait for page to be ready and video data to load
          await new Promise(resolve => setTimeout(resolve, 500));
          
          // Also wait for ytInitialPlayerResponse to be available/updated
          let attempts = 0;
          const maxAttempts = 10;
          while (attempts < maxAttempts) {
            const currentVideoID = getYouTubeVideoID(url);
            if (window.ytInitialPlayerResponse?.videoDetails?.videoId) {
              const dataVideoID = window.ytInitialPlayerResponse.videoDetails.videoId;
              // If video IDs match, or we can't extract from URL, proceed
              if (!currentVideoID || dataVideoID === currentVideoID) {
                break;
              }
            }
            await new Promise(resolve => setTimeout(resolve, 200));
            attempts++;
          }
        }
        
        const contentType = detectContentType();
        let data = {
          type: contentType,
          url: window.location.href,
          title: document.title,
          content: '',
          metadata: {},
        };

        if (contentType === 'amazon') {
          const product = extractAmazonProduct();
          data.title = product.title || document.title;
          data.content = `Price: ${product.price || 'N/A'}\nRating: ${product.rating || 'N/A'}\n\n${product.description || ''}`;
          data.metadata = {
            price: product.price,
            rating: product.rating,
            asin: product.asin,
            image: product.image,
          };
        } else if (contentType === 'blog') {
          const blog = extractBlogPost();
          data.title = blog.title || document.title;
          data.content = blog.content || '';
          data.metadata = {
            author: blog.author,
            date: blog.date,
            image: blog.image,
          };
        } else if (contentType === 'video') {
          const video = await extractVideoInfo();
          data.title = video.title || document.title;
          // Include description in content if available
          let contentParts = [`Platform: ${video.platform}`];
          if (video.channel) {
            contentParts.push(`Channel: ${video.channel}`);
          }
          if (video.description) {
            contentParts.push(`\n\nDescription:\n${video.description}`);
          }
          data.content = contentParts.join('\n');
          data.metadata = {
            platform: video.platform,
            channel: video.channel,
            thumbnail: video.thumbnail,
            description: video.description,
          };
        } else {
          const page = extractPageContent();
          data.title = page.title;
          data.content = page.content || selectedText;
        }

        // Handle selected text - if it's a todo, format it properly
        if (selectedText) {
          if (detectTodoList(selectedText)) {
            // If it's a todo list, use it as the main content
            data.type = 'text';
            data.title = data.title || 'Todo List';
            data.content = selectedText;
          } else {
            // Otherwise, prepend to existing content
            data.content = selectedText + '\n\n---\n\n' + data.content;
          }
        }

        sendResponse(data);
      } catch (error) {
        console.error('Error extracting content:', error);
        sendResponse({ 
          error: error.message,
          type: 'url',
          url: window.location.href,
          title: document.title,
          content: document.body.innerText.substring(0, 1000)
        });
      }
    })();
    return true; // Keep message channel open for async response
  }
  
  return false;
});

// Track text selection
document.addEventListener('mouseup', () => {
  selectedText = window.getSelection().toString().trim();
});
