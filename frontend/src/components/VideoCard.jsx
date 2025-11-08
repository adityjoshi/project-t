import { Link } from 'react-router-dom';

export default function VideoCard({ item }) {
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

  return (
    <Link
      to={`/items/${item.id}`}
      className="block bg-gradient-to-br from-red-50 to-pink-50 rounded-lg shadow-sm border border-red-200 hover:shadow-lg transition overflow-hidden"
    >
      <div className="relative">
        {item.image_url && (
          <div className="relative h-48 overflow-hidden">
            <img
              src={item.image_url}
              alt={item.title}
              className="w-full h-full object-cover"
              onError={(e) => {
                e.target.style.display = 'none';
              }}
            />
            <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-30">
              <div className="text-white text-5xl">▶️</div>
            </div>
          </div>
        )}
        {/* Only show thumbnail, not embed */}
      </div>
      <div className="p-6">
        <div className="flex items-center gap-2 mb-2">
          <span className="text-lg">🎥</span>
          <h3 className="text-xl font-semibold text-gray-900 line-clamp-2 flex-1">
            {item.title}
          </h3>
        </div>
        {item.summary && (
          <p className="text-gray-600 text-sm mb-4 line-clamp-2">
            {item.summary}
          </p>
        )}
        <div className="flex justify-between items-center text-xs text-gray-500">
          <span>{formatDate(item.created_at)}</span>
          {item.category && (
            <span className="px-2 py-1 bg-red-100 text-red-800 rounded-full text-xs font-medium">
              {item.category}
            </span>
          )}
        </div>
      </div>
    </Link>
  );
}

