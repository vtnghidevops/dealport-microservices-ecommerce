//     id: 1,
//     title: "Radiant Glow Hydrating Serum",
//     description:
//       "Gentle yet effective, our Radiant Boosting Foaming our Radiant Boosting Foaming our Radiant Boosting Foaming",
//     price: 29.99,
//     originalPrice: 39.99,
//     discount: 20,
//     review: {
//       rating: 4.8,
//       count: 345,
//     },
//     imageUrl: `${urlImgTrending}serum.png`,

export interface TrendingProductItem {
  id: number;
  title: string;
  description: string;
  price: number | string;
  originalPrice: number | string;
  discount: number | string;
  review: {
    rating: number;
    count: number;
  };
  imageUrl: string;
}

//       id: 1,
//       image: `${urlImgMenCollection}pants.png`, // Replace with your image path
//       price: "25.95",
//       discount: "20% off",
export interface MenCollectionItem {
  id: number;
  image: string;
  price: number;
  discount: number;
}


// id: 1,
// image:  `${urlImgMenCollection}pants.png`, // Replace with your image path
// price: '25.95',
// discount: '20% off',
export interface MenCardItem {
  id: number;
  image: string;
  price: number;
  discount: number;
}



