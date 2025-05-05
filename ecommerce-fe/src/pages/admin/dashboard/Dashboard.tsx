import React, { useEffect, useState } from "react";
import AdminHeader from "../../../components/admin/layout/AdminHeader";

// Dashboard Components
import { DashboardCard } from "../../../components/admin/dashboard/cards";
import WeeklyReportChart from "../../../components/admin/dashboard/charts/WeeklyReportChart";
import UserActivityChart from "../../../components/admin/dashboard/charts/UserActivityChart";
import TransactionTable from "../../../components/admin/dashboard/tables/TransactionTable";
import TopProductsTable from "../../../components/admin/dashboard/cards/TopProducts";
import BestSellingTable from "../../../components/admin/dashboard/tables/BestSellingTable";
import ProductCategoryList from "../../../components/admin/dashboard/widgets/ProductCategoryList";
import NewProductList from "../../../components/admin/dashboard/widgets/NewProductList";

// Services
import { DashboardService } from "../../../components/admin/dashboard/services/dashboard.service";
import { TransactionService } from "../../../components/admin/dashboard/services/transaction.service";
import { ProductDashboardService } from "../../../components/admin/dashboard/services/product.service";
import { CategoryService } from "../../../services/product/product.service";

// Models
import { DashboardSummary } from "../../../components/admin/dashboard/models/dashboard.model";
import { Transaction } from "../../../components/admin/dashboard/models/transaction.model";
import { BestSellingProductStats } from "../../../components/admin/dashboard/services/product.service";
import { Category } from "@/types/category.model";
import { Product } from "@/types/product.model";

// Icons
import { HiOutlineDotsVertical } from "react-icons/hi";
import { IoFilterSharp } from "react-icons/io5";
import { CiCirclePlus } from "react-icons/ci";

const DashboardAdmin: React.FC = () => {
  const [loading, setLoading] = useState<boolean>(true);
  const [dashboardData, setDashboardData] = useState<DashboardSummary | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [bestSellingProducts, setBestSellingProducts] = useState<BestSellingProductStats[]>([]);
  const [productCategories, setProductCategories] = useState<Category[]>([]);
  const [newProducts, setNewProducts] = useState<Product[]>([]);
  const [topProducts, setTopProducts] = useState<BestSellingProductStats[]>([]);

  useEffect(() => {
    const fetchDashboardData = async () => {
      setLoading(true);
      try {
        // Fetch all dashboard data in parallel
        const [
          dashboardSummary,
          transactionData,
          topSellingProducts,
          categories,
          latestProducts,
        ] = await Promise.all([
          DashboardService.getDashboardSummary(),
          TransactionService.getTransactions(1, 5),
          ProductDashboardService.getTopSellingProducts(5),
          CategoryService.getAllCategories(),
          ProductDashboardService.getNewProducts(5),
        ]);

        // Update state with real data
        setDashboardData(dashboardSummary);
        setTransactions(transactionData.transactions);
        setBestSellingProducts(topSellingProducts);
        setTopProducts(topSellingProducts);
        setProductCategories(categories);
        setNewProducts(latestProducts);

        console.log('Dashboard data loaded successfully', {
          dashboardSummary,
          transactions: transactionData.transactions,
          bestSellingProducts: topSellingProducts,
          categories,
          latestProducts
        });
      } catch (error) {
        console.error("Error fetching dashboard data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  if (loading || !dashboardData) {
    return (
      <div className="flex justify-center items-center h-screen">
        <div className="animate-spin rounded-full h-24 w-24 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="flex bg-neutral-50 ">
      {/* <Sidebar isOpen={true} /> */}
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Dashboard" />
        <main className="p-[1rem]">
          {/* Metric Cards */}
          <div className="flex flex-col md:flex-row gap-[18px] mb-6">
            <DashboardCard
              type="sales"
              title="Total Sales"
              amount={dashboardData?.totalSales?.amount || 0}
              currency={dashboardData.totalSales.currency}
              percentChange={dashboardData.totalSales.percentChange}
              previousAmount={dashboardData.totalSales.previousAmount}
              lastDays={dashboardData.totalSales.lastDays}
            />

            <DashboardCard
              type="orders"
              title="Total Orders"
              count={dashboardData.totalOrders.count}
              percentChange={dashboardData.totalOrders.percentChange}
              previousCount={dashboardData.totalOrders.previousCount}
              lastDays={dashboardData.totalOrders.lastDays}
            />

            <DashboardCard
              type="pendingCanceled"
              title="Pending & Canceled"
              pending={dashboardData.pendingCanceled.pending}
              canceled={dashboardData.pendingCanceled.canceled}
              lastDays={dashboardData.pendingCanceled.lastDays}
            />
          </div>

          {/* Report and User Stats */}
          <div className="flex flex-col md:flex-row gap-[16px] mb-6 mt-[1.25rem]">
            <div className="md:col-span-2 filter drop-shadow-lg bg-white rounded-lg shadow-sm h-[460px] w-[744px] p-[1.25rem]">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-[18px] font-bold">Báo cáo tuần này (Weekly Report)</h3>
                <div className="flex space-x-2 ">
                  <div className="flex bg-aqua-spring items-center justify-center w-[178px] h-[38px] border rounded-[12px]">
                    <button className="px-[12px] py-[8px] w-[83px] h-[30px] text-xs bg-white text-ocean-green rounded-[12px]">
                      This week
                    </button>
                    <button className="p-[12px] py-[8px] text-xs w-[83px] h-[30px] text-gray-500 rounded-[12px]">
                      Last week
                    </button>
                  </div>
                  <button className="text-gray-400">
                    <HiOutlineDotsVertical></HiOutlineDotsVertical>
                  </button>
                </div>
              </div>

              <div className="grid grid-cols-5 gap-2 mb-[3rem] text-center">
                <div className="p-2">
                  <p className="text-xl font-semibold">{(dashboardData.weeklyReport.customers / 1000).toFixed(1)}k</p>
                  <p className="text-xs text-gray-500">Customers</p>
                </div>
                <div className="p-2">
                  <p className="text-xl font-semibold">{(dashboardData.weeklyReport.totalProducts / 1000).toFixed(1)}k</p>
                  <p className="text-xs text-gray-500">Total Products</p>
                </div>
                <div className="p-2">
                  <p className="text-xl font-semibold">{(dashboardData.weeklyReport.stockProducts / 1000).toFixed(1)}k</p>
                  <p className="text-xs text-gray-500">Stock Products</p>
                </div>
                <div className="p-2">
                  <p className="text-xl font-semibold">{(dashboardData.weeklyReport.outOfStock / 1000).toFixed(1)}k</p>
                  <p className="text-xs text-gray-500">Out of Stock</p>
                </div>
                <div className="p-2">
                  <p className="text-xl font-semibold">{(dashboardData.weeklyReport.revenue / 1000).toFixed(1)}k</p>
                  <p className="text-xs text-gray-500">Revenue</p>
                </div>
              </div>

              <WeeklyReportChart data={dashboardData.weeklyReport.chartData} />
            </div>

            <div className="space-y-6">
              <div className="w-[361px] p-[1rem] bg-white rounded-lg shadow-sm filter drop-shadow-lg">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-sm font-medium text-primary">
                    Người dùng hoạt động (Active Users)
                  </h3>
                  <button className="text-gray-400">
                    <HiOutlineDotsVertical></HiOutlineDotsVertical>
                  </button>
                </div>
                <p className="text-2xl font-bold mb-2">{(dashboardData.userStats.totalInLastMinutes / 1000).toFixed(1)}K</p>
                <p className="text-xs text-gray-500 mb-2 mt-[1rem]">
                  Users per minute
                </p>
                <UserActivityChart
                  data={dashboardData.userStats.userPerMinute}
                />
              </div>
            </div>
          </div>

          {/* Transactions and Top Selling */}
          <div className="flex flex-wrap gap-[16px] mb-6 mt-[1.25rem]">
            <div className="filter drop-shadow-lg bg-white rounded-lg shadow-sm p-[1.25rem] w-[720px] h-[420px]">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-[18px] font-bold">Giao dịch (Transactions)</h3>
                <button className="flex items-center px-3 py-1 text-xs bg-ocean-green text-white rounded-lg w-[85px] h-[32px] p-8">
                  <span className="mr-[0.5rem] text-[14px]">Filter</span>
                  <IoFilterSharp className="w-[18px] h-[18px]"></IoFilterSharp>
                </button>
              </div>
              <TransactionTable transactions={transactions} />
              <div className="flex justify-end mt-[2rem] items-center mr-[2rem]">
                <button className="text-sm bg-white text-primary border border-primary rounded-[25px] w-[96px] h-[32px]">
                  Chi tiết (Details)
                </button>
              </div>
            </div>

            <div className="rounded-lg shadow-sm w-[360px]">
              <TopProductsTable products={topProducts} />
            </div>
          </div>

          {/* Best Selling Products and Add New Product*/}
          <div className="flex flex-wrap gap-[16px] mb-6 mt-[1.25rem]">
            <div className="filter drop-shadow-lg bg-white rounded-lg shadow-sm p-[1.25rem] w-[720px] h-[430px]">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-[18px] font-bold">Sản phẩm bán chạy (Best Selling)</h3>
                <button className="flex items-center px-3 py-1 text-xs bg-ocean-green text-white rounded-lg w-[85px] h-[32px] p-8">
                  <span className="mr-[0.5rem] text-[14px]">Filter</span>
                  <IoFilterSharp className="w-[18px] h-[18px]"></IoFilterSharp>
                </button>
              </div>
              <BestSellingTable bestProducts={bestSellingProducts} />
              <div className="flex justify-end mt-[1rem] items-center mr-[1rem]">
                <button className="text-sm bg-white text-primary border border-primary rounded-[25px] w-[96px] h-[32px]">
                  Chi tiết (Details)
                </button>
              </div>
            </div>

            <div className="bg-white rounded-lg p-[1.25rem] max-w-md w-[360px] filter drop-shadow-lg">
              <div className="flex justify-between mb-6">
                <div>
                  <h3 className="text-[18px] font-bold">Thêm sản phẩm (Add Product)</h3>
                  <p className="text-sm text-gray-500">Add your product with category</p>
                </div>
                <button className="rounded-full h-10 w-10 flex justify-center items-center bg-ocean-green text-white text-[24px]">
                  <CiCirclePlus />
                </button>
              </div>

              <div className="mb-4">
                <h4 className="text-[16px] font-semibold mb-3">Danh mục sản phẩm (Categories)</h4>
                <ProductCategoryList categories={productCategories} />
              </div>

              <div>
                <h4 className="text-[16px] font-semibold mb-3">Sản phẩm mới (New Products)</h4>
                <NewProductList products={newProducts} />
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default DashboardAdmin;
