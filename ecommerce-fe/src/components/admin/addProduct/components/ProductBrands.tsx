import React, { ChangeEvent, useState, useEffect } from "react";

interface ProductBrandsProps {
  brand: string;
  onChange: (e: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => void;
}

export const ProductBrands: React.FC<ProductBrandsProps> = ({
  brand,
  onChange,
}) => {
  const [brands, setBrands] = useState<string[]>([
    "Apple",
    "Samsung",
    "Sony",
    "LG",
    "Nike",
    "Adidas",
    "Puma",
    "Dell",
    "HP",
    "Lenovo",
  ]);

  const [searchTerm, setSearchTerm] = useState<string>("");
  const [showDropdown, setShowDropdown] = useState<boolean>(false);
  const [customBrand, setCustomBrand] = useState<boolean>(false);

  useEffect(() => {
    // Check if current brand is not in the list
    if (brand && !brands.includes(brand)) {
      setCustomBrand(true);
    }
  }, [brand, brands]);

  const filteredBrands = brands.filter(b =>
    b.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleBrandSelect = (selectedBrand: string) => {
    const fakeEvent = {
      target: {
        name: "brand",
        value: selectedBrand
      }
    } as ChangeEvent<HTMLInputElement>;

    onChange(fakeEvent);
    setShowDropdown(false);
    setSearchTerm("");
  };

  const toggleCustomBrand = () => {
    setCustomBrand(!customBrand);
    if (!customBrand) {
      setSearchTerm("");
    }
  };

  const handleAddNewBrand = () => {
    if (searchTerm.trim() && !brands.includes(searchTerm.trim())) {
      // Add to brands list
      setBrands([...brands, searchTerm.trim()]);
      // Select it
      handleBrandSelect(searchTerm.trim());
    }
  };

  return (
    <div className="bg-white rounded-lg p-[1.25rem] drop-shadow-sm filter mb-6">
      <h2 className="font-bold text-cyprus text-[22px] mb-[15px]">
        Brand Information
      </h2>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-2">
          Product Brand
        </label>

        {customBrand ? (
          <div className="flex">
            <input
              type="text"
              name="brand"
              value={brand || ""}
              onChange={onChange}
              placeholder="Enter custom brand name"
              className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
            />
            <button
              onClick={toggleCustomBrand}
              className="ml-2 px-3 py-2 border border-neutral-300 rounded-lg bg-neutral-50 hover:bg-neutral-100"
            >
              Select
            </button>
          </div>
        ) : (
          <div className="relative">
            <div
              className="flex items-center justify-between text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 cursor-pointer"
              onClick={() => setShowDropdown(!showDropdown)}
            >
              <span>{brand || "Select a brand"}</span>
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </div>

            {showDropdown && (
              <div className="absolute z-10 w-full mt-1 bg-white border border-neutral-300 rounded-lg shadow-lg">
                <div className="p-2 border-b">
                  <input
                    type="text"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    placeholder="Search brands..."
                    className="w-full p-2 border border-neutral-300 rounded-lg focus:outline-none focus:border-ocean-green"
                  />
                </div>

                <div className="max-h-60 overflow-y-auto">
                  {filteredBrands.length > 0 ? (
                    filteredBrands.map((b, index) => (
                      <div
                        key={index}
                        className="p-2 hover:bg-neutral-100 cursor-pointer"
                        onClick={() => handleBrandSelect(b)}
                      >
                        {b}
                      </div>
                    ))
                  ) : (
                    <div className="p-2 text-center text-neutral-500">
                      No brands found
                      {searchTerm && (
                        <button
                          onClick={handleAddNewBrand}
                          className="ml-2 text-ocean-green hover:underline"
                        >
                          Add "{searchTerm}"
                        </button>
                      )}
                    </div>
                  )}
                </div>

                <div className="p-2 border-t">
                  <button
                    onClick={toggleCustomBrand}
                    className="w-full text-left text-ocean-green hover:underline"
                  >
                    + Add custom brand
                  </button>
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      <div className="mb-4">
        <label className="block text-cyprus font-bold text-[15px] mb-2">
          Brand Origin <span className="text-neutral-500">(Optional)</span>
        </label>
        <input
          type="text"
          name="brandOrigin"
          onChange={onChange}
          placeholder="e.g. United States, Japan, etc."
          className="text-cyprus bg-neutral-50 w-full border border-neutral-300 rounded-lg p-2 focus:outline-none focus:border-ocean-green"
        />
      </div>
    </div>
  );
}; 