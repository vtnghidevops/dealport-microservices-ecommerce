import React from "react";

export interface FilterTab {
  id: string;
  label: string;
  count: number;
}

interface TabFilterProps {
  tabs: FilterTab[];
  activeTab: string;
  onTabChange: (tabId: string) => void;
  loading?: boolean;
  containerClassName?: string;
  tabClassName?: string;
  activeTabClassName?: string;
}

/**
 * A reusable tab filter component
 * Can be used for any filter that needs tab-based selection with counts
 */
const TabFilter: React.FC<TabFilterProps> = ({
  tabs,
  activeTab,
  onTabChange,
  loading = false,
  containerClassName = "bg-aqua-spring w-full h-[40px]",
  tabClassName = "px-[12px] py-[6px] h-[32px] rounded-md whitespace-nowrap",
  activeTabClassName = "bg-white text-black"
}) => {
  return (
    <div className={`gap-[14px] p-4 rounded-lg flex items-center space-x-2 overflow-x-auto ${containerClassName}`}>
      {tabs.map((tab) => (
        <button
          key={tab.id}
          onClick={() => onTabChange(tab.id)}
          className={`${tabClassName} ${
            activeTab === tab.id
              ? activeTabClassName
              : "bg-transparent text-neutral-600"
          }`}
        >
          {tab.label}{" "}
          <span className="text-xs ml-1">
            ({loading ? "..." : tab.count})
          </span>
        </button>
      ))}
    </div>
  );
};

export default TabFilter;