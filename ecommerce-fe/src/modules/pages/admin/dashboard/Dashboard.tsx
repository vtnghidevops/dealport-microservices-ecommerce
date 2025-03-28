import React, { useEffect, useState } from "react";
import Sidebar from "../../../components/admin/layout/Sidebar";
import AdminHeader from "../../../components/admin/layout/AdminHeader";
import { DashboardCard } from "../../../components/admin/dashboard/cards";
import WeeklyReportChart from "../../../components/admin/dashboard/charts/WeeklyReportChart";
import { DashboardService } from "../../../components/sections/AdminDashboard/services/dashboard.service";
import { TransactionService } from "../../../components/sections/AdminDashboard/services/transaction.service";
import { ProductService } from "../../../components/sections/AdminDashboard/services/product.service";
import { DashboardSummary } from "../../../components/sections/AdminDashboard/models/dashboard.model";
import { Transaction } from "../../../components/sections/AdminDashboard/models/transaction.model";
import {
  BestSellingProduct,
  NewProduct,
  ProductCategory,
} from "../../../components/sections/AdminDashboard/models/product.model";
import { HiOutlineDotsVertical } from "react-icons/hi";
import UserActivityChart from "../../../components/admin/dashboard/charts/UserActivityChart";
import TransactionTable from "../../../components/admin/dashboard/tables/TransactionTable";
import { Product } from "../../../components/sections/AdminDashboard/models/product.model";
import { IoFilterSharp } from "react-icons/io5";
import { CiCirclePlus } from "react-icons/ci";
import TopProductsTable from "../../../components/admin/dashboard/cards/TopProducts";
import BestSellingTable from "../../../components/admin/dashboard/tables/BestSellingTable";
import ProductCategoryList from "../../../components/admin/dashboard/widgets/ProductCategoryList";
import NewProductList from "../../../components/admin/dashboard/widgets/NewProductList";
const DashboardAdmin: React.FC = () => {
  const [loading, setLoading] = useState<boolean>(true);
  const [dashboardData, setDashboardData] = useState<DashboardSummary | null>(
    null
  );
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [bestSellingProducts, setBestSellingProducts] = useState<
    BestSellingProduct[]
  >([]);
  const [productCategories, setProductCategories] = useState<ProductCategory[]>(
    []
  );
  const [newProducts, setNewProducts] = useState<NewProduct[]>([]);
  const [topProducts, setTopProducts] = useState<Product[]>([])


  useEffect(() => {
    const fetchDashboardData = async () => {
      setLoading(true);
      try {
        const [
          dashboardSummary,
          transactionData,
          topProducts,
          bestSellingProducts,
          categories,
          latestProducts,
        ] = await Promise.all([
          DashboardService.getDashboardSummary(),
          TransactionService.getTransactions(),
          ProductService.getTopProducts(),
          ProductService.getBestSellingProducts(),
          ProductService.getProductCategories(),
          ProductService.getNewProducts(),
        ]);

        setDashboardData(dashboardSummary);
        setTransactions(transactionData.transactions);
        setTopProducts(topProducts);
        setBestSellingProducts(bestSellingProducts);
        setProductCategories(categories);
        setNewProducts(latestProducts);
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
      <div className="flex border-t-2">
        <Sidebar isOpen={true} />
        <div className="flex-1 flex flex-col ">
          <AdminHeader/>
          <div className="flex-1 flex items-center justify-center">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex bg-gray-50 border-t-2">
      <Sidebar isOpen={true} />
      <div className="flex-1 overflow-auto">
        <AdminHeader />
        <main className="p-[1rem]">
          {/* Metric Cards */}
          <div className="flex flex-col md:flex-row gap-[18px] mb-6">
            <DashboardCard
              type="sales"
              title="Total Sales"
              amount={dashboardData.totalSales.amount}
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
                <h3 className="text-[18px] font-bold">Report for this week</h3>
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
                  <p className="text-xl font-semibold">52k</p>
                  <p className="text-xs text-gray-500">Customers</p>
                </div>
                <div className="p-2">
                  <p className="text-xl font-semibold">3.5k</p>
                  <p className="text-xs text-gray-500">Total Products</p>
                </div>
                <div className="p-2">
                  <p className="text-xl font-semibold">2.5k</p>
                  <p className="text-xs text-gray-500">Stock Products</p>
                </div>
                <div className="p-2">
                  <p className="text-xl font-semibold">0.5k</p>
                  <p className="text-xs text-gray-500">Out of Stock</p>
                </div>
                <div className="p-2">
                  <p className="text-xl font-semibold">250k</p>
                  <p className="text-xs text-gray-500">Revenue</p>
                </div>
              </div>

              <WeeklyReportChart data={dashboardData.weeklyReport.chartData} />
            </div>

            <div className="space-y-6 h-[220px] ">
              <div className="w-[361px] p-[1rem] bg-white rounded-lg shadow-sm filter drop-shadow-lg">
                <div className="flex justify-between items-center mb-4">
                  <h3 className="text-sm font-medium text-primary">
                    Users in last 30 minutes
                  </h3>
                  <button className="text-gray-400">
                    <HiOutlineDotsVertical></HiOutlineDotsVertical>
                  </button>
                </div>
                <p className="text-2xl font-bold mb-2">21.5K</p>
                <p className="text-xs text-gray-500 mb-2 mt-[1rem]">
                  Users per minute
                </p>
                <UserActivityChart
                  data={dashboardData.userStats.userPerMinute}
                />
                {/* Test */}
                <div className="flex justify-between items-center mb-4 mt-[8.5rem]">
                  <h3 className="text-sm font-medium text-primary">
                    Users in last 30 minutes
                  </h3>
                  <button className="text-gray-400">
                    <HiOutlineDotsVertical></HiOutlineDotsVertical>
                  </button>
                </div>
                <p className="text-2xl font-bold mb-2">21.5K</p>
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
            <div className="filter drop-shadow-lg bg-white rounded-lg shadow-sm p-[1.25rem] w-[800px] h-[420px]">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-[18px] font-bold">Transaction</h3>
                <button className="flex items-center px-3 py-1 text-xs bg-ocean-green text-white rounded-lg w-[85px] h-[32px] p-8">
                  <span className="mr-[0.5rem] text-[14px]">Filter</span>
                  <IoFilterSharp className="w-[18px] h-[18px]"></IoFilterSharp>
                </button>
              </div>
              <TransactionTable transactions={transactions} />
              <div className="flex justify-end mt-[2rem] items-center mr-[2rem]">
                <button className="text-sm bg-white text-primary border border-primary rounded-[25px] w-[96px] h-[32px]">
                  Details
                </button>
              </div>
            </div>

            <div className="rounded-lg shadow-sm w-[292px]">
              <TopProductsTable products={topProducts} />
            </div>
          </div>

          {/* Best Selling Products and Add New Product*/}
          <div className="flex flex-wrap gap-[16px] mb-6 mt-[1.25rem]">
            <div className="filter drop-shadow-lg bg-white rounded-lg shadow-sm p-[1.25rem] w-[744px] h-[430px]">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-[18px] font-bold">Best Selling Product</h3>
                <button className="flex items-center px-3 py-1 text-xs bg-ocean-green text-white rounded-lg w-[85px] h-[32px] p-8">
                  <span className="mr-[0.5rem] text-[14px]">Filter</span>
                  <IoFilterSharp className="w-[18px] h-[18px]"></IoFilterSharp>
                </button>
              </div>
              <BestSellingTable bestProducts={bestSellingProducts} />
              <div className="flex justify-end mt-[1rem] items-center mr-[1rem]">
                <button className="text-sm bg-white text-primary border border-primary rounded-[25px] w-[96px] h-[32px]">
                  Details
                </button>
              </div>
            </div>

            <div className="bg-white rounded-lg p-[1.25rem] max-w-md w-[360px] filter drop-shadow-lg">
              <div className="flex items-center justify-between mb-6">
                <h1 className="text-[18px] font-bold text-[#23272E]">
                  Add New Product
                </h1>
                <button className="flex items-center gap-1 text-primary">
                  <span><CiCirclePlus className="w-[20px] h-[20px]"></CiCirclePlus></span>
                  <span className="text-[14px]">Add New</span>
                </button>
              </div>

              <ProductCategoryList categories={productCategories} />
              <NewProductList products={newProducts} />
            </div>
          </div>

         
        </main>
      </div>
    </div>
  );
};

export default DashboardAdmin;
