import { Link } from 'react-router-dom';
import BookCard from './BookCard';
import RecipeCard from './RecipeCard';
import ArticleCard from './ArticleCard';
import VideoCard from './VideoCard';
import ProductCard from './ProductCard';
import TodoCard from './TodoCard';

export default function ItemCard({ item }) {
  // Check if it's a YouTube video (even if type is url/image)
  const isYouTubeVideo = item.source_url && (
    item.source_url.includes('youtube.com') || 
    item.source_url.includes('youtu.be')
  );

  // Route to specialized card components based on type
  switch (item.type) {
    case 'book':
      return <BookCard item={item} />;
    case 'recipe':
      return <RecipeCard item={item} />;
    case 'blog':
    case 'article':
      return <ArticleCard item={item} />;
    case 'video':
      return <VideoCard item={item} />;
    case 'url':
      // If URL is a YouTube video, show as video
      if (isYouTubeVideo) {
        return <VideoCard item={item} />;
      }
      return <ArticleCard item={item} />;
    case 'amazon':
      return <ProductCard item={item} />;
    case 'text':
      // Check if it's a todo list
      if (item.title?.toLowerCase().includes('todo') || 
          item.title?.toLowerCase().includes('to-do') ||
          item.content?.match(/^[-*•]\s|^\d+\.\s|^\[[\sx]\]/im)) {
        return <TodoCard item={item} />;
      }
      return <DefaultItemCard item={item} />;
    case 'image':
      // If image is a screenshot of a YouTube video, show as video
      if (isYouTubeVideo) {
        return <VideoCard item={item} />;
      }
      return <DefaultItemCard item={item} />;
    default:
      // Default card for other types
      return <DefaultItemCard item={item} />;
  }
}

function DefaultItemCard({ item }) {
  const formatDate = (dateString) => {
    if (!dateString) return 'No date';
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) {
        return 'Invalid date';
      }
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      });
    } catch (error) {
      return 'Invalid date';
    }
  };

  return (
    <Link
      to={`/items/${item.id}`}
      className="block bg-white rounded-lg shadow-sm border hover:shadow-md transition overflow-hidden"
    >
      {item.image_url && (
        <img
          src={item.image_url}
          alt={item.title}
          className="w-full h-48 object-cover"
          onError={(e) => {
            e.target.style.display = 'none';
          }}
        />
      )}
      <div className="p-6">
        <h3 className="text-xl font-semibold mb-2 text-gray-900 line-clamp-2">
          {item.title}
        </h3>
      {item.summary && (
        <p className="text-gray-600 text-sm mb-4 line-clamp-3">
          {item.summary}
        </p>
      )}
      {item.tags && item.tags.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-4">
          {item.tags.slice(0, 3).map((tag, idx) => (
            <span
              key={idx}
              className="px-2 py-1 bg-indigo-100 text-indigo-700 text-xs rounded-full"
            >
              {tag}
            </span>
          ))}
          {item.tags.length > 3 && (
            <span className="px-2 py-1 text-gray-500 text-xs">
              +{item.tags.length - 3}
            </span>
          )}
        </div>
      )}
            <div className="flex justify-between items-center text-xs text-gray-500">
              <span>{formatDate(item.created_at)}</span>
              <div className="flex items-center gap-2">
                {item.category && (
                  <span className="px-2 py-1 bg-purple-100 text-purple-700 rounded-full text-xs font-medium">
                    {item.category}
                  </span>
                )}
                <span className="capitalize">{item.type}</span>
              </div>
            </div>
      </div>
    </Link>
  );
}

