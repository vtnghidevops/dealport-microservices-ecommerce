import { 
  BestSellingProduct, 
  Product, 
  ProductCategory,
  NewProduct
} from '../models/product.model';

export const ProductService = {
  getBestSellingProducts: async (): Promise<BestSellingProduct[]> => {
    // Giả lập API call
    await new Promise(resolve => setTimeout(resolve, 400));
    
    return [
      {
        product: {
          id: 1,  
          name: 'Apple iPhone 13',
          category: 'Electronic',
          image: 'images/limited/iphone_13.png',
          itemCode: '#PX2-4567',
          price: 999.00
        },
        totalOrder: 104,
        status: 'Stock',
        price: 999.00
      },
      {
        product: {
          id: 2,
          name: 'Nike Air Jordan',
          category: 'Fashion',
          image: 'images/limited/iphone_13.png',
          itemCode: '#PX2-4567',
          price: 999.00
        },
        totalOrder: 56,
        status: 'Stock out',
        price: 999.00
      },
      {
        product: {
          id: 3,
          name: 'T-shirt',
          category: 'Fashion',
          image: 'images/limited/iphone_13.png',
          itemCode: '#PX2-4567',
          price: 999.00
        },
        totalOrder: 266,
        status: 'Stock',
        price: 999.00
      },
      {
        product: {
          id: 4,
          name: 'Cross Bag',
          category: 'Fashion',
          image: 'images/limited/iphone_13.png',
          itemCode: '#PX2-4567',
          price: 999.00
        },
        totalOrder: 506,
        status: 'Stock',
        price: 999.00
      }
    ];
  },
  
  getTopProducts: async (): Promise<Product[]> => {
    // Giả lập API call
    await new Promise(resolve => setTimeout(resolve, 300));
    
    return [
      {
        id: 1,
        name: 'Apple iPhone 13',
        category: 'Electronic',
        price: 999.00,
        image: 'images/limited/iphone_13.png',
        itemCode: '#PX2-4567'
      },
      {
        id: 2,
        name: 'Nike Air Jordan',
        category: 'Fashion',
        price: 72.40,
        image: 'images/limited/iphone_13.png',
        itemCode: '#PX2-4567'
      },
      {
        id: 3,
        name: 'T-shirt',
        category: 'Fashion',
        price: 35.40,
        image: 'images/limited/iphone_13.png',
        itemCode: '#PX2-4567'
      },
      {
        id: 4,
        name: 'Assorted Cross Bag',
        category: 'Fashion',
        price: 80.00,
        image: 'images/limited/iphone_13.png',
        itemCode: '#PX2-4567'
      }
    ];
  },
  
  getProductCategories: async (): Promise<ProductCategory[]> => {
    // Giả lập API call
    await new Promise(resolve => setTimeout(resolve, 200));
    
    return [
      {
        id: 1,
        name: 'Electronic',
        image: 'images/limited/iphone_13.png'
      },
      {
        id: 2,
        name: 'Fashion',
        image: 'images/limited/iphone_13.png'
      },
      {
        id: 3,
        name: 'Home',
        image: 'images/limited/iphone_13.png'
      }
    ];
  },
  
  getNewProducts: async (): Promise<NewProduct[]> => {
    // Giả lập API call
    await new Promise(resolve => setTimeout(resolve, 250));
    
    return [
      {
        id: 1,
        name: 'Smart Fitness Tracker',
        price: 39.99,
        image: 'images/limited/smart_watch.png'
      },
      {
        id: 2,
        name: 'Leather Wallet',
        price: 19.99,
        image: 'images/limited/smart_watch.png'
      },
      {
        id: 3,
        name: 'Electric Hair Trimmer',
        price: 24.99,
        image: 'images/limited/smart_watch.png'
      }
    ];
  }
};
