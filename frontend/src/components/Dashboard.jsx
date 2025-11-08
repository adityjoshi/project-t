import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import ItemCard from './ItemCard';
import SwipeableView from './SwipeableView';

export default function Dashboard({ items, loading, isSearch, onRefresh }) {
  const [selectedCategory, setSelectedCategory] = useState('');
  const [viewMode, setViewMode] = useState('grid'); // 'grid' or 'swipe'

  // Get unique categories from items
  const categories = useMemo(() => {
    const cats = new Set();
    items.forEach(item => {
      if (item.category) {
        cats.add(item.category);
      }
    });
    return Array.from(cats).sort();
  }, [items]);

  // Filter items by category
  const filteredItems = useMemo(() => {
    if (!selectedCategory) return items;
    return items.filter(item => item.category === selectedCategory);
  }, [items, selectedCategory]);

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
        <div className="flex items-center gap-4">
          {categories.length > 0 && !isSearch && (
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">All Categories</option>
              {categories.map((category) => (
                <option key={category} value={category}>
                  {category}
                </option>
              ))}
            </select>
          )}
          {filteredItems.length > 0 && (
            <div className="flex items-center gap-2 border border-gray-300 rounded-lg overflow-hidden">
              <button
                onClick={() => setViewMode('grid')}
                className={`px-4 py-2 transition ${
                  viewMode === 'grid'
                    ? 'bg-indigo-600 text-white'
                    : 'bg-white text-gray-700 hover:bg-gray-50'
                }`}
                title="Grid View"
              >
                ⊞
              </button>
              <button
                onClick={() => setViewMode('swipe')}
                className={`px-4 py-2 transition ${
                  viewMode === 'swipe'
                    ? 'bg-indigo-600 text-white'
                    : 'bg-white text-gray-700 hover:bg-gray-50'
                }`}
                title="Swipe View"
              >
                ⇄
              </button>
            </div>
          )}
          {isSearch && (
            <button
              onClick={onRefresh}
              className="text-indigo-600 hover:text-indigo-700"
            >
              Clear Search
            </button>
          )}
        </div>
      </div>

      {filteredItems.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg mb-4">
            {isSearch
              ? 'No results found'
              : selectedCategory
              ? `No items in "${selectedCategory}" category`
              : "You haven't captured anything yet"}
          </p>
          {!isSearch && !selectedCategory && (
            <Link
              to="/capture"
              className="text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Capture your first thought →
            </Link>
          )}
          {selectedCategory && (
            <button
              onClick={() => setSelectedCategory('')}
              className="text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Clear filter
            </button>
          )}
        </div>
      ) : viewMode === 'swipe' ? (
        <SwipeableView items={filteredItems} />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredItems.map((item) => (
            <ItemCard key={item.id} item={item} />
          ))}
        </div>
      )}
    </div>
  );
}

