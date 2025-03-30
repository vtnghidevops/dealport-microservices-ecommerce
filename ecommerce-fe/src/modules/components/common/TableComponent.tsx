import React from "react";
import Pagination from "./Pagination";

interface TableColumn<T> {
  header: string;
  key: keyof T | string;
  render?: (item: T) => React.ReactNode;
  align?: "left" | "center" | "right";
  width?: string;
}

export interface TableProps<T> {
  columns: TableColumn<T>[];
  data: T[];
  keyExtractor: (item: T, index: number) => string;
  totalItems: number;
  currentPage: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  loading?: boolean;
  onRowClick?: (item: T) => void;
  isSelectable?: boolean;
  selectedItems?: string[];
  onSelectItem?: (id: string, isSelected: boolean) => void;
  onSelectAll?: (isSelected: boolean) => void;
  tableClassName?: string;
  headerClassName?: string;
  rowClassName?: (item: T) => string;
  emptyMessage?: string;
}

/**
 * A reusable table component that can be used to display any type of data
 * Supports pagination, selection, custom cell rendering, and more
 */
function TableComponent<T>({
  columns,
  data,
  keyExtractor,
  totalItems,
  currentPage,
  pageSize,
  onPageChange,
  loading = false,
  onRowClick,
  isSelectable = false,
  selectedItems = [],
  onSelectItem,
  onSelectAll,
  tableClassName = "",
  headerClassName = "bg-aqua-spring",
  rowClassName = () => "",
  emptyMessage = "No data available"
}: TableProps<T>) {
  const handleRowClick = (item: T) => {
    if (onRowClick) {
      onRowClick(item);
    }
  };

  const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (onSelectAll) {
      onSelectAll(e.target.checked);
    }
  };

  const handleSelectItem = (id: string, isSelected: boolean) => {
    if (onSelectItem) {
      onSelectItem(id, isSelected);
    }
  };

  const renderCell = (item: T, column: TableColumn<T>, index: number) => {
    if (column.render) {
      return column.render(item);
    }

    const key = column.key as keyof T;
    return item[key] as React.ReactNode;
  };

  if (loading) {
    return (
      <div className="animate-pulse">
        <div className="h-8 bg-gray-200 rounded mb-4"></div>
        <div className="h-40 bg-gray-200 rounded"></div>
      </div>
    );
  }

  return (
    <div>
      <div className="overflow-x-auto relative rounded-lg border">
        <table className={`w-full text-sm text-left ${tableClassName}`}>
          <thead className={`text-xs uppercase ${headerClassName}`}>
            <tr className="h-[56px]">
              {isSelectable && (
                <th className="py-4 px-6 text-center">
                  <input
                    type="checkbox"
                    onChange={handleSelectAll}
                    className="h-[16px] w-[16px] rounded border border-ocean-green checked:bg-success checked:border-success appearance-none relative checked:after:content-['✓'] checked:after:text-white checked:after:absolute checked:after:text-xs checked:after:left-[3px]"
                  />
                </th>
              )}
              {columns.map((column, idx) => (
                <th
                  key={idx}
                  className={`py-4 px-6 text-${column.align || "left"} text-[15px] font-medium text-cyprus ${column.width ? `w-[${column.width}]` : ""}`}
                >
                  {column.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.length > 0 ? (
              data.map((item, index) => {
                const key = keyExtractor(item, index);
                const isSelected = selectedItems.includes(key);
                return (
                  <tr
                    key={key}
                    onClick={() => handleRowClick(item)}
                    className={`border-b hover:bg-gray-50 ${rowClassName(item)} ${onRowClick ? "cursor-pointer" : ""} h-[68px]`}
                  >
                    {isSelectable && (
                      <td className="py-4 px-6 text-center" onClick={(e) => e.stopPropagation()}>
                        <input
                          type="checkbox"
                          checked={isSelected}
                          onChange={(e) => handleSelectItem(key, e.target.checked)}
                          className="h-[16px] w-[16px] rounded border border-ocean-green checked:bg-success checked:border-success appearance-none relative checked:after:content-['✓'] checked:after:text-white checked:after:absolute checked:after:text-xs checked:after:left-[3px]"
                        />
                      </td>
                    )}
                    {columns.map((column, idx) => (
                      <td
                        key={`${key}-${idx}`}
                        className={`py-4 px-6 text-${column.align || "left"}`}
                      >
                        {renderCell(item, column, index)}
                      </td>
                    ))}
                  </tr>
                );
              })
            ) : (
              <tr>
                <td
                  colSpan={isSelectable ? columns.length + 1 : columns.length}
                  className="py-4 px-6 text-center text-gray-500"
                >
                  {emptyMessage}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      <Pagination
        currentPage={currentPage}
        totalItems={totalItems}
        pageSize={pageSize}
        onPageChange={onPageChange}
      />
    </div>
  );
}

export default TableComponent;