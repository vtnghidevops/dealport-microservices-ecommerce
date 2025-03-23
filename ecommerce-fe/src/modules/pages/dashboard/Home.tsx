import React from "react";
import HeroBanner from "../../components/sections/HeroBanner";
import NewFashion from "../../components/sections/NewFashion";
import GamingBanner from "../../components/sections/GamingBanner";
import DisplayGrid from "../../components/sections/DisplayGrid";
import BannerShowCase from "../../components/sections/BannerShowcase";
import { TrendingProducts } from "../../components/sections/TrendingProducts";
import CategoryExplorer from "../../components/sections/CategoryExplorer";
import { Category } from "../../models/category.models";
import TopProducts from "../../components/sections/TopProducts";
import { LimitedDeal } from "../../components/sections/LimitedDeal";
import HappyCustomers from "../../components/sections/HappyCustomers";
const urlImgProduct = "images/products/";
const urlImgGaming = "images/gaming/";
const urlImgDisplay = "images/display/";
const urlImgBanner = "images/banner/";
const urlImgTrending = "images/trending/";
const urlImgExploring = "images/exploring/";
const urlImgTop = "images/bestselling/";
const urlLimited = "images/limited/" ;

const products_trending = [
  {
    id: 1,
    title: "Radiant Glow Hydrating Serum",
    description:
      "Gentle yet effective, our Radiant Boosting Foaming our Radiant Boosting Foaming our Radiant Boosting Foaming",
    price: 29.99,
    originalPrice: 39.99,
    discount: 20,
    review: {
      rating: 4.8,
      count: 345,
    },
    imageUrl: `${urlImgTrending}serum.png`,
  },
  {
    id: 2,
    title: "Modern Minimalist Vase",
    description:
      "Track your workouts, heart rate, sleep quality and receive notifications. Water resistant up to 50m with 7-day battery life.",
    price: 40.99,
    originalPrice: 0,
    discount: 0,
    review: {
      rating: 4.6,
      count: 842,
    },
    imageUrl: `${urlImgTrending}vase.png`,
  },
  {
    id: 3,
    title: "FitPro 3000 Smart Watch",
    description:
      "Fast-charging power bank with dual USB ports and USB-C compatibility. Charge multiple devices simultaneously on the go.",
    price: 119.99,
    originalPrice: 0,
    discount: 0,
    review: {
      rating: 4.0,
      count: 2105,
    },
    imageUrl: `${urlImgTrending}smartwatch.png`,
  },
];
// Exploring category
const categories: Category[] = [
  {
    id: 1,
    name: "Grocery",
    image: `${urlImgExploring}grocery.png`,
  },
  {
    id: 2,
    name: "Home",
    image: `${urlImgExploring}/home.png`,
  },
  {
    id: 3,
    name: "Fashion",
    image: `${urlImgExploring}/fashion.png`,
  },
  {
    id: 4,
    name: "Electronic",
    image: `${urlImgExploring}electronic.png`,
  },
  {
    id: 5,
    name: "Toys",
    image: `${urlImgExploring}toys.png`,
  },
  {
    id: 6,
    name: "Grocery",
    image: `${urlImgExploring}grocery.png`,
  },
];
const handleCategoryClick = (category: Category) => {
  console.log("Category clicked:", category);
  // Navigate or perform actions
};

const handleViewAll = () => {
  console.log("View all categories clicked");
  // Navigate to all categories page
};

// Top Products selling
const featuredProducts = [
  {
    id: 1,
    name: "Computer Accessories",
    image: `${urlImgTop}computer_accessories.png`,
    setUpDesign: "row",
    isCommingSoon: true,
    gridSpan: {
      col: 1,
      row: 1
    }
  },
  {
    id: 2,
    name: "Men's Casual Outfit",
    price: "200",
    image: `${urlImgTop}football.png`,
    actionLabel: "Visit store",
    setUpDesign: "row",
    gridSpan: {
      col: 1,
      row: 1
    }
  },
  {
    id: 3,
    name: "Pome Granate Juice",
    price: "49",
    image: `${urlImgTop}juice.png`,
    actionLabel: "Buy now",
    setUpDesign: "row",
    gridSpan: {
      col: 1,
      row: 1
    }
    
  },
  {
    id: 4,
    name: "Dog Food Made With Love",
    image: `${urlImgTop}dog_food.png`,
    badge: "15% OFF",
    actionLabel: "Shop Now",
    setUpDesign: "col",
    gridSpan: {
      col: 1,
      row: 2
    },
  },
  {
    id: 5,
    name: "Security Camera System",
    price: "420",
    image: `${urlImgTop}security_camera_left.png`,
    actionLabel: "Visit store",
    setUpDesign: "double",
    gridSpan: {
      col: 2,
      row: 1
    },
    isImageDouble: true,
    image_double: `${urlImgTop}security_camera_right.png`,
  },
  {
    id: 6,
    name: "Premium Cosmetic Set",
    price: "149",
    image: `${urlImgTop}skincare.png`,
    actionLabel: "Buy Now",
    setUpDesign: "row",
    gridSpan: {
      col: 1,
      row: 1
    }

  }
];

const handleProductClick = (product, index) => {
  console.log("Product clicked:", product);
  // Điều hướng đến trang chi tiết sản phẩm
};

// Limited Time Deal
const limitedProducts= [
  {
    id: 1,
    title: "Samsung Galaxy S24",
    description: "Gentle yet effective, our hydrating serum infuses skin with essential moisture while brightening and plumping for a radiant complexion.",
    price: 29.99,
    originalPrice: 39.99,
    discount: 25,
    review: {
      rating: 4.8,
      count: 345,
    },
    imageUrl: `${urlLimited}samsung_s24.png`,
  },
  {
    id: 2,
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    imageUrl: `${urlLimited}earbud.png`,
  },
  {
    id: 3,
    title: "Winter fashion jacket",
    description: "Delivers extreme hydration and helps strengthen the skin barrier. Perfect for both daytime wear and overnight rejuvenation. [[2]]",
    price: 39.99,
    originalPrice: 49.99,
    discount: 20,
    review: {
      rating: 4.9,
      count: 412,
    },
    imageUrl: `${urlLimited}winter_jacket.png`,
  },
  {
    id: 4,
    title: "New Balance 574 Senekers",
    description: "Illuminating serum infused with glow-boosting intelligent botanicals that leave the skin visibly radiant while reducing redness. [[4]]",
    price: 44.99,
    originalPrice: 59.99,
    discount: 25,
    review: {
      rating: 4.6,
      count: 278,
    },
    imageUrl: `${urlLimited}senekers.png`,
  },
  {
    id: 5,
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    imageUrl: `${urlLimited}earbud.png`,
  },
  {
    id: 6,
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    imageUrl: `${urlLimited}earbud.png`,
  },
  {
    id: 7,
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    imageUrl: `${urlLimited}earbud.png`,
  }
];

export default function Home() {
  return (
    <div>
      <HeroBanner></HeroBanner>
      <div>
        {/* Ads */}
        <div className="relative flex items-center justify-center w-full h-[550px]">
          <div className="absolute left-[3%] top-0">
            <NewFashion
              title={"New Year! New Fashion"}
              imgURL={`${urlImgProduct}new-fashion.png`}
              ButtonType={"secondary"}
            ></NewFashion>
          </div>
          <div className="top-0 absolute left-[37%]">
            <GamingBanner
              title="Gaming accessories"
              items={[
                {
                  subtitle: "Headsets",
                  imgURL: `${urlImgGaming}headsets.png`,
                  href: "#",
                },
                {
                  subtitle: "Mouse",
                  imgURL: `${urlImgGaming}mouse.png`,
                  href: "#",
                },
                {
                  subtitle: "Controller",
                  imgURL: `${urlImgGaming}controller.png`,
                  href: "#",
                },
                {
                  subtitle: "Chair",
                  imgURL: `${urlImgGaming}chair.png`,
                  href: "#",
                },
              ]}
            ></GamingBanner>
          </div>
          <div className="absolute left-[71%] top-0 ">
            {/* image, title, price, discount_img, ctaText, buttonType */}
            <DisplayGrid
              products={[
                {
                  image: `${urlImgDisplay}be-winner.png`,
                },
                {
                  image: `${urlImgDisplay}redmi-y3.png`,
                  buttonType: "gradient",
                },
                {
                  image: `${urlImgDisplay}ambilighttv.png`,
                  buttonType: "secondary",
                  price: "750.99",
                  title: "Philips 4K Ambilight TV",
                  discount_img: `${urlImgDisplay}discount_img.png`,
                },
              ]}
            ></DisplayGrid>
          </div>
          <div className="absolute w-full top-[58%]">
            {/* img, buttonType, buttonText, hasMore */}
            <BannerShowCase
              categories={[
                {
                  image: `${urlImgBanner}trousers_fashion.png`,
                  buttonType: "",
                  buttonText: "",
                  hasMore: false,
                },
                {
                  image: `${urlImgBanner}watchmen_fashion.png`,
                  buttonType: "gray",
                  buttonText: "Shop Now",
                  hasMore: false,
                },
                {
                  image: `${urlImgBanner}denim_fashion.png`,
                  buttonType: "",
                  buttonText: "",
                  hasMore: true,
                },
                {
                  image: `${urlImgBanner}dometic.png`,
                  buttonType: "black",
                  buttonText: "Shop Now",
                  hasMore: false,
                },
              ]}
            />
          </div>
        </div>

        {/* Trending Products */}
        <div className="relative w-full h-max pl-[6rem]">
          <TrendingProducts products_trending={products_trending} />
        </div>

        {/* Category */}
        <div className="relative w-full pl-[6rem] pr-[4.3rem] mt-[3rem]">
          <CategoryExplorer
            categories={categories}
            onViewAllClick={handleViewAll}
            onItemClick={handleCategoryClick}
          />
        </div>

        {/* Best selling */}
        <div className="relative w-full pl-[6rem] pr-[4.3rem] mt-[3rem]">
          <TopProducts
            products={featuredProducts}
            onViewAllClick={handleViewAll}
            onItemClick={handleProductClick}
          />
        </div>

        {/* Limited time deal */}
        <div className="relative w-full pl-[6rem] pr-[4.3rem] mt-[3rem]">
            <LimitedDeal products={limitedProducts}></LimitedDeal>
        </div>

        {/* Happy Customers */}
        <div>
            <HappyCustomers></HappyCustomers>
        </div>
      </div>
    </div>
  );
}
