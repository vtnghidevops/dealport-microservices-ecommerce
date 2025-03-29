export interface Transaction {
  id: number;
  customerId: string;
  orderDate: string;
  status: 'Paid' | 'Pending' | 'Canceled' | 'Processing';
  amount: number;
}

export interface TransactionListResponse {
  transactions: Transaction[];
  total: number;
}
