import { Link } from 'react-router-dom';
import ItemCard from './ItemCard';

export default function Dashboard({ items, loading, isSearch, onRefresh }) {
  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">
          {isSearch ? 'Search Results' : 'Your Knowledge Base'}
        </h1>
        {isSearch && (
          <button
            onClick={onRefresh}
            className="text-indigo-600 hover:text-indigo-700"
          >
            Clear Search
          </button>
        )}
      </div>

      {items.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg mb-4">
            {isSearch
              ? 'No results found'
              : "You haven't captured anything yet"}
          </p>
          {!isSearch && (
            <Link
              to="/capture"
              className="text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Capture your first thought →
            </Link>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {items.map((item) => (
            <ItemCard key={item.id} item={item} />
          ))}
        </div>
      )}
    </div>
  );
}

