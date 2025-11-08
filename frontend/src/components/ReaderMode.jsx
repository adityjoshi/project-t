import { useState } from 'react';

export default function ReaderMode({ content, title }) {
  const [fontSize, setFontSize] = useState('base');
  const [fontFamily, setFontFamily] = useState('sans');

  const fontSizes = {
    sm: 'text-sm',
    base: 'text-base',
    lg: 'text-lg',
    xl: 'text-xl',
  };

  const fontFamilies = {
    sans: 'font-sans',
    serif: 'font-serif',
    mono: 'font-mono',
  };

  return (
    <div className="bg-white rounded-lg shadow-lg border">
      {/* Reader Mode Controls */}
      <div className="border-b p-4 flex items-center justify-between bg-gray-50">
        <h2 className="text-lg font-semibold text-gray-700">Reader Mode</h2>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <label className="text-sm text-gray-600">Size:</label>
            <select
              value={fontSize}
              onChange={(e) => setFontSize(e.target.value)}
              className="border rounded px-2 py-1 text-sm"
            >
              <option value="sm">Small</option>
              <option value="base">Medium</option>
              <option value="lg">Large</option>
              <option value="xl">Extra Large</option>
            </select>
          </div>
          <div className="flex items-center gap-2">
            <label className="text-sm text-gray-600">Font:</label>
            <select
              value={fontFamily}
              onChange={(e) => setFontFamily(e.target.value)}
              className="border rounded px-2 py-1 text-sm"
            >
              <option value="sans">Sans</option>
              <option value="serif">Serif</option>
              <option value="mono">Mono</option>
            </select>
          </div>
        </div>
      </div>

      {/* Reader Content */}
      <div className="p-8 md:p-12 max-w-3xl mx-auto">
        <h1 className="text-4xl font-bold mb-6 text-gray-900">{title}</h1>
        <div
          className={`${fontSizes[fontSize]} ${fontFamilies[fontFamily]} leading-relaxed text-gray-800 prose prose-lg max-w-none`}
          style={{
            lineHeight: '1.8',
            wordSpacing: '0.05em',
          }}
        >
          {content.split('\n\n').map((paragraph, idx) => (
            <p key={idx} className="mb-6">
              {paragraph.trim()}
            </p>
          ))}
        </div>
      </div>
    </div>
  );
}

