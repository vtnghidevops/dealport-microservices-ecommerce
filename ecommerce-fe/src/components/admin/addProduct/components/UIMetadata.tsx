import React, { useState } from "react";
import { Product } from "@/types/product.model";

interface UIMetadataProps {
  product: Product;
  onChange: (metadata: any) => void;
}

export const UIMetadata: React.FC<UIMetadataProps> = ({
  product,
  onChange,
}) => {
  const [jsonError, setJsonError] = useState<string | null>(null);
  const [metadataText, setMetadataText] = useState<string>(
    JSON.stringify(product.uiMetadata || {}, null, 2)
  );

  // Initial metadata templates for different product types
  const templates = {
    empty: {},
    product: {
      color: "#D9F3D9",
      featured: true,
      position: "main",
      display_options: {
        show_rating: true,
        show_price: true,
        show_compare: true
      },
      custom_fields: []
    },
    collection: {
      layout: "grid",
      items_per_page: 12,
      sort_options: ["price", "name", "date"],
      filters: {
        price: true,
        color: true,
        size: true,
        brand: true
      }
    }
  };

  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;
    setMetadataText(value);

    try {
      // Attempt to parse JSON
      const parsedJson = JSON.parse(value);
      setJsonError(null);
      onChange(parsedJson);
    } catch (err) {
      setJsonError("Invalid JSON format");
      // Don't update the actual product.ui_metadata until valid JSON
    }
  };

  const applyTemplate = (template: keyof typeof templates) => {
    const templateData = JSON.stringify(templates[template], null, 2);
    setMetadataText(templateData);
    setJsonError(null);
    onChange(templates[template]);
  };

  return (
    <div className="bg-white rounded-lg p-[1.25rem] drop-shadow-sm filter mb-6">
      <h2 className="font-bold text-cyprus text-[22px] mb-[15px]">
        UI Metadata <span className="text-neutral-500 text-sm">(Optional)</span>
      </h2>

      <div className="mb-4">
        <div className="flex justify-between items-center mb-2">
          <label className="block text-cyprus font-bold text-[15px]">
            Custom UI Settings (JSON)
          </label>
          <div className="flex space-x-2">
            <button
              type="button"
              onClick={() => applyTemplate('empty')}
              className="text-xs bg-gray-100 hover:bg-gray-200 text-gray-700 px-2 py-1 rounded"
            >
              Empty
            </button>
            <button
              type="button"
              onClick={() => applyTemplate('product')}
              className="text-xs bg-gray-100 hover:bg-gray-200 text-gray-700 px-2 py-1 rounded"
            >
              Product Template
            </button>
            <button
              type="button"
              onClick={() => applyTemplate('collection')}
              className="text-xs bg-gray-100 hover:bg-gray-200 text-gray-700 px-2 py-1 rounded"
            >
              Collection Template
            </button>
          </div>
        </div>

        <textarea
          value={metadataText}
          onChange={handleTextChange}
          rows={10}
          placeholder='{
  "color": "#D9F3D9",
  "featured": true,
  "position": "main"
}'
          className={`text-cyprus bg-neutral-50 w-full border ${jsonError ? "border-red-500" : "border-neutral-300"
            } rounded-lg p-2 font-mono text-sm focus:outline-none focus:border-ocean-green`}
        ></textarea>

        {jsonError && (
          <p className="text-red-500 text-xs mt-1">{jsonError}</p>
        )}

        <p className="text-xs text-gray-500 mt-1">
          Custom JSON metadata for UI rendering. This can include display preferences, feature flags, or any custom frontend settings.
        </p>

        <div className="mt-4">
          <h3 className="text-sm font-medium text-cyprus mb-2">Current Metadata Preview:</h3>
          <div className="bg-gray-100 p-3 rounded text-xs overflow-auto max-h-[150px]">
            <pre>{JSON.stringify(product.uiMetadata || {}, null, 2)}</pre>
          </div>
        </div>
      </div>
    </div>
  );
}; 