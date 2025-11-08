import { Link } from 'react-router-dom';

export default function ProductCard({ item }) {
  const formatDate = (dateString) => {
    if (!dateString) return 'No date';
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return 'Invalid date';
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      });
    } catch (error) {
      return 'Invalid date';
    }
  };

  // Extract price from content or metadata
  const extractPrice = () => {
    if (item.content) {
      const priceMatch = item.content.match(/Price[:\s]+[₹$€£]?([\d,]+(?:\.\d{2})?)/i);
      if (priceMatch) {
        return priceMatch[0].replace(/Price[:\s]+/i, '');
      }
    }
    return null;
  };

  // Extract rating from content
  const extractRating = () => {
    if (item.content) {
      const ratingMatch = item.content.match(/Rating[:\s]+([\d.]+)/i);
      if (ratingMatch) {
        return parseFloat(ratingMatch[1]);
      }
    }
    return null;
  };

  const price = extractPrice();
  const rating = extractRating();

  return (
    <Link
      to={`/items/${item.id}`}
      className="block bg-gradient-to-br from-blue-50 to-indigo-50 rounded-lg shadow-sm border border-blue-200 hover:shadow-lg transition overflow-hidden"
    >
      {item.image_url && (
        <div className="relative h-48 overflow-hidden bg-white">
          <img
            src={item.image_url}
            alt={item.title}
            className="w-full h-full object-contain p-4"
            onError={(e) => {
              e.target.style.display = 'none';
            }}
          />
          <div className="absolute top-2 right-2 text-xl">🛍️</div>
        </div>
      )}
      <div className="p-6">
        <h3 className="text-lg font-bold text-gray-900 mb-2 line-clamp-2">
          {item.title}
        </h3>
        {price && (
          <div className="mb-3">
            <span className="text-2xl font-bold text-indigo-600">{price}</span>
          </div>
        )}
        {rating && (
          <div className="flex items-center gap-1 mb-3">
            <span className="text-yellow-500">⭐</span>
            <span className="text-sm font-semibold text-gray-700">{rating.toFixed(1)}</span>
          </div>
        )}
        {item.summary && (
          <p className="text-gray-600 text-sm mb-4 line-clamp-2">
            {item.summary}
          </p>
        )}
        <div className="flex items-center justify-between text-xs text-gray-500">
          <span>{formatDate(item.created_at)}</span>
          {item.category && (
            <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-xs font-medium">
              {item.category}
            </span>
          )}
        </div>
      </div>
    </Link>
  );
}

