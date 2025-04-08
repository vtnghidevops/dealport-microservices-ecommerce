// components/admin/product/services/product.service.ts
import { Product } from '../models/product.model';

class ProductService {
  private apiUrl = '/api/products';

  async getProducts(): Promise<Product[]> {
    try {
      const response = await fetch(this.apiUrl);
      return await response.json();
    } catch (error) {
      console.error('Error fetching products:', error);
      return [];
    }
  }

  async getProductById(id: string): Promise<Product | null> {
    try {
      const response = await fetch(`${this.apiUrl}/${id}`);
      return await response.json();
    } catch (error) {
      console.error(`Error fetching product ${id}:`, error);
      return null;
    }
  }

  async createProduct(product: Omit<Product, 'id'>): Promise<Product | null> {
    try {
      const response = await fetch(this.apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(product),
      });
      return await response.json();
    } catch (error) {
      console.error('Error creating product:', error);
      return null;
    }
  }

  async updateProduct(id: string, product: Partial<Product>): Promise<Product | null> {
    try {
      const response = await fetch(`${this.apiUrl}/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(product),
      });
      return await response.json();
    } catch (error) {
      console.error(`Error updating product ${id}:`, error);
      return null;
    }
  }

  async deleteProduct(id: string): Promise<boolean> {
    try {
      await fetch(`${this.apiUrl}/${id}`, {
        method: 'DELETE',
      });
      return true;
    } catch (error) {
      console.error(`Error deleting product ${id}:`, error);
      return false;
    }
  }
}

export const productService = new ProductService();