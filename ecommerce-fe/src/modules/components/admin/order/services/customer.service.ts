import { Customer } from '../models/customer.model';

export const fetchCustomerDetails = async (customerId: string): Promise<Customer | null> => {
  try {
    // In a real application, this would be an API call
    const response = await fetch(`/api/customers/${customerId}`);
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch customer details:', error);
    return null;
  }
};