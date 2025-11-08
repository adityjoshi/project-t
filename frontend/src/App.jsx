import { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import Dashboard from './components/Dashboard';
import CaptureForm from './components/CaptureForm';
import ItemDetail from './components/ItemDetail';
import SearchBar from './components/SearchBar';
import { itemsAPI, searchAPI } from './services/api';

function App() {
  const [items, setItems] = useState([]);
  const [searchResults, setSearchResults] = useState(null);
  const [loading, setLoading] = useState(false);

  const refreshItems = async () => {
    setLoading(true);
    try {
      const response = await itemsAPI.getAll();
      setItems(response.data || []);
    } catch (error) {
      console.error('Failed to fetch items:', error);
      setItems([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    refreshItems();
  }, []);

  const handleSearch = async (query) => {
    if (!query.trim()) {
      setSearchResults(null);
      return;
    }
    setLoading(true);
    try {
      const response = await searchAPI.search(query);
      // Search results come as SearchResult objects with {item, similarity_score}
      // Extract just the items for display
      const results = (response.data || []).map(result => 
        result.item || result // Handle both formats
      );
      setSearchResults(results);
    } catch (error) {
      console.error('Search failed:', error);
      setSearchResults([]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <nav className="bg-white shadow-sm border-b">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between h-16">
              <div className="flex items-center">
                <Link to="/" className="text-2xl font-bold text-indigo-600">
                  🧠 Synapse
                </Link>
              </div>
              <div className="flex items-center space-x-4">
                <SearchBar onSearch={handleSearch} />
                <Link
                  to="/capture"
                  className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition"
                >
                  + Capture
                </Link>
              </div>
            </div>
          </div>
        </nav>

        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <Routes>
            <Route
              path="/"
              element={
                <Dashboard
                  items={searchResults || items}
                  loading={loading}
                  isSearch={searchResults !== null}
                  onRefresh={refreshItems}
                />
              }
            />
            <Route
              path="/capture"
              element={<CaptureForm onSuccess={refreshItems} />}
            />
            <Route
              path="/items/:id"
              element={<ItemDetail onDelete={refreshItems} />}
            />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;

