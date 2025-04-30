// Updated OrderTable component to match the design
import React from "react";
import { Order, OrderStatus } from "@/services/user/order.service";
import { CiDeliveryTruck } from "react-icons/ci";

interface OrderTableProps {
  orders: Order[];
  onStatusChange?: (orderId: string, status: string) => void;
  onViewDetails?: (orderId: string) => void;
}

export const OrderTable: React.FC<OrderTableProps> = ({
  orders,
  onStatusChange,
  onViewDetails,
}) => {
  const getStatusIcon = (status: OrderStatus) => {
    switch (status) {
      case OrderStatus.Delivered:
        return (
          <div className="flex items-center text-green-500">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Delivered</span>
          </div>
        );
      case OrderStatus.Pending:
        return (
          <div className="flex items-center text-amber-500">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Pending</span>
          </div>
        );
      case OrderStatus.Shipped:
        return (
          <div className="flex items-center">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Shipped</span>
          </div>
        );
      case OrderStatus.Cancelled:
        return (
          <div className="flex items-center text-red-500">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Cancelled</span>
          </div>
        );
      case OrderStatus.Paid:
        return (
          <div className="flex items-center text-green-500">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Paid</span>
          </div>
        );
      case OrderStatus.Processing:
        return (
          <div className="flex items-center text-blue-500">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Processing</span>
          </div>
        );
      case OrderStatus.Refunded:
        return (
          <div className="flex items-center text-orange-500">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Refunded</span>
          </div>
        );
      default:
        return <span>{status || 'Unknown'}</span>;
    }
  };

  const getPaymentStatusIcon = (paymentStatus: string) => {
    if (!paymentStatus) {
      return (
        <div className="flex items-center">
          <span className="w-2 h-2 bg-gray-500 rounded-full mr-2"></span>
          <span>Unknown</span>
        </div>
      );
    }

    switch (paymentStatus.toLowerCase()) {
      case "paid":
      case "completed":
        return (
          <div className="flex items-center">
            <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
            <span>Paid</span>
          </div>
        );
      case "pending":
        return (
          <div className="flex items-center">
            <span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
            <span>Pending</span>
          </div>
        );
      case "processing":
        return (
          <div className="flex items-center">
            <span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
            <span>Processing</span>
          </div>
        );
      case "failed":
        return (
          <div className="flex items-center">
            <span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
            <span>Failed</span>
          </div>
        );
      case "refunded":
        return (
          <div className="flex items-center">
            <span className="w-2 h-2 bg-orange-500 rounded-full mr-2"></span>
            <span>Refunded</span>
          </div>
        );
      default:
        return (
          <div className="flex items-center">
            <span className="w-2 h-2 bg-gray-500 rounded-full mr-2"></span>
            <span>{paymentStatus || "Unknown"}</span>
          </div>
        );
    }
  };

  // Format date string
  const formatDate = (dateString: string) => {
    if (!dateString) return 'N/A';

    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return 'Invalid date';

      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    } catch (error) {
      return 'Invalid date';
    }
  };

  // Xử lý hiển thị an toàn sản phẩm đầu tiên
  const getFirstProductDisplay = (order: Order) => {
    if (!order?.items || !Array.isArray(order.items) || order.items.length === 0) {
      return {
        name: "No products",
        imageUrl: null,
        hasMoreItems: false
      };
    }

    const firstItem = order.items[0];

    return {
      name: firstItem.name || "Unnamed product",
      imageUrl: firstItem.imageUrl || null,
      hasMoreItems: order.items.length > 1,
      additionalCount: order.items.length - 1
    };
  };

  // Get formatted order total
  const getOrderTotal = (order: Order): string => {
    if (!order) return '0.00';

    try {
      if (order.totals?.total !== undefined && order.totals.total !== null) {
        return order.totals.total.toFixed(2);
      } else if (typeof order.total === 'number') {
        return order.total.toFixed(2);
      }
      return '0.00';
    } catch (error) {
      console.error('Error formatting order total:', error);
      return '0.00';
    }
  };

  // Get order number or ID
  const getOrderIdentifier = (order: Order): string => {
    if (!order) return 'Unknown';

    if (order.orderNumber) {
      return order.orderNumber;
    }
    return (order.id || "Unknown").toString().substring(0, 8);
  };

  // Check if order is valid for display
  const isValidOrder = (order: any): boolean => {
    return order && typeof order === 'object' && order.id;
  };

  // Get payment status based on order and payment information
  const getOrderPaymentStatus = (order: Order): string => {
    // Nếu order status là paid thì payment status cũng phải là paid
    if (order.status === OrderStatus.Paid) {
      return 'paid';
    }

    // Nếu không thì lấy payment status từ order
    return order.paymentInfo?.status || order.paymentStatus || '';
  };

  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse bg-white rounded-lg">
        <thead className="bg-aqua-spring">
          <tr className="h-[56px]">
            <th className="py-4 px-6 text-center text-[15px] font-medium text-cyprus">No.</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Order ID</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Products</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Date</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Total</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Payment</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Status</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Actions</th>
          </tr>
        </thead>
        <tbody>
          {Array.isArray(orders) && orders.length > 0 ? (
            orders.map((order, index) => {
              if (!isValidOrder(order)) {
                return null; // Skip invalid orders
              }

              const firstProduct = getFirstProductDisplay(order);
              const paymentStatus = getOrderPaymentStatus(order);

              return (
                <tr key={order.id || index} className="border-b h-[68px]">
                  <td className="py-4 px-6 ">
                    <div className="flex items-center justify-center">
                      <input
                        type="checkbox"
                        className="mr-2 h-[16px] w-[16px] rounded border border-ocean-green checked:bg-success checked:border-success appearance-none relative checked:after:content-['✓'] checked:after:text-white checked:after:absolute checked:after:text-xs checked:after:left-[3px]"
                      />
                      <span>{index + 1}</span>
                    </div>
                  </td>
                  <td className="py-4 px-6 font-medium text-[15px]">
                    #{getOrderIdentifier(order)}
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex items-center">
                      {firstProduct.imageUrl && (
                        <div className="border border-neutral-200 w-[40px] h-[40px] mr-3 rounded flex items-center justify-center overflow-hidden">
                          <img
                            src={firstProduct.imageUrl}
                            alt={firstProduct.name}
                            className="object-contain p-1"
                          />
                        </div>
                      )}
                      <span className="text-[15px] max-w-[140px]">
                        {firstProduct.name}
                        {firstProduct.hasMoreItems ? ` +${firstProduct.additionalCount} more` : ''}
                      </span>
                    </div>
                  </td>
                  <td className="py-4 px-6 text-[15px]">
                    {formatDate(order.createdAt || '')}
                  </td>
                  <td className="py-4 px-6 text-[15px]">
                    ${getOrderTotal(order)}
                  </td>
                  <td className="py-4 px-6 text-[15px]">
                    {getPaymentStatusIcon(paymentStatus)}
                  </td>
                  <td className="py-4 px-6 text-[15px]">
                    {order.status ? getStatusIcon(order.status as OrderStatus) : 'Unknown'}
                  </td>
                  <td className="py-4 px-6 text-[15px]">
                    <div className="flex space-x-2">
                      <button
                        onClick={() => onViewDetails && order.id && onViewDetails(order.id)}
                        className="bg-blue-50 hover:bg-blue-100 text-blue-600 px-2 py-1 rounded"
                      >
                        View
                      </button>
                      {order.status && order.status !== OrderStatus.Delivered && order.status !== OrderStatus.Cancelled && order.id && (
                        <select
                          className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded focus:ring-blue-500 focus:border-blue-500 px-2 py-1"
                          value={order.status}
                          onChange={(e) => onStatusChange && onStatusChange(order.id, e.target.value)}
                        >
                          <option value={OrderStatus.Pending}>Pending</option>
                          <option value={OrderStatus.Processing}>Processing</option>
                          <option value={OrderStatus.Paid}>Paid</option>
                          <option value={OrderStatus.Shipped}>Shipped</option>
                          <option value={OrderStatus.Delivered}>Delivered</option>
                          <option value={OrderStatus.Cancelled}>Cancelled</option>
                          <option value={OrderStatus.Refunded}>Refunded</option>
                        </select>
                      )}
                    </div>
                  </td>
                </tr>
              );
            })
          ) : (
            <tr>
              <td colSpan={8} className="py-4 px-6 text-center text-gray-500">
                No orders available
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};
// Updated OrderTable component to use the reusable TableComponent
// import React from "react";
// import { Order, OrderStatus } from "../models/order.model";
// import { CiDeliveryTruck } from "react-icons/ci";
// import TableComponent from "../../../common/TableComponent";

// interface OrderTableProps {
//   orders: Order[];
//   onStatusChange?: (orderId: string, status: OrderStatus) => void;
//   onViewDetails?: (orderId: string) => void;
//   totalItems: number;
//   currentPage: number;
//   pageSize: number;
//   onPageChange: (page: number) => void;
//   loading?: boolean;
// }

// export const OrderTable: React.FC<OrderTableProps> = ({
//   orders,
//   onStatusChange,
//   onViewDetails,
//   totalItems,
//   currentPage,
//   pageSize,
//   onPageChange,
//   loading = false,
// }) => {
//   const getStatusIcon = (status: OrderStatus) => {
//     switch (status) {
//       case "Delivered":
//         return (
//           <div className="flex items-center text-green-500">
//             <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
//             <span>Delivered</span>
//           </div>
//         );
//       case "Pending":
//         return (
//           <div className="flex items-center text-amber-500">
//             <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
//             <span>Pending</span>
//           </div>
//         );
//       case "Shipped":
//         return (
//           <div className="flex items-center">
//             <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
//             <span>Shipped</span>
//           </div>
//         );
//       case "Cancelled":
//         return (
//           <div className="flex items-center text-red-500">
//             <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
//             <span>Cancelled</span>
//           </div>
//         );
//       default:
//         return <span>{status}</span>;
//     }
//   };

//   const getPaymentStatusIcon = (paymentStatus: "Paid" | "Unpaid") => {
//     switch (paymentStatus) {
//       case "Paid":
//         return (
//           <div className="flex items-center">
//             <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
//             <span>Paid</span>
//           </div>
//         );
//       case "Unpaid":
//         return (
//           <div className="flex items-center">
//             <span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
//             <span>Unpaid</span>
//           </div>
//         );
//       default:
//         return null;
//     }
//   };

//   // Prepare normalized data for the table
//   const normalizedOrders = orders.flatMap((order) =>
//     order.products.map((product, idx) => ({
//       ...order,
//       product,
//       productIndex: idx,
//       uniqueId: `${order.id}-${idx}`
//     }))
//   );

//   // Define columns for the table
//   const columns = [
//     {
//       header: "Order Id",
//       key: "orderId",
//       render: (row: any) => (
//         <span className="font-medium text-[15px]">#{row.orderId}</span>
//       ),
//     },
//     {
//       header: "Product",
//       key: "product",
//       render: (row: any) => (
//         <div className="flex items-center">
//           <div className="border border-neutral-200 w-[40px] h-[40px] mr-3 rounded flex items-center justify-center overflow-hidden">
//             <img
//               src={row.product.productImage}
//               alt={row.product.productName}
//               className="object-contain"
//             />
//           </div>
//           <span className="text-[15px] max-w-[140px]">{row.product.productName}</span>
//         </div>
//       ),
//     },
//     {
//       header: "Date",
//       key: "date",
//       render: (row: any) => <span className="text-[15px]">{row.date}</span>,
//     },
//     {
//       header: "Price",
//       key: "price",
//       render: (row: any) => <span className="text-[15px]">{row.product.price.toFixed(2)}</span>,
//     },
//     {
//       header: "Payment",
//       key: "paymentStatus",
//       render: (row: any) => getPaymentStatusIcon(row.paymentStatus),
//     },
//     {
//       header: "Status",
//       key: "status",
//       render: (row: any) => getStatusIcon(row.status),
//     },
//   ];

//   return (
//     <TableComponent
//       columns={columns}
//       data={normalizedOrders}
//       keyExtractor={(item) => item.uniqueId}
//       totalItems={totalItems}
//       currentPage={currentPage}
//       pageSize={pageSize}
//       onPageChange={onPageChange}
//       loading={loading}
//       isSelectable={true}
//       headerClassName="bg-aqua-spring"
//       rowClassName={() => "h-[68px]"}
//     />
//   );
// };