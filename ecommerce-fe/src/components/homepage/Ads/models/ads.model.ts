
// image: `${urlImgBanner}trousers_fashion.png`,
// buttonType: "",
// buttonText: "",
// hasMore: false,
export interface BannerShowCaseItem {
  id: string,
  image_url: string,
  buttonType?: string,
  buttonText?: string,
  hasMore: boolean
};

// image: `${urlImgDisplay}ambilighttv.png`,
// buttonType: "secondary",
// price: "750.99",
// title: "Philips 4K Ambilight TV",
// discount_img: `${urlImgDisplay}discount_img.png`,
export interface DisplayItem {
  id: string,
  image_url: string,
  buttonType?: string,
  price?: string | number,
  title?: string,
  discount_img?: string
}

// subtitle: "Headsets",
// imgURL: `${urlImgGaming}headsets.png`,
// href: "#",
export interface GamingItem {
  id: string,
  subtitle: string,
  image_url: string,
  href: string
}

// title={"New Year! New Fashion"}
// imgURL={`${urlImgProduct}new-fashion.png`}
// ButtonType={"secondary"}
export interface NewFashionItem {
  id: string,
  title: string,
  image_url: string,
  buttonType: string
}