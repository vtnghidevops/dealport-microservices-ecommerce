import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';

const ScrollToTop = () => {
  const { pathname } = useLocation();

  useEffect(() => {
    // Khi đường dẫn URL thay đổi, cuộn lên đầu trang
    window.scrollTo({
      top: 0,
      behavior: 'smooth' // Hoặc 'auto' nếu bạn muốn cuộn ngay lập tức
    });
  }, [pathname]);

  return null; // Component này không render gì cả
};

export default ScrollToTop;