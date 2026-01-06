import { useEffect, useState } from "react";
import { PieChart } from '@mui/x-charts/PieChart';

export const valueFormatter = (item) => `${item.value}%`;

export default function PieChartSummary() {
  const [chartData, setChartData] = useState([]);
  const baseUrl = sessionStorage.getItem("baseUrl");
  const token = sessionStorage.getItem("token");

  useEffect(() => {
    fetchDataFromBackend();
  }, []);

  const fetchDataFromBackend = async () => {
    const to = Math.floor(Date.now() / 1000);
    
    const now = new Date();
    const lastMonthDate = new Date(
      now.getFullYear(),
      now.getMonth() - 1,
      now.getDate()
    );

    const from = Math.floor(lastMonthDate.getTime() / 1000);

    try {
      const res = await fetch(`${baseUrl}/api/v1/summary?from=${from}&to=${to}`, {
        headers: {
          "Authorization": `Bearer ${token}`
        }
      });

      const data = await res.json();

      // Backend MUST return something like:
      // [{ label: "Apples", value: 20 }, { label: "Bananas", value: 30 }]

      const formatted = data.ok.map(item => ({
        label: item.category,
        value: item.amount
      }));


      setChartData(formatted);

    } catch (error) {
      console.error("Error loading chart data:", error);
    }
  };

  return (
    <div style={{ width: 400, height: 400 }}>
      <PieChart
        series={[
        {
          data: chartData,
          highlightScope: { fade: 'global', highlight: 'item' },
          faded: { innerRadius: 30, additionalRadius: -30, color: 'gray' },
          valueFormatter: (item) => `${item.value} €`,
        },
        ]}
        height={300}
        width={300}
      />
    </div>
  );

}
