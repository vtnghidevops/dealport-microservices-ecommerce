// components/admin/product/add/components/ProductFeatures.tsx
import React, { useState } from "react";

interface ProductFeaturesProps {
  features: string[];
  onChange: (features: string[]) => void;
}

export const ProductFeatures: React.FC<ProductFeaturesProps> = ({
  features,
  onChange,
}) => {
  const [newFeature, setNewFeature] = useState<string>("");

  const addFeature = () => {
    if (newFeature.trim()) {
      onChange([...features, newFeature.trim()]);
      setNewFeature("");
    }
  };

  const removeFeature = (index: number) => {
    const updatedFeatures = [...features];
    updatedFeatures.splice(index, 1);
    onChange(updatedFeatures);
  };

  return (
    <div className="bg-white rounded-lg p-[1.25rem] drop-shadow-sm filter mb-6">
      <h2 className="font-bold text-cyprus text-[22px] mb-[15px]">
        Product Features
      </h2>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-12 mt-[1rem]">
          Add Product Features
        </label>
        <div className="flex">
          <input
            type="text"
            value={newFeature}
            onChange={(e) => setNewFeature(e.target.value)}
            placeholder="Enter a product feature"
            className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-l-lg p-2 focus:outline-none focus:border-ocean-green"
          />
          <button
            onClick={addFeature}
            className="bg-ocean-green text-white px-4 rounded-r-lg hover:bg-opacity-90"
            type="button"
          >
            Add
          </button>
        </div>
      </div>

      {features.length > 0 && (
        <div className="mt-4">
          <h3 className="font-semibold text-cyprus mb-2">Features:</h3>
          <ul className="bg-neutral-50 border border-neutral-200 rounded-lg p-3">
            {features.map((feature, index) => (
              <li key={index} className="flex justify-between items-center py-2 border-b border-neutral-200 last:border-0">
                <span>{feature}</span>
                <button
                  onClick={() => removeFeature(index)}
                  className="text-red-500 hover:text-red-700"
                >
                  <svg
                    className="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};