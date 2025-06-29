export interface CustomerSummaryCardProps {
  title: string;
  value: number;
  growthRate: number;
  period: string;
  color: 'primary' | 'success' | 'warning' | 'error';
}