// components/user/OrderCard.tsx
import React, { useState } from 'react';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { FiChevronDown, FiChevronUp } from 'react-icons/fi';

interface OrderItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
  image?: string;
}

interface Order {
  id: string;
  date: string;
  status: string;
  total: number;
  items: OrderItem[];
}

interface OrderCardProps {
  order: Order;
}

const OrderCard: React.FC<OrderCardProps> = ({ order }) => {
  const [expanded, setExpanded] = useState(false);
  
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };
  
  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'processing':
        return 'bg-yellow-100 text-yellow-800';
      case 'shipped':
        return 'bg-blue-100 text-blue-800';
      case 'delivered':
        return 'bg-green-100 text-green-800';
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };
  
  return (
    <Card className='p-[10px] mb-5'>
      <CardHeader className="pb-2">
        <div className="flex flex-col md:flex-row justify-between">
          <div>
            <CardTitle className="text-lg font-bold">Order #{order.id}</CardTitle>
            <p className="text-[15px] text-gray-500">Placed on {formatDate(order.date)}</p>
          </div>
          <div className="flex items-center mt-2 md:mt-0">
            <Badge className={`${getStatusColor(order.status)} capitalize`}>
              {order.status}
            </Badge>
          </div>
        </div>
      </CardHeader>
      
      <CardContent className="pb-2">
        <div className="flex justify-between mb-2">
          <span className="text-[15px] text-gray-500">Items: {order.items.length}</span>
          <span className="font-bold">${order.total.toFixed(2)}</span>
        </div>
        
        {expanded && (
          <div className="mt-4 space-y-3 border-t pt-3">
            {order.items.map(item => (
              <div key={item.id} className="flex items-center">
                {item.image && (
                  <div className="mr-3">
                    <img 
                      src={item.image} 
                      alt={item.name}
                      className="w-[30px] h-[30px] object-cover rounded"
                    />
                  </div>
                )}
                <div className="flex-1">
                  <p className="font-medium">{item.name}</p>
                  <div className="text-sm text-gray-500">
                    {item.quantity} x ${item.price.toFixed(2)}
                  </div>
                </div>
                <div className="font-bold">
                  ${(item.price * item.quantity).toFixed(2)}
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
      
      <CardFooter className="pt-2 flex justify-between">
        <Button 
          variant="ghost" 
          size="sm"
          onClick={() => setExpanded(!expanded)}
          className="flex items-center text-[15px]"
        >
          {expanded ? (
            <>
              <FiChevronUp className="mr-1 !w-[15px] !h-[15px]" /> Hide Details
            </>
          ) : (
            <>
              <FiChevronDown className="mr-1 !w-[15px] !h-[15px]" /> View Details
            </>
          )}
        </Button>
        <Button variant="outline" className='bg-[#0496FF] text-white hover:bg-blue-600 w-[120px] h-[40px] text-[14px] p-2'>
          Track Order
        </Button>
      </CardFooter>
    </Card>
  );
};

export default OrderCard;