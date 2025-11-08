import { Link } from 'react-router-dom';

export default function BookCard({ item }) {
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
      className="block bg-gradient-to-br from-amber-50 to-orange-50 rounded-lg shadow-sm border border-amber-200 hover:shadow-lg transition overflow-hidden"
    >
      <div className="flex">
        {item.image_url && (
          <div className="w-32 flex-shrink-0">
            <img
              src={item.image_url}
              alt={item.title}
              className="w-full h-full object-cover"
              onError={(e) => {
                e.target.style.display = 'none';
              }}
            />
          </div>
        )}
        <div className="flex-1 p-6">
          <div className="flex items-start justify-between mb-2">
            <h3 className="text-xl font-bold text-gray-900 line-clamp-2 flex-1">
              {item.title}
            </h3>
            <span className="ml-2 text-2xl">📚</span>
          </div>
          {item.summary && (
            <p className="text-gray-700 text-sm mb-4 line-clamp-2">
              {item.summary}
            </p>
          )}
          <div className="flex items-center justify-between text-xs text-gray-600">
            <span>{formatDate(item.created_at)}</span>
            {item.category && (
              <span className="px-2 py-1 bg-amber-100 text-amber-800 rounded-full text-xs font-medium">
                {item.category}
              </span>
            )}
          </div>
        </div>
      </div>
    </Link>
  );
}

