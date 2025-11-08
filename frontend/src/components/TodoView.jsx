export default function TodoView({ item }) {
  // Parse todo items from content
  const parseTodos = (content) => {
    if (!content) return [];
    const lines = content.split('\n').map(l => l.trim()).filter(l => l);
    const todos = [];
    
    for (const line of lines) {
      // Match various todo formats
      const todoMatch = line.match(/^[-*•]\s*(?:\[[\sx]\]\s*)?(.+)$|^\d+\.\s*(?:\[[\sx]\]\s*)?(.+)$|^\[[\sx]\]\s*(.+)$/i);
      if (todoMatch) {
        const todoText = todoMatch[1] || todoMatch[2] || todoMatch[3];
        const isCompleted = line.toLowerCase().includes('[x]') || line.toLowerCase().includes('✓') || line.toLowerCase().includes('done');
        todos.push({ text: todoText, completed: isCompleted });
      }
    }
    
    // If no structured todos found, treat each non-empty line as a todo
    if (todos.length === 0 && lines.length > 0) {
      return lines.map(line => ({ text: line, completed: false }));
    }
    
    return todos;
  };

  const todos = parseTodos(item.content);
  const completedCount = todos.filter(t => t.completed).length;
  const totalCount = todos.length;
  const progressPercentage = totalCount > 0 ? (completedCount / totalCount) * 100 : 0;

  return (
    <div className="bg-gradient-to-br from-yellow-50 to-orange-50 rounded-lg border border-yellow-200 p-8">
      <div className="max-w-3xl mx-auto">
        <div className="flex items-center gap-3 mb-6">
          <span className="text-4xl">✅</span>
          <h1 className="text-3xl font-bold text-gray-900">{item.title || 'To-Do List'}</h1>
        </div>

        {totalCount > 0 && (
          <div className="mb-8 bg-white rounded-lg p-6 shadow-sm">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">Progress</h2>
              <span className="text-lg font-bold text-indigo-600">
                {completedCount} / {totalCount}
              </span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-4 mb-2">
              <div
                className="bg-gradient-to-r from-green-500 to-green-600 h-4 rounded-full transition-all duration-300"
                style={{ width: `${progressPercentage}%` }}
              />
            </div>
            <p className="text-sm text-gray-600">
              {progressPercentage.toFixed(0)}% completed
            </p>
          </div>
        )}

        <div className="bg-white rounded-lg p-6 shadow-sm">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Tasks</h2>
          {todos.length > 0 ? (
            <div className="space-y-3">
              {todos.map((todo, idx) => (
                <div
                  key={idx}
                  className={`flex items-start gap-3 p-3 rounded-lg border-2 transition ${
                    todo.completed
                      ? 'bg-green-50 border-green-200'
                      : 'bg-white border-gray-200 hover:border-yellow-300'
                  }`}
                >
                  <span
                    className={`text-2xl mt-0.5 ${
                      todo.completed ? 'text-green-600' : 'text-gray-400'
                    }`}
                  >
                    {todo.completed ? '✓' : '○'}
                  </span>
                  <span
                    className={`text-base flex-1 ${
                      todo.completed
                        ? 'line-through text-gray-500'
                        : 'text-gray-800'
                    }`}
                  >
                    {todo.text}
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-gray-500 italic">
              <p>No structured todo items found. Displaying raw content:</p>
              <div className="mt-4 prose max-w-none">
                <p className="text-gray-700 whitespace-pre-wrap">{item.content}</p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

