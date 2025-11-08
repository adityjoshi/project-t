import { Link } from 'react-router-dom';

export default function TodoCard({ item }) {
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

  // Parse todo items from content
  const parseTodos = (content) => {
    if (!content) return [];
    const lines = content.split('\n').map(l => l.trim()).filter(l => l);
    const todos = [];
    
    for (const line of lines) {
      // Match various todo formats: "- item", "* item", "• item", "1. item", "[ ] item", "- [ ] item"
      const todoMatch = line.match(/^[-*•]\s*(?:\[[\sx]\]\s*)?(.+)$|^\d+\.\s*(?:\[[\sx]\]\s*)?(.+)$|^\[[\sx]\]\s*(.+)$/i);
      if (todoMatch) {
        const todoText = todoMatch[1] || todoMatch[2] || todoMatch[3];
        const isCompleted = line.toLowerCase().includes('[x]') || line.toLowerCase().includes('✓') || line.toLowerCase().includes('done');
        todos.push({ text: todoText, completed: isCompleted });
      }
    }
    
    // If no structured todos found, treat each line as a todo
    if (todos.length === 0 && lines.length > 0) {
      return lines.slice(0, 5).map(line => ({ text: line, completed: false }));
    }
    
    return todos.slice(0, 5); // Show max 5 todos in card
  };

  const todos = parseTodos(item.content);
  const completedCount = todos.filter(t => t.completed).length;
  const totalCount = todos.length;

  return (
    <Link
      to={`/items/${item.id}`}
      className="block bg-gradient-to-br from-yellow-50 to-orange-50 rounded-lg shadow-sm border border-yellow-200 hover:shadow-lg transition overflow-hidden"
    >
      <div className="p-6">
        <div className="flex items-center gap-2 mb-4">
          <span className="text-2xl">✅</span>
          <h3 className="text-xl font-bold text-gray-900 line-clamp-1">
            {item.title || 'To-Do List'}
          </h3>
        </div>
        
        {todos.length > 0 ? (
          <>
            <div className="space-y-2 mb-4">
              {todos.map((todo, idx) => (
                <div key={idx} className="flex items-start gap-2">
                  <span className={`mt-1 ${todo.completed ? 'text-green-600' : 'text-gray-400'}`}>
                    {todo.completed ? '✓' : '○'}
                  </span>
                  <span className={`text-sm flex-1 ${todo.completed ? 'line-through text-gray-500' : 'text-gray-700'}`}>
                    {todo.text}
                  </span>
                </div>
              ))}
            </div>
            {totalCount > 0 && (
              <div className="mb-4">
                <div className="flex items-center justify-between text-xs text-gray-600 mb-1">
                  <span>Progress</span>
                  <span>{completedCount}/{totalCount} completed</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-green-500 h-2 rounded-full transition-all"
                    style={{ width: `${(completedCount / totalCount) * 100}%` }}
                  />
                </div>
              </div>
            )}
          </>
        ) : (
          <p className="text-gray-600 text-sm mb-4 line-clamp-3">
            {item.summary || item.content}
          </p>
        )}
        
        <div className="flex items-center justify-between text-xs text-gray-500">
          <span>{formatDate(item.created_at)}</span>
          {item.category && (
            <span className="px-2 py-1 bg-yellow-100 text-yellow-800 rounded-full text-xs font-medium">
              {item.category}
            </span>
          )}
        </div>
      </div>
    </Link>
  );
}

