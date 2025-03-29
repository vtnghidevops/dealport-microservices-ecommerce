// Updated OrderTable component to match the design
import React from "react";
import { Order, OrderStatus } from "../models/order.model";
import { CiDeliveryTruck } from "react-icons/ci";

interface OrderTableProps {
  orders: Order[];
  onStatusChange?: (orderId: string, status: OrderStatus) => void;
  onViewDetails?: (orderId: string) => void;
}

export const OrderTable: React.FC<OrderTableProps> = ({
  orders,
  onStatusChange,
  onViewDetails,
}) => {
  const getStatusIcon = (status: OrderStatus) => {
    switch (status) {
      case "Delivered":
        return (
          <div className="flex items-center text-green-500">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Delivered</span>
          </div>
        );
      case "Pending":
        return (
          <div className="flex items-center text-amber-500">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Pending</span>
          </div>
        );
      case "Shipped":
        return (
          <div className="flex items-center">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Shipped</span>
          </div>
        );
      case "Cancelled":
        return (
          <div className="flex items-center text-red-500">
            <CiDeliveryTruck className="h-[20px] w-[20px] mr-1"></CiDeliveryTruck>
            <span>Cancelled</span>
          </div>
        );
      default:
        return <span>{status}</span>;
    }
  };

  const getPaymentStatusIcon = (paymentStatus: "Paid" | "Unpaid") => {
    switch (paymentStatus) {
      case "Paid":
        return (
          <div className="flex items-center">
            <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
            <span>Paid</span>
          </div>
        );
      case "Unpaid":
        return (
          <div className="flex items-center">
            <span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
            <span>Unpaid</span>
          </div>
        );
      default:
        return null;
    }
  };

  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse bg-white rounded-lg">
        <thead className="bg-aqua-spring">
          <tr className="h-[56px]">
            <th className="py-4 px-6 text-center text-[15px] font-medium text-cyprus">No.</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Order Id</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Product</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Date</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Price</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Payment</th>
            <th className="py-4 px-6 text-left text-[15px] font-medium text-cyprus">Status</th>
          </tr>
        </thead>
        <tbody>
          {orders.map((order, index) =>
            order.products.map((product, productIndex) => (
              <tr key={`${order.id}-${productIndex}`} className="border-b h-[68px]">
                <td className="py-4 px-6 ">
                  <div className="flex items-center justify-center">
                  <input
  type="checkbox"
  className="mr-2 h-[16px] w-[16px] rounded border border-ocean-green checked:bg-success checked:border-success appearance-none relative checked:after:content-['✓'] checked:after:text-white checked:after:absolute checked:after:text-xs  checked:after:left-[3px]"
/>
                    <span>{index + 1}</span>
                  </div>
                </td>
                <td className="py-4 px-6 font-medium text-[15px]">
                  #{order.orderId}
                </td>
                <td className="py-4 px-6">
                  <div className="flex items-center">
                    <div className="border border-neutral-200 w-[40px] h-[40px] mr-3 rounded flex items-center justify-center overflow-hidden">
                      <img
                        src={product.productImage}
                        alt={product.productName}
                        className="object-contain"
                      />
                    </div>
                    <span className="text-[15px] max-w-[140px]">{product.productName}</span>
                  </div>
                </td>
                <td className="py-4 px-6 text-[15px]">
                  {order.date}
                </td>
                <td className="py-4 px-6 text-[15px]">
                  {product.price.toFixed(2)}
                </td>
                <td className="py-4 px-6 text-[15px]">
                  {getPaymentStatusIcon(order.paymentStatus)}
                </td>
                <td className="py-4 px-6 text-[15px]">
                  {getStatusIcon(order.status)}
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
      
      
    </div>
  );
};