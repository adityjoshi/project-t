import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { itemsAPI } from '../services/api';
import RelatedItems from './RelatedItems';
import ReaderMode from './ReaderMode';
import RecipeView from './RecipeView';
import TodoView from './TodoView';
import ProductView from './ProductView';

export default function ItemDetail({ onDelete }) {
  const { id } = useParams();
  const navigate = useNavigate();
  const [item, setItem] = useState(null);
  const [loading, setLoading] = useState(true);
  const [deleting, setDeleting] = useState(false);
  const [readerMode, setReaderMode] = useState(false);

  useEffect(() => {
    fetchItem();
  }, [id]);

  const fetchItem = async () => {
    setLoading(true);
    try {
      const response = await itemsAPI.getById(id);
      setItem(response.data);
    } catch (error) {
      console.error('Failed to fetch item:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this item?')) {
      return;
    }

    setDeleting(true);
    try {
      await itemsAPI.delete(id);
      if (onDelete) onDelete();
      navigate('/');
    } catch (error) {
      console.error('Failed to delete item:', error);
      alert('Failed to delete item');
    } finally {
      setDeleting(false);
    }
  };


  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  if (!item) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">Item not found</p>
        <Link to="/" className="text-indigo-600 hover:text-indigo-700 mt-4 inline-block">
          ← Back to Dashboard
        </Link>
      </div>
    );
  }

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  // Check if it's a YouTube video (even if type is url/image)
  const isYouTubeVideo = item.source_url && (
    item.source_url.includes('youtube.com') || 
    item.source_url.includes('youtu.be')
  );

  // Check if it's a todo list
  const isTodo = item.type === 'text' && (
    item.title?.toLowerCase().includes('todo') || 
    item.title?.toLowerCase().includes('to-do') ||
    item.content?.match(/^[-*•]\s|^\d+\.\s|^\[[\sx]\]/im)
  );

  // Show specialized views for certain content types
  if (item.type === 'recipe') {
    return (
      <div className="max-w-6xl mx-auto">
        <Link
          to="/"
          className="text-indigo-600 hover:text-indigo-700 mb-4 inline-block"
        >
          ← Back to Dashboard
        </Link>
        <div className="mb-4 flex justify-between items-center">
          <button
            onClick={handleDelete}
            disabled={deleting}
            className="px-4 py-2 text-red-600 border border-red-300 rounded-lg hover:bg-red-50 disabled:opacity-50"
          >
            {deleting ? 'Deleting...' : 'Delete'}
          </button>
        </div>
        <RecipeView item={item} />
        <RelatedItems itemId={id} />
      </div>
    );
  }

  if (item.type === 'amazon') {
    return (
      <div className="max-w-6xl mx-auto">
        <Link
          to="/"
          className="text-indigo-600 hover:text-indigo-700 mb-4 inline-block"
        >
          ← Back to Dashboard
        </Link>
        <div className="mb-4 flex justify-between items-center">
          <button
            onClick={handleDelete}
            disabled={deleting}
            className="px-4 py-2 text-red-600 border border-red-300 rounded-lg hover:bg-red-50 disabled:opacity-50"
          >
            {deleting ? 'Deleting...' : 'Delete'}
          </button>
        </div>
        <ProductView item={item} />
        <RelatedItems itemId={id} />
      </div>
    );
  }

  if (isTodo) {
    return (
      <div className="max-w-6xl mx-auto">
        <Link
          to="/"
          className="text-indigo-600 hover:text-indigo-700 mb-4 inline-block"
        >
          ← Back to Dashboard
        </Link>
        <div className="mb-4 flex justify-between items-center">
          <button
            onClick={handleDelete}
            disabled={deleting}
            className="px-4 py-2 text-red-600 border border-red-300 rounded-lg hover:bg-red-50 disabled:opacity-50"
          >
            {deleting ? 'Deleting...' : 'Delete'}
          </button>
        </div>
        <TodoView item={item} />
        <RelatedItems itemId={id} />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <Link
        to="/"
        className="text-indigo-600 hover:text-indigo-700 mb-4 inline-block"
      >
        ← Back to Dashboard
      </Link>

      <div className="bg-white rounded-lg shadow-sm border p-8">
        <div className="flex justify-between items-start mb-6">
          <div className="flex-1">
            <h1 className="text-3xl font-bold mb-2">{item.title}</h1>
                <div className="flex items-center space-x-4 text-sm text-gray-500">
                  <span>{formatDate(item.created_at)}</span>
                  {item.category && (
                    <span className="px-2 py-1 bg-purple-100 text-purple-700 rounded-full text-xs font-medium">
                      {item.category}
                    </span>
                  )}
                  <span className="capitalize">• {item.type}</span>
                </div>
          </div>
          <button
            onClick={handleDelete}
            disabled={deleting}
            className="px-4 py-2 text-red-600 border border-red-300 rounded-lg hover:bg-red-50 disabled:opacity-50"
          >
            {deleting ? 'Deleting...' : 'Delete'}
          </button>
        </div>

        {item.source_url && (
          <div className="mb-6">
            <a
              href={item.source_url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-indigo-600 hover:text-indigo-700 break-all"
            >
              {item.source_url}
            </a>
          </div>
        )}

        {/* Show embedded video player for YouTube videos */}
        {(item.embed_html || isYouTubeVideo) && (
          <div className="mb-6">
            {item.embed_html ? (
              <div className="rounded-lg overflow-hidden shadow-lg">
                <div className="relative pb-[56.25%] h-0 overflow-hidden bg-black">
                  <div
                    className="absolute top-0 left-0 w-full h-full"
                    dangerouslySetInnerHTML={{ __html: item.embed_html }}
                  />
                </div>
              </div>
            ) : isYouTubeVideo ? (
              // Generate embed for YouTube if not already present
              (() => {
                const videoId = item.source_url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]+)/)?.[1];
                if (videoId) {
                  return (
                    <div className="rounded-lg overflow-hidden shadow-lg">
                      <div className="relative pb-[56.25%] h-0 overflow-hidden bg-black">
                        <iframe
                          className="absolute top-0 left-0 w-full h-full"
                          src={`https://www.youtube.com/embed/${videoId}?rel=0`}
                          frameBorder="0"
                          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                          allowFullScreen
                          title={item.title}
                        />
                      </div>
                    </div>
                  );
                }
                return null;
              })()
            ) : null}
          </div>
        )}

        {item.image_url && !item.embed_html && !isYouTubeVideo && (
          <div className="mb-6">
            <img
              src={item.image_url}
              alt={item.title}
              className="w-full max-w-2xl rounded-lg shadow-md"
              onError={(e) => {
                e.target.style.display = 'none';
              }}
            />
          </div>
        )}

        {item.summary && (
          <div className="mb-6 p-4 bg-indigo-50 rounded-lg">
            <h2 className="font-semibold text-indigo-900 mb-2">Summary</h2>
            <p className="text-indigo-800">{item.summary}</p>
          </div>
        )}

        {item.tags && item.tags.length > 0 && (
          <div className="mb-6">
            <h2 className="font-semibold text-gray-700 mb-2">Tags</h2>
            <div className="flex flex-wrap gap-2">
              {item.tags.map((tag, idx) => (
                <span
                  key={idx}
                  className="px-3 py-1 bg-indigo-100 text-indigo-700 text-sm rounded-full"
                >
                  {tag}
                </span>
              ))}
            </div>
          </div>
        )}

        {/* Reader Mode Toggle for Articles (not for videos) */}
        {(item.type === 'blog' || item.type === 'url' || item.type === 'article') && item.type !== 'video' && !isYouTubeVideo && (
          <div className="mb-4 flex justify-end">
            <button
              onClick={() => setReaderMode(!readerMode)}
              className="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition"
            >
              {readerMode ? '📄 Standard View' : '📖 Reader Mode'}
            </button>
          </div>
        )}

        {readerMode && (item.type === 'blog' || item.type === 'url' || item.type === 'article') && item.type !== 'video' && !isYouTubeVideo ? (
          <ReaderMode content={item.content} title={item.title} />
        ) : (
          <div className="mb-6">
            <h2 className="font-semibold text-gray-700 mb-2">Content</h2>
            <div className="prose max-w-none">
              <p className="text-gray-700 whitespace-pre-wrap">{item.content}</p>
            </div>
          </div>
        )}
      </div>

      <RelatedItems itemId={id} />
    </div>
  );
}

