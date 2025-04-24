// components/admin/product/services/product.service.ts
import axios from 'axios';
import { Product } from '@/types/product.model';

// Tạo một instance của axios với cấu hình chung
const api = axios.create({
  baseURL: import.meta.env.VITE_PUBLIC_PRODUCT_API_URL || 'http://localhost:8082/api/v1',
  headers: {
    'Content-Type': 'application/json',
  }
});

// Interceptor để thêm token vào header cho mỗi request
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Interface cho định dạng ProductImage response từ backend
interface ProductImageResponse {
  url: string;
  is_primary?: boolean;
  display_order?: number;
  product_id?: number;
}

class ProductService {
  // API endpoints
  private apiUrl = '/products';
  private categoryUrl = '/categories';

  async getProducts(page = 1, pageSize = 10, filters = {}): Promise<{ products: Product[], total: number }> {
    try {
      // Lọc bỏ các giá trị undefined hoặc null từ filters
      const filteredParams = Object.fromEntries(
        Object.entries(filters).filter(([_, value]) => value !== undefined && value !== null)
      );

      const params = {
        page: page.toString(),
        page_size: pageSize.toString(),
        ...filteredParams
      };

      const response = await api.get(this.apiUrl, { params });
      const data = response.data;

      // Transform products to ensure imgSlider is always a string[]
      const products = (data.products || data.items || data.data || []).map((product: any) => {
        return this.normalizeProductResponse(product);
      });

      return {
        products,
        total: data.total || data.total_count || data.count || 0
      };
    } catch (error) {
      console.error('Error fetching products:', error);
      return { products: [], total: 0 };
    }
  }

  async getProductById(id: number): Promise<Product | null> {
    try {
      const response = await api.get(`${this.apiUrl}/${id}`);
      const data = response.data;
      const product = data.product || data.data || data;

      // Normalize the product to ensure correct format for frontend
      return this.normalizeProductResponse(product);
    } catch (error) {
      console.error(`Error fetching product ${id}:`, error);
      return null;
    }
  }

  // Normalize a product response from backend to ensure imgSlider is a string[]
  private normalizeProductResponse(product: any): Product {
    if (!product) return product;

    // Handle imgSlider if it's an array of complex objects
    if (product.imgSlider && Array.isArray(product.imgSlider)) {
      // Check if imgSlider contains complex objects
      if (product.imgSlider.length > 0 && typeof product.imgSlider[0] === 'object') {
        const urlArray = product.imgSlider.map((item: ProductImageResponse) => item.url);
        product.imgSlider = urlArray;

        // Also ensure image_url is a string, not an object
        if (product.image_url && typeof product.image_url === 'object' && 'url' in product.image_url) {
          product.image_url = (product.image_url as ProductImageResponse).url;
        }
      }
    } else {
      product.imgSlider = [];
    }

    return product;
  }

  // Method to upload images
  async uploadImages(files: File[], productId?: number): Promise<string[]> {
    try {
      console.log('🔍 Starting image upload for', files.length, 'files', productId ? `with product ID: ${productId}` : 'without product ID');
      const uploadUrls: string[] = [];

      // Check if backend is available
      const isBackendAvailable = await this.isBackendAvailable();
      console.log('🔍 Backend available for upload?', isBackendAvailable);

      if (!isBackendAvailable) {
        console.warn('⚠️ Backend upload service not available, using local file preview');
        // Create URLs for local file preview
        const localUrls = await Promise.all(
          files.map(file => new Promise<string>((resolve) => {
            const reader = new FileReader();
            reader.onloadend = () => {
              resolve(reader.result as string);
            };
            reader.readAsDataURL(file);
          }))
        );
        return localUrls;
      }

      // Validate product ID
      if (!productId) {
        console.error('⚠️ No product ID provided for image upload - the backend requires a valid numeric ID');
        throw new Error('Missing product ID for image upload');
      }

      // Ensure productId is a number (not a string) before using it
      if (typeof productId === 'string') {
        productId = parseInt(productId);
        if (isNaN(productId)) {
          console.error('⚠️ Invalid product ID format:', productId);
          throw new Error('Invalid product ID format');
        }
      }

      console.log(`🔍 Using product ID for upload: ${productId} (type: ${typeof productId})`);

      // Get base URL from environment or default to localhost
      const apiBaseUrl = import.meta.env.VITE_PUBLIC_PRODUCT_API_URL || 'http://localhost:8082/api/v1';
      // Extract the base URL without /api/v1
      const baseUrl = apiBaseUrl.replace(/\/api\/v1$/, '');
      console.log('🔍 Base URL for image upload:', baseUrl);

      // Upload each file individually to the server
      for (const file of files) {
        try {
          console.log(`🔍 Preparing to upload file: ${file.name} (${file.size} bytes, type: ${file.type})`);

          // Create form data for file upload
          const formData = new FormData();
          formData.append('image', file);
          formData.append('isPrimary', uploadUrls.length === 0 ? 'true' : 'false'); // First image is primary

          // Log FormData contents
          console.log('🔍 FormData created with keys:', [...formData.keys()]);

          // Create upload endpoint URL with numeric ID
          const uploadEndpoint = `${this.apiUrl}/${productId}/images`;
          console.log('🔍 Upload endpoint:', uploadEndpoint);
          console.log('🔍 Full upload URL:', api.defaults.baseURL + uploadEndpoint);

          // Use the endpoint with product ID
          const response = await api.post(uploadEndpoint, formData, {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          });

          console.log('🔍 Image upload response status:', response.status);
          console.log('🔍 Image upload response data:', response.data);

          // Get image URL from response
          let imageUrl = '';

          // Handle the standard response format first
          if (response.data && response.data.data && response.data.data.url) {
            imageUrl = response.data.data.url;
            console.log('🔍 Found URL in response.data.data.url:', imageUrl);
          } else if (response.data && response.data.url) {
            imageUrl = response.data.url;
            console.log('🔍 Found URL in response.data.url:', imageUrl);
          } else if (response.data && typeof response.data === 'string') {
            // Direct string URL in response
            imageUrl = response.data;
            console.log('🔍 Found URL in response.data (string):', imageUrl);
          } else {
            // Try to find any URL or path in the response
            const responseStr = JSON.stringify(response.data);
            const urlMatches = responseStr.match(/"(\/[^"]+)"/);
            if (urlMatches && urlMatches[1]) {
              imageUrl = urlMatches[1];
              console.log('🔍 Extracted URL from response JSON:', imageUrl);
            }
          }

          // Ensure URL is properly formatted
          if (imageUrl) {
            // Add the base URL if necessary
            if (imageUrl.startsWith('/')) {
              // Handle both /api/products/images/ and /uploads/products paths
              if (imageUrl.includes('/api/products/images/') || imageUrl.includes('/uploads/products/')) {
                imageUrl = `${baseUrl}${imageUrl}`;
                console.log('🔍 Converted relative URL to absolute:', imageUrl);
              }
            }

            console.log('✅ Final image URL:', imageUrl);
            uploadUrls.push(imageUrl);
          } else {
            console.warn('⚠️ Could not extract image URL from response', response.data);

            // Use a placeholder URL or throw an error
            throw new Error('Failed to get image URL from server response');
          }
        } catch (error) {
          console.error(`❌ Error uploading file ${file.name}:`, error);
          // Continue with next file rather than failing completely
        }
      }

      // Return all successful uploads
      if (uploadUrls.length === 0) {
        console.error('❌ No files were successfully uploaded');
        throw new Error('Failed to upload images to server');
      }

      console.log('✅ Successfully uploaded', uploadUrls.length, 'images:', uploadUrls);
      return uploadUrls;
    } catch (error) {
      console.error('❌ Error in uploadImages:', error);
      throw error;
    }
  }

  // Helper method to check if backend is available
  private async isBackendAvailable(): Promise<boolean> {
    try {
      console.log('🔍 Checking if backend is available...');

      // Mặc định, set availability là false, đợi kết quả check
      let isAvailable = false;

      // Get base URL from environment or default to localhost
      const apiBaseUrl = import.meta.env.VITE_PUBLIC_PRODUCT_API_URL || 'http://localhost:8082/api/v1';
      // Extract the base URL without /api/v1
      const baseUrl = apiBaseUrl.replace(/\/api\/v1$/, '');

      // Log API URL for debugging
      console.log('🔍 Checking backend at base URL:', baseUrl);
      console.log('🔍 API URL:', apiBaseUrl);

      try {
        // Try to fetch health endpoint first at the BASE URL (not API URL)
        const healthResponse = await axios.get(`${baseUrl}/health`, { timeout: 3000 });
        console.log('🔍 Health endpoint response:', healthResponse.status, healthResponse.data);
        isAvailable = healthResponse.status === 200;
      } catch (healthError) {
        console.warn('⚠️ Health endpoint check failed:', healthError);

        // If health endpoint fails, try categories endpoint as fallback
        try {
          const categoriesResponse = await api.get('/categories', { timeout: 3000 });
          console.log('🔍 Categories endpoint response:', categoriesResponse.status);
          isAvailable = categoriesResponse.status === 200;
        } catch (categoriesError) {
          console.warn('⚠️ Categories endpoint check failed:', categoriesError);

          // Last attempt: try to access the static file server endpoint
          try {
            const staticResponse = await fetch(`${baseUrl}/api/products/images/test.jpg`, {
              method: 'HEAD'
            });
            console.log('🔍 Static file server response:', staticResponse.status);
            isAvailable = staticResponse.ok;
          } catch (staticError) {
            console.warn('⚠️ Static file server check failed:', staticError);
            isAvailable = false;
          }
        }
      }

      console.log('🔍 Backend availability result:', isAvailable);
      return isAvailable;
    } catch (error) {
      console.warn('⚠️ Backend availability check failed with error:', error);
      return false;
    }
  }

  async createProduct(product: Partial<Product>): Promise<Product | null> {
    try {
      console.log('Original product data:', product);

      // First create the product - preserve the original data but remove id
      const productObj = { ...product } as any; // Use any to handle potential snake_case props

      // Remove id property and any potential snake_case duplicates
      delete productObj.id;
      delete productObj.original_price;
      delete productObj.stock_quantity;

      // Convert categoryId from string to number
      if (productObj.categoryId) {
        productObj.categoryId = parseInt(String(productObj.categoryId), 10);
      }

      // Ensure images are properly formatted for the backend
      if (productObj.imgSlider && Array.isArray(productObj.imgSlider) && productObj.imgSlider.length > 0) {
        // Convert any complex objects in imgSlider to simple URLs
        productObj.imgSlider = productObj.imgSlider.map((img: any) => {
          if (typeof img === 'object' && img !== null && 'url' in img) {
            return img.url;
          }
          return img;
        });

        // Ensure we have a primary image
        productObj.image_url = productObj.imgSlider[0];

        // Create proper image objects for the backend
        productObj.images = productObj.imgSlider.map((url: string, index: number) => ({
          url: url,
          is_primary: index === 0,
          display_order: index
        }));
      }

      console.log('Creating product with data:', productObj);
      const response = await api.post(this.apiUrl, productObj);

      // Debug log the full response
      console.log('Full API response:', {
        status: response.status,
        statusText: response.statusText,
        headers: response.headers,
        data: response.data
      });

      // Attempt to extract the product ID from various response formats
      let productId = null;

      if (response.data) {
        if (typeof response.data === 'object') {
          // Try all possible locations for the product ID
          const possibleIdFields = ['id', 'product_id', 'productId', 'ID'];

          // Check direct properties on response.data
          for (const field of possibleIdFields) {
            if (response.data[field] !== undefined) {
              productId = response.data[field];
              console.log(`Found product ID in response.data.${field}:`, productId);
              break;
            }
          }

          // Check if ID is in a nested 'product' object
          if (!productId && response.data.product) {
            for (const field of possibleIdFields) {
              if (response.data.product[field] !== undefined) {
                productId = response.data.product[field];
                console.log(`Found product ID in response.data.product.${field}:`, productId);
                break;
              }
            }
          }

          // Check if ID is in a nested 'data' object
          if (!productId && response.data.data) {
            for (const field of possibleIdFields) {
              if (response.data.data[field] !== undefined) {
                productId = response.data.data[field];
                console.log(`Found product ID in response.data.data.${field}:`, productId);
                break;
              }
            }
          }
        } else if (typeof response.data === 'number') {
          // Some APIs might return the ID directly as a number
          productId = response.data;
          console.log('Found product ID as direct number response:', productId);
        }
      }

      // Final check if we found an ID
      if (productId) {
        console.log('✅ Successfully extracted product ID:', productId);

        // Get the full product data
        const productData = await this.getProductById(productId);
        if (productData) {
          return productData;
        }

        // If getProductById fails but we have product data in the response, use that
        if (response.data.product) {
          const responseProduct = this.normalizeProductResponse(response.data.product);
          console.log('Using product data from response:', responseProduct);
          return responseProduct;
        } else if (response.data.data) {
          const responseProduct = this.normalizeProductResponse(response.data.data);
          console.log('Using product data from response.data:', responseProduct);
          return responseProduct;
        } else {
          // If we only have the ID, construct a minimal product object
          console.log('Constructing minimal product with ID:', productId);
          return {
            ...product,
            id: productId
          } as Product;
        }
      } else {
        console.error('❌ Could not extract product ID from response:', response.data);

        // Last attempt: try to find any object that might be the product
        if (response.data && typeof response.data === 'object') {
          if (response.status >= 200 && response.status < 300) {
            console.log('Response status indicates success, returning data as-is with warning');

            // Return whatever we got, but try to normalize it first
            const normalizedData = this.normalizeProductResponse(response.data);
            return normalizedData;
          }
        }

        return null;
      }
    } catch (error) {
      console.error('Error creating product:', error);
      if (axios.isAxiosError(error) && error.response) {
        console.error('Server response:', error.response.status, error.response.data);
      }
      throw error;
    }
  }

  async updateProduct(id: number, product: Partial<Product>): Promise<Product | null> {
    try {
      console.log(`🔄 Updating product ${id} with data:`, JSON.stringify(product, null, 2));

      // Lấy thông tin product hiện tại để kết hợp với dữ liệu cập nhật
      const currentProduct = await this.getProductById(id);
      if (!currentProduct) {
        throw new Error(`Product with ID ${id} not found`);
      }

      console.log('📊 Current product data from server:', currentProduct);

      // Kết hợp dữ liệu hiện tại với dữ liệu cập nhật
      const mergedProduct = {
        ...currentProduct,
        ...product
      };

      console.log('🔄 Merged product data:', mergedProduct);

      // Chuẩn hóa imgSlider nếu có
      if (mergedProduct.imgSlider && Array.isArray(mergedProduct.imgSlider)) {
        const normalizedImgUrls = mergedProduct.imgSlider.map(img => {
          if (typeof img === 'object' && img !== null && 'url' in img) {
            return (img as any).url;
          }
          return img;
        });

        console.log('📊 Normalized imgSlider for update:', normalizedImgUrls);

        // Cập nhật lại product.imgSlider với dạng đã chuẩn hóa
        mergedProduct.imgSlider = normalizedImgUrls;
      }

      // Kiểm tra xem có cần cập nhật images không
      const hasImageUpdates = mergedProduct.imgSlider && mergedProduct.imgSlider.length > 0;

      // Chuẩn bị dữ liệu product cho cập nhật (bao gồm tất cả trường)
      const productData = this.prepareProductData(mergedProduct);

      // Kiểm tra lại có đủ các trường bắt buộc không
      if (!productData.name || !productData.slug || !productData.categoryId) {
        console.error('❌ Missing required fields for update:', {
          name: productData.name,
          slug: productData.slug,
          categoryId: productData.categoryId
        });
      }

      console.log('📤 Final update data:', JSON.stringify(productData, null, 2));

      // FINAL JSON CHECK before API call - Stringify and parse to detect any circular structures
      const jsonCheck = JSON.stringify(productData);
      const parsedData = JSON.parse(jsonCheck);

      // Check imgSlider format in the parsed data one last time
      if (parsedData.imgSlider && !parsedData.imgSlider.every((url: any) => typeof url === 'string')) {
        console.error('❌ FINAL CHECK FAILED: imgSlider still contains non-string values after JSON stringify!');
        parsedData.imgSlider = parsedData.imgSlider.map((item: any) =>
          typeof item === 'object' && item !== null && item.url ? item.url : String(item)
        );
        console.log('🛠️ Fixed imgSlider at final step:', parsedData.imgSlider);
      }

      // Ensure image_url is string
      if (parsedData.image_url && typeof parsedData.image_url === 'object') {
        parsedData.image_url = parsedData.imgSlider[0] || '';
      }

      // Use the cleaned data
      const fixedData = parsedData;

      try {
        // Gửi yêu cầu cập nhật
        const response = await api.put(`${this.apiUrl}/${id}`, fixedData);
        let updatedProduct = response.data.product || response.data.data || response.data;

        // Normalize the product to ensure correct format
        updatedProduct = this.normalizeProductResponse(updatedProduct);

        return updatedProduct;
      } catch (error) {
        console.error('❌ Error updating product:', error);
        if (axios.isAxiosError(error) && error.response) {
          console.error('❌ Server response:', error.response.status, error.response.data);
          throw new Error(error.response.data?.message || error.response.data?.error || `Failed with status: ${error.response.status}`);
        }
        throw error;
      }
    } catch (error) {
      console.error(`Error updating product ${id}:`, error);
      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data?.message || error.response.data?.error || `Failed with status: ${error.response.status}`);
      }
      throw error;
    }
  }

  async deleteProduct(id: number): Promise<boolean> {
    try {
      const response = await api.delete(`${this.apiUrl}/${id}`);
      return response.status >= 200 && response.status < 300;
    } catch (error) {
      console.error(`Error deleting product ${id}:`, error);
      return false;
    }
  }

  async getCategories(): Promise<any[]> {
    try {
      console.log('Fetching categories from:', this.categoryUrl);

      const response = await api.get(this.categoryUrl);
      const data = response.data;
      const categories = data.categories || data.items || data.data || data;

      console.log('Categories fetched:', categories);

      // Chuẩn hóa dữ liệu trả về
      return Array.isArray(categories) ? categories.map((category: any) => ({
        id: category.id || category.ID,
        name: category.name || category.Name,
        slug: category.slug || category.Slug,
        description: category.description || category.Description || '',
        image_url: category.image_url || category.imageUrl || category.ImageURL || '',
      })) : [];
    } catch (error) {
      console.error('Error fetching categories:', error);
      // Trả về mảng trống nhưng thêm một vài category mẫu để test UI
      return [
        { id: 1, name: "Electronics", slug: "electronics" },
        { id: 2, name: "Clothing", slug: "clothing" },
        { id: 3, name: "Home & Kitchen", slug: "home-kitchen" },
        { id: 4, name: "Books", slug: "books" },
        { id: 5, name: "Sports & Outdoors", slug: "sports-outdoors" }
      ];
    }
  }

  private prepareProductData(product: Partial<Product>): any {
    console.log('🔍 prepareProductData input:', product);

    // Kiểm tra fields bắt buộc
    if (!product.name) {
      console.error('❌ Missing required field: name');
    }
    if (!product.slug) {
      console.warn('⚠️ Missing slug, will be generated from name');
    }
    if (!product.categoryId || product.categoryId === '0') {
      console.error('❌ Missing required field: categoryId');
    }
    if (typeof product.price === 'undefined' || product.price === null) {
      console.error('❌ Missing required field: price');
    }

    // Convert frontend product model to backend expected format
    const productData: any = {
      // Các trường cơ bản (basic fields) - đảm bảo trường bắt buộc có giá trị mặc định
      name: product.name || 'Unnamed Product', // Trường bắt buộc, không được để trống
      description: product.description || '',
      type: product.type || 'normal',
      price: parseFloat(String(product.price || 0)),
      originalPrice: parseFloat(String(product.originalPrice || 0)),
      discount: parseFloat(String(product.discount || 0)),
      slug: product.slug || this.slugify(product.name || 'unnamed-product'),
      categoryId: parseInt(String(product.categoryId || 0)),
      categorySlug: product.categorySlug || '',
      stockQuantity: typeof product.stockQuantity === 'string' && product.stockQuantity === 'Unlimited'
        ? 999999
        : parseInt(String(product.stockQuantity || 0)),
      brand: product.brand || '',
      tags: product.tags || [],
      features: product.features || [],
    };

    // Xử lý shipping_info đúng format (Process shipping_info in correct format)
    if (product.shippingInfo) {
      productData.shippingInfo = {
        courier: product.shippingInfo.courier || '',
        local: product.shippingInfo.local || '',
        ups: product.shippingInfo.ups || '',
        global: product.shippingInfo.global || '',
      };
    } else {
      productData.shippingInfo = {
        courier: '',
        local: '',
        ups: '',
        global: ''
      };
    }

    // Xử lý UI metadata (Process UI metadata)
    if (product.uiMetadata) {
      try {
        const jsonStr = JSON.stringify(product.uiMetadata);
        productData.uiMetadata = JSON.parse(jsonStr);
      } catch (e) {
        console.warn('Invalid UI metadata format', e);
        productData.uiMetadata = {};
      }
    }

    // Xử lý hình ảnh (Process images)
    if (product.imgSlider && Array.isArray(product.imgSlider)) {
      // 1. Đảm bảo imgSlider luôn là mảng string URLs đơn giản
      const imgUrls = product.imgSlider.map(img => {
        if (typeof img === 'object' && img !== null && img && 'url' in img) {
          return (img as any).url;
        }
        return img;
      });

      // 2. Đặt imgSlider là mảng URLs đơn giản
      productData.imgSlider = imgUrls;

      // 3. Đặt image_url là URL đầu tiên (string đơn giản)
      productData.image_url = imgUrls.length > 0 ? imgUrls[0] : '';

      // 4. Tạo mảng images cho backend nếu cần
      productData.images = imgUrls.map((url, index) => ({
        url: url,
        is_primary: index === 0,
        display_order: index,
        // Only include product_id for existing products (updates), not for new products
        ...(product.id ? { product_id: product.id } : {})
      }));
    } else {
      // Đảm bảo các trường hình ảnh luôn được khởi tạo
      productData.images = [];
      productData.imgSlider = [];
      productData.image_url = '';
    }

    // Add debugging info
    console.log('🔄 Prepared product data details:');
    console.log('  • name:', productData.name, typeof productData.name);
    console.log('  • description:', productData.description ? 'set' : 'empty');
    console.log('  • categoryId:', productData.categoryId, typeof productData.categoryId);
    console.log('  • categorySlug:', productData.categorySlug, typeof productData.categorySlug);
    console.log('  • price:', productData.price, typeof productData.price);
    console.log('  • slug:', productData.slug, typeof productData.slug);
    console.log('  • images:', productData.images?.length || 0, 'objects');
    console.log('  • imgSlider:', productData.imgSlider?.length || 0, 'URLs (string array)');
    console.log('  • image_url:', typeof productData.image_url, productData.image_url || 'Not set');

    // Kiểm tra imgSlider có đúng là mảng string URLs
    if (productData.imgSlider && productData.imgSlider.length > 0) {
      const isAllStrings = productData.imgSlider.every((item: any) => typeof item === 'string');
      console.log('  • imgSlider format check:', isAllStrings ? '✅ All strings' : '❌ Contains non-string items');
    }

    return productData;
  }

  // Helper method to create slugs from product names
  private slugify(text: string): string {
    return text
      .toString()
      .toLowerCase()
      .replace(/\s+/g, '-')           // Replace spaces with -
      .replace(/[^\w\-]+/g, '')       // Remove all non-word chars
      .replace(/\-\-+/g, '-')         // Replace multiple - with single -
      .replace(/^-+/, '')             // Trim - from start of text
      .replace(/-+$/, '');            // Trim - from end of text
  }

  // Phương thức patch chỉ cập nhật các trường cụ thể thay vì toàn bộ sản phẩm
  async patchProduct(id: number, fields: Partial<Product>): Promise<Product | null> {
    try {
      console.log(`🔄 Attempting to patch product ${id} with specific fields:`, JSON.stringify(fields, null, 2));

      // Ensure ID is valid
      if (!id || isNaN(Number(id))) {
        throw new Error('Invalid product ID for patch operation');
      }

      // Chuẩn hóa imgSlider nếu có (normalize imgSlider if present)
      if (fields.imgSlider && Array.isArray(fields.imgSlider)) {
        fields.imgSlider = fields.imgSlider.map(img => {
          if (typeof img === 'object' && img !== null && 'url' in img) {
            return (img as any).url;
          }
          return img;
        });

        // Nếu có imgSlider, tự động cập nhật image_url là ảnh đầu tiên (if imgSlider exists, automatically update image_url with first image)
        if (fields.imgSlider.length > 0) {
          fields.image_url = fields.imgSlider[0];
        }

        console.log('📊 Normalized imgSlider for update:', fields.imgSlider);
      }

      // Chuẩn bị các trường cần cập nhật (prepare fields for update)
      const patchData: any = {};

      // Xử lý từng trường riêng lẻ với đúng định dạng (process each field with correct format)
      Object.entries(fields).forEach(([key, value]) => {
        switch (key) {
          case 'categoryId':
            if (value !== undefined) patchData[key] = parseInt(String(value));
            break;
          case 'price':
          case 'originalPrice':
          case 'discount':
            if (value !== undefined) patchData[key] = parseFloat(String(value));
            break;
          case 'stockQuantity':
            if (value !== undefined) {
              patchData[key] = typeof value === 'string' && value === 'Unlimited'
                ? 999999
                : parseInt(String(value));
            }
            break;
          case 'imgSlider':
            // imgSlider đã được chuẩn hóa ở trên (imgSlider has been normalized above)
            if (value && Array.isArray(value)) {
              patchData[key] = value;

              // If we're updating images, also update the primary image_url
              if (value.length > 0) {
                patchData.image_url = value[0];
              }
            }
            break;
          default:
            // Các trường khác giữ nguyên (keep other fields as is)
            if (value !== undefined) {
              patchData[key] = value;
            }
        }
      });

      console.log('📤 PATCH data prepared:', JSON.stringify(patchData, null, 2));

      try {
        // Thử gửi yêu cầu PATCH (try PATCH request)
        console.log('🔄 Sending PATCH request to:', `${this.apiUrl}/${id}`);
        const response = await api.patch(`${this.apiUrl}/${id}`, patchData);
        console.log('✅ PATCH request successful, response:', response.data);

        // Extract product data from response
        let updatedProduct = response.data.product || response.data.data || response.data;

        // Normalize the product to ensure correct format
        updatedProduct = this.normalizeProductResponse(updatedProduct);

        console.log('✅ Normalized updated product:', updatedProduct);
        return updatedProduct;
      } catch (patchError) {
        // Nếu PATCH không được hỗ trợ, chuyển sang sử dụng PUT (if PATCH not supported, fall back to PUT)
        console.warn('⚠️ PATCH method failed:', patchError);

        if (axios.isAxiosError(patchError) && patchError.response) {
          console.error('⚠️ Server response:', patchError.response.status, patchError.response.data);

          // If it's a 405 Method Not Allowed or other error suggesting PATCH isn't supported
          if (patchError.response.status === 405 || patchError.response.status === 501) {
            console.log('🔄 Falling back to PUT method (PATCH not supported)...');
            return this.updateProduct(id, fields);
          }

          // For other errors, throw with the server message
          throw new Error(patchError.response.data?.message ||
            patchError.response.data?.error ||
            `Server error: ${patchError.response.status}`);
        }

        // For non-Axios errors or if we can't extract a message
        console.log('🔄 Falling back to PUT method due to unknown error...');
        return this.updateProduct(id, fields);
      }
    } catch (error) {
      console.error(`❌ Error patching product ${id}:`, error);

      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data?.message ||
          error.response.data?.error ||
          `Failed with status: ${error.response.status}`);
      }
      throw error;
    }
  }
}

export const productService = new ProductService();