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

// Extract video information
function extractVideoInfo() {
  const video = {
    title: '',
    channel: '',
    platform: '',
    thumbnail: '',
  };

  const url = window.location.href;

  // YouTube
  if (url.includes('youtube.com') || url.includes('youtu.be')) {
    video.platform = 'YouTube';
    video.title = document.querySelector('h1.ytd-watch-metadata yt-formatted-string, h1.ytd-video-primary-info-renderer')?.innerText || document.title;
    video.channel = document.querySelector('#channel-name a, .ytd-channel-name a')?.innerText || '';
    const thumb = document.querySelector('meta[property="og:image"]');
    if (thumb) video.thumbnail = thumb.getAttribute('content');
  }
  // Vimeo
  else if (url.includes('vimeo.com')) {
    video.platform = 'Vimeo';
    video.title = document.querySelector('h1')?.innerText || document.title;
    const thumb = document.querySelector('meta[property="og:image"]');
    if (thumb) video.thumbnail = thumb.getAttribute('content');
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
  } else if (request.action === 'extractContent') {
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
      const video = extractVideoInfo();
      data.title = video.title || document.title;
      data.content = `Platform: ${video.platform}\n${video.channel ? `Channel: ${video.channel}\n` : ''}`;
      data.metadata = {
        platform: video.platform,
        channel: video.channel,
        thumbnail: video.thumbnail,
      };
    } else {
      const page = extractPageContent();
      data.title = page.title;
      data.content = page.content || selectedText;
    }

    // Add selected text if available
    if (selectedText) {
      data.content = selectedText + '\n\n---\n\n' + data.content;
    }

    sendResponse(data);
  }
  
  return true; // Keep message channel open for async response
});

// Track text selection
document.addEventListener('mouseup', () => {
  selectedText = window.getSelection().toString().trim();
});
