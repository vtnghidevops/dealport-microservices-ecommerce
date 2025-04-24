// components/admin/product/add/components/ProductTags.tsx
import React, { useState } from "react";

interface ProductTagsProps {
  tags: string[];
  onChange: (tags: string[]) => void;
}

export const ProductTags: React.FC<ProductTagsProps> = ({
  tags,
  onChange,
}) => {
  const [tagInput, setTagInput] = useState<string>("");

  const addTag = () => {
    if (tagInput.trim()) {
      onChange([...tags, tagInput.trim()]);
      setTagInput("");
    }
  };

  const removeTag = (index: number) => {
    const updatedTags = [...tags];
    updatedTags.splice(index, 1);
    onChange(updatedTags);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      addTag();
    }
  };

  return (
    <div className="bg-white rounded-lg p-[1.25rem] drop-shadow-sm filter mb-6 mt-5">
      <h2 className="font-bold text-cyprus text-[22px] mb-[15px]">
        Product Tags
      </h2>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Add Tags
        </label>
        <div className="flex">
          <input
            type="text"
            value={tagInput}
            onChange={(e) => setTagInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Enter a tag and press Enter"
            className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-l-lg p-2 focus:outline-none focus:border-ocean-green"
          />
          <button
            onClick={addTag}
            className="bg-ocean-green text-white px-4 rounded-r-lg hover:bg-opacity-90"
            type="button"
          >
            Add
          </button>
        </div>
        <p className="text-sm text-gray-500 mt-1">Press Enter to add multiple tags</p>
      </div>

      {tags.length > 0 && (
        <div className="flex flex-wrap gap-2 mt-3">
          {tags.map((tag, index) => (
            <div 
              key={index} 
              className="bg-neutral-100 text-cyprus px-3 py-1 rounded-full flex items-center"
            >
              <span>{tag}</span>
              <button
                onClick={() => removeTag(index)}
                className="ml-2 text-gray-500 hover:text-red-500"
              >
                <svg className="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                    clipRule="evenodd"
                  />
                </svg>
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};