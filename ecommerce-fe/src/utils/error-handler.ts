/**
 * Error message handler utility
 * Converts technical error messages to user-friendly messages
 * Includes Vietnamese translations for B1 level vocabulary
 */

import { AxiosError } from "axios";

/**
 * Maps common RPC error messages to user-friendly messages with Vietnamese translations
 * @param errorMsg - The original error message
 * @returns A user-friendly error message
 */
export const formatErrorMessage = (errorMsg: string): string => {
  if (!errorMsg) return 'An unknown error occurred'; // Một lỗi không xác định đã xảy ra

  // Handle RPC errors
  if (errorMsg.includes('rpc error:')) {
    // Authentication errors
    if (errorMsg.includes('code = Unauthenticated') && errorMsg.includes('invalid credentials')) {
      return 'Invalid login information. Please check your email and password.'; // Thông tin đăng nhập không hợp lệ. Vui lòng kiểm tra email và mật khẩu của bạn.
    }
    if (errorMsg.includes('code = Unauthenticated')) {
      return 'You need to log in to use this feature.'; // Bạn cần đăng nhập để sử dụng tính năng này.
    }

    // Authorization errors
    if (errorMsg.includes('code = PermissionDenied')) {
      return 'You do not have permission to perform this action.'; // Bạn không có quyền thực hiện hành động này.
    }

    // Not found errors
    if (errorMsg.includes('code = NotFound')) {
      return 'The requested information could not be found.'; // Không tìm thấy thông tin bạn yêu cầu.
    }

    // Invalid argument errors
    if (errorMsg.includes('code = InvalidArgument')) {
      return 'Invalid information. Please check your inputs.'; // Thông tin không hợp lệ. Vui lòng kiểm tra đầu vào của bạn.
    }

    // Already exists errors
    if (errorMsg.includes('code = AlreadyExists')) {
      if (errorMsg.includes('email')) {
        return 'This email is already in use. Please use a different email.'; // Email này đã được sử dụng. Vui lòng sử dụng email khác.
      }
      return 'This information already exists in the system.'; // Thông tin này đã tồn tại trong hệ thống.
    }

    // Aborted errors
    if (errorMsg.includes('code = Aborted')) {
      return 'The operation was cancelled. Please try again.'; // Thao tác đã bị hủy. Vui lòng thử lại.
    }

    // Deadline exceeded errors
    if (errorMsg.includes('code = DeadlineExceeded')) {
      return 'The operation timed out. Please check your connection and try again.'; // Thao tác đã hết thời gian. Vui lòng kiểm tra kết nối của bạn và thử lại.
    }

    // Internal errors
    if (errorMsg.includes('code = Internal')) {
      return 'A system error occurred. Please try again later.'; // Đã xảy ra lỗi hệ thống. Vui lòng thử lại sau.
    }

    // Unavailable errors
    if (errorMsg.includes('code = Unavailable')) {
      return 'The service is currently unavailable. Please try again later.'; // Dịch vụ hiện không khả dụng. Vui lòng thử lại sau.
    }

    // Generic RPC error
    return 'An error occurred while connecting to the server. Please try again.'; // Đã xảy ra lỗi khi kết nối đến máy chủ. Vui lòng thử lại.
  }

  // Handle common domain-specific errors
  if (errorMsg.includes('email already exists')) {
    return 'This email is already in use. Please use a different email.'; // Email này đã được sử dụng. Vui lòng sử dụng email khác.
  }

  if (errorMsg.includes('password')) {
    if (errorMsg.includes('too short')) {
      return 'Password is too short. Please use at least 8 characters.'; // Mật khẩu quá ngắn. Vui lòng sử dụng ít nhất 8 ký tự.
    }
    if (errorMsg.includes('too weak')) {
      return 'Password is too weak. Please combine letters, numbers, and special characters.'; // Mật khẩu quá yếu. Vui lòng kết hợp chữ cái, số và ký tự đặc biệt.
    }
    if (errorMsg.includes('incorrect')) {
      return 'Incorrect password. Please check again.'; // Mật khẩu không chính xác. Vui lòng kiểm tra lại.
    }
    return 'There is an issue with your password. Please check again.'; // Có vấn đề với mật khẩu của bạn. Vui lòng kiểm tra lại.
  }

  // Payment errors
  if (errorMsg.includes('payment') || errorMsg.includes('card')) {
    if (errorMsg.includes('declined')) {
      return 'Your payment was declined. Please try a different payment method.'; // Thanh toán của bạn bị từ chối. Vui lòng thử phương thức thanh toán khác.
    }
    if (errorMsg.includes('expired')) {
      return 'Your card has expired. Please update your payment information.'; // Thẻ của bạn đã hết hạn. Vui lòng cập nhật thông tin thanh toán.
    }
    return 'There was an issue with your payment. Please try again.'; // Đã xảy ra sự cố với thanh toán của bạn. Vui lòng thử lại.
  }

  // Product errors
  if (errorMsg.includes('product')) {
    if (errorMsg.includes('out of stock')) {
      return 'This product is out of stock. Please check back later.'; // Sản phẩm này đã hết hàng. Vui lòng quay lại sau.
    }
    if (errorMsg.includes('no longer available')) {
      return 'This product is no longer available.'; // Sản phẩm này không còn nữa.
    }
    return 'There was an issue with the product. Please try again.'; // Đã xảy ra sự cố với sản phẩm. Vui lòng thử lại.
  }

  // Network errors
  if (errorMsg.includes('network') || errorMsg.includes('connection')) {
    return 'Network connection error. Please check your internet and try again.'; // Lỗi kết nối mạng. Vui lòng kiểm tra internet của bạn và thử lại.
  }

  // Server errors
  if (errorMsg.includes('server') || errorMsg.includes('500')) {
    return 'Server error. Our team has been notified. Please try again later.'; // Lỗi máy chủ. Đội ngũ của chúng tôi đã được thông báo. Vui lòng thử lại sau.
  }

  // Return the original message if no specific formatting is found
  return errorMsg;
};

/**
 * Extracts meaningful error message from various error objects
 * @param error - The error object from API calls
 * @returns A formatted error message
 */
export const extractErrorMessage = (error: unknown): string => {
  if (!error) return 'An unknown error occurred'; // Một lỗi không xác định đã xảy ra

  // Handle Error objects
  if (error instanceof Error) {
    return formatErrorMessage(error.message);
  }

  // Handle Axios errors
  if (error instanceof AxiosError && error.response && error.response.data) {
    const responseData = error.response.data as { message?: string; error?: string };

    // Try to extract error message from response data
    if (responseData.message) {
      return formatErrorMessage(responseData.message);
    }

    if (responseData.error) {
      return formatErrorMessage(responseData.error);
    }

    // For more complex error responses
    if (typeof responseData === 'object') {
      return 'An error occurred. Please try again.'; // Đã xảy ra lỗi. Vui lòng thử lại.
    }

    // If response data is a string
    if (typeof responseData === 'string') {
      return formatErrorMessage(responseData);
    }
  }

  // Handle string errors
  if (typeof error === 'string') {
    return formatErrorMessage(error);
  }

  // Handle other cases
  return 'An unknown error occurred. Please try again.'; // Một lỗi không xác định đã xảy ra. Vui lòng thử lại.
};

// Export default for convenience
export default {
  formatErrorMessage,
  extractErrorMessage
}; 