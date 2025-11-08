import { Link } from 'react-router-dom';

export default function ItemCard({ item }) {
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
          <span className="capitalize">{item.type}</span>
        </div>
      </div>
    </Link>
  );
}

