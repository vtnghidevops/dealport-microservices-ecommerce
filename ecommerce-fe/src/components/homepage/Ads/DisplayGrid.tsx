import React, { useState, useEffect } from "react";
import { ButtonType, getButtonClass } from "@/utils/buttonUtils";
import { DisplayItem } from "./models/ads.model";
import { DisplayService } from "../../../services/product/ads.service";
import { Link } from "react-router-dom";
// import Loading from "@/components/shared/Loading";
import { DisplayGridSkeleton } from "@/components/ui/skeletons";
interface CardMain {
  id?: number;
  imageUrl: string;
  buttonType: string;
  price?: string | number;
  title?: string;
  discountImg?: string;
  type?: string;
  productSlug?: string;
  categorySlug?: string;
}

interface CardFirst {
  imageUrl: string;
  type?: string;
  productSlug?: string;
  categorySlug?: string;
}

interface CardSecond {
  imageUrl: string;
  buttonType: string;
  type?: string;
  productSlug?: string;
  categorySlug?: string;
}

const DisplayGrid: React.FC = () => {
  // fetch data
  const [displayData, setDisplayData] = useState<DisplayItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  // fetch data trending Products
  useEffect(() => {
    const fetchDisplayData = async () => {
      try {
        const data = await DisplayService.getDataDisplay();
        setDisplayData(data);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching: ", error);
        setLoading(false);
      }
    };
    fetchDisplayData();
  }, []);
  //console.log("item displaygrid:", displayData)
  if (loading) {
    return <DisplayGridSkeleton />;
  }

  return (
    <div className="h-[328px] w-[392px] mt-[-5%] text-black rounded-2xl relative">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-4 ">
        {displayData.map((product, index) => (
          <DisplayCard key={index} {...product} />
        ))}
      </div>
    </div>
  );
};
const DisplayCard: React.FC<DisplayItem> = ({
  imageUrl,
  title,
  price,
  discountImg,
  buttonType,
  type,
  productSlug,
  categorySlug,
}) => {
  // Get navigation path based on type and slug
  // const getNavigationPath = () => {
  //   if (!productSlug) return "#";
  //   if (type === "product") {
  //     // If we have a product with categorySlug, use the proper format
  //     if (productSlug) {
  //       return `/category/${categorySlug}/${productSlug}`;
  //     }
  //     // Fallback to products route if no category is available
  //     return `/products/${productSlug}`;
  //   }
  //   // For category type
  //   return `/category/${categorySlug}`;
  // };

  // const navigationPath = getNavigationPath();

  // Kiểm tra các props để quyết định render card nào (Check props to decide which card to render)
  if (discountImg) {
    // Nếu có title và price, render CardMain (card đầy đủ thông tin)
    // (If there is a title and price, render CardMain with full information)
    return (
      <div className="w-[392px] h-[165px] rounded-lg overflow-hidden shadow-md transition-all hover:shadow-xl">
        <CardMain
          imageUrl={imageUrl}
          title={title}
          price={price}
          discountImg={discountImg}
          buttonType={buttonType || ""}
          type={type}
          productSlug={productSlug}
          categorySlug={categorySlug}
        />
      </div>
    );
  } else if (buttonType) {
    // Nếu có buttonType nhưng không có title/price, render CardSecond
    // (If there is buttonType but no title/price, render CardSecond)
    return (
      <div className="rounded-lg h-[146px] w-[188px] overflow-hidden shadow-md transition-all hover:shadow-xl">
        <CardSecond
          imageUrl={imageUrl}
          buttonType={buttonType}
          type={type}
          productSlug={productSlug}
          categorySlug={categorySlug}
        />
      </div>
    );
  } else {
    // Mặc định render CardFirst (chỉ có hình ảnh và link)
    // (Default render CardFirst - only image and link)
    return (
      <div className="rounded-lg h-[146px] w-[188px] overflow-hidden shadow-md transition-all">
        <CardFirst
          imageUrl={imageUrl}
          type={type}
          productSlug={productSlug}
          categorySlug={categorySlug}
        />
      </div>
    );
  }
};

const CardFirst: React.FC<CardFirst> = ({ imageUrl, type, productSlug, categorySlug }) => {
  // Determine navigation path based on type
  let navigationPath = "#";
  if (type === "product" && productSlug) {
    if (categorySlug) {
      navigationPath = `/category/${categorySlug}/${productSlug}`;
    } else {
      navigationPath = `/products/${productSlug}`;
    }
  } else if (type === "category" && categorySlug) {
    navigationPath = `/category/${categorySlug}`;
  }

  return (
    <Link to={navigationPath} className="relative h-[100%] block">
      <img src={imageUrl} className="absolute rounded-xl" alt="Banner" />
      <span className="absolute text-white text-[10px] bottom-0 ml-[1rem] mb-[0.5rem]">
        More Detail
      </span>
    </Link>
  );
};

const CardSecond: React.FC<CardSecond> = ({ imageUrl, buttonType, type, productSlug, categorySlug }) => {
  // Determine navigation path based on type
  let navigationPath = "#";
  if (type === "product" && productSlug) {
    if (categorySlug) {
      navigationPath = `/category/${categorySlug}/${productSlug}`;
    } else {
      navigationPath = `/products/${productSlug}`;
    }
  } else if (type === "category" && categorySlug) {
    navigationPath = `/category/${categorySlug}`;
  }

  return (
    <Link to={navigationPath} className="block relative">
      <img src={imageUrl} alt="Product" />
      <div className="absolute bottom-2 left-1">
        <div className="rounded-3xl">
          <button
            className={`${getButtonClass(
              buttonType as ButtonType || ""
            )} text-white w-[5rem] h-[1.5rem] text-[8px]`}
          >
            Shop Now
          </button>
        </div>
      </div>
    </Link>
  );
};

const CardMain: React.FC<CardMain> = ({
  imageUrl,
  title,
  price,
  discountImg,
  buttonType,
  type,
  productSlug,
  categorySlug
}) => {
  // Determine navigation path based on type
  let navigationPath = "#";
  if (type === "product" && productSlug) {
    if (categorySlug) {
      navigationPath = `/category/${categorySlug}/${productSlug}`;
    } else {
      navigationPath = `/products/${productSlug}`;
    }
  } else if (type === "category" && categorySlug) {
    navigationPath = `/category/${categorySlug}`;
  }

  return (
    <Link to={navigationPath} className="relative h-full w-full block border border-gray-300 shadow-[0_4px_6px_-1px_rgba(0,0,0,0.1),5px_2px_3px_-3px_rgba(0,0,0,0.05),-5px_2px_3px_-3px_rgba(0,0,0,0.05)]">
      <div className="absolute left-[-1rem] max-w-[14rem]">
        <img src={imageUrl} alt={title || "Product"} />
      </div>
      <div className="absolute max-w-[60px] top-[10%] left-[42%]">
        {discountImg && <img src={discountImg} alt="Discount" />}
      </div>
      <div className="absolute top-[27%] right-[8%]">
        <h3>{title}</h3>
        <span>${price}</span>
      </div>
      <div className="absolute right-[7%] block bottom-5">
        <div className="rounded-3xl">
          <button
            className={`${getButtonClass(
              buttonType as ButtonType
            )} w-[7rem] h-[2rem] text-[10px]`}
          >
            Shop Now
          </button>
        </div>
      </div>
    </Link>
  );
};

export default DisplayGrid;
