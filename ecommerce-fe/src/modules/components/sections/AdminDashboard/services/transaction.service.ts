import { Transaction, TransactionListResponse } from '../models/transaction.model';

export const TransactionService = {
  getTransactions: async (page: number = 1, limit: number = 10): Promise<TransactionListResponse> => {
    // Giả lập API call
    await new Promise(resolve => setTimeout(resolve, 300));
    
    const transactions: Transaction[] = [
      {
        id: 1,
        customerId: '#6545',
        orderDate: '01 Oct | 11:29 am',
        status: 'Paid',
        amount: 64
      },
      {
        id: 2,
        customerId: '#5412',
        orderDate: '01 Oct | 11:29 am',
        status: 'Pending',
        amount: 557
      },
      {
        id: 3,
        customerId: '#6622',
        orderDate: '01 Oct | 11:29 am',
        status: 'Paid',
        amount: 156
      },
      {
        id: 4,
        customerId: '#6462',
        orderDate: '01 Oct | 11:29 am',
        status: 'Paid',
        amount: 265
      },
      {
        id: 5,
        customerId: '#6462',
        orderDate: '01 Oct | 11:29 am',
        status: 'Paid',
        amount: 265
      }
    ];
    
    return {
      transactions,
      total: 100 // Tổng số để phân trang
    };
  }
};
