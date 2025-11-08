import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import ItemCard from './ItemCard';

export default function SwipeableView({ items }) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [startX, setStartX] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
  const [offset, setOffset] = useState(0);
  const [hasSwiped, setHasSwiped] = useState(false);
  const cardRef = useRef(null);
  const navigate = useNavigate();

  const currentItem = items[currentIndex];

  // Handle touch start
  const handleTouchStart = (e) => {
    setStartX(e.touches[0].clientX);
    setIsDragging(true);
  };

  // Handle touch move
  const handleTouchMove = (e) => {
    if (!isDragging) return;
    const currentX = e.touches[0].clientX;
    const diff = currentX - startX;
    setOffset(diff);
  };

  // Handle touch end
  const handleTouchEnd = () => {
    if (!isDragging) return;
    setIsDragging(false);
    
    const threshold = 100; // Minimum swipe distance
    const movedDistance = Math.abs(offset);
    
    if (movedDistance > threshold) {
      // Mark that we swiped to prevent card click
      setHasSwiped(true);
      // Prevent card click if we swiped
      if (offset > threshold && currentIndex > 0) {
        // Swipe right - go to previous
        goToPrevious();
      } else if (offset < -threshold && currentIndex < items.length - 1) {
        // Swipe left - go to next
        goToNext();
      }
      // Reset offset and hasSwiped after a short delay
      setTimeout(() => {
        setOffset(0);
        setHasSwiped(false);
      }, 100);
    } else {
      // Small movement - treat as click, not swipe
      setOffset(0);
      setHasSwiped(false);
    }
  };

  // Handle mouse drag
  const handleMouseDown = (e) => {
    // Don't start dragging if clicking on a link or button
    if (e.target.closest('a') || e.target.closest('button')) {
      return;
    }
    setStartX(e.clientX);
    setIsDragging(true);
  };

  const handleMouseMove = (e) => {
    if (!isDragging) return;
    const diff = e.clientX - startX;
    setOffset(diff);
  };

  const handleMouseUp = (e) => {
    if (!isDragging) return;
    setIsDragging(false);
    
    const threshold = 100;
    const movedDistance = Math.abs(offset);
    
    if (movedDistance > threshold) {
      // Mark that we swiped to prevent card click
      setHasSwiped(true);
      if (offset > threshold && currentIndex > 0) {
        goToPrevious();
      } else if (offset < -threshold && currentIndex < items.length - 1) {
        goToNext();
      }
      // Reset offset and hasSwiped after a short delay
      setTimeout(() => {
        setOffset(0);
        setHasSwiped(false);
      }, 100);
    } else {
      // Small movement - treat as click, not swipe
      setOffset(0);
      setHasSwiped(false);
    }
  };
  

  // Navigation functions
  const goToNext = () => {
    if (currentIndex < items.length - 1) {
      setCurrentIndex(currentIndex + 1);
    }
  };

  const goToPrevious = () => {
    if (currentIndex > 0) {
      setCurrentIndex(currentIndex - 1);
    }
  };

  // Keyboard navigation
  useEffect(() => {
    const handleKeyPress = (e) => {
      if (e.key === 'ArrowLeft') {
        goToPrevious();
      } else if (e.key === 'ArrowRight') {
        goToNext();
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [currentIndex, items.length]);

  // Reset to first item when items change
  useEffect(() => {
    setCurrentIndex(0);
  }, [items.length]);

  if (items.length === 0) {
    return (
      <div className="flex items-center justify-center h-96">
        <p className="text-gray-500">No items to display</p>
      </div>
    );
  }

  return (
    <div className="relative w-full max-w-6xl mx-auto">
      {/* Progress indicator */}
      <div className="mb-6">
        <div className="flex items-center justify-between text-sm text-gray-600 mb-2">
          <span className="font-medium">{currentIndex + 1} of {items.length}</span>
          <span className="text-gray-500">{Math.round(((currentIndex + 1) / items.length) * 100)}%</span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
          <div
            className="bg-gradient-to-r from-indigo-500 to-indigo-600 h-2 rounded-full transition-all duration-300 ease-out"
            style={{ width: `${((currentIndex + 1) / items.length) * 100}%` }}
          />
        </div>
      </div>

      {/* Main content area with side buttons */}
      <div className="relative flex items-center gap-4">
        {/* Left navigation button */}
        <button
          onClick={goToPrevious}
          disabled={currentIndex === 0}
          className={`flex-shrink-0 w-14 h-14 rounded-full flex items-center justify-center transition-all duration-200 shadow-lg ${
            currentIndex === 0
              ? 'bg-gray-100 text-gray-300 cursor-not-allowed opacity-50'
              : 'bg-white text-indigo-600 hover:bg-indigo-50 hover:scale-110 border-2 border-indigo-200'
          }`}
          aria-label="Previous item"
        >
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 19l-7-7 7-7"
            />
          </svg>
        </button>

        {/* Card container */}
        <div
          className="flex-1 relative h-[650px]"
          onTouchStart={handleTouchStart}
          onTouchMove={handleTouchMove}
          onTouchEnd={handleTouchEnd}
          onMouseDown={handleMouseDown}
          onMouseMove={handleMouseMove}
          onMouseUp={handleMouseUp}
          onMouseLeave={handleMouseUp}
        >
          <div
            ref={cardRef}
            className="absolute inset-0 transition-transform duration-300 ease-out"
            style={{
              transform: `translateX(${offset}px)`,
              cursor: isDragging ? 'grabbing' : 'pointer',
            }}
          >
            <div 
              className="h-full flex items-center justify-center" 
              style={{ pointerEvents: isDragging ? 'none' : 'auto' }}
            >
              <div className="w-full max-w-2xl transform scale-105 transition-transform duration-200 hover:scale-[1.08]">
                {/* Render card - Link components will handle navigation */}
                <ItemCard item={currentItem} />
              </div>
            </div>
          </div>
        </div>

        {/* Right navigation button */}
        <button
          onClick={goToNext}
          disabled={currentIndex === items.length - 1}
          className={`flex-shrink-0 w-14 h-14 rounded-full flex items-center justify-center transition-all duration-200 shadow-lg ${
            currentIndex === items.length - 1
              ? 'bg-gray-100 text-gray-300 cursor-not-allowed opacity-50'
              : 'bg-white text-indigo-600 hover:bg-indigo-50 hover:scale-110 border-2 border-indigo-200'
          }`}
          aria-label="Next item"
        >
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 5l7 7-7 7"
            />
          </svg>
        </button>
      </div>

      {/* Dot navigation */}
      <div className="flex items-center justify-center gap-2 mt-6">
        {items.map((_, idx) => (
          <button
            key={idx}
            onClick={() => setCurrentIndex(idx)}
            className={`rounded-full transition-all duration-200 ${
              idx === currentIndex
                ? 'w-8 h-2 bg-indigo-600'
                : 'w-2 h-2 bg-gray-300 hover:bg-gray-400'
            }`}
            aria-label={`Go to item ${idx + 1}`}
          />
        ))}
      </div>

      {/* Swipe hint */}
      {items.length > 1 && (
        <div className="mt-4 text-center text-sm text-gray-500">
          <p className="flex items-center justify-center gap-2">
            <span>👆</span>
            <span>Swipe left/right or use arrow keys to navigate</span>
          </p>
        </div>
      )}
    </div>
  );
}

