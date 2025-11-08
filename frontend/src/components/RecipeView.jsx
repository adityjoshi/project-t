export default function RecipeView({ item }) {
  // Parse recipe content to extract ingredients and instructions
  const parseRecipe = (content) => {
    const lines = content.split('\n').map(l => l.trim()).filter(l => l);
    
    let ingredients = [];
    let instructions = [];
    let currentSection = null;
    
    for (const line of lines) {
      const lower = line.toLowerCase();
      
      // Detect section headers
      if (lower.includes('ingredient') || lower === 'ingredients:') {
        currentSection = 'ingredients';
        continue;
      }
      if (lower.includes('instruction') || lower.includes('direction') || lower === 'instructions:' || lower === 'directions:') {
        currentSection = 'instructions';
        continue;
      }
      if (lower.includes('method') || lower === 'method:') {
        currentSection = 'instructions';
        continue;
      }
      
      // Add to appropriate section
      if (currentSection === 'ingredients') {
        if (line.match(/^\d+/) || line.match(/^[•\-\*]/) || line.match(/^[a-z]/i)) {
          ingredients.push(line.replace(/^[•\-\*]\s*/, ''));
        }
      } else if (currentSection === 'instructions') {
        if (line.match(/^\d+\./) || line.match(/^[•\-\*]/) || line.match(/^[A-Z]/)) {
          instructions.push(line.replace(/^\d+\.\s*/, '').replace(/^[•\-\*]\s*/, ''));
        }
      } else {
        // Try to auto-detect
        if (line.match(/^\d+\s*(cup|tbsp|tsp|oz|lb|g|kg|ml|l)/i)) {
          ingredients.push(line);
        } else if (line.match(/^\d+\./) || line.length > 50) {
          instructions.push(line.replace(/^\d+\.\s*/, ''));
        }
      }
    }
    
    return { ingredients, instructions };
  };

  const { ingredients, instructions } = parseRecipe(item.content || '');

  return (
    <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-lg border border-green-200 p-8">
      <div className="max-w-3xl mx-auto">
        <div className="flex items-center gap-3 mb-6">
          <span className="text-4xl">🍳</span>
          <h1 className="text-3xl font-bold text-gray-900">{item.title}</h1>
        </div>

        {item.image_url && (
          <div className="mb-8">
            <img
              src={item.image_url}
              alt={item.title}
              className="w-full max-w-2xl rounded-lg shadow-lg"
              onError={(e) => {
                e.target.style.display = 'none';
              }}
            />
          </div>
        )}

        <div className="grid md:grid-cols-2 gap-8">
          {/* Ingredients */}
          <div className="bg-white rounded-lg p-6 shadow-sm">
            <h2 className="text-2xl font-bold mb-4 text-gray-900 flex items-center gap-2">
              <span>🥘</span> Ingredients
            </h2>
            {ingredients.length > 0 ? (
              <ul className="space-y-2">
                {ingredients.map((ingredient, idx) => (
                  <li key={idx} className="flex items-start gap-2 text-gray-700">
                    <span className="text-green-600 mt-1">•</span>
                    <span>{ingredient}</span>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-gray-500 italic">No ingredients detected. Check the full content below.</p>
            )}
          </div>

          {/* Instructions */}
          <div className="bg-white rounded-lg p-6 shadow-sm">
            <h2 className="text-2xl font-bold mb-4 text-gray-900 flex items-center gap-2">
              <span>📝</span> Instructions
            </h2>
            {instructions.length > 0 ? (
              <ol className="space-y-4">
                {instructions.map((instruction, idx) => (
                  <li key={idx} className="flex gap-3 text-gray-700">
                    <span className="flex-shrink-0 w-8 h-8 bg-green-100 text-green-700 rounded-full flex items-center justify-center font-bold">
                      {idx + 1}
                    </span>
                    <span className="flex-1">{instruction}</span>
                  </li>
                ))}
              </ol>
            ) : (
              <p className="text-gray-500 italic">No instructions detected. Check the full content below.</p>
            )}
          </div>
        </div>

        {/* Full Content (fallback) */}
        {ingredients.length === 0 && instructions.length === 0 && (
          <div className="mt-8 bg-white rounded-lg p-6 shadow-sm">
            <h2 className="text-xl font-bold mb-4 text-gray-900">Full Recipe</h2>
            <div className="prose max-w-none">
              <p className="text-gray-700 whitespace-pre-wrap">{item.content}</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

