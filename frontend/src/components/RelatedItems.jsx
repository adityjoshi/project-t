import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { itemsAPI } from '../services/api';

export default function RelatedItems({ itemId }) {
  const [related, setRelated] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchRelated();
  }, [itemId]);

  const fetchRelated = async () => {
    setLoading(true);
    try {
      const response = await itemsAPI.getRelated(itemId);
      setRelated(response.data);
    } catch (error) {
      console.error('Failed to fetch related items:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="mt-8">
        <h2 className="text-2xl font-bold mb-4">Related Items</h2>
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  if (related.length === 0) {
    return null;
  }

  return (
    <div className="mt-8">
      <h2 className="text-2xl font-bold mb-4">Related Items</h2>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {related.map((rel) => (
          <Link
            key={rel.item.id}
            to={`/items/${rel.item.id}`}
            className="block bg-white rounded-lg shadow-sm border hover:shadow-md transition p-4"
          >
            <div className="flex justify-between items-start mb-2">
              <h3 className="font-semibold text-gray-900 line-clamp-2">
                {rel.item.title}
              </h3>
              <span className="text-xs text-gray-500 ml-2">
                {Math.round(rel.similarity_score * 100)}%
              </span>
            </div>
            {rel.item.summary && (
              <p className="text-sm text-gray-600 line-clamp-2">
                {rel.item.summary}
              </p>
            )}
            {rel.item.tags && rel.item.tags.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-2">
                {rel.item.tags.slice(0, 2).map((tag, idx) => (
                  <span
                    key={idx}
                    className="px-2 py-0.5 bg-gray-100 text-gray-600 text-xs rounded"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            )}
          </Link>
        ))}
      </div>
    </div>
  );
}

