import { Link } from 'react-router-dom';

export default function RecipeCard({ item }) {
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

  // Try to extract key info from content
  const extractRecipeInfo = (content) => {
    const prepMatch = content.match(/(?:prep|preparation)[:\s]+(\d+[\s\w]+)/i);
    const cookMatch = content.match(/(?:cook|baking)[:\s]+(\d+[\s\w]+)/i);
    const servingsMatch = content.match(/(?:serves|servings|yields)[:\s]+(\d+)/i);
    return {
      prepTime: prepMatch ? prepMatch[1] : null,
      cookTime: cookMatch ? cookMatch[1] : null,
      servings: servingsMatch ? servingsMatch[1] : null,
    };
  };

  const recipeInfo = extractRecipeInfo(item.content || '');

  return (
    <Link
      to={`/items/${item.id}`}
      className="block bg-gradient-to-br from-green-50 to-emerald-50 rounded-lg shadow-sm border border-green-200 hover:shadow-lg transition overflow-hidden"
    >
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
          <div className="absolute top-2 right-2 text-2xl">🍳</div>
        </div>
      )}
      <div className="p-6">
        <h3 className="text-xl font-bold text-gray-900 mb-2 line-clamp-2">
          {item.title}
        </h3>
        {(recipeInfo.prepTime || recipeInfo.cookTime || recipeInfo.servings) && (
          <div className="flex gap-4 mb-3 text-sm text-gray-600">
            {recipeInfo.prepTime && (
              <span>⏱️ Prep: {recipeInfo.prepTime}</span>
            )}
            {recipeInfo.cookTime && (
              <span>🔥 Cook: {recipeInfo.cookTime}</span>
            )}
            {recipeInfo.servings && (
              <span>👥 Serves: {recipeInfo.servings}</span>
            )}
          </div>
        )}
        {item.summary && (
          <p className="text-gray-700 text-sm mb-4 line-clamp-2">
            {item.summary}
          </p>
        )}
        <div className="flex items-center justify-between text-xs text-gray-600">
          <span>{formatDate(item.created_at)}</span>
          {item.category && (
            <span className="px-2 py-1 bg-green-100 text-green-800 rounded-full text-xs font-medium">
              {item.category}
            </span>
          )}
        </div>
      </div>
    </Link>
  );
}

