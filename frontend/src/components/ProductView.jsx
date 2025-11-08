export default function ProductView({ item }) {
  // Extract price from content
  const extractPrice = () => {
    if (item.content) {
      const priceMatch = item.content.match(/Price[:\s]+([₹$€£]?[\d,]+(?:\.\d{2})?)/i);
      if (priceMatch) {
        return priceMatch[1];
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

  // Extract other product details
  const extractDetails = () => {
    const details = {};
    if (item.content) {
      // Extract brand
      const brandMatch = item.content.match(/Brand[:\s]+([^\n]+)/i);
      if (brandMatch) details.brand = brandMatch[1].trim();
      
      // Extract availability
      const availMatch = item.content.match(/Availability[:\s]+([^\n]+)/i);
      if (availMatch) details.availability = availMatch[1].trim();
      
      // Extract features
      const featuresMatch = item.content.match(/Features?[:\s]+([^\n]+)/i);
      if (featuresMatch) details.features = featuresMatch[1].trim();
    }
    return details;
  };

  const price = extractPrice();
  const rating = extractRating();
  const details = extractDetails();

  return (
    <div className="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-lg border border-blue-200 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-start gap-3 mb-6">
          <span className="text-4xl">🛍️</span>
          <div className="flex-1">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">{item.title}</h1>
            {details.brand && (
              <p className="text-lg text-gray-600">Brand: {details.brand}</p>
            )}
          </div>
        </div>

        <div className="grid md:grid-cols-2 gap-8">
          {/* Product Image */}
          {item.image_url && (
            <div className="bg-white rounded-lg p-6 shadow-sm">
              <img
                src={item.image_url}
                alt={item.title}
                className="w-full h-auto rounded-lg"
                onError={(e) => {
                  e.target.style.display = 'none';
                }}
              />
            </div>
          )}

          {/* Product Details */}
          <div className="space-y-6">
            {price && (
              <div className="bg-white rounded-lg p-6 shadow-sm">
                <h2 className="text-sm font-semibold text-gray-600 mb-2">Price</h2>
                <p className="text-4xl font-bold text-indigo-600">{price}</p>
              </div>
            )}

            {rating && (
              <div className="bg-white rounded-lg p-6 shadow-sm">
                <h2 className="text-sm font-semibold text-gray-600 mb-2">Rating</h2>
                <div className="flex items-center gap-2">
                  <span className="text-3xl text-yellow-500">⭐</span>
                  <span className="text-2xl font-bold text-gray-800">{rating.toFixed(1)}</span>
                  <span className="text-gray-500">/ 5.0</span>
                </div>
              </div>
            )}

            {details.availability && (
              <div className="bg-white rounded-lg p-6 shadow-sm">
                <h2 className="text-sm font-semibold text-gray-600 mb-2">Availability</h2>
                <p className="text-lg text-gray-800">{details.availability}</p>
              </div>
            )}

            {item.source_url && (
              <div className="bg-white rounded-lg p-6 shadow-sm">
                <a
                  href={item.source_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center px-6 py-3 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition font-semibold"
                >
                  View Product →
                </a>
              </div>
            )}
          </div>
        </div>

        {/* Full Content */}
        {item.content && (
          <div className="mt-8 bg-white rounded-lg p-6 shadow-sm">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Product Details</h2>
            <div className="prose max-w-none">
              <p className="text-gray-700 whitespace-pre-wrap">{item.content}</p>
            </div>
          </div>
        )}

        {item.summary && (
          <div className="mt-6 bg-indigo-50 rounded-lg p-6">
            <h2 className="font-semibold text-indigo-900 mb-2">Summary</h2>
            <p className="text-indigo-800">{item.summary}</p>
          </div>
        )}
      </div>
    </div>
  );
}

