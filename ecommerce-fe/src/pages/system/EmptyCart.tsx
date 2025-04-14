import { useNavigate } from 'react-router-dom';

const EmptyCart: React.FC = () => {
  const navigate = useNavigate();
  return (
    <div className="min-h-[80vh] bg-white flex items-center justify-center px-5 lg:px-0">
      <div className="w-[415px] flex flex-col items-center justify-center mx-auto">
        {/* Image Container */}
        <div className="w-full mb-16">
          <div className="max-w-[312px] h-[280px] relative flex justify-center items-center mx-auto">
            <img
              src="/images/system/empty_cart.png"
              alt="empty_cart"
              width={312}
              height={160}
              className="object-contain w-full h-full"
            />
          </div>
        </div>

        {/* Content Container */}
        <div className="flex flex-col items-center gap-6 w-full">
          <div className="text-center max-w-[412px] w-full">
            <p className="text-base leading-6 tracking-wider font-sans mb-5">
              You haven't added any items to your wishlist yet. Browse our products and save your favorite items here!
            </p>
          </div>
          <div>
            <button 
              onClick={() => navigate('/')} 
              className="bg-[#8AC732] text-white font-sans w-[146px] h-[48px] rounded-[100px] font-medium text-sm hover:bg-[#7AB52B] transition-colors"
            >
              Go Shopping
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default EmptyCart;